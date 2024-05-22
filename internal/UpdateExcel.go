package internal

import (
	"fmt"
	"github.com/pterm/pterm"
	"github.com/tealeg/xlsx"
	"log"
	"time"
)

// Open an existing Excel file, add a new sheet with current date and time, and populate it with data.
func UpdateExcel(data []InterfaceData, filename string, logger *pterm.Logger) error {
	// Open the existing Excel file.
	file, err := xlsx.OpenFile(filename)
	if err != nil {
		return err // Return the error if the file cannot be opened.
	}

	// Format the current date and time as 'DDMMYY_HHMM' and create a sheet name with it.
	dateTime := time.Now().Format("02012006") // Format: DDMMYYYY
	sheetName := fmt.Sprintf("Audit %s", dateTime)
	sheet, err := file.AddSheet(sheetName)
	if err != nil {
		logger.Fatal("Failed to add sheet", logger.Args("Reason", err))
		log.Fatalf("Failed to add sheet: %v", err) // Log and return the error if a new sheet cannot be added.
		return err
	}

	// Column headers based on the InterfaceData struct fields.
	headers := []string{"Switch Name", "Interface", "SLOT", "PORT", "TYPE", "Port Status", "VLAN", "Duplex", "SPEED", "Port Description"}
	headerRow := sheet.AddRow()
	for _, header := range headers {
		cell := headerRow.AddCell()
		cell.Value = header
	}

	// Populate the new sheet with data from the InterfaceData slice.
	for _, ci := range data {
		row := sheet.AddRow()
		row.AddCell().Value = ci.Node
		row.AddCell().Value = ci.Interface
		row.AddCell().Value = ci.Slot // Add Slot to Excel
		row.AddCell().Value = ci.Port // Add Port to Excel
		row.AddCell().Value = ci.Type
		row.AddCell().Value = ci.Status
		row.AddCell().Value = ci.VLAN
		row.AddCell().Value = ci.Duplex
		row.AddCell().Value = ci.Speed
		// Check if Description is blank and set a default value if it is
		description := ci.Description
		if description == "" {
			description = "Unallocated" // Set default value "Unallocated" as per the Baseline
		}
		row.AddCell().Value = description
	}

	// Save the updated Excel file.
	err = file.Save(filename)
	if err != nil {
		log.Fatalf("Failed to save Excel file: %v", err) // Log and return the error if the file cannot be saved.
		return err
	}

	log.Printf("Excel file '%s' updated successfully with new sheet '%s'.", filename, sheetName)
	return nil // Return nil to indicate success without errors.
}
