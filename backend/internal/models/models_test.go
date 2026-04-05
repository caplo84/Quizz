package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func strPtr(s string) *string {
	return &s
}

// ---- UserAnswers tests ----

func TestUserAnswers_Value(t *testing.T) {
	ua := UserAnswers{
		"1": "choice_a",
		"2": true,
	}

	val, err := ua.Value()
	assert.NoError(t, err)
	assert.NotNil(t, val)

	// The driver.Value should be JSON bytes
	b, ok := val.([]byte)
	assert.True(t, ok)

	var decoded map[string]interface{}
	err = json.Unmarshal(b, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, "choice_a", decoded["1"])
}

func TestUserAnswers_Scan_Bytes(t *testing.T) {
	original := UserAnswers{"q1": "a", "q2": "b"}
	b, _ := json.Marshal(original)

	var ua UserAnswers
	err := ua.Scan(b)
	assert.NoError(t, err)
	assert.Equal(t, "a", ua["q1"])
	assert.Equal(t, "b", ua["q2"])
}

func TestUserAnswers_Scan_Nil(t *testing.T) {
	var ua UserAnswers
	err := ua.Scan(nil)
	assert.NoError(t, err)
	assert.NotNil(t, ua)
}

func TestUserAnswers_Scan_NonBytes(t *testing.T) {
	var ua UserAnswers
	err := ua.Scan("not bytes")
	assert.NoError(t, err)
}

// ---- Attempt helper methods ----

func TestAttempt_CalculatePercentage_Normal(t *testing.T) {
	a := &Attempt{
		TotalScore:       7,
		MaxPossibleScore: 10,
	}
	assert.InDelta(t, 70.0, a.CalculatePercentage(), 0.001)
}

func TestAttempt_CalculatePercentage_ZeroMax(t *testing.T) {
	a := &Attempt{
		TotalScore:       5,
		MaxPossibleScore: 0,
	}
	assert.Equal(t, 0.0, a.CalculatePercentage())
}

func TestAttempt_CalculatePercentage_Perfect(t *testing.T) {
	a := &Attempt{
		TotalScore:       10,
		MaxPossibleScore: 10,
	}
	assert.InDelta(t, 100.0, a.CalculatePercentage(), 0.001)
}

func TestAttempt_GetTimeTaken_InProgress(t *testing.T) {
	a := &Attempt{
		StartedAt:   time.Now().Add(-30 * time.Second),
		CompletedAt: nil,
	}
	taken := a.GetTimeTaken()
	assert.GreaterOrEqual(t, taken, 29)
	assert.LessOrEqual(t, taken, 32)
}

func TestAttempt_GetTimeTaken_Completed(t *testing.T) {
	start := time.Now().Add(-2 * time.Minute)
	completed := start.Add(90 * time.Second)
	a := &Attempt{
		StartedAt:   start,
		CompletedAt: &completed,
	}
	assert.Equal(t, 90, a.GetTimeTaken())
}

func TestAttempt_ToResponse_WithAnswers(t *testing.T) {
	now := time.Now()
	completed := now.Add(5 * time.Minute)
	seconds := 300
	a := &Attempt{
		ID:               1,
		QuizID:           2,
		UserName:         strPtr("Alice"),
		StartedAt:        now,
		CompletedAt:      &completed,
		TotalScore:       8,
		MaxPossibleScore: 10,
		PercentageScore:  80.0,
		TimeTakenSeconds: &seconds,
		Status:           "completed",
		Answers:          UserAnswers{"1": "a"},
	}

	resp := a.ToResponse(true)
	assert.Equal(t, uint(1), resp.ID)
	assert.Equal(t, uint(2), resp.QuizID)
	assert.Equal(t, strPtr("Alice"), resp.UserName)
	assert.Equal(t, 8, resp.TotalScore)
	assert.Equal(t, 10, resp.MaxPossibleScore)
	assert.InDelta(t, 80.0, resp.PercentageScore, 0.001)
	assert.Equal(t, "completed", resp.Status)
	assert.NotNil(t, resp.Answers)
}

func TestAttempt_ToResponse_WithoutAnswers(t *testing.T) {
	a := &Attempt{
		ID:      1,
		QuizID:  2,
		Answers: UserAnswers{"1": "a"},
	}
	resp := a.ToResponse(false)
	assert.Nil(t, resp.Answers)
}

// ---- Topic helper methods ----

