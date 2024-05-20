package internal

import (
	"strings"
)

// Split interface into slot and port
func ParseSlotAndPort(interfaceName string) (string, string) {
	parts := strings.Split(interfaceName, "/")
	if len(parts) > 1 {
		return parts[0][2:], parts[1] // Assumes interface names start with "Gi", "Te", etc.
	}
	return "", "" // Return empty strings if not parsable
}
