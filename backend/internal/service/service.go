package service

import (
	"gotocard-backend/internal/models"
	"gotocard-backend/internal/repository"
)

type UserService interface {
	CreateUser(req *models.CreateUserRequest) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
	ListUsers() ([]models.User, error)
}

type CategoryService interface {
	CreateCategory(category *models.Category) error
	GetCategoryByID(id uint) (*models.Category, error)
	GetCategoryByName(name string) (*models.Category, error)
	UpdateCategory(category *models.Category) error
	DeleteCategory(id uint) error
	ListCategories() ([]models.Category, error)
}

type CreditCardService interface {
	CreateCreditCard(card *models.CreditCard) error
	GetCreditCardByID(id uint) (*models.CreditCard, error)
	UpdateCreditCard(card *models.CreditCard) error
	DeleteCreditCard(id uint) error
	ListCreditCards() ([]models.CreditCard, error)
	GetActiveCards() ([]models.CreditCard, error)
}

type SpendingService interface {
	AddSpending(userID uint, req *models.SpendingRequest) error
	GetUserSpending(userID uint) ([]models.UserSpending, error)
	GetUserSpendingByCategory(userID, categoryID uint) ([]models.UserSpending, error)
	GetUserSpendingByMonth(userID uint, month, year int) ([]models.UserSpending, error)
	UpdateSpending(spending *models.UserSpending) error
	DeleteSpending(id uint) error
}

type RecommendationService interface {
	GenerateRecommendations(userID uint) ([]models.RecommendationResponse, error)
	GetRecommendationsByUser(userID uint) ([]models.RecommendationResponse, error)
	GetRecommendationsByCategory(userID, categoryID uint) ([]models.RecommendationResponse, error)
	RefreshRecommendations(userID uint) error
}

type ScrapingService interface {
	ScrapeCardData() error
	ScrapeCardDataBySource(source string) error
	UpdateCardDatabase() error
}

type Services struct {
	User           UserService
	Category       CategoryService
	CreditCard     CreditCardService
	Spending       SpendingService
	Recommendation RecommendationService
	Scraping       ScrapingService
}

func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		User:           NewUserService(repos),
		Category:       NewCategoryService(repos),
		CreditCard:     NewCreditCardService(repos),
		Spending:       NewSpendingService(repos),
		Recommendation: NewRecommendationService(repos),
		Scraping:       NewScrapingService(repos),
	}
}