func TestTopic_ToResponse(t *testing.T) {
	topic := &Topic{
		ID:          1,
		Name:        "Programming",
		Slug:        "programming",
		Description: strPtr("Prog topics"),
		IconURL:     strPtr("http://icon.url"),
		Quizzes:     []Quiz{{}, {}},
	}

	resp := topic.ToResponse()
	assert.Equal(t, uint(1), resp.ID)
	assert.Equal(t, "Programming", resp.Name)
	assert.Equal(t, "programming", resp.Slug)
	assert.Equal(t, strPtr("Prog topics"), resp.Description)
	assert.Equal(t, strPtr("http://icon.url"), resp.IconURL)
	assert.Equal(t, 2, resp.QuizCount)
}

func TestTopic_ToResponse_NoQuizzes(t *testing.T) {
	topic := &Topic{
		ID:   2,
		Name: "Science",
		Slug: "science",
	}
	resp := topic.ToResponse()
	assert.Equal(t, 0, resp.QuizCount)
}

// ---- Quiz helper methods ----

func TestQuiz_ToListResponse(t *testing.T) {
	desc := "A quiz"
	quiz := &Quiz{
		ID:               10,
		Title:            "Go Basics",
		Slug:             "go-basics",
		Description:      &desc,
		DifficultyLevel:  "easy",
		TimeLimitMinutes: 15,
		TotalQuestions:   5,
		Topic: Topic{
			Name: "Programming",
			Slug: "programming",
		},
	}

	resp := quiz.ToListResponse()
	assert.Equal(t, uint(10), resp.ID)
	assert.Equal(t, "Go Basics", resp.Title)
	assert.Equal(t, "go-basics", resp.Slug)
	assert.Equal(t, "easy", resp.DifficultyLevel)
	assert.Equal(t, 15, resp.TimeLimitMinutes)
	assert.Equal(t, 5, resp.TotalQuestions)
	assert.Equal(t, "Programming", resp.TopicName)
	assert.Equal(t, "programming", resp.TopicSlug)
}

func TestQuiz_ToDetailResponse(t *testing.T) {
	quiz := &Quiz{
		ID:              1,
		Title:           "Python Quiz",
		Slug:            "python-quiz",
		DifficultyLevel: "medium",
		Topic: Topic{
			ID:   3,
			Name: "Python",
			Slug: "python",
		},
	}

	resp := quiz.ToDetailResponse()
	assert.Equal(t, uint(1), resp.ID)
	assert.Equal(t, "Python Quiz", resp.Title)
	assert.Equal(t, uint(3), resp.Topic.ID)
	assert.Equal(t, "Python", resp.Topic.Name)
}

// ---- Question helper methods ----

func TestQuestion_ToResponse(t *testing.T) {
	q := &Question{
		ID:           5,
		QuestionText: "What is Go?",
		QuestionType: "multiple_choice",
		Points:       2,
		OrderIndex:   1,
		Choices: []Choice{
			{ID: 1, ChoiceText: "A language", IsCorrect: true, OrderIndex: 1},
			{ID: 2, ChoiceText: "A game", IsCorrect: false, OrderIndex: 2},
		},
	}

	resp := q.ToResponse()
	assert.Equal(t, uint(5), resp.ID)
	assert.Equal(t, "What is Go?", resp.QuestionText)
	assert.Equal(t, "multiple_choice", resp.QuestionType)
	assert.Equal(t, 2, resp.Points)
	assert.Len(t, resp.Choices, 2)
	// IsCorrect should NOT be in ChoiceResponse
	assert.Equal(t, uint(1), resp.Choices[0].ID)
	assert.Equal(t, "A language", resp.Choices[0].ChoiceText)
}

func TestQuestion_ToResponse_NoChoices(t *testing.T) {
	q := &Question{
		ID:           3,
		QuestionText: "True or false?",
		QuestionType: "true_false",
		Points:       1,
		OrderIndex:   2,
		Choices:      []Choice{},
	}
	resp := q.ToResponse()
	assert.Empty(t, resp.Choices)
}

// ---- TableName tests ----

func TestTableNames(t *testing.T) {
	assert.Equal(t, "topics", Topic{}.TableName())
	assert.Equal(t, "quizzes", Quiz{}.TableName())
	assert.Equal(t, "questions", Question{}.TableName())
	assert.Equal(t, "choices", Choice{}.TableName())
	assert.Equal(t, "attempts", Attempt{}.TableName())
}
