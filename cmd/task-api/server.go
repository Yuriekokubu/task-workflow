package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Yuriekokubu/workflow/internal/middleware/LogMiddleware"
	"github.com/Yuriekokubu/workflow/internal/routes"
	"github.com/Yuriekokubu/workflow/lib"
)

func init() {
	log.Println("Init function started")

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	} else {
		log.Println("Loaded .env file successfully")
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "2024"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Panic("DATABASE_URL environment variable not set")
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}

	db.AutoMigrate(&middleware.LogEntry{})

	r := gin.Default()

	corsConfig := cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
			"http://localhost:8000",
			"http://127.0.0.1:8000",
			"http://localhost:8080",
			"http://127.0.0.1:8080",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           12 * time.Hour,
	}

	r.Use(cors.New(corsConfig))
	r.Use(middleware.LogMiddleware(db))

	routes.RegisterUserRoutes(r, db)
	routes.RegisterItemRoutes(r, db)

	// Set up the HTTP server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	lib.StartServer(server)
}
