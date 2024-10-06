package repositories

import (
	"database/sql"
	"fmt"
	"message-automation/src/models"
	"message-automation/src/repositories/queries"
)

type MessageRepository struct {
	DB *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {

	return &MessageRepository{DB: db}
}

func (r *MessageRepository) RetrieveSentMessages(limit int) ([]models.Message, error) {
	rows, err := r.DB.Query(queries.GetSentMessagesQuery(limit))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var message models.Message
		err := rows.Scan(&message.Id, &message.Recipient, &message.Content, &message.Status, &message.CreatedAt, &message.SentAt)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (r *MessageRepository) GetUnsentMessages(limit int) ([]models.Message, error) {
	query := queries.GetUnsentMessagesQuery(limit)
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var message models.Message
		err := rows.Scan(&message.Id, &message.Recipient, &message.Content)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (r *MessageRepository) UpdateMessageStatus(message models.Message) error {
	var err error
	result, err := r.DB.Exec(queries.UpdateMessageStatus(message.Id, "Sent"))
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("message not found")
	}

	return nil
}
