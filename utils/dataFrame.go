package utils

import (
	"encoding/csv"
	"errors"
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
		//fmt.Println("Error opening file:", err)
		return csvDataBase, err
	}

	var reader = csv.NewReader(file)

	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" { // 如果读取到文件末尾，err 会是 io.EOF，可以正常退出循环
				break
			}
			//fmt.Println("Error reading record:", err)
			return csvDataBase, err
		}
		csvDataBase.csvData = append(csvDataBase.csvData, record)
	}
	err = file.Close()
	if err != nil {
		return csvDataBase, err
	}
	return csvDataBase, nil
}

func (c *CSVDataBase) SaveCSV() error {
	var file, err = os.OpenFile(c.filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	var writer = csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range c.csvData {
		fmt.Println(record)
		if err := writer.Write(record); err != nil {
			//fmt.Println("Error writing record:", err)
			return err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *CSVDataBase) GetRowData(RowIndex int) ([]string, error) {
	if RowIndex >= 0 && RowIndex < len(c.csvData) {
		return c.csvData[RowIndex], nil
	} else {
		return nil, errors.New(fmt.Sprintf("Row index %d is out of range", RowIndex))
	}
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
			return "", errors.New("key not found: " + key)
		}
		return c.csvData[RowIndex][index], nil
	} else {
		return "", errors.New(fmt.Sprintf("Row index %d is out of range", RowIndex))
	}
}

func (c *CSVDataBase) GetAllData() ([][]string, error) {
	return c.csvData, nil
}

func (c *CSVDataBase) SetRowData(RowIndex int, data []string) error {
	if RowIndex >= 0 && RowIndex < len(c.csvData) {
		c.csvData[RowIndex] = data
		return nil
	} else {
		return errors.New(fmt.Sprintf("Row index %d is out of range", RowIndex))
	}
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
			//fmt.Println("Error getting cell data, cannot find key: ", key)
			return errors.New("key not found: " + key)
		}
		c.csvData[RowIndex][index] = data
		return nil
	} else {
		return errors.New(fmt.Sprintf("Row index %d is out of range", RowIndex))
	}
}

type DataFrame interface {
	GetRowData(RowIndex int) ([]string, error)
	GetCellData(key string, RowIndex int) (string, error)
	GetAllData() ([][]string, error)
	SetRowData(RowIndex int, data []string) error
	SetCellData(key string, RowIndex int, data string) error

	OpenCSV(path string) error
	SaveCSV() error
}
