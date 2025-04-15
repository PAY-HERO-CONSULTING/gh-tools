package utils

import "strings"

func StringSliceContains(list []string, val string) bool {

	_, exists := stringsSliceContains(list, val, true)

	return exists
}

func StringSliceContainsCaseInsensitive(list []string, val string) bool {

	_, exists := stringsSliceContains(list, val, false)

	return exists
}

func StringSliceIndex(list []string, val string) int {

	index, _ := stringsSliceContains(list, val, true)

	return index
}

func StringSliceIndexCaseInsensitive(list []string, val string) int {

	index, _ := stringsSliceContains(list, val, false)

	return index
}

func stringsSliceContains(list []string, val string, caseSensitivity bool) (int, bool) {

	for index, listValue := range list {

		if !caseSensitivity {
			listValue = strings.ToLower(listValue)
			val = strings.ToLower(val)
		}

		if listValue == val {
			return index, true
		}
	}

	return -1, false
}
