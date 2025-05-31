package main

import (
	"auth-service/internal/config"
	"auth-service/internal/handler"
	"auth-service/internal/middleware"
	"auth-service/internal/model"
	"auth-service/internal/repository"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	config.LoadEnv()

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// migrate schema
	db.AutoMigrate(&model.User{})

	r := gin.Default()
	authRepo := repository.NewUserRepository(db)
	authHandler := handler.NewAuthHandler(authRepo);

	auth := r.Group("/")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	secure := r.Group("/")
	secure.Use(middleware.JWTAuth("access"))
	{
		secure.GET("/me", authHandler.Me)
		secure.GET("/verify-token", authHandler.VerifyToken)
	}

	secureRefresh := r.Group("/")
	secureRefresh.Use(middleware.JWTAuth("refresh"))
	{
		secureRefresh.POST("/refresh-token", authHandler.RefreshToken)
	}

	log.Fatal(r.Run(":" + os.Getenv("PORT")))
}
