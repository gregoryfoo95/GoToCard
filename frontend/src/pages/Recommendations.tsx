import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from 'react-query';
import { recommendationAPI, userAPI } from '../services/api';
import { Recommendation } from '../types';
import toast from 'react-hot-toast';

const Recommendations: React.FC = () => {
  const queryClient = useQueryClient();
  const [selectedUserId, setSelectedUserId] = useState<number>(1); // Default user for demo

  // Fetch users for selection
  const { data: usersResponse, isLoading: usersLoading } = useQuery(
    'users',
    userAPI.getAll
  );

  // Fetch existing recommendations
  const { data: recommendationsResponse, isLoading: recommendationsLoading, refetch } = useQuery(
    ['recommendations', selectedUserId],
    () => recommendationAPI.getByUser(selectedUserId),
    {
      enabled: selectedUserId > 0,
    }
  );

  // Generate new recommendations
  const generateRecommendationsMutation = useMutation(
    (userId: number) => recommendationAPI.generate(userId),
    {
      onSuccess: () => {
        toast.success('New recommendations generated!');
        queryClient.invalidateQueries(['recommendations', selectedUserId]);
        refetch();
      },
      onError: (error: any) => {
        toast.error(error.response?.data?.message || 'Failed to generate recommendations');
      },
    }
  );

  const users = usersResponse?.users || [];
  const recommendations = recommendationsResponse?.recommendations || [];

  const handleGenerateRecommendations = () => {
    generateRecommendationsMutation.mutate(selectedUserId);
  };

  const getRewardTypeDisplay = (rec: Recommendation) => {
    // Find the best benefit for display
    const card = rec.card;
    if (card.card_benefits && card.card_benefits.length > 0) {
      const benefit = card.card_benefits.find(b => b.category_id === rec.category.id);
      if (benefit) {
        if (benefit.cashback_rate > 0) {
          return {
            type: 'Cashback',
            rate: `${benefit.cashback_rate}%`,
            color: 'bg-green-100 text-green-800'
          };
        }
        if (benefit.miles_rate > 0) {
          return {
            type: 'Miles',
            rate: `${benefit.miles_rate}x`,
            color: 'bg-blue-100 text-blue-800'
          };
        }
        if (benefit.points_rate > 0) {
          return {
            type: 'Points',
            rate: `${benefit.points_rate}x`,
            color: 'bg-purple-100 text-purple-800'
          };
        }
      }
    }
    return {
      type: 'Reward',
      rate: 'Available',
      color: 'bg-gray-100 text-gray-800'
    };
  };

  const getScoreColor = (score: number) => {
    if (score >= 80) return 'text-green-600';
    if (score >= 60) return 'text-yellow-600';
    return 'text-red-600';
  };

  const getScoreBg = (score: number) => {
    if (score >= 80) return 'bg-green-100';
    if (score >= 60) return 'bg-yellow-100';
    return 'bg-red-100';
  };

  if (usersLoading) {
    return (
      <div className="max-w-6xl mx-auto">
        <h1 className="text-3xl font-bold text-gray-900 mb-6">Credit Card Recommendations</h1>
        <div className="flex justify-center py-12">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-6xl mx-auto">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-gray-900">Credit Card Recommendations</h1>
        <button
          onClick={handleGenerateRecommendations}
          disabled={generateRecommendationsMutation.isLoading}
          className="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {generateRecommendationsMutation.isLoading ? 'Generating...' : 'Generate New Recommendations'}
        </button>
      </div>

      {/* User Selection */}
      {users.length > 0 && (
        <div className="bg-white rounded-lg shadow p-6 mb-6">
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Select User (Demo)
          </label>
          <select
            value={selectedUserId}
            onChange={(e) => setSelectedUserId(parseInt(e.target.value))}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            {users.map((user) => (
              <option key={user.id} value={user.id}>
                {user.name} ({user.email})
              </option>
            ))}
          </select>
        </div>
      )}

      {/* Recommendations Info */}
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-6">
        <div className="flex">
          <svg className="w-6 h-6 text-blue-600 mr-3 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          </svg>
          <div>
            <h3 className="text-lg font-medium text-blue-900 mb-2">How Recommendations Work</h3>
            <ul className="text-sm text-blue-800 space-y-1">
              <li>• Based on your spending patterns across different categories</li>
              <li>• Considers cashback rates, miles earning, and annual fees</li>
              <li>• Factors in minimum spending requirements and caps</li>
              <li>• Calculates net annual benefit for each card-category combination</li>
              <li>• Higher scores indicate better value for your spending pattern</li>
            </ul>
          </div>
        </div>
      </div>

      {/* Recommendations List */}
      {recommendationsLoading ? (
        <div className="flex justify-center py-12">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        </div>
      ) : recommendations.length > 0 ? (
        <div className="space-y-6">
          {recommendations.map((recommendation, index) => {
            const rewardDisplay = getRewardTypeDisplay(recommendation);
            return (
              <div key={`${recommendation.card.id}-${recommendation.category.id}`} className="bg-white rounded-lg shadow-lg border border-gray-200 overflow-hidden">
                <div className="p-6">
                  {/* Header with ranking */}
                  <div className="flex items-start justify-between mb-4">
                    <div className="flex items-center">
                      <div className="w-8 h-8 bg-blue-600 text-white rounded-full flex items-center justify-center font-bold text-sm mr-4">
                        #{index + 1}
                      </div>
                      <div>
                        <h3 className="text-xl font-bold text-gray-900">{recommendation.card.name}</h3>
                        <p className="text-sm text-gray-600">{recommendation.card.bank}</p>
                      </div>
                    </div>
                    <div className={`px-3 py-1 rounded-full ${getScoreBg(recommendation.score)}`}>
                      <span className={`font-bold ${getScoreColor(recommendation.score)}`}>
                        Score: {recommendation.score}/100
                      </span>
                    </div>
                  </div>

                  {/* Category and Reward Info */}
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
                    <div className="flex items-center p-3 bg-gray-50 rounded-lg">
                      <span className="text-2xl mr-3">{recommendation.category.icon}</span>
                      <div>
                        <p className="font-medium text-gray-900">{recommendation.category.name}</p>
                        <p className="text-sm text-gray-600">Category</p>
                      </div>
                    </div>

                    <div className="flex items-center p-3 bg-gray-50 rounded-lg">
                      <div className={`p-2 rounded-lg mr-3 ${rewardDisplay.color}`}>
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"></path>
                        </svg>
                      </div>
                      <div>
                        <p className="font-medium text-gray-900">{rewardDisplay.rate}</p>
                        <p className="text-sm text-gray-600">{rewardDisplay.type}</p>
                      </div>
                    </div>

                    <div className="flex items-center p-3 bg-gray-50 rounded-lg">
                      <div className="p-2 bg-green-100 rounded-lg mr-3">
                        <svg className="w-5 h-5 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"></path>
                        </svg>
                      </div>
                      <div>
                        <p className="font-medium text-gray-900">${recommendation.estimated_reward.toFixed(2)}</p>
                        <p className="text-sm text-gray-600">Monthly Reward</p>
                      </div>
                    </div>
                  </div>

                  {/* Card Details */}
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4 text-sm">
                    <div>
                      <span className="text-gray-500">Annual Fee:</span>
                      <span className="ml-2 font-medium">
                        {recommendation.card.annual_fee === 0 ? 'Free' : `$${recommendation.card.annual_fee}`}
                      </span>
                    </div>
                    <div>
                      <span className="text-gray-500">Min Income:</span>
                      <span className="ml-2 font-medium">${recommendation.card.min_income.toLocaleString()}</span>
                    </div>
                    <div>
                      <span className="text-gray-500">Card Type:</span>
                      <span className={`ml-2 px-2 py-1 rounded text-xs font-medium ${
                        recommendation.card.card_type === 'visa' ? 'bg-blue-100 text-blue-800' :
                        recommendation.card.card_type === 'mastercard' ? 'bg-orange-100 text-orange-800' :
                        'bg-gray-100 text-gray-800'
                      }`}>
                        {recommendation.card.card_type.toUpperCase()}
                      </span>
                    </div>
                  </div>

                  {/* Recommendation Reason */}
                  <div className="bg-blue-50 rounded-lg p-4">
                    <h4 className="font-medium text-blue-900 mb-2">Why This Card?</h4>
                    <p className="text-blue-800 text-sm">{recommendation.reason}</p>
                  </div>

                  {/* Welcome Bonus */}
                  {recommendation.card.welcome_bonus && (
                    <div className="mt-4 p-3 bg-green-50 border border-green-200 rounded-lg">
                      <div className="flex items-center">
                        <svg className="w-5 h-5 text-green-600 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 8v13m0-13V6a2 2 0 112 2h-2zm0 0V5.5A2.5 2.5 0 109.5 8H12zm-7 4h14M5 12a2 2 0 110-4h14a2 2 0 110 4M5 12v7a2 2 0 002 2h10a2 2 0 002-2v-7"></path>
                        </svg>
                        <span className="font-medium text-green-800">Welcome Bonus: {recommendation.card.welcome_bonus}</span>
                      </div>
                    </div>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      ) : (
        <div className="text-center py-12">
          <svg className="w-16 h-16 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"></path>
          </svg>
          <h3 className="text-lg font-medium text-gray-900 mb-2">No Recommendations Yet</h3>
          <p className="text-gray-600 mb-4">
            Add some spending data first, then generate recommendations to see personalized credit card suggestions.
          </p>
          <button
            onClick={handleGenerateRecommendations}
            disabled={generateRecommendationsMutation.isLoading}
            className="bg-blue-600 text-white px-6 py-2 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
          >
            {generateRecommendationsMutation.isLoading ? 'Generating...' : 'Generate Recommendations'}
          </button>
        </div>
      )}
    </div>
  );
};

export default Recommendations; 