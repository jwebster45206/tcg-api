package storage

import (
	"context"
	"fmt"
	"reflect"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/models"
	"github.com/jwebster45206/tcg-api/internal/query"
)

// ImageCard operations
func (m *MySQLStorage) ListImageCards(ctx context.Context, filters []query.Filter, sorts []query.SortOption, pageSize int, pageNum int) ([]*models.ImageCard, error) {
	// Start building the query
	query := squirrel.Select(
		"uuid",
		"name",
		"description",
		"front_image_url",
		"back_image_url",
		"created_at",
		"updated_at").
		From("cards").
		PlaceholderFormat(squirrel.Question)

	// Mandatory filters
	query = m.applyImageCardSystemFilters(query)

	// User filters
	for _, filter := range filters {
		if validatedQuery, ok := m.applyValidatedFilter(query, filter, models.ImageCardQueryConfig); ok {
			query = validatedQuery
		}
		// Silently skip invalid filters
	}

	// Apply sorting
	for _, sort := range sorts {
		if validatedQuery, ok := m.applyValidatedSort(query, sort, models.ImageCardQueryConfig); ok {
			query = validatedQuery
		}
		// Silently skip invalid sorts
	}

	if pageSize > 0 {
		query = query.Limit(uint64(pageSize))
		if pageNum > 1 {
			offset := uint64((pageNum - 1) * pageSize)
			query = query.Offset(offset)
		}
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := m.readerDB.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

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

func (m *MySQLStorage) GetImageCard(ctx context.Context, id uuid.UUID) (*models.ImageCard, error) {
	filters := []query.Filter{
		{
			Column:   "id",
			Operator: query.OpEqual,
			Value:    id,
		},
	}

	cards, err := m.ListImageCards(ctx, filters, []query.SortOption{}, 1, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get image card: %w", err)
	}
	if len(cards) == 0 {
		return nil, ErrNotFound
	}
	return cards[0], nil
}

func (m *MySQLStorage) CreateImageCard(ctx context.Context, imageCard models.ImageCard) (*models.ImageCard, error) {
	// Generate UUID if not provided
	if imageCard.ID == uuid.Nil {
		imageCard.ID = uuid.New()
	}

	// Build insert query - pass UUID bytes directly
	query := squirrel.Insert("cards").
		Columns(
			"uuid",
			"name",
			"description",
			"front_image_url",
			"back_image_url",
			"card_type_id",
			"deleted",
		).
		Values(
			imageCard.ID[:],
			imageCard.Name,
			imageCard.Description,
			imageCard.FrontImageURL,
			imageCard.BackImageURL,
			models.TypeImageCardID,
			false,
		).
		PlaceholderFormat(squirrel.Question)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build insert query: %w", err)
	}

	_, err = m.writerDB.ExecContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to insert image card: %w", err)
	}

	// Return the created card by fetching it
	return m.GetImageCard(ctx, imageCard.ID)
}

func (m *MySQLStorage) UpdateImageCard(ctx context.Context, imageCard models.ImageCard) (*models.ImageCard, error) {
	// Build update query
	query := squirrel.Update("cards").
		Set("name", imageCard.Name).
		Set("description", imageCard.Description).
		Set("front_image_url", imageCard.FrontImageURL).
		Set("back_image_url", imageCard.BackImageURL).
		Where(squirrel.Eq{"uuid": imageCard.ID[:]}).
		Where(squirrel.Eq{"card_type_id": models.TypeImageCardID}).
		Where(squirrel.Eq{"deleted": false}).
		PlaceholderFormat(squirrel.Question)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build update query: %w", err)
	}

	result, err := m.writerDB.ExecContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update image card: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, ErrNotFound
	}

	// Return the updated card by fetching it
	return m.GetImageCard(ctx, imageCard.ID)
}

// DeleteImageCard soft-deletes an image card
func (m *MySQLStorage) DeleteImageCard(ctx context.Context, id uuid.UUID) error {
	query := squirrel.Update("cards").
		Set("deleted", true).
		Where(squirrel.Eq{"uuid": id[:]}).
		Where(squirrel.Eq{"card_type_id": models.TypeImageCardID}).
		Where(squirrel.Eq{"deleted": false}).
		PlaceholderFormat(squirrel.Question)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}
	result, err := m.writerDB.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete image card: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// applyImageCardSystemFilters applies mandatory system filters for ImageCard operations
