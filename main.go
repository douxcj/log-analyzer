package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// Constants for our "Service Level Indicators"
const (
	logFile        = "access.log"
	errorThreshold = 3
	alertWindow    = 10 * time.Second
)

func main() {
	// 1. Ensure the log file exists
	ensureLogFileExists(logFile)

	fmt.Printf("🚀 Log Monitor started on %s\n", logFile)
	fmt.Printf("Threshold: %d errors | Window: %s\n", errorThreshold, alertWindow)

	// 2. Start the tailing process
	tailAndAnalyze(logFile)
}

func tailAndAnalyze(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		// there is error
		log.Fatalf("Critical Failure: Could not open file: %v", err) // log.Fatalf terminates program
	}

	defer file.Close() //Go guarantees that file.Close() will run when the tailAndAnalyze function finishes, regardless of how it finishes.

	// Move to the end of the file so we only see NEW logs
	// 0 is the offset, 2 means "relative to the end" (io.SeekEnd).
	file.Seek(0, 2)
	// Create a buffered reader for efficient line-by-line reading.
	reader := bufio.NewReader(file)

	errorCount := 0
	// Timer to reset error count (simulating a sliding window)
	// NewTicker creates a channel that delivers a "tick" every 10 seconds (alertWindow).
	ticker := time.NewTicker(alertWindow)

	for { // infinite loop
		select { // = switch for channels. It allows a goroutine to wait on multiple communication operations
		// 1. Check if the ticker fired (event-driven)
		// ticker.C receives the current time (time.Time) every time the duration passes.
		case <-ticker.C:
			// Reset error count every 10 seconds
			if errorCount > 0 {
				fmt.Printf("--- Window Reset: Cleared %d errors ---\n", errorCount)
				errorCount = 0
			}
		// 2. If ticker didn't fire, run this immediately (non-blocking)
		default:
			// ReadString reads only up to the next newline ('\n').
			// variable 'line' holds just ONE log entry, not the whole 10s chunk.
			// The reader maintains its position, so the next call reads the NEXT line.
			line, err := reader.ReadString('\n')
			if err != nil {
				// No new line yet, wait briefly
				time.Sleep(500 * time.Millisecond)
				continue
			}

			// 3. Logic: Check for errors
			cleanLine := strings.ToUpper(strings.TrimSpace(line))
			if strings.Contains(cleanLine, "ERROR") || strings.Contains(cleanLine, "500") {
				errorCount++
				fmt.Printf("⚠️  Detected Error: %s (Total in window: %d)\n", cleanLine, errorCount)

				if errorCount >= errorThreshold {
					triggerAlert(errorCount)
					errorCount = 0 // Reset after alerting
				}
			}
		}
	}
}

func triggerAlert(count int) {
	fmt.Println("\n#############################################")
	fmt.Printf("🚨 ALERT: High Error Rate Detected!\n")
	fmt.Printf("Observed %d errors in the last window.\n", count)
	fmt.Println("#############################################\n")
}

func ensureLogFileExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.WriteFile(path, []byte(""), 0644) // owner has full read/write access, but everyone else can only read it. standard permission setting for text files like logs.
	}
}
