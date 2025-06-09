package repository

import (
	"gotocard-backend/internal/models"
	"gorm.io/gorm"
)

type creditCardRepository struct {
	db *gorm.DB
}

func NewCreditCardRepository(db *gorm.DB) CreditCardRepository {
	return &creditCardRepository{db: db}
}

func (r *creditCardRepository) Create(card *models.CreditCard) error {
	return r.db.Create(card).Error
}

func (r *creditCardRepository) GetByID(id uint) (*models.CreditCard, error) {
	var card models.CreditCard
	err := r.db.Preload("CardBenefits").Preload("CardBenefits.Category").First(&card, id).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

func (r *creditCardRepository) GetByBankAndName(bank, name string) (*models.CreditCard, error) {
	var card models.CreditCard
	err := r.db.Where("bank = ? AND name = ?", bank, name).First(&card).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

func (r *creditCardRepository) Update(card *models.CreditCard) error {
	return r.db.Save(card).Error
}

func (r *creditCardRepository) Delete(id uint) error {
	return r.db.Delete(&models.CreditCard{}, id).Error
}

func (r *creditCardRepository) List() ([]models.CreditCard, error) {
	var cards []models.CreditCard
	err := r.db.Preload("CardBenefits").Preload("CardBenefits.Category").Find(&cards).Error
	return cards, err
}

func (r *creditCardRepository) GetActiveCards() ([]models.CreditCard, error) {
	var cards []models.CreditCard
	err := r.db.Where("is_active = ?", true).Preload("CardBenefits").Preload("CardBenefits.Category").Find(&cards).Error
	return cards, err
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) GetByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetByName(name string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("name = ?", name).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}

func (r *categoryRepository) List() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

type cardBenefitRepository struct {
	db *gorm.DB
}

func NewCardBenefitRepository(db *gorm.DB) CardBenefitRepository {
	return &cardBenefitRepository{db: db}
}

func (r *cardBenefitRepository) Create(benefit *models.CardBenefit) error {
	return r.db.Create(benefit).Error
}

func (r *cardBenefitRepository) GetByID(id uint) (*models.CardBenefit, error) {
	var benefit models.CardBenefit
	err := r.db.Preload("Card").Preload("Category").First(&benefit, id).Error
	if err != nil {
		return nil, err
	}
	return &benefit, nil
}

func (r *cardBenefitRepository) GetByCardID(cardID uint) ([]models.CardBenefit, error) {
	var benefits []models.CardBenefit
	err := r.db.Where("card_id = ?", cardID).Preload("Category").Find(&benefits).Error
	return benefits, err
}

func (r *cardBenefitRepository) GetByCardAndCategory(cardID, categoryID uint) (*models.CardBenefit, error) {
	var benefit models.CardBenefit
	err := r.db.Where("card_id = ? AND category_id = ?", cardID, categoryID).First(&benefit).Error
	if err != nil {
		return nil, err
	}
	return &benefit, nil
}

func (r *cardBenefitRepository) Update(benefit *models.CardBenefit) error {
	return r.db.Save(benefit).Error
}

func (r *cardBenefitRepository) Delete(id uint) error {
	return r.db.Delete(&models.CardBenefit{}, id).Error
}

func (r *cardBenefitRepository) List() ([]models.CardBenefit, error) {
	var benefits []models.CardBenefit
	err := r.db.Preload("Card").Preload("Category").Find(&benefits).Error
	return benefits, err
}

type userSpendingRepository struct {
	db *gorm.DB
}

func NewUserSpendingRepository(db *gorm.DB) UserSpendingRepository {
	return &userSpendingRepository{db: db}
}

func (r *userSpendingRepository) Create(spending *models.UserSpending) error {
	return r.db.Create(spending).Error
}

func (r *userSpendingRepository) GetByID(id uint) (*models.UserSpending, error) {
	var spending models.UserSpending
	err := r.db.Preload("User").Preload("Category").First(&spending, id).Error
	if err != nil {
		return nil, err
	}
	return &spending, nil
}

func (r *userSpendingRepository) GetByUserID(userID uint) ([]models.UserSpending, error) {
	var spendings []models.UserSpending
	err := r.db.Where("user_id = ?", userID).Preload("Category").Find(&spendings).Error
	return spendings, err
}

func (r *userSpendingRepository) GetByUserAndCategory(userID, categoryID uint) ([]models.UserSpending, error) {
	var spendings []models.UserSpending
	err := r.db.Where("user_id = ? AND category_id = ?", userID, categoryID).Preload("Category").Find(&spendings).Error
	return spendings, err
}

func (r *userSpendingRepository) GetByUserAndMonth(userID uint, month, year int) ([]models.UserSpending, error) {
	var spendings []models.UserSpending
	err := r.db.Where("user_id = ? AND month = ? AND year = ?", userID, month, year).Preload("Category").Find(&spendings).Error
	return spendings, err
}

func (r *userSpendingRepository) Update(spending *models.UserSpending) error {
	return r.db.Save(spending).Error
}

func (r *userSpendingRepository) Delete(id uint) error {
	return r.db.Delete(&models.UserSpending{}, id).Error
}

type recommendationRepository struct {
	db *gorm.DB
}

func NewRecommendationRepository(db *gorm.DB) RecommendationRepository {
	return &recommendationRepository{db: db}
}

func (r *recommendationRepository) Create(recommendation *models.Recommendation) error {
	return r.db.Create(recommendation).Error
}

func (r *recommendationRepository) GetByID(id uint) (*models.Recommendation, error) {
	var recommendation models.Recommendation
	err := r.db.Preload("User").Preload("Category").Preload("Card").First(&recommendation, id).Error
	if err != nil {
		return nil, err
	}
	return &recommendation, nil
}

func (r *recommendationRepository) GetByUserID(userID uint) ([]models.Recommendation, error) {
	var recommendations []models.Recommendation
	err := r.db.Where("user_id = ?", userID).Preload("Category").Preload("Card").Find(&recommendations).Error
	return recommendations, err
}

func (r *recommendationRepository) GetByUserAndCategory(userID, categoryID uint) ([]models.Recommendation, error) {
	var recommendations []models.Recommendation
	err := r.db.Where("user_id = ? AND category_id = ?", userID, categoryID).Preload("Category").Preload("Card").Find(&recommendations).Error
	return recommendations, err
}

func (r *recommendationRepository) Update(recommendation *models.Recommendation) error {
	return r.db.Save(recommendation).Error
}

func (r *recommendationRepository) Delete(id uint) error {
	return r.db.Delete(&models.Recommendation{}, id).Error
}

func (r *recommendationRepository) DeleteByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.Recommendation{}).Error
} 