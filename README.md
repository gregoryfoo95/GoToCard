# GoToCard - Credit Card Recommendation App

A comprehensive full-stack application that helps users find the best credit cards for their spending patterns using data scraped from SingSaver and MoneySmart Singapore.

## Architecture

### Backend (Golang)
- **Clean Architecture**: Repository → Service → Controller pattern
- **GORM**: PostgreSQL ORM with auto-migrations
- **Gin**: HTTP web framework with middleware
- **Colly**: Web scraping for card data from SingSaver and MoneySmart
- **Input Validation**: Server-side validation using go-playground/validator
- **SQL Injection Protection**: Parameterized queries via GORM

### Frontend (React TypeScript)
- **Modern React**: Hooks, TypeScript, and functional components
- **React Query**: Data fetching and caching
- **React Hook Form**: Form validation with Yup schemas
- **Tailwind CSS**: Utility-first styling
- **React Router**: Client-side routing
- **Axios**: HTTP client with interceptors

### Database (PostgreSQL)
- **Relational Design**: Users, Categories, Credit Cards, Card Benefits, User Spending, Recommendations
- **Foreign Key Constraints**: Data integrity
- **Indexing**: Optimized queries

## Features

- ✅ User management with email validation
- ✅ Spending category tracking
- ✅ Credit card database with benefits
- ✅ Web scraping from Singapore financial sites
- ✅ Smart recommendation algorithm
- ✅ RESTful API with comprehensive validation
- ✅ Responsive modern UI
- ✅ Docker containerization

## Getting Started

### Prerequisites
- Docker and Docker Compose
- Node.js 18+ (for local development)
- Go 1.21+ (for local development)

### Quick Start with Docker

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd GoToCard
   ```

2. **Start all services**
   ```bash
   docker-compose up --build
   ```

3. **Access the applications**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Database: localhost:5432

### Local Development

#### Backend Setup
```bash
cd backend
go mod download
go run cmd/main.go
```

#### Frontend Setup
```bash
cd frontend
npm install
npm start
```

## API Endpoints

### Users
- `POST /api/v1/users` - Create user
- `GET /api/v1/users` - List users
- `GET /api/v1/users/{id}` - Get user by ID

### Categories
- `GET /api/v1/categories` - List spending categories
- `POST /api/v1/categories` - Create category

### Credit Cards
- `GET /api/v1/cards` - List all credit cards
- `GET /api/v1/cards/{id}` - Get card details

### Spending
- `POST /api/v1/users/{userId}/spending` - Add spending record
- `GET /api/v1/users/{userId}/spending` - Get user spending

### Recommendations
- `POST /api/v1/users/{userId}/recommendations/generate` - Generate recommendations
- `GET /api/v1/users/{userId}/recommendations` - Get saved recommendations

### Admin
- `POST /api/v1/admin/scrape` - Trigger card data scraping

## Database Schema

### Core Tables
- `users`: User accounts
- `categories`: Spending categories (Dining, Groceries, etc.)
- `credit_cards`: Credit card details
- `card_benefits`: Card benefits per category
- `user_spending`: User spending records
- `recommendations`: Generated recommendations

## Recommendation Algorithm

The system analyzes user spending patterns and calculates the best credit cards based on:
- **Cashback/Points/Miles rates** per category
- **Annual fees** vs. estimated rewards
- **Spending caps** and minimum requirements
- **Net benefit calculation** over 12 months

## Security Features

- Input validation on both frontend and backend
- SQL injection prevention via GORM
- CORS configuration
- Request/Response interceptors
- Environment variable configuration

## Development Features

- **Hot Reload**: Backend (Air) and Frontend (React)
- **Type Safety**: TypeScript throughout
- **Error Handling**: Comprehensive error responses
- **Logging**: Structured logging
- **Validation**: Form validation with error messages

## Environment Variables

### Backend
```env
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=gotocard
JWT_SECRET=your-secret-key
SERVER_PORT=8080
```

### Frontend
```env
REACT_APP_API_URL=http://localhost:8080
```

## Testing

### Backend Testing
```bash
cd backend
go test ./...
```

### Frontend Testing
```bash
cd frontend
npm test
```

## Deployment

The application is containerized and can be deployed using:
```bash
docker-compose -f docker-compose.prod.yml up --build
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License.

## Future Enhancements

- [ ] User authentication with JWT
- [ ] Card comparison features
- [ ] Spending analytics dashboard
- [ ] Mobile app
- [ ] Machine learning recommendations
- [ ] Integration with bank APIs
- [ ] Real-time card updates
- [ ] Social features and reviews