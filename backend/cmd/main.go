package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"gotocard-backend/internal/config"
	"gotocard-backend/internal/controller"
	"gotocard-backend/internal/models"
	"gotocard-backend/internal/repository"
	"gotocard-backend/internal/service"
	"gotocard-backend/pkg/validator"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Database connection
	db, err := cfg.ConnectDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Connected to database")

	// Auto-migrate database
	if err := db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.CreditCard{},
		&models.CardBenefit{},
		&models.UserSpending{},
		&models.Recommendation{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed")

	// Clean existing sample data and seed fresh real data
	cleanAndSeedRealData(db)

	// Initialize repositories, services, and controllers
	repos := repository.NewRepositories(db)
	services := service.NewServices(repos)
	v := validator.NewValidator()
	controllers := controller.NewControllers(services, v)

	// Setup routes
	router := setupRoutes(controllers)

	// Start server
	log.Printf("Server starting on %s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Fatal(router.Run(cfg.Server.Host + ":" + cfg.Server.Port))
}

func setupRoutes(controllers *controller.Controllers) *gin.Engine {
	router := gin.Default()

	// Add CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := router.Group("/api/v1")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// User routes
		api.POST("/users", controllers.User.CreateUser)
		api.GET("/users", controllers.User.ListUsers)
		api.GET("/users/:id", controllers.User.GetUser)

		// Category routes
		api.POST("/categories", controllers.Category.CreateCategory)
		api.GET("/categories", controllers.Category.ListCategories)

		// Credit card routes
		api.GET("/cards", controllers.CreditCard.ListCreditCards)
		api.GET("/cards/:id", controllers.CreditCard.GetCreditCard)

		// Spending routes
		api.POST("/spending/users/:userId", controllers.Spending.AddSpending)
		api.GET("/spending/users/:userId", controllers.Spending.GetUserSpending)

		// Recommendation routes
		api.POST("/recommendations/users/:userId/generate", controllers.Recommendation.GenerateRecommendations)
		api.GET("/recommendations/users/:userId", controllers.Recommendation.GetRecommendations)

		// Admin routes
		admin := api.Group("/admin")
		{
			admin.POST("/scrape", controllers.Scraping.ScrapeCardData)
		}
	}

	return router
}

func cleanAndSeedRealData(db *gorm.DB) {
	log.Println("Cleaning existing data including curated cards...")

	// Clean all existing data in proper order (respecting foreign key constraints)
	if err := db.Exec("DELETE FROM recommendations").Error; err != nil {
		log.Printf("Warning: Error cleaning recommendations: %v", err)
	}

	if err := db.Exec("DELETE FROM user_spendings").Error; err != nil {
		log.Printf("Warning: Error cleaning user_spendings: %v", err)
	}

	if err := db.Exec("DELETE FROM card_benefits").Error; err != nil {
		log.Printf("Warning: Error cleaning card_benefits: %v", err)
	}

	if err := db.Exec("DELETE FROM credit_cards").Error; err != nil {
		log.Printf("Warning: Error cleaning credit_cards: %v", err)
	}

	if err := db.Exec("DELETE FROM users").Error; err != nil {
		log.Printf("Warning: Error cleaning users: %v", err)
	}

	// Reset auto-increment counters
	sequences := []string{
		"credit_cards_id_seq",
		"card_benefits_id_seq",
		"users_id_seq",
		"user_spendings_id_seq",
		"recommendations_id_seq",
	}

	for _, seq := range sequences {
		if err := db.Exec(fmt.Sprintf("ALTER SEQUENCE %s RESTART WITH 1", seq)).Error; err != nil {
			log.Printf("Warning: Error resetting %s sequence: %v", seq, err)
		}
	}

	log.Println("All existing data cleanup completed - ready for fresh scraped data only")

	// Seed only essential categories (these should remain)
	seedCategories(db)

	log.Println("Categories seeded")

	// Use scraping service to populate real card data (no more curated data)
	repos := repository.NewRepositories(db)
	scrapingService := service.NewScrapingService(repos)

	log.Println("Starting real data scraping from web sources only...")
	if err := scrapingService.ScrapeCardData(); err != nil {
		log.Printf("Warning: Scraping failed: %v", err)
		log.Println("Will continue with available data...")
	} else {
		log.Println("Real data scraping completed successfully")
	}

	// Create a demo user for testing
	createDemoUser(db)
}

func seedCategories(db *gorm.DB) {
	log.Println("Seeding initial categories...")

	categories := []models.Category{
		{Name: "Dining", Description: "Restaurant meals, food delivery, cafes"},
		{Name: "Groceries", Description: "Supermarket purchases, food shopping"},
		{Name: "Petrol", Description: "Fuel, gas stations"},
		{Name: "Shopping", Description: "Retail purchases, department stores"},
		{Name: "Transport", Description: "Public transport, ride-hailing, parking"},
		{Name: "Travel", Description: "Hotels, flights, vacation expenses"},
		{Name: "Entertainment", Description: "Movies, concerts, gaming, streaming"},
		{Name: "Healthcare", Description: "Medical expenses, pharmacy, dental"},
		{Name: "Bills", Description: "Utilities, phone, internet, insurance"},
		{Name: "Online", Description: "E-commerce, online subscriptions"},
	}

	for _, category := range categories {
		var existingCategory models.Category
		result := db.Where("name = ?", category.Name).First(&existingCategory)
		if result.Error != nil {
			if err := db.Create(&category).Error; err != nil {
				log.Printf("Failed to create category %s: %v", category.Name, err)
			} else {
				log.Printf("Created category: %s", category.Name)
			}
		}
	}

	log.Println("Categories seeding completed")
}

func createDemoUser(db *gorm.DB) {
	log.Println("Creating demo user...")

	demoUser := &models.User{
		Name:  "Demo User",
		Email: "demo@gotocard.sg",
	}

	// Check if demo user already exists
	var existingUser models.User
	result := db.Where("email = ?", demoUser.Email).First(&existingUser)
	if result.Error != nil {
		if err := db.Create(demoUser).Error; err != nil {
			log.Printf("Failed to create demo user: %v", err)
		} else {
			log.Printf("Created demo user: %s", demoUser.Name)
		}
	} else {
		log.Printf("Demo user already exists: %s", existingUser.Name)
	}
}
