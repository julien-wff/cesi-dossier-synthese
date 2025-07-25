package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"time"
)

type TimingKey struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TelemetryStats struct {
	TotalParses             int                `json:"totalParses"`
	TotalParsesOverLastWeek int                `json:"totalParsesOverLastWeek"`
	UniqueUsers             int                `json:"uniqueUsers"`
	UniqueUsersOverLastWeek int                `json:"uniqueUsersOverLastWeek"`
	ErrorRate               float64            `json:"errorRate"`
	ErrorsOverLastWeek      int                `json:"errorsOverLastWeek"`
	AveragePdfSizeKb        float64            `json:"averagePdfSizeKb"`
	MaxPdfSizeKb            int64              `json:"maxPdfSizeKb"`
	AverageParseTime        float64            `json:"averageParseTime"`
	AverageParseTime95th    float64            `json:"averageParseTime95th"`
	UAOSs                   map[string]int64   `json:"UAOSs"`
	UABrowsers              map[string]int64   `json:"UABrowsers"`
	UAPlatforms             map[string]int64   `json:"UAPlatforms"`
	LatestSuccessfulParses  []parseTelemetry   `json:"latestSuccessfulParses"`
	LatestFailedParses      []parseTelemetry   `json:"latestFailedParses"`
	TimingAverageDuration   map[string]float64 `json:"timingAverageDuration"`
	TimingKeys              []TimingKey        `json:"timingKeys"`
}

// ReadTelemetry reads the raw content of the log file and returns it as a slice of parseTelemetry structs
// It locks the telemetryMutex to ensure thread safety while reading the file.
func ReadTelemetry() (*[]parseTelemetry, error) {
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

	// Read file content line by line
	scanner := bufio.NewScanner(file)
	var telemetry []parseTelemetry
	for scanner.Scan() {
		var entry parseTelemetry
		line := scanner.Text()
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			return nil, fmt.Errorf("failed to unmarshal telemetry entry: %w", err)
		}
		telemetry = append(telemetry, entry)
	}

	return &telemetry, nil
}

