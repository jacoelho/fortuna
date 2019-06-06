package fortuna

import (
	"crypto/cipher"
	"errors"
	"hash"

	"github.com/jacoelho/fortuna/pkg/counter"
)

const (
	byteSize           = 8
	keySize            = 32
	maxBytesPerRequest = 1 << 20
)

var (
	_ PRNG = (*Generator)(nil)
)

// NewCipher will return a new cipher with key
type NewCipher func(key []byte) (cipher.Block, error)

// Generator represents a random data generator
// 9.4
type Generator struct {
	key     []byte
	hash    hash.Hash
	cipher  cipher.Block
	counter counter.Counter

	newCipher NewCipher

	buffer []byte
}

// NewGenerator creates a new generator using an hash and a block cipher function
func NewGenerator(hash hash.Hash, f NewCipher) (*Generator, error) {
	g := &Generator{
		key:       make([]byte, keySize, keySize),
		hash:      hash,
		newCipher: f,
	}

	if err := g.updateCipher(); err != nil {
		return nil, err
	}

	ctr, err := counter.New(uint(g.cipher.BlockSize()) * byteSize)
	if err != nil {
		return nil, err
	}

	g.counter = ctr

	g.buffer = make([]byte, len(g.counter), len(g.counter))

	return g, nil
}

// Seed will seed generator following rand.Source interface
func (g *Generator) Seed(seed int64) {
	if _, err := g.hash.Write(int64ToBytes(seed)); err != nil {
		panic(err)
	}

	if err := g.reseed(); err != nil {
		panic(err)
	}
}

// SeedBytes will seed generator using bytes allowing more complex seeding
func (g *Generator) SeedBytes(seed []byte) {
	if _, err := g.hash.Write(seed); err != nil {
		panic(err)
	}

	if err := g.reseed(); err != nil {
		panic(err)
	}
}

// reseed creates a new key using hash function
// 9.4.2
func (g *Generator) reseed() error {
	g.key = g.hash.Sum(nil)

	g.counter.Increment()

	return g.updateCipher()
}

// updateCipher creates a cipher block using a new key
func (g *Generator) updateCipher() error {
	var err error

	g.cipher, err = g.newCipher(g.key)
	if err != nil {
		return err
	}

	g.hash.Reset()

	if _, err = g.hash.Write(g.key); err != nil {
		return err
	}

	return nil
}

// rekey reads enough random data to create a new key
// 9.4.4 "switch to a new key"
func (g *Generator) rekey() error {
	for i := keySize / g.cipher.BlockSize(); i > 0; i-- {
		g.readBlock(g.key[g.cipher.BlockSize()*i:])
	}

	return g.updateCipher()
}

// readBlock reads at most cipher block size random data
func (g *Generator) readBlock(out []byte) int {
	var n int
	if len(out) > g.cipher.BlockSize() {
		n = g.cipher.BlockSize()
	} else {
		n = len(out)
	}

	g.cipher.Encrypt(g.buffer, g.counter)
	g.counter.Increment()

	copy(out, g.buffer[0:n])
	return n
}

// Read reads random data from generator
func (g *Generator) Read(p []byte) (int, error) {
	if g.counter.Value() == 0 {
		return 0, errors.New("attempted to read an initialized generator")
	}

	want := len(p)
	read := 0
	rekeyBytesLeft := maxBytesPerRequest

	for read < want {
		n := g.readBlock(p[read:])
		read += n
		rekeyBytesLeft -= n

		if rekeyBytesLeft < 0 {
			if err := g.rekey(); err != nil {
				return read, err
			}

			rekeyBytesLeft = maxBytesPerRequest
		}
	}

	if err := g.rekey(); err != nil {
		return read, err
	}

	return read, nil
}

// Int63 returns a non-negative pseudo-random 63-bit integer as an int64.
func (g *Generator) Int63() int64 {
	buffer := make([]byte, 8)
	if _, err := g.Read(buffer); err != nil {
		panic(err)
	}

	return parseInt64(buffer)
}

func parseInt64(bytes []byte) (n int64) {
	shift := uint(0)
	for _, b := range bytes {
		n |= int64(b&0x7F) << shift
		shift += 7
	}
	return
}

func int64ToBytes(v int64) []byte {
	b := make([]byte, 8)

	b[0] = byte(v >> 56)
	b[1] = byte(v >> 48)
	b[2] = byte(v >> 40)
	b[3] = byte(v >> 32)
	b[4] = byte(v >> 24)
	b[5] = byte(v >> 16)
	b[6] = byte(v >> 8)
	b[7] = byte(v)

	return b
}
