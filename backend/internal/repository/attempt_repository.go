package repository

import (
    "context"
    "github.com/caplo84/quizz-backend/internal/models"
    "gorm.io/gorm"
)

type attemptRepository struct {
    db *gorm.DB
}

func NewAttemptRepository(db *gorm.DB) AttemptRepository {
    return &attemptRepository{db: db}
}

func (r *attemptRepository) CreateAttempt(ctx context.Context, attempt *models.Attempt) error {
    return r.db.WithContext(ctx).Create(attempt).Error
}

func (r *attemptRepository) GetAttemptByID(ctx context.Context, id uint) (*models.Attempt, error) {
    var attempt models.Attempt
    if err := r.db.WithContext(ctx).First(&attempt, id).Error; err != nil {
        return nil, err
    }
    return &attempt, nil
}

func (r *attemptRepository) UpdateAttemptAnswers(ctx context.Context, attempt *models.Attempt) error {
    return r.db.WithContext(ctx).Save(attempt).Error
}