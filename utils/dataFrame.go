package utils

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
)

type CSVDataBase struct {
	file    os.File
	reader  *csv.Reader
	csvData [][]string
}

func (c *CSVDataBase) OpenDB(path string) error {
	// 打开CSV文件
	var file, err = os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	c.file = *file

	// 创建一个CSV reader
	c.reader = csv.NewReader(file)

	c.csvData = make([][]string, 0)

	// 逐行读取CSV内容
	for {
		record, err := c.reader.Read()
		if err != nil {
			// 如果读取到文件末尾，err 会是 io.EOF，可以正常退出循环
			if err.Error() == "EOF" {
				break
			}
			fmt.Println("Error reading record:", err)
			return err
		}
		c.csvData = append(c.csvData, record)
		// 处理每行的内容
		fmt.Println("Record:", record)
	}
	return nil
}

func (c *CSVDataBase) CloseDB() error {
	return c.file.Close()
}

func (c *CSVDataBase) GetRowData(RawIndex int) ([]string, error) {
	return c.csvData[RawIndex], nil
}

func (c *CSVDataBase) GetCellData(key string, RawIndex int) (string, error) {
	found := false
	index := 0
	for index = range c.csvData[0] {
		if c.csvData[index][index] == key {
			found = true
			break
		}
	}
	if !found {
		fmt.Println("Error getting cell data, cannot find key: ", key)
		return "", errors.New("key not found")
	}
	return c.csvData[RawIndex][index], nil
}

func (c *CSVDataBase) SetRowData(RawIndex int, data []string) error {
	c.csvData[RawIndex] = data
	return nil
}

func (c *CSVDataBase) SetCellData(key string, RawIndex int, data string) error {
	found := false
	index := 0
	for index = range c.csvData[0] {
		if c.csvData[index][index] == key {
			found = true
			break
		}
	}
	if !found {
		fmt.Println("Error getting cell data, cannot find key: ", key)
		return errors.New("key not found")
	}
	c.csvData[RawIndex][index] = data
	return nil
}

type DataFrame interface {
	GetRowData(RawIndex int) ([]string, error)
	GetCellData(key string, RawIndex int) (string, error)
	SetRowData(RawIndex int, data []string) error
	SetCellData(key string, RawIndex int, data string) error

	OpenDB(path string) error
	CloseDB() error
}
