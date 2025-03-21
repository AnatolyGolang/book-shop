package repositories

import (
	rm "book-shop/internal/app/repositories/models"
	sm "book-shop/internal/app/services/models"
	"book-shop/internal/pkg/postgres"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type BookRepositoryImpl struct {
	db *postgres.DBConnection
}

func NewBookRepository(db *postgres.DBConnection) *BookRepositoryImpl {
	return &BookRepositoryImpl{db: db}
}

func (r *BookRepositoryImpl) CreateBook(ctx context.Context, book sm.DomainBook) (rm.Book, error) {
	query := `INSERT INTO books (title, author, category_id, price, amount, year) 
              VALUES ($1, $2, $3, $4, 10, $5) 
              RETURNING id, title, author, category_id, price, amount, year, created_at, updated_at`

	var newBook rm.Book
	err := r.db.QueryRow(ctx, query, book.Title, book.Author, book.CategoryID, book.Price, book.Year).Scan(
		&newBook.ID, &newBook.Title, &newBook.Author, &newBook.CategoryID,
		&newBook.Price, &newBook.Amount, &newBook.Year, &newBook.CreatedAt, &newBook.UpdatedAt,
	)
	if err != nil {
		return rm.Book{}, fmt.Errorf("failed to create book: %w", err)
	}

	return newBook, nil
}

func (r *BookRepositoryImpl) GetBook(ctx context.Context, id int) (rm.Book, error) {
	if id == 0 {
		return rm.Book{}, fmt.Errorf("id can not be 0")
	}

	var book rm.Book
	query := `SELECT id, title, author, category_id, price, amount, year, created_at, updated_at 
              	FROM books 
              WHERE id = $1`

	row := r.db.QueryRow(ctx, query, id)
	err := row.Scan(
		&book.ID, &book.Title, &book.Author, &book.CategoryID,
		&book.Price, &book.Amount, &book.Year, &book.CreatedAt, &book.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return rm.Book{}, fmt.Errorf("book not found")
		}
		return rm.Book{}, fmt.Errorf("failed to get a book: %w", err)
	}

	return book, nil
}

func (r *BookRepositoryImpl) UpdateBook(ctx context.Context, id int, book sm.DomainBook) (rm.Book, error) {
	query := `UPDATE books SET title = $1, author = $2, category_id = $3, price = $4, amount = $5, year = $6, updated_at = NOW()
              WHERE id = $7 
              RETURNING id, title, author, category_id, price, amount, year, created_at, updated_at`

	var updatedBook rm.Book
	err := r.db.QueryRow(ctx, query, book.Title, book.Author, book.CategoryID, book.Price, book.Amount, book.Year, id).Scan(
		&updatedBook.ID, &updatedBook.Title, &updatedBook.Author, &updatedBook.CategoryID,
		&updatedBook.Price, &updatedBook.Amount, &updatedBook.Year, &updatedBook.CreatedAt, &updatedBook.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return rm.Book{}, fmt.Errorf("book not found")
		}
		return rm.Book{}, fmt.Errorf("failed to update book: %w", err)
	}

	return updatedBook, nil
}

func (r *BookRepositoryImpl) DeleteBook(ctx context.Context, id int) error {
	if id == 0 {
		return fmt.Errorf("id can not be 0")
	}

	query := `DELETE FROM books 
              	WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("book not found")
	}

	return nil
}

func (r *BookRepositoryImpl) GetBooksByCategories(ctx context.Context, categoryIDs []int, limit int, offset int) ([]rm.Book, int, error) {
	if len(categoryIDs) == 0 {
		return nil, 0, fmt.Errorf("categoryIDs cannot be empty")
	}

	query := fmt.Sprintf(
		`SELECT id, title, author, category_id, price, amount, year, created_at, updated_at 
        FROM books 
        WHERE category_id = ANY($1) 
        ORDER BY id 
        LIMIT $2 OFFSET $3`,
	)

	rows, err := r.db.Query(ctx, query, categoryIDs, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get books: %w", err)
	}
	defer rows.Close()

	var books []rm.Book
	for rows.Next() {
		var book rm.Book
		err := rows.Scan(
			&book.ID, &book.Title, &book.Author, &book.CategoryID,
			&book.Price, &book.Amount, &book.Year, &book.CreatedAt, &book.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan book: %w", err)
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating books: %w", err)
	}

	countQuery := `SELECT COUNT(*) FROM books WHERE category_id = ANY($1)`
	var total int
	err = r.db.QueryRow(ctx, countQuery, categoryIDs).Scan(&total)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, 0, nil
		}
		return nil, 0, fmt.Errorf("failed to count books: %w", err)
	}

	return books, total, nil
}
