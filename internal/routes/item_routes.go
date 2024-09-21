package routes

import (
	"github.com/Yuriekokubu/workflow/internal/item"
	"github.com/Yuriekokubu/workflow/internal/middleware/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"os"
)

// RegisterItemRoutes registers routes for item-related operations
func RegisterItemRoutes(r *gin.Engine, db *gorm.DB) {
	itemController := item.NewController(db)

	items := r.Group("/items")
	{
		items.GET("/:id", itemController.GetItemByID)
		items.PUT("/:id", itemController.UpdateItemByID)
		items.PATCH("/:id", itemController.UpdateItemStatus)
		items.DELETE("/:id", itemController.DeleteItem)

		// Protect the following routes with JWT auth
		items.Use(auth.Guard(os.Getenv("JWT_SECRET")))
		{
			items.POST("", itemController.CreateItem)
			items.GET("", itemController.FindItems)
			items.DELETE("/delete", itemController.DeleteItems)
		}
	}
}
