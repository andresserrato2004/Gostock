<template>
  <div class="bg-white shadow rounded-lg mb-0 border border-gray-100 h-full flex flex-col">
    <div class="p-6 border-b border-gray-100 flex-shrink-0 flex flex-col md:flex-row md:items-center md:justify-between bg-gray-50 rounded-t-lg">
      <div>
          <h2 class="text-xl font-bold text-gray-800">Top Market Opportunities</h2>
          <p class="text-sm text-gray-500 mt-1">Algorithm (Upside + Consensus + Momentum)</p>
      </div>
      
      <div class="flex items-center gap-4 mt-4 md:mt-0">
          <div class="flex items-center gap-2">
              <label for="limit" class="text-sm font-medium text-gray-700">Show:</label>
              <select 
                id="limit" 
                v-model="currentLimit" 
                @change="handleLimitChange"
                class="block w-full rounded-md border-0 py-1.5 pl-3 pr-8 text-gray-900 ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-indigo-600 sm:text-sm sm:leading-6 cursor-pointer"
              >
                  <option :value="3">Top 3</option>
                  <option :value="5">Top 5</option>
                  <option :value="7">Top 7</option>
                  <option :value="10">Top 10</option>
                  <option :value="20">Top 20</option>
              </select>
          </div>

          <div v-if="loading">
              <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800 animate-pulse">
                  Updating...
              </span>
          </div>
      </div>
    </div>

    <div class="flex-1 overflow-auto min-h-0 relative">
        <div v-if="loading && recommendations.length === 0" class="flex justify-center items-center h-full">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>

        <div v-else-if="!loading && recommendations.length === 0" class="flex justify-center items-center h-full p-8 text-center text-gray-500">
            <p>No recommendations available derived from current market data.</p>
        </div>
        <!-- tengo que ponerle los filtros por busqueda -->
        <table v-else class="min-w-full leading-normal">
            <thead class="sticky top-0 z-10 bg-white">
            <tr>
            <th class="px-5 py-3 border-b-2 border-gray-200 bg-gray-50 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">
              Score (0-100)
            </th>
            <th class="px-5 py-3 border-b-2 border-gray-200 bg-gray-50 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">
              Ticker
            </th>
            <th class="px-5 py-3 border-b-2 border-gray-200 bg-gray-50 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">
              Price Targets
            </th>
            <th class="px-5 py-3 border-b-2 border-gray-200 bg-gray-50 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">
              Upside Potential
            </th>
            <th class="px-5 py-3 border-b-2 border-gray-200 bg-gray-50 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">
              Analysis Factors
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in recommendations" :key="item.ticker" class="hover:bg-gray-50 transition-colors duration-200">
            
            <td class="px-5 py-4 border-b border-gray-200 bg-white text-sm">
                <div class="flex items-center">
                    <div 
                        class="flex items-center justify-center w-12 h-12 rounded-full border-2 font-bold text-sm shadow-sm"
                        :class="getScoreClass(item.final_score)"
                    >
                        {{ item.final_score.toFixed(0) }}
                    </div>
                </div>
            </td>

            <td class="px-5 py-4 border-b border-gray-200 bg-white text-sm">
                <div class="flex flex-col">
                    <span class="text-gray-900 font-black text-lg">{{ item.ticker }}</span>
                    <span class="text-gray-500 text-xs font-medium">{{ item.company }}</span>
                </div>
            </td>

            <td class="px-5 py-4 border-b border-gray-200 bg-white text-sm">
                 <div class="flex flex-col space-y-1">
                    <div class="text-xs text-gray-500 flex justify-between w-32">
                        <span>Current:</span>
                        <span class="font-mono font-medium text-gray-900">{{ formatCurrency(item.current_price) }}</span>
                    </div>
                    <div class="text-xs text-blue-600 flex justify-between w-32 bg-blue-50 px-1 rounded">
                        <span>Target:</span>
                        <span class="font-mono font-bold">{{ formatCurrency(item.target_price) }}</span>
                    </div>
                </div>
            </td>

            <td class="px-5 py-4 border-b border-gray-200 bg-white text-sm">
                <span class="font-mono font-bold text-green-600 text-base">
                    +{{ formatPercent(item.upside) }}
                </span>
            </td>

            <td class="px-5 py-4 border-b border-gray-200 bg-white text-sm">
                <div class="flex flex-col gap-1">
                    <p class="text-gray-800 text-xs bg-gray-100 px-2 py-1 rounded border border-gray-200">
                        {{ item.reason }}
                    </p>
                    <div class="flex gap-2 text-[10px] text-gray-500 mt-1">
                        <span class="bg-purple-50 text-purple-700 px-1 rounded border border-purple-100">Momentum: {{ item.momentum.toFixed(1) }}%</span>
                        <span class="bg-indigo-50 text-indigo-700 px-1 rounded border border-indigo-100">Conviction: {{ item.conviction.toFixed(1) }}</span>
                    </div>
                </div>
            </td>

          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useRecommendations } from '../composables/useRecommendations';

const { 
    recommendations, 
    loading, 
    init, 
    formatCurrency, 
    formatPercent,
    getScoreClass,
    currentLimit,
    handleLimitChange
} = useRecommendations();

onMounted(() => {
    init();
});
</script>