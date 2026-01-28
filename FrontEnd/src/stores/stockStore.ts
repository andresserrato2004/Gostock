import { defineStore } from 'pinia';
import { ref } from 'vue';
import api from '../services/api';

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

interface Meta {
    limit: number;
    page: number;
    total: number;
}

interface FetchStocksParams {
    page?: number;
    limit?: number;
    search?: string;
    sort_by?: string;
    order?: string;
}

export const useStockStore = defineStore('stock', () => {
    
    const stocks = ref<Stock[]>([]);
    const meta = ref<Meta>({ limit: 10, page: 1, total: 0 });
    const isStocksLoading = ref(false);
    const stocksError = ref<string | null>(null);
    const currentParams = ref<FetchStocksParams>({ 
        page: 1, 
        limit: 10, 
        sort_by: 'target_to_num', 
        order: 'desc' 
    });

    const recommendations = ref<Recommendation[]>([]);
    const isRecsLoading = ref(false);
    const recsError = ref<string | null>(null);

    const fetchStocks = async (params: FetchStocksParams) => {
        isStocksLoading.value = true;
        stocksError.value = null;
        currentParams.value = { ...currentParams.value, ...params };
        
        try {
            const response = await api.getStocks(currentParams.value);
            console.log('Respuesta del Back end (Stocks):', response.data);
            stocks.value = response.data.data;
            meta.value = response.data.meta;
        } catch (err: any) {
            stocksError.value = err.message || 'Failed to fetch stocks';
        } finally {
            isStocksLoading.value = false;
        }
    };

    const fetchRecommendations = async (limit: number = 3) => {
        isRecsLoading.value = true;
        recsError.value = null;
        try {
            const response = await api.getRecommendations({ limit });
            console.log('Respuesta del Back end (Recommendations):', response.data);
            recommendations.value = response.data.data;
        } catch (err: any) {
            console.error("Failed to fetch recommendations", err);
            recsError.value = err.message || 'Failed to fetch recommendations';
        } finally {
            isRecsLoading.value = false;
        }
    };

    const setPage = (page: number) => {
        fetchStocks({ page });
    };

    const setSearch = (search: string) => {
        fetchStocks({ search, page: 1 });
    };

    const setSort = (sortBy: string, order: string) => {
        fetchStocks({ sort_by: sortBy, order, page: 1 });
    };

    return {
        stocks,
        meta,
        currentParams,
        isStocksLoading,
        stocksError,
        
        recommendations,
        isRecsLoading,
        recsError,

        fetchStocks,
        fetchRecommendations,
        setPage,
        setSearch,
        setSort
    };
});
