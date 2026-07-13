package service

import (
	"api-golang/internal/model"
	"context"
	"errors"
	"testing"
)

type mockStore struct {
	getBooksFunc    func(ctx context.Context) ([]*model.Book, error)
	getBookByIDFunc func(ctx context.Context, id int) (*model.Book, error)
	createBookFunc  func(ctx context.Context, book *model.Book) (*model.Book, error)
	updateBookFunc  func(ctx context.Context, id int, book model.Book) (*model.Book, error)
	deleteBookFunc  func(ctx context.Context, id int) error
}

func (m *mockStore) GetBooks(ctx context.Context) ([]*model.Book, error) {
	return m.getBooksFunc(ctx)
}

func (m *mockStore) GetBookByID(ctx context.Context, id int) (*model.Book, error) {
	return m.getBookByIDFunc(ctx, id)
}

func (m *mockStore) CreateBook(ctx context.Context, book *model.Book) (*model.Book, error) {
	return m.createBookFunc(ctx, book)
}

func (m *mockStore) UpdateBook(ctx context.Context, id int, book model.Book) (*model.Book, error) {
	return m.updateBookFunc(ctx, id, book)
}

func (m *mockStore) DeleteBook(ctx context.Context, id int) error {
	return m.deleteBookFunc(ctx, id)
}

