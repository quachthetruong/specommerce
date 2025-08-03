import React, { useState } from 'react';
import { campaignService } from '../services/api';
import { CreateCampaignRequest } from '../types';

interface CreateCampaignFormProps {
  onSuccess: (message: string) => void;
  onError: (message: string) => void;
  onCreated: () => void;
}

const CreateCampaignForm: React.FC<CreateCampaignFormProps> = ({ onSuccess, onError, onCreated }) => {
  const [loading, setLoading] = useState(false);
  const [campaignForm, setCampaignForm] = useState<CreateCampaignRequest>({
    name: '',
    description: '',
    start_time: '',
    end_time: '',
    total_reward: 20,
    min_order_amount: 200,
    max_tracked_orders: 50,
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
      const formattedCampaign = {
        ...campaignForm,
        start_time: new Date(campaignForm.start_time).toISOString(),
        end_time: new Date(campaignForm.end_time).toISOString(),
      };

      await campaignService.createCampaign(formattedCampaign);
      onSuccess('Campaign created successfully!');
      onCreated();
    } catch (err) {
      onError('Failed to create campaign');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <h2 className="text-lg font-semibold mb-4">Create iPhone Campaign</h2>
      
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
          {loading ? 'Creating...' : 'Create Campaign'}
        </button>
      </form>
    </div>
  );
};

export default CreateCampaignForm;