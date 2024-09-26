package database

import "database/sql"

type Database interface {
	Connect(connectionString string) error
	Close() error
	Insert(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Update(query string, args ...interface{}) (sql.Result, error)
	Delete(query string, args ...interface{}) (sql.Result, error)
}
