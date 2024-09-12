package domain

import "time"

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	UserName  string    `json:"userName"`
	Role      Role      `json:"role"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
type UserFilter struct {
	ID       string
	UserName string
}

// UserService defines the methods for managing users
type UserService interface {
	CreateUser(user *User) error
	GetUser(userFilter UserFilter) (*User, error)
}
