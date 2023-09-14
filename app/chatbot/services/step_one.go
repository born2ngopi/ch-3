package services

import (
	"context"
	"fmt"
	"github.com/born2ngopi/chatbot-3/app/chatbot/request"
	"github.com/born2ngopi/chatbot-3/app/wati"
	"time"
)

func (s *chatbotService) StepOne(ctx context.Context, chatbotStep ChatbotStep, req request.WebhookRequest, back bool) error {

	if !back && req.Text != "Beli Perdana Tri" {
		return nil
	}

	var templates = []string{
		"generasi_happy_tanggerang_v3",
		"generasi_happy_lampung_v3",
		"-",
	}

	for i := 0; i < 3; i++ {
		if i != 2 {
			if err := s.wati.SendTemplate(
				ctx,
				req.WaID,
				templates[i],
				"Media",
				[]wati.WatiParameters{},
			); err != nil {
				fmt.Println("error dimari", i)
				// server error
				return err
			}
		} else {
			// send list interactive button
			payload := wati.WatiInteractiveListPayload{
				Header: "",
				Footer: "",
				Body: `*Entry Pass Bot*

--------------------
Perdana Tri yang akan kamu dapatkan adalah paket Happy dengan kuota sebesar 5GB
--------------------

Dengan kamu membeli perdana Happy secara otomatis kamu akan mendapatkan entry pass

Keseterdiaan dibawah ini pertanggal ` + "some-date" + `

Silahkan pilih Lokasi
`,
				ButtonText: "Pilih Entry Pass",
				Sections: []wati.WatiInteractiveSections{
					{
						Title: "Pilih Entry Pass",
						Rows: []wati.WatiInteractiveRows{
							{
								Title:       "TANGGERANG",
								Description: "GENERASI HAPPY - TANGGERANG",
							},
							{
								Title:       "LAMPUNG",
								Description: "GENERASI HAPPY - LAMPUNG",
							},
							{
								Title:       "Kembali Ke Menu",
								Description: "Kembali Ke Menu",
							},
							{
								Title:       "Help Desk",
								Description: "Help Desk",
							},
						},
					},
				},
			}

			if err := s.wati.SendInteractiveList(ctx, req.WaID, payload); err != nil {
				return err
			}
		}

		time.Sleep(1500 * time.Millisecond)
	}

	// store to redis
	payload := ChatbotStep{
		CurrentStep: 2,
	}

	if err := s.deleteAndSetRedis(ctx, req.WaID, payload); err != nil {
		// server error
		return err
	}

	return nil
}
