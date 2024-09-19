package routes

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Yuriekokubu/workflow/internal/middleware/auth"
	"github.com/Yuriekokubu/workflow/internal/item"
	"github.com/Yuriekokubu/workflow/internal/user"
	"github.com/Yuriekokubu/workflow/version"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	itemController := item.NewController(db)
	userController := user.NewController(db, os.Getenv("JWT_SECRET"))

	items := r.Group("/items")
	{
		items.GET("/:id", itemController.GetItemByID)
		items.PUT("/:id", itemController.UpdateItemByID)
		items.PATCH("/:id", itemController.UpdateItemStatus)
		items.DELETE("/:id", itemController.DeleteItem)
		items.Use(auth.Guard(os.Getenv("JWT_SECRET")))

		{
			items.POST("", itemController.CreateItem)
			items.GET("", itemController.FindItems)
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
	r.POST("/register", userController.Register)

	r.GET("/test", func(ctx *gin.Context) {
		for i := 0; i < 10; i++ {
			log.Println(1)
			time.Sleep(1 * time.Second)
		}
		ctx.JSON(200, "test response")
	})
}
