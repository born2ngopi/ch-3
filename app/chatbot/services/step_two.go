package services

import (
	"context"
	"github.com/born2ngopi/chatbot-3/app/chatbot/request"
	"github.com/born2ngopi/chatbot-3/app/wati"
)

func (s *chatbotService) StepTwo(ctx context.Context, chatbotStep ChatbotStep, req request.WebhookRequest, back bool) error {

	text := req.Text

	if text == "OK, Confirmation" {

		if err := s.wati.SendTemplate(
			ctx,
			req.WaID,
			"submission_form",
			"Text",
			[]wati.WatiParameters{},
		); err != nil {
			return err
		}

		payload := ChatbotStep{
			CurrentStep: 3,
		}

		if err := s.deleteAndSetRedis(ctx, req.WaID, payload); err != nil {
			// server error
			return err
		}
		return nil
	}

	return nil
}
