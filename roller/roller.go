package roller

import (
	"math/rand"
	"regexp"
)

// given the size of a die, roll a single die and return the result
func singleRoll(size int) int {
	return rand.Intn(size) + 1
}

// split a roll string into useful parts (array of dice sizes, array of dice quantities, array of modifiers, array of operands)
func parseRollString(rollString string) ([]int, []int, []string, []string, error) {
	reOperands := regexp.MustCompile(`([+-*/])`)
	dicePairs := reOperands.Split(rollString, -1)
	operands := reOperands.FindAllString(rollString, -1)

}
