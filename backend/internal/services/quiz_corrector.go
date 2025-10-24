package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/caplo84/quizz-backend/internal/models"
	"github.com/caplo84/quizz-backend/internal/services/datasources"
	"gorm.io/gorm"
)

type CorrectionOptions struct {
	QuizSlug            string
	DryRun              bool
	BatchSize           int
	ConfidenceThreshold float64
	Verbose             bool
}

type ValidationResult struct {
	IsValid          bool
	Issues           []string
	NonEmptyChoices  int
	CorrectChoices   int
	TotalChoiceCount int
}

type QuestionCorrectionDetail struct {
	QuestionID  uint     `json:"question_id"`
	QuizSlug    string   `json:"quiz_slug"`
	Action      string   `json:"action"`
	Source      string   `json:"source,omitempty"`
	Confidence  float64  `json:"confidence,omitempty"`
	Issues      []string `json:"issues,omitempty"`
	Error       string   `json:"error,omitempty"`
	CorrectedAt string   `json:"corrected_at,omitempty"`
}

type CorrectionReport struct {
	Timestamp        time.Time                  `json:"timestamp"`
	Duration         string                     `json:"duration"`
	DryRun           bool                       `json:"dry_run"`
	TotalProcessed   int                        `json:"total_processed"`
	TotalFixed       int                        `json:"total_fixed"`
	TotalSkipped     int                        `json:"total_skipped"`
	TotalFailed      int                        `json:"total_failed"`
	BySource         map[string]int             `json:"by_source"`
	ByConfidence     map[string]int             `json:"by_confidence"`
	EstimatedAPICost float64                    `json:"estimated_api_cost"`
	Details          []QuestionCorrectionDetail `json:"details,omitempty"`
}

type CorrectionPayload struct {
	QuestionText string
	Explanation  *string
	Choices      []string
	CorrectIndex int
	Source       string
	Confidence   float64
}

type QuizCorrectorService struct {
	db          *gorm.DB
	github      *datasources.GitHubClient
	ai          *AIAnswerService
	parsedCache map[string][]datasources.ParsedQuestion
}

type QuestionCorrector interface {
	CorrectAllQuizzes(ctx context.Context, opts CorrectionOptions) (*CorrectionReport, error)
}

func NewQuizCorrectorService(db *gorm.DB, github *datasources.GitHubClient, ai *AIAnswerService) *QuizCorrectorService {
	return &QuizCorrectorService{
		db:          db,
		github:      github,
		ai:          ai,
		parsedCache: make(map[string][]datasources.ParsedQuestion),
	}
}

func (s *QuizCorrectorService) CorrectAllQuizzes(ctx context.Context, opts CorrectionOptions) (*CorrectionReport, error) {
	started := time.Now()
	if opts.BatchSize <= 0 {
		opts.BatchSize = 100
	}
	if opts.ConfidenceThreshold <= 0 {
		opts.ConfidenceThreshold = 0.7
	}

	report := &CorrectionReport{
		Timestamp:    started.UTC(),
		DryRun:       opts.DryRun,
		BySource:     map[string]int{"parsed": 0, "ai-extracted": 0, "ai-generated": 0},
		ByConfidence: map[string]int{"high_0.9_1.0": 0, "medium_0.7_0.9": 0, "low_0.0_0.7": 0},
	}

	questions, err := s.loadQuestions(ctx, opts.QuizSlug)
	if err != nil {
		return nil, err
	}
	report.TotalProcessed = len(questions)

	for i := 0; i < len(questions); i += opts.BatchSize {
		end := i + opts.BatchSize
		if end > len(questions) {
			end = len(questions)
		}

		batch := questions[i:end]
		batchFailures := 0
		for idx := range batch {
			question := &batch[idx]
			validation := s.ValidateQuestion(question)
			if validation.IsValid {
				report.TotalSkipped++
				continue
			}

			payload, fixErr := s.buildCorrectionPayload(ctx, question, validation)
			if fixErr != nil {
				report.TotalFailed++
				batchFailures++
				report.Details = append(report.Details, QuestionCorrectionDetail{
					QuestionID: question.ID,
					QuizSlug:   question.Quiz.Slug,
					Action:     "failed",
					Issues:     validation.Issues,
					Error:      fixErr.Error(),
				})
				continue
			}

			if payload.Confidence < opts.ConfidenceThreshold {
				report.TotalFailed++
				batchFailures++
				report.Details = append(report.Details, QuestionCorrectionDetail{
					QuestionID: question.ID,
					QuizSlug:   question.Quiz.Slug,
					Action:     "failed",
					Source:     payload.Source,
					Confidence: payload.Confidence,
					Issues:     validation.Issues,
					Error:      "confidence below threshold",
				})
				continue
			}

			if applyErr := s.applyCorrection(ctx, question, payload, opts.DryRun); applyErr != nil {
				report.TotalFailed++
				batchFailures++
				report.Details = append(report.Details, QuestionCorrectionDetail{
					QuestionID: question.ID,
					QuizSlug:   question.Quiz.Slug,
					Action:     "failed",
					Source:     payload.Source,
					Confidence: payload.Confidence,
					Issues:     validation.Issues,
					Error:      applyErr.Error(),
				})
				continue
			}

			report.TotalFixed++
			report.BySource[payload.Source]++
			s.addConfidenceBucket(report, payload.Confidence)
			report.Details = append(report.Details, QuestionCorrectionDetail{
				QuestionID:  question.ID,
				QuizSlug:    question.Quiz.Slug,
				Action:      "fixed",
				Source:      payload.Source,
				Confidence:  payload.Confidence,
				Issues:      validation.Issues,
				CorrectedAt: time.Now().UTC().Format(time.RFC3339),
			})
		}

		if len(batch) > 0 && float64(batchFailures)/float64(len(batch)) > 0.2 {
			log.Printf("warning: batch failure ratio exceeded 20%% (%d/%d)", batchFailures, len(batch))
		}
	}

	report.Duration = time.Since(started).String()
	report.EstimatedAPICost = estimateCost(report.BySource)
	return report, nil
}

