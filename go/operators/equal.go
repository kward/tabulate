package operators

func EqualSlices(x, y []interface{}) bool {
	// Special cases.
  switch {
	case len(x) != len(y):
		return false
	}
	for i, v := range x {
		if v != y[i] {
			return false
		}
	}
	return true
}

func EqualSlicesOfByte(x, y []byte) bool {
	// Special cases.
	switch {
	case len(x) != len(y):
		return false
	}
	for i, v := range x {
		if v != y[i] {
			return false
		}
	}
	return true
}

func EqualSlicesOf2Byte(x, y [][]byte) bool {
	// Special cases.
	switch {
	case len(x) != len(y):
		return false
	}
	for i, iv := range x {
		jv := y[i]
		switch {
		case !EqualSlicesOfByte(iv, jv):
			return false
		}
	}
	return true
}

func EqualSlicesOfInt(x, y []int) bool {
	// Special cases.
	switch {
	case len(x) != len(y):
		return false
	}
	for i, v := range x {
		if v != y[i] {
			return false
		}
	}
	return true
}

func EqualSlicesOfString(x, y []string) bool {
	// Special cases.
	switch {
	case len(x) != len(y):
		return false
	}
	for i, v := range x {
		if v != y[i] {
			return false
		}
	}
	return true
}
