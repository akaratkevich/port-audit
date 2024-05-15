package subFunctions

import (
	"flag"
	"fmt"
	"os"
)

// Return pointers to the values of the flags.
func SetupFlags() (username *string, password *string, inventoryFile *string, baseFile *bool, err error) {
	// Define flags
	username = flag.String("u", "", "Username for device access")
	password = flag.String("p", "", "Password for device access")
	inventoryFile = flag.String("f", "", "Inventory file path")
	baseFile = flag.Bool("base", false, "Create initial Excel file with a baseline sheet")

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	// Parse the command line flags
	flag.Parse()

	// Validate the input flags
	err = validateFlags(username, password, inventoryFile)
	return

}

// Check if necessary flags are provided; return an error if any are missing.
func validateFlags(username, password, inventoryFile *string) error {
	// Validate required flags
	if *username == "" {
		return fmt.Errorf("error: Username is required. Please provide a username with --u (eg. --u admin)")
	}
	if *password == "" {
		return fmt.Errorf("error: Password is required. Please provide a password with --p (eg. --p password)")
	}
	if *inventoryFile == "" {
		return fmt.Errorf("error: Inventory file is required. Please provide a file with --f (eg. --f ./Inventory.yml)")
	}

	return nil
}
