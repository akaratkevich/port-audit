package internal

import (
	"github.com/tealeg/xlsx"
	"sync"
)

// Read data from a given Excel sheet and returns a slice of InterfaceData.
func ReadExcelData(sheet *xlsx.Sheet, headerMap map[string]int, ch chan<- []InterfaceData, wg *sync.WaitGroup) {
	defer wg.Done()

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

	ch <- data
}

// Get the value of a cell by index with a fallback for missing cells.
func getCellValue(row *xlsx.Row, headerMap map[string]int, header string) string {
	if idx, ok := headerMap[header]; ok && idx < len(row.Cells) {
		return row.Cells[idx].String()
	}
	return ""
}