// These filters ensure data isolation and business rules are always enforced
func (m *MySQLStorage) applyImageCardSystemFilters(queryBuilder squirrel.SelectBuilder) squirrel.SelectBuilder {
	return queryBuilder.
		Where(squirrel.Eq{"card_type_id": 1}). // Only image cards
		Where(squirrel.Eq{"deleted": false})   // Only non-deleted records
}

// applyValidatedFilter applies a filter only if it's in the allowed list
func (m *MySQLStorage) applyValidatedFilter(queryBuilder squirrel.SelectBuilder, filter query.Filter, config query.QueryConfig) (squirrel.SelectBuilder, bool) {
	// Check if filter is allowed and get database column
	dbColumn, allowed := config.GetFilterDBColumn(filter.Column)
	if !allowed {
		return queryBuilder, false
	}

	// Apply the filter with the mapped database column
	validatedFilter := query.Filter{
		Column:   dbColumn,
		Operator: filter.Operator,
		Value:    filter.Value,
	}

	return m.applyFilter(queryBuilder, validatedFilter), true
}

// applyValidatedSort applies a sort only if it's in the allowed list
func (m *MySQLStorage) applyValidatedSort(queryBuilder squirrel.SelectBuilder, sort query.SortOption, config query.QueryConfig) (squirrel.SelectBuilder, bool) {
	// Check if sort field is allowed and get database column
	dbColumn, allowed := config.GetSortDBColumn(sort.Field)
	if !allowed {
		return queryBuilder, false
	}

	// Apply the sort with the mapped database column
	direction := "ASC"
	if sort.Desc {
		direction = "DESC"
	}
	return queryBuilder.OrderBy(dbColumn + " " + direction), true
}

// applyFilter applies a single filter to the Squirrel query (now only used internally with validated filters)
func (m *MySQLStorage) applyFilter(queryBuilder squirrel.SelectBuilder, filter query.Filter) squirrel.SelectBuilder {
	if filter.Value == nil {
		// Handle NULL checks
		switch filter.Operator {
		case query.OpEqual:
			return queryBuilder.Where(squirrel.Eq{filter.Column: nil})
		case query.OpNotEqual:
			return queryBuilder.Where(squirrel.NotEq{filter.Column: nil})
		default:
			// For other operators with nil values, skip the filter
			return queryBuilder
		}
	}

	// Special handling for UUID fields stored as BINARY(16)
	if uuidValue, ok := filter.Value.(uuid.UUID); ok {
		// Convert UUID to []byte for direct comparison with BINARY(16) columns
		uuidBytes := uuidValue[:]
		switch filter.Operator {
		case query.OpEqual:
			return queryBuilder.Where(squirrel.Eq{filter.Column: uuidBytes})
		case query.OpNotEqual:
			return queryBuilder.Where(squirrel.NotEq{filter.Column: uuidBytes})
		default:
			// For other operators with UUID values, skip the filter
			return queryBuilder
		}
	}

	// Check if value is a slice (for IN/NOT IN operations)
	val := reflect.ValueOf(filter.Value)
	if val.Kind() == reflect.Slice {
		switch filter.Operator {
		case query.OpEqual:
			return queryBuilder.Where(squirrel.Eq{filter.Column: filter.Value})
		case query.OpNotEqual:
			return queryBuilder.Where(squirrel.NotEq{filter.Column: filter.Value})
		default:
			// For other operators with slice values, skip the filter
			return queryBuilder
		}
	}

	// Handle single values with various operators
	switch filter.Operator {
	case query.OpEqual:
		return queryBuilder.Where(squirrel.Eq{filter.Column: filter.Value})
	case query.OpNotEqual:
		return queryBuilder.Where(squirrel.NotEq{filter.Column: filter.Value})
	case query.OpGreaterThan:
		return queryBuilder.Where(squirrel.Gt{filter.Column: filter.Value})
	case query.OpGreaterEqual:
		return queryBuilder.Where(squirrel.GtOrEq{filter.Column: filter.Value})
	case query.OpLessThan:
		return queryBuilder.Where(squirrel.Lt{filter.Column: filter.Value})
	case query.OpLessEqual:
		return queryBuilder.Where(squirrel.LtOrEq{filter.Column: filter.Value})
	case query.OpLike:
		return queryBuilder.Where(squirrel.Like{filter.Column: filter.Value})
	case query.OpNotLike:
		return queryBuilder.Where(squirrel.NotLike{filter.Column: filter.Value})
	default:
		// Unknown operator, skip the filter
		return queryBuilder
	}
}
