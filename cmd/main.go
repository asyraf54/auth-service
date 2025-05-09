package main

import (
	"auth-service/internal/config"
	"auth-service/internal/handler"
	"auth-service/internal/middleware"
	"auth-service/internal/model"
	"auth-service/internal/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {
	config.LoadEnv()

	dsn := os.Getenv("DATABASE_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// migrate schema
	db.AutoMigrate(&model.User{})

	r := gin.Default()
	authRepo := repository.NewUserRepository(db)

	auth := r.Group("/")
	{
		auth.POST("/register", handler.Register(authRepo))
		auth.POST("/login", handler.Login(authRepo))
	}

	secure := r.Group("/")
	secure.Use(middleware.JWTAuth())
	{
		secure.GET("/me", handler.Me)
		secure.GET("/verify-token", handler.VerifyToken)
	}

	log.Fatal(r.Run(":" + os.Getenv("PORT")))
}
