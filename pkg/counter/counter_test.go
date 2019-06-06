package counter_test

import (
	"testing"

	"github.com/jacoelho/fortuna/pkg/counter"
)

func TestCounter(t *testing.T) {
	c, err := counter.New(128)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	for i := 0; i < 789; i++ {
		c.Increment()
	}

	if c.Value() != 789 {
		t.Fatalf("unexpected counter value %d", c.Value())
	}
}
