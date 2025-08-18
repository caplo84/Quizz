package services

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/caplo84/quizz-backend/internal/cache"
    "github.com/caplo84/quizz-backend/internal/logger"
    "github.com/caplo84/quizz-backend/internal/models"
    "github.com/caplo84/quizz-backend/internal/repository"
)

type QuizService interface {
    GetQuizzesByTopic(ctx context.Context, topicID uint) ([]models.Quiz, error)
    GetQuizBySlug(ctx context.Context, slug string) (*models.Quiz, error)
    GetQuizQuestions(ctx context.Context, quizID uint) ([]models.Question, error)
    CreateQuiz(ctx context.Context, quiz *models.Quiz) error
    UpdateQuiz(ctx context.Context, quiz *models.Quiz) error
    DeleteQuiz(ctx context.Context, id uint) error
}

type quizService struct {
    repo  repository.QuizRepository
    cache cache.Cache  // Add cache dependency
}

func NewQuizService(repo repository.QuizRepository, cache cache.Cache) QuizService {
    return &quizService{
        repo:  repo,
        cache: cache,  // Initialize cache
    }
}

func (s *quizService) GetQuizzesByTopic(ctx context.Context, topicID uint) ([]models.Quiz, error) {
    start := time.Now()
    
    logger.Log.WithContext(ctx).WithFields(logger.Fields{
        "operation": "get_quizzes_by_topic",
        "topic_id":  topicID,
    }).Debug("Starting quizzes retrieval by topic")
    
    quizzes, err := s.repo.GetQuizzesByTopic(ctx, topicID)
    if err != nil {
        logger.Log.WithContext(ctx).WithError(err).WithFields(logger.Fields{
            "operation": "get_quizzes_by_topic",
            "topic_id":  topicID,
        }).Error("Failed to retrieve quizzes by topic")
        return nil, err
    }
    
    logger.Log.WithContext(ctx).WithFields(logger.Fields{
        "operation":    "get_quizzes_by_topic",
        "topic_id":     topicID,
        "quiz_count":   len(quizzes),
        "duration_ms":  time.Since(start).Milliseconds(),
    }).Info("Quizzes retrieved successfully")
    
    return quizzes, nil
}

func (s *quizService) GetQuizBySlug(ctx context.Context, slug string) (*models.Quiz, error) {
    start := time.Now()

    logger.Log.WithContext(ctx).WithFields(logger.Fields{
        "operation": "get_quiz_by_slug",
        "slug":      slug,
    }).Debug("Starting quiz retrieval")

    // Try cache first
    cacheKey := fmt.Sprintf("quiz:%s", slug)
    if cachedQuiz, err := s.cache.Get(ctx, cacheKey); err == nil {
        logger.Log.LogCacheOperation(ctx, "get", cacheKey, true, time.Since(start))

        var quiz models.Quiz
        if err := json.Unmarshal(cachedQuiz, &quiz); err == nil {
            return &quiz, nil
        }
    }

    // Cache miss, get from database
    logger.Log.LogCacheOperation(ctx, "get", cacheKey, false, time.Since(start))

    quiz, err := s.repo.GetQuizBySlug(ctx, slug)
    if err != nil {
        logger.Log.WithContext(ctx).WithError(err).WithFields(logger.Fields{
            "operation": "get_quiz_by_slug",
            "slug":      slug,
        }).Error("Failed to retrieve quiz from database")
        return nil, err
    }

    // Cache the result
    if quizData, err := json.Marshal(quiz); err == nil {
        s.cache.Set(ctx, cacheKey, quizData, 15*time.Minute)
    }

    logger.Log.WithContext(ctx).WithFields(logger.Fields{
        "operation":   "get_quiz_by_slug",
        "slug":        slug,
        "quiz_id":     quiz.ID,
        "duration_ms": time.Since(start).Milliseconds(),
    }).Info("Quiz retrieved successfully")

    return quiz, nil
}

