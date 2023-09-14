package services

import (
	"context"
	"github.com/born2ngopi/chatbot-3/app/chatbot/request"
	"github.com/born2ngopi/chatbot-3/app/wati"
	"strings"
)

func (s *chatbotService) StepFive(ctx context.Context, chatbotStep ChatbotStep, req request.WebhookRequest, back bool) error {

	if strings.ToLower(req.Text) == "kembali ke menu awal" {

		payload := ChatbotStep{
			CurrentStep: 1,
		}

		if err := s.deleteAndSetRedis(ctx, req.WaID, payload); err != nil {
			// server error
			return err
		}

		return s.StepOne(ctx, payload, req, true)
	}

	if req.Type != "media" && req.Data == "" || req.Text == "" {
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
			"finish_message_v4",
			"Text",
			[]wati.WatiParameters{},
		); err != nil {
			return err
		}

	}

	return nil
}
