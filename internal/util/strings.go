package util

import "strings"

func SplitLinesNonEmpty(text, delimeter string) []string {
	tokens := strings.Split(text, delimeter)
	return RemoveEmptyStrings(tokens)
}

func RemoveEmptyStrings(input []string) []string {
	var output []string
	for _, v := range input {
		if v != "" {
			output = append(output, v)
		}
	}
	return output
}