// ComputeTelemetryStats computes statistics from the telemetry data
func ComputeTelemetryStats(telemetry *[]parseTelemetry) TelemetryStats {
	if telemetry == nil || len(*telemetry) == 0 {
		return TelemetryStats{}
	}

	data := *telemetry
	oneWeekAgo := time.Now().AddDate(0, 0, -7)

	// Initialize counters and accumulators
	totalRecords := len(data)
	lastWeekRecords := 0
	totalErrors := 0
	errorsOverLastWeek := 0
	var totalPdfSizeKb int64
	var maxPdfSizeKb int64
	var totalParseDuration float64
	successfulParses := 0

	// Use maps for efficient unique user tracking
	uniqueIPs := make(map[string]struct{})
	uniqueIPsLastWeek := make(map[string]struct{})

	// Maps for user agent statistics
	uaOSs := make(map[string]int64)
	uaBrowsers := make(map[string]int64)
	uaPlatforms := make(map[string]int64)

	// Timing statistics - accumulate durations by timing name
	timingDurations := make(map[string][]float64)
	var mostCompleteTimings []TimingKey
	maxSeenTimingsCount := 0

	// Collect parse durations for percentile calculation
	var parseDurations []float64

	// Single pass through data for efficiency
	for _, entry := range data {
		// Track unique users
		uniqueIPs[entry.ClientIP] = struct{}{}

		// Parse timestamp string to time.Time
		timestamp, err := time.Parse(time.DateTime, entry.Timestamp)
		isRecentEntry := false
		if err == nil {
			isRecentEntry = timestamp.After(oneWeekAgo)
		}

		if isRecentEntry {
			lastWeekRecords++
			uniqueIPsLastWeek[entry.ClientIP] = struct{}{}
		}

		// Track errors
		if !entry.Success {
			totalErrors++
			if isRecentEntry {
				errorsOverLastWeek++
			}
		} else {
			// Calculate total duration for successful parses
			var entryDuration float64
			for _, timing := range entry.Timings {
				entryDuration += timing.Duration
				// Accumulate timing durations by name
				timingDurations[timing.Name] = append(timingDurations[timing.Name], timing.Duration)
			}
			totalParseDuration += entryDuration
			parseDurations = append(parseDurations, entryDuration)
			successfulParses++

			// Track the entry with the most complete timings
			if len(entry.Timings) > maxSeenTimingsCount {
				maxSeenTimingsCount = len(entry.Timings)
				mostCompleteTimings = make([]TimingKey, 0, len(entry.Timings))
				for _, timing := range entry.Timings {
					mostCompleteTimings = append(mostCompleteTimings, TimingKey{
						Name:        timing.Name,
						Description: timing.Description,
					})
				}
			}
		}

		// Track PDF sizes
		totalPdfSizeKb += entry.ContentLengthKB
		if entry.ContentLengthKB > maxPdfSizeKb {
			maxPdfSizeKb = entry.ContentLengthKB
		}

		// Track user agent statistics
		if entry.UserAgent.OS != "" {
			uaOSs[entry.UserAgent.OS]++
		}
		if entry.UserAgent.Browser != "" {
			uaBrowsers[entry.UserAgent.Browser]++
		}
		if entry.UserAgent.Platform != "" {
			uaPlatforms[entry.UserAgent.Platform]++
		}
	}

	// Calculate timing averages
	timingAverageDuration := make(map[string]float64)
	for timingName, durations := range timingDurations {
		if len(durations) > 0 {
			var sum float64
			for _, duration := range durations {
				sum += duration
			}
			timingAverageDuration[timingName] = sum / float64(len(durations))
		}
	}

	// Calculate derived statistics
	var errorRate float64
	if totalRecords > 0 {
		errorRate = float64(math.Round(float64(totalErrors)/float64(totalRecords)*10000)) / 100
	}

	var averagePdfSizeKb float64
	if totalRecords > 0 {
		averagePdfSizeKb = float64(totalPdfSizeKb) / float64(totalRecords)
	}

	var averageParseTime float64
	if successfulParses > 0 {
		averageParseTime = totalParseDuration / float64(successfulParses)
	}

	// Calculate 95th percentile parse time
	var averageParseTime95th float64
	if len(parseDurations) > 0 {
		sort.Slice(parseDurations, func(i, j int) bool {
			return parseDurations[i] < parseDurations[j]
		})

		percentile95Index := int(math.Floor(float64(len(parseDurations)) * 0.95))
		if percentile95Index < len(parseDurations) {
			top5Percent := parseDurations[percentile95Index:]
			var sum float64
			for _, duration := range top5Percent {
				sum += duration
			}
			if len(top5Percent) > 0 {
				averageParseTime95th = sum / float64(len(top5Percent))
			}
		}
	}

	// Handle NaN values
	if math.IsNaN(errorRate) {
		errorRate = 0
	}
	if math.IsNaN(averagePdfSizeKb) {
		averagePdfSizeKb = 0
	}
	if math.IsNaN(averageParseTime) {
		averageParseTime = 0
	}
	if math.IsNaN(averageParseTime95th) {
		averageParseTime95th = 0
	}

	// Prepare latest successful and error parses
	const maxLatestParses = 15
	latestSuccessfulParses := make([]parseTelemetry, 0, maxLatestParses)
	latestFailedParses := make([]parseTelemetry, 0, maxLatestParses)
	for i := len(data) - 1; i >= 0 && (len(latestSuccessfulParses) < maxLatestParses || len(latestFailedParses) < maxLatestParses); i-- {
		entry := data[i]
		if entry.Success && len(latestSuccessfulParses) < maxLatestParses {
			latestSuccessfulParses = append(latestSuccessfulParses, entry)
		} else if !entry.Success && len(latestFailedParses) < maxLatestParses {
			latestFailedParses = append(latestFailedParses, entry)
		}
	}

	return TelemetryStats{
		TotalParses:             totalRecords,
		TotalParsesOverLastWeek: lastWeekRecords,
		UniqueUsers:             len(uniqueIPs),
		UniqueUsersOverLastWeek: len(uniqueIPsLastWeek),
		ErrorRate:               errorRate,
		ErrorsOverLastWeek:      errorsOverLastWeek,
		AveragePdfSizeKb:        averagePdfSizeKb,
		MaxPdfSizeKb:            maxPdfSizeKb,
		AverageParseTime:        averageParseTime,
		AverageParseTime95th:    averageParseTime95th,
		UAOSs:                   uaOSs,
		UABrowsers:              uaBrowsers,
		UAPlatforms:             uaPlatforms,
		LatestSuccessfulParses:  latestSuccessfulParses,
		LatestFailedParses:      latestFailedParses,
		TimingAverageDuration:   timingAverageDuration,
		TimingKeys:              mostCompleteTimings,
	}
}