func (s *QuizCorrectorService) ValidateQuestion(q *models.Question) ValidationResult {
	result := ValidationResult{IsValid: true, Issues: []string{}}

	if strings.TrimSpace(q.QuestionText) == "" {
		result.IsValid = false
		result.Issues = append(result.Issues, "question_text is empty")
	}

	result.TotalChoiceCount = len(q.Choices)
	for _, choice := range q.Choices {
		if strings.TrimSpace(choice.ChoiceText) != "" {
			result.NonEmptyChoices++
		}
		if choice.IsCorrect {
			result.CorrectChoices++
		}
	}

	if q.QuestionType == "multiple_choice" || q.QuestionType == "" {
		if result.NonEmptyChoices < 2 {
			result.IsValid = false
			result.Issues = append(result.Issues, "insufficient non-empty choices")
		}
		if result.CorrectChoices != 1 {
			result.IsValid = false
			result.Issues = append(result.Issues, fmt.Sprintf("expected exactly 1 correct answer, got %d", result.CorrectChoices))
		}
	}

	return result
}

func (s *QuizCorrectorService) ExtractAnswerFromSource(ctx context.Context, q *models.Question) (*CorrectionPayload, error) {
	parsed, matchConfidence, err := s.findParsedQuestion(ctx, q)
	if err == nil && parsed != nil && len(parsed.Choices) > 1 && parsed.Correct >= 0 && parsed.Correct < len(parsed.Choices) {
		payload := &CorrectionPayload{
			QuestionText: strings.TrimSpace(parsed.Text),
			Choices:      dedupeAndTrim(parsed.Choices),
			CorrectIndex: parsed.Correct,
			Source:       "parsed",
			Confidence:   matchConfidence,
		}
		if q.Explanation != nil {
			payload.Explanation = q.Explanation
		}
		return payload, nil
	}

	if s.ai == nil || !s.ai.IsConfigured() {
		return nil, fmt.Errorf("no parsed source match and ai service is not configured")
	}

	nonEmptyChoices := extractNonEmptyChoiceTexts(q.Choices)
	if len(nonEmptyChoices) >= 2 {
		syntheticMarkdown := buildSyntheticMarkdown(q.QuestionText, nonEmptyChoices)
		extracted, extractErr := s.ai.ExtractAnswerFromMarkdown(ctx, q.QuestionText, syntheticMarkdown)
		if extractErr == nil && extracted != nil && extracted.Found {
			correctIndex := indexByNormalizedMatch(nonEmptyChoices, extracted.CorrectAnswer)
			if correctIndex >= 0 {
				choices := nonEmptyChoices
				if len(choices) < 4 {
					wrong, wrongErr := s.ai.GenerateWrongAnswers(ctx, choices[correctIndex], q.QuestionText, q.Quiz.Topic.Name)
					if wrongErr == nil {
						for _, candidate := range wrong {
							if len(choices) >= 4 {
								break
							}
							if indexByNormalizedMatch(choices, candidate) == -1 {
								choices = append(choices, candidate)
							}
						}
					}
				}
				return &CorrectionPayload{
					QuestionText: strings.TrimSpace(q.QuestionText),
					Choices:      choices,
					CorrectIndex: correctIndex,
					Source:       "ai-extracted",
					Confidence:   extracted.Confidence,
				}, nil
			}
		}
	}

	answerSet, err := s.ai.GenerateAnswerSet(ctx, q.QuestionText, q.Quiz.Topic.Name, q.Quiz.DifficultyLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to generate answer set: %w", err)
	}

	choices := []string{answerSet.CorrectAnswer}
	choices = append(choices, answerSet.WrongAnswers...)
	payload := &CorrectionPayload{
		QuestionText: strings.TrimSpace(q.QuestionText),
		Choices:      dedupeAndTrim(choices),
		CorrectIndex: 0,
		Source:       "ai-generated",
		Confidence:   answerSet.Confidence,
	}
	if strings.TrimSpace(answerSet.Explanation) != "" {
		explanation := strings.TrimSpace(answerSet.Explanation)
		payload.Explanation = &explanation
	}

	return payload, nil
}

