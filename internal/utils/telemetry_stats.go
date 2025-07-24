package utils

import (
	"fmt"
	"io"
	"os"
)

// ReadTelemetry reads the raw content of the log file and returns it as a byte slice
func ReadTelemetry() (*[]byte, error) {
	// Lock Mutex
	telemetryMutex.Lock()
	defer telemetryMutex.Unlock()

	// Open log file
	file, err := os.Open(logFile)
	if err != nil {
		return nil, err
	}

	// Close file when done
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("error while closing log file", err)
		}
	}(file)

	// Read file content
	telemetry, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return &telemetry, nil
}
