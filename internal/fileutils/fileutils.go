package fileutils

import (
	"bufio"
	"fmt"
	"os"
)

const MaxParagraphSize = 1_000

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func splitString(input string) []string {
	maxLength := MaxParagraphSize
	inputLength := len(input)
	if inputLength <= maxLength {
		return []string{input}
	}
	segmentCount := inputLength / maxLength
	if inputLength%maxLength != 0 {
		segmentCount++
	}
	var segments []string
	for i := 0; i < segmentCount; i++ {
		start := i * maxLength
		end := start + maxLength

		if end > inputLength {
			end = inputLength
		}
		if end-start < 40 {
			// probably nonsense or unusably small text
			continue
		}

		segments = append(segments, input[start:end])
	}
	return segments
}

func AppendNewParagraph(filename, text string) error {
	// Open file in append mode, or create it if it doesn't exist
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()
	// Append new paragraph with four blank lines and a separator
	if _, err := file.WriteString("\n\n----------\n\n" + text); err != nil {
		return err
	}
	return nil
}

func GatherParagraphs(filePath string) ([]string, error) {
	paragraphs := []string{}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed opening %s : %w", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	nextParagraph := ""
	for scanner.Scan() {
		line := scanner.Text()
		nextParagraph += line
		if len(nextParagraph) > MaxParagraphSize {
			paragraphs = append(paragraphs, splitString(nextParagraph)...)
			nextParagraph = ""
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error during file scan: %s %v", filePath, err)
	}
	return paragraphs, nil
}
