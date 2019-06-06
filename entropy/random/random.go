package random

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"sync"
	"time"
)

const urandomDevice = "/dev/urandom"

// Random provides entropy by reading /dev/urandom device
type Random struct {
	Interval time.Duration
	f        io.Reader
	tmp      bytes.Buffer
	mu       sync.Mutex
}

// Entropy reads random data from /dev/urandom
func (r *Random) Entropy() ([]byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.Interval == 0 {
		r.Interval = time.Second
	}

	time.Sleep(r.Interval)

	if r.f == nil {
		f, err := os.Open(urandomDevice)
		if err != nil {
			return nil, err
		}

		r.f = bufio.NewReader(f)
	}

	r.tmp.Reset()
	_, err := io.CopyN(&r.tmp, r.f, 32)
	if err != nil {
		return nil, err
	}

	return r.tmp.Bytes(), nil
}
