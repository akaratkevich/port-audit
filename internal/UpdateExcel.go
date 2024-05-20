package internal

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"log"
	"time"
)

// Open an existing Excel file, add a new sheet with current date and time, and populate it with data.
func UpdateExcel(data []InterfaceData, filename string) error {
	// Open the existing Excel file.
	file, err := xlsx.OpenFile(filename)
	if err != nil {
		return err // Return the error if the file cannot be opened.
	}

	// Format the current date and time as 'DDMMYY_HHMM' and create a sheet name with it.
	dateTime := time.Now().Format("020106_1504") // Format: DDMMYY_HHMM
	sheetName := fmt.Sprintf("Audit %s", dateTime)
	sheet, err := file.AddSheet(sheetName)
	if err != nil {
		log.Fatalf("Failed to add sheet: %v", err) // Log and return the error if a new sheet cannot be added.
		return err
	}

	// Column headers based on the InterfaceData struct fields.
	headers := []string{"Switch Name", "Interface", "SLOT", "PORT", "Description", "Status", "VLAN", "Duplex", "Speed", "Type"}
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
		row.AddCell().Value = ci.Description
		row.AddCell().Value = ci.Status
		row.AddCell().Value = ci.VLAN
		row.AddCell().Value = ci.Duplex
		row.AddCell().Value = ci.Speed
		row.AddCell().Value = ci.Type
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
