package service

import (
	"fmt"
	"math"
	"sort"

	"gotocard-backend/internal/models"
	"gotocard-backend/internal/repository"
)

type recommendationService struct {
	repos *repository.Repositories
}

func NewRecommendationService(repos *repository.Repositories) RecommendationService {
	return &recommendationService{repos: repos}
}

func (s *recommendationService) GenerateRecommendations(userID uint) ([]models.RecommendationResponse, error) {
	// Get user spending data
	spendings, err := s.repos.UserSpending.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user spending: %w", err)
	}

	// Get all active credit cards
	cards, err := s.repos.CreditCard.GetActiveCards()
	if err != nil {
		return nil, fmt.Errorf("failed to get active cards: %w", err)
	}

	// Calculate spending by category
	categorySpending := make(map[uint]float64)
	for _, spending := range spendings {
		categorySpending[spending.CategoryID] += spending.Amount
	}

	// Generate recommendations for each category with spending
	var recommendations []models.RecommendationResponse
	for categoryID, totalSpent := range categorySpending {
		categoryRecs := s.calculateBestCardsForCategory(categoryID, totalSpent, cards)
		recommendations = append(recommendations, categoryRecs...)
	}

	// Sort by score descending
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	// Save top recommendations to database
	err = s.saveRecommendations(userID, recommendations)
	if err != nil {
		return nil, fmt.Errorf("failed to save recommendations: %w", err)
	}

	// Return top 10 recommendations
	if len(recommendations) > 10 {
		recommendations = recommendations[:10]
	}

	return recommendations, nil
}

func (s *recommendationService) calculateBestCardsForCategory(categoryID uint, monthlySpent float64, cards []models.CreditCard) []models.RecommendationResponse {
	var recommendations []models.RecommendationResponse

	category, err := s.repos.Category.GetByID(categoryID)
	if err != nil {
		return recommendations
	}

	for _, card := range cards {
		// Find benefits for this category
		var bestBenefit *models.CardBenefit
		for _, benefit := range card.CardBenefits {
			if benefit.CategoryID == categoryID {
				bestBenefit = &benefit
				break
			}
		}

		if bestBenefit == nil {
			continue // No benefits for this category
		}

		// Calculate expected reward
		reward := s.calculateReward(monthlySpent, bestBenefit)
		
		// Calculate score (considering annual fee)
		annualReward := reward * 12
		netBenefit := annualReward - card.AnnualFee
		score := s.calculateScore(netBenefit, card.AnnualFee, bestBenefit)

		// Generate reason
		reason := s.generateReason(bestBenefit, reward, card.AnnualFee)

		recommendations = append(recommendations, models.RecommendationResponse{
			Card:            card,
			Category:        *category,
			Score:           score,
			EstimatedReward: reward,
			Reason:          reason,
		})
	}

	return recommendations
}

func (s *recommendationService) calculateReward(monthlySpent float64, benefit *models.CardBenefit) float64 {
	if monthlySpent < benefit.MinSpend {
		return 0
	}

	spentAmount := monthlySpent
	if benefit.Cap > 0 && monthlySpent > benefit.Cap {
		spentAmount = benefit.Cap
	}

	// Calculate based on best rate available
	if benefit.CashbackRate > 0 {
		return spentAmount * (benefit.CashbackRate / 100)
	}
	
	if benefit.PointsRate > 0 {
		// Assume 1 point = $0.01 (can be adjusted)
		return spentAmount * (benefit.PointsRate / 100) * 0.01
	}
	
	if benefit.MilesRate > 0 {
		// Assume 1 mile = $0.015 (can be adjusted)
		return spentAmount * (benefit.MilesRate / 100) * 0.015
	}

	return 0
}

func (s *recommendationService) calculateScore(netBenefit, annualFee float64, benefit *models.CardBenefit) float64 {
	// Base score from net benefit
	score := netBenefit

	// Bonus for higher reward rates
	if benefit.CashbackRate > 0 {
		score += benefit.CashbackRate * 10
	}
	if benefit.PointsRate > 0 {
		score += benefit.PointsRate * 8
	}
	if benefit.MilesRate > 0 {
		score += benefit.MilesRate * 12
	}

	// Penalty for annual fee
	score -= annualFee * 0.5

	// Normalize score to 0-100 range
	score = math.Max(0, math.Min(100, score))
	
	return math.Round(score*100) / 100
}

func (s *recommendationService) generateReason(benefit *models.CardBenefit, monthlyReward, annualFee float64) string {
	annualReward := monthlyReward * 12
	netBenefit := annualReward - annualFee

	reason := fmt.Sprintf("Earn %.2f%% ", benefit.CashbackRate)
	if benefit.PointsRate > 0 {
		reason = fmt.Sprintf("Earn %.1fx points ", benefit.PointsRate)
	}
	if benefit.MilesRate > 0 {
		reason = fmt.Sprintf("Earn %.1fx miles ", benefit.MilesRate)
	}

	reason += fmt.Sprintf("on this category. Expected monthly reward: $%.2f", monthlyReward)
	
	if annualFee > 0 {
		reason += fmt.Sprintf(", Annual fee: $%.0f", annualFee)
	}
	
	reason += fmt.Sprintf(", Net annual benefit: $%.2f", netBenefit)

	return reason
}

func (s *recommendationService) saveRecommendations(userID uint, recommendations []models.RecommendationResponse) error {
	// Clear existing recommendations
	err := s.repos.Recommendation.DeleteByUserID(userID)
	if err != nil {
		return err
	}

	// Save new recommendations (top 10)
	count := len(recommendations)
	if count > 10 {
		count = 10
	}

	for i := 0; i < count; i++ {
		rec := &models.Recommendation{
			UserID:          userID,
			CategoryID:      recommendations[i].Category.ID,
			CardID:          recommendations[i].Card.ID,
			Score:           recommendations[i].Score,
			EstimatedReward: recommendations[i].EstimatedReward,
			Reason:          recommendations[i].Reason,
		}
		
		err := s.repos.Recommendation.Create(rec)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *recommendationService) GetRecommendationsByUser(userID uint) ([]models.RecommendationResponse, error) {
	recs, err := s.repos.Recommendation.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []models.RecommendationResponse
	for _, rec := range recs {
		responses = append(responses, models.RecommendationResponse{
			ID:              rec.ID,
			Card:            rec.Card,
			Category:        rec.Category,
			Score:           rec.Score,
			EstimatedReward: rec.EstimatedReward,
			Reason:          rec.Reason,
		})
	}

	return responses, nil
}

func (s *recommendationService) GetRecommendationsByCategory(userID, categoryID uint) ([]models.RecommendationResponse, error) {
	recs, err := s.repos.Recommendation.GetByUserAndCategory(userID, categoryID)
	if err != nil {
		return nil, err
	}

	var responses []models.RecommendationResponse
	for _, rec := range recs {
		responses = append(responses, models.RecommendationResponse{
			ID:              rec.ID,
			Card:            rec.Card,
			Category:        rec.Category,
			Score:           rec.Score,
			EstimatedReward: rec.EstimatedReward,
			Reason:          rec.Reason,
		})
	}

	return responses, nil
}

func (s *recommendationService) RefreshRecommendations(userID uint) error {
	_, err := s.GenerateRecommendations(userID)
	return err
} 