package sqlrepo

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// NewDB opens a Postgres connection, verifies it with a Ping, and returns
// the *sql.DB handle. The caller is responsible for calling db.Close().
func NewDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
