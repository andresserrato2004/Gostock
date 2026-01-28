import { ref, computed } from 'vue';
import { useStockStore } from '../stores/stockStore';
import { storeToRefs } from 'pinia';


/**
 * Composable para gestionar la lógica de Recomendaciones de Mercado.
 * * Actúa como una capa intermedia entre el `StockStore` y la vista de StockList, proporcionando:
 * 1. Estado reactivo sincronizado con Pinia.
 * 2. Formateadores de datos específicos para la UI (Moneda, Porcentajes).
 * 3. Lógica de presentación condicional (Clases CSS según Score).
 * 4. Lógica de paginación y navegación.
 * 5. Lógica de búsqueda y ordenamiento.
 * 
 */
export function useStockList() {

    const store = useStockStore();
    const { stocks, meta, isStocksLoading: loading, currentParams } = storeToRefs(store);
    const searchQuery = ref('');
    const totalPages = computed(() => Math.ceil(meta.value.total / meta.value.limit));

    const visiblePages = computed(() => {
        const total = totalPages.value;
        const current = meta.value.page;
        const delta = 2;
        const range: number[] = [];
        const rangeWithDots: (number | string)[] = [];
        let l: number;
       
        for (let i = 1; i <= total; i++) {
            if (i === 1 || i === total || (i >= current - delta && i <= current + delta)) {
                range.push(i);
            }
        }
       
        for (const i of range) {
            if (l!) {
                if (i - l === 2) {
                    rangeWithDots.push(l + 1);
                } else if (i - l !== 1) {
                    rangeWithDots.push('...');
                }
            }
            rangeWithDots.push(i);
            l = i;
        }
        return rangeWithDots;
    });


    const goToPage = (page: number | string) => {
        if(typeof page === 'number' && page >=1 && page <= totalPages.value) {
            store.setPage(page);
        }
    };

    const nextPage = () => {
        if (meta.value.page * meta.value.limit < meta.value.total) {
            store.setPage(meta.value.page + 1);
        }
    };
    
    const prevPage = () => {
        if (meta.value.page > 1) {
            store.setPage(meta.value.page - 1);
        }
    };

    const sort = (sortBy: string) => {
        if (currentParams.value.sort_by === sortBy) {
            if(currentParams.value.order === 'asc') {
                store.setSort(sortBy, 'desc');
            } else {
                store.setSort('', 'asc');
            }
        } else {
            store.setSort(sortBy, 'asc');
        }
    };

    function debounceFn(func: Function, wait: number) {
        let timeout: any;
        return function(...args: any[]) {
            clearTimeout(timeout);
            timeout = setTimeout(() => func(...args), wait);
        };
    }

    const handleSearch = debounceFn(() => {
        store.setSearch(searchQuery.value);
    }, 300);

    const init = () => {
        store.fetchStocks({});
    };

    return {
        stocks,
        meta,
        loading,
        currentParams,
        searchQuery,
        visiblePages,
        totalPages,
        
        handleSearch,
        sort,
        goToPage,
        nextPage,
        prevPage,
        init
    };
}