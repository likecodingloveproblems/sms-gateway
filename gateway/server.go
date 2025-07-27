package gateway

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	myredis "github.com/likecodingloveproblems/sms-gateway/pkg/redis"
	_ "net/http"
)

type APIServer struct{}

func NewAPIServer() *APIServer {
	return &APIServer{}
}

func (APIServer) Serve() {
	client := myredis.NewClient(0)
	repository := NewRepository(client)
	gateway := NewService(repository)

	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	// SMS
	e.POST("/sms", SendSMS(gateway, ValidateSendMessageRequest))
	// REPORTING

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
