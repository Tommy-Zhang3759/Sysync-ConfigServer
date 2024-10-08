package utils

import (
	"encoding/csv"
	"fmt"
	"os"
)

type CSVDataBase struct {
	filePath string
	csvData  [][]string
}

func OpenCSV(path string) (CSVDataBase, error) {
	var csvDataBase = CSVDataBase{
		filePath: path,
		csvData:  make([][]string, 0),
	}

	var file, err = os.Open(path)
	if err != nil {
		return csvDataBase, fmt.Errorf("failed to open CSV file %s: %w", path, err)
	}

	var reader = csv.NewReader(file)

	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" { // 如果读取到文件末尾，err 会是 io.EOF，可以正常退出循环
				break
			}
			return csvDataBase, fmt.Errorf("error reading CSV file: %w", err)
		}
		csvDataBase.csvData = append(csvDataBase.csvData, record)
	}

	if err = file.Close(); err != nil {
		return csvDataBase, fmt.Errorf("error closing CSV file: %w", err)
	}

	return csvDataBase, nil
}

func (c *CSVDataBase) SaveCSV() error {
	var file, err = os.OpenFile(c.filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}

	var writer = csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range c.csvData {
		fmt.Println(record)
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record to CSV file: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("error writing CSV file: %w", err)
	}

	if err = file.Close(); err != nil {
		return fmt.Errorf("error closing file after writing: %w", err)
	}

	return nil
}

func (c *CSVDataBase) GetRowData(RowIndex int) ([]string, error) {
	if RowIndex >= 0 && RowIndex < len(c.csvData) {
		return c.csvData[RowIndex], nil
	}
	return nil, fmt.Errorf("row index %d is out of range", RowIndex)
}

func (c *CSVDataBase) GetCellData(key string, RowIndex int) (string, error) {
	if RowIndex >= 0 && RowIndex < len(c.csvData) {
		found := false
		index := 0
		for index = range c.csvData[0] {
			if c.csvData[0][index] == key {
				found = true
				break
			}
		}
		if !found {
			return "", fmt.Errorf("key not found: %s", key)
		}
		return c.csvData[RowIndex][index], nil
	}
	return "", fmt.Errorf("row index %d is out of range", RowIndex)
}

func (c *CSVDataBase) GetAllData() ([][]string, error) {
	return c.csvData, nil
}

func (c *CSVDataBase) SetRowData(RowIndex int, data []string) error {
	if RowIndex >= 0 && RowIndex < len(c.csvData) {
		c.csvData[RowIndex] = data
		return nil
	}
	return fmt.Errorf("row index %d is out of range", RowIndex)
}

func (c *CSVDataBase) SetCellData(key string, RowIndex int, data string) error {
	if RowIndex >= 0 && RowIndex < len(c.csvData) {
		found := false
		index := 0
		for index = range c.csvData[0] {
			if c.csvData[0][index] == key {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("key not found: %s", key)
		}
		c.csvData[RowIndex][index] = data
		return nil
	}
	return fmt.Errorf("row index %d is out of range", RowIndex)
}

type CSVDataFrame interface {
	GetRowData(RowIndex int) ([]string, error)
	GetCellData(key string, RowIndex int) (string, error)
	GetAllData() ([][]string, error)
	SetRowData(RowIndex int, data []string) error
	SetCellData(key string, RowIndex int, data string) error

	OpenCSV(path string) error
	SaveCSV() error
}
