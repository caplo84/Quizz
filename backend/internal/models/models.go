package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Topic represents a quiz topic/category
type Topic struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:100;not null;uniqueIndex:idx_topics_name" validate:"required,min=2,max=100"`
	Slug        string    `json:"slug" gorm:"size:100;not null;uniqueIndex:idx_topics_slug" validate:"required,min=2,max=100"`
	Description *string   `json:"description,omitempty" gorm:"type:text"`
	IconURL     *string   `json:"icon_url,omitempty" gorm:"size:255"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Quizzes []Quiz `json:"quizzes,omitempty" gorm:"foreignKey:TopicID;constraint:OnDelete:CASCADE"`
}

// Quiz represents a quiz with questions
type Quiz struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	Title            string    `json:"title" gorm:"size:200;not null" validate:"required,min=3,max=200"`
	Slug             string    `json:"slug" gorm:"size:200;not null;uniqueIndex" validate:"required,min=3,max=200"`
	Description      *string   `json:"description,omitempty" gorm:"type:text"`
	TopicID          uint      `json:"topic_id" gorm:"not null;index" validate:"required"`
	DifficultyLevel  string    `json:"difficulty_level" gorm:"size:20;default:medium;check:difficulty_level IN ('easy','medium','hard')" validate:"oneof=easy medium hard"`
	TimeLimitMinutes int       `json:"time_limit_minutes" gorm:"default:30" validate:"min=1,max=180"`
	TotalQuestions   int       `json:"total_questions" gorm:"default:0"`
	IsActive         bool      `json:"is_active" gorm:"default:true"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// Relationships
	Topic     Topic      `json:"topic" gorm:"foreignKey:TopicID;constraint:OnDelete:CASCADE"`
	Questions []Question `json:"questions,omitempty" gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE"`
	Attempts  []Attempt  `json:"attempts,omitempty" gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE"`
}

// Question represents a quiz question
type Question struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	QuizID       uint      `json:"quiz_id" gorm:"not null;index" validate:"required"`
	QuestionText string    `json:"question_text" gorm:"type:text;not null" validate:"required,min=10"`
	QuestionType string    `json:"question_type" gorm:"size:20;default:multiple_choice;check:question_type IN ('multiple_choice','true_false','text')" validate:"oneof=multiple_choice true_false text"`
	Points       int       `json:"points" gorm:"default:1" validate:"min=1,max=100"`
	Explanation  *string   `json:"explanation,omitempty" gorm:"type:text"`
	OrderIndex   int       `json:"order_index" gorm:"not null" validate:"required,min=1"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relationships
	Quiz    Quiz     `json:"quiz" gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE"`
	Choices []Choice `json:"choices,omitempty" gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE"`
}

// Choice represents a multiple choice answer option
type Choice struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	QuestionID  uint      `json:"question_id" gorm:"not null;index" validate:"required"`
	ChoiceText  string    `json:"choice_text" gorm:"type:text;not null" validate:"required,min=1"`
	IsCorrect   bool      `json:"is_correct" gorm:"default:false"`
	OrderIndex  int       `json:"order_index" gorm:"not null" validate:"required,min=1"`
	Explanation *string   `json:"explanation,omitempty" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Question Question `json:"question" gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE"`
}

// UserAnswers represents the JSON structure for storing user answers
type UserAnswers map[string]interface{}

// Value implements the driver.Valuer interface for GORM
func (ua UserAnswers) Value() (driver.Value, error) {
	return json.Marshal(ua)
}

// Scan implements the sql.Scanner interface for GORM
func (ua *UserAnswers) Scan(value interface{}) error {
	if value == nil {
		*ua = make(UserAnswers)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, ua)
}

// Attempt represents a user's quiz attempt
type Attempt struct {
	ID               uint        `json:"id" gorm:"primaryKey"`
	QuizID           uint        `json:"quiz_id" gorm:"not null;index" validate:"required"`
	UserIdentifier   *string     `json:"user_identifier,omitempty" gorm:"size:255;index"`
	UserName         *string     `json:"user_name,omitempty" gorm:"size:100"`
	StartedAt        time.Time   `json:"started_at" gorm:"default:CURRENT_TIMESTAMP"`
	CompletedAt      *time.Time  `json:"completed_at,omitempty"`
	TotalScore       int         `json:"total_score" gorm:"default:0"`
	MaxPossibleScore int         `json:"max_possible_score" gorm:"default:0"`
	PercentageScore  float64     `json:"percentage_score" gorm:"type:decimal(5,2);default:0.00"`
	TimeTakenSeconds *int        `json:"time_taken_seconds,omitempty"`
	Status           string      `json:"status" gorm:"size:20;default:in_progress;check:status IN ('in_progress','completed','abandoned')" validate:"oneof=in_progress completed abandoned"`
	Answers          UserAnswers `json:"answers" gorm:"type:jsonb"`
	Score            int         `json:"score" gorm:"default:0"`            // Add this for compatibility
	IsCompleted      bool        `json:"is_completed" gorm:"default:false"` // Add this field
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`

	// Relationships
	Quiz Quiz `json:"quiz" gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE"`
}

