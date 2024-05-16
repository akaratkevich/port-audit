package internal

import "regexp"

var showInterfaceDescriptionIOSXR = regexp.MustCompile(
	`^(?P<Interface>(Te|GigabitEthernet|TenGigE|Eth|Ge|Ethernet)[\w/]+)\s+is\s+` +
		`(?P<Status>up|down|administratively down|admin down|connected|disabled)\s+` +
		`,\s+` +
		`(?P<Protocol>up|down)\s+` +
		`(?P<Description>\#\#.*?\#\#)\s*` +
		`.*$`,
)
