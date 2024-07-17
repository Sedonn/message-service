package postgresql

import (
	"context"
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

// Messages возвращает данные о всех сообщениях.
func (r *Repository) Messages(ctx context.Context, limit uint, offset uint) ([]models.Message, error) {
	var messages []models.Message
	if tx := r.db.Limit(int(limit)).Offset(int(offset)).Find(&messages); tx.Error != nil {
		return nil, tx.Error
	}

	return messages, nil
}

// ProcessedMessages возвращает данные только обработанных сообщений.
func (r *Repository) ProcessedMessages(ctx context.Context, limit uint, offset uint) ([]models.Message, error) {
	var messages []models.Message
	tx := r.db.
		Where("processed_at IS NOT NULL").
		Limit(int(limit)).
		Offset(int(offset)).
		Find(&messages)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return messages, nil
}

// UnprocessedMessages возвращает данные только необработанных сообщений.
func (r *Repository) UnprocessedMessages(ctx context.Context, limit uint, offset uint) ([]models.Message, error) {
	var messages []models.Message
	tx := r.db.
		Where("processed_at IS NULL").
		Limit(int(limit)).
		Offset(int(offset)).
		Find(&messages)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return messages, nil
}

// SaveMessage сохраняет данные нового сообщения.
func (r *Repository) SaveMessage(ctx context.Context, m models.Message) (uint64, error) {
	if tx := r.db.Create(&m); tx.Error != nil {
		return 0, tx.Error
	}

	return m.ID, nil
}

// makeDSN создает строку подключения к базе данных на основе текущей конфигурации.
func makeDSN(cfg *config.DBConfig) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Host, cfg.User, cfg.Password, cfg.Database, cfg.Port)
}
