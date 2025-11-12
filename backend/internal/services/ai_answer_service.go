package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type ExtractedAnswer struct {
	CorrectAnswer string
	Found         bool
	Confidence    float64
	SourceExcerpt string
}

type GeneratedAnswerSet struct {
	CorrectAnswer string
	WrongAnswers  []string
	Explanation   string
	Confidence    float64
}

type QuestionQualityReview struct {
	Severity   string
	Comment    string
	Issues     []string
	Confidence float64
}

type AIAnswerService struct {
	client    *resty.Client
	provider  string
	baseURL   string
	apiKey    string
	model     string
	maxTokens int
}

func NewAIAnswerServiceFromEnv() *AIAnswerService {
	provider := strings.ToLower(strings.TrimSpace(os.Getenv("AI_PROVIDER")))
	if provider == "" {
		provider = "ollama"
	}

	apiKey := ""
	model := ""
	baseURL := ""
	accountID := ""

	switch provider {
	case "ollama":
		baseURL = strings.TrimSpace(os.Getenv("OLLAMA_BASE_URL"))
		if baseURL == "" {
			baseURL = "http://localhost:11434"
		}
		model = strings.TrimSpace(os.Getenv("OLLAMA_MODEL"))
		if model == "" {
			model = "qwen2.5:7b"
		}
	case "cloudflare":
		apiKey = strings.TrimSpace(os.Getenv("CLOUDFLARE_API_TOKEN"))
		if apiKey == "" {
			apiKey = strings.TrimSpace(os.Getenv("CF_API_TOKEN"))
		}
		accountID = strings.TrimSpace(os.Getenv("CLOUDFLARE_ACCOUNT_ID"))
		if accountID == "" {
			accountID = strings.TrimSpace(os.Getenv("CF_ACCOUNT_ID"))
		}
		baseURL = normalizeCloudflareBaseURL(strings.TrimSpace(os.Getenv("CLOUDFLARE_AI_BASE_URL")), accountID)
		model = normalizeCloudflareModel(strings.TrimSpace(os.Getenv("CLOUDFLARE_AI_MODEL")))
		if model == "" {
			model = "@cf/meta/llama-3.1-8b-instruct"
		}
	default:
		// Fallback to local ollama for unknown provider values
		provider = "ollama"
		baseURL = strings.TrimSpace(os.Getenv("OLLAMA_BASE_URL"))
		if baseURL == "" {
			baseURL = "http://localhost:11434"
		}
		model = strings.TrimSpace(os.Getenv("OLLAMA_MODEL"))
		if model == "" {
			model = "qwen2.5:7b"
		}
	}

	client := resty.New().
		SetBaseURL(strings.TrimRight(baseURL, "/")).
		SetTimeout(60 * time.Second).
		SetRetryCount(2)

	return &AIAnswerService{
		client:    client,
		provider:  provider,
		baseURL:   strings.TrimRight(baseURL, "/"),
		apiKey:    apiKey,
		model:     model,
		maxTokens: 700,
	}
}

func normalizeCloudflareBaseURL(raw, accountID string) string {
	baseURL := strings.TrimSpace(raw)
	if baseURL == "" {
		if accountID == "" {
			return ""
		}
		return fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/ai", accountID)
	}

	baseURL = strings.TrimRight(baseURL, "/")
	if idx := strings.Index(baseURL, "/run/"); idx >= 0 {
		baseURL = baseURL[:idx]
	}
	baseURL = strings.TrimSuffix(baseURL, "/run")

	if regexp.MustCompile(`/accounts/[^/]+$`).MatchString(baseURL) {
		baseURL += "/ai"
	}

	if idx := strings.Index(baseURL, "/ai/"); idx >= 0 {
		baseURL = baseURL[:idx+3]
	}

	return strings.TrimRight(baseURL, "/")
}

func normalizeCloudflareModel(raw string) string {
	model := strings.TrimSpace(raw)
	if model == "" {
		return ""
	}

	if idx := strings.Index(model, "/run/"); idx >= 0 {
		model = model[idx+len("/run/"):]
	}
	model = strings.TrimPrefix(model, "/run/")

	if idx := strings.Index(model, "?"); idx >= 0 {
		model = model[:idx]
	}

	return strings.TrimSpace(strings.Trim(model, "/"))
}

