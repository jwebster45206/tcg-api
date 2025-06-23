package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jwebster45206/tcg-api/internal/config"
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
