package product

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var (
	ErrNotFound  = errors.New("product not found")
	ErrInvalidId = errors.New("provided product id is not valid")
)

func List(ctx context.Context, db *sqlx.DB) ([]Product, error) {
	list := []Product{}

	if err := db.SelectContext(ctx, &list, "SELECT * FROM products"); err != nil {
		return nil, err
	}

	return list, nil
}

func Single(ctx context.Context, db *sqlx.DB, id string) (Product, error) {
	var p Product

	if _, err := uuid.Parse(id); err != nil {
		return p, ErrInvalidId
	}

	if err := db.GetContext(ctx, &p, "SELECT * FROM products WHERE product_id = $1", id); err != nil {
		if err == sql.ErrNoRows {
			return p, ErrNotFound
		}

		return p, err
	}

	p.Name = "123"

	return p, nil
}

func Create(ctx context.Context, db *sqlx.DB, np NewProduct, now time.Time) (Product, error) {
	p := Product{
		ID:          uuid.New().String(),
		Name:        np.Name,
		Cost:        np.Cost,
		Quantity:    np.Quantity,
		DateCreated: now,
		DateUpdated: now,
	}

	q := `INSERT INTO products (product_id, name, cost, quantity, date_created, date_updated) VALUES ($1, $2, $3, $4, $5, $6)`

	if _, err := db.ExecContext(ctx, q, p.ID, p.Name, p.Cost, p.Quantity, p.DateCreated, p.DateUpdated); err != nil {
		return Product{}, errors.Wrapf(err, "insert into products: %v", np)
	}

	return p, nil
}
