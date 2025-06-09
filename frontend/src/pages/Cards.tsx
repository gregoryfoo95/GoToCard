import React from 'react';
import { useQuery } from 'react-query';
import { creditCardAPI } from '../services/api';
import { CreditCard } from '../types';

const Cards: React.FC = () => {
  const { data: cardsResponse, isLoading, error } = useQuery(
    'credit-cards',
    creditCardAPI.getAll,
    {
      staleTime: 5 * 60 * 1000, // 5 minutes
    }
  );

  if (isLoading) {
    return (
      <div className="max-w-6xl mx-auto">
        <h1 className="text-3xl font-bold text-gray-900 mb-6">Credit Cards</h1>
        <div className="flex justify-center py-12">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="max-w-6xl mx-auto">
        <h1 className="text-3xl font-bold text-gray-900 mb-6">Credit Cards</h1>
        <div className="bg-red-50 border border-red-200 rounded-lg p-4">
          <p className="text-red-600">Error loading credit cards. Please try again later.</p>
        </div>
      </div>
    );
  }

  const cards = cardsResponse?.cards || [];

  return (
    <div className="max-w-6xl mx-auto">
      <h1 className="text-3xl font-bold text-gray-900 mb-6">Credit Cards</h1>
      
      <div className="mb-6">
        <p className="text-gray-600">
          Browse through our comprehensive database of credit cards from Singapore's leading banks.
        </p>
        <p className="text-sm text-gray-500 mt-2">
          Found {cards.length} credit cards
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {cards.map((card: CreditCard) => (
          <div key={card.id} className="bg-white rounded-lg shadow-md hover:shadow-lg transition-shadow duration-300">
            <div className="p-6">
              {/* Card Header */}
              <div className="flex items-start justify-between mb-4">
                <div>
                  <h3 className="text-lg font-semibold text-gray-900 mb-1">
                    {card.name}
                  </h3>
                  <p className="text-sm text-gray-600">{card.bank}</p>
                </div>
                <div className="flex flex-col items-end">
                  <span className={`px-2 py-1 text-xs rounded-full font-medium ${
                    card.card_type === 'visa' ? 'bg-blue-100 text-blue-800' :
                    card.card_type === 'mastercard' ? 'bg-orange-100 text-orange-800' :
                    'bg-gray-100 text-gray-800'
                  }`}>
                    {card.card_type.toUpperCase()}
                  </span>
                </div>
              </div>

              {/* Card Description */}
              <p className="text-sm text-gray-600 mb-4 line-clamp-2">
                {card.description}
              </p>

              {/* Card Details */}
              <div className="space-y-2 mb-4">
                <div className="flex justify-between text-sm">
                  <span className="text-gray-500">Annual Fee:</span>
                  <span className="font-medium">
                    {card.annual_fee === 0 ? 'Free' : `$${card.annual_fee}`}
                  </span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-gray-500">Min Income:</span>
                  <span className="font-medium">${card.min_income.toLocaleString()}</span>
                </div>
                {card.welcome_bonus && (
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-500">Welcome Bonus:</span>
                    <span className="font-medium text-green-600">{card.welcome_bonus}</span>
                  </div>
                )}
              </div>

              {/* Action Button */}
              <button className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 transition-colors duration-200 text-sm font-medium">
                View Details
              </button>
            </div>
          </div>
        ))}
      </div>

      {cards.length === 0 && (
        <div className="text-center py-12">
          <p className="text-gray-500">No credit cards found.</p>
        </div>
      )}
    </div>
  );
};

export default Cards; 