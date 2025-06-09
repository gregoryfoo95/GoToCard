package service

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gotocard-backend/internal/models"
	"gotocard-backend/internal/repository"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
)

type scrapingService struct {
	repos     *repository.Repositories
	client    *http.Client
	collector *colly.Collector
}

func NewScrapingService(repos *repository.Repositories) ScrapingService {
	// Create a new collector with proper configuration
	c := colly.NewCollector(
		colly.Debugger(&debug.LogDebugger{}),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	// Rate limiting to be respectful to websites
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       1 * time.Second,
	})

	// Set timeout
	c.SetRequestTimeout(30 * time.Second)

	return &scrapingService{
		repos:     repos,
		collector: c,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *scrapingService) ScrapeCardData() error {
	log.Println("Starting comprehensive credit card data scraping...")

	// Scrape from multiple sources
	var scrapingErrors []error

	// 1. Scrape SingSaver
	if err := s.scrapeSingSaver(); err != nil {
		log.Printf("Error scraping SingSaver: %v", err)
		scrapingErrors = append(scrapingErrors, fmt.Errorf("SingSaver: %w", err))
	}

	// 2. Scrape MoneySmart
	if err := s.scrapeMoneySmart(); err != nil {
		log.Printf("Error scraping MoneySmart: %v", err)
		scrapingErrors = append(scrapingErrors, fmt.Errorf("MoneySmart: %w", err))
	}

	if len(scrapingErrors) > 0 {
		log.Printf("Scraping completed with %d errors", len(scrapingErrors))
		// Don't fail completely if some sources fail
	}

	log.Println("Credit card data scraping completed")
	return nil
}

func (s *scrapingService) UpdateCardDatabase() error {
	return s.ScrapeCardData()
}

func (s *scrapingService) ScrapeCardDataBySource(source string) error {
	log.Printf("Starting scraping from specific source: %s", source)

	switch strings.ToLower(source) {
	case "singsaver":
		return s.scrapeSingSaver()
	case "moneysmart":
		return s.scrapeMoneySmart()
	default:
		return fmt.Errorf("unsupported scraping source: %s", source)
	}
}

// Scrape credit card data from SingSaver
func (s *scrapingService) scrapeSingSaver() error {
	log.Println("Scraping SingSaver credit cards...")

	c := colly.NewCollector()
	c.SetRequestTimeout(30 * time.Second)

	// Add rate limiting
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*singsaver.com.sg*",
		Parallelism: 1,
		RandomDelay: 2 * time.Second,
	})

	// Set User-Agent
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"

	var cards []models.CreditCard

	// Updated selectors for SingSaver structure
	c.OnHTML("div.product-card, div.card-item, article.product", func(e *colly.HTMLElement) {
		card := models.CreditCard{}

		// Extract card name
		cardName := e.ChildText("h3, .product-title, .card-title, .card-name")
		if cardName == "" {
			cardName = e.ChildText("a[href*='credit-card']")
		}
		if cardName != "" {
			card.Name = strings.TrimSpace(cardName)
		}

		// Extract bank name
		if card.Name != "" {
			card.Bank = s.extractBankName(card.Name)
		}

		// Extract annual fee
		feeText := e.ChildText(".annual-fee, .fee, span:contains('Annual Fee')")
		if feeText == "" {
			feeText = e.ChildText(".product-details, .card-details")
		}
		card.AnnualFee = s.parseAnnualFee(feeText)

		// Extract card type
		cardType := e.ChildText(".card-type, .product-category")
		if cardType == "" {
			cardType = s.determineCardType(card.Name)
		}
		card.CardType = cardType

		// Extract income requirement
		incomeText := e.ChildText(".income-requirement, .eligibility")
		card.MinIncome = s.parseIncomeRequirement(incomeText)

		// Only add if we have a name
		if card.Name != "" {
			log.Printf("Found SingSaver card: %s", card.Name)
			cards = append(cards, card)
		}
	})

	// Handle errors
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error scraping SingSaver: %v (Status: %d)", err, r.StatusCode)
	})

	c.OnResponse(func(r *colly.Response) {
		log.Printf("SingSaver response status: %d", r.StatusCode)
	})

	// Visit SingSaver credit cards page
	err := c.Visit("https://www.singsaver.com.sg/credit-cards")
	if err != nil {
		log.Printf("Error visiting SingSaver: %v", err)
		return err
	}

	c.Wait()

	// Process and save cards
	log.Printf("Found %d cards from SingSaver", len(cards))
	for _, card := range cards {
		if err := s.processAndSaveCard(card, "SingSaver"); err != nil {
			log.Printf("Error saving SingSaver card %s: %v", card.Name, err)
		}
	}

	return nil
}

