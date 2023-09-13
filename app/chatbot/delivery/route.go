package delivery

import (
	"github.com/born2ngopi/chatbot-3/app/chatbot/delivery/handler"
	"github.com/labstack/echo/v4"
)

func SetupRoute(route *echo.Group, h handler.ChatbotHandler) {
	route.POST("/webhook", h.Webhook)
}
