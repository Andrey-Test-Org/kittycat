package author

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

func sampleAuthor() CreateInput {
	return CreateInput{
		FullName:  "Alice Authoress",
		Country:   "Imagineland",
		Birthdate: time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC),
		Bio:       "Loves cats and footnotes.",
	}
}

func TestService_Create(t *testing.T) {
	t.Parallel()
	svc := newTestService(t)
	a, err := svc.Create(context.Background(), sampleAuthor())
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if a.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if a.Slug() != "alice-authoress" {
		t.Fatalf("expected alice-authoress, got %s", a.Slug())
	}
}

func TestService_Create_Validation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		mutate  func(*CreateInput)
		wantErr error
	}{
		{name: "blank name", mutate: func(in *CreateInput) { in.FullName = "" }, wantErr: ErrInvalidName},
		{name: "blank country", mutate: func(in *CreateInput) { in.Country = "" }, wantErr: ErrInvalidCountry},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			svc := newTestService(t)
			in := sampleAuthor()
			tc.mutate(&in)
			_, err := svc.Create(context.Background(), in)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("want %v, got %v", tc.wantErr, err)
			}
		})
	}
}

func TestService_GetUpdate(t *testing.T) {
	t.Parallel()
	svc := newTestService(t)
	created, err := svc.Create(context.Background(), sampleAuthor())
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

	newCountry := "Otherland"
	upd, err := svc.Update(context.Background(), created.ID, UpdateInput{Country: &newCountry})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if upd.Country != "Otherland" {
		t.Fatalf("expected Otherland, got %s", upd.Country)
	}
}

func TestService_List(t *testing.T) {
	t.Parallel()
	svc := newTestService(t)
	for i := 0; i < 5; i++ {
		if _, err := svc.Create(context.Background(), sampleAuthor()); err != nil {
			t.Fatalf("Create: %v", err)
		}
	}

	all, err := svc.List(context.Background(), 0, 10)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(all) != 5 {
		t.Fatalf("expected 5, got %d", len(all))
	}
}
