package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/models"
	"github.com/jwebster45206/tcg-api/internal/query"
	"github.com/jwebster45206/tcg-api/internal/storage"
)

// mustParseTime is a helper function for tests that parses a date string and panics on error
func mustParseTime(dateStr string) time.Time {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		panic(fmt.Sprintf("failed to parse time %s: %v", dateStr, err))
	}
	return t
}

func TestImageCardsHandler_ListCards(t *testing.T) {
	req, err := http.NewRequest("GET", "/image-cards", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Create handler with dependencies
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewImageCardsHandler(mockStorage, logger)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var cards []interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &cards); err != nil {
		t.Errorf("Could not parse response body: %v", err)
	}

	// Should return an empty list since mock storage starts empty
	if len(cards) != 0 {
		t.Errorf("Expected empty card list, got %d cards", len(cards))
	}
}

func TestImageCardsHandler_GetCard(t *testing.T) {
	// Test with valid UUID
	cardID := uuid.New()

	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewImageCardsHandler(mockStorage, logger)

	// Test getting non-existent card
	req, err := http.NewRequest("GET", "/image-cards/"+cardID.String(), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

	// Test with invalid UUID format
	req, err = http.NewRequest("GET", "/image-cards/invalid-uuid", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestImageCardsHandler_CreateCard(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewImageCardsHandler(mockStorage, logger)

	card := models.ImageCard{
		Name:          "Test Image Card",
		Description:   "A test image card",
		FrontImageURL: "https://example.com/front.jpg",
		BackImageURL:  "https://example.com/back.jpg",
	}

	jsonData, err := json.Marshal(card)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/image-cards", bytes.NewReader(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var createdCard models.ImageCard
	if err := json.Unmarshal(rr.Body.Bytes(), &createdCard); err != nil {
		t.Errorf("Could not parse response body: %v", err)
	}

	if createdCard.Name != card.Name {
		t.Errorf("Expected card name %s, got %s", card.Name, createdCard.Name)
	}

	if createdCard.ID == uuid.Nil {
		t.Error("Expected created card to have an ID")
	}
}

func TestImageCardsHandler_CreateCard_InvalidJSON(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewImageCardsHandler(mockStorage, logger)

	req, err := http.NewRequest("POST", "/image-cards", bytes.NewReader([]byte("invalid json")))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestImageCardsHandler_UpdateCard(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewImageCardsHandler(mockStorage, logger)

	// First create a card
	card := models.ImageCard{
		Name:          "Original Image Card",
		Description:   "Original description",
		FrontImageURL: "https://example.com/front.jpg",
		BackImageURL:  "https://example.com/back.jpg",
	}

	createdCard, err := mockStorage.CreateImageCard(context.Background(), card)
	if err != nil {
		t.Fatal(err)
	}

	// Update the card
	updatedCard := models.ImageCard{
		ID:            createdCard.ID,
		Name:          "Updated Image Card",
		Description:   "Updated description",
		FrontImageURL: "https://example.com/new-front.jpg",
		BackImageURL:  "https://example.com/new-back.jpg",
	}

	jsonData, err := json.Marshal(updatedCard)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/image-cards/"+createdCard.ID.String(), bytes.NewReader(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var responseCard models.ImageCard
	if err := json.Unmarshal(rr.Body.Bytes(), &responseCard); err != nil {
		t.Errorf("Could not parse response body: %v", err)
	}

	if responseCard.Name != updatedCard.Name {
		t.Errorf("Expected card name %s, got %s", updatedCard.Name, responseCard.Name)
	}
}

func TestImageCardsHandler_DeleteCard(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewImageCardsHandler(mockStorage, logger)

	// First create a card
	card := models.ImageCard{
		Name:          "Card to Delete",
		Description:   "This card will be deleted",
		FrontImageURL: "https://example.com/front.jpg",
		BackImageURL:  "https://example.com/back.jpg",
	}

	createdCard, err := mockStorage.CreateImageCard(context.Background(), card)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("DELETE", "/image-cards/"+createdCard.ID.String(), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	// Verify the card was deleted
	_, err = mockStorage.GetImageCard(context.Background(), createdCard.ID)
	if err == nil {
		t.Error("Expected card to be deleted, but it still exists")
	}
}

func TestImageCardsHandler_MethodNotAllowed(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewImageCardsHandler(mockStorage, logger)

	req, err := http.NewRequest("PATCH", "/image-cards", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}

func TestParseFilters(t *testing.T) {
	tests := []struct {
		name        string
		queryParams string
		expected    []query.Filter
		expectError bool
	}{
		{
			name:        "simple filter with equals",
			queryParams: "filter[name]=john",
			expected: []query.Filter{
				{Column: "name", Operator: query.OpEqual, Value: "john"},
			},
			expectError: false,
		},
		{
			name:        "filter with operator",
			queryParams: "filter[name][like]=john",
			expected: []query.Filter{
				{Column: "name", Operator: query.OpLike, Value: "john"},
			},
			expectError: false,
		},
		{
			name:        "multiple values for same field",
			queryParams: "filter[name]=john&filter[name]=jane",
			expected: []query.Filter{
				{Column: "name", Operator: query.OpEqual, Value: []interface{}{"john", "jane"}},
			},
			expectError: false,
		},
		{
			name:        "UUID filter",
			queryParams: "filter[id]=123e4567-e89b-12d3-a456-426614174000",
			expected: []query.Filter{
				{Column: "uuid", Operator: query.OpEqual, Value: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")},
			},
			expectError: false,
		},
		{
			name:        "null value filter",
			queryParams: "filter[name]=null",
			expected: []query.Filter{
				{Column: "name", Operator: query.OpEqual, Value: nil},
			},
			expectError: false,
		},
		{
			name:        "multiple different filters",
			queryParams: "filter[name]=john&filter[created_at][gte]=2023-01-01",
			expected: []query.Filter{
				{Column: "name", Operator: query.OpEqual, Value: "john"},
				{Column: "created_at", Operator: query.OpGreaterEqual, Value: mustParseTime("2023-01-01")},
			},
			expectError: false,
		},
		{
			name:        "disallowed field",
			queryParams: "filter[secret]=value",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "invalid operator",
			queryParams: "filter[name][invalid]=value",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "no filters",
			queryParams: "",
			expected:    []query.Filter{},
			expectError: false,
		},
		{
			name:        "non-filter parameters ignored",
			queryParams: "page=1&filter[name]=john&limit=10",
			expected: []query.Filter{
				{Column: "name", Operator: query.OpEqual, Value: "john"},
			},
			expectError: false,
		},
		{
			name:        "invalid UUID format",
			queryParams: "filter[id]=invalid-uuid",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "invalid datetime format",
			queryParams: "filter[created_at]=invalid-date",
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request with query parameters
			req := httptest.NewRequest("GET", "/image-cards?"+tt.queryParams, nil)

			filters, err := ParseFilters(req, models.ImageCardQueryConfig)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(filters) != len(tt.expected) {
				t.Errorf("Expected %d filters, got %d", len(tt.expected), len(filters))
				return
			}

			// Sort filters by column for consistent comparison
			sort.Slice(filters, func(i, j int) bool {
				return filters[i].Column < filters[j].Column
			})
			sort.Slice(tt.expected, func(i, j int) bool {
				return tt.expected[i].Column < tt.expected[j].Column
			})

			for i, expected := range tt.expected {
				actual := filters[i]
				if actual.Column != expected.Column {
					t.Errorf("Filter %d: expected column %s, got %s", i, expected.Column, actual.Column)
				}
				if actual.Operator != expected.Operator {
					t.Errorf("Filter %d: expected operator %s, got %s", i, expected.Operator, actual.Operator)
				}

				// Compare values - handle arrays specially
				if !compareFilterValues(actual.Value, expected.Value) {
					t.Errorf("Filter %d: expected value %v, got %v", i, expected.Value, actual.Value)
				}
			}
		})
	}
}

// compareFilterValues compares two filter values, handling arrays and different types
func compareFilterValues(actual, expected interface{}) bool {
	// Handle nil values
	if actual == nil && expected == nil {
		return true
	}
	if actual == nil || expected == nil {
		return false
	}

	// Handle arrays
	actualSlice, actualIsSlice := actual.([]interface{})
	expectedSlice, expectedIsSlice := expected.([]interface{})

	if actualIsSlice && expectedIsSlice {
		if len(actualSlice) != len(expectedSlice) {
			return false
		}

		// Sort both slices as strings for comparison
		actualStrs := make([]string, len(actualSlice))
		expectedStrs := make([]string, len(expectedSlice))

		for i, v := range actualSlice {
			actualStrs[i] = fmt.Sprintf("%v", v)
		}
		for i, v := range expectedSlice {
			expectedStrs[i] = fmt.Sprintf("%v", v)
		}

		sort.Strings(actualStrs)
		sort.Strings(expectedStrs)

		for i := range actualStrs {
			if actualStrs[i] != expectedStrs[i] {
				return false
			}
		}
		return true
	}

	if actualIsSlice || expectedIsSlice {
		return false // One is slice, other isn't
	}

	// Handle regular values
	return fmt.Sprintf("%v", actual) == fmt.Sprintf("%v", expected)
}

func TestParseSorts(t *testing.T) {
	tests := []struct {
		name        string
		queryParams string
		expected    []query.SortOption
		wantError   bool
		errorMsg    string
	}{
		{
			name:        "no sort parameter",
			queryParams: "",
			expected:    []query.SortOption{},
			wantError:   false,
		},
		{
			name:        "single sort ascending",
			queryParams: "sort=created_at",
			expected: []query.SortOption{
				{Field: "created_at", Desc: false},
			},
			wantError: false,
		},
		{
			name:        "single sort descending",
			queryParams: "sort=-created_at",
			expected: []query.SortOption{
				{Field: "created_at", Desc: true},
			},
			wantError: false,
		},
		{
			name:        "multiple sorts",
			queryParams: "sort=created_at,-updated_at",
			expected: []query.SortOption{
				{Field: "created_at", Desc: false},
				{Field: "updated_at", Desc: true},
			},
			wantError: false,
		},
		{
			name:        "disallowed sort field",
			queryParams: "sort=secret_field",
			expected:    nil,
			wantError:   true,
			errorMsg:    "Field 'secret_field' is not allowed for sorting",
		},
		{
			name:        "mixed allowed and disallowed",
			queryParams: "sort=created_at,secret_field",
			expected:    nil,
			wantError:   true,
			errorMsg:    "Field 'secret_field' is not allowed for sorting",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request with query parameters
			req := httptest.NewRequest("GET", "/image-cards?"+tt.queryParams, nil)

			sorts, err := ParseSorts(req, models.ImageCardQueryConfig)

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(sorts) != len(tt.expected) {
				t.Errorf("Expected %d sorts, got %d", len(tt.expected), len(sorts))
				return
			}

			for i, expected := range tt.expected {
				if sorts[i].Field != expected.Field || sorts[i].Desc != expected.Desc {
					t.Errorf("Sort %d: expected %+v, got %+v", i, expected, sorts[i])
				}
			}
		})
	}
}

func TestParsePagination(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		expectedOffset int
		expectedLimit  int
		wantError      bool
		errorMsg       string
	}{
		{
			name:           "no pagination parameters",
			queryParams:    "",
			expectedOffset: 0,
			expectedLimit:  50,
			wantError:      false,
		},
		{
			name:           "page-based pagination",
			queryParams:    "page=2&page_size=25",
			expectedOffset: 25,
			expectedLimit:  25,
			wantError:      false,
		},
		{
			name:           "page only (default page size)",
			queryParams:    "page=3",
			expectedOffset: 100,
			expectedLimit:  50,
			wantError:      false,
		},
		{
			name:           "offset-based pagination",
			queryParams:    "offset=10&limit=20",
			expectedOffset: 10,
			expectedLimit:  20,
			wantError:      false,
		},
		{
			name:           "offset only (default limit)",
			queryParams:    "offset=15",
			expectedOffset: 15,
			expectedLimit:  50,
			wantError:      false,
		},
		{
			name:        "invalid page",
			queryParams: "page=0",
			wantError:   true,
			errorMsg:    "Page must be a positive integer",
		},
		{
			name:        "invalid page size",
			queryParams: "page=1&page_size=150",
			wantError:   true,
			errorMsg:    "Page size must be between 1 and 100",
		},
		{
			name:        "invalid offset",
			queryParams: "offset=-5",
			wantError:   true,
			errorMsg:    "Offset must be a non-negative integer",
		},
		{
			name:        "invalid limit",
			queryParams: "limit=150",
			wantError:   true,
			errorMsg:    "Limit must be between 1 and 100",
		},
		{
			name:        "non-numeric page",
			queryParams: "page=abc",
			wantError:   true,
			errorMsg:    "Page must be a positive integer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request with query parameters
			req := httptest.NewRequest("GET", "/image-cards?"+tt.queryParams, nil)

			offset, limit, err := ParsePagination(req)

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if offset != tt.expectedOffset {
				t.Errorf("Expected offset %d, got %d", tt.expectedOffset, offset)
			}

			if limit != tt.expectedLimit {
				t.Errorf("Expected limit %d, got %d", tt.expectedLimit, limit)
			}
		})
	}
}
