package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"

	"piggy.com/internal/db/repo"
	"piggy.com/internal/handlers"
	"piggy.com/internal/middleware"
	"piggy.com/internal/piggyservice"
)

func main() {
		// Load .env file
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

		// Build DB URL from env
	dbUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
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
	protected.GET("/dashboard/stats", handlers.GetDashboardStats)
	fmt.Println("Server running on port 8080")
	route.Run()
}