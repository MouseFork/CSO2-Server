package sqlite

import (
	"database/sql"
)

func InitDatabase(file string) (*sql.DB, error) {
	DB, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}
	return DB, nil
}
