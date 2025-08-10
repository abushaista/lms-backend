package main

import (
	"log"
	"os"

	"github.com/abushaista/lms-backend/delivery/http"
	"github.com/abushaista/lms-backend/infrastructure/config"
	"github.com/abushaista/lms-backend/infrastructure/database"
	"github.com/abushaista/lms-backend/internal/domain"
	"github.com/abushaista/lms-backend/internal/repository"
	"github.com/abushaista/lms-backend/internal/usecase"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

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

	if err := db.AutoMigrate(&domain.Book{}, &domain.Category{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	rBook := repository.NewGormBookRepository(db)
	ucBook := usecase.NewBookUsecase(rBook)
	rCategory := repository.NewGormCategoryRepository(db)
	ucCategory := usecase.NewCategoryUseCase(rCategory)
	e := echo.New()

	api := e.Group("/api")
	http.NewBookHandler(api, ucBook)
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
