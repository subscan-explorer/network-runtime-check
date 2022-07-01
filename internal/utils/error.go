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
	if strings.Contains(errStr, "deadline exceeded") {
		return "deadline"
	}
	if e := strings.Split(errStr, "err:"); len(e) != 0 {
		return strings.TrimSpace(e[len(e)-1])
	}
	return errStr
}
