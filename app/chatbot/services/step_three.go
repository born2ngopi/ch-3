package services

import (
	"context"
	"github.com/born2ngopi/chatbot-3/app/chatbot/request"
	"github.com/born2ngopi/chatbot-3/app/wati"
	"strings"
	"time"
)

func (s *chatbotService) StepThree(ctx context.Context, chatbotStep ChatbotStep, req request.WebhookRequest, back bool) error {

	text := req.Text
	text = strings.ToLower(text)

	// submission_form_invalid

	//check message valid
	var isValid = true
	if !strings.Contains(text, "nama") {
		isValid = false
	}
	// alamat pengiriman
	if !strings.Contains(text, "alamat") && !strings.Contains(text, "alamat pengiriman") {
		isValid = false
	}
	// no hp
	if !strings.Contains(text, "no") && !strings.Contains(text, "no hp") && !strings.Contains(text, "no.hp:") && !strings.Contains(text, "no hp") {
		isValid = false
	}

	if !isValid {
		if err := s.wati.SendTemplate(
			ctx,
			req.WaID,
			"submission_form_invalid",
			"",
			[]wati.WatiParameters{},
		); err != nil {
			return err
		}
		return nil
	}

	now := time.Now()

	futureTime := now.AddDate(0, 0, 1)
	wib := time.FixedZone("WIB", 7*60*60)
	futureTime = futureTime.In(wib)

	if err := s.wati.SendTemplate(
		ctx,
		req.WaID,
		"payment_message_v3",
		"Media",
		[]wati.WatiParameters{
			{
				Name:  "end_date",
				Value: futureTime.Format("2006-01-02 15:04:05"),
			},
		},
	); err != nil {
		return err
	}

	payload := ChatbotStep{
		CurrentStep: 4,
	}

	if err := s.deleteAndSetRedis(ctx, req.WaID, payload); err != nil {
		// server error
		return err
	}
	return nil

}
