package operators

import "testing"

func TestEqualSlicesOfByte(t *testing.T) {
	// Equal.
	if !EqualSlicesOfByte([]byte{}, []byte{}) {
		t.Errorf("EqualSlicesOfByte([]byte): [] == []")
	}
	if !EqualSlicesOfByte([]byte{1, 2, 3}, []byte{1, 2, 3}) {
		t.Errorf("EqualSlicesOfByte([]byte): [1, 2, 3] == [1, 2, 3]")
	}
	// Not equal.
	if EqualSlicesOfByte([]byte{1, 2, 3}, []byte{}) {
		t.Errorf("EqualSlicesOfByte([]byte): [1, 2, 3] != []")
	}
	if EqualSlicesOfByte([]byte{1, 2, 3}, []byte{4, 5, 6}) {
		t.Errorf("EqualSlicesOfByte([]byte): [1, 2, 3] != [4, 5, 6]")
	}
}

func TestEqualSlicesOf2Byte(t *testing.T) {
	// Equal.
	if !EqualSlicesOf2Byte([][]byte{}, [][]byte{}) {
		t.Errorf("EqualSlicesOf2Byte([][]byte): [][] == [][]")
	}
	if !EqualSlicesOf2Byte([][]byte{{1, 2}, {3, 4}}, [][]byte{{1, 2}, {3, 4}}) {
		t.Errorf("EqualSlicesOf2Byte([][]byte): [1, 2][3, 4] == [1, 2][3, 4]")
	}
	// Not equal.
	if EqualSlicesOf2Byte([][]byte{{1, 2}, {3, 4}}, [][]byte{}) {
		t.Errorf("EqualSlicesOf2Byte([][]byte): [1, 2][3, 4] == []")
	}
	if EqualSlicesOf2Byte([][]byte{{1, 2}}, [][]byte{{1, 2}, {3, 4}}) {
		t.Errorf("EqualSlicesOf2Byte([][]byte): [1, 2][3, 4] == [1, 2][3, 4]")
	}
}

func TestEqualSlicesOfInt(t *testing.T) {
	// Equal.
	if !EqualSlicesOfInt([]int{}, []int{}) {
		t.Errorf("EqualSlicesOfInt([]int): [] == []")
	}
	if !EqualSlicesOfInt([]int{1, 2, 3}, []int{1, 2, 3}) {
		t.Errorf("EqualSlicesOfInt([]int): [1, 2, 3] == [1, 2, 3]")
	}
	// Not equal.
	if EqualSlicesOfInt([]int{1, 2, 3}, []int{}) {
		t.Errorf("EqualSlicesOfInt([]int): [1, 2, 3] != []")
	}
	if EqualSlicesOfInt([]int{1, 2, 3}, []int{4, 5, 6}) {
		t.Errorf("EqualSlicesOfInt([]int): [1, 2, 3] != [4, 5, 6]")
	}
}

func TestEqualSliceOfString(t *testing.T) {
	// Equal.
	if !EqualSlicesOfString([]string{}, []string{}) {
		t.Errorf("EqualSlicesOfString([]string): [] == []")
	}
	if !EqualSlicesOfString([]string{"1", "2", "3"}, []string{"1", "2", "3"}) {
		t.Errorf("EqualSlicesOfString([]string): [1, 2, 3] == [1, 2, 3]")
	}
	// Not equal.
	if EqualSlicesOfString([]string{"1", "2", "3"}, []string{}) {
		t.Errorf("EqualSlicesOfString([]string): [1, 2, 3] != []")
	}
	if EqualSlicesOfString([]string{"1", "2", "3"}, []string{"4", "5", "6"}) {
		t.Errorf("EqualSlicesOfString([]string): [1, 2, 3] != [4, 5, 6]")
	}
}
