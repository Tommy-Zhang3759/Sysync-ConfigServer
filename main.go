package main

import (
	"ConfigServer/utils"
	"fmt"
)

func init() {

}

func main() {
	//var mvg sync.WaitGroup

	//go webUI.StartServer("8080", webUI.Handler)
	//
	////mvg.Wait()
	//// var clients []clientManage.Client
	//for {
	//	time.Sleep(time.Hour)
	//}

	csv, err := utils.OpenCSV("./test/example.csv")
	if err != nil {
		println(err.Error())
		return
	}
	fmt.Println(csv.GetAllData())
	fmt.Println(csv.GetRowData(1))
	fmt.Println(csv.GetRowData(2))
	fmt.Println(csv.GetRowData(4))
	fmt.Println(csv.GetCellData("Name", 1))
	fmt.Println(csv.GetCellData("Age", 2))
	fmt.Println(csv.GetCellData("City", 3))
	err = csv.SetRowData(1, []string{"Alice", "25", "New York"})
	if err != nil {
		return
	}
	err = csv.SetCellData("Name", 2, "Peler")
	if err != nil {
		return
	}
	fmt.Println(csv.GetAllData())
	fmt.Println(csv.GetRowData(1))
	fmt.Println(csv.GetRowData(2))
	fmt.Println(csv.GetRowData(3))
	fmt.Println(csv.GetCellData("Name", 1))
	fmt.Println(csv.GetCellData("Age", 2))
	fmt.Println(csv.GetCellData("City", 3))
	err = csv.SaveCSV()
	if err != nil {
		println(err.Error())
		return
	}
}
