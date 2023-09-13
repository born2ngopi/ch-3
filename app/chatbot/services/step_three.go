package services

import (
	"context"
	"fmt"
	"github.com/born2ngopi/chatbot-3/app/chatbot/request"
	"github.com/born2ngopi/chatbot-3/app/wati"
	"strconv"
	"time"
)

func (s *chatbotService) StepThree(ctx context.Context, chatbotStep ChatbotStep, req request.WebhookRequest, back bool) error {

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
		return nil
	} else if text == "Kembali Ke Menu" {
		req.ListReply.Title = chatbotStep.Location

		if err := s.deleteAndSetRedis(ctx, req.WaID, ChatbotStep{
			CurrentStep: 2,
			Location:    chatbotStep.Location,
		}); err != nil {
			return err
		}

		return s.StepTwo(ctx, chatbotStep, req, true)
	} else if isValid := s.buyValidationText.MatchString(text); isValid {
		// get number
		fmt.Println("masuk sini numberByt")
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

		// send syarat dan ketentuan umum
		if err := s.wati.SendTemplate(
			ctx,
			req.WaID,
			"term_condition_1",
			"Text",
			[]wati.WatiParameters{},
		); err != nil {
			// server error
			return err
		}
		time.Sleep(600 * time.Millisecond)

		if err := s.wati.SendTemplate(
			ctx,
			req.WaID,
			"term_condition_2",
			"Text",
			[]wati.WatiParameters{},
		); err != nil {
			// server error
			return err
		}
		time.Sleep(600 * time.Millisecond)

		total := (number * 25000) + 1000
		// send detail order
		if err := s.wati.SendTemplate(
			ctx,
			req.WaID,
			"package_detail_v2",
			"Text",
			[]wati.WatiParameters{
				{
					Name:  "entry_pass_name",
					Value: chatbotStep.Location,
				},
				{
					Name:  "total_items",
					Value: strconv.Itoa(number),
				},
				{
					Name:  "platform_fee",
					Value: s.numberToIDRCurrency(1000),
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
			// server error
			return err
		}

		payload := ChatbotStep{
			CurrentStep: 4,
			Location:    chatbotStep.Location,
		}

		fmt.Println("masuk sini send redis")
		if err := s.deleteAndSetRedis(ctx, req.WaID, payload); err != nil {
			// server error
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
