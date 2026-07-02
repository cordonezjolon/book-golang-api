package transport

import (
	"api-golang/internal/model"
	"api-golang/internal/service"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type BookHandler struct {
	service service.BookService
}

func NewBookHandler(service service.BookService) *BookHandler {
	return &BookHandler{service: service}
}

func (h *BookHandler) HandleBooks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	switch r.Method {
	case http.MethodGet:
		books, err := h.service.GetBooks(ctx)
		if err != nil {
			http.Error(w, "Failed to get books", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(books); err != nil {
			http.Error(w, "Failed to encode books", http.StatusInternalServerError)
			return
		}

	case http.MethodPost:
		var book model.Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		createdBook, err := h.service.CreateBook(ctx, &book)
		if err != nil {
			if errors.Is(err, service.ErrValidation) {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				log.Printf("create book: %v", err)
				http.Error(w, "Failed to create book", http.StatusInternalServerError)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(createdBook); err != nil {
			http.Error(w, "Failed to encode created book", http.StatusInternalServerError)
			return
		}

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *BookHandler) HandleBookByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}
	if id == 0 {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		book, err := h.service.GetBookByID(ctx, id)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Book not found", http.StatusNotFound)
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(book); err != nil {
			http.Error(w, "Failed to encode book", http.StatusInternalServerError)
			return
		}

	case http.MethodPut:
		var book model.Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		updatedBook, err := h.service.UpdateBook(ctx, id, book)
		if err != nil {
			if errors.Is(err, service.ErrValidation) {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				log.Printf("create book: %v", err)
				http.Error(w, "Failed to create book", http.StatusInternalServerError)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(updatedBook); err != nil {
			http.Error(w, "Failed to encode updated book", http.StatusInternalServerError)
			return
		}

	case http.MethodDelete:
		if err := h.service.DeleteBook(ctx, id); err != nil {
			http.Error(w, "Failed to delete book", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
