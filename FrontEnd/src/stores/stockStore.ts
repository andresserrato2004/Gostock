import type { Stock, Meta, Recommendation, FetchStocksParams } from '../types';
import { defineStore } from 'pinia';
import { ref } from 'vue';
import api from '../services/api';


/**
 * @store StockStore
 * @description Gestión centralizada de activos financieros y recomendaciones algorítmicas.
 * Maneja la paginación, búsqueda y ordenamiento de stocks, así como la obtención
 * de oportunidades de mercado (recommendations).
 */

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
    const recsLimit = ref(7);

    const clearErrors = () => {
        stocksError.value = null;
        recsError.value = null;
    };

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
            console.error("Failed to fetch stocks", err);
            if (err.code === 'ERR_NETWORK' || err.message === 'Network Error') {
                 stocksError.value = 'We are having trouble connecting to our servers. Please check your internet connection regarding the backend or try again later.';
            } else if (err.response && err.response.status >= 500) {
                 stocksError.value = 'Our servers are currently experiencing issues. Please try again later.';
            } else if (err.response && err.response.status === 429) {
                 stocksError.value = 'You are making too many requests. Please wait a moment before trying again.';
            } else {
                 stocksError.value = err.response?.data?.message || err.message || 'Failed to fetch stocks';
            }
        } finally {
            isStocksLoading.value = false;
        }
    };

    const fetchRecommendations = async (limit?: number) => {
        isRecsLoading.value = true;
        recsError.value = null;
        if (limit) recsLimit.value = limit;

        try {
            const response = await api.getRecommendations({ limit: recsLimit.value });
            console.log('Respuesta del Back end (Recommendations):', response.data);
            recommendations.value = response.data.data;
        } catch (err: any) {
            console.error("Failed to fetch recommendations", err);
            if (err.code === 'ERR_NETWORK' || err.message === 'Network Error') {
                 recsError.value = 'We are having trouble connecting to our servers. Please check your internet connection regarding the backend or try again later.';
            } else if (err.response && err.response.status >= 500) {
                 recsError.value = 'Our servers are currently experiencing issues. Please try again later.';
            } else if (err.response && err.response.status === 429) {
                 recsError.value = 'You are making too many requests. Please wait a moment before trying again.';
            } else {
                recsError.value = err.response?.data?.message || err.message || 'Failed to fetch recommendations';
            }
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
        recsLimit,
        
        clearErrors,
        fetchStocks,
        fetchRecommendations,
        setPage,
        setSearch,
        setSort
    };
});