func (s *quizService) GetQuizQuestions(ctx context.Context, quizID uint) ([]models.Question, error) {
    start := time.Now()
    
    logger.Log.WithContext(ctx).WithFields(logger.Fields{
        "operation": "get_quiz_questions",
        "quiz_id":   quizID,
    }).Debug("Starting quiz questions retrieval")
    
    questions, err := s.repo.GetQuizQuestions(ctx, quizID)
    if err != nil {
        logger.Log.WithContext(ctx).WithError(err).WithFields(logger.Fields{
            "operation": "get_quiz_questions",
            "quiz_id":   quizID,
        }).Error("Failed to retrieve quiz questions")
        return nil, err
    }
    
    logger.Log.WithContext(ctx).WithFields(logger.Fields{
        "operation":      "get_quiz_questions",
        "quiz_id":        quizID,
        "question_count": len(questions),
        "duration_ms":    time.Since(start).Milliseconds(),
    }).Info("Quiz questions retrieved successfully")
    
    return questions, nil
}

func (s *quizService) CreateQuiz(ctx context.Context, quiz *models.Quiz) error {
    start := time.Now()
    
    logger.Log.WithContext(ctx).WithFields(logger.Fields{
        "operation": "create_quiz",
        "quiz_title": quiz.Title,
        "quiz_slug":  quiz.Slug,
    }).Info("Creating new quiz")
    
    err := s.repo.CreateQuiz(ctx, quiz)
    if err != nil {
        logger.Log.WithContext(ctx).WithError(err).WithFields(logger.Fields{
            "operation": "create_quiz",
            "quiz_title": quiz.Title,
        }).Error("Failed to create quiz")
        return err
    }
    
    logger.Log.WithContext(ctx).WithFields(logger.Fields{
        "operation":   "create_quiz",
        "quiz_id":     quiz.ID,
        "quiz_title":  quiz.Title,
        "duration_ms": time.Since(start).Milliseconds(),
    }).Info("Quiz created successfully")
    
    return nil
}

func (s *quizService) UpdateQuiz(ctx context.Context, quiz *models.Quiz) error {
    start := time.Now()
    
    logger.Log.WithContext(ctx).WithFields(logger.Fields{
        "operation": "update_quiz",
        "quiz_id":   quiz.ID,
        "quiz_title": quiz.Title,
    }).Info("Updating quiz")
    
    err := s.repo.UpdateQuiz(ctx, quiz)
    if err != nil {
        logger.Log.WithContext(ctx).WithError(err).WithFields(logger.Fields{
            "operation": "update_quiz",
            "quiz_id":   quiz.ID,
        }).Error("Failed to update quiz")
        return err
    }
    
    // Invalidate cache
    cacheKey := fmt.Sprintf("quiz:%s", quiz.Slug)
    s.cache.Delete(ctx, cacheKey)
    
    logger.Log.WithContext(ctx).WithFields(logger.Fields{
        "operation":   "update_quiz",
        "quiz_id":     quiz.ID,
        "duration_ms": time.Since(start).Milliseconds(),
    }).Info("Quiz updated successfully")
    
    return nil
}

func (s *quizService) DeleteQuiz(ctx context.Context, id uint) error {
    start := time.Now()
    
    logger.Log.WithContext(ctx).WithFields(logger.Fields{
        "operation": "delete_quiz",
        "quiz_id":   id,
    }).Info("Deleting quiz")
    
    // Get quiz first to invalidate cache
    quiz, _ := s.repo.GetQuizByID(ctx, id)
    
    err := s.repo.DeleteQuiz(ctx, id)
    if err != nil {
        logger.Log.WithContext(ctx).WithError(err).WithFields(logger.Fields{
            "operation": "delete_quiz",
            "quiz_id":   id,
        }).Error("Failed to delete quiz")
        return err
    }
    
    // Invalidate cache if we got the quiz
    if quiz != nil {
        cacheKey := fmt.Sprintf("quiz:%s", quiz.Slug)
        s.cache.Delete(ctx, cacheKey)
    }
    
    logger.Log.WithContext(ctx).WithFields(logger.Fields{
        "operation":   "delete_quiz",
        "quiz_id":     id,
        "duration_ms": time.Since(start).Milliseconds(),
    }).Info("Quiz deleted successfully")
    
    return nil
}