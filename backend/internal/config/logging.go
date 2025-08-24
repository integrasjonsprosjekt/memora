package config

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

type LogLevel int

const (
	LogLevelError LogLevel = iota
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
)

var logLevels = map[string]LogLevel{
	"error":   LogLevelError,
	"warn":    LogLevelWarn,
	"info":    LogLevelInfo,
	"debug":   LogLevelDebug,
}

var logLevelNames = slices.Collect(maps.Keys(logLevels))

// String returns the string representation of the log level
func (l LogLevel) String() (string, error) {
	for name, level := range logLevels {
		if level == l {
			return name, nil
		}
	}
	return "", fmt.Errorf("invalid log level: %d", l)
}

// ParseLogLevel converts a string to a LogLevel, returns error for invalid levels
func ParseLogLevel(level string) (LogLevel, error) {
	if logLevel, ok := logLevels[strings.ToLower(level)]; ok {
		return logLevel, nil
	}
	return LogLevelInfo, fmt.Errorf("invalid log level: %s (valid options: %s)", level, strings.Join(logLevelNames, ", "))
}
