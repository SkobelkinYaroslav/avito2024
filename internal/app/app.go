package app

import (
	"avito2024/internal/handler"
	"avito2024/internal/repository"
	"avito2024/internal/service"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func init() {
	fmt.Println(os.Getwd())
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %q", err)
	}
}

func SetupRouter() *gin.Engine {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	r := gin.Default()
	handler.NewHandler(svc, r)

	return r
}
