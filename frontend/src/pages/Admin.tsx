import React, { useState } from 'react';
import { useMutation, useQueryClient } from 'react-query';
import toast from 'react-hot-toast';

const Admin: React.FC = () => {
  const queryClient = useQueryClient();
  const [isScrapingAll, setIsScrapingAll] = useState(false);
  const [scrapingStatus, setScrapingStatus] = useState<string>('');

  // Scraping mutation
  const scrapingMutation = useMutation(
    async (source?: string) => {
      const url = source 
        ? `http://localhost:8080/api/v1/admin/scrape?source=${source}`
        : 'http://localhost:8080/api/v1/admin/scrape';
      
      const response = await fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
      });
      
      if (!response.ok) {
        throw new Error('Scraping failed');
      }
      
      return response.json();
    },
    {
      onSuccess: () => {
        toast.success('Scraping completed successfully!');
        setScrapingStatus('✅ Completed successfully');
        queryClient.invalidateQueries('cards');
        setIsScrapingAll(false);
      },
      onError: (error: any) => {
        toast.error(`Scraping failed: ${error.message}`);
        setScrapingStatus('❌ Failed');
        setIsScrapingAll(false);
      },
    }
  );

  const handleScrapeAll = () => {
    setIsScrapingAll(true);
    setScrapingStatus('🔄 Scraping in progress...');
    scrapingMutation.mutate(undefined);
  };

  const handleScrapeSource = (source: string) => {
    setScrapingStatus(`🔄 Scraping ${source}...`);
    scrapingMutation.mutate(source);
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-4xl mx-auto">
        <div className="bg-white shadow-lg rounded-lg overflow-hidden">
          <div className="bg-gradient-to-r from-blue-600 to-purple-600 px-6 py-4">
            <h1 className="text-2xl font-bold text-white">Admin Dashboard</h1>
            <p className="text-blue-100 mt-2">Manage credit card data scraping and system operations</p>
          </div>

          <div className="p-6">
            {/* Scraping Status */}
            <div className="mb-8">
              <h2 className="text-xl font-semibold text-gray-800 mb-4">Scraping Status</h2>
              <div className="bg-gray-50 rounded-lg p-4">
                <div className="flex items-center justify-between">
                  <span className="text-gray-600">Current Status:</span>
                  <span className="font-medium">
                    {scrapingStatus || '⏸️ Ready'}
                  </span>
                </div>
              </div>
            </div>

            {/* Scraping Controls */}
            <div className="mb-8">
              <h2 className="text-xl font-semibold text-gray-800 mb-4">Data Scraping</h2>
              <div className="grid md:grid-cols-2 gap-6">
                
                {/* Full Scraping */}
                <div className="bg-gradient-to-br from-green-50 to-emerald-50 border border-green-200 rounded-lg p-6">
                  <h3 className="text-lg font-semibold text-green-800 mb-2">Full Scraping</h3>
                  <p className="text-green-600 text-sm mb-4">
                    Scrape credit card data from all sources (SingSaver and MoneySmart)
                  </p>
                  <button
                    onClick={handleScrapeAll}
                    disabled={isScrapingAll || scrapingMutation.isLoading}
                    className="w-full bg-green-600 hover:bg-green-700 disabled:bg-gray-400 disabled:cursor-not-allowed text-white font-medium py-2 px-4 rounded-lg transition-colors duration-200"
                  >
                    {isScrapingAll ? '🔄 Scraping...' : '🚀 Start Full Scraping'}
                  </button>
                </div>

                {/* Individual Sources */}
                <div className="bg-gradient-to-br from-blue-50 to-cyan-50 border border-blue-200 rounded-lg p-6">
                  <h3 className="text-lg font-semibold text-blue-800 mb-2">Individual Sources</h3>
                  <p className="text-blue-600 text-sm mb-4">
                    Scrape from specific data sources individually
                  </p>
                  <div className="space-y-2">
                    <button
                      onClick={() => handleScrapeSource('singsaver')}
                      disabled={scrapingMutation.isLoading}
                      className="w-full bg-blue-500 hover:bg-blue-600 disabled:bg-gray-400 disabled:cursor-not-allowed text-white font-medium py-2 px-3 rounded text-sm transition-colors duration-200"
                    >
                      SingSaver
                    </button>
                    <button
                      onClick={() => handleScrapeSource('moneysmart')}
                      disabled={scrapingMutation.isLoading}
                      className="w-full bg-purple-500 hover:bg-purple-600 disabled:bg-gray-400 disabled:cursor-not-allowed text-white font-medium py-2 px-3 rounded text-sm transition-colors duration-200"
                    >
                      MoneySmart
                    </button>
                  </div>
                </div>
              </div>
            </div>

            {/* Scraping Information */}
            <div className="mb-8">
              <h2 className="text-xl font-semibold text-gray-800 mb-4">Scraping Information</h2>
              <div className="bg-gray-50 rounded-lg p-6">
                <div className="grid md:grid-cols-3 gap-6">
                  <div className="text-center">
                    <div className="text-2xl font-bold text-blue-600">2</div>
                    <div className="text-gray-600 text-sm">Data Sources</div>
                  </div>
                  <div className="text-center">
                    <div className="text-2xl font-bold text-green-600">24/7</div>
                    <div className="text-gray-600 text-sm">Available</div>
                  </div>
                  <div className="text-center">
                    <div className="text-2xl font-bold text-purple-600">Smart</div>
                    <div className="text-gray-600 text-sm">Parsing</div>
                  </div>
                </div>
              </div>
            </div>

            {/* Data Sources Info */}
            <div className="mb-8">
              <h2 className="text-xl font-semibold text-gray-800 mb-4">Data Sources</h2>
              <div className="space-y-4">
                <div className="flex items-center justify-between p-4 bg-white border border-gray-200 rounded-lg">
                  <div className="flex items-center">
                    <div className="w-3 h-3 bg-green-500 rounded-full mr-3"></div>
                    <div>
                      <div className="font-medium text-gray-800">SingSaver</div>
                      <div className="text-sm text-gray-600">www.singsaver.com.sg</div>
                    </div>
                  </div>
                  <span className="text-sm text-gray-500">Credit Cards, Cashback, Miles</span>
                </div>

                <div className="flex items-center justify-between p-4 bg-white border border-gray-200 rounded-lg">
                  <div className="flex items-center">
                    <div className="w-3 h-3 bg-green-500 rounded-full mr-3"></div>
                    <div>
                      <div className="font-medium text-gray-800">MoneySmart</div>
                      <div className="text-sm text-gray-600">www.moneysmart.sg</div>
                    </div>
                  </div>
                  <span className="text-sm text-gray-500">Credit Cards, Reviews</span>
                </div>
              </div>
            </div>

            {/* Features */}
            <div>
              <h2 className="text-xl font-semibold text-gray-800 mb-4">Scraping Features</h2>
              <div className="grid md:grid-cols-2 gap-4">
                <div className="flex items-center p-3 bg-green-50 rounded-lg">
                  <div className="w-8 h-8 bg-green-100 rounded-full flex items-center justify-center mr-3">
                    <svg className="w-4 h-4 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7"></path>
                    </svg>
                  </div>
                  <span className="text-green-800 font-medium">Rate Limiting</span>
                </div>

                <div className="flex items-center p-3 bg-blue-50 rounded-lg">
                  <div className="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center mr-3">
                    <svg className="w-4 h-4 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                    </svg>
                  </div>
                  <span className="text-blue-800 font-medium">Smart Parsing</span>
                </div>

                <div className="flex items-center p-3 bg-purple-50 rounded-lg">
                  <div className="w-8 h-8 bg-purple-100 rounded-full flex items-center justify-center mr-3">
                    <svg className="w-4 h-4 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4"></path>
                    </svg>
                  </div>
                  <span className="text-purple-800 font-medium">Duplicate Detection</span>
                </div>

                <div className="flex items-center p-3 bg-orange-50 rounded-lg">
                  <div className="w-8 h-8 bg-orange-100 rounded-full flex items-center justify-center mr-3">
                    <svg className="w-4 h-4 text-orange-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
                    </svg>
                  </div>
                  <span className="text-orange-800 font-medium">Real-time Updates</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Admin; 