package gateway

import (
	"github.com/go-ozzo/ozzo-validation/v4"
)

func ValidateSendMessageRequest(req SendMessageRequest) error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.UserId, validation.Required),
		validation.Field(&req.Recipient, validation.Required),
		validation.Field(&req.Text, validation.Required),
		validation.Field(&req.Type,
			validation.In("", "normal", "express").Error("must be empty, 'normal' or 'express'")),
	)
}
