package repository

import (
    "context"
    "github.com/caplo84/quizz-backend/internal/models"
    "gorm.io/gorm"
)

// Implementation struct
type quizRepository struct {
    db *gorm.DB
}

func NewQuizRepository(db *gorm.DB) QuizRepository {
    return &quizRepository{db: db}
}

func (r *quizRepository) GetQuizzesByTopic(ctx context.Context, topicID uint) ([]models.Quiz, error) {
    var quizzes []models.Quiz
    if err := r.db.WithContext(ctx).Where("topic_id = ?", topicID).Find(&quizzes).Error; err != nil {
        return nil, err
    }
    return quizzes, nil
}

func (r *quizRepository) GetQuizBySlug(ctx context.Context, slug string) (*models.Quiz, error) {
    var quiz models.Quiz
    if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&quiz).Error; err != nil {
        return nil, err
    }
    return &quiz, nil
}

func (r *quizRepository) GetQuizQuestions(ctx context.Context, quizID uint) ([]models.Question, error) {
    var questions []models.Question
    if err := r.db.WithContext(ctx).Where("quiz_id = ?", quizID).Find(&questions).Error; err != nil {
        return nil, err
    }
    return questions, nil
}

func (r *quizRepository) CreateQuiz(ctx context.Context, quiz *models.Quiz) error {
    return r.db.WithContext(ctx).Create(quiz).Error
}

func (r *quizRepository) UpdateQuiz(ctx context.Context, quiz *models.Quiz) error {
    return r.db.WithContext(ctx).Save(quiz).Error
}

func (r *quizRepository) DeleteQuiz(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).Delete(&models.Quiz{}, id).Error
}

func (r *quizRepository) GetQuizByID(ctx context.Context, id uint) (*models.Quiz, error) {
    var quiz models.Quiz
    
    err := r.db.WithContext(ctx).
        Preload("Topic").
        Preload("Questions").
        Preload("Questions.Choices").
        Where("id = ? AND is_active = ?", id, true).
        First(&quiz).Error
    
    if err != nil {
        return nil, err
    }
    
    return &quiz, nil
}