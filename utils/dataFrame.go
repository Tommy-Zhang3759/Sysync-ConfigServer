package utils

import (
	"encoding/csv"
	"fmt"
	"os"
)

type CSVDataBase struct {
	file   os.File
	reader *csv.Reader
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

	csvData := make([][]string, 0)

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
		csvData = append(csvData, record)
		// 处理每行的内容
		fmt.Println("Record:", record)
	}
	return nil
}

func (c *CSVDataBase) CloseDB() error {
	return c.file.Close()
}

func (c *CSVDataBase) GetRowData(key string, keyColNum int) []string {
	return nil
}

func (c *CSVDataBase) GetCellData(key string, keyColNum int, valueColNum int) []string {
	return nil
}

func (c *CSVDataBase) SaveRowData(key string, keyColNum int) error {
	return nil
}

func (c *CSVDataBase) SaveCellData(key string, keyColNum int, valueColNum int) error {
	return nil
}

type DataFrame interface {
	GetRowData(key string, keyColNum int) []string
	GetCellData(key string, keyColNum int, valueColNum int) []string
	SaveRowData(key string, keyColNum int) error
	SaveCellData(key string, keyColNum int, valueColNum int) error

	OpenDB(path string) error
	CloseDB() error
}
