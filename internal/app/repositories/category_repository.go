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

type CategoryRepositoryImpl struct {
	db *postgres.DBConnection
}

func NewCategoryRepository(db *postgres.DBConnection) *CategoryRepositoryImpl {
	return &CategoryRepositoryImpl{db: db}
}

func (r *CategoryRepositoryImpl) CreateCategory(ctx context.Context, category sm.DomainCategory) (rm.Category, error) {
	query := `INSERT INTO categories (name) 
              VALUES ($1) 
              RETURNING id, name, created_at, updated_at`

	var newCategory rm.Category
	err := r.db.QueryRow(ctx, query, category.Name).Scan(
		&newCategory.Id, &newCategory.Name, &newCategory.CreatedAt, &newCategory.UpdatedAt,
	)
	if err != nil {
		return rm.Category{}, fmt.Errorf("failed to create category: %w", err)
	}

	return newCategory, nil
}

func (r *CategoryRepositoryImpl) GetCategory(ctx context.Context, id int) (rm.Category, error) {
	if id == 0 {
		return rm.Category{}, fmt.Errorf("id can not be 0")
	}

	var category rm.Category
	query := `SELECT id, name, created_at, updated_at 
              	FROM categories 
              WHERE id = $1`

	row := r.db.QueryRow(ctx, query, id)
	err := row.Scan(
		&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return rm.Category{}, fmt.Errorf("category not found")
		}
		return rm.Category{}, fmt.Errorf("failed to get a category: %w", err)
	}

	return category, nil
}

func (r *CategoryRepositoryImpl) UpdateCategory(ctx context.Context, id int, category sm.DomainCategory) (rm.Category, error) {
	query := `UPDATE categories SET name = $1, updated_at = NOW()
              WHERE id = $2 
              RETURNING id, name, created_at, updated_at`

	var updatedCategory rm.Category
	err := r.db.QueryRow(ctx, query, category.Name, id).Scan(
		&updatedCategory.Id, &updatedCategory.Name, &updatedCategory.CreatedAt, &updatedCategory.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return rm.Category{}, fmt.Errorf("category not found")
		}
		return rm.Category{}, fmt.Errorf("failed to update book: %w", err)
	}

	return updatedCategory, nil
}

func (r *CategoryRepositoryImpl) DeleteCategory(ctx context.Context, id int) error {
	if id == 0 {
		return fmt.Errorf("id can not be 0")
	}

	query := `DELETE FROM categories 
              	WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}

func (r *CategoryRepositoryImpl) GetCategories(ctx context.Context) ([]rm.Category, error) {
	query := `SELECT id, name, created_at, updated_at 
        	  FROM categories  
        	  ORDER BY id`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get books: %w", err)
	}
	defer rows.Close()

	var categories []rm.Category
	for rows.Next() {
		var category rm.Category
		err := rows.Scan(
			&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating categories: %w", err)
	}

	return categories, nil
}
