package handler

import (
	"fmt"
	"github.com/born2ngopi/chatbot-3/app/chatbot/request"
	"github.com/born2ngopi/chatbot-3/app/chatbot/services"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type ChatbotHandler interface {
	Webhook(c echo.Context) error
}

type chatbotHandler struct {
	service services.ChatbotService
}

func NewChatbotHandler(service services.ChatbotService) ChatbotHandler {
	return &chatbotHandler{
		service: service,
	}
}

type Response struct {
	Succsess bool
	Message  string
}

func (h *chatbotHandler) Webhook(c echo.Context) error {

	var payload request.WebhookRequest
	if err := c.Bind(&payload); err != nil {
		fmt.Println("bind error", err.Error())
		return c.JSON(http.StatusBadRequest, Response{
			Succsess: false,
			Message:  fmt.Sprintf("%v", err),
		})
	}

	if err := h.service.Webhook(c.Request().Context(), payload); err != nil {
		log.Println(err)
		fmt.Println("response error")
		return c.JSON(500, Response{
			Succsess: false,
			Message:  fmt.Sprintf("Internal Server Error: %v", err),
		})
	}

	return c.JSON(200, Response{
		Succsess: true,
		Message:  "Success",
	})
}
