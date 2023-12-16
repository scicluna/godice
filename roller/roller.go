package roller

import (
	"godice/roller/validators"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

func RollDiceString(rollString string) (*RollResult, error) {
	dicePairs, operands, parseErr := parseRollString(rollString)

	if parseErr != nil {
		return nil, parseErr
	}

	sizes, quantities, specials := parseDicePairs(dicePairs)

	// Create results array and populate it with the results of each roll
	var groupedResults [][][]int
	for i := 0; i < len(sizes); i++ {
		rollResult := rollSet(sizes[i], quantities[i], specials[i])
		groupedResults = append(groupedResults, rollResult)
	}

	// Use the operands to combine the results into a total
	var totals []int
	for _, groupResult := range groupedResults {
		setTotal := 0
		for _, result := range groupResult {
			setTotal += sumRolls(result)
		}
		totals = append(totals, setTotal)
	}

	var grandTotal int
	if len(operands) == 0 {
		// If there are no operands, just sum up all totals
		grandTotal = sumRolls(totals)
	} else {
		grandTotal = calculateWithOrderOfOperations(totals, operands)
	}

	// Append all the results + the total to a string and return
	HTMLProps := buildHTMLProps(groupedResults, totals, sizes, operands, grandTotal)
	return &HTMLProps, nil
}

// given the size of a die, roll a single die and return the result
func singleRoll(size int, special string) []int {
	var result []int
	exploding := false

	if special == "!" {
		exploding = true
	}

	for {
		roll := rand.Intn(size) + 1
		result = append(result, roll)

		if roll == size && exploding {
			continue
		} else {
			break
		}
	}

	return result
}

func rollSet(size int, quantity int, special string) [][]int {
	var rollResults [][]int
	for i := 0; i < quantity; i++ {
		if size != 0 {
			rollResult := singleRoll(size, special)
			rollResults = append(rollResults, rollResult)
		} else {
			// Handle simple modifiers as individual results
			rollResults = append(rollResults, []int{quantity})
			break // No need to loop since it's not a dice roll
		}
	}
	return rollResults
}

// split a roll string into useful parts (array of dice sizes, array of dice quantities, array of modifiers, array of operands)
func parseRollString(rollString string) ([]string, []string, error) {
	if err := validators.RollValidator(rollString); err != nil {
		return nil, nil, err
	}
	reOperands := regexp.MustCompile(`([\+\-\*/])`)
	dicePairs := reOperands.Split(rollString, -1)
	operands := reOperands.FindAllString(rollString, -1)

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
