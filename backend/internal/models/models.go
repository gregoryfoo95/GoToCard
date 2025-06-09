package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Name      string         `json:"name" gorm:"not null" validate:"required,min=2,max=100"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Category struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null" validate:"required,min=2,max=50"`
	Description string         `json:"description"`
	Icon        string         `json:"icon"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreditCard struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name" gorm:"not null" validate:"required,min=2,max=100"`
	Bank         string         `json:"bank" gorm:"not null" validate:"required,min=2,max=50"`
	CardType     string         `json:"card_type" gorm:"not null" validate:"required,oneof=visa mastercard amex"`
	AnnualFee    float64        `json:"annual_fee" gorm:"default:0"`
	ImageURL     string         `json:"image_url"`
	Description  string         `json:"description"`
	MinIncome    float64        `json:"min_income" gorm:"default:0"`
	WelcomeBonus string         `json:"welcome_bonus"`
	SourceURL    string         `json:"source_url"`
	IsActive     bool           `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	CardBenefits []CardBenefit `json:"card_benefits" gorm:"foreignKey:CardID"`
}

type CardBenefit struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	CardID        uint           `json:"card_id" gorm:"not null"`
	CategoryID    uint           `json:"category_id" gorm:"not null"`
	CashbackRate  float64        `json:"cashback_rate" gorm:"default:0"`
	PointsRate    float64        `json:"points_rate" gorm:"default:0"`
	MilesRate     float64        `json:"miles_rate" gorm:"default:0"`
	Cap           float64        `json:"cap" gorm:"default:0"`
	MinSpend      float64        `json:"min_spend" gorm:"default:0"`
	Description   string         `json:"description"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Card     CreditCard `json:"card" gorm:"foreignKey:CardID"`
	Category Category   `json:"category" gorm:"foreignKey:CategoryID"`
}

type UserSpending struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id" gorm:"not null"`
	CategoryID uint           `json:"category_id" gorm:"not null"`
	Amount     float64        `json:"amount" gorm:"not null" validate:"required,min=0"`
	Month      int            `json:"month" gorm:"not null" validate:"required,min=1,max=12"`
	Year       int            `json:"year" gorm:"not null" validate:"required,min=2020"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	User     User     `json:"user" gorm:"foreignKey:UserID"`
	Category Category `json:"category" gorm:"foreignKey:CategoryID"`
}

type Recommendation struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	UserID       uint           `json:"user_id" gorm:"not null"`
	CategoryID   uint           `json:"category_id" gorm:"not null"`
	CardID       uint           `json:"card_id" gorm:"not null"`
	Score        float64        `json:"score" gorm:"not null"`
	EstimatedReward float64     `json:"estimated_reward"`
	Reason       string         `json:"reason"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	User     User       `json:"user" gorm:"foreignKey:UserID"`
	Category Category   `json:"category" gorm:"foreignKey:CategoryID"`
	Card     CreditCard `json:"card" gorm:"foreignKey:CardID"`
}

// Request/Response DTOs
type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
}

type SpendingRequest struct {
	CategoryID uint    `json:"category_id" validate:"required"`
	Amount     float64 `json:"amount" validate:"required,min=0"`
	Month      int     `json:"month" validate:"required,min=1,max=12"`
	Year       int     `json:"year" validate:"required,min=2020"`
}

type RecommendationRequest struct {
	UserID   uint `json:"user_id" validate:"required"`
	CategoryID uint `json:"category_id,omitempty"`
}

type RecommendationResponse struct {
	ID              uint    `json:"id"`
	Card            CreditCard `json:"card"`
	Category        Category `json:"category"`
	Score           float64 `json:"score"`
	EstimatedReward float64 `json:"estimated_reward"`
	Reason          string  `json:"reason"`
} 