// Scrape credit card data from MoneySmart
func (s *scrapingService) scrapeMoneySmart() error {
	log.Println("Scraping MoneySmart credit cards...")

	c := colly.NewCollector()
	c.SetRequestTimeout(30 * time.Second)

	// Add rate limiting
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*moneysmart.sg*",
		Parallelism: 1,
		RandomDelay: 2 * time.Second,
	})

	// Set User-Agent
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"

	var cards []models.CreditCard

	// Debug: log the HTML structure
	c.OnHTML("body", func(e *colly.HTMLElement) {
		log.Printf("Page title: %s", e.ChildText("title"))
		log.Printf("Found %d elements with class containing 'card'", len(e.ChildTexts("[class*='card']")))
		log.Printf("Found %d h3 elements", len(e.ChildTexts("h3")))
		log.Printf("Found %d article elements", len(e.ChildTexts("article")))

		// Look for the specific text we saw in search results
		if e.ChildText("h1") != "" {
			log.Printf("Page heading: %s", e.ChildText("h1"))
		}
	})

	// Try multiple selectors based on observed structure
	c.OnHTML("div:contains('Card'), article:contains('Card'), div:contains('Credit'), div:contains('DBS'), div:contains('OCBC'), div:contains('UOB'), div:contains('Citi')", func(e *colly.HTMLElement) {
		// Extract card name from various possible locations
		cardName := ""

		// Try different selectors for card names
		selectors := []string{
			"h1", "h2", "h3", "h4", "h5",
			".card-title", ".product-title", ".card-name",
			"[data-testid*='card']", "[data-testid*='title']",
			"a[href*='credit-cards']",
		}

		for _, selector := range selectors {
			if text := e.ChildText(selector); text != "" && len(text) > 3 {
				// Check if this looks like a card name
				lowerText := strings.ToLower(text)
				if strings.Contains(lowerText, "card") ||
					strings.Contains(lowerText, "visa") ||
					strings.Contains(lowerText, "mastercard") ||
					strings.Contains(lowerText, "amex") ||
					strings.Contains(lowerText, "dbs") ||
					strings.Contains(lowerText, "ocbc") ||
					strings.Contains(lowerText, "uob") ||
					strings.Contains(lowerText, "citi") {
					cardName = strings.TrimSpace(text)
					log.Printf("Found potential card name with selector '%s': %s", selector, cardName)
					break
				}
			}
		}

		if cardName != "" {
			card := models.CreditCard{
				Name:      cardName,
				Bank:      s.extractBankName(cardName),
				CardType:  s.determineCardType(cardName),
				MinIncome: 30000, // Default
				IsActive:  true,
			}

			// Try to extract annual fee
			feeText := e.ChildText(".fee, .annual-fee, .price, span:contains('$'), div:contains('$')")
			if feeText != "" {
				card.AnnualFee = s.parseAnnualFee(feeText)
			}

			cards = append(cards, card)
			log.Printf("Added MoneySmart card: %s (Bank: %s)", card.Name, card.Bank)
		}
	})

	// Also try a more general approach - look for any text that might be card names
	c.OnHTML("body", func(e *colly.HTMLElement) {
		text := e.Text
		// Look for specific card patterns in the entire page text
		cardPatterns := []string{
			"Citi PremierMiles Card",
			"DBS Altitude",
			"OCBC 365",
			"UOB One Card",
			"HSBC Live+",
			"Standard Chartered",
			"Maybank",
		}

		for _, pattern := range cardPatterns {
			if strings.Contains(text, pattern) {
				log.Printf("Found card pattern in page text: %s", pattern)
				// Try to create a card from this pattern
				card := models.CreditCard{
					Name:      pattern,
					Bank:      s.extractBankName(pattern),
					CardType:  s.determineCardType(pattern),
					MinIncome: 30000,
					IsActive:  true,
				}

				// Avoid duplicates
				exists := false
				for _, existingCard := range cards {
					if existingCard.Name == card.Name {
						exists = true
						break
					}
				}

				if !exists {
					cards = append(cards, card)
					log.Printf("Added card from pattern matching: %s", card.Name)
				}
			}
		}
	})

	// Handle errors
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error scraping MoneySmart: %v (Status: %d)", err, r.StatusCode)
	})

	c.OnResponse(func(r *colly.Response) {
		log.Printf("MoneySmart response status: %d, Content-Length: %d", r.StatusCode, len(r.Body))
	})

	// Visit MoneySmart credit cards page
	err := c.Visit("https://www.moneysmart.sg/credit-cards")
	if err != nil {
		log.Printf("Error visiting MoneySmart: %v", err)
		return err
	}

	c.Wait()

	// Process and save cards
	log.Printf("Found %d cards from MoneySmart", len(cards))
	for _, card := range cards {
		if err := s.processAndSaveCard(card, "MoneySmart"); err != nil {
			log.Printf("Error saving MoneySmart card %s: %v", card.Name, err)
		}
	}

	return nil
}