func TestService_GetBooks(t *testing.T) {
	t.Run("returns books from store", func(t *testing.T) {
		want := []*model.Book{{ID: 1, Title: "Go", Author: "Rob"}}
		store := &mockStore{
			getBooksFunc: func(ctx context.Context) ([]*model.Book, error) {
				return want, nil
			},
		}
		s := NewService(store)

		got, err := s.GetBooks(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 1 || got[0] != want[0] {
			t.Fatalf("got %v, want %v", got, want)
		}
	})

	t.Run("propagates store error", func(t *testing.T) {
		wantErr := errors.New("db down")
		store := &mockStore{
			getBooksFunc: func(ctx context.Context) ([]*model.Book, error) {
				return nil, wantErr
			},
		}
		s := NewService(store)

		got, err := s.GetBooks(context.Background())
		if !errors.Is(err, wantErr) {
			t.Fatalf("got err %v, want %v", err, wantErr)
		}
		if got != nil {
			t.Fatalf("got %v, want nil", got)
		}
	})
}

func TestService_GetBookByID(t *testing.T) {
	t.Run("returns book from store", func(t *testing.T) {
		want := &model.Book{ID: 1, Title: "Go", Author: "Rob"}
		store := &mockStore{
			getBookByIDFunc: func(ctx context.Context, id int) (*model.Book, error) {
				if id != 1 {
					t.Fatalf("got id %d, want 1", id)
				}
				return want, nil
			},
		}
		s := NewService(store)

		got, err := s.GetBookByID(context.Background(), 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
	})

	t.Run("propagates store error", func(t *testing.T) {
		wantErr := errors.New("not found")
		store := &mockStore{
			getBookByIDFunc: func(ctx context.Context, id int) (*model.Book, error) {
				return nil, wantErr
			},
		}
		s := NewService(store)

		_, err := s.GetBookByID(context.Background(), 99)
		if !errors.Is(err, wantErr) {
			t.Fatalf("got err %v, want %v", err, wantErr)
		}
	})
}

func TestService_CreateBook(t *testing.T) {
	tests := []struct {
		name                string
		book                *model.Book
		wantErr             bool
		wantErrIsValidation bool
	}{
		{
			name:    "valid book",
			book:    &model.Book{Title: "Go", Author: "Rob"},
			wantErr: false,
		},
		{
			name:                "missing title",
			book:                &model.Book{Title: "", Author: "Rob"},
			wantErr:             true,
			wantErrIsValidation: true,
		},
		{
			name:                "missing author",
			book:                &model.Book{Title: "Go", Author: ""},
			wantErr:             true,
			wantErrIsValidation: true,
		},
		{
			name:                "missing both",
			book:                &model.Book{Title: "", Author: ""},
			wantErr:             true,
			wantErrIsValidation: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storeCalled := false
			store := &mockStore{
				createBookFunc: func(ctx context.Context, book *model.Book) (*model.Book, error) {
					storeCalled = true
					book.ID = 1
					return book, nil
				},
			}
			s := NewService(store)

			got, err := s.CreateBook(context.Background(), tt.book)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.wantErrIsValidation && !errors.Is(err, ErrValidation) {
					t.Fatalf("got err %v, want ErrValidation", err)
				}
				if storeCalled {
					t.Fatal("store should not be called for invalid input")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !storeCalled {
				t.Fatal("expected store to be called")
			}
			if got.ID != 1 {
				t.Fatalf("got ID %d, want 1", got.ID)
			}
		})
	}

	t.Run("propagates store error", func(t *testing.T) {
		wantErr := errors.New("insert failed")
		store := &mockStore{
			createBookFunc: func(ctx context.Context, book *model.Book) (*model.Book, error) {
				return nil, wantErr
			},
		}
		s := NewService(store)

		_, err := s.CreateBook(context.Background(), &model.Book{Title: "Go", Author: "Rob"})
		if !errors.Is(err, wantErr) {
			t.Fatalf("got err %v, want %v", err, wantErr)
		}
	})
}

func TestService_UpdateBook(t *testing.T) {
	tests := []struct {
		name                string
		book                model.Book
		wantErr             bool
		wantErrIsValidation bool
	}{
		{
			name:    "valid book",
			book:    model.Book{Title: "Go", Author: "Rob"},
			wantErr: false,
		},
		{
			name:                "missing title",
			book:                model.Book{Title: "", Author: "Rob"},
			wantErr:             true,
			wantErrIsValidation: true,
		},
		{
			name:                "missing author",
			book:                model.Book{Title: "Go", Author: ""},
			wantErr:             true,
			wantErrIsValidation: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storeCalled := false
			store := &mockStore{
				updateBookFunc: func(ctx context.Context, id int, book model.Book) (*model.Book, error) {
					storeCalled = true
					book.ID = id
					return &book, nil
				},
			}
			s := NewService(store)

			got, err := s.UpdateBook(context.Background(), 5, tt.book)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.wantErrIsValidation && !errors.Is(err, ErrValidation) {
					t.Fatalf("got err %v, want ErrValidation", err)
				}
				if storeCalled {
					t.Fatal("store should not be called for invalid input")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !storeCalled {
				t.Fatal("expected store to be called")
			}
			if got.ID != 5 {
				t.Fatalf("got ID %d, want 5", got.ID)
			}
		})
	}

	t.Run("propagates store error", func(t *testing.T) {
		wantErr := errors.New("update failed")
		store := &mockStore{
			updateBookFunc: func(ctx context.Context, id int, book model.Book) (*model.Book, error) {
				return nil, wantErr
			},
		}
		s := NewService(store)

		_, err := s.UpdateBook(context.Background(), 5, model.Book{Title: "Go", Author: "Rob"})
		if !errors.Is(err, wantErr) {
			t.Fatalf("got err %v, want %v", err, wantErr)
		}
	})
}

func TestService_DeleteBook(t *testing.T) {
	t.Run("deletes book via store", func(t *testing.T) {
		var gotID int
		store := &mockStore{
			deleteBookFunc: func(ctx context.Context, id int) error {
				gotID = id
				return nil
			},
		}
		s := NewService(store)

		if err := s.DeleteBook(context.Background(), 7); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if gotID != 7 {
			t.Fatalf("got id %d, want 7", gotID)
		}
	})

	t.Run("propagates store error", func(t *testing.T) {
		wantErr := errors.New("delete failed")
		store := &mockStore{
			deleteBookFunc: func(ctx context.Context, id int) error {
				return wantErr
			},
		}
		s := NewService(store)

		err := s.DeleteBook(context.Background(), 7)
		if !errors.Is(err, wantErr) {
			t.Fatalf("got err %v, want %v", err, wantErr)
		}
	})
}
