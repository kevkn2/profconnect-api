package main

import (
	"log"
	"time"

	authHandler "profconnect-api/internal/adapter/handler/auth"
	professorHandler "profconnect-api/internal/adapter/handler/professor"
	projectHandler "profconnect-api/internal/adapter/handler/project"
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
	projectUsecase "profconnect-api/internal/usecase/project"
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
	projectRepository := repository.NewProjectRepository(queries)
	projectApplicationRepository := repository.NewProjectApplicationRepository(queries)
	projectMemberRepository := repository.NewProjectMemberRepository(queries)
	projectInvitationRepository := repository.NewProjectInvitationRepository(queries)

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

	// Professor use cases
	professorProfileUC := professorUsecase.NewProfileUsecase(professorRepository)
	createProjectUC := professorUsecase.NewCreateProjectUsecase(projectRepository, professorRepository)
	listApplicationsByProjectUC := professorUsecase.NewListApplicationsByProjectUsecase(projectRepository, projectApplicationRepository, professorRepository)
	reviewApplicationUC := professorUsecase.NewReviewApplicationUsecase(projectRepository, projectApplicationRepository, projectMemberRepository, professorRepository)
	sendInvitationUC := professorUsecase.NewSendInvitationUsecase(projectRepository, projectInvitationRepository, projectMemberRepository, projectApplicationRepository, professorRepository, studentRepository)
	listInvitationsByProjectUC := professorUsecase.NewListInvitationsByProjectUsecase(projectRepository, projectInvitationRepository, professorRepository)
	cancelInvitationUC := professorUsecase.NewCancelInvitationUsecase(projectRepository, projectInvitationRepository, professorRepository)
	removeMemberUC := professorUsecase.NewRemoveMemberUsecase(projectRepository, projectMemberRepository, professorRepository)
	listStudentsUC := professorUsecase.NewListStudentsUsecase(studentRepository)

	// Student use cases
	studentProfileUC := studentUsecase.NewProfileUsecase(studentRepository)
	applyProjectUC := studentUsecase.NewApplyProjectUsecase(projectRepository, projectApplicationRepository, studentRepository)
	withdrawApplicationUC := studentUsecase.NewWithdrawApplicationUsecase(projectApplicationRepository, studentRepository)
	listMyApplicationsUC := studentUsecase.NewListMyApplicationsUsecase(projectApplicationRepository, studentRepository)
	listApplicationsByProjectIDUC := studentUsecase.NewListApplicationsPerIDUsecase(projectApplicationRepository, studentRepository)
	listMyInvitationsUC := studentUsecase.NewListMyInvitationsUsecase(projectInvitationRepository, studentRepository)
	respondInvitationUC := studentUsecase.NewRespondInvitationUsecase(projectRepository, projectInvitationRepository, projectMemberRepository, studentRepository)
	listMyProjectsUC := studentUsecase.NewListMyProjectsUsecase(projectMemberRepository, studentRepository)
	leaveProjectUC := studentUsecase.NewLeaveProjectUsecase(projectRepository, projectMemberRepository, studentRepository)

	// Generic project use cases
	listProjectsUC := projectUsecase.NewListProjectsUsecase(projectRepository)
	getProjectUC := projectUsecase.NewGetProjectUsecase(projectRepository)
	listMembersByProjectUC := projectUsecase.NewListMembersByProjectUsecase(projectRepository, projectMemberRepository)

	// Handlers
	authH := authHandler.New(
		registerAdminUC,
		registerProfessorUC,
		registerStudentUC,
		loginUC,
		refreshUC,
	)
	professorH := professorHandler.New(
		professorProfileUC,
		createProjectUC,
		listApplicationsByProjectUC,
		reviewApplicationUC,
		sendInvitationUC,
		listInvitationsByProjectUC,
		cancelInvitationUC,
		removeMemberUC,
		listStudentsUC,
	)
	studentH := studentHandler.New(
		studentProfileUC,
		applyProjectUC,
		withdrawApplicationUC,
		listMyApplicationsUC,
		listApplicationsByProjectIDUC,
		listMyInvitationsUC,
		respondInvitationUC,
		listMyProjectsUC,
		leaveProjectUC,
	)
	projectH := projectHandler.New(
		listProjectsUC,
		getProjectUC,
		listMembersByProjectUC,
	)

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

	r := router.NewRouter(app, authH, professorH, studentH, projectH, authMW)
	r.RegisterGeneralRoutes()
	r.RegisterSwaggerRoutes()
	r.RegisterAuthRoutes()
	r.RegisterStudentRoutes()
	r.RegisterProfessorRoutes()
	r.RegisterProjectRoutes()

	app.Listen(":3001")
}
