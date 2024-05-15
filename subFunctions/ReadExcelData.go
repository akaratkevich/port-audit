package subFunctions

import "github.com/tealeg/xlsx"

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
	for _, row := range sheet.Rows[1:] { // Skip header row
		if len(row.Cells) < 8 {
			continue
		}
		d := InterfaceData{
			Node:        row.Cells[0].String(),
			Interface:   row.Cells[1].String(),
			Description: row.Cells[2].String(),
			Status:      row.Cells[3].String(),
			VLAN:        row.Cells[4].String(),
			Duplex:      row.Cells[5].String(),
			Speed:       row.Cells[6].String(),
			Type:        row.Cells[7].String(),
		}
		data = append(data, d)
	}
	return data, nil
}
