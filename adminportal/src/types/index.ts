// Common types
export interface BaseResponse<T> {
  data: T;
}

export interface PaginationMetadata {
  total: number;
  page_size: number;
  page_number: number;
  total_pages: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  metadata: PaginationMetadata;
}

// Order types
export interface Order {
  id: string;
  customer_id: string;
  customer_name: string;
  total_amount: number;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface SearchOrdersParams {
  page?: number;
  size?: number;
  sort?: string;
}

// Payment types
export interface Payment {
  id: string;
  order_id: string;
  customer_id: string;
  total_amount: number;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface SearchPaymentsParams {
  page?: number;
  size?: number;
  sort?: string;
}

// Campaign types
export interface Campaign {
  id: number;
  name: string;
  type: string;
  description: string;
  policy: {
    total_reward: number;
    min_order_amount: number;
    max_tracked_orders: number;
  };
  start_time: string;
  end_time: string;
  created_at: string;
  updated_at: string;
}

export interface CreateCampaignRequest {
  name: string;
  description: string;
  start_time: string;
  end_time: string;
  total_reward: number;
  min_order_amount: number;
  max_tracked_orders: number;
}

export interface UpdateCampaignRequest {
  id: number;
  name: string;
  description: string;
  start_time: string;
  end_time: string;
  total_reward: number;
  min_order_amount: number;
  max_tracked_orders: number;
}

export interface IphoneWinner {
  customer_id: string;
  customer_name: string;
  first_order_time: string;
  max_total_order_amount: number;
}