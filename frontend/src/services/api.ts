import axios from 'axios';
import {
  Category,
  CreditCard,
  CreateUserRequest,
  SpendingRequest,
  UserResponse,
  UsersResponse,
  CategoriesResponse,
  CardsResponse,
  SpendingsResponse,
  RecommendationsResponse
} from '../types';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: `${API_BASE_URL}/api/v1`,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor for adding auth tokens if needed
api.interceptors.request.use(
  (config) => {
    // Add auth token if available
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error('API Error:', error.response?.data || error.message);
    return Promise.reject(error);
  }
);

// User API
export const userAPI = {
  create: (userData: CreateUserRequest): Promise<UserResponse> =>
    api.post('/users', userData).then(res => res.data),
    
  getById: (id: number): Promise<UserResponse> =>
    api.get(`/users/${id}`).then(res => res.data),
    
  getAll: (): Promise<UsersResponse> =>
    api.get('/users').then(res => res.data),
};

// Category API
export const categoryAPI = {
  getAll: (): Promise<CategoriesResponse> =>
    api.get('/categories').then(res => res.data),
    
  create: (categoryData: Partial<Category>): Promise<{ category: Category }> =>
    api.post('/categories', categoryData).then(res => res.data),
};

// Credit Card API
export const creditCardAPI = {
  getAll: (): Promise<CardsResponse> =>
    api.get('/cards').then(res => res.data),
    
  getById: (id: number): Promise<{ card: CreditCard }> =>
    api.get(`/cards/${id}`).then(res => res.data),
};

// Spending API
export const spendingAPI = {
  add: (userId: number, spendingData: SpendingRequest): Promise<{ message: string }> =>
    api.post(`/spending/users/${userId}`, spendingData).then(res => res.data),
    
  getUserSpending: (userId: number): Promise<SpendingsResponse> =>
    api.get(`/spending/users/${userId}`).then(res => res.data),
};

// Recommendation API
export const recommendationAPI = {
  generate: (userId: number): Promise<RecommendationsResponse> =>
    api.post(`/recommendations/users/${userId}/generate`).then(res => res.data),
    
  getByUser: (userId: number): Promise<RecommendationsResponse> =>
    api.get(`/recommendations/users/${userId}`).then(res => res.data),
};

// Admin API
export const adminAPI = {
  scrapeCards: (): Promise<{ message: string }> =>
    api.post('/admin/scrape').then(res => res.data),
};

export default api; 