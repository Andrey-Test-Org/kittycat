package book

import (
	"context"
	"errors"
	"testing"
	"time"
)

func newTestService(t *testing.T) *Service {
	t.Helper()
	return NewService(NewInMemoryRepository())
}

func sampleInput() CreateInput {
	return CreateInput{
		ISBN:        "978-0-13-468599-1",
		Title:       "The Go Programming Language",
		AuthorID:    "auth_abc",
		PriceCents:  3999,
		Currency:    "usd",
		PublishedAt: time.Date(2015, 11, 5, 0, 0, 0, 0, time.UTC),
		Genre:       "Programming",
		PageCount:   380,
		Description: "An authoritative guide.",
	}
}

func TestService_Create(t *testing.T) {
	t.Parallel()
	svc := newTestService(t)
	in := sampleInput()
	b, err := svc.Create(context.Background(), in)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if b.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if b.Currency != "USD" {
		t.Fatalf("expected USD, got %s", b.Currency)
	}
	if b.ISBN != "9780134685991" {
		t.Fatalf("isbn not normalized: %s", b.ISBN)
	}
}

func TestService_Create_Validation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		mutate  func(*CreateInput)
		wantErr error
	}{
		{
			name:    "blank title",
			mutate:  func(in *CreateInput) { in.Title = "  " },
			wantErr: ErrInvalidTitle,
		},
		{
			name:    "negative price",
			mutate:  func(in *CreateInput) { in.PriceCents = -1 },
			wantErr: ErrInvalidPrice,
		},
		{
			name:    "unsupported currency",
			mutate:  func(in *CreateInput) { in.Currency = "ZZZ" },
			wantErr: ErrInvalidCurrency,
		},
		{
			name:    "zero pages",
			mutate:  func(in *CreateInput) { in.PageCount = 0 },
			wantErr: ErrInvalidPageCount,
		},
		{
			name:    "bad isbn",
			mutate:  func(in *CreateInput) { in.ISBN = "abc" },
			wantErr: ErrInvalidISBN,
		},
		{
			name:    "blank genre",
			mutate:  func(in *CreateInput) { in.Genre = "" },
			wantErr: ErrInvalidGenre,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			svc := newTestService(t)
			in := sampleInput()
			tc.mutate(&in)
			_, err := svc.Create(context.Background(), in)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("want %v, got %v", tc.wantErr, err)
			}
		})
	}
}

func TestService_GetUpdateDelete(t *testing.T) {
	t.Parallel()
	svc := newTestService(t)
	created, err := svc.Create(context.Background(), sampleInput())
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	got, err := svc.Get(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.ID != created.ID {
		t.Fatalf("expected %s, got %s", created.ID, got.ID)
	}

	newPrice := int64(4999)
	updated, err := svc.Update(context.Background(), created.ID, UpdateInput{PriceCents: &newPrice})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.PriceCents != 4999 {
		t.Fatalf("expected 4999, got %d", updated.PriceCents)
	}

	if err := svc.Delete(context.Background(), created.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := svc.Get(context.Background(), created.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestService_ListAndSearch(t *testing.T) {
	t.Parallel()
	svc := newTestService(t)
	titles := []string{"Go Cookbook", "Learning Go", "Effective Go", "The Pragmatic Programmer"}
	for _, title := range titles {
		in := sampleInput()
		in.Title = title
		if _, err := svc.Create(context.Background(), in); err != nil {
			t.Fatalf("Create %s: %v", title, err)
		}
	}

	all, err := svc.List(context.Background(), 0, 10)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(all) != 4 {
		t.Fatalf("expected 4 books, got %d", len(all))
	}

	hits, err := svc.Search(context.Background(), "Go", 10)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(hits) != 3 {
		t.Fatalf("expected 3 Go books, got %d", len(hits))
	}
}

func TestService_Search_InvalidQuery(t *testing.T) {
	t.Parallel()
	svc := newTestService(t)
	_, err := svc.Search(context.Background(), "   ", 10)
	if !errors.Is(err, ErrInvalidQuery) {
		t.Fatalf("expected ErrInvalidQuery, got %v", err)
	}
}

func TestService_CountForAuthor(t *testing.T) {
	t.Parallel()
	svc := newTestService(t)
	in := sampleInput()
	in.AuthorID = "auth_x"
	if _, err := svc.Create(context.Background(), in); err != nil {
		t.Fatalf("Create: %v", err)
	}
	in.Title = "Second Title"
	if _, err := svc.Create(context.Background(), in); err != nil {
		t.Fatalf("Create: %v", err)
	}
	in.AuthorID = "auth_y"
	in.Title = "Third Title"
	if _, err := svc.Create(context.Background(), in); err != nil {
		t.Fatalf("Create: %v", err)
	}

	n, err := svc.CountForAuthor(context.Background(), "auth_x")
	if err != nil {
		t.Fatalf("CountForAuthor: %v", err)
	}
	if n != 2 {
		t.Fatalf("expected 2, got %d", n)
	}
}
