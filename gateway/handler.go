package gateway

import (
	"errors"
	errmsg "github.com/likecodingloveproblems/sms-gateway/pkg/err_msg"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ValidateSendMessage func(request SendMessageRequest) error

func SendSMS(service Service, validator ValidateSendMessage) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req SendMessageRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid request format",
			})
		}

		if err := validator(req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}

		message := mapToMessage(req)
		err := service.ProcessMessage(c.Request().Context(), message)
		if errors.Is(err, errmsg.ErrNotEnoughBudget) {
			return c.JSON(http.StatusPaymentRequired, map[string]string{
				"error": err.Error(),
			})
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusAccepted, map[string]interface{}{
			"recipient": req.Recipient,
			"text":      message.Text,
			"status":    "queued",
			"type":      message.Type,
		})
	}
}
