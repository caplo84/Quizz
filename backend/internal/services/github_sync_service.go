package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/caplo84/quizz-backend/internal/models"
	"github.com/caplo84/quizz-backend/internal/repository"
	"github.com/caplo84/quizz-backend/internal/services/datasources"
	"github.com/caplo84/quizz-backend/pkg/utils"
)

// GitHubSyncService interface defines the contract for GitHub synchronization
type GitHubSyncService interface {
	SyncFromGitHub(ctx context.Context) error
	GetRateLimit(ctx context.Context) (*datasources.RateLimit, error)
}

// SyncConfig holds configuration for the sync process
type SyncConfig struct {
	OnlyProgrammingLanguages bool
	SpecificCategories       []string
}

type GitHubSyncServiceImpl struct {
	githubClient *datasources.GitHubClient
	quizRepo     repository.QuizRepository
	topicRepo    repository.TopicRepository
	config       SyncConfig
}

func NewGitHubSyncService(
	githubClient *datasources.GitHubClient,
	quizRepo repository.QuizRepository,
	topicRepo repository.TopicRepository,
) GitHubSyncService {
	// Check if specific categories are configured via environment
	configuredCategories := os.Getenv("SYNC_CATEGORIES")
	onlyProgramming := os.Getenv("SYNC_ONLY_PROGRAMMING") != "false" // Default to true

	config := SyncConfig{
		OnlyProgrammingLanguages: onlyProgramming,
	}

	// Parse configured categories if provided
	if configuredCategories != "" {
		categories := strings.Split(configuredCategories, ",")
		for i, cat := range categories {
			categories[i] = strings.TrimSpace(cat)
		}
		config.SpecificCategories = categories
		config.OnlyProgrammingLanguages = false // Override if specific categories are set
	}

	return &GitHubSyncServiceImpl{
		githubClient: githubClient,
		quizRepo:     quizRepo,
		topicRepo:    topicRepo,
		config:       config,
	}
}

func (s *GitHubSyncServiceImpl) SyncFromGitHub(ctx context.Context) error {
	log.Println("Starting GitHub sync...")

	// Step 1: Check authentication
	if err := s.githubClient.ValidateAuth(ctx); err != nil {
		return fmt.Errorf("GitHub authentication failed: %w", err)
	}
	log.Println("✅ GitHub authentication successful")

	// Step 2: Check rate limits
	rateLimit, err := s.githubClient.GetRateLimit(ctx)
	if err != nil {
		return fmt.Errorf("failed to check rate limit: %w", err)
	}
	log.Printf("📊 GitHub rate limit: %d/%d remaining", rateLimit.Remaining, rateLimit.Total)

	if rateLimit.Remaining < 100 {
		return fmt.Errorf("rate limit too low: %d remaining (need at least 100)", rateLimit.Remaining)
	}

	// Step 3: Fetch categories (topics)
	categories, err := s.githubClient.FetchCategories(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch categories: %w", err)
	}

	log.Printf("Found %d categories", len(categories))

	var categoriesToProcess []datasources.Category

	// Step 4: Apply filtering based on configuration
	if len(s.config.SpecificCategories) > 0 {
		// Use configured specific categories
		log.Printf("Using configured categories: %v", s.config.SpecificCategories)
		for _, configCat := range s.config.SpecificCategories {
			for _, cat := range categories {
				if strings.EqualFold(cat.Name, configCat) {
					categoriesToProcess = append(categoriesToProcess, cat)
					break
				}
			}
		}
		log.Printf("Found %d configured categories", len(categoriesToProcess))
	} else if s.config.OnlyProgrammingLanguages {
		// Filter for programming languages
		categoriesToProcess = s.filterProgrammingLanguages(categories)
		log.Printf("Found %d programming language categories", len(categoriesToProcess))
	} else {
		// Use all categories
		categoriesToProcess = categories
		log.Printf("Using all %d categories", len(categoriesToProcess))
	}

	if len(categoriesToProcess) == 0 {
		log.Println("No categories to process after filtering")
		return nil
	}

	// Step 5: Process each filtered category
	successCount := 0
	for _, category := range categoriesToProcess {
		if err := s.processCategoryQuizzes(ctx, category.Name); err != nil {
			log.Printf("❌ Error processing category %s: %v", category.Name, err)
			continue // Continue with other categories
		}
		successCount++
	}

	log.Printf("✅ GitHub sync completed successfully! Processed %d/%d categories", successCount, len(categoriesToProcess))
	return nil
}

