package item

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"github.com/Yuriekokubu/workflow/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Controller struct {
	Service Service
}

func NewController(db *gorm.DB) Controller {
	return Controller{
		Service: NewService(db),
	}
}

type ApiError struct {
	Field  string
	Reason string
}

func msgForTag(tag, param string) string {
	switch tag {
	case "required":
		return "จำเป็นต้องกรอกข้อมูลนี้"
	case "email":
		return "Invalid email"
	case "gt":
		return fmt.Sprintf("Number must greater than %v", param)
	case "gte":
		return fmt.Sprintf("Number must greater than or equal %v", param)
	}
	return ""
}

func getValidationErrors(err error) []ApiError {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]ApiError, len(ve))
		for i, fe := range ve {
			out[i] = ApiError{fe.Field(), msgForTag(fe.Tag(), fe.Param())}
		}
		return out
	}
	return nil
}

func (controller Controller) CreateItem(ctx *gin.Context) {
	// Bind
	var request model.RequestItem

	if err := ctx.Bind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": getValidationErrors(err),
		})
		return
	}

	// Create item
	item, err := controller.Service.Create(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Response
	ctx.JSON(http.StatusCreated, item)
}

func (controller Controller) FindItems(ctx *gin.Context) {
	// Bind query parameters
	var (
		request model.RequestFindItem
	)

	if err := ctx.BindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	// Find
	items, err := controller.Service.Find(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, items)
}

func (controller Controller) GetItemByID(ctx *gin.Context) {
	// Get the ID from the URL path
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID",
		})
		return
	}

	// Fetch the item by ID using the service
	item, err := controller.Service.FindByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "Item not found",
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		}
		return
	}

	// Return the item in the response
	ctx.JSON(http.StatusOK, item)
}

func (controller Controller) UpdateItemStatus(ctx *gin.Context) {
	// Bind
	var (
		request model.RequestUpdateItem
	)

	if err := ctx.Bind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	// Path param
	id, _ := strconv.Atoi(ctx.Param("id"))

	// Update status
	item, err := controller.Service.UpdateStatus(uint(id), request.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, item)
}

func (controller Controller) UpdateItemByID(ctx *gin.Context) {
	// Parse the ID from the path
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid item ID",
		})
		return
	}

	// Bind the request body to the model
	var request model.RequestItem
	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": getValidationErrors(err),
		})
		return
	}

	// Update the item
	item, err := controller.Service.UpdateItemByID(uint(id), request)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "Item not found",
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		}
		return
	}

	// Return the updated item
	ctx.JSON(http.StatusOK, item)
}

func (controller Controller) DeleteItem(ctx *gin.Context) {
	// Get the item ID from the path
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid item ID"})
		return
	}

	// Delete the item
	if err := controller.Service.DeleteItemByID(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Respond with a success message
	ctx.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}
