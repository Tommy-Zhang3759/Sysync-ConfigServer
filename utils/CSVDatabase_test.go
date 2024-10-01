package utils

import (
	"encoding/csv"
	"os"
	"testing"
)

const testCSVFile = "./test.csv"

// 测试打开CSV文件
func TestOpenCSV(t *testing.T) {
	// 创建临时 CSV 文件
	createTestCSV(t, testCSVFile, [][]string{
		{"Name", "Age", "City"},
		{"Alice", "30", "New York"},
		{"Bob", "25", "Los Angeles"},
	})

	defer removeTestCSV(t, testCSVFile)

	// 打开 CSV 文件
	db, err := OpenCSV(testCSVFile)
	if err != nil {
		t.Fatalf("Failed to open CSV file: %v", err)
	}

	// 验证读取的数据
	if len(db.csvData) != 3 {
		t.Errorf("Expected 3 rows, but got %d", len(db.csvData))
	}

	if db.csvData[1][0] != "Alice" {
		t.Errorf("Expected 'Alice' in first column, but got %s", db.csvData[1][0])
	}
}

// 测试保存CSV文件
func TestSaveCSV(t *testing.T) {
	// 创建一个新的 CSV 数据库实例
	db := CSVDataBase{
		filePath: testCSVFile,
		csvData: [][]string{
			{"Name", "Age", "City"},
			{"Charlie", "35", "Chicago"},
		},
	}

	// 保存 CSV 文件
	err := db.SaveCSV()
	if err != nil {
		t.Fatalf("Failed to save CSV file: %v", err)
	}

	// 打开并验证 CSV 文件
	openedDb, err := OpenCSV(testCSVFile)
	if err != nil {
		t.Fatalf("Failed to open CSV file after saving: %v", err)
	}

	if len(openedDb.csvData) != 2 {
		t.Errorf("Expected 2 rows, but got %d", len(openedDb.csvData))
	}
	defer removeTestCSV(t, testCSVFile)
}

// 测试获取一行数据
func TestGetRowData(t *testing.T) {
	db := createMockCSVDatabase()

	// 获取第一行数据
	row, err := db.GetRowData(1)
	if err != nil {
		t.Fatalf("Failed to get row data: %v", err)
	}

	if row[0] != "Alice" {
		t.Errorf("Expected 'Alice', but got %s", row[0])
	}

	// 测试越界
	_, err = db.GetRowData(10)
	if err == nil {
		t.Errorf("Expected an out of range error, but got nil")
	}
}

// 测试获取单元格数据
func TestGetCellData(t *testing.T) {
	db := createMockCSVDatabase()

	// 获取单元格数据
	cell, err := db.GetCellData("City", 1)
	if err != nil {
		t.Fatalf("Failed to get cell data: %v", err)
	}

	if cell != "New York" {
		t.Errorf("Expected 'New York', but got %s", cell)
	}

	// 测试越界
	_, err = db.GetCellData("City", 10)
	if err == nil {
		t.Errorf("Expected an out of range error, but got nil")
	}

	// 测试键不存在
	_, err = db.GetCellData("Nonexistent", 1)
	if err == nil {
		t.Errorf("Expected a key not found error, but got nil")
	}
}

// 测试设置行数据
func TestSetRowData(t *testing.T) {
	db := createMockCSVDatabase()

	// 更新第一行数据
	newRow := []string{"Alice", "31", "San Francisco"}
	err := db.SetRowData(1, newRow)
	if err != nil {
		t.Fatalf("Failed to set row data: %v", err)
	}

	// 验证更新是否成功
	row, _ := db.GetRowData(1)
	if row[2] != "San Francisco" {
		t.Errorf("Expected 'San Francisco', but got %s", row[2])
	}

	// 测试越界
	err = db.SetRowData(10, newRow)
	if err == nil {
		t.Errorf("Expected an out of range error, but got nil")
	}
}

// 测试设置单元格数据
func TestSetCellData(t *testing.T) {
	db := createMockCSVDatabase()

	// 更新单元格数据
	err := db.SetCellData("City", 1, "Boston")
	if err != nil {
		t.Fatalf("Failed to set cell data: %v", err)
	}

	// 验证更新是否成功
	cell, _ := db.GetCellData("City", 1)
	if cell != "Boston" {
		t.Errorf("Expected 'Boston', but got %s", cell)
	}

	// 测试键不存在
	err = db.SetCellData("Nonexistent", 1, "Boston")
	if err == nil {
		t.Errorf("Expected a key not found error, but got nil")
	}

	// 测试越界
	err = db.SetCellData("City", 10, "Boston")
	if err == nil {
		t.Errorf("Expected an out of range error, but got nil")
	}
}

// 辅助函数，创建模拟的CSV数据
func createMockCSVDatabase() *CSVDataBase {
	return &CSVDataBase{
		filePath: "mock.csv",
		csvData: [][]string{
			{"Name", "Age", "City"},
			{"Alice", "30", "New York"},
			{"Bob", "25", "Los Angeles"},
		},
	}
}

// 创建临时CSV文件用于测试
func createTestCSV(t *testing.T, path string, data [][]string) {
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create test CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range data {
		if err := writer.Write(record); err != nil {
			t.Fatalf("Failed to write to test CSV file: %v", err)
		}
	}
}

// 删除测试CSV文件
func removeTestCSV(t *testing.T, path string) {
	err := os.Remove(path)
	if err != nil {
		t.Fatalf("Failed to remove test CSV file: %v", err)
	}
}
