package models

import (
	"github.com/google/uuid"
	"time"
)

type Message struct {
	Id        uuid.UUID `json:"id"`
	Recipient string    `json:"recipient"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	SentAt    time.Time `json:"sentAt"`
}
