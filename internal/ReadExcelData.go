package internal

import "github.com/tealeg/xlsx"

// ReadExcelData reads data from a given Excel sheet and returns a slice of InterfaceData.
func ReadExcelData(sheet *xlsx.Sheet) ([]InterfaceData, error) {
	var data []InterfaceData
	var headerMap = make(map[string]int)

	// Assume the first row is the header
	for i, cell := range sheet.Rows[0].Cells {
		headerMap[cell.String()] = i
	}

	// Map of expected headers to struct fields
	headerMapping := map[string]string{
		"Switch Name":      "Node",
		"Interface":        "Interface",
		"SLOT":             "Slot",
		"PORT":             "Port",
		"TYPE":             "Type",
		"Port Status":      "Status",
		"VLAN":             "VLAN",
		"Duplex":           "Duplex",
		"SPEED":            "Speed",
		"Port Description": "Description",
	}

	for _, row := range sheet.Rows[1:] { // Skip the header row
		entry := InterfaceData{
			Node:        getCellValue(row, headerMap, headerMapping["Switch Name"]),
			Interface:   getCellValue(row, headerMap, headerMapping["Interface"]),
			Slot:        getCellValue(row, headerMap, headerMapping["SLOT"]),
			Port:        getCellValue(row, headerMap, headerMapping["PORT"]),
			Type:        getCellValue(row, headerMap, headerMapping["TYPE"]),
			Status:      getCellValue(row, headerMap, headerMapping["Port Status"]),
			VLAN:        getCellValue(row, headerMap, headerMapping["VLAN"]),
			Duplex:      getCellValue(row, headerMap, headerMapping["Duplex"]),
			Speed:       getCellValue(row, headerMap, headerMapping["SPEED"]),
			Description: getCellValue(row, headerMap, headerMapping["Port Description"]),
		}
		data = append(data, entry)
	}
	return data, nil
}

// Get the value of a cell by index with a fallback for missing cells.
func getCellValue(row *xlsx.Row, headerMap map[string]int, header string) string {
	if idx, ok := headerMap[header]; ok && idx < len(row.Cells) {
		return row.Cells[idx].String()
	}
	return ""
}