func (s *QuizCorrectorService) GenerateMissingChoices(ctx context.Context, q *models.Question, correct string) ([]string, error) {
	if s.ai == nil || !s.ai.IsConfigured() {
		return nil, fmt.Errorf("ai service is not configured")
	}
	return s.ai.GenerateWrongAnswers(ctx, correct, q.QuestionText, q.Quiz.Topic.Name)
}

func (s *QuizCorrectorService) UpdateQuestion(ctx context.Context, q *models.Question, payload *CorrectionPayload, dryRun bool) error {
	return s.applyCorrection(ctx, q, payload, dryRun)
}

func (s *QuizCorrectorService) loadQuestions(ctx context.Context, quizSlug string) ([]models.Question, error) {
	var questions []models.Question
	query := s.db.WithContext(ctx).
		Model(&models.Question{}).
		Preload("Choices").
		Preload("Quiz").
		Preload("Quiz.Topic").
		Order("questions.id ASC")

	if strings.TrimSpace(quizSlug) != "" {
		query = query.Joins("JOIN quizzes ON quizzes.id = questions.quiz_id").Where("quizzes.slug = ?", quizSlug)
	}

	if err := query.Find(&questions).Error; err != nil {
		return nil, fmt.Errorf("failed to load questions: %w", err)
	}

	return questions, nil
}

func (s *QuizCorrectorService) buildCorrectionPayload(ctx context.Context, q *models.Question, validation ValidationResult) (*CorrectionPayload, error) {
	_ = validation
	payload, err := s.ExtractAnswerFromSource(ctx, q)
	if err != nil {
		return nil, err
	}
	if len(payload.Choices) < 2 {
		return nil, fmt.Errorf("correction payload has insufficient choices")
	}
	if payload.CorrectIndex < 0 || payload.CorrectIndex >= len(payload.Choices) {
		return nil, fmt.Errorf("correction payload has invalid correct index")
	}
	return payload, nil
}

func (s *QuizCorrectorService) applyCorrection(ctx context.Context, q *models.Question, payload *CorrectionPayload, dryRun bool) error {
	if dryRun {
		return nil
	}

	now := time.Now().UTC()
	oldChoicesBytes, _ := json.Marshal(q.Choices)
	newChoicesPreview := make([]map[string]interface{}, 0, len(payload.Choices))
	for i, choiceText := range payload.Choices {
		newChoicesPreview = append(newChoicesPreview, map[string]interface{}{
			"choice_text": choiceText,
			"is_correct":  i == payload.CorrectIndex,
			"order_index": i + 1,
		})
	}
	newChoicesBytes, _ := json.Marshal(newChoicesPreview)

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{
			"updated_at": now,
		}
		if strings.TrimSpace(payload.QuestionText) != "" {
			updates["question_text"] = strings.TrimSpace(payload.QuestionText)
		}
		if payload.Explanation != nil {
			updates["explanation"] = payload.Explanation
		}

		if err := tx.Model(&models.Question{}).Where("id = ?", q.ID).Updates(updates).Error; err != nil {
			return fmt.Errorf("failed to update question: %w", err)
		}

		if err := tx.Where("question_id = ?", q.ID).Delete(&models.Choice{}).Error; err != nil {
			return fmt.Errorf("failed to clear old choices: %w", err)
		}

		for i, choiceText := range payload.Choices {
			choice := &models.Choice{
				QuestionID:   q.ID,
				ChoiceText:   strings.TrimSpace(choiceText),
				IsCorrect:    i == payload.CorrectIndex,
				OrderIndex:   i + 1,
				AnswerSource: payload.Source,
				AIConfidence: payload.Confidence,
				CorrectedAt:  &now,
			}
			if err := tx.Create(choice).Error; err != nil {
				return fmt.Errorf("failed to create corrected choice: %w", err)
			}
		}

		_ = tx.Exec(
			`INSERT INTO correction_audit_log (question_id, action, old_data, new_data, confidence, created_by) VALUES (?, ?, ?, ?, ?, ?)`,
			q.ID,
			"question_corrected",
			string(oldChoicesBytes),
			string(newChoicesBytes),
			payload.Confidence,
			"quiz-corrector",
		).Error

		return nil
	})
}

