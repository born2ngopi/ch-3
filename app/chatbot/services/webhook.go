package services

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/born2ngopi/chatbot-3/app/chatbot/request"
	"github.com/born2ngopi/chatbot-3/app/wati"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func (s *chatbotService) Webhook(ctx context.Context, req request.WebhookRequest) error {

	value, err := s.redis.Get(ctx, req.WaID).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return err
		}
		// first message
		if !strings.Contains(strings.ToLower(req.Text), "beli") {
			return nil
		}

		// send image, and list price list
		if err := s.wati.SendTemplate(
			ctx,
			req.WaID,
			"generasi_happy_lampung_v3",
			"Media",
			[]wati.WatiParameters{},
		); err != nil {

			return err
		}
		time.Sleep(1500 * time.Millisecond)

		interactivePayload := wati.WatiInteractiveListPayload{
			Header: "",
			Footer: "",
			Body: `Selamat kamu sudah terhubung dengan GENERASI HAPPY FESTIVAL LAMPUNG.

Ini adalah layanan otomatis.

Untuk melanjutkan, klik tombol beli dibawah ini.`,
			ButtonText: "Pilih Entry Pass",
			Sections: []wati.WatiInteractiveSections{
				{
					Title: "Pilih Entry Pass",
					Rows: []wati.WatiInteractiveRows{
						{
							Title:       "Beli 10 Perdana",
							Description: "@25.000 x 10 Entry Pass",
						},
						{
							Title:       "Beli 6 Perdana",
							Description: "@25.000 x 6 Entry Pass",
						},
						{
							Title:       "Beli 4 Perdana",
							Description: "@25.000 x 4 Entry Pass",
						},
						{
							Title:       "Beli 3 Perdana",
							Description: "@25.000 x 3 Entry Pass",
						},
						{
							Title:       "Beli 2 Perdana",
							Description: "@25.000 x 2 Entry Pass",
						},
						{
							Title:       "Beli 1 Perdana",
							Description: "@25.000 x 1 Entry Pass",
						},
					},
				},
			},
		}

		if err := s.wati.SendInteractiveList(ctx, req.WaID, interactivePayload); err != nil {
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
	}
	//
	return nil
}

func (s *chatbotService) InvalidMessage(ctx context.Context, req request.WebhookRequest) error {
	// send invalid message
	if err := s.wati.SendTemplate(
		ctx,
		req.WaID,
		"invalid_message",
		"",
		[]wati.WatiParameters{},
	); err != nil {
		return err
	}
	return nil
}
