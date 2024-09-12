package domain

// AuthService defines the methods for authentication
type AuthService interface {
	GenerateToken(user *User) (string, error)
	ValidateToken(token string) (*User, error)
}
