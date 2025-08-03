import React, { useState, useEffect } from 'react';
import { campaignService } from '../services/api';
import { IphoneWinner, Campaign } from '../types';
import CreateCampaignForm from '../components/CreateCampaignForm';
import UpdateCampaignForm from '../components/UpdateCampaignForm';

const CampaignsPage: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [initialLoading, setInitialLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [winners, setWinners] = useState<IphoneWinner[]>([]);
  const [showWinners, setShowWinners] = useState(false);
  const [existingCampaign, setExistingCampaign] = useState<Campaign | null>(null);
  const [isUpdateMode, setIsUpdateMode] = useState(false);

  useEffect(() => {
    const checkExistingCampaign = async () => {
      try {
        setInitialLoading(true);
        const response = await campaignService.getIphoneCampaign();
        setExistingCampaign(response.data);
        setIsUpdateMode(true);
      } catch (err) {
        setIsUpdateMode(false);
        setExistingCampaign(null);
      } finally {
        setInitialLoading(false);
      }
    };

    checkExistingCampaign();
  }, []);

  const handleSuccess = (message: string) => {
    setSuccess(message);
    setError(null);
  };

  const handleError = (message: string) => {
    setError(message);
    setSuccess(null);
  };

  const handleCampaignCreated = () => {
    window.location.reload();
  };

  const handleGetWinners = async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await campaignService.getIphoneWinners();
      setWinners(response.data);
      setShowWinners(true);
    } catch (err) {
      setError('Failed to fetch winners');
    } finally {
      setLoading(false);
    }
  };

  if (initialLoading) return <div>Loading...</div>;

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Campaigns</h1>
      
      {/* Global messages */}
      {error && <div className="p-3 bg-red-100 text-red-700 rounded">{error}</div>}
      {success && <div className="p-3 bg-green-100 text-green-700 rounded">{success}</div>}
      
      {/* Create/Update Campaign */}
      {isUpdateMode && existingCampaign ? (
        <UpdateCampaignForm 
          campaign={existingCampaign}
          onSuccess={handleSuccess}
          onError={handleError}
        />
      ) : (
        <CreateCampaignForm 
          onSuccess={handleSuccess}
          onError={handleError}
          onCreated={handleCampaignCreated}
        />
      )}

      {/* iPhone Inventory */}
      {existingCampaign && (
        <div className="bg-white rounded-lg shadow p-6 mb-6">
          <h2 className="text-lg font-semibold mb-4">iPhone Inventory</h2>
          <div className="grid grid-cols-3 gap-4">
            <div className="bg-blue-50 p-4 rounded-lg text-center">
              <p className="text-sm text-gray-600">Total iPhones</p>
              <p className="text-2xl font-bold text-blue-600">{existingCampaign.policy.total_reward}</p>
            </div>
            <div className="bg-orange-50 p-4 rounded-lg text-center">
              <p className="text-sm text-gray-600">Winners So Far</p>
              <p className="text-2xl font-bold text-orange-600">{winners.length}</p>
            </div>
            <div className="bg-green-50 p-4 rounded-lg text-center">
              <p className="text-sm text-gray-600">iPhones Left</p>
              <p className="text-2xl font-bold text-green-600">
                {existingCampaign.policy.total_reward - winners.length}
              </p>
            </div>
          </div>
        </div>
      )}

      {/* Get Winners */}
      <div className="bg-white rounded-lg shadow p-6">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-lg font-semibold">iPhone Campaign Winners</h2>
          <button
            onClick={handleGetWinners}
            disabled={loading}
            className="bg-blue-500 text-white px-4 py-2 rounded disabled:opacity-50"
          >
            {loading ? 'Loading...' : 'Get Winners'}
          </button>
        </div>

        {showWinners && (
          <div className="bg-white rounded-lg shadow">
            <table className="w-full">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-4 py-2 text-left">Customer ID</th>
                  <th className="px-4 py-2 text-left">Customer Name</th>
                  <th className="px-4 py-2 text-left">First Order</th>
                  <th className="px-4 py-2 text-left">Max Amount</th>
                </tr>
              </thead>
              <tbody>
                {winners.map((winner, index) => (
                  <tr key={`${winner.customer_id}-${index}`} className="border-t">
                    <td className="px-4 py-2">{winner.customer_id}</td>
                    <td className="px-4 py-2">{winner.customer_name}</td>
                    <td className="px-4 py-2">{new Date(winner.first_order_time).toLocaleDateString()}</td>
                    <td className="px-4 py-2">${winner.max_total_order_amount}</td>
                  </tr>
                ))}
              </tbody>
            </table>
            
            {winners.length === 0 && (
              <div className="p-8 text-center text-gray-500">No winners found</div>
            )}
          </div>
        )}
      </div>
    </div>
  );
};

export default CampaignsPage;