package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

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