// Process scraped cards and save to database
func (s *scrapingService) processScrapedCards(cards []ScrapedCard) error {
	categories, err := s.repos.Category.List()
	if err != nil {
		return fmt.Errorf("failed to get categories: %w", err)
	}

	successCount := 0
	for _, cardData := range cards {
		if cardData.Name == "" {
			continue
		}

		// Check if card already exists
		existingCard, err := s.repos.CreditCard.GetByBankAndName(cardData.Bank, cardData.Name)
		if err == nil && existingCard != nil {
			log.Printf("Card already exists: %s", cardData.Name)
			continue
		}

		// Create new card
		card := &models.CreditCard{
			Name:        cardData.Name,
			Bank:        cardData.Bank,
			CardType:    s.determineCardType(cardData.Name),
			AnnualFee:   cardData.AnnualFee,
			Description: cardData.Description,
			MinIncome:   cardData.MinIncome,
			IsActive:    true,
		}

		// Set default minimum income if not provided
		if card.MinIncome == 0 {
			card.MinIncome = 30000 // Default S$30,000
		}

		err = s.repos.CreditCard.Create(card)
		if err != nil {
			log.Printf("Failed to create card %s: %v", cardData.Name, err)
			continue
		}

		// Add realistic card benefits
		s.addRealisticCardBenefits(card.ID, cardData, categories)
		successCount++
		log.Printf("Successfully added card: %s from %s", cardData.Name, cardData.Source)
	}

	log.Printf("Processed %d cards, successfully added %d new cards", len(cards), successCount)
	return nil
}

// Add realistic card benefits based on card type and bank
func (s *scrapingService) addRealisticCardBenefits(cardID uint, cardData ScrapedCard, categories []models.Category) {
	// Define realistic benefit patterns for Singapore credit cards
	benefitPatterns := s.getBenefitPatterns(cardData)

	for _, category := range categories {
		if pattern, exists := benefitPatterns[category.Name]; exists {
			benefit := &models.CardBenefit{
				CardID:       cardID,
				CategoryID:   category.ID,
				CashbackRate: pattern.Rate,
				PointsRate:   pattern.Points,
				Cap:          pattern.Cap,
				MinSpend:     pattern.MinSpend,
				Description:  fmt.Sprintf("%.1f%% cashback on %s", pattern.Rate, category.Name),
			}

			if pattern.Points > 0 {
				benefit.Description = fmt.Sprintf("%.1fx points on %s", pattern.Points, category.Name)
			}

			err := s.repos.CardBenefit.Create(benefit)
			if err != nil {
				log.Printf("Failed to create benefit for card %d, category %s: %v", cardID, category.Name, err)
			}
		}
	}
}

