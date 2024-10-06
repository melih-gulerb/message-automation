package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"message-automation/src/models"
	"message-automation/src/models/base"
	"message-automation/src/repositories/queries"
)

type MessageRepository struct {
	DB *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {

	return &MessageRepository{DB: db}
}

// GetSentMessages acquires sent message with filtering the limit and messageId if provided
func (r *MessageRepository) GetSentMessages(limit int, messageId string) ([]models.Message, error) {
	rows, err := r.DB.Query(queries.GetSentMessagesQuery(limit, messageId))
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

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

// GetUnsentMessages acquires unsent message with filtering the limit if provided
func (r *MessageRepository) GetUnsentMessages(limit int) ([]models.Message, error) {
	query := queries.GetUnsentMessagesQuery(limit)
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

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

// UpdateMessageStatus updates the status of a message using a transaction if provided
func (r *MessageRepository) UpdateMessageStatus(message models.Message, transaction *sql.Tx) error {
	var err error
	var result sql.Result

	if transaction != nil {
		result, err = transaction.Exec(queries.UpdateMessageStatus(message.Id, "Sent"))
	} else {
		result, err = r.DB.Exec(queries.UpdateMessageStatus(message.Id, "Sent"))
	}

	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		base.Log(fmt.Sprintf("Couldn't find any message to update"))
	}

	return nil
}

func (r *MessageRepository) BeginTransaction(ctx context.Context) (*sql.Tx, error) {
	transaction, err := r.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
