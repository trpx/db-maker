package utils

import "strings"

func FilterNonEmptyLines(lines []string) (nonEmptyLines []string) {
	for _, el := range lines {
		el = strings.TrimSpace(el)
		if el != "" {
			nonEmptyLines = append(nonEmptyLines, el)
		}
	}
	return nonEmptyLines
}
