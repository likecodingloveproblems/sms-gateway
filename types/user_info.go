package types

type UserInfo struct {
	ID   uint64 `json:"id"`
	Role Role   `json:"role"`
}

type Role uint8

const (
	RoleAdmin    Role = iota + 1 // admin = 1
	RoleOwner                    // restaurant owner = 2
	RoleCustomer                 // customer = 3
)
