export interface User {
  id: number;
  name: string;
  email: string;
  created_at: string;
  updated_at: string;
}

export interface Category {
  id: number;
  name: string;
  description: string;
  icon: string;
  created_at: string;
  updated_at: string;
}

export interface CreditCard {
  id: number;
  name: string;
  bank: string;
  card_type: string;
  annual_fee: number;
  image_url?: string;
  description: string;
  min_income: number;
  welcome_bonus?: string;
  is_active: boolean;
  card_benefits?: CardBenefit[];
  created_at: string;
  updated_at: string;
}

export interface CardBenefit {
  id: number;
  card_id: number;
  category_id: number;
  cashback_rate: number;
  points_rate: number;
  miles_rate: number;
  cap: number;
  min_spend: number;
  description: string;
  category: Category;
  created_at: string;
  updated_at: string;
}

export interface UserSpending {
  id: number;
  user_id: number;
  category_id: number;
  amount: number;
  month: number;
  year: number;
  category: Category;
  created_at: string;
  updated_at: string;
}

export interface Recommendation {
  id: number;
  card: CreditCard;
  category: Category;
  score: number;
  estimated_reward: number;
  reason: string;
}

// Request DTOs
export interface CreateUserRequest {
  name: string;
  email: string;
}

export interface SpendingRequest {
  category_id: number;
  amount: number;
  month: number;
  year: number;
}

// API Response types
export interface ApiResponse<T> {
  message?: string;
  data?: T;
  error?: string;
}

export interface UserResponse extends ApiResponse<User> {
  user?: User;
}

export interface UsersResponse extends ApiResponse<User[]> {
  users?: User[];
}

export interface CategoriesResponse extends ApiResponse<Category[]> {
  categories?: Category[];
}

export interface CardsResponse extends ApiResponse<CreditCard[]> {
  cards?: CreditCard[];
}

export interface SpendingsResponse extends ApiResponse<UserSpending[]> {
  spendings?: UserSpending[];
}

export interface RecommendationsResponse extends ApiResponse<Recommendation[]> {
  recommendations?: Recommendation[];
} 