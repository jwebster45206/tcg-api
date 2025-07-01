// Package handlers contains HTTP request handlers for the TCG API
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/deckdef"
	"github.com/jwebster45206/tcg-api/internal/query"
)

const (
	errStrBadRequest = "bad_request"
	errStrValidation = "validation_error"
	errStrInternal   = "internal_error"
	errStrNotFound   = "not_found"
)

// ParseFilters parses filter parameters from an HTTP request using filter array syntax
// Format: filter[field][operator]=value or filter[field]=value (defaults to '=' operator)
// Returns a slice of filters or an error if validation fails
func ParseFilters(r *http.Request, config query.QueryConfig) ([]query.Filter, error) {
	var filters []query.Filter

	// Group values by field and operator combination
	filterMap := make(map[string]map[string][]string)

	for param, values := range r.URL.Query() {
		if !strings.HasPrefix(param, "filter[") || !strings.HasSuffix(param, "]") {
			continue
		}

		// Parse filter[field] or filter[field][operator]
		content := param[7 : len(param)-1] // Remove "filter[" and "]"

		var field, operator string

		// Check if we have nested brackets for operator
		if strings.Contains(content, "][") {
			parts := strings.Split(content, "][")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid filter format")
			}
			field = parts[0]
			operator = parts[1]
		} else {
			field = content
			operator = "=" // Default
		}

		if !config.IsFilterAllowed(field) {
			return nil, fmt.Errorf("filter field '%s' not allowed", field)
		}

		if !isValidOperator(operator) {
			return nil, fmt.Errorf("invalid operator '%s' for field '%s'", operator, field)
		}

		if filterMap[field] == nil {
			filterMap[field] = make(map[string][]string)
		}
		filterMap[field][operator] = append(filterMap[field][operator], values...)
	}

	for field, operatorMap := range filterMap {
		fieldType := config.GetFieldType(field)

		for operatorStr, values := range operatorMap {
			operator := normalizeOperator(operatorStr)

			// Convert values to appropriate type
			var value any
			if len(values) == 1 {
				convertedValue, err := convertFilterValue(values[0], fieldType)
				if err != nil {
					return nil, fmt.Errorf("invalid value for field '%s': %v", field, err)
				}
				value = convertedValue
			} else if len(values) > 1 {
				// Multiple values - create array
				convertedValues := make([]interface{}, len(values))
				for i, v := range values {
					convertedValue, err := convertFilterValue(v, fieldType)
					if err != nil {
						return nil, fmt.Errorf("invalid value for field '%s': %v", field, err)
					}
					convertedValues[i] = convertedValue
				}
				value = convertedValues
			}

			filters = append(filters, query.Filter{
				Column:   field, // Keep API field name, let storage layer handle DB mapping
				Operator: operator,
				Value:    value,
			})
		}
	}
	return filters, nil
}

// ParseSorts parses sort parameters from an HTTP request
// Format: sort=field1,-field2,field3 (- prefix for descending)
// Returns a slice of sort options or an error if validation fails
func ParseSorts(r *http.Request, config query.QueryConfig) ([]query.SortOption, error) {
	var sorts []query.SortOption

	sortParam := r.URL.Query().Get("sort")
	if sortParam == "" {
		return sorts, nil
	}

	// Split by comma
	sortFields := strings.Split(sortParam, ",")

	for _, field := range sortFields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}

		var sortField string
		var desc bool

		// Check for descending prefix
		if strings.HasPrefix(field, "-") {
			desc = true
			sortField = field[1:]
		} else {
			desc = false
			sortField = field
		}

		// Validate field is allowed
		if !config.IsSortAllowed(sortField) {
			return nil, fmt.Errorf("field '%s' is not allowed for sorting", sortField)
		}

		dbColumn, _ := config.GetSortDBColumn(sortField)
		sorts = append(sorts, query.SortOption{
			Field: dbColumn,
			Desc:  desc,
		})
	}

	return sorts, nil
}

