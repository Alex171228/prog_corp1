package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// Domain layer

type WordStats struct {
	TotalWords  int
	Occurrences int
}

type FileContent string

// Interfaces

type WordAnalyzer interface {
	Analyze(text FileContent, searchWord string) WordStats
}

type FileReader interface {
	Read(path string) (FileContent, error)
}

// Implementations

type SimpleWordAnalyzer struct{}

func (a SimpleWordAnalyzer) Analyze(text FileContent, searchWord string) WordStats {
	if isEmpty(text) || isEmptyString(searchWord) {
		return WordStats{}
	}

	words := splitWords(string(text))

	return WordStats{
		TotalWords:  len(words),
		Occurrences: countMatches(words, searchWord),
	}
}

// Helpers

func isEmpty(text FileContent) bool {
	return strings.TrimSpace(string(text)) == ""
}

func isEmptyString(s string) bool {
	return strings.TrimSpace(s) == ""
}

func splitWords(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

func countMatches(words []string, searchWord string) int {
	count := 0
	for _, word := range words {
		if strings.EqualFold(word, searchWord) {
			count++
		}
	}
	return count
}

// Infrastructure

type OSFileReader struct{}

func (r OSFileReader) Read(path string) (FileContent, error) {
	if strings.TrimSpace(path) == "" {
		return "", errors.New("file path cannot be empty")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read file: %w", err)
	}

	return FileContent(data), nil
}

// Application layer

type WordAnalysisService struct {
	reader   FileReader
	analyzer WordAnalyzer
}

func NewWordAnalysisService(r FileReader, a WordAnalyzer) *WordAnalysisService {
	return &WordAnalysisService{reader: r, analyzer: a}
}

func (s *WordAnalysisService) AnalyzeFile(path, word string) (WordStats, error) {
	content, err := s.reader.Read(path)
	if err != nil {
		return WordStats{}, err
	}

	return s.analyzer.Analyze(content, word), nil
}

// Presentation

type Formatter interface {
	Format(stats WordStats, filePath, word string) string
	FormatError(err error) string
}

type ConsoleFormatter struct{}

func (f ConsoleFormatter) Format(stats WordStats, filePath, word string) string {
	return fmt.Sprintf(`=== Results ===
File: %s
Word: "%s"
Total words: %d
Occurrences: %d
`, filePath, word, stats.TotalWords, stats.Occurrences)
}

func (f ConsoleFormatter) FormatError(err error) string {
	return fmt.Sprintf("Error: %v", err)
}

// CLI

type CLI struct {
	service   *WordAnalysisService
	formatter Formatter
}

func NewCLI(s *WordAnalysisService, f Formatter) *CLI {
	return &CLI{service: s, formatter: f}
}

func (cli *CLI) Run(args []string) error {
	input, err := parseArgs(args)
	if err != nil {
		return err
	}

	stats, err := cli.service.AnalyzeFile(input.filePath, input.searchWord)
	if err != nil {
		return err
	}

	fmt.Print(cli.formatter.Format(stats, input.filePath, input.searchWord))
	return nil
}

// Input

type Input struct {
	filePath   string
	searchWord string
}

func parseArgs(args []string) (*Input, error) {
	if len(args) < 3 {
		return nil, errors.New("usage: <program> <file_path> <search_word>")
	}

	return &Input{
		filePath:   args[1],
		searchWord: args[2],
	}, nil
}

// Main

func main() {
	reader := OSFileReader{}
	analyzer := SimpleWordAnalyzer{}

	service := NewWordAnalysisService(reader, analyzer)
	formatter := ConsoleFormatter{}
	cli := NewCLI(service, formatter)

	if err := cli.Run(os.Args); err != nil {
		fmt.Println(formatter.FormatError(err))
		os.Exit(1)
	}
}
