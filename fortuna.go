package fortuna

import (
	"crypto/aes"
	"crypto/sha256"
	"errors"
	"io"
	"math/rand"
	"sync"
	"time"
)

const (
	// NumPools number of pools 9.5.2
	NumPools = 32
	// MinPoolSize minimum pool size 9.5.5
	MinPoolSize = 64
)

var (
	_ PRNG = (*Fortuna)(nil)
)

// PRNG represents a prng source
// To get random data use rand.Source or io.Reader interface
type PRNG interface {
	rand.Source
	io.Reader
}

// Fortuna is a cryptographically secure pseudorandom number generator
// Devised by Bruce Schneier and Niels Ferguson
type Fortuna struct {
	pools     []*Pool
	generator *Generator

	seededCount    uint64
	lastSeededTime time.Time

	mu sync.Mutex
}

// New creates and initializes an unseed fortuna prng
func New() (*Fortuna, error) {
	p := make([]*Pool, NumPools)

	for idx := range p {
		p[idx] = NewPool()
	}

	g, err := NewGenerator(sha256.New(), aes.NewCipher)
	if err != nil {
		return nil, err
	}

	return &Fortuna{
		pools:          p,
		generator:      g,
		lastSeededTime: time.Now().Add(-time.Minute),
	}, nil
}

// Read reads random data from generator
func (f *Fortuna) Read(p []byte) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.shouldReseed() {
		if err := f.reseed(); err != nil {
			return -1, err
		}
	}

	if f.seededCount == 0 {
		return -1, errors.New("fortuna not initialized")
	}

	return f.generator.Read(p)
}

// Seed will seed following rand.Source interface
func (f *Fortuna) Seed(seed int64) {
	if err := f.AddRandomEvent(0, 0, int64ToBytes(seed)); err != nil {
		panic(err)
	}
}

// SeedBytes will seed using bytes allowing more complex seeding
func (f *Fortuna) SeedBytes(seed []byte) {
	if err := f.AddRandomEvent(0, 0, seed); err != nil {
		panic(err)
	}
}

// Int63 returns a non-negative pseudo-random 63-bit integer as an int64.
func (f *Fortuna) Int63() int64 {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.shouldReseed() {
		if err := f.reseed(); err != nil {
			panic(err)
		}
	}

	if f.seededCount == 0 {
		panic("fortuna not initialized")
	}

	return f.generator.Int63()
}

func (f *Fortuna) shouldReseed() bool {
	return f.pools[0].size > MinPoolSize && time.Now().After(f.lastSeededTime)
}

func (f *Fortuna) reseed() error {
	for i := uint(0); i < NumPools; i++ {
		if f.seededCount%(1<<i) == 0 {
			data, err := f.pools[i].Sum()
			if err != nil {
				return err
			}

			if _, err := f.generator.hash.Write(data); err != nil {
				return err
			}
		}
	}

	if err := f.generator.reseed(); err != nil {
		return err
	}

	f.seededCount++
	f.lastSeededTime = time.Now().Add(100 * time.Millisecond)

	return nil
}

// AddRandomEvent adds an `event` from source `id` to pool `pooID`
// 9.5.6
func (f *Fortuna) AddRandomEvent(id int, poolID int, event []byte) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.pools[poolID].Write([]byte{byte(id), byte(len(event))})
	f.pools[poolID].Write(event)

	if poolID == 0 {
		return f.reseed()
	}

	return nil
}
