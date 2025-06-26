package storage

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/models"
	"github.com/jwebster45206/tcg-api/internal/query"
)

func (m *MySQLStorage) ListPlayingCards(ctx context.Context, filters []query.Filter, sorts []query.SortOption, pageSize int, pageNum int) ([]*models.PlayingCard, error) {
	// Start building the query with joins
	queryBuilder := squirrel.Select(
		"c.uuid",
		"pc.suit",
		"pc.ranking",
		"c.front_image_url",
		"c.back_image_url",
		"c.created_at",
		"c.updated_at").
		From("cards c").
		Join("playing_cards pc ON c.id = pc.card_id").
		PlaceholderFormat(squirrel.Question)

	queryBuilder = m.applyPlayingCardSystemFilters(queryBuilder)

	for _, filter := range filters {
		if validatedQuery, ok := m.applyValidatedFilter(queryBuilder, filter, models.PlayingCardQueryConfig); ok {
			queryBuilder = validatedQuery
		}
	}
	for _, sort := range sorts {
		if validatedQuery, ok := m.applyValidatedSort(queryBuilder, sort, models.PlayingCardQueryConfig); ok {
			queryBuilder = validatedQuery
		}
	}
	if pageSize > 0 {
		queryBuilder = queryBuilder.Limit(uint64(pageSize))
		if pageNum > 1 {
			offset := uint64((pageNum - 1) * pageSize)
			queryBuilder = queryBuilder.Offset(offset)
		}
	}

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := m.readerDB.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var playingCards []*models.PlayingCard
	for rows.Next() {
		playingCard := &models.PlayingCard{}
		err := rows.Scan(
			&playingCard.ID,
			&playingCard.Suite,
			&playingCard.Ranking,
			&playingCard.FrontImageURL,
			&playingCard.BackImageURL,
			&playingCard.CreatedAt,
			&playingCard.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan playing card: %w", err)
		}
		playingCards = append(playingCards, playingCard)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return playingCards, nil
}

// applyPlayingCardSystemFilters applies mandatory system filters for PlayingCard operations
// These filters ensure data isolation and business rules are always enforced
func (m *MySQLStorage) applyPlayingCardSystemFilters(queryBuilder squirrel.SelectBuilder) squirrel.SelectBuilder {
	return queryBuilder.
		Where(squirrel.Eq{"c.card_type_id": 2}). // Only playing cards
		Where(squirrel.Eq{"c.deleted": false})   // Only non-deleted records
}

func (m *MySQLStorage) GetPlayingCard(ctx context.Context, id uuid.UUID) (*models.PlayingCard, error) {
	filters := []query.Filter{
		{
			Column:   "id",
			Operator: query.OpEqual,
			Value:    id,
		},
	}

	cards, err := m.ListPlayingCards(ctx, filters, []query.SortOption{}, 1, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get playing card: %w", err)
	}
	if len(cards) == 0 {
		return nil, fmt.Errorf("playing card not found")
	}
	return cards[0], nil
}

func (m *MySQLStorage) CreatePlayingCard(ctx context.Context, card models.PlayingCard) (*models.PlayingCard, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) UpdatePlayingCard(ctx context.Context, card models.PlayingCard) (*models.PlayingCard, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) DeletePlayingCard(ctx context.Context, id uuid.UUID) error {
	return fmt.Errorf("not implemented")
}
