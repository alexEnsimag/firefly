package words

import "testing"

func TestIsMinimumSize(t *testing.T) {
	tests := []struct {
		input    string
		minSize  int
		expected bool
	}{
		{"hello", 3, true},
		{"hi", 3, false},
		{"", 0, true},
		{"a", 1, true},
		{"a", 2, false},
	}

	for _, tt := range tests {
		filter := IsMinimumSize(tt.minSize)
		result := filter(tt.input)
		if result != tt.expected {
			t.Errorf("IsMinimumSize(%d)(%q) = %v; want %v", tt.minSize, tt.input, result, tt.expected)
		}
	}
}

func TestIsAlpha(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"hello", true},
		{"HelloWorld", true},
		{"hello123", false},
		{"hello world", false},
		{"", true},
		{"abcABC", true},
		{"abc-ABC", false},
	}

	filter := IsAlpha()
	for _, tt := range tests {
		result := filter(tt.input)
		if result != tt.expected {
			t.Errorf("IsAlpha()(%q) = %v; want %v", tt.input, result, tt.expected)
		}
	}
}