// DTOs for API requests and responses

// TopicResponse represents topic data for API responses
type TopicResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description,omitempty"`
	IconURL     *string `json:"icon_url,omitempty"`
	QuizCount   int     `json:"quiz_count,omitempty"`
}

// QuizListResponse represents quiz data for list endpoints
type QuizListResponse struct {
	ID               uint    `json:"id"`
	Title            string  `json:"title"`
	Slug             string  `json:"slug"`
	Description      *string `json:"description,omitempty"`
	DifficultyLevel  string  `json:"difficulty_level"`
	TimeLimitMinutes int     `json:"time_limit_minutes"`
	TotalQuestions   int     `json:"total_questions"`
	TopicName        string  `json:"topic_name"`
	TopicSlug        string  `json:"topic_slug"`
}

// QuizDetailResponse represents detailed quiz data
type QuizDetailResponse struct {
	ID               uint          `json:"id"`
	Title            string        `json:"title"`
	Slug             string        `json:"slug"`
	Description      *string       `json:"description,omitempty"`
	DifficultyLevel  string        `json:"difficulty_level"`
	TimeLimitMinutes int           `json:"time_limit_minutes"`
	TotalQuestions   int           `json:"total_questions"`
	Topic            TopicResponse `json:"topic"`
}

// QuestionResponse represents question data for API responses (without correct answers)
type QuestionResponse struct {
	ID           uint             `json:"id"`
	QuestionText string           `json:"question_text"`
	QuestionType string           `json:"question_type"`
	Points       int              `json:"points"`
	OrderIndex   int              `json:"order_index"`
	Choices      []ChoiceResponse `json:"choices,omitempty"`
}

// ChoiceResponse represents choice data for API responses (without is_correct flag)
type ChoiceResponse struct {
	ID         uint   `json:"id"`
	ChoiceText string `json:"choice_text"`
	OrderIndex int    `json:"order_index"`
}

// AttemptCreateRequest represents the request to start a new quiz attempt
type AttemptCreateRequest struct {
	UserIdentifier *string `json:"user_identifier,omitempty"`
	UserName       *string `json:"user_name,omitempty" validate:"omitempty,min=2,max=100"`
}

// AttemptSubmitRequest represents the request to submit quiz answers
type AttemptSubmitRequest struct {
	Answers UserAnswers `json:"answers" validate:"required"`
}

// AttemptResponse represents attempt data for API responses
type AttemptResponse struct {
	ID               uint        `json:"id"`
	QuizID           uint        `json:"quiz_id"`
	UserName         *string     `json:"user_name,omitempty"`
	StartedAt        time.Time   `json:"started_at"`
	CompletedAt      *time.Time  `json:"completed_at,omitempty"`
	TotalScore       int         `json:"total_score"`
	MaxPossibleScore int         `json:"max_possible_score"`
	PercentageScore  float64     `json:"percentage_score"`
	TimeTakenSeconds *int        `json:"time_taken_seconds,omitempty"`
	Status           string      `json:"status"`
	Answers          UserAnswers `json:"answers,omitempty"`
}

// Admin DTOs for quiz management

// QuizCreateRequest represents the request to create a new quiz
type QuizCreateRequest struct {
	Title            string  `json:"title" validate:"required,min=3,max=200"`
	Slug             string  `json:"slug" validate:"required,min=3,max=200"`
	Description      *string `json:"description,omitempty"`
	TopicID          uint    `json:"topic_id" validate:"required"`
	DifficultyLevel  string  `json:"difficulty_level" validate:"oneof=easy medium hard"`
	TimeLimitMinutes int     `json:"time_limit_minutes" validate:"min=1,max=180"`
}

// QuizUpdateRequest represents the request to update a quiz
type QuizUpdateRequest struct {
	Title            *string `json:"title,omitempty" validate:"omitempty,min=3,max=200"`
	Description      *string `json:"description,omitempty"`
	DifficultyLevel  *string `json:"difficulty_level,omitempty" validate:"omitempty,oneof=easy medium hard"`
	TimeLimitMinutes *int    `json:"time_limit_minutes,omitempty" validate:"omitempty,min=1,max=180"`
	IsActive         *bool   `json:"is_active,omitempty"`
}

// QuestionCreateRequest represents the request to create a new question
type QuestionCreateRequest struct {
	QuestionText string                `json:"question_text" validate:"required,min=10"`
	QuestionType string                `json:"question_type" validate:"oneof=multiple_choice true_false text"`
	Points       int                   `json:"points" validate:"min=1,max=100"`
	Explanation  *string               `json:"explanation,omitempty"`
	OrderIndex   int                   `json:"order_index" validate:"required,min=1"`
	Choices      []ChoiceCreateRequest `json:"choices,omitempty" validate:"dive"`
}