// ParsePagination parses pagination parameters from an HTTP request
// Supports both offset-based (offset, limit) and page-based (page, page_size) pagination
// Returns offset, limit values or an error if validation fails
func ParsePagination(r *http.Request) (offset, limit int, err error) {
	// Default values
	limit = 50
	offset = 0

	// Try page-based pagination first
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			return 0, 0, fmt.Errorf("page must be a positive integer")
		}

		pageSize := 50 // Default page size
		if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
			pageSize, err = strconv.Atoi(pageSizeStr)
			if err != nil || pageSize < 1 || pageSize > 100 {
				return 0, 0, fmt.Errorf("page size must be between 1 and 100")
			}
		}

		offset = (page - 1) * pageSize
		limit = pageSize
	} else {
		// Try offset-based pagination
		if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
			offset, err = strconv.Atoi(offsetStr)
			if err != nil || offset < 0 {
				return 0, 0, fmt.Errorf("offset must be a non-negative integer")
			}
		}

		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			limit, err = strconv.Atoi(limitStr)
			if err != nil || limit < 1 || limit > 100 {
				return 0, 0, fmt.Errorf("limit must be between 1 and 100")
			}
		}
	}

	return offset, limit, nil
}

// ErrorResponse represents an error response structure
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// writeJSONResponse safely writes a JSON response with proper error handling
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(jsonData)
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

// isValidOperator checks if the given operator string is valid
func isValidOperator(op string) bool {
	switch strings.ToLower(op) {
	case "=", "eq", "!=", "ne", "not_equal", ">", "gt", ">=", "gte", "greater_equal", "<", "lt", "<=", "lte", "less_equal", "like", "not_like":
		return true
	default:
		return false
	}
}

// normalizeOperator converts operator aliases to standard FilterOperator
func normalizeOperator(op string) query.FilterOperator {
	switch strings.ToLower(op) {
	case "=", "eq":
		return query.OpEqual
	case "!=", "ne", "not_equal":
		return query.OpNotEqual
	case ">", "gt":
		return query.OpGreaterThan
	case ">=", "gte", "greater_equal":
		return query.OpGreaterEqual
	case "<", "lt":
		return query.OpLessThan
	case "<=", "lte", "less_equal":
		return query.OpLessEqual
	case "like":
		return query.OpLike
	case "not_like":
		return query.OpNotLike
	default:
		return query.FilterOperator(op) // Fallback to original
	}
}

// convertFilterValue converts a string value to appropriate type
// convertFilterValue converts a string value to the appropriate type based on field type
func convertFilterValue(value string, fieldType query.FieldType) (interface{}, error) {
	// Handle null/nil
	if value == "" || strings.ToLower(value) == "null" {
		return nil, nil
	}

	switch fieldType {
	case query.FieldTypeUUID:
		return uuid.Parse(value)
	case query.FieldTypeInt:
		return strconv.Atoi(value)
	case query.FieldTypeFloat:
		return strconv.ParseFloat(value, 64)
	case query.FieldTypeBool:
		return strconv.ParseBool(value)
	case query.FieldTypeDateTime:
		// Try RFC3339 first, then other common formats
		if t, err := time.Parse(time.RFC3339, value); err == nil {
			return t, nil
		}
		if t, err := time.Parse("2006-01-02", value); err == nil {
			return t, nil
		}
		if t, err := time.Parse("2006-01-02 15:04:05", value); err == nil {
			return t, nil
		}
		return nil, fmt.Errorf("invalid datetime format: %s", value)
	case query.FieldTypeString:
		fallthrough
	default:
		return value, nil
	}
}

// deckCardsChangedFromInput compares current deck cards with a new card collection input to determine if they've changed
// Returns true if the cards are different (different cards, quantities, or order)
func deckCardsChangedFromInput(current []*deckdef.CardWithQuantity, newCollection *deckdef.CardCollectionInput) bool {
	// If both are empty, no change
	if (len(current) == 0) && (newCollection == nil || len(newCollection.Items) == 0) {
		return false
	}

	// If different number of unique cards, changed
	if len(current) != len(newCollection.Items) {
		return true
	}

	// Create maps for comparison (card ID -> quantity)
	currentMap := make(map[uuid.UUID]int)
	for _, cardWithQty := range current {
		currentMap[cardWithQty.Card.GetID()] = cardWithQty.Quantity
	}

	newMap := make(map[uuid.UUID]int)
	for _, cardInput := range newCollection.Items {
		newMap[cardInput.Card.ID] = cardInput.Quantity
	}

	// Compare the maps
	for cardID, currentQty := range currentMap {
		newQty, exists := newMap[cardID]
		if !exists || currentQty != newQty {
			return true
		}
	}

	return false
}
