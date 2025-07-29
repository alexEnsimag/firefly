package words

import "unicode"

type Filter func(string) bool

func IsMinimumSize(minSize int) Filter {
	return func(w string) bool {
		return len(w) >= minSize
	}
}

func IsAlpha() Filter {
	return func(w string) bool {
		for _, r := range w {
			if !unicode.IsLetter(r) {
				return false
			}
		}
		return true
	}
}
