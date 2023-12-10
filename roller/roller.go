package roller

import (
	"godice/roller/validators"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

func RollDiceString(rollString string) (string, error) {
	dicePairs, operands, parseErr := parseRollString(rollString)

	if parseErr != nil {
		return "", parseErr
	}

	sizes, quantities, specials := parseDicePairs(dicePairs)

	//create results array and populate it with the results of each roll
	var total int
	var results []int

	//append all the results + the total to a string and return
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

func parseDicePairs(dicePairs []string) ([]int, []int, []string) {
	var sizes []int
	var quantities []int
	var specials []string

	for _, pair := range dicePairs {
		brokenPair := strings.Split(pair, "d")

		if len(brokenPair) == 1 {
			dieSize, _ := strconv.Atoi(brokenPair[0])

			sizes = append(sizes, dieSize)
			quantities = append(quantities, 0)
			specials = append(specials, "")
		} else {
			detectedSpecial, despecializedPair := specialParse(brokenPair)
			dieSize, _ := strconv.Atoi(despecializedPair[1])
			dieQuantity, _ := strconv.Atoi(despecializedPair[0])

			sizes = append(sizes, dieSize)
			quantities = append(quantities, dieQuantity)
			specials = append(specials, detectedSpecial)
		}
	}
	return sizes, quantities, specials
}

func specialParse(brokenPair []string) (string, []string) {
	//check for exploding dice (!)
	if strings.Contains(brokenPair[1], "!") {
		brokenPair[1] = strings.Replace(brokenPair[1], "!", "", -1)
		return "!", brokenPair
	}

	return "", brokenPair
}
