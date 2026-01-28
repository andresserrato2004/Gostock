export interface Stock {
  id: number;
  ticker: string;
  company: string;
  brokerage: string;
  action: string;
  rating_from: string;
  rating_to: string;
  target_to?: string;   
  target_from?: string; 
  target_from_num?: number;
  target_to_num?: number;
  created_at: string;
}

export interface Meta {
  limit: number;
  page: number;
  total: number;
}

export interface StockResponse {
  data: Stock[];
  meta: Meta;
}

export interface RecommendationResponse {
    data: Stock[];
}

export interface FetchStocksParams {
  page?: number;
  limit?: number;
  search?: string;
  sort_by?: string;
  order?: 'asc' | 'desc';
}
