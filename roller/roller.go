package roller

import (
	"math/rand"
)

func Roll(sides int) int {
	return rand.Intn(sides) + 1
}
