package internal

import (
	"fmt"
	"github.com/tealeg/xlsx"
)

/*
Extract interface data from an Excel sheet and convert it into a slice of InterfaceData structs.

Parameters:
  - sheet *xlsx.Sheet: The Excel sheet from which data will be read.

Returns:
  - []InterfaceData: A slice containing all the interface data extracted from the Excel sheet.
  - error: Returns an error if any issues occur during the data extraction process.
*/

func ReadExcelData(sheet *xlsx.Sheet) ([]InterfaceData, error) {
	var data []InterfaceData
	var headerMap = make(map[string]int)

	for i, cell := range sheet.Rows[0].Cells {
		headerMap[cell.String()] = i
	}

	// Required headers
	requiredHeaders := []string{"Switch Name", "SLOT", "PORT", "TYPE", "Port Status", "VLAN", "Duplex", "SPEED", "Port Description"} //"Interface" is not in the ref sheet

	// Check if all required headers are present
	for _, header := range requiredHeaders {
		if _, ok := headerMap[header]; !ok {
			return nil, fmt.Errorf("Missing required header: %s", header)
		}
	}

	for _, row := range sheet.Rows[1:] { // Skip the header row
		entry := InterfaceData{
			Node:        getCellValue(row, headerMap["Switch Name"]),
			Interface:   getCellValue(row, headerMap["Interface"]),
			Slot:        getCellValue(row, headerMap["SLOT"]),
			Port:        getCellValue(row, headerMap["PORT"]),
			Description: getCellValue(row, headerMap["Port Description"]),
			Status:      getCellValue(row, headerMap["Port Status"]),
			Speed:       getCellValue(row, headerMap["SPEED"]),
			Duplex:      getCellValue(row, headerMap["Duplex"]),
			VLAN:        getCellValue(row, headerMap["VLAN"]),
			Type:        getCellValue(row, headerMap["TYPE"]),
		}
		data = append(data, entry)
	}
	return data, nil
}

// Get the value of a cell by index with a fallback for missing cells.
func getCellValue(row *xlsx.Row, index int) string {
	if index < len(row.Cells) {
		return row.Cells[index].String()
	}
	return ""
}