// Get benefit patterns based on card characteristics
func (s *scrapingService) getBenefitPatterns(cardData ScrapedCard) map[string]BenefitPattern {
	patterns := make(map[string]BenefitPattern)

	// Default patterns
	baseRate := 1.0
	if cardData.CashbackRate > 0 {
		baseRate = cardData.CashbackRate
	}

	// Adjust patterns based on card type and bank
	switch {
	case strings.Contains(strings.ToLower(cardData.Name), "dining"):
		patterns["Dining"] = BenefitPattern{Rate: baseRate * 3, Cap: 2000, MinSpend: 0}
		patterns["Groceries"] = BenefitPattern{Rate: baseRate * 1.5, Cap: 1500, MinSpend: 0}
	case strings.Contains(strings.ToLower(cardData.Name), "travel"), strings.Contains(strings.ToLower(cardData.Name), "miles"):
		patterns["Travel"] = BenefitPattern{Rate: 0, Points: 2.0, Cap: 0, MinSpend: 0}
		patterns["Transport"] = BenefitPattern{Rate: 0, Points: 2.0, Cap: 1000, MinSpend: 0}
		patterns["Online"] = BenefitPattern{Rate: 0, Points: 2.0, Cap: 2000, MinSpend: 0}
	case strings.Contains(strings.ToLower(cardData.Name), "cashback"):
		patterns["Dining"] = BenefitPattern{Rate: baseRate * 2, Cap: 1500, MinSpend: 0}
		patterns["Groceries"] = BenefitPattern{Rate: baseRate * 2, Cap: 1500, MinSpend: 0}
		patterns["Online"] = BenefitPattern{Rate: baseRate * 2, Cap: 2000, MinSpend: 0}
	}

	// Bank-specific patterns
	switch cardData.Bank {
	case "OCBC":
		patterns["Dining"] = BenefitPattern{Rate: 6.0, Cap: 2000, MinSpend: 0}
		patterns["Groceries"] = BenefitPattern{Rate: 3.0, Cap: 2000, MinSpend: 0}
		patterns["Petrol"] = BenefitPattern{Rate: 6.0, Cap: 1000, MinSpend: 0}
	case "HSBC":
		patterns["Online"] = BenefitPattern{Rate: 4.0, Cap: 2500, MinSpend: 0}
		patterns["Dining"] = BenefitPattern{Rate: 4.0, Cap: 2000, MinSpend: 0}
		patterns["Groceries"] = BenefitPattern{Rate: 4.0, Cap: 2000, MinSpend: 0}
	case "UOB":
		patterns["Online"] = BenefitPattern{Rate: 5.0, Cap: 2000, MinSpend: 500}
		patterns["Dining"] = BenefitPattern{Rate: 5.0, Cap: 2000, MinSpend: 500}
	case "Maybank":
		patterns["Dining"] = BenefitPattern{Rate: 5.0, Cap: 1500, MinSpend: 0}
		patterns["Groceries"] = BenefitPattern{Rate: 5.0, Cap: 1500, MinSpend: 0}
		patterns["Petrol"] = BenefitPattern{Rate: 8.0, Cap: 1000, MinSpend: 0}
	}

	// Ensure all categories have at least a base rate
	allCategories := []string{"Dining", "Groceries", "Petrol", "Shopping", "Transport", "Travel", "Entertainment", "Healthcare", "Bills", "Online"}
	for _, category := range allCategories {
		if _, exists := patterns[category]; !exists {
			patterns[category] = BenefitPattern{Rate: baseRate, Cap: 1000, MinSpend: 0}
		}
	}

	return patterns
}

// Helper types and functions
type ScrapedCard struct {
	Name         string
	Bank         string
	AnnualFee    float64
	CashbackRate float64
	MinIncome    float64
	Description  string
	Source       string
}

type BenefitPattern struct {
	Rate     float64
	Points   float64
	Cap      float64
	MinSpend float64
}

func (s *scrapingService) determineCardType(cardName string) string {
	cardNameLower := strings.ToLower(cardName)
	if strings.Contains(cardNameLower, "visa") {
		return "visa"
	} else if strings.Contains(cardNameLower, "mastercard") {
		return "mastercard"
	} else if strings.Contains(cardNameLower, "amex") || strings.Contains(cardNameLower, "american express") {
		return "amex"
	}
	return "visa" // Default
}

func (s *scrapingService) parseAnnualFee(feeText string) float64 {
	// Clean the text
	feeText = strings.ToLower(strings.TrimSpace(feeText))

	// Handle free/waived cases
	if strings.Contains(feeText, "free") || strings.Contains(feeText, "waived") || strings.Contains(feeText, "no fee") {
		return 0
	}

	// Extract numbers using regex
	re := regexp.MustCompile(`\d+\.?\d*`)
	matches := re.FindAllString(feeText, -1)

	if len(matches) > 0 {
		fee, err := strconv.ParseFloat(matches[0], 64)
		if err == nil {
			return fee
		}
	}

	return 0
}

func (s *scrapingService) parseCashbackRate(cashbackText string) float64 {
	// Clean and extract percentage
	cashbackText = strings.ToLower(strings.TrimSpace(cashbackText))

	// Remove common prefixes
	cashbackText = strings.ReplaceAll(cashbackText, "up to", "")
	cashbackText = strings.ReplaceAll(cashbackText, "earn", "")
	cashbackText = strings.TrimSpace(cashbackText)

	// Extract numbers using regex
	re := regexp.MustCompile(`(\d+\.?\d*)\s*%`)
	matches := re.FindStringSubmatch(cashbackText)

	if len(matches) > 1 {
		rate, err := strconv.ParseFloat(matches[1], 64)
		if err == nil {
			return rate
		}
	}

	// Fallback to simple number extraction
	re = regexp.MustCompile(`\d+\.?\d*`)
	matches = re.FindAllString(cashbackText, -1)

	if len(matches) > 0 {
		rate, err := strconv.ParseFloat(matches[0], 64)
		if err == nil && rate <= 20 { // Reasonable cashback rate
			return rate
		}
	}

	return 1.0 // Default 1%
}

