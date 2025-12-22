//go:build examples
// +build examples

package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nas-ai/api/src/config"
	"github.com/nas-ai/api/src/database"
	"github.com/nas-ai/api/src/logger"
	"github.com/nas-ai/api/src/models"
	"github.com/nas-ai/api/src/repository"
)

// Example demonstrating sqlx usage
func main() {
	// Setup
	log := logger.NewSlogLogger("info", "development")
	cfg, _ := config.LoadConfig()
	db, _ := database.NewPostgresConnectionX(cfg, log)
	defer db.Close()

	ctx := context.Background()
	userRepo := repository.NewUserRepositoryX(db, log)

	// Example 1: Single row query (GetContext)
	user, err := userRepo.FindByEmail(ctx, "user@example.com")
	if err != nil {
		log.Error("Failed to find user", slog.String("error", err.Error()))
		return
	}
	if user != nil {
		fmt.Printf("Found user: %s (%s)\n", user.Username, user.Email)
	}

	// Example 2: Batch query (SelectContext with IN)
	ids := []string{"user_1", "user_2", "user_3"}
	users, err := userRepo.FindByIDs(ctx, ids)
	if err != nil {
		log.Error("Failed to find users", slog.String("error", err.Error()))
		return
	}
	fmt.Printf("Found %d users\n", len(users))

	// Example 3: Named query (NamedExecContext)
	user = &models.User{
		ID:       "user_123",
		Username: "newusername",
		Email:    "newemail@example.com",
	}
	err = userRepo.UpdateUser(ctx, user)
	if err != nil {
		log.Error("Failed to update user", slog.String("error", err.Error()))
		return
	}
	fmt.Println("User updated successfully")

	// Example 4: Pagination (SelectContext)
	limit, offset := 10, 0
	paginatedUsers, err := userRepo.List(ctx, limit, offset)
	if err != nil {
		log.Error("Failed to list users", slog.String("error", err.Error()))
		return
	}
	fmt.Printf("Listed %d users (page 1)\n", len(paginatedUsers))

	// Example 5: Count query
	count, err := userRepo.Count(ctx)
	if err != nil {
		log.Error("Failed to count users", slog.String("error", err.Error()))
		return
	}
	fmt.Printf("Total users in database: %d\n", count)

	// Example 6: Transaction
	err = db.WithTransaction(ctx, func(tx *database.TxFunc) error {
		// Multiple operations in transaction
		// If any fails, all are rolled back
		user1 := &models.User{
			Username: "user1",
			Email:    "user1@example.com",
		}

		user2 := &models.User{
			Username: "user2",
			Email:    "user2@example.com",
		}

		// Both creates succeed or both fail
		if _, err := userRepo.CreateUser(ctx, user1.Username, user1.Email, "hash1"); err != nil {
			return err // Triggers rollback
		}

		if _, err := userRepo.CreateUser(ctx, user2.Username, user2.Email, "hash2"); err != nil {
			return err // Triggers rollback
		}

		return nil // Triggers commit
	})

	if err != nil {
		log.Error("Transaction failed", slog.String("error", err.Error()))
		return
	}

	fmt.Println("Transaction completed successfully")
}
