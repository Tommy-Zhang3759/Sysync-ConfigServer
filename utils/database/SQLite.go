package DataFrame

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// SQLite 实现了 Database 接口
type SQLite struct {
	db *sql.DB
}

// Connect 连接 SQLite 数据库
func (s *SQLite) Connect(connectionString string) error {
	var err error
	s.db, err = sql.Open("sqlite3", connectionString)
	if err != nil {
		return err
	}
	return s.db.Ping()
}

// Close 关闭数据库连接
func (s *SQLite) Close() error {
	return s.db.Close()
}

// Insert 插入数据
func (s *SQLite) Insert(query string, args ...interface{}) (sql.Result, error) {
	statement, err := s.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer func(statement *sql.Stmt) {
		_ = statement.Close()
	}(statement)
	return statement.Exec(args...)
}

// Query 查询数据
func (s *SQLite) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return s.db.Query(query, args...)
}

// Update 更新数据
func (s *SQLite) Update(query string, args ...interface{}) (sql.Result, error) {
	statement, err := s.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer func(statement *sql.Stmt) {
		_ = statement.Close()
	}(statement)
	return statement.Exec(args...)
}

// Delete 删除数据
func (s *SQLite) Delete(query string, args ...interface{}) (sql.Result, error) {
	statement, err := s.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer func(statement *sql.Stmt) {
		_ = statement.Close()
	}(statement)
	return statement.Exec(args...)
}
