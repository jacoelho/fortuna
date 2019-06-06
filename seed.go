package fortuna

import (
	"errors"
	"io"
	"io/ioutil"

	"github.com/google/renameio"
)

// SeedFile represents a seed file
// 9.6
type SeedFile struct {
	src  PRNG
	dest string
}

// NewSeedFile creates a new seed file represention
func NewSeedFile(src PRNG, filename string) *SeedFile {
	return &SeedFile{
		src:  src,
		dest: filename,
	}
}

// Write seed to a file
// 9.6.1
func (s *SeedFile) Write() error {
	t, err := renameio.TempFile("", s.dest)
	if err != nil {
		return err
	}

	defer t.Cleanup()

	if _, err := io.CopyN(t, s.src, 64); err != nil {
		return err
	}

	return t.CloseAtomicallyReplace()
}

// Read seed file and truncate it to avoid reusing
// 9.6.2
func (s *SeedFile) Read() (b []byte, err error) {
	b, err = ioutil.ReadFile(s.dest)
	if err != nil {
		return
	}

	if len(b) < 64 {
		err = errors.New("failed to read 64 bytes")
	}

	err = ioutil.WriteFile(s.dest, []byte(""), 0600)

	return
}
