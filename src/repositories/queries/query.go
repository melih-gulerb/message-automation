package queries

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
)

var (
	GetSentMessagesQuery = func(limit int, messageId string) string {
		query := []string{"SELECT"}

		if limit != 0 {
			query = append(query, fmt.Sprintf("TOP %d", limit))
		}

		query = append(query, `Id, Recipient, Content, Status, CreatedAt, SentAt 
                          FROM Messages 
                          WHERE status = 'Sent'`)

		if messageId != "" {
			query = append(query, fmt.Sprintf("AND Id = '%s'", messageId))
		}

		query = append(query, "ORDER BY CreatedAt DESC")

		return strings.Join(query, "\n")
	}

	GetUnsentMessagesQuery = func(limit int) string {
		query := []string{"SELECT"}
		if limit != 0 {
			query = append(query, fmt.Sprintf("TOP %d", limit))
		}
		query = append(query, `Id, Recipient, Content
                          FROM Messages
                          WHERE status = 'Unsent'
                          ORDER BY CreatedAt ASC`)

		return strings.Join(query, "\n")
	}

	UpdateMessageStatus = func(messageId uuid.UUID, status string) string {
		return fmt.Sprintf(`UPDATE Messages SET Status = 'Sent', SentAt = GETDATE() WHERE Id = '%s'`, messageId)
	}
)
