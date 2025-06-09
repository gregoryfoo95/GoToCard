import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from 'react-query';
import { categoryAPI, spendingAPI, userAPI } from '../services/api';
import { SpendingRequest } from '../types';
import toast from 'react-hot-toast';

const Spending: React.FC = () => {
  const queryClient = useQueryClient();
  const [selectedUserId, setSelectedUserId] = useState<number>(1); // Default user for demo
  const [formData, setFormData] = useState<SpendingRequest>({
    category_id: 0,
    amount: 0,
    month: new Date().getMonth() + 1,
    year: new Date().getFullYear(),
  });

  // Fetch data
  const { data: categoriesResponse, isLoading: categoriesLoading } = useQuery(
    'categories',
    categoryAPI.getAll
  );

  const { data: usersResponse, isLoading: usersLoading } = useQuery(
    'users',
    userAPI.getAll
  );

  const { data: spendingsResponse, isLoading: spendingsLoading } = useQuery(
    ['user-spending', selectedUserId],
    () => spendingAPI.getUserSpending(selectedUserId),
    {
      enabled: selectedUserId > 0,
    }
  );

  // Add spending mutation
  const addSpendingMutation = useMutation(
    (data: { userId: number; spending: SpendingRequest }) =>
      spendingAPI.add(data.userId, data.spending),
    {
      onSuccess: () => {
        toast.success('Spending added successfully!');
        queryClient.invalidateQueries(['user-spending', selectedUserId]);
        setFormData({
          category_id: 0,
          amount: 0,
          month: new Date().getMonth() + 1,
          year: new Date().getFullYear(),
        });
      },
      onError: (error: any) => {
        toast.error(error.response?.data?.message || 'Failed to add spending');
      },
    }
  );

  const categories = categoriesResponse?.categories || [];
  const users = usersResponse?.users || [];
  const spendings = spendingsResponse?.spendings || [];

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (formData.category_id === 0 || formData.amount <= 0) {
      toast.error('Please select a category and enter a valid amount');
      return;
    }
    addSpendingMutation.mutate({ userId: selectedUserId, spending: formData });
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: name === 'amount' ? parseFloat(value) || 0 : parseInt(value) || 0,
    }));
  };

  // Calculate total spending
  const totalSpending = spendings.reduce((sum, spending) => sum + spending.amount, 0);
  const currentMonthSpendings = spendings.filter(
    s => s.month === new Date().getMonth() + 1 && s.year === new Date().getFullYear()
  );
  const currentMonthTotal = currentMonthSpendings.reduce((sum, spending) => sum + spending.amount, 0);

  // Group spendings by category
  const spendingsByCategory = spendings.reduce((acc, spending) => {
    const categoryName = spending.category.name;
    if (!acc[categoryName]) {
      acc[categoryName] = { total: 0, count: 0, icon: spending.category.icon };
    }
    acc[categoryName].total += spending.amount;
    acc[categoryName].count += 1;
    return acc;
  }, {} as Record<string, { total: number; count: number; icon: string }>);

  if (categoriesLoading || usersLoading) {
    return (
      <div className="max-w-6xl mx-auto">
        <h1 className="text-3xl font-bold text-gray-900 mb-6">My Spending</h1>
        <div className="flex justify-center py-12">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-6xl mx-auto">
      <h1 className="text-3xl font-bold text-gray-900 mb-6">My Spending</h1>

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

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Add Spending Form */}
        <div className="lg:col-span-1">
          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-xl font-bold text-gray-900 mb-4">Add Spending</h2>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Category
                </label>
                <select
                  name="category_id"
                  value={formData.category_id}
                  onChange={handleInputChange}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
                >
                  <option value={0}>Select a category</option>
                  {categories.map((category) => (
                    <option key={category.id} value={category.id}>
                      {category.icon} {category.name}
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Amount ($)
                </label>
                <input
                  type="number"
                  name="amount"
                  value={formData.amount || ''}
                  onChange={handleInputChange}
                  step="0.01"
                  min="0"
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="0.00"
                  required
                />
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Month
                  </label>
                  <select
                    name="month"
                    value={formData.month}
                    onChange={handleInputChange}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  >
                    {Array.from({ length: 12 }, (_, i) => (
                      <option key={i + 1} value={i + 1}>
                        {new Date(0, i).toLocaleString('default', { month: 'long' })}
                      </option>
                    ))}
                  </select>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Year
                  </label>
                  <select
                    name="year"
                    value={formData.year}
                    onChange={handleInputChange}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  >
                    {Array.from({ length: 5 }, (_, i) => {
                      const year = new Date().getFullYear() - 2 + i;
                      return (
                        <option key={year} value={year}>
                          {year}
                        </option>
                      );
                    })}
                  </select>
                </div>
              </div>

              <button
                type="submit"
                disabled={addSpendingMutation.isLoading}
                className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {addSpendingMutation.isLoading ? 'Adding...' : 'Add Spending'}
              </button>
            </form>
          </div>
        </div>

        {/* Spending Overview */}
        <div className="lg:col-span-2">
          {/* Summary Cards */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center">
                <div className="p-2 bg-blue-100 rounded-lg">
                  <svg className="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"></path>
                  </svg>
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600">Total Spending</p>
                  <p className="text-2xl font-bold text-gray-900">${totalSpending.toFixed(2)}</p>
                </div>
              </div>
            </div>

            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center">
                <div className="p-2 bg-green-100 rounded-lg">
                  <svg className="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                  </svg>
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600">This Month</p>
                  <p className="text-2xl font-bold text-gray-900">${currentMonthTotal.toFixed(2)}</p>
                </div>
              </div>
            </div>
          </div>

          {/* Spending by Category */}
          <div className="bg-white rounded-lg shadow p-6 mb-6">
            <h3 className="text-lg font-bold text-gray-900 mb-4">Spending by Category</h3>
            <div className="space-y-3">
              {Object.entries(spendingsByCategory).map(([categoryName, data]) => (
                <div key={categoryName} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                  <div className="flex items-center">
                    <span className="text-2xl mr-3">{data.icon}</span>
                    <div>
                      <p className="font-medium text-gray-900">{categoryName}</p>
                      <p className="text-sm text-gray-600">{data.count} transactions</p>
                    </div>
                  </div>
                  <p className="font-bold text-gray-900">${data.total.toFixed(2)}</p>
                </div>
              ))}
            </div>
            {Object.keys(spendingsByCategory).length === 0 && (
              <p className="text-gray-500 text-center py-4">No spending data yet</p>
            )}
          </div>

          {/* Recent Transactions */}
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-bold text-gray-900 mb-4">Recent Transactions</h3>
            {spendingsLoading ? (
              <div className="flex justify-center py-4">
                <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
              </div>
            ) : (
              <div className="space-y-3">
                {spendings.slice(0, 10).map((spending) => (
                  <div key={spending.id} className="flex items-center justify-between p-3 border border-gray-200 rounded-lg">
                    <div className="flex items-center">
                      <span className="text-xl mr-3">{spending.category.icon}</span>
                      <div>
                        <p className="font-medium text-gray-900">{spending.category.name}</p>
                        <p className="text-sm text-gray-600">
                          {new Date(0, spending.month - 1).toLocaleString('default', { month: 'long' })} {spending.year}
                        </p>
                      </div>
                    </div>
                    <p className="font-bold text-gray-900">${spending.amount.toFixed(2)}</p>
                  </div>
                ))}
                {spendings.length === 0 && (
                  <p className="text-gray-500 text-center py-4">No transactions yet</p>
                )}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default Spending; 