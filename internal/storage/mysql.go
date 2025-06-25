package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

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

	writerDB, err := retryConnect(writerDSN, logger, "writer")
	if err != nil {
		return nil, err
	}

	storage := &MySQLStorage{
		writerDB: writerDB,
		readerDB: writerDB, // Default to writer for reads
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

		readerDB, err := retryConnect(readerDSN, logger, "reader")
		if err != nil {
			logger.Warn("Failed to connect to reader database, will use writer for reads",
				slog.Any("error", err))
		} else {
			storage.readerDB = readerDB
			logger.Info("Successfully connected to MySQL read replica")
		}
	}

	logger.Info("Successfully connected to MySQL writer database")
	return storage, nil
}

// retryConnect attempts to connect to a database with exponential backoff
func retryConnect(dsn string, logger *slog.Logger, dbType string) (*sql.DB, error) {
	maxRetries := 10
	initialDelay := 1 * time.Second
	maxDelay := 30 * time.Second

	var db *sql.DB
	var err error

	for i := range maxRetries {
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			logger.Warn(fmt.Sprintf("Failed to open %s database connection", dbType),
				slog.Int("attempt", i+1),
				slog.Int("max_retries", maxRetries),
				slog.Any("error", err))
		} else {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			err = db.PingContext(ctx)
			cancel()

			if err == nil {
				logger.Info(fmt.Sprintf("Successfully connected to %s database", dbType),
					slog.Int("attempt", i+1))
				return db, nil
			}
			logger.Warn(fmt.Sprintf("Failed to ping %s database", dbType),
				slog.Int("attempt", i+1),
				slog.Int("max_retries", maxRetries),
				slog.Any("error", err))
			db.Close()
		}

		if i < maxRetries-1 {
			delay := time.Duration(1<<uint(i)) * initialDelay
			if delay > maxDelay {
				delay = maxDelay
			}
			logger.Info(fmt.Sprintf("Retrying %s database connection", dbType),
				slog.Duration("delay", delay))
			time.Sleep(delay)
		}
	}

	return nil, fmt.Errorf("failed to connect to %s database after %d attempts: %w", dbType, maxRetries, err)
}

func (ms *MySQLStorage) Ping(ctx context.Context) error {
	if err := ms.writerDB.PingContext(ctx); err != nil {
		return fmt.Errorf("writer database ping failed: %w", err)
	}
	if ms.readerDB != nil && ms.readerDB != ms.writerDB {
		if err := ms.readerDB.PingContext(ctx); err != nil {
			return fmt.Errorf("reader database ping failed: %w", err)
		}
	}
	return nil
}
