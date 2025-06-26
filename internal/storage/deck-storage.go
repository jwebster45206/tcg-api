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
		"dt.name as deck_type",
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
		err := rows.Scan(
			&deck.ID,
			&deck.Name,
			&deck.DeckType,
			&deck.SleeveImageURL,
			&deck.CreatedAt,
			&deck.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan deck: %w", err)
		}
		// Cards are not loaded in list operation - they can be requested separately
		decks = append(decks, deck)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return decks, nil
}

func (m *MySQLStorage) GetDeck(ctx context.Context, id uuid.UUID) (*models.Deck, error) {
	filters := []query.Filter{
		{
			Column:   "id",
			Operator: query.OpEqual,
			Value:    id,
		},
	}

	decks, err := m.ListDecks(ctx, filters, []query.SortOption{}, 1, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get deck: %w", err)
	}
	if len(decks) == 0 {
		return nil, ErrNotFound
	}
	return decks[0], nil
}

func (m *MySQLStorage) CreateDeck(ctx context.Context, deck models.Deck) (*models.Deck, error) {
	if deck.ID == uuid.Nil {
		deck.ID = uuid.New()
	}

	deckType := deck.DeckType
	if deckType == "" {
		deckType = models.DeckTypeStandard
	}

	query := squirrel.Insert("decks").
		Columns(
			"uuid",
			"name",
			"deck_type_id",
			"sleeve_image_url",
		).
		Values(
			deck.ID[:], // Convert UUID to bytes
			deck.Name,
			squirrel.Select("id").From("deck_types").Where(squirrel.Eq{"name": deckType}),
			deck.SleeveImageURL,
		).
		PlaceholderFormat(squirrel.Question)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build insert query: %w", err)
	}

	_, err = m.writerDB.ExecContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to insert deck: %w", err)
	}

	return m.GetDeck(ctx, deck.ID)
}

func (m *MySQLStorage) UpdateDeck(ctx context.Context, deck models.Deck) (*models.Deck, error) {
	deckType := deck.DeckType
	if deckType == "" {
		deckType = models.DeckTypeStandard
	}
	deckTypeSubquery := squirrel.Select("id").From("deck_types").Where(squirrel.Eq{"name": deckType})

	query := squirrel.Update("decks").
		Where(squirrel.Eq{"uuid": deck.ID[:]}).
		Set("name", deck.Name).
		Set("sleeve_image_url", deck.SleeveImageURL).
		Set("deck_type_id", deckTypeSubquery).
		PlaceholderFormat(squirrel.Question)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build update query: %w", err)
	}
	result, err := m.writerDB.ExecContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update deck: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, ErrNotFound
	}
	return m.GetDeck(ctx, deck.ID)
}

func (m *MySQLStorage) DeleteDeck(ctx context.Context, id uuid.UUID) error {
	query := squirrel.Update("decks").
		Set("deleted", true).
		Where(squirrel.Eq{"uuid": id[:]}).
		Where(squirrel.Eq{"deleted": false}).
		PlaceholderFormat(squirrel.Question)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}
	result, err := m.writerDB.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete deck: %w", err)
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
