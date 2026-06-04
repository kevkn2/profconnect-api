package main

import (
	"context"
	"log"

	"profconnect-api/internal/config"
	"profconnect-api/internal/database/sqlc/generated"
	inputoutput "profconnect-api/internal/domain/input_output"
	"profconnect-api/internal/infrastructure/adapter/crypto"
	"profconnect-api/internal/infrastructure/adapter/postgres"
	"profconnect-api/internal/infrastructure/adapter/repository"
	authUsecase "profconnect-api/internal/usecase/auth"
)

func main() {
	cfg := config.LoadConfig()
	if cfg.SeedAdminName == "" || cfg.SeedAdminEmail == "" || cfg.SeedAdminPassword == "" {
		log.Fatalf("missing seeder config: SEED_ADMIN_NAME, SEED_ADMIN_EMAIL, and SEED_ADMIN_PASSWORD are required")
	}

	db, err := postgres.Connect(postgres.Config{
		Host:     cfg.Host,
		Port:     cfg.Port,
		User:     cfg.User,
		Password: cfg.Password,
		DBName:   cfg.DBName,
		SSLMode:  cfg.SSLMode,
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	queries := generated.New(db)
	userRepository := repository.NewUserRepository(queries)
	registerCoreUC := authUsecase.NewRegisterCoreUsecase(userRepository, crypto.NewBcryptHasher())
	registerAdminUC := authUsecase.NewRegisterAdminUsecase(registerCoreUC)
	ctx := context.Background()

	created, err := registerAdminUC.Execute(ctx, &inputoutput.RegisterInput{
		Name:     cfg.SeedAdminName,
		Email:    cfg.SeedAdminEmail,
		Password: cfg.SeedAdminPassword,
	})
	if err != nil {
		log.Fatalf("failed to create admin: %v", err)
	}

	log.Printf("admin seeded successfully: user_id=%s", created.UserID)
}
