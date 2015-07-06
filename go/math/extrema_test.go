// Basic max(x, y) routines that aren't provided by default.

package math

import "testing"

func TestMaxInt(t *testing.T) {
	var got, want int

	want = 2

	got = Max(1, 2)
	if got != want {
		t.Errorf("Max(int): got %v, want %v", got, want)
	}

	got = Max(2, 2)
	if got != want {
		t.Errorf("Max(int): got %v, want %v", got, want)
	}

	got = Max(2, 1)
	if got != want {
		t.Errorf("Max(int): got %v, want %v", got, want)
	}
}

func TestMinInt(t *testing.T) {
	var got, want int

	want = 1

	got = Min(1, 2)
	if got != want {
		t.Errorf("Max(int): got %v, want %v", got, want)
	}

	got = Min(1, 1)
	if got != want {
		t.Errorf("Max(int): got %v, want %v", got, want)
	}

	got = Min(2, 1)
	if got != want {
		t.Errorf("Max(int): got %v, want %v", got, want)
	}
}
