import React from 'react';
import { useQuery } from 'react-query';
import { Link } from 'react-router-dom';
import { categoryAPI, creditCardAPI } from '../services/api';

const Dashboard: React.FC = () => {
  const { data: categoriesResponse, isLoading: categoriesLoading } = useQuery(
    'categories',
    categoryAPI.getAll
  );

  const { data: cardsResponse, isLoading: cardsLoading } = useQuery(
    'credit-cards',
    creditCardAPI.getAll
  );

  const categories = categoriesResponse?.categories || [];
  const cards = cardsResponse?.cards || [];
  const isLoading = categoriesLoading || cardsLoading;

  if (isLoading) {
    return (
      <div className="max-w-6xl mx-auto">
        <h1 className="text-3xl font-bold text-gray-900 mb-6">Dashboard</h1>
        <div className="flex justify-center py-12">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        </div>
      </div>
    );
  }

  // Calculate statistics
  const freeCards = cards.filter(card => card.annual_fee === 0).length;
  const premiumCards = cards.filter(card => card.annual_fee > 0).length;

  return (
    <div className="max-w-6xl mx-auto">
      <h1 className="text-3xl font-bold text-gray-900 mb-6">Dashboard</h1>
      
      {/* Statistics Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center">
            <div className="p-2 bg-blue-100 rounded-lg">
              <svg className="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z"></path>
              </svg>
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">Total Cards</p>
              <p className="text-2xl font-bold text-gray-900">{cards.length}</p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center">
            <div className="p-2 bg-green-100 rounded-lg">
              <svg className="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"></path>
              </svg>
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">Free Cards</p>
              <p className="text-2xl font-bold text-gray-900">{freeCards}</p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center">
            <div className="p-2 bg-purple-100 rounded-lg">
              <svg className="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z"></path>
              </svg>
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">Premium Cards</p>
              <p className="text-2xl font-bold text-gray-900">{premiumCards}</p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center">
            <div className="p-2 bg-orange-100 rounded-lg">
              <svg className="w-6 h-6 text-orange-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"></path>
              </svg>
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">Categories</p>
              <p className="text-2xl font-bold text-gray-900">{categories.length}</p>
            </div>
          </div>
        </div>
      </div>

      {/* Quick Actions */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-bold text-gray-900 mb-4">Quick Actions</h2>
          <div className="space-y-3">
            <Link to="/spending" className="w-full text-left px-4 py-3 bg-blue-50 hover:bg-blue-100 rounded-lg transition-colors duration-200">
              <div className="flex items-center">
                <div className="p-2 bg-blue-200 rounded mr-3">
                  <svg className="w-4 h-4 text-blue-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
                  </svg>
                </div>
                <div>
                  <p className="font-medium text-gray-900">Add Spending</p>
                  <p className="text-sm text-gray-600">Track your monthly expenses</p>
                </div>
              </div>
            </Link>
            
            <Link to="/recommendations" className="w-full text-left px-4 py-3 bg-green-50 hover:bg-green-100 rounded-lg transition-colors duration-200">
              <div className="flex items-center">
                <div className="p-2 bg-green-200 rounded mr-3">
                  <svg className="w-4 h-4 text-green-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"></path>
                  </svg>
                </div>
                <div>
                  <p className="font-medium text-gray-900">Get Recommendations</p>
                  <p className="text-sm text-gray-600">Find the best cards for you</p>
                </div>
              </div>
            </Link>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-bold text-gray-900 mb-4">Spending Categories</h2>
          <div className="grid grid-cols-2 gap-3">
            {categories.slice(0, 8).map((category) => (
              <div key={category.id} className="flex items-center p-3 bg-gray-50 rounded-lg">
                <span className="text-2xl mr-3">{category.icon}</span>
                <div>
                  <p className="font-medium text-gray-900 text-sm">{category.name}</p>
                  <p className="text-xs text-gray-600">{category.description}</p>
                </div>
              </div>
            ))}
          </div>
          {categories.length > 8 && (
            <p className="text-sm text-gray-500 mt-3">
              +{categories.length - 8} more categories
            </p>
          )}
        </div>
      </div>

      {/* Recent Cards */}
      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-xl font-bold text-gray-900 mb-4">Featured Credit Cards</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {cards.slice(0, 3).map((card) => (
            <div key={card.id} className="border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow duration-200">
              <div className="flex items-start justify-between mb-2">
                <div>
                  <h3 className="font-semibold text-gray-900">{card.name}</h3>
                  <p className="text-sm text-gray-600">{card.bank}</p>
                </div>
                <span className={`px-2 py-1 text-xs rounded font-medium ${
                  card.annual_fee === 0 ? 'bg-green-100 text-green-800' : 'bg-orange-100 text-orange-800'
                }`}>
                  {card.annual_fee === 0 ? 'Free' : `$${card.annual_fee}`}
                </span>
              </div>
              <p className="text-sm text-gray-600 mb-3">{card.description}</p>
              <div className="text-xs text-gray-500">
                Min Income: ${card.min_income.toLocaleString()}
              </div>
            </div>
          ))}
        </div>
        {cards.length > 3 && (
          <div className="mt-4 text-center">
            <Link to="/cards" className="text-blue-600 hover:text-blue-700 font-medium text-sm">
              View All Cards ({cards.length})
            </Link>
          </div>
        )}
      </div>
    </div>
  );
};

export default Dashboard; 