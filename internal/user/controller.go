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

	if err := ctx.Bind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	username, userID, token, err := controller.Service.Login(request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	fmt.Println("Generated token:", token)

	ctx.SetCookie(
		"token",
		fmt.Sprintf("Bearer %s", token),
		2592000,     
		"/",         
		"localhost", 
		false,       
		true,        
	)

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "Login succeeded",
		"token":    token,
		"username": username,
		"userID":   userID,
	})
}

func (controller Controller) Register(ctx *gin.Context) {
	var request model.RequestRegister

	if err := ctx.Bind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	user, err := controller.Service.Register(request)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "Registration succeeded",
		"username": user.Username,
	})
}
