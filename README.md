# GoStock 

**GoStock Analyst** is a full-stack web application designed to track real-time brokerage ratings and provide algorithmic stock recommendations. It aggregates analyst data (upgrades, downgrades, price targets) and applies a proprietary scoring model Algorithm to identify high-potential market opportunities.

## Features

*   **Real-Time Stock Feed**: View a live feed of stock rating updates from major brokerages (Goldman Sachs, Morgan Stanley, etc.).
*   **Smart Recommendation Engine**: 
    *   **Algorithm**: Ranks stocks based on a weighted formula including Upside Potential, Analyst Consensus, Conviction, and Momentum.
    *   **Top Picks Dashboard**: Filter recommendations by "Top 3", "Top 5", etc.
*   **Interactive UI**:
    *   **Sortable Columns**: Sort by Ticker, Company, Rating, or Target Price.
    *   **Numeric Pagination**

## Tech Stack

### Backend (`/BackEnd`)
Built for performance and concurrency.
*   **Language**: [Go (Golang)](https://go.dev/)
*   **Framework**: [Gin Gonic](https://gin-gonic.com/) (Web Framework)
*   **ORM**: [GORM](https://gorm.io/)
*   **Database**: CockroachDB
*   **Architecture**: Modular (Handlers, Services, Repositories).

### Frontend (`/Frontend`)
Modern, type-safe, and reactive.
*   **Framework**: [Vue 3](https://vuejs.org/) (Composition API + `<script setup>`)
*   **Language**: [TypeScript](https://www.typescriptlang.org/)
*   **State Management**: [Pinia](https://pinia.vuejs.org/) (Setup Stores pattern)
*   **Styling**: [Tailwind CSS](https://tailwindcss.com/)
*   **HTTP Client**: Axios
*   **Build Tool**: [Vite](https://vitejs.dev/)

## Project Structure

```bash
GoStock/
├── BackEnd/              # Go API Server
│   ├── cmd/              # Entry points
│   ├── internal/         # Business logic (Handlers, Models, Services)
│   └── go.mod            # Go dependencies
└── FrontEnd/             # Vue 3 Application
    ├── src/
    │   ├── components/   # UI Components
    │   ├── composables/  # Logic extraction (Hooks)
    │   ├── stores/       # Global State (Pinia)
    │   └── types/        # TypeScript Interfaces
    └── package.json
          
```

## How to Run

### Prerequisites
*   Go 1.25+
*   pnpm
*   CockroachDB

### 1. Backend Setup
```bash
cd BackEnd

go mod download

# Configure Database CockroachDB

# Run Server
go run cmd/server/main.go
```
*Server will start on `http://localhost:8080`*

### 2. Frontend Setup
```bash
cd FrontEnd

# Install dependencies
pnpm install

# Run Development Server
pnpm dev
```
*App will start on `http://localhost:5173`*

## Algorithm Logic

The recommendation engine calculates a **Final Score (0-100)** for each stock based on four key metrics:

1.  **Upside Potential (40%)**: The percentage difference between the current price and the average analyst target price.
2.  **Rating Score (20%)**: Normalized analyst consensus (Buy = 1.0, Hold = 0.5, Sell = 0.0).
3.  **Conviction (20%)**: Tracks the trend of price target adjustments (Are analysts raising or lowering targets?).
4.  **Momentum (20%)**: Recent market price action to gauge buying pressure.

---
© 2026 GoStock con ❤️ andres =3  
