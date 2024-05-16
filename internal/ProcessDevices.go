package internal

import "log"

/*
Manage the processing of a single network device within a goroutine.

Parameters:
  - device Device: A struct containing details about the device.
  - dataChan chan<- InterfaceData: A channel used to send processed interface data back to the main program.
                                   The channel is "send-only" within this function.
  - username string: The username required for SSH authentication.
  - password string: The password required for SSH authentication.

Description:
  - This function logs the beginning of the processing for a specific device.
  - It invokes 'ConnectAndExecute' to establish an SSH connection, execute a command relevant to the device's platform,
    and handle the output. Any occurring errors during connection or execution are logged.
  - After processing, it logs the completion of the operation for the device.
  - This function does not manage concurrency directly (e.g., it does not call 'wg.Done()'); instead, it is designed to
    be managed by a worker pool where concurrency control is handled externally.

Usage:
  - Intended to be run inside a worker goroutine as part of a pool that processes multiple devices concurrently.
  - The worker pool is responsible for managing the lifecycle of goroutines, including the synchronization of their completion.
*/

func ProcessDevice(device Device, dataChan chan<- InterfaceData, username, password string, selectedCommand string) {
	log.Printf("Starting processing for device: %s", device.Host)
	if err := ConnectAndExecute(device, username, password, dataChan, selectedCommand); err != nil {
		log.Printf("Failed to connect or execute on device %s: %v", device.Host, err)
	}
	log.Printf("Completed processing for device: %s", device.Host)
}
