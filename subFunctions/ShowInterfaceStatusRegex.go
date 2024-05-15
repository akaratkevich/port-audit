package subFunctions

import "regexp"

// var showIntStatus = regexp.MustCompile(
//
//	`^` + // Start of the line
//		`(?P<Interface>` + // Begin named group 'Interface'
//		`(GigabitEthernet|TenGigabitEthernet|Eth|Ge|Gi)` + // Match interface types
//		`[0-9/]+` + // Followed by one or more digits or slashes
//		`) ` + // End named group 'Interface'
//		`+` + // One or more spaces
//		`(?P<Description>.+?) ` + // Begin named group 'Description', non-greedy match of any character
//		`+` + // One or more spaces
//		`(?P<Status>up|down|administratively down|admin down|connected|disabled) ` + // Named group 'Status' with specific options
//		`+` + // One or more spaces
//		`(?P<VLAN>` + // Begin named group 'VLAN'
//		`[1-9]|[1-9]\d|1\d\d|[1-3]\d\d\d|40\d\d|409[0-4]|routed|trunked` + // Match VLAN ID 1-4094
//		`) ` + // End named group 'VLAN'
//		`+` + // One or more spaces
//		`(?P<Duplex>full|half|auto) ` + // Named group 'Duplex' matches 'full', 'half', or 'auto'
//		`+` + // One or more spaces
//		`(?P<Speed>\S+) ` + // Named group 'Speed' matches non-whitespace characters
//		`+` + // One or more spaces
//		`(?P<Type>\S+)` + // Named group 'Type' matches non-whitespace characters
//		`$`, // End of the line
//
// )
var showIntStatus = regexp.MustCompile(
	`^(?P<Interface>(GigabitEthernet|TenGigabitEthernet|Eth|Ge|Gi)\d+/\d+)\s+` + // Match interface names such as Gi1/1
		`(?P<Description>.*?)\s+` + // Non-greedy match for Name which might be empty
		`(?P<Status>up|down|administratively down|admin down|connected|disabled)\s+` + // Capture status
		`(?P<VLAN>\d+)\s+` + // VLAN number
		`(?P<Duplex>full|half|auto)\s+` + // Duplex setting
		`(?P<Speed>\S+)\s+` + // Speed, non-whitespace characters (this needs to include auto 10/100/1000BaseT as one field)
		`(?P<Type>.*)$`, // Type which captures until the end of the line
)
