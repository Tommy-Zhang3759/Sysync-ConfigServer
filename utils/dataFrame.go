package utils

import (
	"encoding/csv"
	"errors"
	"os"
)

type CSVDataBase struct {
	file    os.File
	csvData [][]string
}

func (c *CSVDataBase) OpenDB(path string) error {
	var file, err = os.Open(path)
	if err != nil {
		//fmt.Println("Error opening file:", err)
		return err
	}

	c.file = *file
	var reader = csv.NewReader(file)
	c.csvData = make([][]string, 0)

	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" { // 如果读取到文件末尾，err 会是 io.EOF，可以正常退出循环
				break
			}
			//fmt.Println("Error reading record:", err)
			return err
		}
		c.csvData = append(c.csvData, record)
	}
	return nil
}

func (c *CSVDataBase) SaveDB() error {
	var writer = csv.NewWriter(&c.file)
	defer writer.Flush()

	for _, record := range c.csvData {
		if err := writer.Write(record); err != nil {
			//fmt.Println("Error writing record:", err)
			return err
		}
	}
	return nil
}

func (c *CSVDataBase) CloseDB() error {
	return c.file.Close()
}

func (c *CSVDataBase) GetRowData(RowIndex int) ([]string, error) {
	return c.csvData[RowIndex], nil
}

func (c *CSVDataBase) GetCellData(key string, RowIndex int) (string, error) {
	found := false
	index := 0
	for index = range c.csvData[0] {
		if c.csvData[0][index] == key {
			found = true
			break
		}
	}
	if !found {
		//fmt.Println("Error getting cell data, cannot find key: ", key)
		return "", errors.New("key not found: " + key)
	}
	return c.csvData[RowIndex][index], nil
}

func (c *CSVDataBase) GetAllData() ([][]string, error) {
	return c.csvData, nil
}

func (c *CSVDataBase) SetRowData(RowIndex int, data []string) error {
	c.csvData[RowIndex] = data
	return nil
}

func (c *CSVDataBase) SetCellData(key string, RowIndex int, data string) error {
	found := false
	index := 0
	for index = range c.csvData[0] {
		if c.csvData[0][index] == key {
			found = true
			break
		}
	}
	if !found {
		//fmt.Println("Error getting cell data, cannot find key: ", key)
		return errors.New("key not found: " + key)
	}
	c.csvData[RowIndex][index] = data
	return nil
}

type DataFrame interface {
	GetRowData(RowIndex int) ([]string, error)
	GetCellData(key string, RowIndex int) (string, error)
	GetAllData() ([][]string, error)
	SetRowData(RowIndex int, data []string) error
	SetCellData(key string, RowIndex int, data string) error

	OpenDB(path string) error
	SaveDB() error
	CloseDB() error
}
