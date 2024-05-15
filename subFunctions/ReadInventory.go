package subFunctions

import (
	"fmt"
	"github.com/pterm/pterm"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

// Device struct
type Device struct {
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	Platform  string `yaml:"platform"`
	Transport string `yaml:"transport"`
}

type Inventory struct {
	Devices []Device `yaml:"devices"`
}

/*
Read a YAML file containing device inventory information and unmarshal it into an Inventory struct.

Parameters:
  - filename string: The path to the YAML file that contains the inventory data.

Returns:
  - *Inventory: A pointer to the Inventory struct that holds all the parsed data from the YAML file.
  - error: An error object that indicates if the reading or unmarshalling process failed.
*/

func ReadInventory(filename string, logger *pterm.Logger) (*Inventory, error) {
	log.Println("Reading inventory file...") // log to the file
	var inventory Inventory
	data, err := os.ReadFile(filename) // Read the YAML file
	if err != nil {
		return nil, fmt.Errorf("failed to read inventory file: %v", err)
	}

	log.Println("Unmarshalling inventory data...")
	err = yaml.Unmarshal(data, &inventory) // Parse the YAML data into the Inventory struct
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal inventory data: %v", err)
	}
	log.Printf("Successfully loaded inventory: %d devices ready for processing.", len(inventory.Devices))               // log to the file
	logger.Trace("Inventory loaded: ready for device processing.", logger.Args("Device count", len(inventory.Devices))) // log to the screen
	return &inventory, nil
}
