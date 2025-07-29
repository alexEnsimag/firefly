package app

import (
	"context"
	"encoding/json"
	"firefly/alex/internal/words"
	"log/slog"
	"os"
	"testing"
)

func TestEssayWordCountWithEndgUrlsShort(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	wordBank, err := words.NewBank("../../data/words.txt", words.IsMinimumSize(3), words.IsAlpha())
	if err != nil {
		t.Fatalf("Failed to load word bank: %v", err)
	}

	ewc := NewEssayWordCount(logger, "../../data/endg-urls-short", wordBank)
	result, err := ewc.Run(ctx, 10, 10, 1)
	if err != nil {
		t.Fatalf("EssayWordCount.Run failed: %v", err)
	}

	expected := []WordCount{
		{"the", 168},
		{"for", 75},
		{"and", 70},
		{"that", 56},
		{"while", 34},
		{"space", 27},
		{"will", 25},
		{"this", 24},
		{"its", 23},
		{"not", 22},
	}

	if len(result) != len(expected) {
		b, _ := json.MarshalIndent(result, "", "  ")
		t.Fatalf("Expected %d results, got %d. Actual: %s", len(expected), len(result), string(b))
	}

	for i := range expected {
		if result[i] != expected[i] {
			b, _ := json.MarshalIndent(result, "", "  ")
			t.Errorf("Mismatch at index %d: got %+v, want %+v.\nFull result: %s", i, result[i], expected[i], string(b))
		}
	}
}

func TestNormalizeWord(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Hello", "hello"},
		{"Hello!", "hello"},
		{"co-operate", "co-operate"},
		{"don't", "don't"}, // apostrophe preserved
		{"123abc", "123abc"},
		{"foo_bar", "foo_bar"},
		{"Café.", "café"},
		{"end.", "end"},
		{"(parentheses)", "parentheses"},
		{"rock'n'roll'", "rock'n'roll'"}, // apostrophes preserved
		{"'quoted'", "quoted'"},          // leading/trailing apostrophes preserved
	}

	for _, tt := range tests {
		got := normalizeWord(tt.input)
		if got != tt.want {
			t.Errorf("normalizeWord(%q) = %q; want %q", tt.input, got, tt.want)
		}
	}
}
