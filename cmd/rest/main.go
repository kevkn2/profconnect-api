package main

import (
	"log"
	"time"

	authHandler "profconnect-api/internal/adapter/handler/auth"
	professorHandler "profconnect-api/internal/adapter/handler/professor"
	studentHandler "profconnect-api/internal/adapter/handler/student"
	"profconnect-api/internal/adapter/middleware"
	"profconnect-api/internal/adapter/router"
	"profconnect-api/internal/config"
	"profconnect-api/internal/database/sqlc/generated"
	"profconnect-api/internal/infrastructure/adapter/crypto"
	"profconnect-api/internal/infrastructure/adapter/postgres"
	"profconnect-api/internal/infrastructure/adapter/repository"
	"profconnect-api/internal/infrastructure/adapter/token"
	authUsecase "profconnect-api/internal/usecase/auth"
	professorUsecase "profconnect-api/internal/usecase/professor"
	studentUsecase "profconnect-api/internal/usecase/student"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/static"
)

const (
	accessTokenTTL  = 24 * time.Hour
	refreshTokenTTL = 7 * 24 * time.Hour
	tokenIssuer     = "profconnect-api"
)

func main() {
	cfg := config.LoadConfig()

	db, err := postgres.Connect(postgres.Config{
		Host:     cfg.Host,
		Port:     cfg.Port,
		User:     cfg.User,
		Password: cfg.Password,
		DBName:   cfg.DBName,
		SSLMode:  cfg.SSLMode,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	queries := generated.New(db)

	// Repositories
	userRepository := repository.NewUserRepository(queries)
	professorRepository := repository.NewProfessorRepository(queries)
	studentRepository := repository.NewStudentsRepository(queries)

	// Infrastructure services (port implementations)
	passwordHasher := crypto.NewBcryptHasher()
	tokenService := token.NewJWTService(cfg.JWTSecret, accessTokenTTL, refreshTokenTTL, tokenIssuer)

	// Auth use cases
	registerCoreUC := authUsecase.NewRegisterCoreUsecase(userRepository, passwordHasher)
	registerAdminUC := authUsecase.NewRegisterAdminUsecase(registerCoreUC)
	registerProfessorUC := authUsecase.NewRegisterProfessorUsecase(registerCoreUC, professorRepository)
	registerStudentUC := authUsecase.NewRegisterStudentUsecase(registerCoreUC, studentRepository)
	loginUC := authUsecase.NewLoginUsecase(userRepository, passwordHasher, tokenService)
	refreshUC := authUsecase.NewRefreshUsecase(userRepository, tokenService)

	// Role-specific use cases
	professorProfileUC := professorUsecase.NewProfileUsecase(professorRepository)
	studentProfileUC := studentUsecase.NewProfileUsecase(studentRepository)

	// Handlers
	authH := authHandler.New(
		registerAdminUC,
		registerProfessorUC,
		registerStudentUC,
		loginUC,
		refreshUC,
	)
	professorH := professorHandler.New(professorProfileUC)
	studentH := studentHandler.New(studentProfileUC)

	// Middleware
	authMW := middleware.NewAuth(tokenService)

	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${method} ${path} | ${latency} | ${ip}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	app.Use("/docs", static.New("./docs"))

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))

	r := router.NewRouter(app, authH, professorH, studentH, authMW)
	r.RegisterGeneralRoutes()
	r.RegisterSwaggerRoutes()
	r.RegisterAuthRoutes()
	r.RegisterStudentRoutes()
	r.RegisterProfessorRoutes()

	app.Listen(":3001")
}
