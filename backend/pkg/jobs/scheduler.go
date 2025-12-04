// Package jobs provides background job scheduling and execution using robfig/cron.
package jobs

import (
	"context"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

// Job represents a scheduled job with cron expression support.
type Job struct {
	Name     string
	CronExpr string // Cron expression (e.g., "*/30 * * * *" for every 30 minutes)
	Handler  func(ctx context.Context) error
	running  bool
	mu       sync.Mutex
}

// Scheduler manages background jobs using robfig/cron.
type Scheduler struct {
	cron    *cron.Cron
	jobs    []*Job
	ctx     context.Context
	cancel  context.CancelFunc
	running bool
	mu      sync.Mutex
}

// NewScheduler creates a new job scheduler with robfig/cron.
func NewScheduler() *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		cron:   cron.New(cron.WithSeconds()),
		jobs:   make([]*Job, 0),
		ctx:    ctx,
		cancel: cancel,
	}
}

// AddJob adds a job to the scheduler with a cron expression.
func (s *Scheduler) AddJob(job *Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create a wrapper function that handles context and concurrency
	wrappedHandler := s.createJobWrapper(job)

	_, err := s.cron.AddFunc(job.CronExpr, wrappedHandler)
	if err != nil {
		return err
	}

	s.jobs = append(s.jobs, job)
	return nil
}

// createJobWrapper creates a wrapper function for the job that handles
// context cancellation and prevents concurrent execution.
func (s *Scheduler) createJobWrapper(job *Job) func() {
	return func() {
		job.mu.Lock()
		if job.running {
			job.mu.Unlock()
			log.Warn().Str("job", job.Name).Msg("Job already running, skipping")
			return
		}
		job.running = true
		job.mu.Unlock()

		defer func() {
			job.mu.Lock()
			job.running = false
			job.mu.Unlock()
		}()

		// Check if scheduler context is cancelled
		select {
		case <-s.ctx.Done():
			return
		default:
		}

		start := time.Now()
		if err := job.Handler(s.ctx); err != nil {
			log.Error().Err(err).Str("job", job.Name).Msg("Job failed")
		} else {
			log.Debug().Str("job", job.Name).Dur("duration", time.Since(start)).Msg("Job completed")
		}
	}
}

// Start starts the scheduler.
func (s *Scheduler) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	log.Info().Int("job_count", len(s.jobs)).Msg("Starting job scheduler")
	s.cron.Start()
}

// Stop stops the scheduler gracefully.
func (s *Scheduler) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.running = false
	s.mu.Unlock()

	log.Info().Msg("Stopping job scheduler")
	s.cancel()
	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Info().Msg("Job scheduler stopped")
}

// IsRunning returns whether the scheduler is currently running.
func (s *Scheduler) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

// JobCount returns the number of registered jobs.
func (s *Scheduler) JobCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.jobs)
}

// GetJobs returns a copy of the registered jobs.
func (s *Scheduler) GetJobs() []*Job {
	s.mu.Lock()
	defer s.mu.Unlock()
	jobs := make([]*Job, len(s.jobs))
	copy(jobs, s.jobs)
	return jobs
}

// CreateDefaultJobs creates the default set of background jobs.
// These are stubs that log their execution - implement actual logic as needed.
// Uses standard cron expressions with seconds field: sec min hour day month weekday
func CreateDefaultJobs() []*Job {
	return []*Job{
		{
			Name:     "OddsSync",
			CronExpr: "0 */30 * * * *", // Every 30 minutes
			Handler:  oddsSyncHandler,
		},
		{
			Name:     "StockSync",
			CronExpr: "*/15 * * * * *", // Every 15 seconds
			Handler:  stockSyncHandler,
		},
		{
			Name:     "MatchStatusUpdate",
			CronExpr: "0 * * * * *", // Every minute
			Handler:  matchStatusUpdateHandler,
		},
		{
			Name:     "NewsSync",
			CronExpr: "0 */15 * * * *", // Every 15 minutes
			Handler:  newsSyncHandler,
		},
		{
			Name:     "SentimentAnalysis",
			CronExpr: "0 */30 * * * *", // Every 30 minutes
			Handler:  sentimentAnalysisHandler,
		},
		{
			Name:     "AlertChecker",
			CronExpr: "*/30 * * * * *", // Every 30 seconds
			Handler:  alertCheckerHandler,
		},
		{
			Name:     "ValueBetCalculator",
			CronExpr: "0 0 * * * *", // Every hour
			Handler:  valueBetCalculatorHandler,
		},
		{
			Name:     "AnalyticsAggregation",
			CronExpr: "0 0 * * * *", // Every hour
			Handler:  analyticsAggregationHandler,
		},
	}
}

