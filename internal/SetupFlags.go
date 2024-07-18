package internal

import (
	"flag"
	"fmt"
	"os"
)

// SetupFlags parses the command-line flags and returns their values.
func SetupFlags() (username, password, filePath *string, baseFile, generateInv, usageGuide *bool, err error) {
	// Define flags
	usageGuide = flag.Bool("usage", false, "Display the usage guide")
	username = flag.String("u", "", "Username for device access")
	password = flag.String("p", "", "Password for device access")
	filePath = flag.String("f", "", "File path")
	baseFile = flag.Bool("base", false, "Create initial Excel file with a baseline sheet")
	generateInv = flag.Bool("gen", false, "Generate a YAML inventory file from a list of devices")

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	// Parse the command-line flags
	flag.Parse()

	// Return if only the usage guide flag is set
	if *usageGuide {
		return
	}

	// If the generate inventory flag is set and an inventory file is provided, return
	if *generateInv && *filePath != "" {
		return
	}

	// Validate the input flags for all other cases
	err = validateFlags(username, password, filePath)
	return
}

// validateFlags checks if necessary flags are provided and returns an error if any are missing.
func validateFlags(username, password, filePath *string) error {
	// Validate required flags
	if *username == "" {
		return fmt.Errorf("error: Username is required. Please provide a username with --u (e.g., --u admin)")
	}
	if *password == "" {
		return fmt.Errorf("error: Password is required. Please provide a password with --p (e.g., --p password)")
	}
	if *filePath == "" {
		return fmt.Errorf("error: Inventory file is required. Please provide a file with --f (e.g., --f ./Inventory.yml)")
	}

	return nil
}
