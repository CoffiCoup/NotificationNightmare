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

// SaveTAProfile: Handles TA Bio Input -> TA_Profiles.xlsx
func SaveTAProfile(p models.Profile) error {
	mu.Lock()
	defer mu.Unlock()

	f, err := excelize.OpenFile("TA_Profiles.xlsx")
	if err != nil {
		f = excelize.NewFile()
	}
	defer f.Close()

	rows, _ := f.GetRows("Sheet1")
	newRow := len(rows) + 1

	// Columns: A=ID, B=Name, C=Title, D=Bio, E=PhotoURL
	f.SetCellValue("Sheet1", fmt.Sprintf("A%d", newRow), p.ComputingID)
	f.SetCellValue("Sheet1", fmt.Sprintf("B%d", newRow), p.Name)
	f.SetCellValue("Sheet1", fmt.Sprintf("C%d", newRow), p.Title)
	f.SetCellValue("Sheet1", fmt.Sprintf("D%d", newRow), p.Bio)
	f.SetCellValue("Sheet1", fmt.Sprintf("E%d", newRow), p.PhotoURL)

	return f.SaveAs("TA_Profiles.xlsx")
}

// GetAllProfiles: Reads the Excel file to display TAs on the website
func GetAllProfiles() ([]models.Profile, error) {
	// Make sure the path points to where the file is (internal/data/)
	f, err := excelize.OpenFile("internal/data/TA_Profiles.xlsx")
	if err != nil {
		return nil, fmt.Errorf("could not open profiles file: %v", err)
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	var profiles []models.Profile

	for i, row := range rows {
		// 1. SKIP THE HEADER ROW (Row 1 in Excel is index 0 in Go)
		if i == 0 {
			continue
		}

		// 2. BASIC LENGTH CHECK
		if len(row) < 4 {
			continue
		}

		// 3. OPTIONAL PHOTO URL CHECK
		photo := ""
		if len(row) > 4 {
			photo = row[4]
		}

		profiles = append(profiles, models.Profile{
			ComputingID: row[0],
			Name:        row[1],
			Title:       row[2],
			Bio:         row[3],
			PhotoURL:    photo,
		})
	}

	return profiles, nil
}
