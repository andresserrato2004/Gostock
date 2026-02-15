# Frontend Client 

The user interface for the **GoStock** platform. Built with **Vue 3** and **TypeScript**.

## Tech Stack

*   **Framework**: [Vue 3](https://vuejs.org/) (Composition API, `<script setup>`)
*   **State Management**: [Pinia](https://pinia.vuejs.org/) (Setup Stores Pattern)
*   **Language**: [TypeScript](https://www.typescriptlang.org/) (Strict typing)
*   **Styling**: [Tailwind CSS](https://tailwindcss.com/) (Utility-first)
*   **HTTP Client**: [Axios](https://axios-http.com/)
*   **Build Tool**: [Vite](https://vitejs.dev/) (Fast HMR)

## Directory Structure

```bash
src/
├── assets/             
├── components/         
│   ├── StockList.vue   
│   ├── RecommendationList.vue 
│   ├── AlgorithmExplanation.vue 
│   └── GlobalErrorAlert.vue
├── composables/        
│   ├── useStockList.ts
│   └── useRecommendations.ts
├── services/           
│   └── api.ts         
├── stores/            
│   └── stockStore.ts  
├── types/             
│   └── index.ts
├── App.vue            
└── main.ts            
```

## Development Setup

### Prerequisites
*   Node.js (v18+)
*   pnpm 

### Installation
```bash
pnpm install
```

### Run Locally
```bash
# Start Vite Development Server
pnpm dev
```
The app will launch at `http://localhost:5173`.
