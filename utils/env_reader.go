package utils

import (
	"os"
	"strings"
)

func ReadPasswordsFromEnv() (adminPasswords []string, userPasswords []string) {
	adminPasswords = ReadNonEmptyLinesFromEnv("DB_ADMIN_PASSWORDS")
	userPasswords = ReadNonEmptyLinesFromEnv("DB_USER_PASSWORDS")
	return adminPasswords, userPasswords
}

// Reads non-empty pass lines from env variable
func ReadNonEmptyLinesFromEnv(envVar string) []string {
	value := os.Getenv(envVar)
	lines := strings.Split(value, "\n")
	nonEmptyLines := FilterNonEmptyLines(lines)
	return nonEmptyLines
}