func (s *GitHubSyncServiceImpl) processCategoryQuizzes(ctx context.Context, categoryName string) error {
	log.Printf("Processing category: %s", categoryName)

	// Step 1: Ensure topic exists in database
	topic, err := s.ensureTopicExists(ctx, categoryName)
	if err != nil {
		return fmt.Errorf("failed to ensure topic exists: %w", err)
	}

	// Step 2: Fetch questions for this category
	params := datasources.FetchParams{
		Category: categoryName,
		Limit:    100, // Adjust as needed
	}

	questions, err := s.githubClient.FetchQuestions(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to fetch questions for category %s: %w", categoryName, err)
	}

	if len(questions) == 0 {
		log.Printf("No questions found for category: %s", categoryName)
		return nil
	}

	// Step 3: Group questions by quiz (by external reference)
	quizGroups := s.groupQuestionsByQuiz(questions)

	// Step 4: Save each quiz to database
	for quizIdentifier, quizQuestions := range quizGroups {
		if err := s.saveQuizToDB(ctx, topic, quizIdentifier, quizQuestions); err != nil {
			log.Printf("Error saving quiz %s: %v", quizIdentifier, err)
			continue
		}
	}

	return nil
}

func (s *GitHubSyncServiceImpl) ensureTopicExists(ctx context.Context, categoryName string) (*models.Topic, error) {
	// Try to find existing topic
	topic, err := s.topicRepo.GetBySlug(ctx, utils.GenerateSlug(categoryName))
	if err == nil {
		return topic, nil // Topic already exists
	}

	// Create new topic with improved naming
	displayName := s.generateTopicDisplayName(categoryName)
	newTopic := &models.Topic{
		Name:        displayName,
		Slug:        utils.GenerateSlug(categoryName), // Keep original for consistency
		Description: stringPtr(fmt.Sprintf("Programming quiz questions for %s", displayName)),
		IsActive:    true,
	}

	if err := s.topicRepo.Create(ctx, newTopic); err != nil {
		return nil, fmt.Errorf("failed to create topic: %w", err)
	}

	log.Printf("📁 Created new topic: %s (slug: %s)", displayName, newTopic.Slug)
	return newTopic, nil
}

func (s *GitHubSyncServiceImpl) groupQuestionsByQuiz(questions []datasources.Question) map[string][]datasources.Question {
	quizGroups := make(map[string][]datasources.Question)

	for _, question := range questions {
		// Use external reference (file URL) as quiz identifier
		quizIdentifier := question.ExternalReference
		if quizIdentifier == "" {
			quizIdentifier = question.Category // Fallback to category
		}

		quizGroups[quizIdentifier] = append(quizGroups[quizIdentifier], question)
	}

	return quizGroups
}

func (s *GitHubSyncServiceImpl) saveQuizToDB(ctx context.Context, topic *models.Topic, quizIdentifier string, questions []datasources.Question) error {
	if len(questions) == 0 {
		return nil
	}

	// Check if quiz already exists
	existingQuiz, err := s.quizRepo.GetQuizByExternalID(ctx, questions[0].ExternalID)
	if err == nil && existingQuiz != nil {
		log.Printf("Quiz already exists, skipping: %s", quizIdentifier)
		return nil
	}

	// Create quiz title from first question or identifier
	quizTitle := s.generateQuizTitle(questions[0], quizIdentifier)

	// Create quiz
	quiz := &models.Quiz{
		Title:             quizTitle,
		Slug:              utils.GenerateSlug(quizTitle),
		Description:       stringPtr(fmt.Sprintf("Quiz imported from GitHub: %s", quizIdentifier)),
		TopicID:           topic.ID,
		DifficultyLevel:   "medium", // Default, you can enhance this later
		TimeLimitMinutes:  30,
		TotalQuestions:    len(questions),
		IsActive:          true,
		Source:            stringPtr("github"),
		ExternalReference: stringPtr(quizIdentifier),
		ExternalID:        stringPtr(questions[0].ExternalID),
		LastSyncedAt:      timePtr(time.Now()),
	}

	// Save quiz
	if err := s.quizRepo.Create(ctx, quiz); err != nil {
		return fmt.Errorf("failed to create quiz: %w", err)
	}

	// Save questions
	for i, q := range questions {
		modelQuestion := &models.Question{
			QuizID:            quiz.ID,
			QuestionText:      q.Text,
			QuestionType:      "multiple_choice",
			Points:            1,
			OrderIndex:        i + 1,
			IsActive:          true,
			Source:            stringPtr("github"),
			ExternalReference: stringPtr(q.ExternalReference),
			ExternalID:        stringPtr(q.ExternalID),
		}

		if err := s.quizRepo.CreateQuestion(ctx, modelQuestion); err != nil {
			log.Printf("Error creating question: %v", err)
			continue
		}

		// Save choices
		for j, choice := range q.Choices {
			modelChoice := &models.Choice{
				QuestionID: modelQuestion.ID,
				ChoiceText: choice,
				IsCorrect:  j == q.Correct,
				OrderIndex: j + 1,
			}

			if err := s.quizRepo.CreateChoice(ctx, modelChoice); err != nil {
				log.Printf("Error creating choice: %v", err)
			}
		}
	}

	log.Printf("Saved quiz: %s with %d questions", quizTitle, len(questions))
	return nil
}

