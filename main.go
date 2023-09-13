package main

import (
	"github.com/born2ngopi/chatbot-3/app/chatbot/delivery"
	"github.com/born2ngopi/chatbot-3/app/chatbot/delivery/handler"
	"github.com/born2ngopi/chatbot-3/app/chatbot/services"
	"github.com/born2ngopi/chatbot-3/app/wati"
	"github.com/born2ngopi/chatbot-3/config"
	"github.com/labstack/echo/v4"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {

	conf := config.Init()
	//db := conf.ConnectDatabase()
	_redis := conf.ConnectRedis()
	_wati := wati.NewClient()

	e := echo.New()

	//service
	chatbotService := services.NewChatbotService(_redis, _wati)

	//handler
	chatbotHandler := handler.NewChatbotHandler(chatbotService)

	e.GET("/version", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"version": "1.0.0",
		})
	})

	g := e.Group("/api/v1")
	delivery.SetupRoute(g, chatbotHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
