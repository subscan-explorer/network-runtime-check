package utils

import "strings"

func ErrorReduction(err error) string {
	errStr := err.Error()
	if strings.Contains(errStr, "timeout") {
		return "timeout"
	}
	if strings.Contains(errStr, "cancel") {
		return "cancel"
	}
	return errStr
}
