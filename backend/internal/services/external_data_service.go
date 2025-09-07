package services

import (
	"context"
	"fmt"
	
	"github.com/caplo84/quizz-backend/internal/services/datasources"
	"github.com/caplo84/quizz-backend/internal/repository"
)

type NewExternalDataService(
	githubClient *datasources.GitHubClient,
	quizRepo repository.QuizRepository,
) *ExternalDataService {
	return &ExternalDataService{
		githubClient: githubClient,
		quizRepo:	 quizRepo,
	}
}

func (s *ExternalDataService) SyncFromGithub(ctx context.Context) error {
	// Check rate limit
	rateLimit, err := s.githubClient.GetRateLimit(ctx)
	if err != nil {
		return fmt.Errorf("failed to check rate limit: %w", err)
	}

	if rateLimit.Remainint < 100 {
		return fmt.Errorf("rate limit too low: %d remaining", rateLimit.Remaining)
	}

	// Fetch and process questions
	question, err := s.githubClient.FetchQuestions(ctx, datasources.FetchParms{
		Limit: 100,
	})
	if err != nil {
		return fmt.Errorf("failed to fetch questions: %w", err)
	}

	// Save to database
	for _, question := range questions {
		// Convert to your internal model and save
        // Implementation depends on your repository interface
    }

	return nil
}