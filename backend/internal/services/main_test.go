package services

import (
	"os"
	"testing"

	"github.com/caplo84/quizz-backend/internal/logger"
)

func TestMain(m *testing.M) {
	// Initialize logger so service code that calls logger.Log doesn't panic
	_ = logger.InitializeLogger("error", "text")
	os.Exit(m.Run())
}
