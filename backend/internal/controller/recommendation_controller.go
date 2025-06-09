package controller

import (
	"net/http"
	"strconv"

	"gotocard-backend/internal/models"
	"gotocard-backend/internal/service"
	"gotocard-backend/pkg/validator"

	"github.com/gin-gonic/gin"
)

type RecommendationController struct {
	services  *service.Services
	validator *validator.Validator
}

func NewRecommendationController(services *service.Services, validator *validator.Validator) *RecommendationController {
	return &RecommendationController{
		services:  services,
		validator: validator,
	}
}

func (c *RecommendationController) GenerateRecommendations(ctx *gin.Context) {
	idParam := ctx.Param("userId")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	recommendations, err := c.services.Recommendation.GenerateRecommendations(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":         "Recommendations generated successfully",
		"recommendations": recommendations,
	})
}

func (c *RecommendationController) GetRecommendations(ctx *gin.Context) {
	idParam := ctx.Param("userId")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	recommendations, err := c.services.Recommendation.GetRecommendationsByUser(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"recommendations": recommendations})
}

type SpendingController struct {
	services  *service.Services
	validator *validator.Validator
}

func NewSpendingController(services *service.Services, validator *validator.Validator) *SpendingController {
	return &SpendingController{
		services:  services,
		validator: validator,
	}
}

func (c *SpendingController) AddSpending(ctx *gin.Context) {
	idParam := ctx.Param("userId")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req models.SpendingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := c.validator.Validate(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.services.Spending.AddSpending(uint(userID), &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Spending added successfully"})
}

func (c *SpendingController) GetUserSpending(ctx *gin.Context) {
	idParam := ctx.Param("userId")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	spendings, err := c.services.Spending.GetUserSpending(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"spendings": spendings})
}

type CreditCardController struct {
	services  *service.Services
	validator *validator.Validator
}

func NewCreditCardController(services *service.Services, validator *validator.Validator) *CreditCardController {
	return &CreditCardController{
		services:  services,
		validator: validator,
	}
}

func (c *CreditCardController) ListCreditCards(ctx *gin.Context) {
	cards, err := c.services.CreditCard.ListCreditCards()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch credit cards"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"cards": cards})
}

func (c *CreditCardController) GetCreditCard(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card ID"})
		return
	}

	card, err := c.services.CreditCard.GetCreditCardByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Credit card not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"card": card})
}

type ScrapingController struct {
	services  *service.Services
	validator *validator.Validator
}

func NewScrapingController(services *service.Services, validator *validator.Validator) *ScrapingController {
	return &ScrapingController{
		services:  services,
		validator: validator,
	}
}

func (c *ScrapingController) ScrapeCardData(ctx *gin.Context) {
	source := ctx.Query("source")

	var err error
	if source != "" {
		// Scrape from specific source
		err = c.services.Scraping.ScrapeCardDataBySource(source)
	} else {
		// Scrape from all sources
		err = c.services.Scraping.ScrapeCardData()
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Card data scraping completed successfully"})
}
