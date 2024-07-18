package internal

var CiscoPortAuditUsageGuide = `
Cisco Port Audit
--------------------------------------
Ensure you have the Port Allocation Spreadsheet in the home directory (*unless you are creating a new spreadsheet)
- File name: PortAudit.xlsx
- Sheet name: Baseline (*used for comparison)

Follow these steps to use the tool:

Example: port-audit -u admin -p admin123 -f inventory.yml

1. Enter SSH username and password (only SSH is currently supported).
2. Provide the file path to the inventory file.
3. Select the command to run on the devices from the interactive menu (show interface description/show interface status).
4. The application will read the inventory file and execute the selected command on each device.
5. The results will be logged, and the Excel file will be updated/created.
6. Difference report files will be generated per device.

Note:
- The inventory file can be generated using --gen flag (Create YAML Inventory File).

Example of inventory file (YAML format):
--------------------------------------
devices:
  - host: r1
    port: 22
    platform: ios
    transport: ssh
  - host: r2
    port: 22
    platform: ios
    transport: ssh
--------------------------------------

Flags:
  -base
        Create initial Excel file with a baseline sheet
  -f string
        File path
  -gen
        Generate a YAML inventory file from a list of devices
  -p string
        Password for device access
  -u string
        Username for device access
  -usage
        Display the usage guide

`
