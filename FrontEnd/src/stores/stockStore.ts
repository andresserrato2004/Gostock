import { defineStore } from 'pinia';
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

interface State {
    stocks: Stock[];
    meta: Meta;
    recommendations: Recommendation[]; 
    loading: boolean;
    error: string | null;
    currentParams: FetchStocksParams;
}

export const useStockStore = defineStore('stock', {
    state: (): State => ({
        stocks: [],
        meta: { limit: 10, page: 1, total: 0 },
        recommendations: [], // Inicialmente vacío
        loading: false,
        error: null,
        currentParams: { 
            page: 1, 
            limit: 10, 
            sort_by: 'target_to_num', 
            order: 'desc' 
        },
    }),
    
    actions: {
        async fetchStocks(params: FetchStocksParams) {
            this.loading = true;
            this.error = null;
            this.currentParams = { ...this.currentParams, ...params };
            try {
                const response = await api.getStocks(this.currentParams);
                console.log('Respuesta del Back end (Stocks):', response.data);
                this.stocks = response.data.data;
                this.meta = response.data.meta;
            } catch (err: any) {
                this.error = err.message || 'Failed to fetch stocks';
            } finally {
                this.loading = false;
            }
        },

        // --- NUEVA ACCIÓN PARA TRAER RECOMENDACIONES ---
        async fetchRecommendations(limit: number = 3) {
            // Nota: Podrías querer un loading separado para recomendaciones si quieres que carguen independiente de la tabla
            this.loading = true; 
            try {
                const response = await api.getRecommendations({ limit });
                console.log('Respuesta del Back end (Recommendations):', response.data);
                this.recommendations = response.data.data;
            } catch (err: any) {
                console.error("Failed to fetch recommendations", err);
                this.error = err.message || 'Failed to fetch recommendations';
            } finally {
                this.loading = false;
            }
        },

        setPage(page: number) {
            this.fetchStocks({ page });
        },

        setSearch(search: string) {
            this.fetchStocks({ search, page: 1 });
        },

        setSort(sortBy: string, order: string) {
            this.fetchStocks({ sort_by: sortBy, order, page: 1 });
        }
    },
});
