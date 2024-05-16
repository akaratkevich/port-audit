package internal

import "regexp"

var showIntRe = regexp.MustCompile(
	`^(?P<Interface>(Bundle-Ether|GigabitEthernet|TenGigE|Eth|Ge|Ethernet)[\w/]+)\s+is\s+` +
		`(?P<Status>up|down|administratively down|admin down|connected|disabled)` +
		`.*?` +
		`admin state is (?P<AdminState>up|down),` +
		`.*?` +
		`Description: (?P<Description>\#\#.*\#\#)\s+` +
		`MTU (?P<MTU>\d+) bytes, BW (?P<BW>\d+) Kbit` +
		`.*$`,
)