// ChoiceCreateRequest represents the request to create a new choice
type ChoiceCreateRequest struct {
	ChoiceText  string  `json:"choice_text" validate:"required,min=1"`
	IsCorrect   bool    `json:"is_correct"`
	OrderIndex  int     `json:"order_index" validate:"required,min=1"`
	Explanation *string `json:"explanation,omitempty"`
}

// Pagination and filtering

// PaginationParams represents common pagination parameters
type PaginationParams struct {
	Page     int `json:"page" form:"page" validate:"min=1"`
	PageSize int `json:"page_size" form:"page_size" validate:"min=1,max=100"`
}

// QuizFilterParams represents filtering parameters for quiz queries
type QuizFilterParams struct {
	PaginationParams
	TopicID         *uint   `json:"topic_id,omitempty" form:"topic_id"`
	DifficultyLevel *string `json:"difficulty_level,omitempty" form:"difficulty_level" validate:"omitempty,oneof=easy medium hard"`
	IsActive        *bool   `json:"is_active,omitempty" form:"is_active"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

// TableName methods for custom table names (if needed)
func (Topic) TableName() string {
	return "topics"
}

func (Quiz) TableName() string {
	return "quizzes"
}

func (Question) TableName() string {
	return "questions"
}

func (Choice) TableName() string {
	return "choices"
}

func (Attempt) TableName() string {
	return "attempts"
}

// Helper methods

// CalculatePercentage calculates the percentage score
func (a *Attempt) CalculatePercentage() float64 {
	if a.MaxPossibleScore == 0 {
		return 0.0
	}
	return (float64(a.TotalScore) / float64(a.MaxPossibleScore)) * 100.0
}

// GetTimeTaken calculates time taken in seconds
func (a *Attempt) GetTimeTaken() int {
	if a.CompletedAt == nil {
		return int(time.Since(a.StartedAt).Seconds())
	}
	return int(a.CompletedAt.Sub(a.StartedAt).Seconds())
}

// ToResponse converts Attempt to AttemptResponse
func (a *Attempt) ToResponse(includeAnswers bool) AttemptResponse {
	response := AttemptResponse{
		ID:               a.ID,
		QuizID:           a.QuizID,
		UserName:         a.UserName,
		StartedAt:        a.StartedAt,
		CompletedAt:      a.CompletedAt,
		TotalScore:       a.TotalScore,
		MaxPossibleScore: a.MaxPossibleScore,
		PercentageScore:  a.PercentageScore,
		TimeTakenSeconds: a.TimeTakenSeconds,
		Status:           a.Status,
	}

	if includeAnswers {
		response.Answers = a.Answers
	}

	return response
}

// ToResponse converts Topic to TopicResponse
func (t *Topic) ToResponse() TopicResponse {
	return TopicResponse{
		ID:          t.ID,
		Name:        t.Name,
		Slug:        t.Slug,
		Description: t.Description,
		IconURL:     t.IconURL,
		QuizCount:   len(t.Quizzes),
	}
}

// ToListResponse converts Quiz to QuizListResponse
func (q *Quiz) ToListResponse() QuizListResponse {
	return QuizListResponse{
		ID:               q.ID,
		Title:            q.Title,
		Slug:             q.Slug,
		Description:      q.Description,
		DifficultyLevel:  q.DifficultyLevel,
		TimeLimitMinutes: q.TimeLimitMinutes,
		TotalQuestions:   q.TotalQuestions,
		TopicName:        q.Topic.Name,
		TopicSlug:        q.Topic.Slug,
	}
}

// ToDetailResponse converts Quiz to QuizDetailResponse
func (q *Quiz) ToDetailResponse() QuizDetailResponse {
	return QuizDetailResponse{
		ID:               q.ID,
		Title:            q.Title,
		Slug:             q.Slug,
		Description:      q.Description,
		DifficultyLevel:  q.DifficultyLevel,
		TimeLimitMinutes: q.TimeLimitMinutes,
		TotalQuestions:   q.TotalQuestions,
		Topic:            q.Topic.ToResponse(),
	}
}

// ToResponse converts Question to QuestionResponse (hides correct answers)
func (q *Question) ToResponse() QuestionResponse {
	choices := make([]ChoiceResponse, len(q.Choices))
	for i, choice := range q.Choices {
		choices[i] = ChoiceResponse{
			ID:         choice.ID,
			ChoiceText: choice.ChoiceText,
			OrderIndex: choice.OrderIndex,
		}
	}

	return QuestionResponse{
		ID:           q.ID,
		QuestionText: q.QuestionText,
		QuestionType: q.QuestionType,
		Points:       q.Points,
		OrderIndex:   q.OrderIndex,
		Choices:      choices,
	}
}
