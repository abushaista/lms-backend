package main

import (
	"log"
	"os"

	"github.com/abushaista/lms-backend/delivery/http"
	libMiddleWare "github.com/abushaista/lms-backend/delivery/middleware"
	_ "github.com/abushaista/lms-backend/docs"
	"github.com/abushaista/lms-backend/infrastructure/config"
	"github.com/abushaista/lms-backend/infrastructure/database"
	"github.com/abushaista/lms-backend/infrastructure/logger"
	validatorInfra "github.com/abushaista/lms-backend/infrastructure/validator"
	"github.com/abushaista/lms-backend/internal/domain"
	"github.com/abushaista/lms-backend/internal/repository"
	"github.com/abushaista/lms-backend/internal/usecase"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Library Management API
// @version 1.0
// @description Simple Library Management System in Go using Clean Architecture, Echo, MySQL, GORM, and Swagger
// @host localhost:8080
// @BasePath /
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading it, continuing with env variables")
	}

	cfg := config.LoadConfig()

	db, err := database.NewGormDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&domain.Book{}, &domain.Category{}, &domain.User{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}
	rootLogger := logger.NewLogger()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(libMiddleWare.CorrelationMiddleware)
	e.Use(libMiddleWare.AttachRequestLogger(rootLogger))
	e.Validator = validatorInfra.NewEchoValidator()
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	rUser := repository.NewGormUserRepository(db)
	ucAuth := usecase.NewAuthUseCase(rUser, cfg.JWTSecret)

	public := e.Group("api")
	http.NewAuthHandler(public, ucAuth, rootLogger)

	ucClean := usecase.NewCleanUpUsecase()

	http.NewCleanUpHandler(public, ucClean, rootLogger)

	api := e.Group("/api")

	api.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(cfg.JWTSecret),
	}))

	rBook := repository.NewGormBookRepository(db)
	ucBook := usecase.NewBookUsecase(rBook)
	rCategory := repository.NewGormCategoryRepository(db)
	ucCategory := usecase.NewCategoryUseCase(rCategory)

	http.NewBookHandler(api, ucBook, rootLogger)
	http.NewCategoryHandler(api, ucCategory)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("starting server on :%s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}

}
