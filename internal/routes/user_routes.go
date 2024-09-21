package routes

import (
	"github.com/Yuriekokubu/workflow/internal/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"os"
)

// RegisterUserRoutes registers routes for user-related operations
func RegisterUserRoutes(r *gin.Engine, db *gorm.DB) {
	userController := user.NewController(db, os.Getenv("JWT_SECRET"))

	// User routes
	r.POST("/login", userController.Login)
	r.POST("/signup", userController.Register)
}
