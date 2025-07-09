package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/mileusna/useragent"
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
	UserAgent       userAgent       `json:"userAgent"`
}

type userAgent struct {
	OS       string `json:"os"`
	Browser  string `json:"browser"`
	Platform string `json:"platform"`
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

const unknownUserAgent = "Unknown"

// parseUserAgent parses the User-Agent header from the request and returns a userAgent struct
func parseUserAgent(rawUA string) userAgent {
	if rawUA == "" {
		return userAgent{
			OS:       unknownUserAgent,
			Browser:  unknownUserAgent,
			Platform: unknownUserAgent,
		}
	}

	ua := useragent.Parse(rawUA)

	platform := unknownUserAgent
	if ua.Desktop {
		platform = "Desktop"
	} else if ua.Mobile {
		platform = "Mobile"
	} else if ua.Tablet {
		platform = "Tablet"
	} else if ua.Bot {
		platform = "Bot"
	}

	if ua.OS == "" {
		ua.OS = unknownUserAgent
	}

	if ua.Name == "" {
		ua.Name = unknownUserAgent
	}

	return userAgent{
		OS:       ua.OS,
		Browser:  ua.Name,
		Platform: platform,
	}
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
	clientIpHasher := sha256.New()
	clientIpHasher.Write([]byte(clientIp))

	ua := parseUserAgent(req.UserAgent())

	return appendLog(parseTelemetry{
		Success:         error == nil,
		Timestamp:       time.Now().Format(time.DateTime),
		ClientIP:        hex.EncodeToString(clientIpHasher.Sum(nil)),
		Source:          req.URL.Path,
		ContentLengthKB: req.ContentLength / 1e3,
		Timings:         timingsElements,
		Error:           errorMessage,
		UserAgent:       ua,
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
