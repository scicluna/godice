package roller

import (
	"math/rand"
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
	HTMLProps := buildHTMLProps(groupedResults, totals, sizes, operands, specials, grandTotal)
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
