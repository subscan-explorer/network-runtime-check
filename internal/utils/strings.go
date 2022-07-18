package utils

import (
	"encoding/hex"
	"strings"
)

func MaxLenArrString(arr []string) int {
	var maxLen = 0
	for _, s := range arr {
		l := len(s)
		if maxLen < l {
			maxLen = l
		}
	}
	return maxLen
}

func HexToBytes(s string) []byte {
	s = strings.TrimPrefix(s, "0x")
	c := make([]byte, hex.DecodedLen(len(s)))
	_, _ = hex.Decode(c, []byte(s))
	return c
}

func Reverse[T any](arr []T) {
	i := 0
	j := len(arr) - 1
	for i < j {
		arr[i], arr[j] = arr[j], arr[i]
		i++
		j--
	}
}

func SliceEqual[T comparable](a []T, b []T) bool {
	var maps = make(map[T]int)
	for _, t := range a {
		maps[t]++
	}
	for _, t := range b {
		maps[t]--
	}
	for _, v := range maps {
		if v != 0 {
			return false
		}
	}
	return true
}
