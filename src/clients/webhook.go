package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type WebhookClient struct {
	Client *http.Client
	URL    string
	Token  string
}

func NewWebhookClient(url, token string) *WebhookClient {
	return &WebhookClient{
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
		URL:   url,
		Token: token,
	}
}

type MessagePayload struct {
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
}

type WebhookResponse struct {
	Message   string `json:"message"`
	MessageId string `json:"messageId"`
}

func (w *WebhookClient) SendMessage(recipient, content, messageId string) (*WebhookResponse, error) {
	payload := MessagePayload{
		Recipient: recipient,
		Content:   content,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := w.URL + fmt.Sprintf("/%s", w.Token)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := w.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusAccepted {
		return nil, errors.New("failed to send message")
	}

	var webhookResp WebhookResponse
	if err := json.NewDecoder(resp.Body).Decode(&webhookResp); err != nil {
		return nil, err
	}
	// Mocked because the lack of CustomAction in webhook.site
	webhookResp.MessageId = messageId

	return &webhookResp, nil
}
