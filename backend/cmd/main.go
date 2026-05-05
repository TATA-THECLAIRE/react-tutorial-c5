package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"piggy.com/internal/db/repo"
	"piggy.com/internal/handlers"
	"piggy.com/internal/piggyservice"
	"piggy.com/internal/middleware"
	
)

func main() {
	route := gin.Default()

	// Configure Cors
	route.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))



	// Healthcheck
	route.GET("/api/v1/healthcheck", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Healthy!",
		})
	})

	// Initialize repo and apply migrations
	ctx := context.Background()
	dbUrl := "postgres://piggy:secret@127.0.0.1:5432/piggy?sslmode=disable"
	
	dbConn, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		panic(err)
	}
	fmt.Println("Database connection established!")
	repostory := repo.NewRepository(dbConn)
	if err :=repo.MigrateUp(dbUrl, "./internal/db/migrations", zerolog.Nop().With().Logger());err !=nil{
		panic(err)
	}

	// Initialize service
	appService := piggyservice.NewService(repostory)
	handlers := handlers.NewHandler(appService)

		// Auth routes
	route.POST("/api/v1/auth/register", handlers.Register)
	route.POST("/api/v1/auth/login", handlers.Login)

	// Define application endpoints
	protected := route.Group("/api/v1", middleware.RequireAuth())
	protected.POST("/transactions", handlers.CreateTransaction)
	protected.GET("/transactions", handlers.GetTransactions)
	fmt.Println("Server running on port 8080")
	route.Run()
}