package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// timingElement represents a single step in the process timing.
type timingElement struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Duration    float64 `json:"duration"` // In milliseconds
}

// ProcessTiming represents the duration of each step, with a list of timingElement.
type ProcessTiming struct {
	Elements      []timingElement `json:"elements"`
	lastTimestamp int64
}

// NewProcessTiming creates a new ProcessTiming object, initializing the timestamp to the current time.
func NewProcessTiming() *ProcessTiming {
	return &ProcessTiming{
		Elements:      []timingElement{},
		lastTimestamp: time.Now().UnixNano(),
	}
}

// AddElement adds a new entry to the ProcessTiming object.
// It calculates the duration between the last element and the current time.
// The timestamp is then reset to the current time.
func (pt *ProcessTiming) AddElement(name, description string) {
	currentTimestamp := time.Now().UnixNano()
	duration := currentTimestamp - pt.lastTimestamp
	pt.lastTimestamp = currentTimestamp
	pt.Elements = append(pt.Elements, timingElement{
		Name:        name,
		Description: description,
		Duration:    float64(duration) / 1e6,
	})
}

// String returns a string representation of the ProcessTiming object.
// It displays the duration of each step, as well as the total duration.
func (pt *ProcessTiming) String() string {
	result := ""
	sum := 0.0
	for _, element := range pt.Elements {
		result += element.Description + ": " + strconv.FormatFloat(element.Duration, 'f', -1, 64) + "ms\n"
		sum += element.Duration
	}
	result += "Total: " + strconv.FormatFloat(sum, 'f', 2, 64) + "ms"
	return result
}

// ServerTiming returns a string that can be used in the Server-Timing header,
// indicating the duration of each step of the process based on the ProcessTiming object.
//
// Format: <name>;desc="<description>";dur=<duration>
func (pt *ProcessTiming) ServerTiming() string {
	result := make([]string, 0)

	for _, element := range pt.Elements {
		result = append(result, fmt.Sprintf(
			"%s;desc=\"%s\";dur=%f",
			element.Name,
			element.Description,
			element.Duration,
		))
	}

	return strings.Join(result, ", ")
}
