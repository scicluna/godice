package roller

type DiceRoll struct {
	Value    int
	RollType string // "min", "med", or "max"
}

type DiceSet struct {
	Rolls   []DiceRoll
	Total   int
	Special string // "!", or ""
}

type RollResult struct {
	Operands   []string
	Sets       []DiceSet
	GrandTotal int
}

func buildHTMLProps(groupedResults [][][]int, totals []int, sizes []int, operands []string, specials []string, grandTotal int) RollResult {
	/*
		If grouped Results are [ [ [1,2,3] ], [[1,2,3],[5,6]]]
		Then we need to have a struct array that shows [{total: 5, dicesets: [[{val:1, min}, {val:2, med}, {val: 3, med}]]}, {total: 16, dicesets:[[{val:1, min}, {val:2, med}, {val: 3, med}], [{val:5, med}, {val:6, max}]] or something like that. what
	*/
	var diceSets []DiceSet
	var operandSets []string
	for i, outerGroup := range groupedResults {
		var diceSet DiceSet
		var minRoll = 1
		var maxRoll = sizes[i]
		diceSet.Total = totals[i]
		diceSet.Special = specials[i]

		if i == len(operands) {
			operandSets = append(operandSets, "")
		} else {
			operandSets = append(operandSets, operands[i])
		}

		for _, innerGroup := range outerGroup {
			var diceRolls []DiceRoll
			for _, roll := range innerGroup {
				var rollType string
				if roll == minRoll {
					rollType = "min"
				} else if roll == maxRoll {
					rollType = "max"
				} else {
					rollType = "med"
				}
				diceRoll := DiceRoll{
					Value:    roll,
					RollType: rollType,
				}
				diceRolls = append(diceRolls, diceRoll)
			}
			diceSet.Rolls = append(diceSet.Rolls, diceRolls...)
		}
		diceSets = append(diceSets, diceSet)
	}
	return RollResult{
		Sets:       diceSets,
		Operands:   operandSets,
		GrandTotal: grandTotal,
	}
}
