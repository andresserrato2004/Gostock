package services

import (
	"BackEnd/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseCurrency(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{"Valid Currency", "$150.00", 150.0},
		{"Valid No Symbol", "200.50", 200.5},
		{"Invalid String", "abc", 0.0},
		{"Empty String", "", 0.0},
		{"With Whitespace", " $ 120.00 ", 120.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCurrency(tt.input)
			if got != tt.expected {
				t.Errorf("parseCurrency(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestCleanData(t *testing.T) {
	s := &StockService{}
	stocks := []models.Stock{
		{TargetFrom: "$100.00", TargetTo: "$110.00"},
		{TargetFrom: "invalid", TargetTo: "50.5"},
	}

	s.cleanData(stocks)

	if stocks[0].TargetFromNum != 100.0 {
		t.Errorf("Expected 100.0, got %f", stocks[0].TargetFromNum)
	}
	if stocks[0].TargetToNum != 110.0 {
		t.Errorf("Expected 110.0, got %f", stocks[0].TargetToNum)
	}
	if stocks[1].TargetFromNum != 0.0 {
		t.Errorf("Expected 0.0 for invalid input, got %f", stocks[1].TargetFromNum)
	}
	if stocks[1].TargetToNum != 50.5 {
		t.Errorf("Expected 50.5, got %f", stocks[1].TargetToNum)
	}
}

func TestFetchWithRetries(t *testing.T) {
	// Mock Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"items": [{"ticker": "AAPL"}], "next_page": ""}`))
	}))
	defer server.Close()

	svc := NewStockService(nil)
	svc.Client = server.Client()

	var target models.ApiResponse
	err := svc.fetchWithRetries2(server.URL, "test-token", &target, 1)

	if err != nil {
		t.Fatalf("fetchWithRetries failed: %v", err)
	}

	if len(target.Data) != 1 {
		t.Errorf("Expected 1 item, got %d", len(target.Data))
	}
	if len(target.Data) > 0 && target.Data[0].Ticker != "AAPL" {
		t.Errorf("Expected ticker AAPL, got %s", target.Data[0].Ticker)
	}
}

func TestFetchWithRetries_FailAfterRetries(t *testing.T) {
	// Este test valida que falle después de N intentos si el server da error siempre
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	svc := NewStockService(nil)
	svc.Client = server.Client()

	var target models.ApiResponse
	// Usamos 1 reintento para no esperar mucho tiempo en el test (duerme 1 segundo)
	err := svc.fetchWithRetries2(server.URL, "token", &target, 1)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}
