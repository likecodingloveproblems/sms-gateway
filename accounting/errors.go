package accounting

import "errors"

var (
	ErrUserNotFound   = errors.New("user_not_found")
	ErrInternalServer = errors.New("internal_server_error")
)
