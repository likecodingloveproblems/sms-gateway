package generalmappers

import (
	"fmt"

	"git.gocasts.ir/remenu/beehive/types"
)

func MapStringToRole(roleString string) (types.Role, error) {
	var role types.Role

	switch roleString {
	case "admin":
		role = types.RoleAdmin
	case "owner":
		role = types.RoleOwner
	case "customer":
		role = types.RoleCustomer
	default:
		return types.RoleOwner, fmt.Errorf("unknown role: %s", roleString)
	}

	return role, nil
}
