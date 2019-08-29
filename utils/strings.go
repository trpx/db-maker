package utils

import "strings"

func FilterNonEmptyLines(lines []string) (nonEmptyLines []string) {
	for _, el := range lines {
		el = strings.TrimSpace(el)
		if el != "" {
			nonEmptyLines = append(nonEmptyLines, el)
		}
	}
	if len(nonEmptyLines) == 0 {
		Panicf("pass file '%s' is empty")
	}
	return nonEmptyLines
}
