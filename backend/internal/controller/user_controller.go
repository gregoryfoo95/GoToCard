package controller

import (
	"net/http"
	"strconv"

	"gotocard-backend/internal/models"
	"gotocard-backend/internal/service"
	"gotocard-backend/pkg/validator"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	services  *service.Services
	validator *validator.Validator
}

func NewUserController(services *service.Services, validator *validator.Validator) *UserController {
	return &UserController{
		services:  services,
		validator: validator,
	}
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var req models.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := c.validator.Validate(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.services.User.CreateUser(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}

func (c *UserController) GetUser(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := c.services.User.GetUserByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func (c *UserController) ListUsers(ctx *gin.Context) {
	users, err := c.services.User.ListUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"users": users})
}

type CategoryController struct {
	services  *service.Services
	validator *validator.Validator
}

func NewCategoryController(services *service.Services, validator *validator.Validator) *CategoryController {
	return &CategoryController{
		services:  services,
		validator: validator,
	}
}

func (c *CategoryController) CreateCategory(ctx *gin.Context) {
	var category models.Category
	if err := ctx.ShouldBindJSON(&category); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := c.validator.Validate(&category); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.services.Category.CreateCategory(&category)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":  "Category created successfully",
		"category": category,
	})
}

func (c *CategoryController) ListCategories(ctx *gin.Context) {
	categories, err := c.services.Category.ListCategories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"categories": categories})
} 