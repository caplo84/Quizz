package repository

import (
	"context"
	"errors"
	"github.com/caplo84/quizz-backend/internal/models"
	"gorm.io/gorm"
)

var errTopicRepositoryDBUnavailable = errors.New("topic repository: database connection is not initialized")

type topicRepository struct {
	db *gorm.DB
}

func NewTopicRepository(db *gorm.DB) TopicRepository {
	return &topicRepository{db: db}
}

func (r *topicRepository) GetAllTopics(ctx context.Context) ([]models.Topic, error) {
	if r == nil || r.db == nil {
		return nil, errTopicRepositoryDBUnavailable
	}

	var topics []models.Topic
	if err := r.db.WithContext(ctx).
		Table("topics").
		Select(`
			topics.*,
			COUNT(DISTINCT quizzes.id) AS active_quiz_count,
			COUNT(questions.id) AS active_question_count
		`).
		Joins("LEFT JOIN quizzes ON quizzes.topic_id = topics.id AND quizzes.is_active = ?", true).
		Joins("LEFT JOIN questions ON questions.quiz_id = quizzes.id AND questions.is_active = ?", true).
		Group("topics.id").
		Order("topics.name ASC").
		Scan(&topics).Error; err != nil {
		return nil, err
	}
	return topics, nil
}

func (r *topicRepository) GetTopicByID(ctx context.Context, id uint) (*models.Topic, error) {
	if r == nil || r.db == nil {
		return nil, errTopicRepositoryDBUnavailable
	}

	var topic models.Topic
	if err := r.db.WithContext(ctx).First(&topic, id).Error; err != nil {
		return nil, err
	}
	return &topic, nil
}

func (r *topicRepository) CreateTopic(ctx context.Context, topic *models.Topic) error {
	if r == nil || r.db == nil {
		return errTopicRepositoryDBUnavailable
	}

	return r.db.WithContext(ctx).Create(topic).Error
}

func (r *topicRepository) GetBySlug(ctx context.Context, slug string) (*models.Topic, error) {
	if r == nil || r.db == nil {
		return nil, errTopicRepositoryDBUnavailable
	}

	var topic models.Topic
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&topic).Error; err != nil {
		return nil, err
	}
	return &topic, nil
}

func (r *topicRepository) Create(ctx context.Context, topic *models.Topic) error {
	if r == nil || r.db == nil {
		return errTopicRepositoryDBUnavailable
	}

	return r.db.WithContext(ctx).Create(topic).Error
}

func (r *topicRepository) UpdateTopic(ctx context.Context, topic *models.Topic) error {
	if r == nil || r.db == nil {
		return errTopicRepositoryDBUnavailable
	}

	return r.db.WithContext(ctx).Save(topic).Error
}

func (r *topicRepository) DeleteTopic(ctx context.Context, id uint) error {
	if r == nil || r.db == nil {
		return errTopicRepositoryDBUnavailable
	}

	return r.db.WithContext(ctx).Delete(&models.Topic{}, id).Error
}
