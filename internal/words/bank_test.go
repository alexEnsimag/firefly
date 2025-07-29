package words

import (
	"strings"
	"testing"
)

func TestLoadWords(t *testing.T) {
	data := "apple\nbanana\ncarrot\n"
	reader := strings.NewReader(data)

	wordsMap, err := LoadWords(reader, IsMinimumSize(6))
	if err != nil {
		t.Fatalf("loadWords failed: %v", err)
	}

	tests := []struct {
		word     string
		expected bool
	}{
		{"apple", false},
		{"banana", true},
		{"carrot", true},
		{"durian", false},
		{"", false},
	}

	for _, tt := range tests {
		_, ok := wordsMap[tt.word]
		if ok != tt.expected {
			t.Errorf("loadWords: word %q present = %v; want %v", tt.word, ok, tt.expected)
		}
	}
}

func TestBankIsValid(t *testing.T) {
	bank := &Bank{words: map[string]struct{}{
		"apple":  {},
		"banana": {},
	}}

	tests := []struct {
		word     string
		expected bool
	}{
		{"apple", true},
		{"banana", true},
		{"carrot", false},
		{"", false},
	}

	for _, tt := range tests {
		got := bank.IsValid(tt.word)
		if got != tt.expected {
			t.Errorf("IsValid(%q) = %v; want %v", tt.word, got, tt.expected)
		}
	}
}
