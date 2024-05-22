package internal

import (
	"github.com/tealeg/xlsx"
)

// ReadExcelData reads data from a given Excel sheet and returns a slice of InterfaceData.
func ReadExcelData(sheet *xlsx.Sheet, headerMap map[string]int) []InterfaceData {
	var data []InterfaceData

	for _, row := range sheet.Rows[1:] { // Skip the header row
		entry := InterfaceData{
			Node:        getCellValue(row, headerMap, "Switch Name"),
			Slot:        getCellValue(row, headerMap, "SLOT"),
			Port:        getCellValue(row, headerMap, "PORT"),
			Description: getCellValue(row, headerMap, "Port Description"),
			Status:      getCellValue(row, headerMap, "Port Status"),
			Speed:       getCellValue(row, headerMap, "SPEED"),
			Duplex:      getCellValue(row, headerMap, "Duplex"),
			VLAN:        getCellValue(row, headerMap, "VLAN"),
			Type:        getCellValue(row, headerMap, "TYPE"),
		}
		data = append(data, entry)
	}

	return data
}

// Get the value of a cell by index with a fallback for missing cells.
func getCellValue(row *xlsx.Row, headerMap map[string]int, header string) string {
	if idx, ok := headerMap[header]; ok && idx < len(row.Cells) {
		return row.Cells[idx].String()
	}
	return ""
}

// GetHeaderMap reads the headers from the first row of the sheet and returns a map of header names to their indices.
func getHeaderMap(sheet *xlsx.Sheet) map[string]int {
	headerMap := make(map[string]int)
	for i, cell := range sheet.Rows[0].Cells {
		headerMap[cell.String()] = i
	}
	return headerMap
}
