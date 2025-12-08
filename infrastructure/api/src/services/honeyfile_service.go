package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/nas-ai/api/src/repository"
	"github.com/sirupsen/logrus"
)

// HoneyfileService manages honeyfile intrusion detection with RAM caching
type HoneyfileService struct {
	repo      *repository.HoneyfileRepository
	encSvc    *EncryptionService
	logger    *logrus.Logger
	cache     map[string]bool
	cacheLock sync.RWMutex
}

// NewHoneyfileService creates a new service and loads cache
func NewHoneyfileService(repo *repository.HoneyfileRepository, encSvc *EncryptionService, logger *logrus.Logger) *HoneyfileService {
	s := &HoneyfileService{
		repo:   repo,
		encSvc: encSvc,
		logger: logger,
		cache:  make(map[string]bool),
	}

	// Load cache at startup
	if err := s.ReloadCache(context.Background()); err != nil {
		logger.WithError(err).Warn("Failed to load honeyfile cache at startup")
	}

	return s
}

// ReloadCache loads all honeyfile paths from DB into RAM cache
func (s *HoneyfileService) ReloadCache(ctx context.Context) error {
	paths, err := s.repo.GetAllPaths(ctx)
	if err != nil {
		return fmt.Errorf("reload cache: %w", err)
	}

	s.cacheLock.Lock()
	defer s.cacheLock.Unlock()

	s.cache = make(map[string]bool)
	for _, p := range paths {
		cleanPath := filepath.Clean(p)
		s.cache[cleanPath] = true
	}

	s.logger.WithField("count", len(s.cache)).Info("Honeyfile cache loaded")
	return nil
}

// IsHoneyfile checks if a path is in the honeyfile cache (fast RAM lookup)
func (s *HoneyfileService) IsHoneyfile(rawPath string) bool {
	cleanPath := filepath.Clean(rawPath)

	s.cacheLock.RLock()
	defer s.cacheLock.RUnlock()

	return s.cache[cleanPath]
}

// CheckAndTrigger checks if path is honeyfile and triggers lockdown if true
// Returns true if honeyfile was triggered (vault is now locked)
func (s *HoneyfileService) CheckAndTrigger(rawPath string) bool {
	cleanPath := filepath.Clean(rawPath)

	// Fast RAM check (no DB call!)
	s.cacheLock.RLock()
	isHoney := s.cache[cleanPath]
	s.cacheLock.RUnlock()

	if isHoney {
		// ALARM!
		s.logger.WithField("path", cleanPath).Error("🕷️ HONEYFILE ACCESSED - INITIATING LOCKDOWN")

		// Async DB update (don't block the panic)
		go func() {
			if err := s.repo.IncrementTrigger(context.Background(), cleanPath); err != nil {
				s.logger.WithError(err).Error("Failed to record honeyfile trigger")
			}
		}()

		// THE KILL SWITCH - Wipe all keys from RAM
		if err := s.encSvc.Lock(); err != nil {
			s.logger.WithError(err).Error("Failed to lock vault during honeyfile panic!")
		}

		return true
	}
	return false
}

// Create adds a new honeyfile with optional fake content generation
func (s *HoneyfileService) Create(ctx context.Context, rawPath, fileType string, createdBy *uuid.UUID) (*repository.Honeyfile, error) {
	cleanPath := filepath.Clean(rawPath)

	// 1. Create DB entry
	honeyfile, err := s.repo.Create(ctx, cleanPath, fileType, createdBy)
	if err != nil {
		return nil, err
	}

	// 2. Create physical file if it doesn't exist
	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		content := s.generateFakeContent(fileType)
		if err := os.WriteFile(cleanPath, content, 0644); err != nil {
			s.logger.WithError(err).Warn("Failed to create physical honeyfile")
			// Don't fail - DB entry is still valid
		} else {
			s.logger.WithField("path", cleanPath).Info("Physical honeyfile created with fake content")
		}
	}

	// 3. Refresh cache
	s.cacheLock.Lock()
	s.cache[cleanPath] = true
	s.cacheLock.Unlock()

	return honeyfile, nil
}

// Delete removes a honeyfile marker (does NOT delete physical file)
func (s *HoneyfileService) Delete(ctx context.Context, rawPath string) error {
	cleanPath := filepath.Clean(rawPath)

	if err := s.repo.Delete(ctx, cleanPath); err != nil {
		return err
	}

	// Update cache
	s.cacheLock.Lock()
	delete(s.cache, cleanPath)
	s.cacheLock.Unlock()

	return nil
}

// ListAll returns all honeyfiles
func (s *HoneyfileService) ListAll(ctx context.Context) ([]repository.Honeyfile, error) {
	return s.repo.ListAll(ctx)
}

// generateFakeContent creates convincing fake content based on file type
func (s *HoneyfileService) generateFakeContent(fileType string) []byte {
	switch fileType {
	case "finance":
		return []byte(`# Bitcoin Wallet Backup
# Generated: 2024-01-15
# WARNING: Keep this file secure!

Wallet Address: bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh
Private Key: L5BmPp4Ry9K8H4Q9TghPQbU4VYdnKMqFVWHTh6MZqz7kJ3wN8xYP

Recovery Seed Phrase:
1. abandon  2. ability  3. able     4. about
5. above    6. absent   7. absorb   8. abstract
9. absurd   10. abuse   11. access  12. accident

Balance: 2.4587 BTC
Last Transaction: 2024-01-14T18:32:00Z
`)

	case "it":
		return []byte(`# SSH Private Key - ROOT ACCESS
# Server: nas-prod-01.internal
# User: root

-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAACFwAAAAdz
c2gtcnNhAAAAAwEAAQAAAgEAzKHJyNw7vGTpN3hM0zLx9Ke55bT9n8Xz9Tf0P3k1
FAKE_KEY_DO_NOT_USE_FAKE_KEY_DO_NOT_USE_FAKE_KEY_DO_NOT_USE_FAKE_KEY
YmFzZTY0LWVuY29kZWQta2V5LWRhdGEtaGVyZS1mb3ItdGhlLWhvbmV5cG90LXRy
YXAtdGhpcy1pcy1ub3QtYS1yZWFsLXNzaC1rZXktaXQtaXMtYS1kZWNveS10cmFw
-----END OPENSSH PRIVATE KEY-----

# AWS Credentials (Production)
AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY

# Database Credentials
DB_HOST=mysql-prod.internal
DB_USER=admin
DB_PASS=Pr0d_Sup3r_S3cr3t_P4ss!
`)

	case "private":
		return []byte(`Meine geheimen Passwörter - NICHT TEILEN!
==========================================

Online Banking (Sparkasse):
  Benutzer: max.mustermann
  PIN: 84729
  TAN-Liste: Im Tresor

Amazon:
  Email: max.mustermann@gmail.com
  Passwort: MeinHund2019!

Netflix:
  max.mustermann@gmail.com / Netflix123

Facebook:
  max.mustermann@gmail.com / Sommer2020!

Router (FritzBox):
  Admin / fritzbox4ever

Haustür Code: 1234#5678
Tresor Kombination: 38-12-54
`)

	default:
		return []byte(`# Confidential Document
# Classification: RESTRICTED
# Do not distribute

This document contains sensitive information.
Access is logged and monitored.

[Content redacted for security]
`)
	}
}
