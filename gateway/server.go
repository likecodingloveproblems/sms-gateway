package gateway

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type APIServer struct {
}

type SMSRequest struct {
	Phone   string `json:"phone" validate:"required"`
	Message string `json:"message" validate:"required,min=1,max=160"`
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	// SMS
	e.POST("/sms", sendSMS)
	// REPORTING

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

func sendSMS(c echo.Context) error {
	return processSMS(c, "normal")
}

func processSMS(c echo.Context, smsType string) error {
	var req SMSRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{
			"error": err.Error(),
		})
	}

	// TODO: Add actual SMS sending logic with queueing
	return c.JSON(http.StatusAccepted, map[string]interface{}{
		"status":  "queued",
		"type":    smsType,
		"phone":   req.Phone,
		"message": req.Message,
		"warning": "SMS not actually sent - demo mode",
	})
}
