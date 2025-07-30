package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/query"
	"github.com/jwebster45206/tcg-api/pkg/deckdef"
)

func (m *MySQLStorage) ListPlayingCards(ctx context.Context, filters []query.Filter, sorts []query.SortOption, pageSize int, pageNum int) ([]*deckdef.PlayingCard, error) {
	queryBuilder := squirrel.Select(
		"c.uuid",
		"c.name",
		"c.description",
		"pc.suit",
		"pc.ranking",
		"c.front_image_url",
		"c.back_image_url",
		"c.created_at",
		"c.updated_at").
		From("cards c").
		Join("playing_cards pc ON c.id = pc.card_id").
		PlaceholderFormat(squirrel.Question)

	queryBuilder = queryBuilder.
		Where(squirrel.Eq{"c.card_type_id": 2}). // Only playing cards
		Where(squirrel.Eq{"c.deleted": false})   // Only non-deleted records

	for _, filter := range filters {
		if validatedQuery, ok := m.applyValidatedFilter(queryBuilder, filter, deckdef.PlayingCardQueryConfig); ok {
			queryBuilder = validatedQuery
		}
	}
	for _, sort := range sorts {
		if validatedQuery, ok := m.applyValidatedSort(queryBuilder, sort, deckdef.PlayingCardQueryConfig); ok {
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
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			m.logger.Warn("Failed to close rows", "error", closeErr)
		}
	}()

	var playingCards []*deckdef.PlayingCard
	for rows.Next() {
		var uuidBytes []byte
		var name, suit string
		var ranking int
		var description, frontImageURL, backImageURL *string
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&uuidBytes,
			&name,
			&description,
			&suit,
			&ranking,
			&frontImageURL,
			&backImageURL,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan playing card: %w", err)
		}

		cardUUID, err := uuid.FromBytes(uuidBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse UUID: %w", err)
		}

		playingCard := &deckdef.PlayingCard{
			ID:            cardUUID,
			Name:          name,
			Description:   safeString(description),
			Suit:          suit,
			Ranking:       ranking,
			FrontImageURL: safeString(frontImageURL),
			BackImageURL:  safeString(backImageURL),
			CreatedAt:     createdAt,
			UpdatedAt:     updatedAt,
		}
		playingCards = append(playingCards, playingCard)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return playingCards, nil
}

func (m *MySQLStorage) GetPlayingCard(ctx context.Context, id uuid.UUID) (*deckdef.PlayingCard, error) {
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

func (m *MySQLStorage) CreatePlayingCard(ctx context.Context, card deckdef.PlayingCard) (*deckdef.PlayingCard, error) {
	if card.ID == uuid.Nil {
		card.ID = uuid.New()
	}

	// Start a transaction for atomic creation
	tx, err := m.writerDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback() // Ignore error - rollback is safe to call multiple times
	}()

	// Insert into cards table first
	cardQuery := squirrel.Insert("cards").
		Columns(
			"uuid",
			"name",
			"description",
			"front_image_url",
			"back_image_url",
			"card_type_id",
		).
		Values(
			card.ID[:],
			card.Name,
			card.Description,
			card.FrontImageURL,
			card.BackImageURL,
			deckdef.TypePlayingCardID,
		).
		PlaceholderFormat(squirrel.Question)

	cardSQL, cardArgs, err := cardQuery.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build cards insert query: %w", err)
	}
	result, err := tx.ExecContext(ctx, cardSQL, cardArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to insert into cards table: %w", err)
	}
	cardID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get inserted card ID: %w", err)
	}

	// Insert into playing_cards lookup table
	playingCardQuery := squirrel.Insert("playing_cards").
		Columns(
			"card_id",
			"suit",
			"ranking",
		).
		Values(
			cardID,
			card.Suit,
			card.Ranking,
		).
		PlaceholderFormat(squirrel.Question)

	playingCardSQL, playingCardArgs, err := playingCardQuery.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build playing_cards insert query: %w", err)
	}
	_, err = tx.ExecContext(ctx, playingCardSQL, playingCardArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to insert into playing_cards table: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return m.GetPlayingCard(ctx, card.ID)
}

func (m *MySQLStorage) UpdatePlayingCard(ctx context.Context, card deckdef.PlayingCard) (*deckdef.PlayingCard, error) {
	if card.ID == uuid.Nil {
		return nil, fmt.Errorf("card ID cannot be nil")
	}

	// Start a transaction for atomic update
	tx, err := m.writerDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback() // Ignore error - rollback is safe to call multiple times
	}()

	// Update cards table
	cardQuery := squirrel.Update("cards").
		Set("name", card.Name).
		Set("description", card.Description).
		Set("front_image_url", card.FrontImageURL).
		Set("back_image_url", card.BackImageURL).
		Where(squirrel.Eq{"uuid": card.ID[:]}).
		Where(squirrel.Eq{"card_type_id": deckdef.TypePlayingCardID}).
		Where(squirrel.Eq{"deleted": false}).
		PlaceholderFormat(squirrel.Question)

	cardSQL, cardArgs, err := cardQuery.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build cards update query: %w", err)
	}
	result, err := tx.ExecContext(ctx, cardSQL, cardArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to update cards table: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("playing card not found")
	}

	// Update playing_cards table - use subquery format
	playingCardQuery := squirrel.Update("playing_cards").
		Set("suit", card.Suit).
		Set("ranking", card.Ranking).
		Where("card_id = (SELECT id FROM cards WHERE uuid = ? AND card_type_id = ? AND deleted = false)",
			card.ID[:], deckdef.TypePlayingCardID).
		PlaceholderFormat(squirrel.Question)

	playingCardSQL, playingCardArgs, err := playingCardQuery.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build playing_cards update query: %w", err)
	}
	_, err = tx.ExecContext(ctx, playingCardSQL, playingCardArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to update playing_cards table: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return m.GetPlayingCard(ctx, card.ID)
}

func (m *MySQLStorage) DeletePlayingCard(ctx context.Context, id uuid.UUID) error {
	query := squirrel.Update("cards").
		Set("deleted", true).
		Where(squirrel.Eq{"uuid": id[:]}).
		Where(squirrel.Eq{"card_type_id": deckdef.TypePlayingCardID}).
		Where(squirrel.Eq{"deleted": false}).
		PlaceholderFormat(squirrel.Question)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}
	result, err := m.writerDB.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete playing card: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("playing card not found")
	}
	return nil
}
