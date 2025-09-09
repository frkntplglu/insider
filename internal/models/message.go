package models

import "time"

type MessageStatus string

const (
	Pending MessageStatus = "pending"
	Sent    MessageStatus = "sent"
	Failed  MessageStatus = "failed"
)

type Message struct {
	Id             string        `gorm:"id" json:"id"`
	RecipientPhone string        `gorm:"recipient_phone" json:"recipientPhone"`
	Content        string        `gorm:"content" json:"content"`
	Status         MessageStatus `gorm:"status" json:"status"`
	CreatedAt      *time.Time    `gorm:"created_at" json:"createdAt"`
	SentAt         *time.Time    `gorm:"sent_at" json:"sentAt"`
}

type MessageSentItem struct {
	MessageId string     `json:"messageId"`
	SentAt    *time.Time `json:"sentAt"`
}
