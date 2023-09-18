package services

import (
	"context"
	"github.com/born2ngopi/chatbot-3/app/chatbot/request"
	"github.com/born2ngopi/chatbot-3/app/wati"
)

func (s *chatbotService) StepFour(ctx context.Context, chatbotStep ChatbotStep, req request.WebhookRequest, back bool) error {

	if req.Type != "media" && req.Data == "" {
		if err := s.wati.SendTemplate(
			ctx,
			req.WaID,
			"invalid_payment",
			"Text",
			[]wati.WatiParameters{},
		); err != nil {
			return err
		}
	} else {
		if err := s.wati.SendTemplate(
			ctx,
			req.WaID,
			"finish_message_v5",
			"Text",
			[]wati.WatiParameters{},
		); err != nil {
			return err
		}
	}

	if err := s.redis.Del(ctx, req.WaID).Err(); err != nil {
		return err
	}

	return nil
}
