package users

import "time"

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	APIKey    string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}