func (s *QuizCorrectorService) findParsedQuestion(ctx context.Context, q *models.Question) (*datasources.ParsedQuestion, float64, error) {
	if s.github == nil {
		return nil, 0, fmt.Errorf("github client is nil")
	}

	category := strings.TrimSpace(q.Quiz.Topic.Slug)
	if category == "" {
		category = strings.TrimSpace(q.Quiz.Topic.Name)
	}
	if category == "" {
		return nil, 0, fmt.Errorf("question has no topic category")
	}

	cacheKey := strings.ToLower(category)
	parsedQuestions, ok := s.parsedCache[cacheKey]
	if !ok {
		fetched, err := s.github.FetchParsedQuestions(ctx, datasources.FetchParams{
			Category: category,
			Limit:    500,
		})
		if err != nil {
			return nil, 0, err
		}
		parsedQuestions = fetched
		s.parsedCache[cacheKey] = parsedQuestions
	}

	normQuestion := normalizeQuestionText(q.QuestionText)
	for i := range parsedQuestions {
		pq := &parsedQuestions[i]
		if q.ExternalID != nil && *q.ExternalID != "" && *q.ExternalID == pq.ExternalID {
			return pq, 0.99, nil
		}
		if q.ExternalReference != nil && *q.ExternalReference != "" && *q.ExternalReference == pq.ExternalReference {
			return pq, 0.97, nil
		}
		if normalizeQuestionText(pq.Text) == normQuestion {
			return pq, 0.94, nil
		}
	}

	return nil, 0, fmt.Errorf("no parsed source match")
}

func (s *QuizCorrectorService) addConfidenceBucket(report *CorrectionReport, confidence float64) {
	if confidence >= 0.9 {
		report.ByConfidence["high_0.9_1.0"]++
		return
	}
	if confidence >= 0.7 {
		report.ByConfidence["medium_0.7_0.9"]++
		return
	}
	report.ByConfidence["low_0.0_0.7"]++
}

func estimateCost(bySource map[string]int) float64 {
	parsed := float64(bySource["parsed"]) * 0.0002
	aiextracted := float64(bySource["ai-extracted"]) * 0.004
	aigenerated := float64(bySource["ai-generated"]) * 0.006
	return parsed + aiextracted + aigenerated
}

func extractNonEmptyChoiceTexts(choices []models.Choice) []string {
	result := make([]string, 0, len(choices))
	for _, c := range choices {
		trimmed := strings.TrimSpace(c.ChoiceText)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return dedupeAndTrim(result)
}

func buildSyntheticMarkdown(question string, options []string) string {
	var sb strings.Builder
	sb.WriteString("#### Question\n")
	sb.WriteString(strings.TrimSpace(question))
	sb.WriteString("\n\n#### Choices\n")
	for _, opt := range options {
		sb.WriteString("- ")
		sb.WriteString(strings.TrimSpace(opt))
		sb.WriteString("\n")
	}
	return sb.String()
}

func indexByNormalizedMatch(values []string, target string) int {
	normTarget := normalizeQuestionText(target)
	for i, value := range values {
		normValue := normalizeQuestionText(value)
		if normValue == normTarget {
			return i
		}
		if strings.Contains(normValue, normTarget) || strings.Contains(normTarget, normValue) {
			return i
		}
	}
	return -1
}

func dedupeAndTrim(values []string) []string {
	result := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		norm := normalizeQuestionText(trimmed)
		if _, exists := seen[norm]; exists {
			continue
		}
		seen[norm] = struct{}{}
		result = append(result, trimmed)
	}
	return result
}

func normalizeQuestionText(input string) string {
	value := strings.ToLower(strings.TrimSpace(input))
	value = strings.ReplaceAll(value, "`", "")
	value = strings.ReplaceAll(value, "\n", " ")
	space := regexp.MustCompile(`\s+`)
	return space.ReplaceAllString(value, " ")
}