// Job handlers (stubs)

func oddsSyncHandler(ctx context.Context) error {
	log.Info().Msg("OddsSync: Fetching odds from bookmakers")
	// TODO: Implement odds fetching from Pinnacle, Bet365, etc.
	// - Fetch odds for upcoming matches
	// - Update odds in database
	// - Broadcast updates via WebSocket
	return nil
}

func stockSyncHandler(ctx context.Context) error {
	log.Debug().Msg("StockSync: Fetching stock prices")
	// TODO: Implement stock price fetching
	// - Fetch prices from market data provider
	// - Update prices in database
	// - Broadcast updates via WebSocket
	return nil
}

func matchStatusUpdateHandler(ctx context.Context) error {
	log.Debug().Msg("MatchStatusUpdate: Checking live match statuses")
	// TODO: Implement match status updates
	// - Check for live matches
	// - Update scores and status
	// - Broadcast updates via WebSocket
	return nil
}

func newsSyncHandler(ctx context.Context) error {
	log.Info().Msg("NewsSync: Fetching news from sources")
	// TODO: Implement news fetching
	// - Fetch from Bloomberg, Reuters, etc.
	// - Parse and store articles
	// - Create embeddings for semantic search
	return nil
}

func sentimentAnalysisHandler(ctx context.Context) error {
	log.Info().Msg("SentimentAnalysis: Analyzing news sentiment")
	// TODO: Implement sentiment analysis
	// - Process unanalyzed articles
	// - Use OpenAI for sentiment detection
	// - Update sentiment scores in database
	return nil
}

func alertCheckerHandler(ctx context.Context) error {
	log.Debug().Msg("AlertChecker: Checking price alerts")
	// TODO: Implement alert checking
	// - Check price alerts against current prices
	// - Trigger notifications for met conditions
	// - Update alert status
	return nil
}

func valueBetCalculatorHandler(ctx context.Context) error {
	log.Info().Msg("ValueBetCalculator: Calculating value bets")
	// TODO: Implement value bet calculation
	// - Get current odds from database
	// - Calculate true probabilities using models
	// - Identify value bets
	// - Store in value_bets table
	return nil
}

func analyticsAggregationHandler(ctx context.Context) error {
	log.Info().Msg("AnalyticsAggregation: Aggregating analytics data")
	// TODO: Implement analytics aggregation
	// - Calculate ROI by league/team/bet type
	// - Update performance metrics
	// - Generate daily/weekly summaries
	return nil
}

// CreateDailyJobs returns jobs that should run once per day.
// Uses standard cron expressions with seconds field.
func CreateDailyJobs() []*Job {
	return []*Job{
		{
			Name:     "DailyPicks",
			CronExpr: "0 0 6 * * *", // Every day at 6:00 AM
			Handler:  dailyPicksHandler,
		},
		{
			Name:     "DataCleanup",
			CronExpr: "0 0 2 * * *", // Every day at 2:00 AM
			Handler:  dataCleanupHandler,
		},
		{
			Name:     "BackupJob",
			CronExpr: "0 0 3 * * *", // Every day at 3:00 AM
			Handler:  backupJobHandler,
		},
	}
}

func dailyPicksHandler(ctx context.Context) error {
	log.Info().Msg("DailyPicks: Generating daily betting picks")
	// TODO: Implement daily picks generation
	// - Analyze upcoming matches
	// - Generate top value bet recommendations
	// - Store in daily_picks table
	return nil
}

func dataCleanupHandler(ctx context.Context) error {
	log.Info().Msg("DataCleanup: Cleaning old data")
	// TODO: Implement data cleanup
	// - Remove old odds history
	// - Archive completed matches
	// - Clean temporary data
	return nil
}

func backupJobHandler(ctx context.Context) error {
	log.Info().Msg("BackupJob: Running database backup")
	// TODO: Implement database backup
	// - Create database dump
	// - Upload to cloud storage
	// - Rotate old backups
	return nil
}
