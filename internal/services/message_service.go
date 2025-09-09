package services

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/frkntplglu/insider/internal/models"
	"github.com/frkntplglu/insider/pkg/logger"
	"github.com/frkntplglu/insider/pkg/ticker"
)

type messageRepository interface {
	RetrieveAllUnsentMessages(ctx context.Context, limit int) ([]models.Message, error)
	UpdateMessage(ctx context.Context, updates map[string]interface{}, message *models.Message) error
}

type redisClient interface {
	GetJson(ctx context.Context, key string, src interface{}) error
	RPush(ctx context.Context, key string, value any) error
	LRange(ctx context.Context, key string, dest interface{}) error
}

type smsClient interface {
	SendSMS(to string, message string) (string, error)
}

type tickerClient interface {
	Start(ctx context.Context)
	Stop()
}

type MessageService struct {
	messageRepo       messageRepository
	redisClient       redisClient
	smsClient         smsClient
	ticker            tickerClient
	mutex             sync.RWMutex
	isRunning         bool
	ctx               context.Context
	cancel            context.CancelFunc
	autoSendingPeriod time.Duration
	redisKey          string
}

func NewMessageService(
	messageRepo messageRepository,
	redisClient redisClient,
	smsClient smsClient,
	redisKey string,
	autoSendingPeriod time.Duration,
) *MessageService {
	return &MessageService{
		messageRepo:       messageRepo,
		redisClient:       redisClient,
		smsClient:         smsClient,
		isRunning:         false,
		autoSendingPeriod: autoSendingPeriod,
		redisKey:          redisKey,
	}
}

func (s *MessageService) GetAllSentMessages(ctx context.Context) ([]models.MessageSentItem, error) {
	var messages []models.MessageSentItem
	err := s.redisClient.LRange(ctx, s.redisKey, &messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *MessageService) StartAutoSending() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.isRunning {
		return errors.New("auto sending is already running")
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())
	newTicker := ticker.NewTimeTicker(s.autoSendingPeriod, s.processPendingMessages)
	s.ticker = &newTicker
	s.isRunning = true

	go newTicker.Start(s.ctx)

	logger.Info("Auto message sending started")

	return nil

}

func (s *MessageService) StopAutoSending() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.isRunning {
		return errors.New("auto sending is not running")
	}

	if s.cancel != nil {
		s.cancel()
	}

	if s.ticker != nil {
		s.ticker.Stop()
	}

	s.isRunning = false

	logger.Info("Auto message sending stopped")
	return nil
}

func (s *MessageService) processPendingMessages() {
	logger.Info("Processing pending messages...")

	messages, err := s.messageRepo.RetrieveAllUnsentMessages(s.ctx, 2)
	if err != nil {
		logger.Error("Error getting pending messages", "error", err)
		return
	}

	if len(messages) == 0 {
		logger.Info("No pending messages found")
		return
	}

	// If our requirements were sending more than two messages, I would consider send messages concurrently. Since we will send just 2 message in 2 minutes sequential is fine.
	for _, message := range messages {
		isSucceed := true
		updates := map[string]interface{}{
			"status": "sent",
		}
		messageId, err := s.sendSingleMessage(message)

		if err != nil {
			updates["status"] = "failed"
			isSucceed = false
			logger.Error("Error sending message", "message_id", message.Id, "error", err)
		}

		now := time.Now()
		updates["sent_at"] = now

		err = s.messageRepo.UpdateMessage(s.ctx, updates, &message)
		if err != nil {
			logger.Error("Error updating message status", "error", err)
		}

		if isSucceed { // only successfully sent message should be cached
			messageSentItem := models.MessageSentItem{
				MessageId: messageId,
				SentAt:    &now,
			}
			err = s.redisClient.RPush(s.ctx, s.redisKey, messageSentItem)
			if err != nil {
				logger.Error("Error caching sent item to Redis", "error", err)
			}
		}
	}
}

func (s *MessageService) sendSingleMessage(message models.Message) (string, error) {
	messageId, err := s.smsClient.SendSMS(message.RecipientPhone, message.Content)
	if err != nil {
		return "", err
	}

	return messageId, nil
}
