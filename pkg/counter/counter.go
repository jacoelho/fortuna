package counter

import "errors"

const byteSize = 8

// Counter represents a byte counter
type Counter []byte

var (
	// ErrCounterInvalidSize represents an invalid counter size
	ErrCounterInvalidSize = errors.New("invalid size")
)

// New creates a new counter with size bits
func New(size uint) (Counter, error) {
	if size%byteSize > 0 {
		return Counter(nil), ErrCounterInvalidSize
	}

	return Counter(make([]byte, size/byteSize, size/byteSize)), nil
}

// Increment increments counter by one
func (c Counter) Increment() {
	for i := 0; i < len(c); i++ {
		c[i]++
		if c[i] != 0 {
			break
		}
	}
}

// Value returns counter current value
func (c Counter) Value() uint {
	if len(c) == 0 {
		return 0
	}

	acc := uint(c[0])
	for i := uint(1); i < uint(len(c)); i++ {
		acc += uint(c[i]) << (8 * i)
	}

	return acc
}
