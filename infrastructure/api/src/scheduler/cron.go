package scheduler

import (
	"fmt"
	"strings"
	"sync"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"

	"github.com/nas-ai/api/src/config"
	"github.com/nas-ai/api/src/services"
)

var (
	mu          sync.Mutex
	cronRunner  *cron.Cron
	backupSvc   *services.BackupService
	cfgRef      *config.Config
	logger      *logrus.Logger
	cronParser  = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	defaultSpec = "0 3 * * *"
)

// StartBackupScheduler starts the cron job that periodically creates and prunes backups.
func StartBackupScheduler(service *services.BackupService, cfg *config.Config) error {
	if service == nil {
		return fmt.Errorf("backup service is required")
	}
	if cfg == nil {
		return fmt.Errorf("config is required")
	}

	mu.Lock()
	defer mu.Unlock()

	backupSvc = service
	cfgRef = cfg
	logger = service.Logger()

	return startLocked()
}

// RestartScheduler restarts the scheduler using the currently configured service and config.
func RestartScheduler() error {
	mu.Lock()
	defer mu.Unlock()

	if backupSvc == nil || cfgRef == nil {
		return fmt.Errorf("backup scheduler not initialized")
	}

	return startLocked()
}

func startLocked() error {
	schedule := strings.TrimSpace(cfgRef.BackupSchedule)
	if schedule == "" {
		schedule = defaultSpec
	}

	if _, err := cronParser.Parse(schedule); err != nil {
		return fmt.Errorf("invalid backup schedule: %w", err)
	}

	if cronRunner != nil {
		ctx := cronRunner.Stop()
		<-ctx.Done()
	}

	cronRunner = cron.New(cron.WithParser(cronParser))

	job := func() {
		runBackupJob(backupSvc, cfgRef, logger)
	}

	if _, err := cronRunner.AddFunc(schedule, job); err != nil {
		return fmt.Errorf("register backup job: %w", err)
	}

	cronRunner.Start()

	if logger != nil {
		logger.WithField("schedule", schedule).Info("backup scheduler started")
	}

	return nil
}

func runBackupJob(svc *services.BackupService, cfg *config.Config, log *logrus.Logger) {
	if svc == nil || cfg == nil {
		return
	}

	if log != nil {
		log.WithFields(logrus.Fields{
			"schedule":  strings.TrimSpace(cfg.BackupSchedule),
			"retention": cfg.BackupRetentionCount,
			"path":      cfg.BackupStoragePath,
		}).Info("running scheduled backup")
	}

	if err := svc.SetBackupPath(cfg.BackupStoragePath); err != nil {
		if log != nil {
			log.WithError(err).Error("backup scheduler: failed to ensure backup path")
		}
		return
	}

	if _, err := svc.CreateBackup(cfg.BackupStoragePath); err != nil {
		if log != nil {
			log.WithError(err).Error("backup scheduler: failed to create backup")
		}
		return
	}

	if err := svc.PruneBackups(cfg.BackupRetentionCount); err != nil {
		if log != nil {
			log.WithError(err).Error("backup scheduler: failed to prune backups")
		}
		return
	}

	if log != nil {
		log.WithFields(logrus.Fields{
			"path":      cfg.BackupStoragePath,
			"retention": cfg.BackupRetentionCount,
		}).Info("backup scheduler: backup completed")
	}
}
