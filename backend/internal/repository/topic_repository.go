package repository

import (
	"context"
	"github.com/caplo84/quizz-backend/internal/models"
	"gorm.io/gorm"
)

type topicRepository struct {
	db *gorm.DB
}

func NewTopicRepository(db *gorm.DB) TopicRepository {
	return &topicRepository{db: db}
}

func (r *topicRepository) GetAllTopics(ctx context.Context) ([]models.Topic, error) {
	var topics []models.Topic
	if err := r.db.WithContext(ctx).Find(&topics).Error; err != nil {
		return nil, err
	}
	return topics, nil
}

func (r *topicRepository) GetTopicByID(ctx context.Context, id uint) (*models.Topic, error) {
	var topic models.Topic
	if err := r.db.WithContext(ctx).First(&topic, id).Error; err != nil {
		return nil, err
	}
	return &topic, nil
}

func (r *topicRepository) CreateTopic(ctx context.Context, topic *models.Topic) error {
	return r.db.WithContext(ctx).Create(topic).Error
}

func (r *topicRepository) GetBySlug(ctx context.Context, slug string) (*models.Topic, error) {
	var topic models.Topic
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&topic).Error; err != nil {
		return nil, err
	}
	return &topic, nil
}

func (r *topicRepository) Create(ctx context.Context, topic *models.Topic) error {
	return r.db.WithContext(ctx).Create(topic).Error
}

func (r *topicRepository) UpdateTopic(ctx context.Context, topic *models.Topic) error {
	return r.db.WithContext(ctx).Save(topic).Error
}

func (r *topicRepository) DeleteTopic(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Topic{}, id).Error
}
