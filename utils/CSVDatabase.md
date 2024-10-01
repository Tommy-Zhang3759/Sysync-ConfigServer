# TestCSV Function Documentation

## Package

```
import "ConfigServer/utils"
```

### Purpose
The `TestCSV` function is a unit test designed to validate the functionality of CSV file operations implemented in the `utils` package. It tests various methods for reading, updating, and saving data in a CSV file.

### Preconditions
The CSV file should contain data structured in a way that matches the expected format.

### Functionality Overview
The `TestCSV` function performs the following operations:

1. Open CSV File:
Calls `utils.OpenCSV` to open the specified CSV file. It will return an `CSVDataBase` object.
   
2. Read All Data:
Retrieves and prints all data from the CSV file using `CSVDataBase.GetAllData()`.

3. Retrieve Specific Row Data:
Calls `CSVDataBase.GetRowData()` to fetch data for specified row indices (1, 2, and 4).Demonstrates error handling when trying to access a non-existent row.

4. Retrieve Specific Cell Data:
Fetches cell data based on column names and row indices using `CSVDataBase.GetCellData()`.

5. Update Row Data:
Uses `CSVDataBase.SetRowData()` to update an entire row with new data.
Handles any errors that occur during this operation.

6. Update Specific Cell Data:
Updates a specific cell's data using `CSVDataBase.SetCellData()` and handles potential errors.
Final Data Display:
Prints all data, specific row data, and specific cell data after modifications.

7. Save Changes:
Saves the updated data back to the CSV file using `CSVDataBase.SaveCSV()`.
Handles any errors during the save operation.
## Code Breakdown
```
func TestCSV(t *testing.T) {
	csv_, err := utils.OpenCSV("./example.csv")
	if err != nil {
		println(err.Error())
		return
	}
	
	// Display all data
	fmt.Println(csv_.GetAllData())
	
	// Access specific rows and handle potential errors
	fmt.Println(csv_.GetRowData(1))
	fmt.Println(csv_.GetRowData(2))
	fmt.Println(csv_.GetRowData(4)) // This may produce an error if row 4 does not exist
	
	// Retrieve cell data by column name and row index
	fmt.Println(csv_.GetCellData("Name", 1))
	fmt.Println(csv_.GetCellData("Age", 2))
	fmt.Println(csv_.GetCellData("City", 3))
	
	// Update a row
	err = csv_.SetRowData(1, []string{"Alice", "25", "New York"})
	if err != nil {
		println(err.Error())
		return
	}
	
	// Update a specific cell
	err = csv_.SetCellData("Name", 2, "Peler")
	if err != nil {
		println(err.Error())
		return
	}
	
	// Display updated data
	fmt.Println(csv_.GetAllData())
	fmt.Println(csv_.GetRowData(1))
	fmt.Println(csv_.GetRowData(2))
	fmt.Println(csv_.GetRowData(3))
	fmt.Println(csv_.GetCellData("Name", 1))
	fmt.Println(csv_.GetCellData("Age", 2))
	fmt.Println(csv_.GetCellData("City", 3))
	
	// Save the modified data back to the CSV file
	err = csv_.SaveCSV()
	if err != nil {
		println(err.Error())
		return
	}
}
```
## Error Handling
The function includes multiple error checks to ensure that:

- The CSV file is opened successfully.
- Accessing specific rows and cells does not result in an out-of-bounds error.
- Updates to row and cell data are executed without issues.
- The final save operation is completed successfully.