package model

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

func ParseRole(role string) (Role, bool) {
	switch Role(role) {
	case RoleAdmin:
		return RoleAdmin, true
	case RoleUser:
		return RoleUser, true
	default:
		return "", false
	}
}
