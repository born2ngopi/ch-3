package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/born2ngopi/chatbot-3/app/chatbot/request"
	"github.com/born2ngopi/chatbot-3/app/wati"
	"github.com/redis/go-redis/v9"
	"strings"
)

func (s *chatbotService) Webhook(ctx context.Context, req request.WebhookRequest) error {
	fmt.Println("check")
	value, err := s.redis.Get(ctx, req.WaID).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return err
		}
		// first message
		if !strings.Contains(strings.ToLower(req.Text), "beli") {
			return nil
		}

		if err := s.wati.SendTemplate(
			ctx,
			req.WaID,
			"first_message_buy_tri",
			"Media",
			[]wati.WatiParameters{},
		); err != nil {
			// server error
			return err
		}

		// set redis
		payload := ChatbotStep{
			CurrentStep: 1,
		}
		b, err := json.Marshal(payload)
		if err != nil {
			// server error
			return err
		}

		if err := s.redis.Set(ctx, req.WaID, string(b), s.chatbotDuration).Err(); err != nil {
			// server error
			return err
		}

		return nil
	}

	var chatbotStep ChatbotStep
	if err := json.Unmarshal([]byte(value), &chatbotStep); err != nil {
		// server error
		return err
	}
	switch chatbotStep.CurrentStep {
	case 1:
		return s.StepOne(ctx, chatbotStep, req, false)
	case 2:
		return s.StepTwo(ctx, chatbotStep, req, false)
	case 3:
		return s.StepThree(ctx, chatbotStep, req, false)
	case 4:
		return s.StepFour(ctx, chatbotStep, req, false)
	case 5:
		return s.StepFive(ctx, chatbotStep, req, false)
	}
	//
	return nil
}
