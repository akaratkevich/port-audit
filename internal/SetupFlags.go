package internal

import (
	"flag"
	"fmt"
	"os"
)

// Return pointers to the values of the flags.
func SetupFlags() (username *string, password *string, filePath *string, baseFile *bool, generateInv *bool, err error) {
	// Define flags
	username = flag.String("u", "", "Username for device access")
	password = flag.String("p", "", "Password for device access")
	filePath = flag.String("f", "", "File path")
	baseFile = flag.Bool("base", false, "Create initial Excel filePath with a baseline sheet")
	generateInv = flag.Bool("gen", false, "Generate a yaml inventory filePath from a list of devices")

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	// Parse the command line flags
	flag.Parse()

	// If the generate inventory flag is set and an inventory file is provided,
	// return immediately to allow for inventory file generation.
	if *generateInv && *filePath != "" {
		// No further validation is needed since we're just generating a file.
		return
	}

	// Validate the input flags for all other cases
	err = validateFlags(username, password, filePath)
	return
}

// Check if necessary flags are provided; return an error if any are missing.
func validateFlags(username, password, file *string) error {
	// Validate required flags
	if *username == "" {
		return fmt.Errorf("error: Username is required. Please provide a username with --u (eg. --u admin)")
	}
	if *password == "" {
		return fmt.Errorf("error: Password is required. Please provide a password with --p (eg. --p password)")
	}
	if *file == "" {
		return fmt.Errorf("error: Inventory file is required. Please provide a file with --f (eg. --f ./Inventory.yml)")
	}

	return nil
}
