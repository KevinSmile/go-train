package w2_error

import (
	"database/sql"
	"github.com/pkg/errors"
)

func getNameById(id int, db *sql.DB) (string, error) {
	var name string
	err := db.QueryRow("SELECT product_name FROM products WHERE id = ?", id).Scan(&name)
	if err == sql.ErrNoRows {
		return name, nil
	}
	return name, errors.Wrap(err, "getNameById failed")
}
