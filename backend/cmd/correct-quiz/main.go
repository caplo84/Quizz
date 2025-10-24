package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/caplo84/quizz-backend/internal/services"
	"github.com/caplo84/quizz-backend/internal/services/datasources"
	"github.com/caplo84/quizz-backend/pkg/utils"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func main() {
	_ = godotenv.Load()

	quizSlug := flag.String("quiz-slug", "", "Specific quiz slug to correct (optional)")
	dryRun := flag.Bool("dry-run", true, "If true, no database changes are written")
	batchSize := flag.Int("batch-size", 100, "Number of questions per processing batch")
	threshold := flag.Float64("ai-confidence-threshold", 0.7, "Minimum confidence required to apply correction")
	verbose := flag.Bool("verbose", false, "Verbose logs")
	flag.Parse()

	cfg, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := connectDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	githubClient := datasources.NewGitHubClient(datasources.GitHubConfig{
		Token:      os.Getenv("GITHUB_TOKEN"),
		Owner:      "Ebazhanov",
		Repository: "linkedin-skill-assessments-quizzes",
		BaseURL:    "https://api.github.com",
	})
	aiService := services.NewAIAnswerServiceFromEnv()
	corrector := services.NewQuizCorrectorService(db, githubClient, aiService)

	opts := services.CorrectionOptions{
		QuizSlug:            *quizSlug,
		DryRun:              *dryRun,
		BatchSize:           *batchSize,
		ConfidenceThreshold: *threshold,
		Verbose:             *verbose,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Hour)
	defer cancel()

	report, err := corrector.CorrectAllQuizzes(ctx, opts)
	if err != nil {
		log.Fatalf("correction failed: %v", err)
	}

	if err := os.MkdirAll("logs", 0o755); err != nil {
		log.Printf("warning: failed to create logs directory: %v", err)
	}

	reportBytes, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal report: %v", err)
	}

	stamp := time.Now().UTC().Format("2006-01-02")
	reportPath := filepath.Join("logs", fmt.Sprintf("correction_report_%s.json", stamp))
	if err := os.WriteFile(reportPath, reportBytes, 0o644); err != nil {
		log.Fatalf("failed to write report file: %v", err)
	}

	logPath := filepath.Join("logs", fmt.Sprintf("quiz_correction_%s.log", stamp))
	if f, openErr := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644); openErr == nil {
		defer f.Close()
		_, _ = fmt.Fprintf(f, "%s correction complete: processed=%d fixed=%d skipped=%d failed=%d dry_run=%t\n",
			time.Now().UTC().Format(time.RFC3339),
			report.TotalProcessed,
			report.TotalFixed,
			report.TotalSkipped,
			report.TotalFailed,
			report.DryRun,
		)
	}

	fmt.Printf("Correction complete. Report written to %s\n", reportPath)
	fmt.Println(string(reportBytes))
}

func connectDB(cfg *utils.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	logger := gormLogger.Default.LogMode(gormLogger.Error)
	if cfg.Server.Environment == "development" {
		logger = gormLogger.Default.LogMode(gormLogger.Warn)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
