package datasources

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
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

// DownloadImage downloads an image from GitHub and saves it locally
// Returns the standardized filename that should be stored in database
func (g *GitHubClient) DownloadImage(ctx context.Context, imageURL, topic string) (string, error) {
	if imageURL == "" {
		return "", nil
	}

	// Create images directory in backend static folder
	imageDir := filepath.Join("static", "quiz-images")
	if err := os.MkdirAll(imageDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create image directory: %w", err)
	}

	// Standardized filename generation: Create consistent naming
	// Format: {topic}_{original_filename}
	// Example: "android_04.jpeg", "css_Q-141.png"
	originalFilename := filepath.Base(imageURL)
	filename := fmt.Sprintf("%s_%s", topic, originalFilename)
	localFilePath := filepath.Join(imageDir, filename)

	// Check if file already exists
	if _, err := os.Stat(localFilePath); err == nil {
		// File already exists, no need to download
		fmt.Printf("Image already exists: %s\n", localFilePath)
		return filename, nil // Return standardized filename
	}

	// Construct GitHub raw URL
	// If imageURL contains a full path (like "css/images/css_q036.png"), use it directly
	// If it's just a filename, prepend the topic
	var rawURL string
	if strings.Contains(imageURL, "/") {
		// Full path provided
		rawURL = fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/%s",
			g.config.Owner, g.config.Repository, imageURL)
	} else {
		// Just filename, use topic
		rawURL = fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/%s/%s",
			g.config.Owner, g.config.Repository, topic, imageURL)
	}

	fmt.Printf("Downloading image: %s\n", rawURL)

	// Download the image
	resp, err := http.Get(rawURL)
	if err != nil {
		return "", fmt.Errorf("failed to download image from %s: %w", rawURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to download image: HTTP %d from %s", resp.StatusCode, rawURL)
	}

	// Create the local file
	out, err := os.Create(localFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create local image file: %w", err)
	}
	defer out.Close()

	// Copy the image data
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save image data: %w", err)
	}

	fmt.Printf("Successfully downloaded image: %s -> %s\n", rawURL, localFilePath)
	return filename, nil // Return standardized filename for database storage
}

// ListImageFiles gets all image files from a specific folder in the repository
func (g *GitHubClient) ListImageFiles(ctx context.Context, folderPath string) ([]string, error) {
	// URL encode the folder path to handle special characters like c++
	encodedPath := strings.ReplaceAll(folderPath, "+", "%2B")

	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s",
		g.config.BaseURL, g.config.Owner, g.config.Repository, encodedPath)

	type GitHubFile struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}

	var files []GitHubFile
	resp, err := g.client.R().
		SetContext(ctx).
		SetResult(&files).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch folder contents: %w", err)
	}

	if resp.StatusCode() == 404 {
		// Folder doesn't exist, return empty list
		return []string{}, nil
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("GitHub API returned status %d for folder %s", resp.StatusCode(), folderPath)
	}

	var imageFiles []string
	imageExtensions := map[string]bool{
		".png": true, ".jpg": true, ".jpeg": true, ".gif": true,
		".svg": true, ".webp": true, ".bmp": true, ".ico": true,
	}

	for _, file := range files {
		if file.Type == "file" {
			// Check if file has image extension
			ext := strings.ToLower(filepath.Ext(file.Name))
			if imageExtensions[ext] {
				imageFiles = append(imageFiles, file.Name)
			}
		}
	}

	return imageFiles, nil
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
		isEnglishQuizFile := isEnglishQuizFile(file.Name)
		if !isEnglishQuizFile {
			continue
		}
		if file.Name == "" || file.Content == "" {
			continue
		}
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

