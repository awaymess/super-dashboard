package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/awaymess/super-dashboard/backend/internal/config"
	"github.com/awaymess/super-dashboard/backend/pkg/jobs"
	"github.com/awaymess/super-dashboard/backend/pkg/logger"
)

func main() {
	// Initialize logger
	logger.Init()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	log.Info().
		Str("env", cfg.Env).
		Msg("Worker starting")

	// Create scheduler
	scheduler := jobs.NewScheduler()

	// Add default jobs
	for _, job := range jobs.CreateDefaultJobs() {
		if err := scheduler.AddJob(job); err != nil {
			log.Error().Err(err).Str("job", job.Name).Msg("Failed to add job")
			continue
		}
		log.Info().Str("job", job.Name).Str("cron", job.CronExpr).Msg("Job registered")
	}

	// Add daily jobs
	for _, job := range jobs.CreateDailyJobs() {
		if err := scheduler.AddJob(job); err != nil {
			log.Error().Err(err).Str("job", job.Name).Msg("Failed to add job")
			continue
		}
		log.Info().Str("job", job.Name).Str("cron", job.CronExpr).Msg("Job registered")
	}

	// Start scheduler
	scheduler.Start()
	log.Info().Int("job_count", scheduler.JobCount()).Msg("Worker started with scheduled jobs")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down worker...")

	// Stop scheduler
	scheduler.Stop()

	log.Info().Msg("Worker shutdown complete")
}
