// Memoiz backend
//
// DB storage.
//
// Rémy Mathieu © 2016

package storage

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage struct {
	Conn *sql.DB
}

// Shared Storage connection.
// TODO(remy): handle a pool of connection
var st Storage

// ----------------------

func DB() *sql.DB {
	return st.Conn
}

// Init opens a PostgreSQL connection with the given connectionString.
func Init(connectionString string) (*sql.DB, error) {
	dbase, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	st.Conn = dbase

	return dbase, st.Conn.Ping()
}
