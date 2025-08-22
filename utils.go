package config

import (
	"fmt"
	"strings"
)

// ParseBool parses a string to boolean with support for various formats
func ParseBool(s string) (bool, error) {
	switch strings.ToLower(s) {
	case "true", "t", "yes", "y", "1", "on":
		return true, nil
	case "false", "f", "no", "n", "0", "off", "":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean format")
	}
}
