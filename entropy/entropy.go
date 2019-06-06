package entropy

// Source represents an entropy source
// 9.5.1
type Source interface {
	Entropy() ([]byte, error)
}

// Accumulator collects real random data from various sources.
// 9.5.1
type Accumulator interface {
	AddRandomEvent(id int, poolID int, data []byte) error
}
