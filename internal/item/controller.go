package item

import (
	"errors"
	"fmt"
	"log"
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
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

func msgForTag(tag, param string) string {
	switch tag {
	case "required":
		return "จำเป็นต้องกรอกข้อมูลนี้" 
	case "email":
		return "Invalid email"
	case "gt":
		return fmt.Sprintf("Number must be greater than %v", param)
	case "gte":
		return fmt.Sprintf("Number must be greater than or equal to %v", param)
	default:
		return "Invalid value"
	}
}

func getValidationErrors(err error) []ApiError {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]ApiError, len(ve))
		for i, fe := range ve {
			out[i] = ApiError{
				Field:  fe.Field(),
				Reason: msgForTag(fe.Tag(), fe.Param()),
			}
		}
		return out
	}
	return nil
}

func respondSuccess(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusOK, gin.H{"message": message})
}

func parseID(ctx *gin.Context) (uint, error) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func handleServiceError(ctx *gin.Context, err error, notFoundMessage string) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": notFoundMessage})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
}


func (controller *Controller) CreateItem(ctx *gin.Context) {
	var request model.RequestItem

	if err := ctx.ShouldBindJSON(&request); err != nil {
		validationErrors := getValidationErrors(err)
		if validationErrors != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Validation failed",
				"errors":  validationErrors,
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("Received request body: %+v", request)

	item, err := controller.Service.Create(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, item)
}

func (controller *Controller) FindItems(ctx *gin.Context) {
	var request model.RequestFindItem

	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid query parameters", "error": err.Error()})
		return
	}

	items, err := controller.Service.Find(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, items)
}

func (controller *Controller) GetItemByID(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	item, err := controller.Service.FindByID(id)
	if err != nil {
		handleServiceError(ctx, err, "Item not found")
		return
	}

	ctx.JSON(http.StatusOK, item)
}

func (controller *Controller) UpdateItemStatus(ctx *gin.Context) {
	var request model.RequestUpdateItem

	if err := ctx.ShouldBindJSON(&request); err != nil {
		validationErrors := getValidationErrors(err)
		if validationErrors != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Validation failed",
				"errors":  validationErrors,
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body", "error": err.Error()})
		return
	}

	id, err := parseID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid item ID"})
		return
	}

	item, err := controller.Service.UpdateStatus(id, request.Status)
	if err != nil {
		handleServiceError(ctx, err, "Failed to update status")
		return
	}

	ctx.JSON(http.StatusOK, item)
}

func (controller *Controller) UpdateItemByID(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid item ID"})
		return
	}

	var request model.RequestItem
	if err := ctx.ShouldBindJSON(&request); err != nil {
		validationErrors := getValidationErrors(err)
		if validationErrors != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Validation failed",
				"errors":  validationErrors,
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body format",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("Received request body: %+v", request)

	item, err := controller.Service.UpdateItemByID(id, request)
	if err != nil {
		handleServiceError(ctx, err, "Error updating item")
		return
	}

	ctx.JSON(http.StatusOK, item)
}

func (controller *Controller) DeleteItem(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid item ID"})
		return
	}

	if err := controller.Service.DeleteItemByID(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	respondSuccess(ctx, "Item deleted successfully")
}
