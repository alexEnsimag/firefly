package words

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Bank struct {
	words map[string]struct{}
}

func NewBank(path string, filters ...Filter) (*Bank, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	allowedWords, err := LoadWords(file, filters...)
	if err != nil {
		return nil, fmt.Errorf("failed to load words: %w", err)
	}

	return &Bank{
		words: allowedWords,
	}, nil
}

func (w *Bank) IsValid(word string) bool {
	_, ok := w.words[word]
	return ok
}

func LoadWords(r io.Reader, filters ...Filter) (map[string]struct{}, error) {
	allowedWords := map[string]struct{}{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		word := scanner.Text()
		if isValidWord(word, filters...) {
			allowedWords[strings.ToLower(word)] = struct{}{}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate words: %w", err)
	}

	return allowedWords, nil
}

func isValidWord(word string, filters ...Filter) bool {
	for _, f := range filters {
		if !f(word) {
			return false
		}
	}
	return true
}
