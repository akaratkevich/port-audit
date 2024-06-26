package internal

import (
	"bufio"
	"log"
	"regexp"
	"strings"
)

// Define regex patterns for different platforms
var (
	RegexInterfaceStatus           = showIntStatus
	RegexInterfaceDescription      = showIntDescription
	RegexInterfaceDescriptionIOSXR = showInterfaceDescriptionIOSXR
)

/*
Scan the output string line by line to extract interface data based on the device's platform.
The extracted data is then sent to a channel for further processing.

Parameters:
  - output string: The raw command output from the device.
  - device Device: A struct that contains details about the device such as host and platform.
  - dataChan chan<- InterfaceData: A channel used to send processed interface data to other parts of the program.
*/

// ProcessOutput determines the platform of the device and parses the output accordingly.
func ProcessOutput(output string, command string, device Device, dataChan chan<- InterfaceData) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		var data *InterfaceData
		switch command {
		case "show interface status":
			data = parseInterfaceStatus(line, RegexInterfaceStatus, device) // Show interface status
		case "show interface description":
			data = parseInterfaceDescription(line, RegexInterfaceDescription, device) // Show interface description
		case "show int description":
			data = parseInterfaceDescriptionIOSXR(line, RegexInterfaceDescriptionIOSXR, device)
		}
		if data != nil {
			// ! DEBUGGING !
			// log.Printf("Sending data to channel for device %s: %+v", device.Host, *data)
			dataChan <- *data
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading command output: %v", err)
	}
}

/*
Use a regular expression to parse a single line of command output and extracts interface data.

Parameters:
  - line string: A single line of text to be parsed.
  - regex *regexp.Regexp: A compiled regular expression used to extract data from the line.
  - device Device: The device from which the output was obtained, used here to populate the Node field in the result.

Returns:
  - *InterfaceData: A pointer to an InterfaceData struct containing the parsed data, or nil if no data matches the regex.
*/

// parseInterfaceStatus parses a single line of output using the given regex.
func parseInterfaceStatus(line string, regex *regexp.Regexp, device Device) *InterfaceData {
	matches := regex.FindStringSubmatch(line)
	if matches == nil {
		return nil
	}

	paramMap := make(map[string]string)
	for i, name := range regex.SubexpNames() {
		if i > 0 && i <= len(matches) && name != "" {
			paramMap[name] = matches[i]
		}
	}
	slot, port := ParseSlotAndPort(paramMap["Interface"])

	return &InterfaceData{
		Node:        device.Host,
		Interface:   paramMap["Interface"],
		Slot:        slot,
		Port:        port,
		Description: paramMap["Description"],
		Status:      paramMap["Status"],
		VLAN:        paramMap["VLAN"],
		Duplex:      paramMap["Duplex"],
		Speed:       paramMap["Speed"],
		Type:        paramMap["Type"],
	}
}

// parseInterfaceDescription parses a single line of output from 'show interfaces description' for IOS devices.
func parseInterfaceDescription(line string, regex *regexp.Regexp, device Device) *InterfaceData {
	matches := regex.FindStringSubmatch(line)
	if matches == nil {
		return nil
	}

	paramMap := make(map[string]string)
	for i, name := range regex.SubexpNames() {
		if i > 0 && i <= len(matches) && name != "" {
			paramMap[name] = matches[i]
		}
	}

	slot, port := ParseSlotAndPort(paramMap["Interface"])
	// Concatenate the 'Status' and 'Protocol' fields for the 'Status' entry of InterfaceData.
	status := paramMap["Status"]
	if protocol := paramMap["Protocol"]; protocol != "" {
		status += " (" + protocol + ")" // Add protocol status in parentheses if it's non-empty.
	}

	// Populate only the available fields from the 'show interfaces description' output
	return &InterfaceData{
		Node:        device.Host,
		Interface:   paramMap["Interface"],
		Slot:        slot,
		Port:        port,
		Description: paramMap["Description"], // Optional, could be empty
		Status:      status,                  // Status + Protocol
		// Note: VLAN, Duplex, Speed, and Type are not available from 'show interfaces description'
	}
}

func parseInterfaceDescriptionIOSXR(line string, regex *regexp.Regexp, device Device) *InterfaceData {
	matches := regex.FindStringSubmatch(line)
	if matches == nil {
		return nil
	}

	paramMap := make(map[string]string)
	for i, name := range regex.SubexpNames() {
		if i > 0 && i <= len(matches) && name != "" {
			paramMap[name] = matches[i]
		}
	}

	status := paramMap["Status"]
	if protocol := paramMap["Protocol"]; protocol != "" {
		status += " (" + protocol + ")" // Add protocol status in parentheses if it's non-empty.
	}

	slot, port := ParseSlotAndPort(paramMap["Interface"])
	// Populate only the available fields from the 'show interfaces description' output
	return &InterfaceData{
		Node:        device.Host,
		Interface:   paramMap["Interface"],
		Slot:        slot,
		Port:        port,
		Description: paramMap["Description"], // This field may be optional and could be empty
		Status:      status,                  // Combines status and line protocol
	}
}
