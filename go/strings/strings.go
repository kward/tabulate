package strings

import "strings"

// SplitNMerged slices s into substrings separated by sep and returns a slice
// of the substrings between those separators. If sep is empty, SplitN splits
// after each UTF-8 sequence. Repeated sep will be merged. The count determines
// the number of substrings to return:
//
//  n > 0: at most n substrings; the last substring will be the unsplit remainder.
//  n == 0: the result is nil (zero substrings)
//  n < 0: all substrings
func SplitNMerged(s string, sep string, n int) []string {
	split := strings.SplitN(s, sep, n)
	merged := make([]string, 0, len(split))
	for _, v := range split {
		if v != "" {
			merged = append(merged, v)
		}
	}
	return merged
}

// Stretches a string to a given length by appending a character to the end.
func Stretch(s string, r rune, l int) string {
	// Special cases.
	if len(s) >= l {
		return s
	}
	return s + strings.Repeat(string(r), l-len(s))
}
