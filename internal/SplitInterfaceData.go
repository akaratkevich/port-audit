package internal

import (
	"fmt"
	"regexp"
)

// Split interface into slot and port
func SplitInterfaceData(interfaceStr string) (slot string, port string, err error) {
	regex := regexp.MustCompile(`\D+(\d+)/(\d+)`)
	matches := regex.FindStringSubmatch(interfaceStr)
	if len(matches) < 3 {
		return "", "", fmt.Errorf("invalid interface format")
	}
	slot = matches[1]
	port = matches[2]
	return slot, port, nil
}
