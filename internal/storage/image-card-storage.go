package storage

import (
	"context"
	"fmt"
	"reflect"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/models"
)

// ImageCard operations
func (m *MySQLStorage) ListImageCards(ctx context.Context, filters []Filter, sorts []SortOption, pageSize int, pageNum int) ([]*models.ImageCard, error) {
	// Start building the query
	query := squirrel.Select("id", "name", "description", "front_image_url", "back_image_url", "created_at", "updated_at").
		From("image_cards").
		PlaceholderFormat(squirrel.Question)

	// Apply filters
	for _, filter := range filters {
		query = m.applyFilter(query, filter)
	}

	// Apply sorting
	for _, sort := range sorts {
		direction := "ASC"
		if sort.Desc {
			direction = "DESC"
		}
		query = query.OrderBy(sort.Field + " " + direction)
	}

	// Apply pagination
	if pageSize > 0 {
		query = query.Limit(uint64(pageSize))
		if pageNum > 1 {
			offset := uint64((pageNum - 1) * pageSize)
			query = query.Offset(offset)
		}
	}

	// Generate SQL and args
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	// Execute query
	rows, err := m.readerDB.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Scan results
	var imageCards []*models.ImageCard
	for rows.Next() {
		imageCard := &models.ImageCard{}
		err := rows.Scan(
			&imageCard.ID,
			&imageCard.Name,
			&imageCard.Description,
			&imageCard.FrontImageURL,
			&imageCard.BackImageURL,
			&imageCard.CreatedAt,
			&imageCard.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan image card: %w", err)
		}
		imageCards = append(imageCards, imageCard)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return imageCards, nil
}

// applyFilter applies a single filter to the Squirrel query
func (m *MySQLStorage) applyFilter(query squirrel.SelectBuilder, filter Filter) squirrel.SelectBuilder {
	if filter.Value == nil {
		// Handle NULL checks
		switch filter.Operator {
		case OpEqual:
			return query.Where(squirrel.Eq{filter.Column: nil})
		case OpNotEqual:
			return query.Where(squirrel.NotEq{filter.Column: nil})
		default:
			// For other operators with nil values, skip the filter
			return query
		}
	}

	// Check if value is a slice (for IN/NOT IN operations)
	val := reflect.ValueOf(filter.Value)
	if val.Kind() == reflect.Slice {
		switch filter.Operator {
		case OpEqual:
			return query.Where(squirrel.Eq{filter.Column: filter.Value})
		case OpNotEqual:
			return query.Where(squirrel.NotEq{filter.Column: filter.Value})
		default:
			// For other operators with slice values, skip the filter
			return query
		}
	}

	// Handle single values with various operators
	switch filter.Operator {
	case OpEqual:
		return query.Where(squirrel.Eq{filter.Column: filter.Value})
	case OpNotEqual:
		return query.Where(squirrel.NotEq{filter.Column: filter.Value})
	case OpGreaterThan:
		return query.Where(squirrel.Gt{filter.Column: filter.Value})
	case OpGreaterEqual:
		return query.Where(squirrel.GtOrEq{filter.Column: filter.Value})
	case OpLessThan:
		return query.Where(squirrel.Lt{filter.Column: filter.Value})
	case OpLessEqual:
		return query.Where(squirrel.LtOrEq{filter.Column: filter.Value})
	case OpLike:
		return query.Where(squirrel.Like{filter.Column: filter.Value})
	case OpNotLike:
		return query.Where(squirrel.NotLike{filter.Column: filter.Value})
	default:
		// Unknown operator, skip the filter
		return query
	}
}

func (m *MySQLStorage) GetImageCard(ctx context.Context, id uuid.UUID) (*models.ImageCard, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) CreateImageCard(ctx context.Context, imageCard models.ImageCard) (*models.ImageCard, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) UpdateImageCard(ctx context.Context, imageCard models.ImageCard) (*models.ImageCard, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) DeleteImageCard(ctx context.Context, id uuid.UUID) error {
	return fmt.Errorf("not implemented")
}
