package services

import (
	"context"
	"fmt"
	"github.com/born2ngopi/chatbot-3/app/chatbot/request"
	"github.com/born2ngopi/chatbot-3/app/wati"
	"time"
)

func (s *chatbotService) StepFour(ctx context.Context, chatbotStep ChatbotStep, req request.WebhookRequest, back bool) error {

	text := req.Text

	if text == "OK, Confirmation" {
		fmt.Println("masuk confirmation")

		now := time.Now()

		futureTime := now.AddDate(0, 0, 1)
		wib := time.FixedZone("WIB", 7*60*60)
		futureTime = futureTime.In(wib)

		if err := s.wati.SendTemplate(
			ctx,
			req.WaID,
			"payment_message_v2",
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
			CurrentStep: 5,
		}

		if err := s.deleteAndSetRedis(ctx, req.WaID, payload); err != nil {
			// server error
			return err
		}
		return nil
	}

	// help desk
	if text == "Help Desk" {
		if err := s.wati.SendTemplate(
			ctx,
			req.WaID,
			"help_desk",
			"",
			[]wati.WatiParameters{},
		); err != nil {
			// server error
			return err
		}

		// delete redis
		if err := s.redis.Del(ctx, req.WaID).Err(); err != nil {
			return err
		}
		return nil
	} else if text == "Kembali Ke Menu" || text == "Back" {
		req.ListReply.Title = chatbotStep.Location

		if err := s.deleteAndSetRedis(ctx, req.WaID, ChatbotStep{
			CurrentStep: 2,
			Location:    chatbotStep.Location,
		}); err != nil {
			return err
		}

		return s.StepTwo(ctx, chatbotStep, req, true)
	} else {
		// send invalid message
		if err := s.wati.SendTemplate(
			ctx,
			req.WaID,
			"invalid_message",
			"Text",
			[]wati.WatiParameters{},
		); err != nil {
			return err
		}
	}

	return nil
}
