package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/config"
	"github.com/jwebster45206/tcg-api/internal/models"
)

// MySQLStorage implements Storage interface using MySQL database
type MySQLStorage struct {
	writerDB *sql.DB
	readerDB *sql.DB // Optional read replica, falls back to writerDB if nil
	logger   *slog.Logger
}

// NewMySQLStorage creates a new MySQL storage instance
func NewMySQLStorage(writerConfig config.MySQLConfig, readerConfig *config.MySQLConfig, logger *slog.Logger) (Storage, error) {
	// Initialize writer database connection
	writerDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=UTC",
		writerConfig.User,
		writerConfig.Password,
		writerConfig.Host,
		writerConfig.Port,
		writerConfig.DBName,
	)

	writerDB, err := sql.Open("mysql", writerDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open writer database connection: %w", err)
	}

	// Test writer connection
	if err := writerDB.PingContext(context.Background()); err != nil {
		writerDB.Close()
		return nil, fmt.Errorf("failed to ping writer database: %w", err)
	}

	storage := &MySQLStorage{
		writerDB: writerDB,
		logger:   logger,
	}

	// Initialize reader database connection if config is provided
	if readerConfig != nil && readerConfig.Host != "" {
		readerDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=UTC",
			readerConfig.User,
			readerConfig.Password,
			readerConfig.Host,
			readerConfig.Port,
			readerConfig.DBName,
		)

		readerDB, err := sql.Open("mysql", readerDSN)
		if err != nil {
			logger.Warn("Failed to open reader database connection, will use writer for reads",
				slog.Any("error", err))
		} else {
			// Test reader connection
			if err := readerDB.PingContext(context.Background()); err != nil {
				logger.Warn("Failed to ping reader database, will use writer for reads",
					slog.Any("error", err))
				readerDB.Close()
			} else {
				storage.readerDB = readerDB
				logger.Info("Successfully connected to MySQL read replica")
			}
		}
	}

	logger.Info("Successfully connected to MySQL writer database")
	return storage, nil
}

// Deck operations - stubs
func (m *MySQLStorage) ListDecks(ctx context.Context, ownerID *uuid.UUID) ([]*models.Deck, error) {
	return nil, fmt.Errorf("not implemented")
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

// ImageCard operations - stubs
func (m *MySQLStorage) ListImageCards(ctx context.Context) ([]*models.ImageCard, error) {
	return nil, fmt.Errorf("not implemented")
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

// PlayingCard operations - stubs
func (m *MySQLStorage) ListPlayingCards(ctx context.Context) ([]*models.PlayingCard, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) GetPlayingCard(ctx context.Context, id uuid.UUID) (*models.PlayingCard, error) {
	return nil, fmt.Errorf("not implemented")
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

// GameCard operations - stubs
func (m *MySQLStorage) ListGameCards(ctx context.Context, cardType string) ([]*models.GameCard, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) GetGameCard(ctx context.Context, id uuid.UUID) (*models.GameCard, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) CreateGameCard(ctx context.Context, card models.GameCard) (*models.GameCard, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) UpdateGameCard(ctx context.Context, card models.GameCard) (*models.GameCard, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) DeleteGameCard(ctx context.Context, id uuid.UUID) error {
	return fmt.Errorf("not implemented")
}
