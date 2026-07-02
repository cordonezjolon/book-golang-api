package service

import (
	"api-golang/internal/model"
	"api-golang/internal/store"
	"context"
	"errors"
	"fmt"
)

var ErrValidation = errors.New("validation error")

type BookService interface {
	GetBooks(ctx context.Context) ([]*model.Book, error)
	GetBookByID(ctx context.Context, id int) (*model.Book, error)
	CreateBook(ctx context.Context, book *model.Book) (*model.Book, error)
	UpdateBook(ctx context.Context, id int, book model.Book) (*model.Book, error)
	DeleteBook(ctx context.Context, id int) error
}

type Service struct {
	store store.Store
}

func NewService(store store.Store) *Service {
	return &Service{
		store: store}
}

func (s *Service) GetBooks(ctx context.Context) ([]*model.Book, error) {

	books, err := s.store.GetBooks(ctx)
	if err != nil {
		return nil, err

	}
	return books, nil

}

func (s *Service) GetBookByID(ctx context.Context, id int) (*model.Book, error) {
	return s.store.GetBookByID(ctx, id)
}

func (s *Service) CreateBook(ctx context.Context, book *model.Book) (*model.Book, error) {
	if book.Title == "" || book.Author == "" {
		return nil, fmt.Errorf("title and author are required: %w", ErrValidation)

	}
	return s.store.CreateBook(ctx, book)
}

func (s *Service) UpdateBook(ctx context.Context, id int, book model.Book) (*model.Book, error) {
	if book.Title == "" || book.Author == "" {
		return nil, fmt.Errorf("title and author are required: %w", ErrValidation)
	}

	return s.store.UpdateBook(ctx, id, book)
}

func (s *Service) DeleteBook(ctx context.Context, id int) error {
	return s.store.DeleteBook(ctx, id)
}
