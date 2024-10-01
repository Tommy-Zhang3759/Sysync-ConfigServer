package DataFrame

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

const testDB = "./test.db"

// 测试数据库连接
func TestSQLite_Connect(t *testing.T) {
	db := &SQLite{}

	// 连接 SQLite 数据库
	err := db.Connect(testDB)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// 确保成功连接
	if err := db.Close(); err != nil {
		t.Fatalf("Failed to close database: %v", err)
	}
}

// 测试插入数据
func TestSQLite_Insert(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	// 插入数据
	insertQuery := "INSERT INTO users (name, age) VALUES (?, ?)"
	_, err := db.Insert(insertQuery, "Alice", 30)
	if err != nil {
		t.Fatalf("Failed to insert data: %v", err)
	}
}

// 测试查询数据
func TestSQLite_Query(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	// 插入测试数据
	insertQuery := "INSERT INTO users (name, age) VALUES (?, ?)"
	_, err := db.Insert(insertQuery, "Bob", 25)
	if err != nil {
		t.Fatalf("Failed to insert data for query test: %v", err)
	}

	// 查询数据
	query := "SELECT id, name, age FROM users WHERE age > ?"
	rows, err := db.Query(query, 20)
	if err != nil {
		t.Fatalf("Failed to query data: %v", err)
	}
	defer rows.Close()

	// 检查查询结果
	var count int
	for rows.Next() {
		var id int
		var name string
		var age int
		if err := rows.Scan(&id, &name, &age); err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}
		count++
		if name != "Bob" {
			t.Errorf("Expected name to be 'Bob', but got %s", name)
		}
	}
	if count != 1 {
		t.Errorf("Expected 1 row, but got %d", count)
	}
}

// 测试更新数据
func TestSQLite_Update(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	// 插入数据
	insertQuery := "INSERT INTO users (name, age) VALUES (?, ?)"
	result, err := db.Insert(insertQuery, "Charlie", 40)
	if err != nil {
		t.Fatalf("Failed to insert data: %v", err)
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Failed to get last insert ID: %v", err)
	}

	// 更新数据
	updateQuery := "UPDATE users SET age = ? WHERE id = ?"
	_, err = db.Update(updateQuery, 45, id)
	if err != nil {
		t.Fatalf("Failed to update data: %v", err)
	}

	// 查询验证更新
	query := "SELECT age FROM users WHERE id = ?"
	row := db.db.QueryRow(query, id)
	var updatedAge int
	err = row.Scan(&updatedAge)
	if err != nil {
		t.Fatalf("Failed to query updated data: %v", err)
	}
	if updatedAge != 45 {
		t.Errorf("Expected age to be 45, but got %d", updatedAge)
	}
}

// 测试删除数据
func TestSQLite_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	// 插入数据
	insertQuery := "INSERT INTO users (name, age) VALUES (?, ?)"
	result, err := db.Insert(insertQuery, "David", 50)
	if err != nil {
		t.Fatalf("Failed to insert data: %v", err)
	}

	// 获取插入的ID
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Failed to get last insert ID: %v", err)
	}

	// 删除数据
	deleteQuery := "DELETE FROM users WHERE id = ?"
	_, err = db.Delete(deleteQuery, id)
	if err != nil {
		t.Fatalf("Failed to delete data: %v", err)
	}

	// 查询验证删除
	query := "SELECT id FROM users WHERE id = ?"
	row := db.db.QueryRow(query, id)
	var deletedID int
	err = row.Scan(&deletedID)
	if err != sql.ErrNoRows {
		t.Errorf("Expected no rows, but got %v", deletedID)
	}
}

// 创建一个测试数据库
func setupTestDB(t *testing.T) *SQLite {
	db := &SQLite{}
	err := db.Connect(testDB)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// 创建测试表
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		age INTEGER
	);`
	_, err = db.Insert(createTableQuery)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	return db
}

// 清理测试数据库
func teardownTestDB(t *testing.T, db *SQLite) {
	err := db.Close()
	if err != nil {
		t.Fatalf("Failed to close test database: %v", err)
	}
	err = os.Remove(testDB) // 删除测试数据库文件
	if err != nil {
		t.Fatalf("Failed to remove test database: %v", err)
	}
}