func (s *AIAnswerService) IsConfigured() bool {
	if s == nil {
		return false
	}

	switch s.provider {
	case "ollama":
		return s.model != ""
	case "cloudflare":
		return s.apiKey != "" && s.baseURL != "" && s.model != ""
	default:
		return false
	}
}

func (s *AIAnswerService) ExtractAnswerFromMarkdown(ctx context.Context, question, markdown string) (*ExtractedAnswer, error) {
	if !s.IsConfigured() {
		return nil, fmt.Errorf("ai service is not configured for provider: %s", s.provider)
	}

	prompt := fmt.Sprintf(`You are validating quiz data.
Given the markdown content below, find the correct answer for the exact question.

Question:
%s

Markdown:
%s

Return ONLY valid JSON:
{
  "correct_answer": "string",
  "found": true,
  "confidence": 0.0,
  "source_excerpt": "string"
}`,
		question,
		markdown,
	)

	type responsePayload struct {
		CorrectAnswer string  `json:"correct_answer"`
		Found         bool    `json:"found"`
		Confidence    float64 `json:"confidence"`
		SourceExcerpt string  `json:"source_excerpt"`
	}

	var payload responsePayload
	if err := s.runJSONPrompt(ctx, prompt, &payload); err != nil {
		return nil, err
	}

	payload.CorrectAnswer = strings.TrimSpace(payload.CorrectAnswer)
	if payload.Confidence < 0 {
		payload.Confidence = 0
	}
	if payload.Confidence > 1 {
		payload.Confidence = 1
	}

	return &ExtractedAnswer{
		CorrectAnswer: payload.CorrectAnswer,
		Found:         payload.Found && payload.CorrectAnswer != "",
		Confidence:    payload.Confidence,
		SourceExcerpt: strings.TrimSpace(payload.SourceExcerpt),
	}, nil
}

func (s *AIAnswerService) GenerateWrongAnswers(ctx context.Context, correctAnswer, question, topic string) ([]string, error) {
	if !s.IsConfigured() {
		return nil, fmt.Errorf("ai service is not configured for provider: %s", s.provider)
	}

	prompt := fmt.Sprintf(`You generate plausible wrong multiple-choice answers.
Question: %s
Topic: %s
Correct answer: %s

Generate exactly 3 plausible but incorrect answers.
Return ONLY valid JSON:
{
  "wrong_answers": ["a", "b", "c"]
}`,
		question,
		topic,
		correctAnswer,
	)

	type responsePayload struct {
		WrongAnswers []string `json:"wrong_answers"`
	}

	var payload responsePayload
	if err := s.runJSONPrompt(ctx, prompt, &payload); err != nil {
		return nil, err
	}

	cleaned := make([]string, 0, len(payload.WrongAnswers))
	seen := map[string]struct{}{}
	for _, item := range payload.WrongAnswers {
		value := strings.TrimSpace(item)
		if value == "" {
			continue
		}
		norm := normalizeQA(value)
		if _, ok := seen[norm]; ok {
			continue
		}
		seen[norm] = struct{}{}
		cleaned = append(cleaned, value)
	}

	if len(cleaned) < 3 {
		return nil, fmt.Errorf("ai returned insufficient wrong answers: got %d", len(cleaned))
	}

	return cleaned[:3], nil
}

