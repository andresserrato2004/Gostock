<template>
  <div class="bg-blue-50 border-t-4 border-blue-500 rounded-b-lg lg:rounded-lg lg:border-t-0 lg:border-l-4 p-5 shadow-sm h-fit">
    <div class="flex items-center justify-between cursor-pointer group" @click="isOpen = !isOpen">
      <div class="flex items-center gap-2">
        <div class="p-1.5 bg-blue-100 rounded-md group-hover:bg-blue-200 transition">
             
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-blue-700" viewBox="0 0 20 20" fill="currentColor">
                <path 
                fill-rule="evenodd" 
                d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" 
                clip-rule="evenodd" 
                />
            </svg>
            
        </div>
        <h3 class="text-base font-bold text-blue-900 leading-tight">
            How it Works
        </h3>
      </div>
       <button class="text-blue-500 hover:text-blue-700 focus:outline-none transition-transform duration-300" :class="{ 'rotate-180': isOpen }">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path 
            stroke-linecap="round" 
            stroke-linejoin="round" 
            stroke-width="2" 
            d="M19 9l-7 7-7-7" 
            />
          </svg>
      </button>
    </div>
    
    <p v-if="!isOpen" class="text-xs text-blue-600 mt-2 ml-1 cursor-pointer hover:underline" @click="isOpen = true">
        Understand the Formula 2.0 scoring logic.
    </p>

    <div v-show="isOpen" class="mt-4 text-gray-700 text-sm space-y-5 animate-fade-in border-t border-blue-200 pt-4">
        
        <div class="bg-white p-3 rounded-md border border-blue-100 shadow-sm">
            <p class="font-bold text-xs text-gray-500 uppercase tracking-widest mb-2">Scoring Formula</p>
            <div class="font-mono bg-blue-50/50 p-2 rounded text-center text-xs text-blue-800 border border-blue-100">
                 Score = (0.4 × Upside) + (0.2 × Rating) + (0.2 × Conviction) + (0.2 × Momentum)
            </div>
        </div>

        <div class="space-y-4">
             
             <div class="group">
                <span class="font-bold text-gray-900 block text-sm group-hover:text-blue-700 transition">Final Score</span>
                <p class="text-xs text-gray-600 mt-0.5 leading-relaxed">The proprietary ranking score.</p>
                <div class="text-[10px] text-gray-500 mt-1 pl-2 border-l-2 border-blue-200">
                    High Upside + Strong Analyst Support + Positive Market Momentum = High Score.
                </div>
            </div>

             <div class="group">
                <span class="font-bold text-gray-900 block text-sm group-hover:text-blue-700 transition">Upside Potential (40%)</span>
                <p class="text-xs text-gray-600 mt-0.5 leading-relaxed">Potential profit % if stock hits target.</p>
                <div class="text-[10px] text-gray-500 mt-1 pl-2 border-l-2 border-blue-200">
                    ex: <strong>115.59</strong> = expects to double.
                </div>
            </div>

            <div class="group">
                <span class="font-bold text-gray-900 block text-sm group-hover:text-blue-700 transition">Rating Score (20%)</span>
                <p class="text-xs text-gray-600 mt-0.5 leading-relaxed">Normalized analyst consensus.</p>
                 <div class="text-[10px] text-gray-500 mt-1 pl-2 border-l-2 border-blue-200">
                    1.0 = Strong Buy, 0.8 = Buy.
                </div>
            </div>

            <div class="group">
                <span class="font-bold text-gray-900 block text-sm group-hover:text-blue-700 transition">Conviction (20%)</span>
                <p class="text-xs text-gray-600 mt-0.5 leading-relaxed">Analyst target price trend.</p>
                 <div class="text-[10px] text-gray-500 mt-1 pl-2 border-l-2 border-blue-200">
                    (+) values = Analysts raising targets.<br>
                    (-) values = Analysts lowering targets.
                </div>
            </div>

             <div class="group">
                <span class="font-bold text-gray-900 block text-sm group-hover:text-blue-700 transition">Momentum (20%)</span>
                <p class="text-xs text-gray-600 mt-0.5 leading-relaxed">Daily market price change.</p>
                 <div class="text-[10px] text-gray-500 mt-1 pl-2 border-l-2 border-blue-200">
                    Measures immediate buying pressure.
                </div>
            </div>

        </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
const isOpen = ref(false);
</script>

<style scoped>
.animate-fade-in {
    animation: fadeIn 0.3s ease-in-out;
}
@keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
}
</style>
