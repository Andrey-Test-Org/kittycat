package users

import "time"

// User is the public representation of a registered account.
// APIKey is never emitted in JSON responses (json:"-").
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	APIKey    string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}
