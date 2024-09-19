package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LogEntry struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    int    `gorm:"type:int4"`
	Method    string `gorm:"type:varchar(10)"`
	URL       string `gorm:"type:text"`
	Action    string `gorm:"type:text"`
	Timestamp time.Time
}

func LogMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userIDUint uint

		if userID, exists := c.Get("user_id"); exists {
			if uid, ok := userID.(uint); ok {
				userIDUint = uid
			}
		}

		entry := LogEntry{
			UserID:    int(userIDUint),
			Method:    c.Request.Method,
			URL:       c.Request.URL.Path,
			Action:    "User action on " + c.Request.URL.Path,
			Timestamp: time.Now(),
		}

		log.Printf("Entry being logged: %+v\n", entry)

		if err := db.Create(&entry).Error; err != nil {
			log.Printf("Failed to log user action: %v", err)
		}

		c.Next()
	}
}
