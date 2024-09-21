package routes

import (
	"log"
	"net/http"
	"time"

	"github.com/Yuriekokubu/workflow/version"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes registers all the routes for the application
func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	RegisterUserRoutes(r, db)
	RegisterItemRoutes(r, db)

	r.GET("/version", func(c *gin.Context) {
		versionID, err := version.GetLatestDBVersion(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"version": versionID})
	})

	r.GET("/test", func(ctx *gin.Context) {
		for i := 0; i < 10; i++ {
			log.Println(1)
			time.Sleep(1 * time.Second)
		}
		ctx.JSON(200, "test response")
	})
}
