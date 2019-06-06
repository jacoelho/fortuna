package oracle

import (
	"errors"
	"math/rand"
)

type Oracle struct {
	rand *rand.Rand
}

func New(r *rand.Rand) *Oracle {
	return &Oracle{
		rand: r,
	}
}

func (o *Oracle) Numbers(min, max, count int) ([]int, error) {
	if count < 0 {
		return nil, errors.New("invalid count provided")
	}

	result := make([]int, 0, count)
	for i := 0; i < count; i++ {
		result = append(result, min+rand.Intn(max-min))
	}

	return result, nil
}

func (o *Oracle) Sequence(min, max, count int) ([]int, error) {
	if count < 0 || (count > max-min) {
		return nil, errors.New("invalid count provided")
	}

	tmp := make(map[int]struct{}, count)
	for count > 0 {
		n := min + rand.Intn(max-min)

		_, found := tmp[n]
		if found {
			continue
		}

		tmp[n] = struct{}{}
		count--
	}

	result := make([]int, 0, count)
	for v := range tmp {
		result = append(result, v)
	}

	return result, nil
}
