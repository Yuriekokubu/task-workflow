package user

import (
	"fmt"
	"github.com/Yuriekokubu/workflow/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller struct {
	Service Service
}

func NewController(db *gorm.DB, secret string) Controller {
	return Controller{
		Service: NewService(db, secret),
	}
}

func (controller Controller) Login(ctx *gin.Context) {
	var request model.RequestLogin

	// Bind the request to the model
	if err := ctx.Bind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Call the service to get the token
	username, userID, token, err := controller.Service.Login(request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Log the token for debugging purposes
	fmt.Println("Generated token:", token)

	// Set the cookie with a more reasonable expiry and appropriate flags
	ctx.SetCookie(
		"token",
		fmt.Sprintf("Bearer %s", token),
		2592000,     // Cookie expiry time in seconds (1 hour here)
		"/",         // Path
		"localhost", // Domain (adjust as needed for production)
		false,       // Secure flag (set to true in production with HTTPS)
		true,        // HttpOnly flag to prevent client-side access
	)

	// Respond to the client
	ctx.JSON(http.StatusOK, gin.H{
		"message":  "Login succeeded",
		"token":    token,
		"username": username,
		"userID":   userID,
	})
}
