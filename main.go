package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func readFile(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", errors.New("пустой путь к файлу")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func splitWords(text string) []string {
	replacer := strings.NewReplacer(
		"\n", " ",
		"\r", " ",
		"\t", " ",
		".", " ",
		",", " ",
		"!", " ",
		"?", " ",
		";", " ",
		":", " ",
	)

	cleanText := replacer.Replace(text)
	return strings.Fields(cleanText)
}

func countTotalWords(text *string) int {
	if text == nil || strings.TrimSpace(*text) == "" {
		return 0
	}

	words := splitWords(*text)
	return len(words)
}

func countOccurrences(text *string, searchWord string) int {
	if text == nil ||
		strings.TrimSpace(*text) == "" ||
		strings.TrimSpace(searchWord) == "" {
		return 0
	}

	words := splitWords(*text)
	count := 0

	for _, word := range words {
		if strings.EqualFold(word, searchWord) {
			count++
		}
	}

	return count
}

func runTest() {
	fmt.Println("\nТестирование")

	testText := "one two one"
	testWord := "one"

	expected := 2
	result := countOccurrences(&testText, testWord)

	fmt.Println("Ожидаемый результат:", expected)
	fmt.Println("Фактический результат:", result)

	if result == expected {
		fmt.Println("Тест пройден успешно")
	} else {
		fmt.Println("Тест не пройден")
	}
}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Использование:")
		fmt.Println("go run main.go <путь_к_файлу> <слово_для_поиска>")
		return
	}

	filePath := os.Args[1]
	searchWord := os.Args[2]

	content, err := readFile(filePath)
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return
	}

	totalWords := countTotalWords(&content)
	occurrences := countOccurrences(&content, searchWord)

	fmt.Println("Общее количество слов в файле:", totalWords)
	fmt.Printf("Количество повторений слова \"%s\": %d\n", searchWord, occurrences)

	runTest()
}
