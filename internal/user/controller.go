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
	var (
		request model.RequestLogin
	)

	if err := ctx.Bind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	token, err := controller.Service.Login(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	fmt.Println("Generated token:", token)

	ctx.SetCookie("token", fmt.Sprintf("Bearer %v", token), 10, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login Succeed",
		"token":   token,
	})
}