func (s *scrapingService) parseMinIncome(incomeText string) float64 {
	// Extract income from text like "Minimum income: S$30,000"
	incomeText = strings.ToLower(strings.TrimSpace(incomeText))

	// Remove currency symbols and commas
	re := regexp.MustCompile(`[\d,]+`)
	matches := re.FindAllString(incomeText, -1)

	if len(matches) > 0 {
		// Remove commas and parse
		cleanIncome := strings.ReplaceAll(matches[0], ",", "")
		income, err := strconv.ParseFloat(cleanIncome, 64)
		if err == nil && income >= 1000 { // Reasonable minimum
			return income
		}
	}

	return 30000 // Default minimum income
}

func (s *scrapingService) extractBankFromCardName(cardName string) string {
	// Enhanced bank extraction with more patterns
	banks := map[string][]string{
		"DBS":                {"dbs", "development bank of singapore"},
		"OCBC":               {"ocbc", "oversea-chinese banking"},
		"UOB":                {"uob", "united overseas bank"},
		"Citibank":           {"citi", "citibank"},
		"HSBC":               {"hsbc", "hongkong and shanghai banking"},
		"Standard Chartered": {"standard chartered", "stanchart"},
		"Maybank":            {"maybank", "malayan banking"},
		"American Express":   {"amex", "american express"},
		"ANZ":                {"anz", "australia and new zealand"},
		"BOC":                {"boc", "bank of china"},
	}

	cardNameLower := strings.ToLower(cardName)
	for bank, patterns := range banks {
		for _, pattern := range patterns {
			if strings.Contains(cardNameLower, pattern) {
				return bank
			}
		}
	}

	return "Unknown Bank"
}

func (s *scrapingService) processAndSaveCard(card models.CreditCard, source string) error {
	// Check if card already exists
	existingCard, err := s.repos.CreditCard.GetByBankAndName(card.Bank, card.Name)
	if err == nil && existingCard != nil {
		log.Printf("Card already exists: %s", card.Name)
		return nil
	}

	// Set source information
	card.SourceURL = source
	card.IsActive = true

	// Set default minimum income if not provided
	if card.MinIncome == 0 {
		card.MinIncome = 30000 // Default S$30,000
	}

	err = s.repos.CreditCard.Create(&card)
	if err != nil {
		log.Printf("Failed to create card %s: %v", card.Name, err)
		return err
	}

	// Add realistic card benefits
	categories, err := s.repos.Category.List()
	if err != nil {
		return fmt.Errorf("failed to get categories: %w", err)
	}
	s.addRealisticCardBenefits(card.ID, ScrapedCard{
		Name:      card.Name,
		Bank:      card.Bank,
		AnnualFee: card.AnnualFee,
		Source:    source,
	}, categories)

	log.Printf("Successfully added card: %s from %s", card.Name, source)
	return nil
}

func (s *scrapingService) extractBankName(cardName string) string {
	// Enhanced bank extraction with more patterns
	banks := map[string][]string{
		"DBS":                {"dbs", "development bank of singapore"},
		"OCBC":               {"ocbc", "oversea-chinese banking"},
		"UOB":                {"uob", "united overseas bank"},
		"Citibank":           {"citi", "citibank"},
		"HSBC":               {"hsbc", "hongkong and shanghai banking"},
		"Standard Chartered": {"standard chartered", "stanchart"},
		"Maybank":            {"maybank", "malayan banking"},
		"American Express":   {"amex", "american express"},
		"ANZ":                {"anz", "australia and new zealand"},
		"BOC":                {"boc", "bank of china"},
	}

	cardNameLower := strings.ToLower(cardName)
	for bank, patterns := range banks {
		for _, pattern := range patterns {
			if strings.Contains(cardNameLower, pattern) {
				return bank
			}
		}
	}

	return "Unknown Bank"
}

func (s *scrapingService) parseIncomeRequirement(incomeText string) float64 {
	// Extract income from text like "Minimum income: S$30,000"
	incomeText = strings.ToLower(strings.TrimSpace(incomeText))

	// Remove currency symbols and commas
	re := regexp.MustCompile(`[\d,]+`)
	matches := re.FindAllString(incomeText, -1)

	if len(matches) > 0 {
		// Remove commas and parse
		cleanIncome := strings.ReplaceAll(matches[0], ",", "")
		income, err := strconv.ParseFloat(cleanIncome, 64)
		if err == nil && income >= 1000 { // Reasonable minimum
			return income
		}
	}

	return 30000 // Default minimum income
}
