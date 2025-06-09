import React from 'react';
import { Link } from 'react-router-dom';

const Home: React.FC = () => {
  return (
    <div className="max-w-4xl mx-auto">
      <div className="text-center py-16">
        <h1 className="text-5xl font-bold text-gray-900 mb-6">
          Find Your Perfect Credit Card
        </h1>
        <p className="text-xl text-gray-600 mb-8 max-w-2xl mx-auto">
          Get personalized credit card recommendations based on your spending patterns. 
          Maximize your rewards and benefits with data-driven insights from Singapore's 
          leading financial platforms.
        </p>
        
        <div className="flex justify-center space-x-4">
          <Link
            to="/dashboard"
            className="bg-primary-600 text-white px-8 py-3 rounded-lg font-semibold hover:bg-primary-700 transition-colors"
          >
            Get Started
          </Link>
          <Link
            to="/cards"
            className="border-2 border-primary-600 text-primary-600 px-8 py-3 rounded-lg font-semibold hover:bg-primary-50 transition-colors"
          >
            Browse Cards
          </Link>
        </div>
      </div>

      <div className="grid md:grid-cols-3 gap-8 py-16">
        <div className="text-center">
          <div className="bg-primary-100 rounded-full w-16 h-16 flex items-center justify-center mx-auto mb-4">
            <span className="text-2xl">📊</span>
          </div>
          <h3 className="text-xl font-semibold mb-2">Smart Analysis</h3>
          <p className="text-gray-600">
            Our algorithm analyzes your spending patterns to recommend the most suitable credit cards for maximum rewards.
          </p>
        </div>

        <div className="text-center">
          <div className="bg-primary-100 rounded-full w-16 h-16 flex items-center justify-center mx-auto mb-4">
            <span className="text-2xl">💳</span>
          </div>
          <h3 className="text-xl font-semibold mb-2">Comprehensive Database</h3>
          <p className="text-gray-600">
            Access up-to-date information on credit cards from major Singapore banks with detailed benefits and features.
          </p>
        </div>

        <div className="text-center">
          <div className="bg-primary-100 rounded-full w-16 h-16 flex items-center justify-center mx-auto mb-4">
            <span className="text-2xl">🎯</span>
          </div>
          <h3 className="text-xl font-semibold mb-2">Personalized Recommendations</h3>
          <p className="text-gray-600">
            Get tailored suggestions that match your lifestyle and spending categories for optimal financial benefits.
          </p>
        </div>
      </div>
    </div>
  );
};

export default Home; 