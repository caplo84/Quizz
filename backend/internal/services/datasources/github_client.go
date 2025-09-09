package datasources

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// min returns the smaller of x or y
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

type GitHubClient struct {
	config GitHubConfig
	client *resty.Client
	parser goldmark.Markdown
}

func NewGitHubClient(config GitHubConfig) *GitHubClient {
	client := resty.New()
	client.SetHeader("Authorization", "token "+config.Token)
	client.SetHeader("Accept", "application/vnd.github.v3+json")
	client.SetTimeout(30 * time.Second)
	client.SetRetryCount(3)

	if config.BaseURL == "" {
		config.BaseURL = "https://api.github.com"
	}
	config.BaseURL = strings.TrimRight(config.BaseURL, "/")

	return &GitHubClient{
		config: config,
		client: client,
		parser: goldmark.New(),
	}
}

func (g *GitHubClient) GetMetadata() SourceMetadata {
	return SourceMetadata{
		Name:         "GitHub",
		Version:      "v3",
		RateLimited:  5000,
		Capabilities: []string{"markdown_parsing", "incremental_sync", "commit_tracking"},
	}
}

func (g *GitHubClient) ValidateAuth(ctx context.Context) error {
	url := fmt.Sprintf("%s/user", g.config.BaseURL)
	resp, err := g.client.R().
		SetContext(ctx).
		Get(url)
	if err != nil {
		return fmt.Errorf("failed to validate GitHub auth: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("GitHub auth failed: %d", resp.StatusCode())
	}

	return nil
}

func (g *GitHubClient) GetRateLimit(ctx context.Context) (*RateLimit, error) {
	url := fmt.Sprintf("%s/rate_limit", g.config.BaseURL)
	resp, err := g.client.R().
		SetContext(ctx).
		Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get rate limit: %w", err)
	}

	var rateLimitResp struct {
		Resources struct {
			Core struct {
				Remaining int   `json:"remaining"`
				Reset     int64 `json:"reset"`
				Limit     int   `json:"limit"`
			} `json:"core"`
		} `json:"resources"`
	}
	if err := json.Unmarshal(resp.Body(), &rateLimitResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rate limit response: %w", err)
	}

	return &RateLimit{
		Remaining: rateLimitResp.Resources.Core.Remaining,
		Reset:     time.Unix(rateLimitResp.Resources.Core.Reset, 0),
		Total:     rateLimitResp.Resources.Core.Limit,
	}, nil
}

func (g *GitHubClient) FetchCategories(ctx context.Context) ([]Category, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/contents",
		g.config.BaseURL, g.config.Owner, g.config.Repository)

	resp, err := g.client.R().
		SetContext(ctx).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	// Try to unmarshal as an array first
	var files []GitHubFile
	if err := json.Unmarshal(resp.Body(), &files); err != nil {
		// If array fails, try to unmarshal as a single object (in case it's a single file response)
		var singleFile GitHubFile
		if err2 := json.Unmarshal(resp.Body(), &singleFile); err2 != nil {
			return nil, fmt.Errorf("failed to parse GitHub response as array (%v) or object (%v). Response: %s",
				err, err2, string(resp.Body()[:min(500, len(resp.Body()))]))
		}
		files = []GitHubFile{singleFile}
	}

	var categories []Category
	for _, file := range files {
		if file.Type == "dir" {
			categories = append(categories, Category{
				Name: file.Name,
			})
		}
	}

	return categories, nil
}

func (g *GitHubClient) FetchQuestions(ctx context.Context, params FetchParams) ([]Question, error) {
	path := ""
	if params.Category != "" {
		path = params.Category
	}

	files, err := g.fetchMarkdownFiles(ctx, path, params.Since)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch markdown files: %w", err)
	}

	var questions []Question
	for _, file := range files {
		quiz, err := g.parseMarkdownQuiz(file)
		if err != nil {
			continue
		}

		for _, q := range quiz.Questions {
			questions = append(questions, Question{
				Text:              q.Text,
				Choices:           q.Choices,
				Correct:           q.Correct,
				Category:          quiz.Category,
				Source:            "github",
				ExternalReference: file.URL,
				ExternalID:        file.SHA,
			})
		}
	}

	return questions, nil
}

func (g *GitHubClient) fetchMarkdownFiles(ctx context.Context, path string, since *time.Time) ([]GitHubFile, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s",
		g.config.BaseURL, g.config.Owner, g.config.Repository, path)

	resp, err := g.client.R().
		SetContext(ctx).
		Get(url)

	if err != nil {
		return nil, err
	}

	var files []GitHubFile
	if err := json.Unmarshal(resp.Body(), &files); err != nil {
		return nil, err
	}

	var markdownFiles []GitHubFile
	for _, file := range files {
		if strings.HasSuffix(file.Name, ".md") {
			content, err := g.fetchFileContent(ctx, file.URL)
			if err != nil {
				continue
			}
			file.Content = content
			markdownFiles = append(markdownFiles, file)
		}
	}

	return markdownFiles, nil
}

