<script setup lang="ts">
import { ref } from 'vue';
import StockList from './components/StockList.vue';
import RecommendationList from './components/RecommendationList.vue';
import AlgorithmExplanation from './components/AlgorithmExplanation.vue';

const activeTab = ref('stocks');
</script>

<template>
  <div class="h-screen flex flex-col bg-gray-100 font-sans overflow-hidden">
    
    <header class="bg-blue-600 shadow-md flex-shrink-0 z-10">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16 items-center">
            
            <div class="flex items-center">
                <div class="flex-shrink-0 flex items-center text-white gap-2">
                     <svg class="h-8 w-8 text-blue-200" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path 
                        stroke-linecap="round" 
                        stroke-linejoin="round" 
                        stroke-width="2" 
                        d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" 
                        />
                    </svg>
                    <span class="font-bold text-xl tracking-tight">GoStock Analyst</span>
                </div>
            </div>

             <div class="hidden md:block">
                <div class="ml-10 flex items-baseline space-x-4">
                    <button 
                        @click="activeTab = 'stocks'"
                        :class="[
                            activeTab === 'stocks' 
                                ? 'bg-blue-700 text-white shadow-inner' 
                                : 'text-blue-100 hover:bg-blue-500 hover:text-white', 
                            'px-3 py-2 rounded-md text-sm font-medium transition duration-150 cursor-pointer focus:outline-none'
                        ]"
                    >
                        Stock List
                    </button>

                    <button 
                        @click="activeTab = 'recommendations'"
                        :class="[
                            activeTab === 'recommendations' 
                                ? 'bg-blue-700 text-white shadow-inner' 
                                : 'text-blue-100 hover:bg-blue-500 hover:text-white', 
                            'px-3 py-2 rounded-md text-sm font-medium transition duration-150 cursor-pointer focus:outline-none'
                        ]"
                    >
                        Recommendations
                    </button>
                </div>
            </div>
            
        </div>
      </div>
    </header>

    <main class="flex-1 flex flex-col max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 w-full min-h-0">
      
        <div class="mb-4 flex-shrink-0">
             <h1 class="text-2xl font-bold text-gray-900 leading-tight">
                {{ activeTab === 'stocks' ? 'Market Data Explorer' : 'Analyst Recommendations' }}
            </h1>
            <p class="mt-1 text-sm text-gray-600">
                {{ activeTab === 'stocks' 
                    ? 'Search, filter, and analyze real-time brokerage ratings and targets.' 
                    : 'Algorithmic picks based on Upside Potential, Consensus, and Momentum.' 
                }}
            </p>
        </div>

        <div v-if="activeTab === 'stocks'" class="flex-1 min-h-0 animate-fade-in-up">
            <StockList />
        </div>

        <div v-if="activeTab === 'recommendations'" class="flex-1 min-h-0 animate-fade-in-up">
             <div class="flex flex-col lg:flex-row gap-6 items-start h-full">
                 <div class="flex-1 w-full min-w-0 h-full">
                     <RecommendationList />
                 </div>
                 
                 <div class="w-full lg:w-96 flex-shrink-0 overflow-y-auto max-h-full pr-2 pb-4">
                     <AlgorithmExplanation />
                 </div>
             </div>
        </div>

    </main>

    <footer class="bg-white border-t border-gray-200 py-4 flex-shrink-0 z-10">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <p class="text-center text-sm text-gray-500">
                esto esta hecho con mucho amor ❤️ por Andres Serrato Camero
            </p>
        </div>
    </footer>

  </div>
</template>

<style>
.animate-fade-in-up {
  animation: fadeInUp 0.5s ease-out;
  height: 100%; 
  display: flex;
  flex-direction: column;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
