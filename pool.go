package fortuna

import (
	"crypto/sha256"
	"hash"
)

// Pool represents an entropy pool
// 9.5.2
type Pool struct {
	hash hash.Hash
	size int
}

// NewPool create a new pool using sha256 as base hash function
func NewPool() *Pool {
	return &Pool{
		hash: sha256.New(),
	}
}

// Write adds more data to the running hash
func (p *Pool) Write(data []byte) (int, error) {
	p.size += len(data)

	return p.hash.Write(data)
}

// Sum returns sha256-d and empties pool
// 9.5.5
func (p *Pool) Sum() ([]byte, error) {
	round1 := p.hash.Sum(nil)
	p.hash.Reset()

	if _, err := p.hash.Write(round1); err != nil {
		return nil, err
	}

	round2 := p.hash.Sum(nil)
	p.hash.Reset()
	p.size = 0

	return round2, nil
}
