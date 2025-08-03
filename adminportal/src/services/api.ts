import axios from 'axios';
import {
  BaseResponse,
  PaginatedResponse,
  Order,
  Payment,
  Campaign,
  CreateCampaignRequest,
  IphoneWinner,
  SearchOrdersParams,
  SearchPaymentsParams,
  UpdateCampaignRequest,
} from '../types';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost';

// Service URLs
const ORDER_SERVICE = `${API_BASE_URL}:8080/api/admin/v1`;
const PAYMENT_SERVICE = `${API_BASE_URL}:8081/api/admin/v1`;
const CAMPAIGN_SERVICE = `${API_BASE_URL}:8082/api/admin/v1`;

// Create axios instances for each service
const orderApi = axios.create({ baseURL: ORDER_SERVICE });
const paymentApi = axios.create({ baseURL: PAYMENT_SERVICE });
const campaignApi = axios.create({ baseURL: CAMPAIGN_SERVICE });

// Order API
export const orderService = {
  searchOrders: async (params: SearchOrdersParams): Promise<PaginatedResponse<Order>> => {
    const response = await orderApi.get('/orders/search', { params });
    // Handle both old (Data/Metadata) and new (data/metadata) API formats
    return {
      data: response.data.data || response.data.Data || [],
      metadata: {
        total: response.data.metadata?.total || response.data.Metadata?.Total || 0,
        page_size: response.data.metadata?.page_size || response.data.Metadata?.PageSize || 0,
        page_number: response.data.metadata?.page_number || response.data.Metadata?.PageNumber || 0,
        total_pages: response.data.metadata?.total_pages || response.data.Metadata?.TotalPages || 0,
      },
    };
  },
};

// Payment API
export const paymentService = {
  searchPayments: async (params: SearchPaymentsParams): Promise<PaginatedResponse<Payment>> => {
    const response = await paymentApi.get('/payments/search', { params });
    // Handle both old (Data/Metadata) and new (data/metadata) API formats
    return {
      data: response.data.data || response.data.Data || [],
      metadata: {
        total: response.data.metadata?.total || response.data.Metadata?.Total || 0,
        page_size: response.data.metadata?.page_size || response.data.Metadata?.PageSize || 0,
        page_number: response.data.metadata?.page_number || response.data.Metadata?.PageNumber || 0,
        total_pages: response.data.metadata?.total_pages || response.data.Metadata?.TotalPages || 0,
      },
    };
  },
};

// Campaign API
export const campaignService = {
  createCampaign: async (campaign: CreateCampaignRequest): Promise<BaseResponse<Campaign>> => {
    const response = await campaignApi.post('/campaigns/iphones', campaign);
    return response.data;
  },

  getIphoneCampaign: async (): Promise<BaseResponse<Campaign>> => {
    const response = await campaignApi.get('/campaigns/iphones');
    return response.data;
  },

  updateCampaign: async (campaign: UpdateCampaignRequest): Promise<BaseResponse<Campaign>> => {
    const response = await campaignApi.put(`/campaigns/iphones/${campaign.id}`, campaign);
    return response.data;
  },
  
  getIphoneWinners: async (): Promise<BaseResponse<IphoneWinner[]>> => {
    const response = await campaignApi.get('/campaigns/iphones/winners');
    return response.data;
  },
};

// Error handling interceptors
[orderApi, paymentApi, campaignApi].forEach(api => {
  api.interceptors.response.use(
    (response) => response,
    (error) => {
      console.error('API Error:', error.response?.data || error.message);
      return Promise.reject(error);
    }
  );
});