func (s *AIAnswerService) GenerateAnswerSet(ctx context.Context, question, topic, difficulty string) (*GeneratedAnswerSet, error) {
	if !s.IsConfigured() {
		return nil, fmt.Errorf("ai service is not configured for provider: %s", s.provider)
	}

	prompt := fmt.Sprintf(`You create high-quality quiz choices.
Question: %s
Topic: %s
Difficulty: %s

Return ONLY valid JSON with one correct answer and exactly 3 wrong answers:
{
  "correct_answer": "string",
  "wrong_answers": ["a", "b", "c"],
  "explanation": "short explanation",
  "confidence": 0.0
}`,
		question,
		topic,
		difficulty,
	)

	type responsePayload struct {
		CorrectAnswer string   `json:"correct_answer"`
		WrongAnswers  []string `json:"wrong_answers"`
		Explanation   string   `json:"explanation"`
		Confidence    float64  `json:"confidence"`
	}

	var payload responsePayload
	if err := s.runJSONPrompt(ctx, prompt, &payload); err != nil {
		return nil, err
	}

	payload.CorrectAnswer = strings.TrimSpace(payload.CorrectAnswer)
	if payload.CorrectAnswer == "" {
		return nil, fmt.Errorf("ai returned empty correct answer")
	}

	wrong, err := s.GenerateWrongAnswers(ctx, payload.CorrectAnswer, question, topic)
	if err == nil {
		payload.WrongAnswers = wrong
	}

	cleanWrong := make([]string, 0, 3)
	seen := map[string]struct{}{normalizeQA(payload.CorrectAnswer): {}}
	for _, item := range payload.WrongAnswers {
		value := strings.TrimSpace(item)
		if value == "" {
			continue
		}
		norm := normalizeQA(value)
		if _, ok := seen[norm]; ok {
			continue
		}
		seen[norm] = struct{}{}
		cleanWrong = append(cleanWrong, value)
	}
	if len(cleanWrong) < 3 {
		return nil, fmt.Errorf("ai returned insufficient unique wrong answers: got %d", len(cleanWrong))
	}

	if payload.Confidence < 0 {
		payload.Confidence = 0
	}
	if payload.Confidence > 1 {
		payload.Confidence = 1
	}
	if payload.Confidence == 0 {
		payload.Confidence = 0.7
	}

	return &GeneratedAnswerSet{
		CorrectAnswer: payload.CorrectAnswer,
		WrongAnswers:  cleanWrong[:3],
		Explanation:   strings.TrimSpace(payload.Explanation),
		Confidence:    payload.Confidence,
	}, nil
}

func (s *AIAnswerService) ReviewQuestionQuality(ctx context.Context, question, topic, difficulty string, choices []string, correctAnswer string) (*QuestionQualityReview, error) {
	if !s.IsConfigured() {
		return nil, fmt.Errorf("ai service is not configured for provider: %s", s.provider)
	}

	choiceLines := make([]string, 0, len(choices))
	for i, c := range choices {
		choiceLines = append(choiceLines, fmt.Sprintf("%d) %s", i+1, strings.TrimSpace(c)))
	}

	prompt := fmt.Sprintf(`You are a strict quiz reviewer.
Review the content quality of a single multiple-choice question and its answers.

Question: %s
Topic: %s
Difficulty: %s
Choices:
%s
Marked correct answer: %s

Rules:
- Evaluate factual correctness, ambiguity, missing context, duplicate options, and whether marked correct answer looks valid.
- Do NOT rewrite content.
- Return severity using one of: good, warning, danger.
- Keep comment concise.
- issues should be a short list of concrete findings.

Return ONLY valid JSON:
{
  "severity": "good",
  "comment": "string",
  "issues": ["string"],
  "confidence": 0.0
}`,
		question,
		topic,
		difficulty,
		strings.Join(choiceLines, "\n"),
		correctAnswer,
	)

	type responsePayload struct {
		Severity   string   `json:"severity"`
		Comment    string   `json:"comment"`
		Issues     []string `json:"issues"`
		Confidence float64  `json:"confidence"`
	}

	var payload responsePayload
	if err := s.runJSONPrompt(ctx, prompt, &payload); err != nil {
		return nil, err
	}

	severity := strings.ToLower(strings.TrimSpace(payload.Severity))
	if severity != "good" && severity != "warning" && severity != "danger" {
		severity = "warning"
	}

	if payload.Confidence < 0 {
		payload.Confidence = 0
	}
	if payload.Confidence > 1 {
		payload.Confidence = 1
	}

	issues := make([]string, 0, len(payload.Issues))
	for _, item := range payload.Issues {
		v := strings.TrimSpace(item)
		if v != "" {
			issues = append(issues, v)
		}
	}

	return &QuestionQualityReview{
		Severity:   severity,
		Comment:    strings.TrimSpace(payload.Comment),
		Issues:     issues,
		Confidence: payload.Confidence,
	}, nil
}

func (s *AIAnswerService) runJSONPrompt(ctx context.Context, prompt string, out interface{}) error {
	switch s.provider {
	case "ollama":
		return s.runJSONPromptOllama(ctx, prompt, out)
	case "cloudflare":
		return s.runJSONPromptCloudflare(ctx, prompt, out)
	default:
		return fmt.Errorf("unsupported ai provider: %s", s.provider)
	}
}

