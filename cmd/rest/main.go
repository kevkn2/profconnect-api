package main

import (
	"log"

	"profconnect-api/internal/adapter/database"
	authHandler "profconnect-api/internal/adapter/handler/auth"
	professorHandler "profconnect-api/internal/adapter/handler/professor"
	studentHandler "profconnect-api/internal/adapter/handler/student"
	"profconnect-api/internal/adapter/routes"
	"profconnect-api/internal/config"
	"profconnect-api/internal/database/sqlc/generated"
	"profconnect-api/internal/infrastructure/adapter/repository"
	"profconnect-api/internal/infrastructure/adapter/service"
	authUsecase "profconnect-api/internal/usecase/auth"
	professorUsecase "profconnect-api/internal/usecase/professor"
	studentUsecase "profconnect-api/internal/usecase/student"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.Connect(database.DBConfig{
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

	// Services
	registerService := service.NewRegisterService(userRepository)

	// Auth use cases
	registerAdminUC := authUsecase.NewRegisterAdminUsecase(registerService)
	registerProfessorUC := authUsecase.NewRegisterProfessorUsecase(registerService, professorRepository)
	registerStudentUC := authUsecase.NewRegisterStudentUsecase(registerService, studentRepository)
	loginUC := authUsecase.NewLoginUsecase(userRepository)
	refreshUC := authUsecase.NewRefreshUsecase(userRepository)

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

	router := routes.NewRouter(app, authH, professorH, studentH)
	router.RegisterGeneralRoutes()
	router.RegisterSwaggerRoutes()
	router.RegisterAuthRoutes()
	router.RegisterStudentRoutes()
	router.RegisterProfessorRoutes()

	app.Listen(":3001")
}
