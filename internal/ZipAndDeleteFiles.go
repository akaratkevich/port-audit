package internal

import (
	"archive/zip"
	"fmt"
	"github.com/pterm/pterm"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Create a zip archive containing all files with "audit_report" in their name located in the working directory,
// then delete the original files after successful zipping.
func ZipAndDeleteFiles(srcDir string, logger *pterm.Logger) (string, error) {
	date := time.Now().Format("02-01-2006")           // Current date
	zipFileName := fmt.Sprintf("report_%s.zip", date) // Name of the zip file
	zipFilePath := filepath.Join(srcDir, zipFileName) // Full path to the new zip file

	// Create the zip file
	newZipFile, err := os.Create(zipFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create zip file: %v", err)
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile) // Initialise a zip writer
	defer zipWriter.Close()

	filesToDelete := []string{} // Track files that need to be deleted after successful zipping

	// Walk the directory and add files to the zip
	err = filepath.Walk(srcDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err // Propagate errors from walking the directory
		}

		if !info.IsDir() && strings.Contains(info.Name(), "audit_report") { // Check for files that include "audit_report"
			fileToZip, err := os.Open(filePath) // Open each file matching the criteria
			if err != nil {
				return fmt.Errorf("failed to open file %s: %v", filePath, err)
			}
			defer fileToZip.Close()

			relPath, err := filepath.Rel(srcDir, filePath) // Determine the relative path to preserve directory structure in the zip
			if err != nil {
				return fmt.Errorf("failed to calculate relative file path: %v", err)
			}

			zipEntry, err := zipWriter.Create(relPath) // Create a new entry in the zip file
			if err != nil {
				return fmt.Errorf("failed to create zip entry for file %s: %v", filePath, err)
			}

			if _, err := io.Copy(zipEntry, fileToZip); err != nil { // Copy file content into the zip entry
				return fmt.Errorf("failed to write file %s to zip: %v", filePath, err)
			}

			filesToDelete = append(filesToDelete, filePath) // Add file path to the deletion list after successful zipping
		}
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to add files to zip: %v", err)
	}

	// Delete the original files after successful zipping
	for _, file := range filesToDelete {
		if err := os.Remove(file); err != nil {
			log.Printf("Failed to delete file %s: %v", file, err)
		} else {
			log.Printf("Deleted file: %s", file)
		}
	}

	log.Printf("Successfully created and cleaned up zip archive: %s", zipFilePath)
	//logger.Info("Difference reports have been compiled, archived and available for download.", logger.Args("file", zipFilePath)) // Log to the screen
	return zipFilePath, nil
}
