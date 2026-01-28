import { ref } from 'vue';
import { useStockStore } from '../stores/stockStore';
import { storeToRefs } from 'pinia';

/**
 * Composable para gestionar la lógica de Recomendaciones de Mercado.
 * * Actúa como una capa intermedia entre el `StockStore` y la vista, proporcionando:
 * 1. Estado reactivo sincronizado con Pinia.
 * 2. Formateadores de datos específicos para la UI (Moneda, Porcentajes).
 * 3. Lógica de presentación condicional (Clases CSS según Score).
 */
export function useRecommendations() {
    const store = useStockStore();
    const { recommendations, isRecsLoading: loading, recsError: error } = storeToRefs(store);

    const currentLimit = ref(7);

    const formatCurrency = (value: number) => {
        return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(value);
    };

    const formatPercent = (value: number) => {
        return `${value.toFixed(2)}%`;
    };

    const getScoreClass = (score: number) => {
        if (score >= 80) return 'bg-green-100 text-green-800 border-green-200';
        if (score >= 50) return 'bg-blue-100 text-blue-800 border-blue-200';
        return 'bg-yellow-100 text-yellow-800 border-yellow-200';
    };

    const init = () => {
        store.fetchRecommendations(currentLimit.value); 
    };

    const handleLimitChange = () => {
        store.fetchRecommendations(currentLimit.value);
    };

    return {
        recommendations,
        loading,
        error,
        init,
        formatCurrency,
        formatPercent,
        getScoreClass,
        currentLimit,
        handleLimitChange
    };
}
