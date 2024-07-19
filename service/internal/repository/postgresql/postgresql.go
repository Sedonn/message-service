package postgresql

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/sedonn/message-service/internal/config"
	"github.com/sedonn/message-service/internal/domain/models"
	"github.com/sedonn/message-service/internal/pkg/logger"
	"github.com/sedonn/message-service/internal/services/message"
)

// Repository содержит методы взаимодействия с базой данных PostgreSQL.
type Repository struct {
	db *gorm.DB
}

var _ message.MessageProvider = (*Repository)(nil)
var _ message.MessageSaver = (*Repository)(nil)

// New создает новый объект репозитория.
func New(cfg *config.Config) (*Repository, error) {
	dsn := makeDSN(&cfg.DB)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 logger.NewGORMLogger(cfg.Env),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.AutoMigrate(&models.Message{}); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &Repository{db: db}, nil
}

// paginate обеспечивает постраничную навигацию в результатах запроса.
func paginate(id, size uint) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(int(size)).Offset(int(id * size))
	}
}

// makeDSN создает строку подключения к базе данных на основе текущей конфигурации.
func makeDSN(cfg *config.DBConfig) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Host, cfg.User, cfg.Password, cfg.Database, cfg.Port)
}
