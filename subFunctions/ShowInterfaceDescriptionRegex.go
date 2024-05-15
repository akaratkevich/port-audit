package subFunctions

import "regexp"

//var showIntDescription = regexp.MustCompile(
//	`^(?P<Interface>(GigabitEthernet|TenGigabitEthernet|Eth|Ge|Gi)\d+/\d+)\s+` + // Capture interface names more flexibly
//		`(?P<Status>[^\s]+)\s+` + // Capture the admin status (e.g., 'admin down', 'up', 'down')
//		`(?P<Protocol>[^\s]+)` + // Capture the protocol status (e.g., 'down', 'up')
//		`\s*(?P<Description>.*)$`, // Capture everything as description after protocol status
//)

var showIntDescription = regexp.MustCompile(`^(?P<Interface>\S+)\s+(?P<Status>admin down|down|up)\s+(?P<Protocol>down|up)\s*(?P<Description>.*)$`)
