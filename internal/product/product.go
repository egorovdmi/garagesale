package product

import (
	"github.com/jmoiron/sqlx"
)

func List(db *sqlx.DB) ([]Product, error) {
	list := []Product{}

	if err := db.Select(&list, "SELECT * FROM products"); err != nil {
		return nil, err
	}

	return list, nil
}

func Single(db *sqlx.DB, id string) (Product, error) {
	var p Product

	err := db.Get(&p, "SELECT * FROM products WHERE id = $1", id)
	return p, err
}
