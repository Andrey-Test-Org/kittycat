// Package author contains the Author domain.
package author

import "time"

// Author represents a writer in the bookstore catalog.
type Author struct {
	ID        string    `json:"id"`
	FullName  string    `json:"fullName"`
	Country   string    `json:"country"`
	Birthdate time.Time `json:"birthdate"`
	Bio       string    `json:"bio,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Slug returns a url-safe slug of the author's full name.
func (a Author) Slug() string {
	out := make([]byte, 0, len(a.FullName))
	for i := 0; i < len(a.FullName); i++ {
		c := a.FullName[i]
		switch {
		case c >= 'a' && c <= 'z':
			out = append(out, c)
		case c >= 'A' && c <= 'Z':
			out = append(out, c+32)
		case c >= '0' && c <= '9':
			out = append(out, c)
		case c == ' ' || c == '-' || c == '_':
			if len(out) > 0 && out[len(out)-1] != '-' {
				out = append(out, '-')
			}
		}
	}
	for len(out) > 0 && out[len(out)-1] == '-' {
		out = out[:len(out)-1]
	}
	return string(out)
}
