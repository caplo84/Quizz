package services

import (
	"context"
	"fmt"
	"log"
	"github.com/caplo84/quizz-backend/internal/cache"
	"github.com/caplo84/quizz-backend/internal/models"
	"github.com/caplo84/quizz-backend/internal/repository"
	"github.com/caplo84/quizz-backend/internal/services/datasources"
)

type AdminService interface {
	CreateQuiz(ctx context.Context, quiz *models.Quiz) error
	UpdateQuiz(ctx context.Context, quiz *models.Quiz) error
	DeleteQuiz(ctx context.Context, id uint) error
	DownloadAllTopicImages(ctx context.Context) error
}

type adminService struct {
	quizRepo     repository.QuizRepository
	cache        cache.Cache
	githubClient *datasources.GitHubClient
	topicRepo    repository.TopicRepository
}

func NewAdminService(quizRepo repository.QuizRepository, cache cache.Cache, githubClient *datasources.GitHubClient, topicRepo repository.TopicRepository) AdminService {
	return &adminService{
		quizRepo:     quizRepo,
		cache:        cache,
		githubClient: githubClient,
		topicRepo:    topicRepo,
	}
}

func (s *adminService) CreateQuiz(ctx context.Context, quiz *models.Quiz) error {
	return s.quizRepo.CreateQuiz(ctx, quiz)
}

func (s *adminService) UpdateQuiz(ctx context.Context, quiz *models.Quiz) error {
	return s.quizRepo.UpdateQuiz(ctx, quiz)
}

func (s *adminService) DeleteQuiz(ctx context.Context, id uint) error {
	return s.quizRepo.DeleteQuiz(ctx, id)
}

func (s *adminService) DownloadAllTopicImages(ctx context.Context) error {
	log.Printf("🚀 Starting download all topic images from repository...")
	
	// Get all topics from database
	allTopics, err := s.topicRepo.GetAllTopics(ctx)
	if err != nil {
		return fmt.Errorf("failed to get topics from database: %w", err)
	}
	
	log.Printf("📊 Found %d topics in database", len(allTopics))
	
	downloadCount := 0
	errorCount := 0
	
	for _, topicModel := range allTopics {
		topic := topicModel.Slug
		log.Printf("🔍 Processing topic: %s", topic)
		
		// Try different common image folder patterns
		imageFolders := []string{
			fmt.Sprintf("%s/images", topic),
			fmt.Sprintf("%s/img", topic),
			fmt.Sprintf("%s/image", topic),
			fmt.Sprintf("%s/assets", topic),
		}
		
		for _, folder := range imageFolders {
			log.Printf("📂 Checking folder: %s", folder)
			
			// Get all image files from this folder
			imageFiles, err := s.githubClient.ListImageFiles(ctx, folder)
			if err != nil {
				log.Printf("⚠️  Could not access folder %s: %v", folder, err)
				continue
			}
			
			if len(imageFiles) == 0 {
				log.Printf("📭 No images found in %s", folder)
				continue
			}
			
			log.Printf("🎯 Found %d images in %s", len(imageFiles), folder)
			
			// Download each image
			for _, imageFile := range imageFiles {
				fullPath := fmt.Sprintf("%s/%s", folder, imageFile)
				log.Printf("📥 Downloading: %s", fullPath)
				
				if standardizedFilename, err := s.githubClient.DownloadImage(ctx, fullPath, topic); err != nil {
					log.Printf("⚠️  Failed to download %s: %v", fullPath, err)
					errorCount++
				} else {
					log.Printf("✅ Downloaded: %s -> %s", imageFile, standardizedFilename)
					downloadCount++
				}
			}
			
			// If found images in this folder, don't try other folder patterns
			if len(imageFiles) > 0 {
				break
			}
		}
	}
	
	log.Printf("🎉 Topic images download completed! Downloaded: %d, Errors: %d", downloadCount, errorCount)
	return nil
}