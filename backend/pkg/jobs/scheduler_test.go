package jobs

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewScheduler(t *testing.T) {
	scheduler := NewScheduler()
	if scheduler == nil {
		t.Fatal("NewScheduler returned nil")
	}
	if scheduler.cron == nil {
		t.Error("Expected cron instance to be initialized")
	}
	if scheduler.jobs == nil {
		t.Error("Expected jobs slice to be initialized")
	}
	if scheduler.ctx == nil {
		t.Error("Expected context to be initialized")
	}
	if scheduler.cancel == nil {
		t.Error("Expected cancel function to be initialized")
	}
}

func TestScheduler_AddJob(t *testing.T) {
	scheduler := NewScheduler()

	job := &Job{
		Name:     "TestJob",
		CronExpr: "* * * * * *", // Every second
		Handler: func(ctx context.Context) error {
			return nil
		},
	}

	err := scheduler.AddJob(job)
	if err != nil {
		t.Errorf("AddJob failed: %v", err)
	}

	if scheduler.JobCount() != 1 {
		t.Errorf("Expected 1 job, got %d", scheduler.JobCount())
	}
}

func TestScheduler_AddJob_InvalidCron(t *testing.T) {
	scheduler := NewScheduler()

	job := &Job{
		Name:     "InvalidJob",
		CronExpr: "invalid cron expression",
		Handler: func(ctx context.Context) error {
			return nil
		},
	}

	err := scheduler.AddJob(job)
	if err == nil {
		t.Error("Expected error for invalid cron expression")
	}
}

func TestScheduler_StartStop(t *testing.T) {
	scheduler := NewScheduler()

	job := &Job{
		Name:     "TestJob",
		CronExpr: "* * * * * *",
		Handler: func(ctx context.Context) error {
			return nil
		},
	}

	err := scheduler.AddJob(job)
	if err != nil {
		t.Fatalf("AddJob failed: %v", err)
	}

	scheduler.Start()

	if !scheduler.IsRunning() {
		t.Error("Expected scheduler to be running after Start")
	}

	// Start again should be no-op
	scheduler.Start()
	if !scheduler.IsRunning() {
		t.Error("Expected scheduler to still be running after second Start")
	}

	scheduler.Stop()

	if scheduler.IsRunning() {
		t.Error("Expected scheduler to be stopped after Stop")
	}

	// Stop again should be no-op
	scheduler.Stop()
	if scheduler.IsRunning() {
		t.Error("Expected scheduler to still be stopped after second Stop")
	}
}

func TestScheduler_JobExecution(t *testing.T) {
	scheduler := NewScheduler()

	var executionCount int32
	var mu sync.Mutex
	executed := false

	job := &Job{
		Name:     "ExecutionTestJob",
		CronExpr: "* * * * * *", // Every second
		Handler: func(ctx context.Context) error {
			atomic.AddInt32(&executionCount, 1)
			mu.Lock()
			executed = true
			mu.Unlock()
			return nil
		},
	}

	err := scheduler.AddJob(job)
	if err != nil {
		t.Fatalf("AddJob failed: %v", err)
	}

	scheduler.Start()
	defer scheduler.Stop()

	// Wait for at least one execution
	time.Sleep(2 * time.Second)

	mu.Lock()
	wasExecuted := executed
	mu.Unlock()

	if !wasExecuted {
		t.Error("Expected job to be executed at least once")
	}

	count := atomic.LoadInt32(&executionCount)
	if count < 1 {
		t.Errorf("Expected execution count >= 1, got %d", count)
	}
}

func TestScheduler_GetJobs(t *testing.T) {
	scheduler := NewScheduler()

	job1 := &Job{
		Name:     "Job1",
		CronExpr: "0 * * * * *",
		Handler: func(ctx context.Context) error {
			return nil
		},
	}
	job2 := &Job{
		Name:     "Job2",
		CronExpr: "0 0 * * * *",
		Handler: func(ctx context.Context) error {
			return nil
		},
	}

	scheduler.AddJob(job1)
	scheduler.AddJob(job2)

	jobs := scheduler.GetJobs()
	if len(jobs) != 2 {
		t.Errorf("Expected 2 jobs, got %d", len(jobs))
	}

	// Verify it's a copy
	jobs[0] = nil
	originalJobs := scheduler.GetJobs()
	if originalJobs[0] == nil {
		t.Error("GetJobs should return a copy, not the original slice")
	}
}

