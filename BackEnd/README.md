# GoStock Backend Documentation & Challenge Report

## 1. Project Overview & Challenge Summary

This project implements the backend component for the "GoStock" application. The goal was to build a system that ingests financial stock data from an external provider, stores it efficiently, and exposes it via a RESTful API with search, filtering, and recommendation capabilities.

### Challenge Requirements Met:
1.  **Data Ingestion**: A robust "Fetcher" service absorbs data from the `api.karenai.click` source. It handles:
    *   Automatic pagination (following `next_page` tokens).
    *   Price parsing (converting string currency fields to comparable numbers).
    *   **Upsert Logic**: Prevents duplicates by updating existing records if a `Ticker` already exists.
2.  **API Development**: A customized REST API built with **Go** and **Gin**.
3.  **Database**: Integration with **CockroachDB** (PostgreSQL-compatible) using **GORM**.
4.  **Features**:
    *   Pagination, Sorting, and Searching.
    *   Specialized Recommendation logic.

---

## 2. Technical Implementation Details

### Architecture
*   **Language**: Go (Golang) 1.23+
*   **Framework**: Gin Web Framework (High performance HTTP web framework)
*   **ORM**: GORM (Go Object Relational Mapping)
*   **Database**: CockroachDB Serverless

### Key Components

#### A. Data Fetcher (`internal/services/fetcher.go`)
The data ingestion engine connects to the external API securely using JWT authentication.
*   **Code Correction**: Fixed a critical bug where `next_page` logic was flawed, ensuring the fetcher consumes the entire dataset, not just the first page.
*   **Data Integrity**: Implemented an `OnConflict` clause for database insertions. This ensures that if we fetch data for "AAPL" today, and fetch it again tomorrow, we update the existing record instead of creating a duplicate.

#### B. Data Models (`internal/models/stock.go`)
*   **Schema**:
    *   `Ticker` (Unique Index): Ensures one record per company.
    *   `TargetFromNum` / `TargetToNum`: Helper fields storing numerical values of price targets for easy sorting/filtering, distinct from the raw string display format.
*   **Fixes**: Added `gorm:"uniqueIndex"` to the Ticker field to enforce database-level integrity, which was critical for the Upsert logic to function.

#### C. Database Configuration (`internal/config/database.go`)
*   **Resilience**: Adjusted the migration logic to log warnings instead of crashing (`log.Fatal`) on minor schema synchronization issues, ensuring high availability even if the DB schema state is slightly drifted.

---

## 3. API Reference

### Base URL: `http://localhost:8080`

### 1. List Stocks
Retrieve a paginated list of stocks with support for searching and sorting.

**Endpoint**: `GET /api/stocks`

**Query Parameters**:
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | int | 1 | Page number. |
| `limit` | int | 10 | Items per page. |
| `search` | string | - | Search term for **Ticker** or **Company Name**. |
| `sort_by` | string | `target_to_num` | Field to sort by (`ticker`, `company`, `rating_to`, `action`). |
| `order` | string | `desc` | Sort order: `asc` or `desc`. |

**Example Request**:
```bash
"http://localhost:8080/api/stocks?page=1&limit=5&search=Apple&sort_by=company&order=asc"
```

**Response Format**:
```json
{
  "data": [
    {
      "id": 105,
      "ticker": "AAPL",
      "company": "Apple",
      "brokerage": "",
      "action": "target raised by",
      "rating_from": "Underweight",
      "rating_to": "Underweight",
      "created_at": "2026-01-24T20:46:31.886Z"
    }
  ],
  "meta": {
    "limit": 5,
    "page": 1,
    "total": 150
  }
}
```

### 2. Get Recommendations (Advanced Algorithm)
Returns a scored and ranked list of stocks combining internal analyst data with real-time market data from Finnhub.

**Endpoint**: `GET /api/recommend`

**Query Parameters**:
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `limit` | int | 10 | Number of top recommendations to return. |

**The Scoring Algorithm (Formula 2.0)**:
The `final_score` (0-100) is a weighted sum designed to find "undervalued winners":
```math
Score = (0.4 \times Upside) + (0.2 \times Rating) + (0.2 \times Conviction) + (0.2 \times Momentum)
```

**Data Glossary & Field Definitions**:

*   **`ticker` / `company`**: Identifiers for the stock (e.g., "NBIS", "Nebius Group").
*   **`current_price`**: Real-time market price fetched live from Finnhub ($97.87).
*   **`target_price`**: The price analysts predict the stock will reach ($211).
*   **`upside`**: The potential profit percentage if the stock reaches the target price.
    *   *Calculation*: `((Target Price - Current Price) / Current Price)`.
    *   *Meaning*: A value of `115.59` means analysts expect the stock to **more than double** (+115%).
*   **`rating_score`**: A normalized value (0.0 to 1.0) representing the consensus rating.
    *   `1.0`: Strong Buy
    *   `0.8`: Buy / Outperform (The model prioritizes these).
*   **`conviction`**: Represents how much analysts "raised" (or lowered) their target price recently.
    *   *Calculation*: `((New Target - Old Target) / Old Target)`.
    *   *Meaning*: A positive value (e.g., `2.43`) means analysts are becoming **more confident/bullish** and have raised their targets.
*   **`momentum`**: The stock's daily price performance (% change today).
    *   *Meaning*: Measures immediate market sentiment. Positive momentum means investors are buying *now*.
*   **`final_score`**: The proprietary ranking score.
    *   *Meaning*: A single number to rank the best opportunities. High Upside + Strong Analyst Support + Positive Market Momentum = High Score.

**Example Request**:
```
"http://localhost:8080/api/recommend?limit=3"
```

**Response Format**:
```json
{
  "data": [
    {
      "ticker": "EEX",
      "company": "Emerald",
      "current_price": 4.72,
      "target_price": 7.7,
      "upside": 63.14,
      "rating_score": 0.8,
      "conviction": -2.53,
      "momentum": 2.61,
      "final_score": 41.27,
      "reason": "Rating: Buy, Action: target lowered by"
    }
  ],
  "meta": null
}
```

---

## 4. How to Run

1.  **Environment Setup**:
    Ensure your `.env` file contains:
    ```
    DATABASE_URL=...
    API_KEY=...
    API_URL=...
    FINNHUB_API_KEY=...
    FINNHUB_API_URL=
    ```

2.  **Start the Server**:
    ```bash
    go run cmd/server/main.go
    ```
    *The server will automatically start ingesting data in the background upon startup.*

3.  **Access**:
    *   API is available at `localhost:8080`.



## Test 
``` bash
go test -v ./...

# coverage

go test -cover ./...
```
### esto es con amor ❤️ andres 
