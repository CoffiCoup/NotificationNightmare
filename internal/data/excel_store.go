package data

import (
	"fmt"
	"sync"

	"notif/internal/models"

	"github.com/xuri/excelize/v2"
)

var mu sync.Mutex

// SaveOHRequest appends a student's request to the OH Spreadsheet
func SaveOHRequest(filePath string, req models.OHRequest) error {
	mu.Lock()
	defer mu.Unlock()

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("could not open excel file: %v", err)
	}
	defer f.Close()

	// Find the next empty row
	rows, _ := f.GetRows("Sheet1")
	newRow := len(rows) + 1

	// Fill columns A, B, and C based on your flowchart
	f.SetCellValue("Sheet1", fmt.Sprintf("A%d", newRow), req.ComputingID)
	f.SetCellValue("Sheet1", fmt.Sprintf("B%d", newRow), req.Reason)
	f.SetCellValue("Sheet1", fmt.Sprintf("C%d", newRow), req.DateTime)

	return f.Save()
}

// SaveAvailability: Populates the Calendar students see
func SaveAvailability(a models.Availability) error {
	mu.Lock()
	defer mu.Unlock()

	f, err := excelize.OpenFile("TA_Availability.xlsx")
	if err != nil {
		f = excelize.NewFile()
	}
	defer f.Close()

	rows, _ := f.GetRows("Sheet1")
	newRow := len(rows) + 1

	// Updated Columns: A=ID, B=Name, C=Date, D=Start, E=Duration, F=Location
	f.SetCellValue("Sheet1", fmt.Sprintf("A%d", newRow), a.ComputingID)
	f.SetCellValue("Sheet1", fmt.Sprintf("B%d", newRow), a.Name) // New Column
	f.SetCellValue("Sheet1", fmt.Sprintf("C%d", newRow), a.Date)
	f.SetCellValue("Sheet1", fmt.Sprintf("D%d", newRow), a.StartTime)
	f.SetCellValue("Sheet1", fmt.Sprintf("E%d", newRow), a.Duration)
	f.SetCellValue("Sheet1", fmt.Sprintf("F%d", newRow), a.Location)

	return f.SaveAs("TA_Availability.xlsx")
}
