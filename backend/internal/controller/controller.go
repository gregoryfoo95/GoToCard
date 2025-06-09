package controller

import (
	"gotocard-backend/internal/service"
	"gotocard-backend/pkg/validator"
)

type Controllers struct {
	User           *UserController
	Category       *CategoryController
	CreditCard     *CreditCardController
	Spending       *SpendingController
	Recommendation *RecommendationController
	Scraping       *ScrapingController
}

func NewControllers(services *service.Services, validator *validator.Validator) *Controllers {
	return &Controllers{
		User:           NewUserController(services, validator),
		Category:       NewCategoryController(services, validator),
		CreditCard:     NewCreditCardController(services, validator),
		Spending:       NewSpendingController(services, validator),
		Recommendation: NewRecommendationController(services, validator),
		Scraping:       NewScrapingController(services, validator),
	}
} 