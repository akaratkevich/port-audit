# port-audit
![Static Badge](https://img.shields.io/badge/Project-IN_PROGRESS-orange) 
![Static Badge](https://img.shields.io/badge/Go-blue) 

| **⚠️ WARNING: This project is still in progress** ⚠️ |

Port-Audit is a tool designed for network administrators and engineers to automate the auditing of network devices. This application leverages SSH for remote execution of network commands, enabling the collection of interface data directly from the network devices.

Capabilities:

Command Execution: Currently, the application supports the following SSH commands for data collection:
- ![Static Badge](https://img.shields.io/badge/COMPLETED-green) [NXOS/IOS - show interface status]
- ![Static Badge](https://img.shields.io/badge/COMPLETED-green) [IOS - show interface description]

Platform Specificity:

The method of screen scraping is employed to parse these outputs, which necessitates distinct handling between different operating systems due to format variations.
At this stage the application is capable of being run against IOS/NXOS devices over SSH to gather the output from the following commands:

                                "NXOS/IOS - show interface status"
                                "IOS - show interface description"

Baseline Comparison:

Users can create a new baseline or upload an existing one to compare the current state of network interfaces against previouslky documented state.

Output and Reporting:

Excel Reporting:
The application generates an Excel spreadsheet summarising the data collected from network devices.
Difference Reports: 
Textual difference reports are produced for each node, detailing deviations from the baseline.
Archiving: 
Text reports are automatically zipped and prepared for download, facilitating easy distribution and review.

Usage Guide:

Mandatory Flags: Execution requires specific flags for operation:
-u: Username for SSH authentication.
-p: Password for SSH authentication.
-f: Path to the YAML file containing the inventory of devices to audit.

Optional Flag:
-base: Indicates whether to create a baseline sheet in the Excel report. This is useful for establishing a reference point for future audits.

Future Development:

Plans are underway to expand support to additional platforms and commands, enhancing the tool's versatility and adaptability to different network environments.
Further discussions are also expected to better integrate the baseline management with broader network management processes, ensuring seamless operational workflows.
Port-Audit aims to streamline network device management, making it more efficient and error-resistant by automating routine checks and documentation processes.