func (s *AIAnswerService) runJSONPromptCloudflare(ctx context.Context, prompt string, out interface{}) error {
	type cloudflareError struct {
		Message string `json:"message"`
	}
	type cloudflareResponse struct {
		Success bool              `json:"success"`
		Errors  []cloudflareError `json:"errors"`
		Result  interface{}       `json:"result"`
	}

	body := map[string]interface{}{
		"prompt":      prompt,
		"max_tokens":  s.maxTokens,
		"temperature": 0.2,
	}

	resp := cloudflareResponse{}
	res, err := s.client.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+s.apiKey).
		SetHeader("content-type", "application/json").
		SetBody(body).
		SetResult(&resp).
		Post("/run/" + strings.TrimPrefix(s.model, "/"))
	if err != nil {
		return fmt.Errorf("cloudflare request failed: %w", err)
	}
	if res.StatusCode() < 200 || res.StatusCode() >= 300 {
		return fmt.Errorf("cloudflare request failed with status %d: %s", res.StatusCode(), string(res.Body()))
	}
	if !resp.Success {
		errMsg := "cloudflare request was not successful"
		if len(resp.Errors) > 0 && strings.TrimSpace(resp.Errors[0].Message) != "" {
			errMsg = strings.TrimSpace(resp.Errors[0].Message)
		}
		return errors.New(errMsg)
	}

	text := strings.TrimSpace(extractCloudflareResultText(resp.Result))
	if text == "" {
		return fmt.Errorf("cloudflare response content is empty")
	}

	jsonPayload := extractJSONObject(text)
	if jsonPayload == "" {
		return fmt.Errorf("cloudflare response did not contain valid json")
	}

	if err := json.Unmarshal([]byte(jsonPayload), out); err != nil {
		return fmt.Errorf("failed to parse cloudflare json payload: %w", err)
	}

	return nil
}

func (s *AIAnswerService) runJSONPromptOllama(ctx context.Context, prompt string, out interface{}) error {
	type ollamaMessage struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	type ollamaResponse struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}

	body := map[string]interface{}{
		"model": s.model,
		"messages": []ollamaMessage{
			{Role: "user", Content: prompt},
		},
		"stream": false,
		"options": map[string]interface{}{
			"temperature": 0.2,
		},
	}

	resp := ollamaResponse{}
	res, err := s.client.R().
		SetContext(ctx).
		SetHeader("content-type", "application/json").
		SetBody(body).
		SetResult(&resp).
		Post("/api/chat")
	if err != nil {
		return fmt.Errorf("ollama request failed: %w", err)
	}
	if res.StatusCode() < 200 || res.StatusCode() >= 300 {
		return fmt.Errorf("ollama request failed with status %d: %s", res.StatusCode(), string(res.Body()))
	}

	text := strings.TrimSpace(resp.Message.Content)
	if text == "" {
		return fmt.Errorf("ollama response content is empty")
	}

	jsonPayload := extractJSONObject(text)
	if jsonPayload == "" {
		return fmt.Errorf("ollama response did not contain valid json")
	}

	if err := json.Unmarshal([]byte(jsonPayload), out); err != nil {
		return fmt.Errorf("failed to parse ollama json payload: %w", err)
	}

	return nil
}

func extractJSONObject(input string) string {
	input = strings.TrimSpace(input)
	if strings.HasPrefix(input, "{") && strings.HasSuffix(input, "}") {
		return input
	}

	re := regexp.MustCompile(`(?s)\{.*\}`)
	match := re.FindString(input)
	return strings.TrimSpace(match)
}

func normalizeQA(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, "`", "")
	re := regexp.MustCompile(`\s+`)
	value = re.ReplaceAllString(value, " ")
	return value
}

func extractCloudflareResultText(result interface{}) string {
	switch v := result.(type) {
	case string:
		return v
	case map[string]interface{}:
		for _, key := range []string{"response", "text", "output_text"} {
			if raw, ok := v[key]; ok {
				if text, ok := raw.(string); ok {
					return text
				}
			}
		}

		if raw, ok := v["messages"]; ok {
			if messages, ok := raw.([]interface{}); ok {
				for i := len(messages) - 1; i >= 0; i-- {
					messageMap, ok := messages[i].(map[string]interface{})
					if !ok {
						continue
					}
					content, ok := messageMap["content"]
					if !ok {
						continue
					}
					if text, ok := content.(string); ok {
						return text
					}
				}
			}
		}
	}

	return ""
}
