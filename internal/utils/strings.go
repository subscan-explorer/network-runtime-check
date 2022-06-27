package utils

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
