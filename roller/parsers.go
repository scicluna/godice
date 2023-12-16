package roller

import (
	"godice/roller/validators"
	"regexp"
	"strconv"
	"strings"
)

// split a roll string into useful parts (array of dice sizes, array of dice quantities, array of modifiers, array of operands)
func parseRollString(rollString string) ([]string, []string, error) {
	// Remove all whitespace from the input string
	processedRollString := strings.Join(strings.Fields(rollString), "")

	if err := validators.RollValidator(processedRollString); err != nil {
		return nil, nil, err
	}

	reOperands := regexp.MustCompile(`([\+\-\*/])`)
	dicePairs := reOperands.Split(processedRollString, -1)
	operands := reOperands.FindAllString(processedRollString, -1)

	return dicePairs, operands, nil
}

// handles the parsing of dice pairs into sizes, quantities, and specials
func parseDicePairs(dicePairs []string) ([]int, []int, []string) {
	var sizes []int
	var quantities []int
	var specials []string

	for _, pair := range dicePairs {
		brokenPair := strings.Split(pair, "d")

		if len(brokenPair) == 1 {
			dieSize, _ := strconv.Atoi(brokenPair[0])

			sizes = append(sizes, 0)
			quantities = append(quantities, dieSize)
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

// looks for special characters to populate the specials array
func specialParse(brokenPair []string) (string, []string) {
	//check for exploding dice (!)
	if strings.Contains(brokenPair[1], "!") && brokenPair[1] != "1" {
		brokenPair[1] = strings.Replace(brokenPair[1], "!", "", -1)
		return "!", brokenPair
	}

	return "", brokenPair
}
