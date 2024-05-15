package main

import (
	"fmt"
	"github.com/pterm/pterm"
	"go-port-audit-v4/subFunctions"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	startTime := time.Now()

	// Create a screen logger
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)

	logger.Trace("Staring the port-audit process...") // Log to the screen

	// Open a file for writing logs.
	logFile, err := os.OpenFile("port-audit-application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()

	// Set the output of logs to the file
	log.SetOutput(logFile)

	// ---- !!! FROM THIS POINT ON, ALL LOG MESSAGES WILL BE WRITTEN TO THE FILE !!! ----

	// 1. Setup and parse command-line arguments
	username, password, inventoryFile, baseFile, err := subFunctions.SetupFlags()
	if err != nil {
		logger.Fatal("Exiting the program due to setup failure", logger.Args("Reason", err)) // Log to the screen
		os.Exit(1)
	} else {
		logger.Trace("Successfully passed the parameters for setup.") // log to the screen
	}

	// 1.5 Setup menu
	// Define command options and their corresponding SSH commands
	options := []string{"NXOS/IOS - show interface status", "IOS - show interface description"}
	commands := map[string]string{
		"NXOS/IOS - show interface status": "show interface status",
		"IOS - show interface description": "show interface description",
	}

	// Interactive menu to select a command
	printer := pterm.DefaultInteractiveSelect.WithOptions(options)
	selectedOption, _ := printer.Show()

	selectedCommand, exists := commands[selectedOption]
	if !exists {
		logger.Fatal("No valid command selected")
		os.Exit(1)
	}
	logger.Info("Selected command:", logger.Args("Command", pterm.Green(selectedCommand)))

	// 2. Read the inventory file
	inventory, err := subFunctions.ReadInventory(*inventoryFile, logger)
	if err != nil {
		log.Printf("Error: Failed to read inventory: %v. Exiting the program due to inventory load failure.", err) // Log to the file
		logger.Fatal("Exiting the program due to inventory load failure.", logger.Args("Reason", err))             // Log to the screen
		os.Exit(1)
	}

	// 3. Setup concurrency
	logger.Trace("Initialising concurrency...") // Log to the screen
	dataChan := make(chan subFunctions.InterfaceData)
	var wg sync.WaitGroup

	// Set the number of workers
	numWorkers := 10
	workQueue := make(chan subFunctions.Device, len(inventory.Devices))

	logger.Trace("Launching worker goroutines for device processing...") // Log to the screen
	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for device := range workQueue {
				subFunctions.ProcessDevice(device, dataChan, *username, *password, selectedCommand)
			}
		}()
	}
	// 4. Distribute work among workers
	logger.Trace("Distributing tasks among workers...", logger.Args("Workers", numWorkers), logger.Args("Work Queue", len(inventory.Devices))) // Log to the screen

	for _, device := range inventory.Devices {
		workQueue <- device
	}
	close(workQueue)

	logger.Trace("Aggregating processed data...") // Log to the screen
	// Use another goroutine to read from the channel and collect data
	allData := make([]subFunctions.InterfaceData, 0)
	go func() {
		for data := range dataChan {
			allData = append(allData, data)
			log.Printf("Received data: %+v", data)
		}
	}()

	wg.Wait()       // Wait for all workers to finish processing
	close(dataChan) // Safely close the data channel

	logger.Trace("Device processing has finished, and data channels are closed. Data collected.", logger.Args("Total interfaces processed", len(allData))) // Log to the screen
	logger.Warn("Check the application log for details on connection issues:", logger.Args("Log File", "port-audit-application.log"))                      // Log to the screen

	log.Printf("All processing goroutines completed: data channels closed successfully, and data collected for %d interfaces.", len(allData)) // Log to file

	// 7. Perform Excel operations based on the command line option.
	logger.Trace("Initiating Excel and data comparison operations, and preparing final reports....") // Log to the screen
	subFunctions.ExcelOperations(allData, *baseFile, logger)

	// 8. Zip the files
	zipPath, err := subFunctions.ZipAndDeleteFiles("./", logger)
	if err != nil {
		logger.Fatal("Failed to zip and delete files: ", logger.Args("error", err))
		os.Exit(1)
	}
	//logger.Info("The Excel file is ready for review.", logger.Args("file", "PortAudit.xlsx")) // Log to the screen

	logger.Trace("Port-audit process completed.")

	// Create a map of interesting stuff.
	filesInfo := map[string]any{
		"Application Log":     "Contains all runtime logs and errors - 'port-audit-application.log'",
		"Excel Data File":     "Compiled interface data - 'PortAudit.xlsx'",
		"Differences Archive": fmt.Sprintf("Zipped reports detailing differences - '%s'", zipPath),
	}

	// Log the comprehensive review message using a formatted string from the map.
	logger.Info("Review the following generated files:", logger.ArgsFromMap(filesInfo))

	// Reporting
	totalNodes := len(inventory.Devices)
	elapsedTime := time.Since(startTime)
	fmt.Println("\n----------------------------------------------------------------")
	pterm.FgLightYellow.Printf("Execution completed for %d devices\n", totalNodes)
	pterm.FgLightYellow.Printf("Execution Time: %s\n", elapsedTime)
	fmt.Println("----------------------------------------------------------------")
}