func TestScheduler_ConcurrentJobPrevention(t *testing.T) {
	scheduler := NewScheduler()

	var executionCount int32
	var concurrentExecutions int32
	var maxConcurrent int32

	job := &Job{
		Name:     "ConcurrentTestJob",
		CronExpr: "* * * * * *", // Every second
		Handler: func(ctx context.Context) error {
			current := atomic.AddInt32(&concurrentExecutions, 1)
			// Track max concurrent executions
			for {
				max := atomic.LoadInt32(&maxConcurrent)
				if current <= max {
					break
				}
				if atomic.CompareAndSwapInt32(&maxConcurrent, max, current) {
					break
				}
			}
			time.Sleep(500 * time.Millisecond) // Simulate work
			atomic.AddInt32(&concurrentExecutions, -1)
			atomic.AddInt32(&executionCount, 1)
			return nil
		},
	}

	err := scheduler.AddJob(job)
	if err != nil {
		t.Fatalf("AddJob failed: %v", err)
	}

	scheduler.Start()
	time.Sleep(3 * time.Second)
	scheduler.Stop()

	max := atomic.LoadInt32(&maxConcurrent)
	if max > 1 {
		t.Errorf("Expected max concurrent executions to be 1, got %d", max)
	}
}

func TestCreateDefaultJobs(t *testing.T) {
	jobs := CreateDefaultJobs()
	if len(jobs) == 0 {
		t.Error("Expected default jobs to be non-empty")
	}

	expectedJobs := []string{
		"OddsSync",
		"StockSync",
		"MatchStatusUpdate",
		"NewsSync",
		"SentimentAnalysis",
		"AlertChecker",
		"ValueBetCalculator",
		"AnalyticsAggregation",
	}

	for _, expected := range expectedJobs {
		found := false
		for _, job := range jobs {
			if job.Name == expected {
				found = true
				if job.CronExpr == "" {
					t.Errorf("Job %s has empty cron expression", expected)
				}
				if job.Handler == nil {
					t.Errorf("Job %s has nil handler", expected)
				}
				break
			}
		}
		if !found {
			t.Errorf("Expected job %s not found in default jobs", expected)
		}
	}
}

func TestCreateDailyJobs(t *testing.T) {
	jobs := CreateDailyJobs()
	if len(jobs) == 0 {
		t.Error("Expected daily jobs to be non-empty")
	}

	expectedJobs := []string{
		"DailyPicks",
		"DataCleanup",
		"BackupJob",
	}

	for _, expected := range expectedJobs {
		found := false
		for _, job := range jobs {
			if job.Name == expected {
				found = true
				if job.CronExpr == "" {
					t.Errorf("Job %s has empty cron expression", expected)
				}
				if job.Handler == nil {
					t.Errorf("Job %s has nil handler", expected)
				}
				break
			}
		}
		if !found {
			t.Errorf("Expected job %s not found in daily jobs", expected)
		}
	}
}

func TestJob_Subscribe(t *testing.T) {
	scheduler := NewScheduler()

	// Verify all default jobs can be added (valid cron expressions)
	for _, job := range CreateDefaultJobs() {
		err := scheduler.AddJob(job)
		if err != nil {
			t.Errorf("Failed to add default job %s: %v", job.Name, err)
		}
	}

	// Verify all daily jobs can be added (valid cron expressions)
	for _, job := range CreateDailyJobs() {
		err := scheduler.AddJob(job)
		if err != nil {
			t.Errorf("Failed to add daily job %s: %v", job.Name, err)
		}
	}

	expectedCount := len(CreateDefaultJobs()) + len(CreateDailyJobs())
	if scheduler.JobCount() != expectedCount {
		t.Errorf("Expected %d jobs, got %d", expectedCount, scheduler.JobCount())
	}
}
