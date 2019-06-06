package fortuna_test

import (
	"crypto/aes"
	"crypto/sha256"
	"io"
	"io/ioutil"
	"testing"
	"time"

	"github.com/jacoelho/fortuna"
)

func TestGeneratorInt(t *testing.T) {
	g, err := fortuna.NewGenerator(sha256.New(), aes.NewCipher)
	if err != nil {
		t.Fatal(err)
	}

	g.Seed(time.Now().UTC().UnixNano())

	for i := 0; i < 100000; i++ {
		if n := g.Int63(); n < 0 {
			t.Fatal(n)
		}
	}
}

func TestGeneratorRead(t *testing.T) {
	g, err := fortuna.NewGenerator(sha256.New(), aes.NewCipher)
	if err != nil {
		t.Fatal(err)
	}

	g.Seed(time.Now().UTC().UnixNano())

	if _, err := io.CopyN(ioutil.Discard, g, 1<<24); err != nil {
		t.Fatalf("unexpected error %v", err)
	}
}

func BenchmarkGenerate(b *testing.B) {
	g, err := fortuna.NewGenerator(sha256.New(), aes.NewCipher)
	if err != nil {
		b.Fatal(err)
	}

	g.Seed(time.Now().UTC().UnixNano())

	buf := make([]byte, 1000)
	for i := 0; i < b.N; i++ {
		n, err := g.Read(buf)
		if err != nil {
			b.Fatal(n, err)
		}
	}
}
