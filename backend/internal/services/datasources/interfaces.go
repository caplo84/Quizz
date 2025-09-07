package datasources

import (
	"context"
	"time"
)

type DataSource interface {
	GetMetadata() SourceMetadata
	FetchCategories(ctx context.Context) ([]Category, error)
	FetchQuestions(ctx context.Context, params FetchParams) ([]Question, error)
	ValidateAuth(ctx context.Context) error
	GetRateLimit(ctx context.Context) (*RateLimit, error)
}

type SourceMetadata struct {
	Name         string
	Version      string
	RateLimited  int
	Capabilities []string
}

type FetchParams struct {
	Category  string
	Difficulty string
	Limit      int
	Offset     int
	Since      *time.Time // Use pointer for optional field
}

type RateLimit struct {
	Remaining int
	Reset     time.Time
	Total     int
}

type Question struct {
	Text              string
	Choices           []string
	Correct           int
	Category          string
	Source            string
	ExternalReference string
	ExternalID        string
}

type Category struct {
	Name        string
	Description string
	Count       int
}