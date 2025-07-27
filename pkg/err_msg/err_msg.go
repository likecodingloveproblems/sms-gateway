package errmsg

import "errors"

var (
	ErrValidationFailed     = errors.New("input validation failed")
	ErrUnexpectedError      = errors.New("unexpected error occurred")
	ErrInvalidRequestFormat = errors.New("invalid request format")
	ErrFailedDecodeBase64   = errors.New("decode data to base 64 failed")
	ErrFailedUnmarshalJson  = errors.New("unmarshal data to JSON failed")
	ErrNotEnoughBudget      = errors.New("not enough budget")
	ErrUnknownMessageType   = errors.New("unknown message type")
	ErrUserNotFound         = errors.New("user_not_found")
)

// Define constant messages generally
const (
	MessageMissingXUserData  = "Missing X-User-Data header"
	MessageInvalidBase64     = "Invalid Base64 data"
	MessageInvalidJsonFormat = "Invalid JSON format"
	ServerError              = "Internal server error"
)