func (s *GitHubSyncServiceImpl) generateQuizTitle(firstQuestion datasources.Question, identifier string) string {
	// Extract filename from identifier (GitHub file path)
	// Remove query parameters first (e.g., ?ref=main)
	cleanIdentifier := strings.Split(identifier, "?")[0]

	parts := strings.Split(cleanIdentifier, "/")
	var filename string
	if len(parts) > 0 {
		filename = parts[len(parts)-1]
		// Remove .md extension and URL encoding
		filename = strings.TrimSuffix(filename, ".md")
		filename = strings.TrimSuffix(filename, ".MD")                    // Handle uppercase
		filename = strings.ReplaceAll(filename, "%E2%80%8B%E2%80%8B", "") // Remove zero-width spaces
		filename = strings.ReplaceAll(filename, "%20", " ")               // Replace URL-encoded spaces
	}

	// If we have a filename, use it; otherwise fall back to category
	var baseName string
	if filename != "" {
		baseName = filename
	} else if firstQuestion.Category != "" {
		baseName = firstQuestion.Category
	} else {
		return "Programming Quiz"
	}

	// Clean up the name and generate a proper title
	return s.generateCleanQuizTitle(baseName)
}

// generateCleanQuizTitle creates clean, professional quiz titles
func (s *GitHubSyncServiceImpl) generateCleanQuizTitle(rawName string) string {
	// Convert to lowercase for processing
	name := strings.ToLower(strings.TrimSpace(rawName))

	// Remove common prefixes/suffixes that create redundancy
	name = strings.TrimSuffix(name, "-quiz")
	name = strings.TrimSuffix(name, "_quiz")
	name = strings.TrimSuffix(name, "quiz")
	name = strings.TrimPrefix(name, "quiz-")
	name = strings.TrimPrefix(name, "quiz_")

	// Handle language codes and variations
	languageMap := map[string]string{
		"-es": " (Spanish)",
		"-fr": " (French)",
		"-it": " (Italian)",
		"-de": " (German)",
		"-pt": " (Portuguese)",
		"-ru": " (Russian)",
		"-ja": " (Japanese)",
		"-ko": " (Korean)",
		"-zh": " (Chinese)",
		"-tr": " (Turkish)",
		"-ua": " (Ukrainian)",
		"-ch": " (Chinese)",
		"_es": " (Spanish)",
		"_fr": " (French)",
		"_it": " (Italian)",
		"_de": " (German)",
		"_pt": " (Portuguese)",
		"_ru": " (Russian)",
		"_ja": " (Japanese)",
		"_ko": " (Korean)",
		"_zh": " (Chinese)",
		"_tr": " (Turkish)",
		"_ua": " (Ukrainian)",
		"_ch": " (Chinese)",
	}

	// Check for language codes and extract them
	languageSuffix := ""
	for code, language := range languageMap {
		if strings.HasSuffix(name, code) {
			languageSuffix = language
			name = strings.TrimSuffix(name, code)
			break
		}
	}

	// After removing language codes, remove any remaining "quiz" references
	name = strings.TrimSuffix(name, "-quiz")
	name = strings.TrimSuffix(name, "_quiz")
	name = strings.TrimSuffix(name, "quiz")
	name = strings.TrimPrefix(name, "quiz-")
	name = strings.TrimPrefix(name, "quiz_")

	// Use centralized technology mapping
	techMap := s.getTechnologyMap()
	validKeys := s.getAllTechnologyKeys()

	// Check for exact technology matches using centralized config
	displayName := ""
	if primaryKey, exists := validKeys[name]; exists {
		if config, found := techMap[primaryKey]; found {
			displayName = config.DisplayName
		}
	}

	// Fallback: clean up the name manually if no mapping found
	if displayName == "" {
		// Replace hyphens and underscores with spaces
		displayName = strings.ReplaceAll(name, "-", " ")
		displayName = strings.ReplaceAll(displayName, "_", " ")

		// Title case each word
		caser := cases.Title(language.Und)
		words := strings.Fields(displayName)
		for i, word := range words {
			words[i] = caser.String(strings.ToLower(word))
		}
		displayName = strings.Join(words, " ")
	}

	// Combine with language suffix if present
	if languageSuffix != "" {
		return fmt.Sprintf("%s%s", displayName, languageSuffix)
	}

	return displayName
}

