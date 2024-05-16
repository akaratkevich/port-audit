package internal

import (
	"bufio"
	"fmt"
	"github.com/pterm/pterm"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

// Read device configuration from a specified file and generate a YAML inventory file.
func GenerateInventory(filePath string, logger *pterm.Logger) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	var devices []Device
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		device := Device{
			Host:      getPart(parts, 0),
			Port:      getPart(parts, 1),
			Platform:  getPart(parts, 2),
			Transport: getPart(parts, 3),
		}
		devices = append(devices, device)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	if err := generateYAMLFromDevices(devices); err != nil {
		return fmt.Errorf("failed to generate YAML: %w", err)
	}

	logger.Info("Inventory file created successfully.", logger.Args("File", "inventory.yml"))
	return nil
}

// Returns the string at the index from parts slice, or a default value if the index is out of bounds.
// Specifically, it returns "22" as a default for the port if not provided and "ssh" as a default for the transport
func getPart(parts []string, index int) string {
	if index < len(parts) && parts[index] != "" {
		return parts[index]
	}
	// Check if the missing part is the port and provide a default value
	if index == 1 {
		return "22"
	}
	if index == 3 {
		return "ssh"
	}
	return ""
}

// Marshals a slice of Device structs into a YAML format and write it to a file.
func generateYAMLFromDevices(devices []Device) error {
	data, err := yaml.Marshal(struct {
		Devices []Device `yaml:"devices"`
	}{Devices: devices})
	if err != nil {
		return fmt.Errorf("failed to marshal devices into YAML: %w", err)
	}

	filePath := "inventory.yml"
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write to file %s: %w", filePath, err)
	}

	return nil
}
