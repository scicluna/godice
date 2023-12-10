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
	var results [][]int
	for i := 0; i < len(sizes); i++ {
		rollResult := rollSet(sizes[i], quantities[i], specials[i])
		results = append(results, rollResult)
	}

	//use the operands to combine the results into a total
	var total []int
	for _, result := range results {
		setTotal := 0
		for _, roll := range result {
			setTotal += roll
		}
		total = append(total, setTotal)
	}

	var grandTotal int
	//if there are no operands, just return the result
	grandTotal += total[0]

	//otherwise, combine the results using the operands
	for i := 0; i < len(operands); i++ {
		switch operands[i] {
		case "+":
			grandTotal += total[i+1]
		case "-":
			grandTotal -= total[i+1]
		case "*":
			grandTotal *= total[i+1]
		case "/":
			if total[i+1] != 0 {
				grandTotal /= total[i+1]
			}
		}
	}

	//append all the results + the total to a string and return
	resultString := buildString(total, results, operands, grandTotal)
	return resultString, nil
}

// given the size of a die, roll a single die and return the result
func singleRoll(size int, special string) int {
	var result int
	exploding := false

	if special == "!" {
		exploding = true
	}

	for {
		roll := rand.Intn(size) + 1
		result += roll

		if roll == size && exploding {
			continue
		} else {
			break
		}
	}

	return result
}

func rollSet(size int, quantity int, special string) []int {
	var rollResults []int
	for i := 0; i < quantity; i++ {
		if size != 0 {
			rollResult := singleRoll(size, special)
			rollResults = append(rollResults, rollResult)
		} else {
			// Handle simple modifiers as individual results
			rollResults = append(rollResults, quantity)
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

func specialParse(brokenPair []string) (string, []string) {
	//check for exploding dice (!)
	if strings.Contains(brokenPair[1], "!") && brokenPair[1] != "1" {
		brokenPair[1] = strings.Replace(brokenPair[1], "!", "", -1)
		return "!", brokenPair
	}

	return "", brokenPair
}

func buildString(totals []int, rawTotals [][]int, operands []string, grandTotal int) string {
	var resultString strings.Builder

	for i, rawTotal := range rawTotals {
		// Append the individual roll results for each set
		resultString.WriteString("(")
		for j, roll := range rawTotal {
			resultString.WriteString(strconv.Itoa(roll))
			if j < len(rawTotal)-1 {
				resultString.WriteString(", ")
			}
		}
		resultString.WriteString(")")

		// Append the operand, if there is one
		if i < len(operands) {
			resultString.WriteString(" " + operands[i] + " ")
		}
	}

	// Append the grand total at the enda
	resultString.WriteString(" = " + strconv.Itoa(grandTotal))

	return resultString.String()
}
