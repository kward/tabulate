package strings

import (
	"testing"

	myoperators "github.com/kward/tabulate/go/operators"
)

func TestSplitNMerged(t *testing.T) {
	t.Parallel()

	var got, want []string

	want = []string{"1", "2", "3"}
	got = SplitNMerged("1 2 3", " ", -1)
	if !myoperators.EqualSlicesOfString(want, got) {
		t.Errorf("SplitNMerged(): want %v, got %v", want, got)
	}

	want = []string{"1", "2", "3"}
	got = SplitNMerged("1 2   3", " ", -1)
	if !myoperators.EqualSlicesOfString(want, got) {
		t.Errorf("SplitNMerged(): want %v, got %v", want, got)
	}

	want = []string{}
	got = SplitNMerged("", " ", -1)
	if !myoperators.EqualSlicesOfString(want, got) {
		t.Errorf("SplitNMerged(): want %v, got %v", want, got)
	}
}

func TestStretch(t *testing.T) {
	t.Parallel()

	var got, have, want string

	// Long.
	have = "str1"
	got = Stretch(have, ' ', 6)
	want = "str1  "
	if want != got {
		t.Errorf("Stretch(%q): got %v, want %v", have, want, got)
	}
	// Just right.
	have = "str2"
	got = Stretch(have, ' ', 4)
	want = "str2"
	if want != got {
		t.Errorf("Stretch(%q): got %v, want %v", have, want, got)
	}
	// Short.
	have = "str3"
	got = Stretch(have, ' ', 1)
	want = "str3"
	if want != got {
		t.Errorf("Stretch(%q): got %v, want %v", have, want, got)
	}
}
