// Package book contains the Book domain: model, repository, service, validators.
package book

import "time"

// Book is a stocked title in the bookstore catalog.
type Book struct {
	ID          string    `json:"id"`
	ISBN        string    `json:"isbn"`
	Title       string    `json:"title"`
	Subtitle    string    `json:"subtitle,omitempty"`
	AuthorID    string    `json:"authorId"`
	PriceCents  int64     `json:"priceCents"`
	Currency    string    `json:"currency"`
	PublishedAt time.Time `json:"publishedAt"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Genre       string    `json:"genre"`
	PageCount   int       `json:"pageCount"`
	Description string    `json:"description,omitempty"`
}

// Summary is a slim projection for list views.
type Summary struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	AuthorID   string `json:"authorId"`
	PriceCents int64  `json:"priceCents"`
	Currency   string `json:"currency"`
	Genre      string `json:"genre"`
}

// ToSummary returns a Summary projection of the receiver.
func (b Book) ToSummary() Summary {
	return Summary{
		ID:         b.ID,
		Title:      b.Title,
		AuthorID:   b.AuthorID,
		PriceCents: b.PriceCents,
		Currency:   b.Currency,
		Genre:      b.Genre,
	}
}
