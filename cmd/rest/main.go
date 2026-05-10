package main

import (
	"log"

	"profconnect-api/internal/adapter/database"
	"profconnect-api/internal/adapter/handler"
	"profconnect-api/internal/adapter/routes"
	"profconnect-api/internal/config"
	"profconnect-api/internal/database/sqlc/generated"
	"profconnect-api/internal/infrastructure/adapter/repository"
	"profconnect-api/internal/infrastructure/adapter/service"
	"profconnect-api/internal/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func main() {
	// Load configuration
	config := config.LoadConfig()

	// Connect to database
	db, err := database.Connect(database.DBConfig{
		Host:     config.Host,
		Port:     config.Port,
		User:     config.User,
		Password: config.Password,
		DBName:   config.DBName,
		SSLMode:  config.SSLMode,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize SQLC queries
	queries := generated.New(db)

	// Initialize repositories (infrastructure adapters)
	userRepository := repository.NewUserRepository(queries)
	professorRepository := repository.NewProfessorRepository(queries)
	studentRepository := repository.NewStudentsRepository(queries)

	// Initialize services
	registerService := service.NewRegisterService(userRepository)

	// Initialize use cases
	registerAdminUsecase := usecase.NewRegisterAdminUsecase(registerService)
	registerProfessorUsecase := usecase.NewRegisterProfessorUsecase(registerService, professorRepository)
	registerStudentUsecase := usecase.NewRegisterStudentUsecase(registerService, studentRepository)
	loginUsecase := usecase.NewLoginUseCase(userRepository)
	professorProfileUsecase := usecase.NewProfessorProfileUsecase(professorRepository)
	studentProfileUsecase := usecase.NewStudentProfileUsecase(studentRepository)

	// Initialize handlers (presentation adapters)
	h := handler.NewHandler(
		registerAdminUsecase,
		registerProfessorUsecase,
		registerStudentUsecase,
		loginUsecase,
		professorProfileUsecase,
		studentProfileUsecase,
	)

	// Initialize Fiber app
	app := fiber.New()

	app.Use("/docs", static.New("./docs"))

	// Enable CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))

	// Initialize router and register routes
	router := routes.NewRouter(app, h)
	router.RegisterGeneralRoutes()
	router.RegisterSwaggerRoutes()

	// Start the server on port 3000
	app.Listen(":3001")
}
