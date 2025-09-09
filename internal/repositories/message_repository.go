package repositories

import (
	"context"

	"github.com/frkntplglu/insider/internal/models"
)

type MessageRepository struct {
	db Database
}

func NewMessageRepository(db Database) *MessageRepository {
	return &MessageRepository{
		db: db,
	}
}

func (r *MessageRepository) RetrieveAllUnsentMessages(ctx context.Context, limit int) ([]models.Message, error) {
	var messages []models.Message
	result := r.db.GetConnection().WithContext(ctx).Model(models.Message{}).Where("status = ?", models.Pending).Order("created_at desc").Limit(limit).Find(&messages)

	if result.Error != nil {
		return []models.Message{}, result.Error
	}

	return messages, nil
}

func (r *MessageRepository) UpdateMessage(ctx context.Context, updates map[string]interface{}, message *models.Message) error {
	result := r.db.GetConnection().WithContext(ctx).Model(message).Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
