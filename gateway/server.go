package gateway

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type APIServer struct {
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	// SMS
	e.POST("/sms", SendSMS)
	// REPORTING

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

