package subFunctions

import (
	"github.com/pterm/pterm"
	"log"
)

const (
	filename = "PortAudit.xlsx"
)

/*
Decide whether to create a new Excel file or update an existing one based on the input (--base flag), and then compares the Excel sheets.

Parameters:
  allData []InterfaceData - A slice containing all the data to be written to or updated in the Excel file.
  baseFile bool - A boolean flag that determines the operation:
                  true to create a new Excel file, false to update an existing file.

This function does not return any value but will halt execution and log a fatal error if any step fails.
*/

func ExcelOperations(allData []InterfaceData, baseFile bool, logger *pterm.Logger) {
	var err error
	if baseFile {
		err = CreateExcel(allData, filename)
	} else {
		err = UpdateExcel(allData, filename)
	}
	if err != nil {
		log.Fatalf("Failed to manage Excel file: %v", err)
	}
	log.Printf("Excel operations completed successfully on '%s'.", filename)

	// Compare data in Excel sheets.
	if err = CompareExcelSheets(filename, logger); err != nil {
		log.Fatalf("Failed during Excel sheet comparison: %v", err)
	}
}
