package repository

import (
	"gotocard-backend/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	List() ([]models.User, error)
}

type CategoryRepository interface {
	Create(category *models.Category) error
	GetByID(id uint) (*models.Category, error)
	GetByName(name string) (*models.Category, error)
	Update(category *models.Category) error
	Delete(id uint) error
	List() ([]models.Category, error)
}

type CreditCardRepository interface {
	Create(card *models.CreditCard) error
	GetByID(id uint) (*models.CreditCard, error)
	GetByBankAndName(bank, name string) (*models.CreditCard, error)
	Update(card *models.CreditCard) error
	Delete(id uint) error
	List() ([]models.CreditCard, error)
	GetActiveCards() ([]models.CreditCard, error)
}

type CardBenefitRepository interface {
	Create(benefit *models.CardBenefit) error
	GetByID(id uint) (*models.CardBenefit, error)
	GetByCardID(cardID uint) ([]models.CardBenefit, error)
	GetByCardAndCategory(cardID, categoryID uint) (*models.CardBenefit, error)
	Update(benefit *models.CardBenefit) error
	Delete(id uint) error
	List() ([]models.CardBenefit, error)
}

type UserSpendingRepository interface {
	Create(spending *models.UserSpending) error
	GetByID(id uint) (*models.UserSpending, error)
	GetByUserID(userID uint) ([]models.UserSpending, error)
	GetByUserAndCategory(userID, categoryID uint) ([]models.UserSpending, error)
	GetByUserAndMonth(userID uint, month, year int) ([]models.UserSpending, error)
	Update(spending *models.UserSpending) error
	Delete(id uint) error
}

type RecommendationRepository interface {
	Create(recommendation *models.Recommendation) error
	GetByID(id uint) (*models.Recommendation, error)
	GetByUserID(userID uint) ([]models.Recommendation, error)
	GetByUserAndCategory(userID, categoryID uint) ([]models.Recommendation, error)
	Update(recommendation *models.Recommendation) error
	Delete(id uint) error
	DeleteByUserID(userID uint) error
}

type Repositories struct {
	User           UserRepository
	Category       CategoryRepository
	CreditCard     CreditCardRepository
	CardBenefit    CardBenefitRepository
	UserSpending   UserSpendingRepository
	Recommendation RecommendationRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:           NewUserRepository(db),
		Category:       NewCategoryRepository(db),
		CreditCard:     NewCreditCardRepository(db),
		CardBenefit:    NewCardBenefitRepository(db),
		UserSpending:   NewUserSpendingRepository(db),
		Recommendation: NewRecommendationRepository(db),
	}
} 