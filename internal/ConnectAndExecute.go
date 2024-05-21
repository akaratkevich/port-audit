package internal

import (
	"log"
	"strconv"
)

/*
Establish an SSH connection to a device and execute a specified command,
then send the processed output to a data channel.

Parameters:
  device Device - The device structure containing the necessary details like Host and Port.
  username string - SSH username for authentication.
  password string - SSH password for authentication.
  dataChan chan<- InterfaceData - A channel to send processed interface data to.

Returns:
  error - Returns an error if any step in the process fails
*/

func ConnectAndExecute(device Device, username, password string, dataChan chan<- InterfaceData, selectedCommand string) error {
	port, err := strconv.Atoi(device.Port)
	if err != nil {
		log.Printf("Error: Invalid port number for host %s: %v", device.Host, err)
		return err
	}

	session, err := InitialiseConnection(device.Host, port, username, password)
	if err != nil {
		log.Printf("Error: SSH connection failed for %s; error: %v", device.Host, err)
		return err
	}
	defer session.Close()
	log.Printf("SSH connection established for %s", device.Host)

	command := selectedCommand
	log.Printf("Executing command on %s: %s", device.Host, command)

	output, err := session.CombinedOutput(command)
	if err != nil {
		log.Printf("Error: Failed to execute command on %s: %v", device.Host, err)
		return err
	}
	log.Printf("Command executed successfully on %s, processing output...", device.Host)

	ProcessOutput(string(output), command, device, dataChan)
	log.Printf("Output processed for %s", device.Host)
	return nil
}
