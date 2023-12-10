package roller

import (
	"godice/roller/validators"
	"math/rand"
	"regexp"
)

func RollDiceString(rollString string) (string, error) {
	dicePairs, rollString, parseErr := parseRollString(rollString)

	if parseErr != nil {
		return "", parseErr
	}

}

// given the size of a die, roll a single die and return the result
func singleRoll(size int) int {
	return rand.Intn(size) + 1
}

// split a roll string into useful parts (array of dice sizes, array of dice quantities, array of modifiers, array of operands)
func parseRollString(rollString string) ([]string, []string, error) {
	if err := validators.RollValidator(rollString); err != nil {
		return nil, nil, err
	}
	reOperands := regexp.MustCompile(`([+-*/])`)
	dicePairs := reOperands.Split(rollString, -1)
	operands := reOperands.FindAllString(rollString, -1)

	return dicePairs, operands, nil
}
