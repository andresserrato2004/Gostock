<script setup lang="ts">
import { useStockStore } from '../stores/stockStore';
import { computed } from 'vue';

const store = useStockStore();

// Computed property to get either error
const errorMessage = computed(() => {
    if (store.stocksError) return store.stocksError;
    if (store.recsError) return store.recsError;
    return null;
});

const closeAlert = () => {
    store.clearErrors();
};
</script>

<template>
  <div v-if="errorMessage" class="fixed top-4 right-4 z-50 animate-bounce-in">
    <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative shadow-md max-w-lg" role="alert">
      <strong class="font-bold block mb-1">Error!</strong>
      <span class="block sm:inline">{{ errorMessage }}</span>
      <span class="absolute top-0 bottom-0 right-0 px-4 py-3" @click="closeAlert">
        <svg class="fill-current h-6 w-6 text-red-500 hover:text-red-800 cursor-pointer" role="button" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20">
            <title>Close</title>
            <path d="M14.348 14.849a1.2 1.2 0 0 1-1.697 0L10 11.819l-2.651 3.029a1.2 1.2 0 1 1-1.697-1.697l2.758-3.15-2.759-3.152a1.2 1.2 0 1 1 1.697-1.697L10 8.183l2.651-3.031a1.2 1.2 0 1 1 1.697 1.697l-2.758 3.152 2.758 3.15a1.2 1.2 0 0 1 0 1.698z"/>
        </svg>
      </span>
    </div>
  </div>
</template>

<style scoped>
@keyframes slideIn {
  from { transform: translateX(100%); opacity: 0; }
  to { transform: translateX(0); opacity: 1; }
}
.animate-bounce-in {
    animation: slideIn 0.3s ease-out forwards;
}
</style>
