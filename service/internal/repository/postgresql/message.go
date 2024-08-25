package postgresql

import (
	"context"

	"github.com/sedonn/message-service/internal/domain/models"
	"gorm.io/gorm"
)

// Messages возвращает данные о всех сообщениях.
func (r *Repository) Messages(ctx context.Context, pageID, pageSize uint) ([]models.Message, error) {
	var messages []models.Message
	tx := r.db.
		WithContext(ctx).
		Scopes(paginate(pageID, pageSize)).
		Find(&messages)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return messages, nil
}

// ProcessedMessages возвращает данные только обработанных сообщений.
func (r *Repository) ProcessedMessages(ctx context.Context, pageID, pageSize uint) ([]models.Message, error) {
	var messages []models.Message
	tx := r.db.
		WithContext(ctx).
		Scopes(paginate(pageID, pageSize), messageProcessed).
		Find(&messages)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return messages, nil
}

// UnprocessedMessages возвращает данные только необработанных сообщений.
func (r *Repository) UnprocessedMessages(ctx context.Context, pageID, pageSize uint) ([]models.Message, error) {
	var messages []models.Message
	tx := r.db.
		WithContext(ctx).
		Scopes(paginate(pageID, pageSize), messageUnprocessed).
		Find(&messages)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return messages, nil
}

// SaveMessage сохраняет данные нового сообщения.
func (r *Repository) SaveMessage(ctx context.Context, m models.Message) (uint64, error) {
	if tx := r.db.WithContext(ctx).Create(&m); tx.Error != nil {
		return 0, tx.Error
	}

	return m.ID, nil
}

// UpdateMessage обновляет данные существующего сообщения.
func (r *Repository) UpdateMessage(ctx context.Context, m models.Message) (models.Message, error) {
	if tx := r.db.WithContext(ctx).Updates(&m); tx.Error != nil {
		return models.Message{}, tx.Error
	}

	return m, nil
}

// messageProcessed фильтрует только обработанные сообщения.
func messageProcessed(db *gorm.DB) *gorm.DB {
	return db.Where("processed_at IS NOT NULL")
}

// messageUnprocessed фильтрует только необработанные сообщения.
func messageUnprocessed(db *gorm.DB) *gorm.DB {
	return db.Where("processed_at IS NULL")
}