func (g *GitHubClient) fetchFileContent(ctx context.Context, url string) (string, error) {
	resp, err := g.client.R().
		SetContext(ctx).
		Get(url)

	if err != nil {
		return "", err
	}

	var file GitHubFile
	if err := json.Unmarshal(resp.Body(), &file); err != nil {
		return "", err
	}

	decoded, err := base64.StdEncoding.DecodeString(file.Content)
	if err != nil {
		return "", err
	}

	return string(decoded), nil
}

func (g *GitHubClient) parseMarkdownQuiz(file GitHubFile) (*ParsedQuiz, error) {
	source := []byte(file.Content)
	doc := g.parser.Parser().Parse(text.NewReader(source))

	// Clean the filename by removing URL parameters and decoding
	cleanFileName := file.Name
	if idx := strings.Index(cleanFileName, "?"); idx != -1 {
		cleanFileName = cleanFileName[:idx]
	}
	// Also remove any URL encoding
	cleanFileName = strings.ReplaceAll(cleanFileName, "%20", " ")

	quiz := &ParsedQuiz{
		Source:     "github",
		ExternalID: file.SHA,
		Category:   strings.TrimSuffix(cleanFileName, ".md"),
	}

	var currentQuestion *ParsedQuestion

	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			switch node := n.(type) {
			case *ast.Heading:
				if node.Level == 1 {
					// Extract title from markdown heading
					quiz.Title = extractTextFromNode(node, string(source))
				} else if node.Level == 4 { // #### Q1. Question text
					if currentQuestion != nil {
						quiz.Questions = append(quiz.Questions, *currentQuestion)
					}
					currentQuestion = &ParsedQuestion{
						Text: extractTextFromNode(node, string(source)),
						Type: "multiple_choice",
					}
				}
			case *ast.List:
				if currentQuestion != nil {
					choices, correct := g.extractChoicesFromList(node, string(source))
					currentQuestion.Choices = choices
					currentQuestion.Correct = correct
				}
			}
		}
		return ast.WalkContinue, nil
	})

	// Add the last question
	if currentQuestion != nil {
		quiz.Questions = append(quiz.Questions, *currentQuestion)
	}

	return quiz, nil
}

func (g *GitHubClient) extractChoicesFromList(list *ast.List, source string) ([]string, int) {
	var choices []string
	correctIndex := -1

	for child := list.FirstChild(); child != nil; child = child.NextSibling() {
		if listItem, ok := child.(*ast.ListItem); ok {
			// Get the full text content of the list item
			text := g.extractFullTextFromNode(listItem, source)

			// Check for correct answer markers: [x] or [X]
			if strings.Contains(text, "[x]") || strings.Contains(text, "[X]") {
				correctIndex = len(choices)
				// Remove the checkbox markers
				text = strings.ReplaceAll(text, "[x]", "")
				text = strings.ReplaceAll(text, "[X]", "")
			} else {
				// Remove the empty checkbox markers
				text = strings.ReplaceAll(text, "[ ]", "")
			}

			// Clean up the text
			text = strings.TrimSpace(text)

			// Skip empty choices
			if text != "" {
				choices = append(choices, text)
			}
		}
	}

	return choices, correctIndex
}

// extractFullTextFromNode extracts all text content from a node, including nested elements
func (g *GitHubClient) extractFullTextFromNode(node ast.Node, source string) string {
	var text strings.Builder

	// Walk through all child nodes to extract text
	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			switch child := n.(type) {
			case *ast.Text:
				segment := child.Segment
				text.Write(segment.Value([]byte(source)))
			case *ast.CodeSpan:
				// Handle inline code
				for c := child.FirstChild(); c != nil; c = c.NextSibling() {
					if textNode, ok := c.(*ast.Text); ok {
						segment := textNode.Segment
						text.WriteString("`")
						text.Write(segment.Value([]byte(source)))
						text.WriteString("`")
					}
				}
			}
		}
		return ast.WalkContinue, nil
	})

	return strings.TrimSpace(text.String())
}

func extractTextFromNode(node ast.Node, source string) string {
	var text strings.Builder

	// Walk through all children to get complete text
	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if textNode, ok := n.(*ast.Text); ok {
				segment := textNode.Segment
				text.Write(segment.Value([]byte(source)))
			}
		}
		return ast.WalkContinue, nil
	})

	return strings.TrimSpace(text.String())
}

func (g *GitHubClient) GetLastCommit(ctx context.Context, path string) (*GitHubCommit, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/commits",
		g.config.BaseURL, g.config.Owner, g.config.Repository)

	resp, err := g.client.R().
		SetContext(ctx).
		SetQueryParam("path", path).
		SetQueryParam("per_page", "1").
		Get(url)

	if err != nil {
		return nil, err
	}

	var commits []GitHubCommit
	if err := json.Unmarshal(resp.Body(), &commits); err != nil {
		return nil, err
	}

	if len(commits) == 0 {
		return nil, fmt.Errorf("no commits found")
	}

	return &commits[0], nil
}
