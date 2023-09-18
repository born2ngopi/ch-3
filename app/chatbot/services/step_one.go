package services

import (
	"context"
	"github.com/born2ngopi/chatbot-3/app/chatbot/request"
	"github.com/born2ngopi/chatbot-3/app/wati"
	"strconv"
)

func (s *chatbotService) StepOne(ctx context.Context, chatbotStep ChatbotStep, req request.WebhookRequest, back bool) error {

	text := req.ListReply.Title

	if isValid := s.buyValidationText.MatchString(text); isValid {

		// get number
		numberBuy := s.findNumberBuy.FindStringSubmatch(text)
		var number int
		if len(numberBuy) > 1 {
			i, err := strconv.Atoi(numberBuy[1])
			if err != nil {
				return err
			}
			number = i
		}

		isValid = func() bool {
			if number == 1 || number == 2 || number == 3 || number == 4 || number == 6 || number == 10 {
				return true
			}
			return false
		}()

		if !isValid {
			return s.InvalidMessage(ctx, req)
		}

		platformFee := 1_000
		total := (number * 25_000) + platformFee

		if err := s.wati.SendTemplate(
			ctx,
			req.WaID,
			"confirmation_message",
			"Text",
			[]wati.WatiParameters{
				{
					Name:  "total_order",
					Value: strconv.Itoa(number),
				},
				{
					Name:  "platform_fee",
					Value: s.numberToIDRCurrency(platformFee),
				},
				{
					Name:  "ongkir_fee",
					Value: "-",
				},
				{
					Name:  "order_amount",
					Value: s.numberToIDRCurrency(total),
				},
			},
		); err != nil {
			return err
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

	return s.InvalidMessage(ctx, req)
}
