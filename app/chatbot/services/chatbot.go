package services

import (
	"context"
	"encoding/json"
	"github.com/born2ngopi/chatbot-3/app/chatbot/request"
	"github.com/born2ngopi/chatbot-3/app/wati"
	"github.com/redis/go-redis/v9"
	"regexp"
	"strconv"
	"time"
)

type (
	ChatbotService interface {
		Webhook(ctx context.Context, req request.WebhookRequest) error
	}

	chatbotService struct {
		redis             *redis.Client
		wati              *wati.WatiApp
		chatbotDuration   time.Duration
		buyValidationText *regexp.Regexp
		findNumberBuy     *regexp.Regexp
	}

	ChatbotStep struct {
		CurrentStep int    `json:"current_step"`
		Location    string `json:"location"`
	}
)

var (
	// change step with template
	mapTemplate = map[int]string{
		0: "first_message",
		1: "second_message",
	}
)

func NewChatbotService(_redis *redis.Client, _wati *wati.WatiApp) ChatbotService {

	// Compile regex pattern
	buyValidationPattern := regexp.MustCompile(`^Beli \d+ Perdana$`)

	// Compile regex pattern
	findNumberBuyPattern := regexp.MustCompile(`Beli (\d+) Perdana`)

	return &chatbotService{
		redis:             _redis,
		wati:              _wati,
		chatbotDuration:   12 * time.Hour,
		buyValidationText: buyValidationPattern,
		findNumberBuy:     findNumberBuyPattern,
	}
}

func (s *chatbotService) deleteAndSetRedis(ctx context.Context, waID string, payload ChatbotStep) error {
	if err := s.redis.Del(ctx, waID).Err(); err != nil {
		return err
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if err := s.redis.Set(ctx, waID, string(b), s.chatbotDuration).Err(); err != nil {
		return err
	}

	return nil
}

func (s *chatbotService) numberToIDRCurrency(number int) string {

	numberStr := strconv.Itoa(number)
	length := len(numberStr)

	if length <= 3 {
		return numberStr
	}

	var formattedNumber string
	commaIndex := length % 3

	if commaIndex == 0 {
		commaIndex = 3
	}

	for i, digit := range numberStr {
		if i == commaIndex {
			formattedNumber += "."
			commaIndex += 3
		}
		formattedNumber += string(digit)
	}

	return "Rp " + formattedNumber

}
