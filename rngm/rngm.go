// Package rngm (random manager) provides a single seeded random source
// shared across the entire application, enabling deterministic output.
package rngm

import "math/rand"

var r *rand.Rand

// Init seeds the shared random source. Call once at program startup.
func Init(seed int64) {
	r = rand.New(rand.NewSource(seed))
}

// Intn returns a non-negative random int in [0, n).
func Intn(n int) int {
	return r.Intn(n)
}
