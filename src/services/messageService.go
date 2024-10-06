package services

import (
	"context"
	"fmt"
	"message-automation/src/clients"
	"message-automation/src/models"
	"message-automation/src/models/base"
	"message-automation/src/repositories"
	"message-automation/src/validators"
	"strconv"
	"strings"
	"sync"
	"time"
)

type MessageService struct {
	MessageRepository   *repositories.MessageRepository
	WebhookClient       *clients.WebhookClient
	RedisClient         *clients.RedisClient
	IsActiveStatus      bool
	MessagePerExecution int
	ExecutionPeriod     time.Duration
	Mu                  *sync.Mutex
}

func NewMessageService(messageRepository *repositories.MessageRepository, webhookClient *clients.WebhookClient,
	redisClient *clients.RedisClient, messagePerExecution int, executionPeriod time.Duration) *MessageService {
	return &MessageService{MessageRepository: messageRepository, WebhookClient: webhookClient, RedisClient: redisClient,
		IsActiveStatus: true, MessagePerExecution: messagePerExecution, ExecutionPeriod: executionPeriod, Mu: &sync.Mutex{}}
}

func (s *MessageService) RetrieveSentMessages(limitQuery, messageId string) []models.Message {
	var err error
	err = validators.ValidateRetrieveSentMessages(limitQuery, messageId)
	if err != nil {
		panic(err)
	}

	limit, _ := strconv.Atoi(limitQuery)
	messages, err := s.MessageRepository.GetSentMessages(limit, messageId)
	if err != nil {
		panic(&base.BadRequestError{
			Message: fmt.Sprint("Failed to retrieve sent messages"),
		})
	}

	return messages
}

// HandleAutomation starts/stops automation
func (s *MessageService) HandleAutomation(isActiveQuery string) {
	var err error
	if err = validators.ValidateHandleAutomation(isActiveQuery, s.IsActiveStatus); err != nil {
		panic(err)
	}

	isActive, _ := strconv.ParseBool(isActiveQuery)

	s.Mu.Lock()
	s.IsActiveStatus = isActive
	s.Mu.Unlock()

	if isActive {
		go s.ExecuteAutomation()
	}
}

// ExecuteAutomation handles message process which will work with goroutine
func (s *MessageService) ExecuteAutomation() {
	defer func() {
		if r := recover(); r != nil {
			base.Log(fmt.Sprintf("Recovered from panic: %v. Retrying execution after 2 minutes", r))
			time.Sleep(2 * time.Minute)

			// Re-execute message automation after panic recovery
			go s.ExecuteAutomation()
		}
	}()

	for {
		if s.IsActiveStatus == false {
			base.Log(fmt.Sprint("Automation deactivated, stopping execution"))
			return
		} else {
			s.processMessages(s.MessagePerExecution)

			base.Log(fmt.Sprintf("Message automation is executed, next start will after %v minutes", s.ExecutionPeriod.Minutes()))
			time.Sleep(s.ExecutionPeriod)
		}
	}
}

// ExecuteAutomationForProjectDeployment handles message process for the deployment
func (s *MessageService) ExecuteAutomationForProjectDeployment() {
	base.Log(fmt.Sprint("Executing all unsent messages with the project deployment"))
	s.processMessages(0)
	base.Log(fmt.Sprint("Executed all unsent messages successfully"))
}

func (s *MessageService) processMessages(limit int) {
	messages, err := s.MessageRepository.GetUnsentMessages(limit)
	if err != nil {
		base.Log(fmt.Sprintf("Error occurred while acquiring unsent messages: %v", err))

		return
	}

	for _, message := range messages {
		s.processSingleMessage(message)
	}
}

func (s *MessageService) processSingleMessage(message models.Message) {
	ctx := context.Background()
	transaction, err := s.MessageRepository.BeginTransaction(ctx)
	if err != nil {
		base.Log(fmt.Sprintf("Error creating transaction for message %s: %v", message.Id.String(), err))
	}

	if err = s.MessageRepository.UpdateMessageStatus(message, transaction); err != nil {
		_ = transaction.Rollback()
		base.Log(fmt.Sprintf("Error updating message status for message %s: %v", message.Id.String(), err))
	}

	messageId, err := s.sendMessageToWebhook(message)
	if err != nil {
		base.Log(fmt.Sprintf("Error sending message %s: %v", message.Id.String(), err))

		_ = transaction.Rollback()
		return
	}

	if err = transaction.Commit(); err != nil {
		base.Log(fmt.Sprintf("Error committing transaction for message %s: %v", message.Id.String(), err))

		return
	}

	if err = s.cacheMessageSendingTime(messageId); err != nil {
		base.Log(fmt.Sprintf("Error caching sending time for message %s: %v", message.Id.String(), err))
	}
}

func (s *MessageService) sendMessageToWebhook(message models.Message) (string, error) {
	ws, err := s.WebhookClient.SendMessage(message.Recipient, message.Content, message.Id.String())
	if err != nil {
		return "", err
	}

	base.Log(fmt.Sprintf("Successfully sent message: %s", message.Id.String()))
	return ws.MessageId, nil
}

func (s *MessageService) cacheMessageSendingTime(messageId string) error {
	cacheKey := strings.ToLower(messageId)
	cacheValue := time.Now().Format("2006-01-02.15.04.05")

	return s.RedisClient.Set(cacheKey, cacheValue)
}
