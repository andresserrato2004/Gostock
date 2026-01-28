<template>
  <div class="bg-white rounded-lg shadow-lg overflow-hidden flex flex-col h-full border border-gray-100">
    
    <div class="px-4 py-4 sm:px-8 border-b flex-shrink-0 flex flex-col md:flex-row justify-between items-center gap-4 bg-gray-50">
        <h2 class="text-xl font-semibold leading-tight text-gray-800">
            Real-Time Stock Feed
        </h2>

        <div class="relative w-full md:w-64">
            <span class="absolute inset-y-0 left-0 flex items-center pl-2">
                <svg 
                viewBox="0 0 24 24" 
                class="h-4 w-4 fill-current text-gray-500"
                >
                    <path
                        d="M10 4a6 6 0 100 12 6 6 0 000-12zm-8 6a8 8 0 1114.32 4.906l5.387 5.387a1 1 0 01-1.414 1.414l-5.387-5.387A8 8 0 012 10z">
                    </path>
                </svg>
            </span>
            <input 
                v-model="searchQuery"
                @input="handleSearch"
                placeholder="Search ticker or company"
                class="appearance-none rounded-md border border-gray-300 border-b block pl-8 pr-6 py-2 w-full bg-white text-sm placeholder-gray-400 text-gray-700 focus:bg-white focus:placeholder-gray-600 focus:text-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-shadow duration-200" 
            />
        </div>
    </div>

    <div class="flex-1 overflow-auto bg-white min-h-0 relative">
      <table class="min-w-full leading-normal">
        <thead class="sticky top-0 z-10">
          <tr>

            <th @click="sort('ticker')" class="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider cursor-pointer hover:bg-gray-200 select-none">
              Ticker
              <span v-if="currentParams.sort_by === 'ticker'">
                {{ currentParams.order === 'asc' ? '↑' : '↓' }}
              </span>
            </th>

            <th @click="sort('company')" class="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider cursor-pointer hover:bg-gray-200 select-none">
              Company
              <span v-if="currentParams.sort_by === 'company'">
                {{ currentParams.order === 'asc' ? '↑' : '↓' }}
              </span>
            </th>
            
            <th @click="sort('brokerage')" class="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider cursor-pointer hover:bg-gray-200 select-none">
              Brokerage
              <span v-if="currentParams.sort_by === 'brokerage'">
                {{ currentParams.order === 'asc' ? '↑' : '↓' }}
              </span>
            </th>
            
            <th @click="sort('action')" class="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider cursor-pointer hover:bg-gray-200 select-none">
              Action
              <span v-if="currentParams.sort_by === 'action'">
                {{ currentParams.order === 'asc' ? '↑' : '↓' }}
              </span>
            </th>
            
            <th @click="sort('rating_to')" class="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider cursor-pointer hover:bg-gray-200 select-none">
              Rating
              <span v-if="currentParams.sort_by === 'rating_to'">
                {{ currentParams.order === 'asc' ? '↑' : '↓' }}
              </span>
            </th>
            
            <th @click="sort('target_to_num')" class="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider cursor-pointer hover:bg-gray-200 select-none">
              Target Price
               <span v-if="currentParams.sort_by === 'target_to_num'">
                 {{ currentParams.order === 'asc' ? '↑' : '↓' }}
               </span>
            </th>

          </tr>
        </thead>
        <tbody>
          <tr v-if="loading && stocks.length === 0">
             <td colspan="6" class="px-5 py-10 border-b border-gray-200 bg-white text-sm text-center">
                <div class="flex justify-center items-center">
                    <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
                </div>
            </td>
          </tr>

          <tr v-else-if="stocks.length === 0">
             <td colspan="6" class="px-5 py-10 border-b border-gray-200 bg-white text-sm text-center text-gray-500">
                No stocks found matching your criteria.
            </td>
          </tr>
          
          <tr v-for="stock in stocks" :key="stock.id" class="hover:bg-gray-50 transition-colors duration-150">
            
            <td class="px-5 py-2 border-b border-gray-200 bg-white text-sm">
                <div class="flex items-center">
                   <div class="ml-3">
                        <p class="text-gray-900 font-bold whitespace-no-wrap">
                            {{ stock.ticker }}
                        </p>
                    </div>
                </div>
            </td>

            <td class="px-5 py-2 border-b border-gray-200 bg-white text-sm">
                <p class="text-gray-900 whitespace-no-wrap">{{ stock.company }}</p>
            </td>

            <td class="px-5 py-2 border-b border-gray-200 bg-white text-sm">
                 <p class="text-gray-900 whitespace-no-wrap">{{ stock.brokerage || '-' }}</p>
            </td>
            
            <td class="px-5 py-2 border-b border-gray-200 bg-white text-sm">
                <span :class="{
                    'relative inline-block px-3 py-1 font-semibold leading-tight rounded-full text-xs': true,
                    'text-green-900 bg-green-200': stock.action && stock.action.toLowerCase().includes('raise') || stock.action.toLowerCase().includes('nitiated'),
                    'text-red-900 bg-red-200': stock.action && stock.action.toLowerCase().includes('low') || stock.action.toLowerCase().includes('sell'),
                    'text-orange-900 bg-orange-200': stock.action && stock.action.toLowerCase().includes('hold'),
                    'text-gray-900 bg-gray-200': !stock.action
                }">
                    <span aria-hidden class="absolute inset-0 opacity-50 rounded-full"></span>
                    <span class="relative">{{ stock.action }}</span>
                </span>
            </td>

            <td class="px-5 py-2 border-b border-gray-200 bg-white text-sm">
                <div class="flex flex-col">
                    <span class="font-semibold text-gray-900">{{ stock.rating_to }}</span>
                    <span class="text-xs text-gray-500" v-if="stock.rating_from && stock.rating_from !== stock.rating_to">
                        was {{ stock.rating_from }}
                    </span>
                </div>
            </td>

             <td class="px-5 py-2 border-b border-gray-200 bg-white text-sm">
                 <div class="flex flex-col">
                    <span class="text-gray-900 font-mono font-bold whitespace-no-wrap">
                        {{ stock.target_to ? `$${stock.target_to}` : '-' }}
                    </span>
                     <span class="text-xs text-gray-500" v-if="stock.target_from && stock.target_from !== stock.target_to">
                        was ${{ stock.target_from }}
                    </span>
                 </div>
            </td>

          </tr>
        </tbody>
      </table>
    </div>

    <div class="px-5 py-3 bg-gray-50 border-t border-gray-200 flex flex-col xs:flex-row items-center xs:justify-between flex-shrink-0">
      <span class="text-xs xs:text-sm text-gray-600 mb-2 xs:mb-0">
        Showing {{ (meta.page - 1) * meta.limit + 1 }} to {{ Math.min(meta.page * meta.limit, meta.total) }} of {{ meta.total }} entries
      </span>
      
      <div class="inline-flex rounded-md shadow-sm isolate">
        <button 
          @click="prevPage" 
          :disabled="meta.page === 1"
          class="relative inline-flex items-center px-2 py-2 text-gray-400 bg-white border border-gray-300 rounded-l-md hover:bg-gray-50 focus:z-10 disabled:opacity-50 disabled:cursor-not-allowed transition">
          <span class="sr-only">Previous</span>
          <svg class="w-5 h-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
            <path 
            fill-rule="evenodd" 
            d="M12.79 5.23a.75.75 0 01-.02 1.06L8.832 10l3.938 3.71a.75.75 0 11-1.04 1.08l-4.5-4.25a.75.75 0 010-1.08l4.5-4.25a.75.75 0 011.06.02z" 
            clip-rule="evenodd" 
            />
          </svg>
        </button>

        <template v-for="(page, index) in visiblePages" :key="index">
            
            <button 
                v-if="typeof page === 'number'"
                @click="goToPage(page)"
                :class="[
                    page === meta.page 
                        ? 'z-10 bg-blue-600 text-white focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600' 
                        : 'text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:outline-offset-0',
                    'relative inline-flex items-center px-4 py-2 text-sm font-semibold border-t border-b border-gray-300 focus:z-20 transition'
                ]"
            >
                {{ page }}
            </button>

            <span 
                v-else
                class="relative inline-flex items-center px-4 py-2 text-sm font-semibold text-gray-700 ring-1 ring-inset ring-gray-300 focus:outline-offset-0 border-t border-b border-gray-300 bg-gray-50"
            >
                ...
            </span>
        </template>

        <button 
          @click="nextPage" 
          :disabled="meta.page >= totalPages"
          class="relative inline-flex items-center px-2 py-2 text-gray-400 bg-white border border-gray-300 rounded-r-md hover:bg-gray-50 focus:z-10 disabled:opacity-50 disabled:cursor-not-allowed transition">
          <span class="sr-only">Next</span>
          <svg 
          class="w-5 h-5" 
          viewBox="0 0 20 20" 
          fill="currentColor" 
          aria-hidden="true"
          >
            <path 
            fill-rule="evenodd" 
            d="M7.21 14.77a.75.75 0 01.02-1.06L11.168 10 7.23 6.29a.75.75 0 111.04-1.08l4.5 4.25a.75.75 0 010 1.08l-4.5 4.25a.75.75 0 01-1.06-.02z" 
            clip-rule="evenodd" 
            />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useStockList } from '../composables/useStockList';

const { 
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
} = useStockList();

onMounted(() => {
    init();
});
</script>

