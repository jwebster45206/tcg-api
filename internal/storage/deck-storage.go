package storage

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/models"
	"github.com/jwebster45206/tcg-api/internal/query"
)

// Deck operations
func (m *MySQLStorage) ListDecks(ctx context.Context, filters []query.Filter, sorts []query.SortOption, pageSize int, pageNum int) ([]*models.Deck, error) {
	// Start building the query with JOIN to get deck type name
	query := squirrel.Select(
		"d.uuid",
		"d.name",
		"d.deck_type_id",
		"dt.name as deck_type_name",
		"d.sleeve_image_url",
		"d.created_at",
		"d.updated_at").
		From("decks d").
		LeftJoin("deck_types dt ON d.deck_type_id = dt.id").
		PlaceholderFormat(squirrel.Question)

	// Apply mandatory system filters
	query = query.Where(squirrel.Eq{"d.deleted": false}) // Only non-deleted records

	// User filters
	for _, filter := range filters {
		if validatedQuery, ok := m.applyValidatedFilter(query, filter, models.DeckQueryConfig); ok {
			query = validatedQuery
		}
		// Silently skip invalid filters
	}

	// Apply sorting
	for _, sort := range sorts {
		if validatedQuery, ok := m.applyValidatedSort(query, sort, models.DeckQueryConfig); ok {
			query = validatedQuery
		}
		// Silently skip invalid sorts
	}

	// Apply pagination
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

	var decks []*models.Deck
	for rows.Next() {
		deck := &models.Deck{}
		var deckTypeName *string
		err := rows.Scan(
			&deck.ID,
			&deck.Name,
			&deck.DeckTypeID,
			&deckTypeName,
			&deck.SleeveImageURL,
			&deck.CreatedAt,
			&deck.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan deck: %w", err)
		}
		deck.DeckTypeName = deckTypeName
		// Cards are not loaded in list operation - they can be requested separately
		decks = append(decks, deck)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return decks, nil
}

func (m *MySQLStorage) GetDeck(ctx context.Context, id uuid.UUID) (*models.Deck, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) CreateDeck(ctx context.Context, deck models.Deck) (*models.Deck, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) UpdateDeck(ctx context.Context, deck models.Deck) (*models.Deck, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) DeleteDeck(ctx context.Context, id uuid.UUID) error {
	return fmt.Errorf("not implemented")
}
