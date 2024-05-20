package internal

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"log"
)

/*
Generate an Excel file with data from a slice of InterfaceData.

Parameters:
  data []InterfaceData - A slice containing the data to be written to the Excel sheet.
  filename string - The path and name of the file where the Excel sheet will be saved.

Returns:
  error - Returns an error if any step of the Excel file creation or saving process fails.
*/

func CreateExcel(data []InterfaceData, filename string) error {
	file := xlsx.NewFile()

	sheetName := fmt.Sprintf("Baseline")
	sheet, err := file.AddSheet(sheetName)
	if err != nil {
		// Log the error and return it.
		log.Fatalf("Failed to add sheet: %v", err)
		return err
	}

	// Column headers for the data to be inserted. These headers correspond to the fields within the InterfaceData struct.
	headers := []string{"Switch Name", "Interface", "SLOT", "PORT", "Description", "Status", "VLAN", "Duplex", "Speed", "Type"}
	headerRow := sheet.AddRow()
	for _, header := range headers {
		// For each header, add a new cell to the row and set its value.
		cell := headerRow.AddCell()
		cell.Value = header
	}

	// Iterate over the slice of InterfaceData to populate the sheet.
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

	// Save the Excel file
	err = file.Save(filename)
	if err != nil {
		log.Fatalf("Failed to save Excel file: %v", err)
		return err
	}

	log.Printf("Excel file saved successfully: %v", filename)
	return nil
}
