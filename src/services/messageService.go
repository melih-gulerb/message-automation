package services

import (
	"fmt"
	"log"
	"message-automation/src/clients"
	"message-automation/src/models"
	"message-automation/src/models/base"
	"message-automation/src/repositories"
	"strconv"
	"time"
)

type MessageService struct {
	MessageRepository *repositories.MessageRepository
	WebhookClient     *clients.WebhookClient
	RedisClient       *clients.RedisClient
	IsActiveStatus    bool
}

func NewMessageService(messageRepository *repositories.MessageRepository, webhookClient *clients.WebhookClient, redisClient *clients.RedisClient) *MessageService {
	return &MessageService{MessageRepository: messageRepository, WebhookClient: webhookClient, RedisClient: redisClient, IsActiveStatus: true}
}

func (s *MessageService) RetrieveSentMessages(limitQuery string) []models.Message {
	var err error
	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		// panic err
	}
	messages, err := s.MessageRepository.RetrieveSentMessages(limit)
	if err != nil {
		// panic err
	}

	return messages
}

func (s *MessageService) HandleAutomation(isActiveQuery string) {
	var err error
	isActive, err := strconv.ParseBool(isActiveQuery)
	if err != nil {
		panic(&base.BadRequestError{Message: fmt.Sprintf("Unable to parse boolean for %s", isActiveQuery)})
	}

	if s.IsActiveStatus == isActive {
		panic(&base.BadRequestError{Message: fmt.Sprintf("Automation is already in requested state")})
	} else {
		s.IsActiveStatus = isActive
		go s.ExecuteAutomation(2)
	}

	return
}

func (s *MessageService) ExecuteAutomation(limit int) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v. Retrying execution after 2 minutes...", r)
			time.Sleep(15 * time.Second)

			// Re-execute message automation after panic recovery
			go s.ExecuteAutomation(limit)
		}
	}()

	for {
		if s.IsActiveStatus == false {
			fmt.Printf("\n[%s] Automation deactivated, stopping execution", time.Now().Format("2006-01-02.15.04.05"))
			return
		} else {
			s.executeAutomation(2)

			fmt.Printf("\n[%s] Message execution will restart after 20 sec", time.Now().Format("2006-01-02.15.04.05"))
			time.Sleep(time.Second * 20)
		}
	}
}

func (s *MessageService) ExecuteAutomationForProjectDeployment() {
	log.Printf("\n[%s] Executing all unsent messages with the project deployment", time.Now().Format("2006-01-02.15.04.05"))
	s.executeAutomation(-1)
	log.Printf("\n[%s] Executed all unsent messages successfully", time.Now().Format("2006-01-02.15.04.05"))
}

func (s *MessageService) executeAutomation(limit int) {
	messages, err := s.MessageRepository.GetUnsentMessages(limit)
	if err != nil {
		fmt.Printf("\n[%s] Error fetching unsent messages: %v", time.Now().Format("2006-01-02.15.04.05"), err)
		return
	}

	for _, message := range messages {
		ws, err := s.WebhookClient.SendMessage(message.Recipient, message.Content)
		if err != nil {
			fmt.Printf("\n[%s] Error sending message to recipient %s: %v", time.Now().Format("2006-01-02.15.04.05"), message.Recipient, err)
			return
		}
		fmt.Printf("\n[%s] Successfully sent message to recipient %s. Webhook response: %v", time.Now().Format("2006-01-02.15.04.05"), message.Recipient, ws)
		if err = s.MessageRepository.UpdateMessageStatus(message); err != nil {
			fmt.Printf("\n[%s] Error updating status for the messageId %s: %v", time.Now().Format("2006-01-02.15.04.05"), message.Id, err)
			return
		}
	}
}
