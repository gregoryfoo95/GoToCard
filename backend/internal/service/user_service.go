package service

import (
	"fmt"
	"gotocard-backend/internal/models"
	"gotocard-backend/internal/repository"
)

type userService struct {
	repos *repository.Repositories
}

func NewUserService(repos *repository.Repositories) UserService {
	return &userService{repos: repos}
}

func (s *userService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	// Check if user already exists
	existingUser, err := s.repos.User.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	user := &models.User{
		Name:  req.Name,
		Email: req.Email,
	}

	err = s.repos.User.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
	user, err := s.repos.User.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	user, err := s.repos.User.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

func (s *userService) UpdateUser(user *models.User) error {
	err := s.repos.User.Update(user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (s *userService) DeleteUser(id uint) error {
	// Check if user exists
	_, err := s.repos.User.GetByID(id)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Delete user recommendations first
	err = s.repos.Recommendation.DeleteByUserID(id)
	if err != nil {
		return fmt.Errorf("failed to delete user recommendations: %w", err)
	}

	err = s.repos.User.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (s *userService) ListUsers() ([]models.User, error) {
	users, err := s.repos.User.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}

type categoryService struct {
	repos *repository.Repositories
}

func NewCategoryService(repos *repository.Repositories) CategoryService {
	return &categoryService{repos: repos}
}

func (s *categoryService) CreateCategory(category *models.Category) error {
	// Check if category already exists
	existing, err := s.repos.Category.GetByName(category.Name)
	if err == nil && existing != nil {
		return fmt.Errorf("category with name %s already exists", category.Name)
	}

	err = s.repos.Category.Create(category)
	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}
	return nil
}

func (s *categoryService) GetCategoryByID(id uint) (*models.Category, error) {
	category, err := s.repos.Category.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}
	return category, nil
}

func (s *categoryService) GetCategoryByName(name string) (*models.Category, error) {
	category, err := s.repos.Category.GetByName(name)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}
	return category, nil
}

func (s *categoryService) UpdateCategory(category *models.Category) error {
	err := s.repos.Category.Update(category)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	return nil
}

func (s *categoryService) DeleteCategory(id uint) error {
	err := s.repos.Category.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return nil
}

func (s *categoryService) ListCategories() ([]models.Category, error) {
	categories, err := s.repos.Category.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	return categories, nil
}

type creditCardService struct {
	repos *repository.Repositories
}

func NewCreditCardService(repos *repository.Repositories) CreditCardService {
	return &creditCardService{repos: repos}
}

func (s *creditCardService) CreateCreditCard(card *models.CreditCard) error {
	err := s.repos.CreditCard.Create(card)
	if err != nil {
		return fmt.Errorf("failed to create credit card: %w", err)
	}
	return nil
}

func (s *creditCardService) GetCreditCardByID(id uint) (*models.CreditCard, error) {
	card, err := s.repos.CreditCard.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("credit card not found: %w", err)
	}
	return card, nil
}

func (s *creditCardService) UpdateCreditCard(card *models.CreditCard) error {
	err := s.repos.CreditCard.Update(card)
	if err != nil {
		return fmt.Errorf("failed to update credit card: %w", err)
	}
	return nil
}

func (s *creditCardService) DeleteCreditCard(id uint) error {
	err := s.repos.CreditCard.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete credit card: %w", err)
	}
	return nil
}

func (s *creditCardService) ListCreditCards() ([]models.CreditCard, error) {
	cards, err := s.repos.CreditCard.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list credit cards: %w", err)
	}
	return cards, err
}

func (s *creditCardService) GetActiveCards() ([]models.CreditCard, error) {
	cards, err := s.repos.CreditCard.GetActiveCards()
	if err != nil {
		return nil, fmt.Errorf("failed to get active cards: %w", err)
	}
	return cards, nil
}

type spendingService struct {
	repos *repository.Repositories
}

func NewSpendingService(repos *repository.Repositories) SpendingService {
	return &spendingService{repos: repos}
}

func (s *spendingService) AddSpending(userID uint, req *models.SpendingRequest) error {
	// Verify user exists
	_, err := s.repos.User.GetByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Verify category exists
	_, err = s.repos.Category.GetByID(req.CategoryID)
	if err != nil {
		return fmt.Errorf("category not found: %w", err)
	}

	spending := &models.UserSpending{
		UserID:     userID,
		CategoryID: req.CategoryID,
		Amount:     req.Amount,
		Month:      req.Month,
		Year:       req.Year,
	}

	err = s.repos.UserSpending.Create(spending)
	if err != nil {
		return fmt.Errorf("failed to add spending: %w", err)
	}
	return nil
}

func (s *spendingService) GetUserSpending(userID uint) ([]models.UserSpending, error) {
	spendings, err := s.repos.UserSpending.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user spending: %w", err)
	}
	return spendings, nil
}

func (s *spendingService) GetUserSpendingByCategory(userID, categoryID uint) ([]models.UserSpending, error) {
	spendings, err := s.repos.UserSpending.GetByUserAndCategory(userID, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user spending by category: %w", err)
	}
	return spendings, nil
}

func (s *spendingService) GetUserSpendingByMonth(userID uint, month, year int) ([]models.UserSpending, error) {
	spendings, err := s.repos.UserSpending.GetByUserAndMonth(userID, month, year)
	if err != nil {
		return nil, fmt.Errorf("failed to get user spending by month: %w", err)
	}
	return spendings, nil
}

func (s *spendingService) UpdateSpending(spending *models.UserSpending) error {
	err := s.repos.UserSpending.Update(spending)
	if err != nil {
		return fmt.Errorf("failed to update spending: %w", err)
	}
	return nil
}

func (s *spendingService) DeleteSpending(id uint) error {
	err := s.repos.UserSpending.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete spending: %w", err)
	}
	return nil
} 