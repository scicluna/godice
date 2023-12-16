package roller

func calculateWithOrderOfOperations(totals []int, operands []string) int {
	// Step 0: Clone totals so we dont modify the originals
	tempTotals := make([]int, len(totals))
	copy(tempTotals, totals)

	// Step 1: Handle multiplication and division first
	for i := 0; i < len(operands); i++ {
		switch operands[i] {
		case "*":
			tempTotals[i] *= tempTotals[i+1]
			tempTotals = append(tempTotals[:i+1], tempTotals[i+2:]...)
			operands = append(operands[:i], operands[i+1:]...)
			i-- // Adjust index since we've removed an element
		case "/":
			if tempTotals[i+1] != 0 {
				tempTotals[i] /= tempTotals[i+1]
			}
			tempTotals = append(tempTotals[:i+1], tempTotals[i+2:]...)
			operands = append(operands[:i], operands[i+1:]...)
			i-- // Adjust index since we've removed an element
		}
	}

	// Step 2: Handle addition and subtraction
	grandTotal := tempTotals[0]
	for i := 0; i < len(operands); i++ {
		switch operands[i] {
		case "+":
			grandTotal += tempTotals[i+1]
		case "-":
			grandTotal -= tempTotals[i+1]
		}
	}

	return grandTotal
}

// Helper function to sum up a slice of integers
func sumRolls(rolls []int) int {
	total := 0
	for _, roll := range rolls {
		total += roll
	}
	return total
}
