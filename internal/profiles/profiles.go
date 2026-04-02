package profiles

import (
	"fmt"
	"notif/internal/models"

	"github.com/xuri/excelize/v2"
)

func GetAllProfiles() ([]models.Profile, error) {
	f, err := excelize.OpenFile("profiles.xlsx")
	if err != nil {

		return nil, fmt.Errorf("could not open profiles file: %v", err)

	}

	defer f.Close()

	rows, err := f.GetRows("Sheet1")

	if err != nil {

		return nil, err

	}

	var profiles []models.Profile

	// Loop through rows (skipping header if you have one, or starting at 0)

	for _, row := range rows {

		if len(row) < 4 {

			continue // Skip incomplete rows

		}

		profiles = append(profiles, models.Profile{

			ComputingID: row[0],

			Name: row[1],

			Title: row[2],

			Bio: row[3],

			// PhotoURL would be row[4] if exists

		})

	}

	return profiles, nil

}
