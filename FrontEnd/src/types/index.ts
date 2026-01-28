/**
 * Algorithm stock recommendation.
 */
export interface Recommendation {
    
    ticker: string; /** Stock ticker (e.g., "AAPL") */
    company: string; /** Company name */
    current_price: number; /** Current market price */
    target_price: number; /** Analyst average target price */
    upside: number; /** Potential gain % */
    rating_score: number; /** Normalized rating (0.0 - 1.0) */
    conviction: number; /** Analyst target trend score */
    momentum: number; /** Market price momentum score */
    final_score: number;    /** 0-100 proprietary ranking */
    reason: string;    /** Recommendation rationale */
}

/**
 * Raw brokerage rating event.
 */
export interface Stock {

    id: number;    /** Unique ID */
    ticker: string; /** Stock ticker */
    company: string; /** Company name */
    brokerage: string; /** Rating firm */
    action: string; /** Rating action (e.g., Upgrade) */
    rating_from: string; /** Previous rating */
    rating_to: string;   /** New rating */
    target_from: string; /** Previous target */
    target_to: string;   /** New target */
    target_from_num: number; /** Parsed target_from */
    target_to_num: number;/** Parsed target_to */
    created_at: string;/** Creation timestamp */
}

/**
 * Pagination metadata.
 */
export interface Meta {
    limit: number;     /** Items per page */
    page: number;    /** Current page */
    total: number; /** Total items */
}

/**
 * Query parameters for Stock List.
 */
export interface FetchStocksParams {
    page?: number;
    limit?: number;
    search?: string;
    sort_by?: string;
    order?: string;
}

/**
 * API Response: Stocks.
 */
export interface StockResponse {
    data: Stock[];
    meta: Meta;
}

/**
 * API Response: Recommendations.
 */
export interface RecommendationResponse {
    data: Recommendation[];
}
