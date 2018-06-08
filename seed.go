package lostinspace

import "math/rand"

type Seed struct {
	Number int64

	perm [256]int
}

func NewSeed(number int64) *Seed {
	seed := new(Seed)

	seed.Number = number
	seed.init()

	return seed
}

func (seed *Seed) init() {
	random := rand.New(rand.NewSource(seed.Number))
	copy(seed.perm[:], random.Perm(256))
}
