package lostinspace

import "math/rand"

type Seed struct {
	number int64
	perm   [256]int
}

func NewSeed(number int64) *Seed {
	seed := new(Seed)

	seed.number = number
	random := rand.New(rand.NewSource(number))
	copy(seed.perm[:], random.Perm(256))

	return seed
}