// GetRateLimit returns the current GitHub API rate limit status
func (s *GitHubSyncServiceImpl) GetRateLimit(ctx context.Context) (*datasources.RateLimit, error) {
	return s.githubClient.GetRateLimit(ctx)
}

// TechConfig holds configuration for technology recognition and display
type TechConfig struct {
	DisplayName string   // How to display this technology
	Aliases     []string // Alternative names for the same technology
	Category    string   // Programming Language, Framework, Database, etc.
}

// getTechnologyMap returns the centralized technology configuration
func (s *GitHubSyncServiceImpl) getTechnologyMap() map[string]TechConfig {
	return TechConfigData
}

// getAllTechnologyKeys returns all valid technology identifiers (primary + aliases)
func (s *GitHubSyncServiceImpl) getAllTechnologyKeys() map[string]string {
	techMap := s.getTechnologyMap()
	allKeys := make(map[string]string)

	// Add primary keys
	for key := range techMap {
		allKeys[strings.ToLower(key)] = key
	}

	// Add aliases
	for primaryKey, config := range techMap {
		for _, alias := range config.Aliases {
			allKeys[strings.ToLower(alias)] = primaryKey
		}
	}

	return allKeys
}

// filterProgrammingLanguages filters categories to include only programming languages
func (s *GitHubSyncServiceImpl) filterProgrammingLanguages(categories []datasources.Category) []datasources.Category {
	validTechKeys := s.getAllTechnologyKeys()
	var filtered []datasources.Category

	for _, category := range categories {
		// Check if category matches any valid technology key
		categoryLower := strings.ToLower(category.Name)
		if primaryKey, exists := validTechKeys[categoryLower]; exists {
			log.Printf("✅ Including technology category: %s (mapped to: %s)", category.Name, primaryKey)
			filtered = append(filtered, category)
		} else {
			log.Printf("⏭️  Skipping non-technology category: %s", category.Name)
		}
	}

	log.Printf("Filtered from %d total categories to %d technology categories", len(categories), len(filtered))
	return filtered
}

// generateTopicDisplayName creates user-friendly topic names from GitHub category names
func (s *GitHubSyncServiceImpl) generateTopicDisplayName(categoryName string) string {
	// Check for custom topic name override from environment
	envKey := fmt.Sprintf("TOPIC_NAME_%s", strings.ToUpper(strings.ReplaceAll(categoryName, "-", "_")))
	if customName := os.Getenv(envKey); customName != "" {
		log.Printf("🎯 Using custom topic name from %s: %s", envKey, customName)
		return customName
	}

	// Use centralized technology mapping
	techMap := s.getTechnologyMap()
	validKeys := s.getAllTechnologyKeys()

	// Clean category name for lookup
	cleanName := strings.ToLower(strings.TrimSpace(categoryName))

	// Check if we have a primary key or alias match
	if primaryKey, exists := validKeys[cleanName]; exists {
		if config, found := techMap[primaryKey]; found {
			return config.DisplayName
		}
	}

	// Fallback: Clean up the original name
	return s.cleanupCategoryName(categoryName)
}

// cleanupCategoryName improves category names when no mapping exists
func (s *GitHubSyncServiceImpl) cleanupCategoryName(categoryName string) string {
	// Replace hyphens and underscores with spaces
	cleaned := strings.ReplaceAll(categoryName, "-", " ")
	cleaned = strings.ReplaceAll(cleaned, "_", " ")

	// Title case each word
	caser := cases.Title(language.Und)
	words := strings.Fields(cleaned)
	for i, word := range words {
		words[i] = caser.String(strings.ToLower(word))
	}

	return strings.Join(words, " ")
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}
