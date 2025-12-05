// Package workers provides background worker implementations for the Super Dashboard.
package workers

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
)

// BackupWorker performs periodic database backups.
type BackupWorker struct {
	interval   time.Duration
	log        zerolog.Logger
	backupPath string
	dbHost     string
	dbPort     string
	dbName     string
	dbUser     string
	dbPassword string
}

// NewBackupWorker creates a new BackupWorker.
func NewBackupWorker(
	interval time.Duration,
	log zerolog.Logger,
	backupPath string,
	dbHost, dbPort, dbName, dbUser, dbPassword string,
) *BackupWorker {
	return &BackupWorker{
		interval:   interval,
		log:        log.With().Str("worker", "backup").Logger(),
		backupPath: backupPath,
		dbHost:     dbHost,
		dbPort:     dbPort,
		dbName:     dbName,
		dbUser:     dbUser,
		dbPassword: dbPassword,
	}
}

// StartBackup starts the backup worker.
func StartBackup(
	ctx context.Context,
	log zerolog.Logger,
	backupPath string,
	dbHost, dbPort, dbName, dbUser, dbPassword string,
) {
	worker := NewBackupWorker(24*time.Hour, log, backupPath, dbHost, dbPort, dbName, dbUser, dbPassword)
	worker.Run(ctx)
}

// Run starts the worker loop.
func (w *BackupWorker) Run(ctx context.Context) {
	w.log.Info().Dur("interval", w.interval).Msg("Starting backup worker")

	// Ensure backup directory exists
	if err := os.MkdirAll(w.backupPath, 0755); err != nil {
		w.log.Error().Err(err).Msg("Failed to create backup directory")
		return
	}

	// Schedule to run at 04:00 daily
	w.runAtScheduledTime(ctx)
}

// runAtScheduledTime runs the worker at a specific time each day.
func (w *BackupWorker) runAtScheduledTime(ctx context.Context) {
	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), 4, 0, 0, 0, now.Location())

		// If it's past 04:00 today, schedule for tomorrow
		if now.After(next) {
			next = next.Add(24 * time.Hour)
		}

		duration := next.Sub(now)
		w.log.Info().
			Time("next_run", next).
			Dur("wait", duration).
			Msg("Backup scheduled")

		select {
		case <-ctx.Done():
			w.log.Info().Msg("Backup worker stopping")
			return
		case <-time.After(duration):
			w.backup(ctx)
		}
	}
}

// backup creates a database backup.
func (w *BackupWorker) backup(ctx context.Context) {
	startTime := time.Now()
	w.log.Info().Msg("Starting database backup")

	// Generate backup filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	backupFile := filepath.Join(w.backupPath, fmt.Sprintf("super_dashboard_%s.sql", timestamp))

	// Use pg_dump to create backup
	cmd := exec.CommandContext(ctx, "pg_dump",
		"-h", w.dbHost,
		"-p", w.dbPort,
		"-U", w.dbUser,
		"-d", w.dbName,
		"-f", backupFile,
		"-F", "p", // Plain SQL format
		"-v",      // Verbose
	)

	// Set password via environment variable
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", w.dbPassword))

	// Run backup
	output, err := cmd.CombinedOutput()
	if err != nil {
		w.log.Error().
			Err(err).
			Str("output", string(output)).
			Msg("Database backup failed")
		return
	}

	// Get file size
	fileInfo, err := os.Stat(backupFile)
	if err != nil {
		w.log.Error().Err(err).Msg("Failed to get backup file info")
	}

	duration := time.Since(startTime)
	w.log.Info().
		Str("file", backupFile).
		Int64("size_bytes", fileInfo.Size()).
		Dur("duration", duration).
		Msg("Database backup completed")

	// Compress backup
	w.compressBackup(ctx, backupFile)

	// Clean up old backups (keep last 7 days)
	w.cleanupOldBackups(ctx)
}

// compressBackup compresses the backup file using gzip.
func (w *BackupWorker) compressBackup(ctx context.Context, backupFile string) {
	w.log.Debug().Str("file", backupFile).Msg("Compressing backup")

	cmd := exec.CommandContext(ctx, "gzip", "-f", backupFile)
	if err := cmd.Run(); err != nil {
		w.log.Error().Err(err).Msg("Failed to compress backup")
		return
	}

	w.log.Info().Str("file", backupFile+".gz").Msg("Backup compressed")
}

// cleanupOldBackups removes backups older than retention period.
func (w *BackupWorker) cleanupOldBackups(ctx context.Context) {
	w.log.Debug().Msg("Cleaning up old backups")

	retentionPeriod := 7 * 24 * time.Hour // Keep 7 days
	cutoff := time.Now().Add(-retentionPeriod)

	files, err := os.ReadDir(w.backupPath)
	if err != nil {
		w.log.Error().Err(err).Msg("Failed to read backup directory")
		return
	}

	deletedCount := 0
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileInfo, err := file.Info()
		if err != nil {
			continue
		}

		if fileInfo.ModTime().Before(cutoff) {
			filePath := filepath.Join(w.backupPath, file.Name())
			if err := os.Remove(filePath); err != nil {
				w.log.Error().
					Err(err).
					Str("file", filePath).
					Msg("Failed to delete old backup")
			} else {
				deletedCount++
				w.log.Debug().
					Str("file", filePath).
					Msg("Deleted old backup")
			}
		}
	}

	if deletedCount > 0 {
		w.log.Info().Int("deleted", deletedCount).Msg("Cleaned up old backups")
	}
}
