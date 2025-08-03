import React, { useState } from 'react';
import { campaignService } from '../services/api';
import { Campaign, UpdateCampaignRequest } from '../types';

interface UpdateCampaignFormProps {
  campaign: Campaign;
  onSuccess: (message: string) => void;
  onError: (message: string) => void;
}

const UpdateCampaignForm: React.FC<UpdateCampaignFormProps> = ({ campaign, onSuccess, onError }) => {
  const [loading, setLoading] = useState(false);
  const [campaignForm, setCampaignForm] = useState<UpdateCampaignRequest>({
    id: campaign.id,
    name: campaign.name,
    description: campaign.description,
    start_time: campaign.start_time.slice(0, 16),
    end_time: campaign.end_time.slice(0, 16),
    total_reward: campaign.policy.total_reward,
    min_order_amount: campaign.policy.min_order_amount,
    max_tracked_orders: campaign.policy.max_tracked_orders,
  });

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value, type } = e.target;
    setCampaignForm(prev => ({
      ...prev,
      [name]: type === 'number' ? Number(value) : value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      const formattedCampaign: UpdateCampaignRequest = {
        ...campaignForm,
        start_time: new Date(campaignForm.start_time).toISOString(),
        end_time: new Date(campaignForm.end_time).toISOString(),
      };

      await campaignService.updateCampaign(formattedCampaign);
      onSuccess('Campaign updated successfully!');
    } catch (err) {
      onError('Failed to update campaign');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <h2 className="text-lg font-semibold mb-4">Update iPhone Campaign</h2>
      
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="grid grid-cols-2 gap-4">
          <input
            type="text"
            name="name"
            placeholder="Campaign Name"
            value={campaignForm.name}
            onChange={handleInputChange}
            required
            className="border rounded px-3 py-2"
          />
          <input
            type="number"
            name="total_reward"
            placeholder="Total Rewards"
            value={campaignForm.total_reward}
            onChange={handleInputChange}
            required
            className="border rounded px-3 py-2"
          />
          <input
            type="datetime-local"
            name="start_time"
            value={campaignForm.start_time}
            onChange={handleInputChange}
            required
            className="border rounded px-3 py-2"
          />
          <input
            type="datetime-local"
            name="end_time"
            value={campaignForm.end_time}
            onChange={handleInputChange}
            required
            className="border rounded px-3 py-2"
          />
          <input
            type="number"
            name="min_order_amount"
            placeholder="Min Order Amount"
            value={campaignForm.min_order_amount}
            onChange={handleInputChange}
            required
            className="border rounded px-3 py-2"
          />
          <input
            type="number"
            name="max_tracked_orders"
            placeholder="Max Tracked Orders"
            value={campaignForm.max_tracked_orders}
            onChange={handleInputChange}
            required
            className="border rounded px-3 py-2"
          />
        </div>
        <textarea
          name="description"
          placeholder="Description"
          value={campaignForm.description}
          onChange={handleInputChange}
          required
          className="w-full border rounded px-3 py-2"
          rows={3}
        />
        <button
          type="submit"
          disabled={loading}
          className="bg-blue-500 text-white px-4 py-2 rounded disabled:opacity-50"
        >
          {loading ? 'Updating...' : 'Update Campaign'}
        </button>
      </form>
    </div>
  );
};

export default UpdateCampaignForm;