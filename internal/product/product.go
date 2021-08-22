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