// FetchParsedQuestions returns ParsedQuestion data with code blocks and images separated
func (g *GitHubClient) FetchParsedQuestions(ctx context.Context, params FetchParams) ([]ParsedQuestion, error) {
	path := ""
	if params.Category != "" {
		path = params.Category
	}

	files, err := g.fetchMarkdownFiles(ctx, path, params.Since)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch markdown files: %w", err)
	}

	var questions []ParsedQuestion
	for _, file := range files {
		isEnglishQuizFile := isEnglishQuizFile(file.Name)
		if !isEnglishQuizFile {
			continue
		}
		if file.Name == "" || file.Content == "" {
			continue
		}
		quiz, err := g.parseMarkdownQuiz(file)
		if err != nil {
			continue
		}

		for _, q := range quiz.Questions {
			// Create a copy with additional metadata
			parsedQ := q
			parsedQ.Category = quiz.Category
			parsedQ.Source = "github"
			parsedQ.ExternalReference = file.URL
			parsedQ.ExternalID = file.SHA
			questions = append(questions, parsedQ)
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

// extractCodeBlock extracts code content and language from a code block node
func (g *GitHubClient) extractCodeBlock(node ast.Node, source []byte) (content string, language string) {
	switch n := node.(type) {
	case *ast.FencedCodeBlock:
		var code []byte
		for i := 0; i < n.Lines().Len(); i++ {
			line := n.Lines().At(i)
			code = append(code, line.Value(source)...)
		}
		return string(code), string(n.Language(source))
	case *ast.CodeBlock:
		var code []byte
		for i := 0; i < n.Lines().Len(); i++ {
			line := n.Lines().At(i)
			code = append(code, line.Value(source)...)
		}
		return string(code), ""
	}
	return "", ""
}

// extractImageFromText extracts image URL and alt text from a text node
func (g *GitHubClient) extractImageFromText(text string) (url string, alt string, cleanText string) {
	// Regex to match ![alt](url) pattern
	imgRegex := regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
	matches := imgRegex.FindStringSubmatch(text)

	if len(matches) == 3 {
		alt = matches[1]
		url = matches[2]
		// Remove the image markdown from text
		cleanText = imgRegex.ReplaceAllString(text, "")
		cleanText = strings.TrimSpace(cleanText)
		return url, alt, cleanText
	}

	return "", "", text
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
	// Match Q1, Q1., Q 1., q1, q1., Q 1, etc., with optional leading spaces
	qNumRegex := regexp.MustCompile(`(?i)^\s*q\s*\d+\.?\s*`)

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
					qText := extractTextFromNode(node, string(source))
					qText = strings.TrimSpace(qText)
					// Remove leading Qn. from question text (robust)
					qText = qNumRegex.ReplaceAllString(qText, "")

					// Extract images from question text
					imageURL, imageAlt, cleanText := g.extractImageFromText(qText)

					currentQuestion = &ParsedQuestion{
						Text: cleanText,
						Type: "multiple_choice",
					}

					// Set image fields if found
					if imageURL != "" {
						currentQuestion.QuestionImageURL = &imageURL
						currentQuestion.QuestionImageAlt = &imageAlt

						// Download the image to local storage and get standardized filename
						if standardizedFilename, err := g.DownloadImage(context.Background(), imageURL, quiz.Category); err != nil {
							fmt.Printf("Warning: Failed to download question image %s: %v\n", imageURL, err)
						} else if standardizedFilename != "" {
							// Update with standardized filename for database storage
							currentQuestion.QuestionImageURL = &standardizedFilename
						}
					}

					// Debug Q36 specifically
					if strings.Contains(cleanText, "Which code snippet would achieve the layout displayed below") {
						fmt.Printf("🔍 Q36 DEBUG - Question text: '%s'\n", cleanText)
						fmt.Printf("🔍 Q36 DEBUG - Original qText: '%s'\n", qText)
						fmt.Printf("🔍 Q36 DEBUG - Image URL: '%s'\n", imageURL)
						fmt.Printf("🔍 Q36 DEBUG - Image Alt: '%s'\n", imageAlt)
					}
				}
			case *ast.Image:
				// Handle direct image nodes
				if currentQuestion != nil && len(currentQuestion.Choices) == 0 {
					if strings.Contains(currentQuestion.Text, "Which code snippet would achieve the layout displayed below") {
						fmt.Printf("🖼️ Q36 DEBUG - Found Image Node!\n")

						// Extract image URL and alt text
						imageURL := string(node.Destination)
						imageAlt := ""

						// Get alt text from image children
						for child := node.FirstChild(); child != nil; child = child.NextSibling() {
							if textNode, ok := child.(*ast.Text); ok {
								imageAlt = string(textNode.Segment.Value([]byte(source)))
								break
							}
						}

						fmt.Printf("🖼️ Q36 DEBUG - Image URL: '%s'\n", imageURL)
						fmt.Printf("🖼️ Q36 DEBUG - Image Alt: '%s'\n", imageAlt)

						// Only set if we don't already have an image
						if currentQuestion.QuestionImageURL == nil {
							// Download the image to local storage and get standardized filename
							if standardizedFilename, err := g.DownloadImage(context.Background(), imageURL, quiz.Category); err != nil {
								fmt.Printf("Warning: Failed to download question image %s: %v\n", imageURL, err)
								// Fallback to original URL
								currentQuestion.QuestionImageURL = &imageURL
							} else if standardizedFilename != "" {
								// Use standardized filename for database storage
								currentQuestion.QuestionImageURL = &standardizedFilename
							}

							if imageAlt != "" {
								currentQuestion.QuestionImageAlt = &imageAlt
							}
						}
					}
				}
			case *ast.FencedCodeBlock, *ast.CodeBlock:
				// Capture code blocks that appear
				if currentQuestion != nil && len(currentQuestion.Choices) == 0 {
					content, lang := g.extractCodeBlock(node, []byte(source))
					if content != "" {
						// Store the first code block found for the question
						if currentQuestion.QuestionCode == nil {
							currentQuestion.QuestionCode = &content
							if lang != "" {
								currentQuestion.QuestionCodeLang = &lang
							}
						}
					}
				}
			case *ast.Paragraph:
				// Check if this paragraph contains an image and we have a current question but no choices yet
				if currentQuestion != nil && len(currentQuestion.Choices) == 0 {
					paragraphText := g.extractFullTextFromNode(node, string(source))
					paragraphText = strings.TrimSpace(paragraphText)

					// Debug logging for Q36
					if strings.Contains(currentQuestion.Text, "layout displayed below") {
						fmt.Printf("DEBUG Q36: Found paragraph: '%s'\n", paragraphText)
						fmt.Printf("DEBUG Q36: Paragraph length: %d chars\n", len(paragraphText))
						// Test if it matches image pattern
						imgRegex := regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
						if imgRegex.MatchString(paragraphText) {
							fmt.Printf("DEBUG Q36: Paragraph matches image pattern!\n")
						} else {
							fmt.Printf("DEBUG Q36: Paragraph does NOT match image pattern\n")
						}
					}

					// Check if this paragraph contains an image
					if imageURL, imageAlt, cleanText := g.extractImageFromText(paragraphText); imageURL != "" {
						// Debug logging for Q36
						if strings.Contains(currentQuestion.Text, "layout displayed below") {
							fmt.Printf("DEBUG Q36: Found image in paragraph - URL: %s, Alt: %s\n", imageURL, imageAlt)
						}

						// Only set if we don't already have an image from the question text
						if currentQuestion.QuestionImageURL == nil {
							// Download the image to local storage and get standardized filename
							if standardizedFilename, err := g.DownloadImage(context.Background(), imageURL, quiz.Category); err != nil {
								fmt.Printf("Warning: Failed to download question image %s: %v\n", imageURL, err)
								// Fallback to original URL
								currentQuestion.QuestionImageURL = &imageURL
							} else if standardizedFilename != "" {
								// Use standardized filename for database storage
								currentQuestion.QuestionImageURL = &standardizedFilename
							}
							currentQuestion.QuestionImageAlt = &imageAlt
						}

						// If there's remaining text after removing the image, append it to the question
						if cleanText != "" {
							if currentQuestion.Text != "" {
								currentQuestion.Text += " " + cleanText
							} else {
								currentQuestion.Text = cleanText
							}
						}
					}
				}
			case *ast.List:
				if currentQuestion != nil && len(currentQuestion.Choices) == 0 {
					// Only process the first list for each question
					choices, correct, choiceCodes, choiceCodeLangs, choiceImageURLs, choiceImageAlts := g.extractChoicesFromList(node, string(source), quiz.Category)

					currentQuestion.Choices = choices
					currentQuestion.Correct = correct
					currentQuestion.ChoiceCodes = choiceCodes
					currentQuestion.ChoiceCodeLangs = choiceCodeLangs
					currentQuestion.ChoiceImageURLs = choiceImageURLs
					currentQuestion.ChoiceImageAlts = choiceImageAlts
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

func (g *GitHubClient) extractChoicesFromList(list *ast.List, source string, category string) ([]string, int, []*string, []*string, []*string, []*string) {
	sourceBytes := []byte(source)
	var choices []string
	var choiceCodes []*string
	var choiceCodeLangs []*string
	var choiceImageURLs []*string
	var choiceImageAlts []*string
	correctIndex := -1

	for child := list.FirstChild(); child != nil; child = child.NextSibling() {
		if listItem, ok := child.(*ast.ListItem); ok {
			var textParts []string
			var choiceCode *string
			var choiceCodeLang *string
			var choiceImageURL *string
			var choiceImageAlt *string

			// Check if this is Q36 by looking for the layout text in the source
			isQ36Debug := strings.Contains(source, "layout displayed below")

			for liChild := listItem.FirstChild(); liChild != nil; liChild = liChild.NextSibling() {
				switch n := liChild.(type) {
				case *ast.FencedCodeBlock:
					content, lang := g.extractCodeBlock(n, sourceBytes)
					if content != "" && choiceCode == nil { // Take first code block only
						choiceCode = &content
						if lang != "" {
							choiceCodeLang = &lang
						}
					}

				case *ast.CodeBlock:
					content, _ := g.extractCodeBlock(n, sourceBytes)
					if content != "" && choiceCode == nil { // Take first code block only
						choiceCode = &content
					}

				default:
					txt := g.extractFullTextFromNode(n, source)
					txt = strings.ReplaceAll(txt, "&shy;", "")
					txt = strings.TrimSpace(txt)
					if txt != "" {
						fmt.Printf("DEBUG: Found text part: '%s'\n", txt)
						textParts = append(textParts, txt)
					}
				}
			}

			// Join all text parts and process
			itemText := strings.Join(textParts, "\n")

			if isQ36Debug {
				fmt.Printf("🔍 Q36 DEBUG - Combined item text: '%s'\n", itemText)
			}

			// Extract images from choice text
			if imageURL, imageAlt, cleanText := g.extractImageFromText(itemText); imageURL != "" {
				choiceImageURL = &imageURL
				choiceImageAlt = &imageAlt
				itemText = cleanText

				// Download the choice image to local storage and get standardized filename
				if standardizedFilename, err := g.DownloadImage(context.Background(), imageURL, category); err != nil {
					fmt.Printf("Warning: Failed to download choice image %s: %v\n", imageURL, err)
					// Keep original URL as fallback
				} else if standardizedFilename != "" {
					// Update with standardized filename for database storage
					choiceImageURL = &standardizedFilename
				}
			}

			// Checkbox detection
			if isQ36Debug {
				fmt.Printf("🔍 Q36 DEBUG - Before checkbox detection: '%s'\n", itemText)
			}
			if strings.Contains(itemText, "[x]") || strings.Contains(itemText, "[X]") {
				if isQ36Debug {
					fmt.Printf("🔍 Q36 DEBUG - Found correct answer (x)\n")
				}
				itemText = strings.ReplaceAll(itemText, "[x]", "")
				itemText = strings.ReplaceAll(itemText, "[X]", "")
				correctIndex = len(choices)
			}
			itemText = strings.ReplaceAll(itemText, "[ ]", "")
			itemText = strings.TrimSpace(itemText)

			if isQ36Debug {
				fmt.Printf("🔍 Q36 DEBUG - After checkbox cleaning: '%s'\n", itemText)
			}

			if itemText != "" {
				if isQ36Debug {
					fmt.Printf("🔍 Q36 DEBUG - Adding choice #%d: '%s'\n", len(choices), itemText)
				}
				choices = append(choices, itemText)
				choiceCodes = append(choiceCodes, choiceCode)
				choiceCodeLangs = append(choiceCodeLangs, choiceCodeLang)
				choiceImageURLs = append(choiceImageURLs, choiceImageURL)
				choiceImageAlts = append(choiceImageAlts, choiceImageAlt)
			} else {
				if isQ36Debug {
					fmt.Printf("🔍 Q36 DEBUG - Skipping empty choice\n")
				}
			}
		}
	}

	// Now consider sibling nodes: sometimes the markdown places details (fenced code
	// blocks) or additional list nodes after the initial list. Walk siblings until
	// the next heading and merge any fenced code block content into the previous
	// choice, and extract additional list items from sibling lists.
	for sib := list.NextSibling(); sib != nil; sib = sib.NextSibling() {
		// Stop when we reach a heading - that's the next question.
		if _, ok := sib.(*ast.Heading); ok {
			break
		}

		switch s := sib.(type) {
		case *ast.FencedCodeBlock:
			// Attach code block content to the last choice (if present)
			if len(choices) == 0 {
				continue
			}
			lastIndex := len(choices) - 1
			content, lang := g.extractCodeBlock(s, sourceBytes)
			if content != "" {
				// If the last choice doesn't have a code block yet, assign it
				if choiceCodes[lastIndex] == nil {
					choiceCodes[lastIndex] = &content
					if lang != "" {
						choiceCodeLangs[lastIndex] = &lang
					}
				}
			}

		case *ast.List:
			// Extract items from this list and append
			for li := s.FirstChild(); li != nil; li = li.NextSibling() {
				if listItem, ok := li.(*ast.ListItem); ok {
					var textParts []string
					var choiceCode *string
					var choiceCodeLang *string
					var choiceImageURL *string
					var choiceImageAlt *string

					for liChild := listItem.FirstChild(); liChild != nil; liChild = liChild.NextSibling() {
						switch n := liChild.(type) {
						case *ast.FencedCodeBlock:
							content, lang := g.extractCodeBlock(n, sourceBytes)
							if content != "" && choiceCode == nil {
								choiceCode = &content
								if lang != "" {
									choiceCodeLang = &lang
								}
							}
						case *ast.CodeBlock:
							content, _ := g.extractCodeBlock(n, sourceBytes)
							if content != "" && choiceCode == nil {
								choiceCode = &content
							}
						default:
							txt := g.extractFullTextFromNode(n, source)
							txt = strings.ReplaceAll(txt, "&shy;", "")
							txt = strings.TrimSpace(txt)
							if txt != "" {
								textParts = append(textParts, txt)
							}
						}
					}

					itemText := strings.Join(textParts, "\n")

					// Extract images from choice text
					if imageURL, imageAlt, cleanText := g.extractImageFromText(itemText); imageURL != "" {
						choiceImageURL = &imageURL
						choiceImageAlt = &imageAlt
						itemText = cleanText

						// Download the choice image to local storage and get standardized filename
						if standardizedFilename, err := g.DownloadImage(context.Background(), imageURL, category); err != nil {
							fmt.Printf("Warning: Failed to download choice image %s: %v\n", imageURL, err)
							// Keep original URL as fallback
						} else if standardizedFilename != "" {
							// Update with standardized filename for database storage
							choiceImageURL = &standardizedFilename
						}
					}

					// Checkbox detection
					if strings.Contains(itemText, "[x]") || strings.Contains(itemText, "[X]") {
						itemText = strings.ReplaceAll(itemText, "[x]", "")
						itemText = strings.ReplaceAll(itemText, "[X]", "")
						correctIndex = len(choices)
					}
					itemText = strings.ReplaceAll(itemText, "[ ]", "")
					itemText = strings.TrimSpace(itemText)
					if itemText != "" {
						choices = append(choices, itemText)
						choiceCodes = append(choiceCodes, choiceCode)
						choiceCodeLangs = append(choiceCodeLangs, choiceCodeLang)
						choiceImageURLs = append(choiceImageURLs, choiceImageURL)
						choiceImageAlts = append(choiceImageAlts, choiceImageAlt)
					}
				}
			}
		default:
			// ignore other node types
		}
	}

	// Add final debug output for Q36
	isQ36 := strings.Contains(source, "layout displayed below")
	if isQ36 {
		fmt.Printf("🔍 Q36 DEBUG - FINAL RESULTS:\n")
		fmt.Printf("🔍 Q36 DEBUG - Choices count: %d\n", len(choices))
		for i, choice := range choices {
			fmt.Printf("🔍 Q36 DEBUG - Choice %d: '%s'\n", i, choice)
		}
		fmt.Printf("🔍 Q36 DEBUG - Correct index: %d\n", correctIndex)
	}

	return choices, correctIndex, choiceCodes, choiceCodeLangs, choiceImageURLs, choiceImageAlts
}

// extractFullTextFromNode extracts only plain text content from a node, excluding code blocks and images
func (g *GitHubClient) extractFullTextFromNode(node ast.Node, source string) string {
	var out strings.Builder
	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch v := n.(type) {
		case *ast.Text:
			out.WriteString(string(v.Segment.Value([]byte(source))))
		case *ast.CodeSpan:
			// Include inline code spans as they're different from code blocks
			if lines := v.Lines(); lines != nil {
				out.WriteString("`" + string(lines.Value([]byte(source))) + "`")
			}
		// Skip code blocks and images - they're handled separately now
		case *ast.FencedCodeBlock, *ast.CodeBlock:
			return ast.WalkSkipChildren, nil
		case *ast.Image:
			return ast.WalkSkipChildren, nil
		}
		return ast.WalkContinue, nil
	})
	return strings.TrimSpace(out.String())
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

func isEnglishQuizFile(quizIdentifier string) bool {
	// Matches _fr, -fr, _de, -de, _ua, -ua, etc. at the end of the identifier (before .md or end of string)
	nonEnglishSuffix := regexp.MustCompile(`(?i)[_-](fr|de|ua|es|ru|it|pl|pt|tr|zh|ja|ko|ar|hi|bn|vi|fa|he|nl|sv|no|da|fi|el|cs|ro|hu|th|id|ms|sk|sl|hr|lt|lv|et|bg|sr|mk|ca|eu|gl|af|sq|sw|zu|xh|st|tn|ts|ss|ve|nr|rw|so|am|om|ti|yo|ig|ha|ee|tw|lg|ln|kg|lu|ny|sn|sh|bs|mt|ga|cy|gd|kw|gv|sm|to|fj|mi|ty|rn|kg|rw|ss|tn|ts|ve|xh|yo|zu)(\.|$)`)
	return !nonEnglishSuffix.MatchString(quizIdentifier)
}
