package database

import (
	"database/sql"
	"fmt"
	"log"

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

func demo() {
	// 1. 连接到 SQLite 数据库，创建数据库文件（如果不存在）
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	// 2. 创建表
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"name" TEXT,
		"age" INTEGER
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("创建表失败: %s", err)
	}
	fmt.Println("表已经创建或存在")

	// 3. 插入数据
	insertUserSQL := `INSERT INTO users(name, age) VALUES (?, ?)`
	_, err = db.Exec(insertUserSQL, "Alice", 25)
	if err != nil {
		log.Fatalf("插入数据失败: %s", err)
	}
	fmt.Println("数据已插入")

	// 4. 查询数据
	query := `SELECT id, name, age FROM users`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("查询数据失败: %s", err)
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	fmt.Println("查询结果：")
	for rows.Next() {
		var id int
		var name string
		var age int
		err = rows.Scan(&id, &name, &age)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", id, name, age)
	}

	// 5. 处理查询中的错误（如果有）
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
