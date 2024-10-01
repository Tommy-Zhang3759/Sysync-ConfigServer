package utils

import (
	"fmt"
	"testing"
)

func TestCSVDataBase(t *testing.T) {
	csv_, err := OpenCSV("./test/example.csv_")
	if err != nil {
		println(err.Error())
		return
	}
	fmt.Println(csv_.GetAllData())
	fmt.Println(csv_.GetRowData(1))
	fmt.Println(csv_.GetRowData(2))
	fmt.Println(csv_.GetRowData(4)) // make an error
	fmt.Println(csv_.GetCellData("Name", 1))
	fmt.Println(csv_.GetCellData("Age", 2))
	fmt.Println(csv_.GetCellData("City", 3))
	err = csv_.SetRowData(1, []string{"Alice", "25", "New York"})
	if err != nil {
		println(err.Error())
		return
	}
	err = csv_.SetCellData("Name", 2, "Peler")
	if err != nil {
		println(err.Error())
		return
	}
	fmt.Println(csv_.GetAllData())
	fmt.Println(csv_.GetRowData(1))
	fmt.Println(csv_.GetRowData(2))
	fmt.Println(csv_.GetRowData(3))
	fmt.Println(csv_.GetCellData("Name", 1))
	fmt.Println(csv_.GetCellData("Age", 2))
	fmt.Println(csv_.GetCellData("City", 3))
	err = csv_.SaveCSV()
	if err != nil {
		println(err.Error())
		return
	}
}
