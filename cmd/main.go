package main

import (
	"fmt"
	"github.com/pterm/pterm"
	"log"
	"os"
	"port-audit/internal"
	"sync"
	"time"
)

func main() {

	// Start the timer
	startTime := time.Now()

	// Create a screen logger
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)

	// Log to the screen starting of the application
	logger.Trace("Staring the port-audit process...\n")

	// Open a file for writing logs.
	logFile, err := os.OpenFile("port-audit-application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()

	// Set the output of logs to the filePath
	log.SetOutput(logFile)

	// FROM THIS POINT ON, ALL LOG MESSAGES WILL BE WRITTEN TO THE FILE

	// Setup and parse command-line arguments
	username, password, filePath, baseFile, generateInv, usageGuide, err := internal.SetupFlags()

	// Check if the usage flag is set and display the usage guide
	if *usageGuide {
		internal.PrintUsageGuide(internal.CiscoPortAuditUsageGuide)
		logger.Info("Displaying usage guide only, exiting the program.")
		log.Printf("Displaying usage guide only, exiting the program: %v", err)
		os.Exit(0)
	}

	// Check if the generate inventory flag is set and a file path is provided
	if *generateInv {
		if *filePath == "" {
			logger.Fatal("File path is required for generating inventory.")
			log.Printf("File path is required for generating inventory.")
			os.Exit(1)
		}
		// Call function to generate inventory file
		internal.GenerateInventory(*filePath, logger)
		return
	}

	// Check for required parameters
	if err != nil {
		logger.Fatal("Exiting the program due to setup failure", logger.Args("Reason", err)) // Log to the screen
		log.Printf("Exiting the program due to setup failure: %v", err)                      // Log to the filePath
		os.Exit(1)
	}

	if *username == "" {
		err = fmt.Errorf("username is required")
		logger.Fatal("Exiting the program due to setup failure", logger.Args("Reason", err)) // Log to the screen
		log.Printf("Exiting the program due to setup failure: %v", err)                      // Log to the filePath
		os.Exit(1)
	}

	if *password == "" {
		err = fmt.Errorf("password is required")
		logger.Fatal("Exiting the program due to setup failure", logger.Args("Reason", err)) // Log to the screen
		log.Printf("Exiting the program due to setup failure: %v", err)                      // Log to the filePath
		os.Exit(1)
	}

	logger.Trace("Successfully passed the parameters for setup.") // log to the screen
	log.Printf("Successfully passed the parameters for setup")    // Log to the filePath

	// Setup menu
	// Define command options and their corresponding SSH commands
	options := []string{"NXOS/IOS - show interface status", "IOS - show interface description", "IOSXR - show interface description"}
	commands := map[string]string{
		"NXOS/IOS - show interface status":   "show interface status",
		"IOS - show interface description":   "show interface description",
		"IOSXR - show interface description": "show int description",
	}

	// Interactive menu to select a command
	printer := pterm.DefaultInteractiveSelect.WithOptions(options)
	selectedOption, _ := printer.Show()

	selectedCommand, exists := commands[selectedOption]
	if !exists {
		logger.Fatal("No valid command selected")
		log.Printf("No valid command selected: %v", err) // Log to the filePath
		os.Exit(1)
	}
	logger.Info("Selected command:", logger.Args("Command", pterm.Green(selectedCommand)))

	// Read the inventory file
	inventory, err := internal.ReadInventory(*filePath, logger)
	if err != nil {
		log.Printf("Error: Failed to read inventory: %v. Exiting the program due to inventory load failure.", err) // Log to the filePath
		logger.Fatal("Exiting the program due to inventory load failure.", logger.Args("Reason", err))             // Log to the screen
		os.Exit(1)
	}

	// Setup concurrency
	logger.Trace("Initialising concurrency...") // Log to the screen
	log.Printf("Initialising concurrency...")   // Log to the filePath
	dataChan := make(chan internal.InterfaceData)
	var wg sync.WaitGroup

	// Counters for success and failure
	var successCounter int
	var failureCounter int
	var mu sync.Mutex

	// Set the number of workers
	numWorkers := 10
	workQueue := make(chan internal.Device, len(inventory.Devices))

	logger.Trace("Launching worker goroutines for device processing...") // Log to the screen
	log.Printf("Launching worker goroutines for device processing...")   // Log to the filePath

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for device := range workQueue {
				internal.ProcessDevice(device, dataChan, *username, *password, selectedCommand, &successCounter, &failureCounter, &mu)
			}
		}()
	}
	// Distribute work among workers
	logger.Trace("Distributing tasks among workers...", logger.Args("Workers", numWorkers), logger.Args("Work Queue", len(inventory.Devices))) // Log to the screen

	for _, device := range inventory.Devices {
		workQueue <- device
	}
	close(workQueue)

	logger.Trace("Aggregating processed data...") // Log to the screen
	log.Printf("Aggregating processed data...")   // Log to the filePath

	// Use another goroutine to read from the channel and collect data
	allData := make([]internal.InterfaceData, 0)
	go func() {
		for data := range dataChan {
			allData = append(allData, data)
			log.Printf("Received data: %+v", data)
		}
	}()

	wg.Wait()       // Wait for all workers to finish processing
	close(dataChan) // Safely close the data channel

	// Check if at least some devices were processed
	if failureCounter > 0 {
		logger.Warn("Some devices encountered connection issues. Please check the application log for detailed error messages.", logger.Args("Total failed connections", failureCounter))
		log.Printf("Some devices encountered connection issues. Total number of failed connections: %d", failureCounter)
	}
	// Check if data collection was successful
	if len(allData) <= 0 {
		log.Printf("Data collection failed: No interface data collected.")
		logger.Error("Data collection failed: No interface data collected.")
	} else {
		log.Printf("All processing goroutines completed: Data channels closed successfully, and data collected for %d interfaces.", len(allData))
		logger.Info("Data collection successful.", logger.Args("Total interfaces collected", len(allData)))
	}

	// Perform Excel operations based on the command line option.
	logger.Trace("Initiating Excel and data comparison operations, and preparing final reports....") // Log to the screen
	log.Printf("Initiating Excel and data comparison operations, and preparing final reports...")    // Log to file
	internal.ExcelOperations(allData, *baseFile, logger)

	// 8. Zip the files
	zipPath, err := internal.ZipAndDeleteFiles("./", logger)
	if err != nil {
		logger.Fatal("Failed to zip and delete files: ", logger.Args("error", err))
		os.Exit(1)
	}
	//logger.Info("The Excel filePath is ready for review.", logger.Args("filePath", "PortAudit.xlsx")) // Log to the screen

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
	pterm.FgLightYellow.Printf("Total %d devices\n", totalNodes)
	pterm.FgLightYellow.Printf("Successful connections: %d\n", successCounter)
	pterm.FgLightYellow.Printf("Failed connections: %d\n", failureCounter)
	pterm.FgLightYellow.Printf("Execution Time: %s\n", elapsedTime)
	fmt.Println("----------------------------------------------------------------")
}
