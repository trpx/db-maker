package utils

import (
	"bufio"
	"os"
)

// Reads non-empty pass file lines and returns them as a []string
func readPassFile(file *string) []string {
	fInfo, err := os.Stat(*file)
	if err != nil {
		if os.IsNotExist(err) {
			Panicf("file '%s' does not exist", *file)
		} else {
			Panicf("couldn't stat file '%s'", *file)
		}
	}
	if fInfo.IsDir() {
		Panicf("path '%s' is a directory", *file)
	}
	lines, err := readLines(*file)
	if err != nil {
		Panicf("couldn't read file '%s' lines: %#v", *file, err)
	}
	return lines
}


// https://stackoverflow.com/questions/5884154/read-text-file-into-string-array-and-write
// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
