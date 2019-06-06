package fortuna

import (
	"sync"

	"github.com/jacoelho/fortuna/entropy"
)

// Manager collects entropy from sources and adds it accumulator
type Manager struct {
	Pools  int
	Logger func(v ...interface{})

	id int
	mu sync.Mutex
}

// Accumulate collects entropy and distributes it to accumulator pools
func (m *Manager) Accumulate(dst entropy.Accumulator, src entropy.Source) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := m.id
	m.id = (m.id + 1) % 255

	go func(sourceID, pools int) {
		var poolID int

		for {
			data, err := src.Entropy()
			if err != nil {
				if m.Logger != nil {
					m.Logger("err", err)
				}
				return
			}

			if err := dst.AddRandomEvent(sourceID, poolID, data); err != nil {
				if m.Logger != nil {
					m.Logger("err", err)
				}
			}

			poolID = (poolID + 1) % pools
		}
	}(id, m.Pools)
}
