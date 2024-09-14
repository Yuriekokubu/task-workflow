package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Yuriekokubu/workflow/internal/auth"
	"github.com/Yuriekokubu/workflow/internal/item"
	"github.com/Yuriekokubu/workflow/internal/user"
	"github.com/Yuriekokubu/workflow/lib"
	"github.com/Yuriekokubu/workflow/version"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	controller := item.NewController(db)
	userController := user.NewController(db, os.Getenv("JWT_SECRET"))

	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:8000",
		"http://127.0.0.1:8000",
	}
	r.Use(cors.New(config))

	items := r.Group("/items")
	{
		items.GET("/:id", controller.GetItemByID)
		items.PUT("/:id", controller.UpdateItemByID)
		items.PATCH("/:id", controller.UpdateItemStatus)
		items.DELETE("/:id", controller.DeleteItem)
		items.Use(auth.Guard(os.Getenv("JWT_SECRET")))
		{
			items.POST("", controller.CreateItem)
			items.GET("", controller.FindItems)
		}
	}

	r.GET("/version", func(c *gin.Context) {
		versionID, err := version.GetLatestDBVersion(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"version": versionID})
	})

	r.POST("/login", userController.Login)

	r.GET("/test", func(ctx *gin.Context) {
		for i := 0; i < 10; i++ {
			log.Println(1)
			time.Sleep(1 * time.Second)
		}
		ctx.JSON(200, "test response")
	})

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	lib.StartServer(server)
}
