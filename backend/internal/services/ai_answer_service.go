package services

import (
	"context"
	"encoding/json"
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
		provider = "anthropic"
	}

	apiKey := ""
	model := ""
	baseURL := ""

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
	case "anthropic":
		apiKey = strings.TrimSpace(os.Getenv("ANTHROPIC_API_KEY"))
		baseURL = strings.TrimSpace(os.Getenv("ANTHROPIC_BASE_URL"))
		if baseURL == "" {
			baseURL = "https://api.anthropic.com"
		}
		model = strings.TrimSpace(os.Getenv("ANTHROPIC_MODEL"))
		if model == "" {
			model = "claude-3-5-sonnet-20241022"
		}
	default:
		// Fallback to anthropic for unknown provider values
		provider = "anthropic"
		apiKey = strings.TrimSpace(os.Getenv("ANTHROPIC_API_KEY"))
		baseURL = "https://api.anthropic.com"
		model = strings.TrimSpace(os.Getenv("ANTHROPIC_MODEL"))
		if model == "" {
			model = "claude-3-5-sonnet-20241022"
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

func (s *AIAnswerService) IsConfigured() bool {
	if s == nil {
		return false
	}

	switch s.provider {
	case "ollama":
		return s.model != ""
	case "anthropic":
		return s.apiKey != "" && s.model != ""
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

func (s *AIAnswerService) runJSONPrompt(ctx context.Context, prompt string, out interface{}) error {
	switch s.provider {
	case "ollama":
		return s.runJSONPromptOllama(ctx, prompt, out)
	case "anthropic":
		return s.runJSONPromptAnthropic(ctx, prompt, out)
	default:
		return fmt.Errorf("unsupported ai provider: %s", s.provider)
	}
}

func (s *AIAnswerService) runJSONPromptAnthropic(ctx context.Context, prompt string, out interface{}) error {
	type anthropicContent struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	type anthropicResponse struct {
		Content []anthropicContent `json:"content"`
	}

	body := map[string]interface{}{
		"model":       s.model,
		"max_tokens":  s.maxTokens,
		"temperature": 0.2,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}

	resp := anthropicResponse{}
	res, err := s.client.R().
		SetContext(ctx).
		SetHeader("x-api-key", s.apiKey).
		SetHeader("anthropic-version", "2023-06-01").
		SetHeader("content-type", "application/json").
		SetBody(body).
		SetResult(&resp).
		Post("/v1/messages")
	if err != nil {
		return fmt.Errorf("anthropic request failed: %w", err)
	}
	if res.StatusCode() < 200 || res.StatusCode() >= 300 {
		return fmt.Errorf("anthropic request failed with status %d: %s", res.StatusCode(), string(res.Body()))
	}

	if len(resp.Content) == 0 {
		return fmt.Errorf("anthropic response content is empty")
	}

	text := strings.TrimSpace(resp.Content[0].Text)
	jsonPayload := extractJSONObject(text)
	if jsonPayload == "" {
		return fmt.Errorf("anthropic response did not contain valid json")
	}

	if err := json.Unmarshal([]byte(jsonPayload), out); err != nil {
		return fmt.Errorf("failed to parse anthropic json payload: %w", err)
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
