package tests

import (
	"fmt"
	"log"
	"notif/internal/data" // Ensure this matches your module name
)

func main() {
	fmt.Println("Attempting to read TA Profiles...")

	// 1. Call the function
	profiles, err := data.GetAllProfiles()

	// 2. Check for errors (like "file not found")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// 3. Print the results
	fmt.Printf("Success! Found %d profiles.\n", len(profiles))
	fmt.Println("--------------------------------------------------")

	for i, p := range profiles {
		fmt.Printf("[%d] ID: %s | Name: %s | Title: %s\n", i+1, p.ComputingID, p.Name, p.Title)
		fmt.Printf("    Bio: %s\n\n", p.Bio)
	}
}
