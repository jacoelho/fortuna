package oracle_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/jacoelho/fortuna/internal/oracle"
)

func TestOracleIntegersWithoutRepetition(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UTC().Unix()))

	o := oracle.New(r)

	t.Fatal(o.Sequence(1, 3, 2))
}

func TestOracleIntegersInterval(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UTC().Unix()))

	o := oracle.New(r)

	for i := 0; i < 1000; i++ {
		if n, _ := o.Numbers(1, 3, 1); n[0] == 1 {
			t.Fatal(n)
		}
	}
}
