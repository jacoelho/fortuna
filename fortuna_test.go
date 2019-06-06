package fortuna_test

import (
	"io"
	"io/ioutil"
	"testing"
	"time"

	"github.com/jacoelho/fortuna"
)

func TestFortunaInt63(t *testing.T) {
	f, err := fortuna.New()
	if err != nil {
		t.Fatal(err)
	}

	f.Seed(time.Now().UTC().UnixNano())

	for i := 0; i < 100000; i++ {
		if n := f.Int63(); n < 0 {
			t.Fatal(n)
		}
	}
}

func TestFortunaInt63AfterReseed(t *testing.T) {
	f, err := fortuna.New()
	if err != nil {
		t.Fatal(err)
	}

	f.Seed(time.Now().UTC().UnixNano())

	// enough to trigger reseed
	time.Sleep(200 * time.Millisecond)

	for i := 0; i < 100000; i++ {
		if n := f.Int63(); n < 0 {
			t.Fatal(n)
		}
	}
}

func TestFortunaRead(t *testing.T) {
	f, err := fortuna.New()
	if err != nil {
		t.Fatal(err)
	}

	f.Seed(time.Now().UTC().UnixNano())

	if _, err := io.CopyN(ioutil.Discard, f, 1<<25); err != nil {
		t.Fatalf("unexpected error %v", err)
	}
}
