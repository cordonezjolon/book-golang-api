package store

import (
	"api-golang/internal/model"
	"context"

	"database/sql"
)

type Store interface {
	GetBooks(ctx context.Context) ([]*model.Book, error)
	GetBookByID(ctx context.Context, id int) (*model.Book, error)
	CreateBook(ctx context.Context, book *model.Book) (*model.Book, error)
	UpdateBook(ctx context.Context, id int, book model.Book) (*model.Book, error)
	DeleteBook(ctx context.Context, id int) error
}

type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &store{db: db}
}

// CreateBook implements [Store].
func (s *store) CreateBook(ctx context.Context, book *model.Book) (*model.Book, error) {

	q := `INSERT INTO books (title, author) VALUES ($1, $2) RETURNING id`

	row := s.db.QueryRowContext(ctx, q, book.Title, book.Author)

	if err := row.Scan(&book.ID); err != nil {

		return nil, err
	}

	return book, nil
}

// DeleteBook implements [Store].
func (s *store) DeleteBook(ctx context.Context, id int) error {
	q := `DELETE FROM books WHERE id = $1`
	result, err := s.db.ExecContext(ctx, q, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil

}

// GetBookByID implements [Store].
func (s *store) GetBookByID(ctx context.Context, id int) (*model.Book, error) {
	q := `SELECT id, title, author FROM books WHERE id = $1`
	book := model.Book{}

	err := s.db.QueryRowContext(ctx, q, id).Scan(&book.ID, &book.Title, &book.Author)

	if err != nil {
		return nil, err
	}

	return &book, nil

}

// GetBooks implements [Store].
func (s *store) GetBooks(ctx context.Context) ([]*model.Book, error) {
	q := `SELECT id, title, author FROM books`
	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := make([]*model.Book, 0)
	for rows.Next() {
		book := model.Book{}
		if err := rows.Scan(&book.ID, &book.Title, &book.Author); err != nil {
			return nil, err
		}
		books = append(books, &book)

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return books, nil
}

// UpdateBook implements [Store].
func (s *store) UpdateBook(ctx context.Context, id int, book model.Book) (*model.Book, error) {

	q := `UPDATE books SET title = $1, author = $2 WHERE id = $3`
	result, err := s.db.ExecContext(ctx, q, book.Title, book.Author, id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	book.ID = id
	return &book, nil

}
