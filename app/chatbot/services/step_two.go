package services

import (
	"context"
	"github.com/born2ngopi/chatbot-3/app/chatbot/request"
	"github.com/born2ngopi/chatbot-3/app/wati"
	"time"
)

func (s *chatbotService) StepTwo(ctx context.Context, chatbotStep ChatbotStep, req request.WebhookRequest, back bool) error {

	text := req.ListReply.Title

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
	} else if text == "Kembali Ke Menu" {
		// back to menu
		if err := s.wati.SendTemplate(
			ctx,
			req.WaID,
			"first_message",
			"first_message",
			[]wati.WatiParameters{},
		); err != nil {
			// server error
			return err
		}

		payload := ChatbotStep{
			CurrentStep: 1,
		}

		if err := s.deleteAndSetRedis(ctx, req.WaID, payload); err != nil {
			// server error
			return err
		}
	} else if text == "TANGGERANG" || text == "LAMPUNG" || back {
		// back to menu
		var entryPass, template string
		if text == "TANGGERANG" {
			entryPass = "GENERASI HAPPY - TANGGERANG"
			template = "generasi_happy_tanggerang"
		} else if back {
			entryPass = chatbotStep.Location
			if entryPass == "GENERASI HAPPY - TANGGERANG" {
				template = "generasi_happy_tanggerang"
			} else {
				template = "generasi_happy_lampung"
			}
		} else {
			entryPass = "GENERASI HAPPY - LAMPUNG"
			template = "generasi_happy_lampung"
		}

		// send image
		if err := s.wati.SendTemplate(
			ctx,
			req.WaID,
			template,
			"Media",
			[]wati.WatiParameters{
				{
					Name:  "entry_pass",
					Value: entryPass,
				},
			},
		); err != nil {
			// server error
			return err
		}
		time.Sleep(600 * time.Millisecond)
		// send text interactive list
		payload := wati.WatiInteractiveListPayload{
			Header: "",
			Footer: "",
			Body: `Entry Pass:
` + entryPass + `

Note:
- Perdana yang sudah dibeli tidak dapat dipindah tangan atau dijual lagi.
- Entry Pass akan dikirimkan soft copy melalui whatsapp.

Perhatian: Perdana yang akan anda dapatkan setelah melakukan pembayaran, merupakan bukti pembelian atas sejumlah Entry Pass.`,
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

		if err := s.wati.SendInteractiveList(ctx, req.WaID, payload); err != nil {
			return err
		}

		if err := s.deleteAndSetRedis(ctx, req.WaID, ChatbotStep{
			CurrentStep: 3,
			Location:    text,
		}); err != nil {
			return err
		}
	} else {
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
	}

	return nil
}
