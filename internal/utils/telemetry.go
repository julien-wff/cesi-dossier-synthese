package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

const logFile = "./data/parser.log"

var mu sync.Mutex

type parseTelemetry struct {
	Success         bool            `json:"success"`
	Timestamp       string          `json:"timestamp"`
	ClientIP        string          `json:"clientIP"`
	Source          string          `json:"source"`
	ContentLengthKB int64           `json:"contentLengthKB"`
	Timings         []TimingElement `json:"timings"`
	Error           *string         `json:"error"`
}

// appendLog serializes the provided content to JSON and appends it to the log file
func appendLog(telemetry parseTelemetry) error {
	// Lock Mutex
	mu.Lock()
	defer mu.Unlock()

	// Create append stream to file
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("error while closing log file", err)
		}
	}(file)

	// Serialize telemetry to JSON
	serializedContent, err := json.Marshal(telemetry)
	if err != nil {
		return err
	}

	// Add line break
	serializedContent = append(serializedContent, '\n')

	// Write telemetry to file
	_, err = file.Write(serializedContent)
	if err != nil {
		return err
	}

	return nil
}

// LogParseTelemetry logs the provided request, timings and error to the log file
func LogParseTelemetry(req *http.Request, timings *ProcessTiming, error *APIError) error {
	var errorMessage *string
	if error != nil {
		errorMessage = &error.Code
	}

	var timingsElements []TimingElement
	if timings != nil {
		timingsElements = timings.Elements
	}

	clientIp, _, _ := net.SplitHostPort(req.RemoteAddr)

	return appendLog(parseTelemetry{
		Success:         error == nil,
		Timestamp:       time.Now().Format(time.UnixDate),
		ClientIP:        clientIp,
		Source:          req.URL.Path,
		ContentLengthKB: req.ContentLength / 1e3,
		Timings:         timingsElements,
		Error:           errorMessage,
	})
}

// ReadTelemetry reads the raw content of the log file and returns it as a byte slice
func ReadTelemetry() (*[]byte, error) {
	// Lock Mutex
	mu.Lock()
	defer mu.Unlock()

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
