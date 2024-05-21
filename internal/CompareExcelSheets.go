package internal

import (
	"fmt"
	"github.com/pterm/pterm"
	"github.com/tealeg/xlsx"
	"log"
	"os"
	"time"
)

// FilterData filters the reference data to include only relevant columns and devices.
func FilterData(refData, newData []InterfaceData) []InterfaceData {
	newNodes := make(map[string]bool)
	for _, d := range newData {
		newNodes[d.Node] = true
	}

	var filteredData []InterfaceData
	for _, d := range refData {
		if newNodes[d.Node] {
			filteredData = append(filteredData, d)
		}
	}
	return filteredData
}

// CompareExcelSheets compares two Excel sheets for differences.
func CompareExcelSheets(filename string, logger *pterm.Logger) error {
	file, err := xlsx.OpenFile(filename)
	if err != nil {
		return fmt.Errorf("Failed to open Excel file: %v", err)
	}

	dateTime := time.Now().Format("02012006") // Generates a timestamp for naming the new audit sheet.
	refSheet := file.Sheet["Baseline"]
	newSheetName := fmt.Sprintf("Audit %s", dateTime)
	newSheet := file.Sheet[newSheetName]

	if refSheet == nil || newSheet == nil {
		return fmt.Errorf("Missing Excel sheets for comparison (reference or new sheet not found)")
	}

	refData, err := ReadExcelData(refSheet) // Read data from the reference sheet.
	if err != nil {
		return fmt.Errorf("Failed to read reference sheet data: %v", err)
	}

	newData, err := ReadExcelData(newSheet) // Read data from the newly created sheet
	if err != nil {
		return fmt.Errorf("Failed to read new sheet data: %v", err)
	}

	filteredRefData := FilterData(refData, newData) // Filter the reference data

	diffCount := compareData(filteredRefData, newData) // Compare data from the two sheets
	log.Printf("Audit completed: %d differences found", diffCount)
	logger.Trace("Completed data comparison.", logger.Args("Differences found", diffCount))
	return nil
}

// CompareData evaluates differences between two slices of InterfaceData (reference data and new data).
func compareData(refData, newData []InterfaceData) int {
	nodeFiles := make(map[string]*os.File)
	statusSummary := make(map[string]map[string]int)        // A nested map to track status summaries per node.
	currentTime := time.Now().Format("02-01-2006 15:04:05") // DD-MM-YYYY HH:MM:SS

	// Prepare files and status summary for nodes found in newData
	for _, d := range newData {
		if _, exists := nodeFiles[d.Node]; !exists {
			reportFile := fmt.Sprintf("audit_report_%s_%s.txt", d.Node, currentTime)
			file, err := os.Create(reportFile)
			if err != nil {
				log.Fatalf("Failed to create report file for node %s: %v", d.Node, err)
			}
			defer file.Close()
			nodeFiles[d.Node] = file
			statusSummary[d.Node] = make(map[string]int) // Initialise status count map for this node

			// Write the initial part of the report
			_, _ = file.WriteString(fmt.Sprintf("Audit Report for %s generated on: %s\n", d.Node, currentTime))
		}
		statusSummary[d.Node][d.Status]++ // Increment count for this status
	}

	// Map reference data for comparison
	refMap := make(map[string]InterfaceData)
	for _, d := range refData {
		key := fmt.Sprintf("%s-%s", d.Node, d.Interface)
		refMap[key] = d
	}

	// Compare new data against reference data and write differences
	diffCount := 0
	for _, d := range newData {
		key := fmt.Sprintf("%s-%s", d.Node, d.Interface)
		ref, exists := refMap[key]
		file, fileExists := nodeFiles[d.Node]

		if exists && fileExists {
			if !dataEquals(ref, d) {
				diffCount++
				diff := fmt.Sprintf("Difference found for Node: %s, Interface: %s\n", d.Node, d.Interface)
				diff += fmt.Sprintf("Reference: %+v\nNew Data: %+v\n", ref, d)
				diff += "-----------------------------------\n"
				_, _ = file.WriteString(diff)
			}
		} else if !exists && fileExists {
			newEntry := fmt.Sprintf("New entry detected for Node: %s, Interface: %s\n", d.Node, d.Interface)
			newEntry += "-----------------------------------\n"
			_, _ = file.WriteString(newEntry)
		}
	}

	// Write the status summary at the top of each report file and close files
	for node, file := range nodeFiles {
		summaryInfo := "\nStatus Summary:\n"
		for status, count := range statusSummary[node] {
			summaryInfo += fmt.Sprintf("%s: %d\n", status, count)
		}
		summaryInfo += "===================================\n"
		file.Seek(0, 0) // Go back to the beginning of the file
		_, _ = file.WriteString(summaryInfo)
		file.Close()
		log.Printf("Differences report for %s saved to 'audit_report_%s.txt'", node, node)
	}

	return diffCount
}

// Check if two InterfaceData objects are identical.
func dataEquals(a, b InterfaceData) bool {
	return a.Type == b.Type && a.Description == b.Description && a.Status == b.Status &&
		a.Speed == b.Speed && a.Duplex == b.Duplex && a.VLAN == b.VLAN && a.Port == b.Port && a.Slot == b.Slot
}
