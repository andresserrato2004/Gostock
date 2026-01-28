export interface Recommendation {
    ticker: string;
    company: string;
    current_price: number;
    target_price: number;
    upside: number;
    rating_score: number;
    conviction: number;
    momentum: number;
    final_score: number;
    reason: string;
}

export interface Stock {
    id: number;
    ticker: string;
    company: string;
    brokerage: string;
    action: string;
    rating_from: string;
    rating_to: string;
    target_from: string;
    target_to: string;
    target_from_num: number;
    target_to_num: number;
    created_at: string;
}

export interface Meta {
    limit: number;
    page: number;
    total: number;
}

export interface FetchStocksParams {
    page?: number;
    limit?: number;
    search?: string;
    sort_by?: string;
    order?: string;
}

export interface StockResponse {
    data: Stock[];
    meta: Meta;
}

export interface RecommendationResponse {
    data: Recommendation[];
}
