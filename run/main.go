package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	rootDir := "../"
	f := func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error accessing path:", err)
			return err
		}

		if strings.Contains(path, ".git") {
			return nil
		}

		if !info.IsDir() && strings.Contains(path, ".go"){
			charCount, lineCount, err := countFileStats(path)
			if err != nil {
				fmt.Println("Error processing file:", err)
				return err
			}
			fmt.Printf("| %-40s | %10d | %10d |\n", path, charCount, lineCount)
		}
		return nil
	}
	fmt.Printf("|  | %s | %s |\n", "char", "line")
	fmt.Println("|-|-|")
	err := filepath.Walk(rootDir, f)

	if err != nil {
		fmt.Println("Error walking the path:", err)
	}
}

// countFileStats returns the number of characters and lines in a file
func countFileStats(filename string) (int, int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	var charCount, lineCount int
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++
		charCount += len(line)
	}

	if err := scanner.Err(); err != nil {
		return 0, 0, err
	}

	return charCount, lineCount, nil
}
