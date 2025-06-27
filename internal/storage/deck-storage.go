package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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

func (m *MySQLStorage) ListDeckCards(ctx context.Context, deckID uuid.UUID) ([]*models.CardWithQuantity, error) {
	var internalDeckID int
	err := m.readerDB.QueryRowContext(ctx, "SELECT id FROM decks WHERE uuid = ?", deckID[:]).Scan(&internalDeckID)
	if err != nil {
		return nil, fmt.Errorf("failed to find deck: %w", err)
	}

	query := squirrel.Select(
		"c.uuid",
		"c.name",
		"c.description",
		"c.card_type_id",
		"ct.name as card_type",
		"c.front_image_url",
		"c.back_image_url",
		"pc.suit",
		"pc.ranking",
		"dc.quantity").
		From("deck_cards dc").
		Join("cards c ON dc.card_id = c.id").
		Join("card_types ct ON c.card_type_id = ct.id").
		LeftJoin("playing_cards pc ON c.card_type_id = 2 AND pc.card_id = c.id").
		Where(squirrel.Eq{"dc.deck_id": internalDeckID}).
		PlaceholderFormat(squirrel.Question)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := m.readerDB.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var cards []*models.CardWithQuantity
	for rows.Next() {
		// Basic fields that every card type has
		var cardUUIDBytes []byte
		var name, cardType string
		var description, frontImageURL, backImageURL *string
		var cardTypeID int
		var quantity int

		// Playing card fields (left-joined, so may be nil)
		var suit *string
		var ranking *int

		err = rows.Scan(
			&cardUUIDBytes,
			&name,
			&description,
			&cardTypeID,
			&cardType,
			&frontImageURL,
			&backImageURL,
			&suit,
			&ranking,
			&quantity,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan card: %w", err)
		}

		cardUUID, err := uuid.FromBytes(cardUUIDBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse UUID: %w", err)
		}

		var cardInterface models.CardInterface
		switch cardType {
		case models.TypePlayingCard:
			if suit == nil || ranking == nil {
				return nil, fmt.Errorf("playing card missing required fields")
			}
			cardInterface = &models.PlayingCard{
				ID:            cardUUID,
				Name:          name,
				Description:   safeString(description),
				FrontImageURL: safeString(frontImageURL),
				BackImageURL:  safeString(backImageURL),
				Suit:          safeString(suit),
				Ranking:       safeInt(ranking),
			}
		case models.TypeImageCard:
			cardInterface = &models.ImageCard{
				ID:            cardUUID,
				Name:          name,
				Description:   safeString(description),
				FrontImageURL: safeString(frontImageURL),
				BackImageURL:  safeString(backImageURL),
			}
		case models.TypeGameCard:
			// For game cards, we'd need to fetch additional fields from game_cards table
			// For now, create a basic game card (this would need expansion when game_cards table is added)
			cardInterface = &models.GameCard{
				ID:            cardUUID,
				Name:          name,
				FrontImageURL: safeString(frontImageURL),
				BackImageURL:  safeString(backImageURL),
			}
		default:
			return nil, fmt.Errorf("unknown card type: %s", cardType)
		}

		cardWithQuantity := &models.CardWithQuantity{
			Card:     cardInterface,
			Quantity: quantity,
		}
		cards = append(cards, cardWithQuantity)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return cards, nil
}

func (m *MySQLStorage) SetDeckCards(ctx context.Context, deckID uuid.UUID, cards []models.DeckCardInput) error {
	// Start transaction
	tx, err := m.writerDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get the deck's internal ID and verify it exists
	var internalDeckID int
	err = tx.QueryRowContext(ctx, "SELECT id FROM decks WHERE uuid = ? AND deleted = false", deckID[:]).Scan(&internalDeckID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return fmt.Errorf("failed to find deck: %w", err)
	}

	// Delete existing cards for this deck
	_, err = tx.ExecContext(ctx, "DELETE FROM deck_cards WHERE deck_id = ?", internalDeckID)
	if err != nil {
		return fmt.Errorf("failed to delete existing deck cards: %w", err)
	}

	// If no cards to insert, commit and return
	if len(cards) == 0 {
		return tx.Commit()
	}

	// Prepare batch insert for new cards
	// First, we need to get internal card IDs from UUIDs and validate existence
	cardUUIDs := make([]interface{}, 0, len(cards))
	cardQuantityMap := make(map[uuid.UUID]int)

	for _, cardInput := range cards {
		cardUUIDs = append(cardUUIDs, cardInput.CardID[:])
		cardQuantityMap[cardInput.CardID] = cardInput.Quantity
	}

	// Build query to get internal card IDs and validate that all cards exist
	placeholders := strings.Repeat("?,", len(cardUUIDs)-1) + "?"
	cardQuery := fmt.Sprintf("SELECT id, uuid FROM cards WHERE uuid IN (%s) AND deleted = false", placeholders)

	cardRows, err := tx.QueryContext(ctx, cardQuery, cardUUIDs...)
	if err != nil {
		return fmt.Errorf("failed to query card IDs: %w", err)
	}
	defer cardRows.Close()

	// Build the insert data and validate that all cards exist
	var insertValues []interface{}
	var valuePlaceholders []string
	foundCards := make(map[uuid.UUID]bool)

	for cardRows.Next() {
		var internalCardID int
		var cardUUIDBytes []byte

		err = cardRows.Scan(&internalCardID, &cardUUIDBytes)
		if err != nil {
			return fmt.Errorf("failed to scan card ID: %w", err)
		}

		cardUUID, err := uuid.FromBytes(cardUUIDBytes)
		if err != nil {
			return fmt.Errorf("failed to parse card UUID: %w", err)
		}

		foundCards[cardUUID] = true
		quantity := cardQuantityMap[cardUUID]
		insertValues = append(insertValues, internalDeckID, internalCardID, quantity)
		valuePlaceholders = append(valuePlaceholders, "(?, ?, ?)")
	}

	if err = cardRows.Err(); err != nil {
		return fmt.Errorf("error iterating over card rows: %w", err)
	}

	// Validate that all requested cards were found
	if len(foundCards) != len(cards) {
		return fmt.Errorf("one or more cards not found")
	}

	// Insert the deck cards
	if len(insertValues) > 0 {
		insertQuery := fmt.Sprintf("INSERT INTO deck_cards (deck_id, card_id, quantity) VALUES %s",
			strings.Join(valuePlaceholders, ", "))

		_, err = tx.ExecContext(ctx, insertQuery, insertValues...)
		if err != nil {
			return fmt.Errorf("failed to insert deck cards: %w", err)
		}
	}

	return tx.Commit()
}
