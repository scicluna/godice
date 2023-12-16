package roller

type DiceRoll struct {
	value    int
	rollType string // "min", "med", or "max"
}

type DiceSet struct {
	rolls []DiceRoll
	total int
}

type RollResult struct {
	sets       []DiceSet
	grandtotal int
}

func buildHTMLProps(groupedResults [][][]int, totals []int, sizes []int, grandTotal int) RollResult {
	/*
		If grouped Results are [ [ [1,2,3] ], [[1,2,3],[5,6]]]
		Then we need to have a struct array that shows [{total: 5, dicesets: [[{val:1, min}, {val:2, med}, {val: 3, med}]]}, {total: 16, dicesets:[[{val:1, min}, {val:2, med}, {val: 3, med}], [{val:5, med}, {val:6, max}]] or something like that. what
	*/
	var diceSets []DiceSet
	for i, outerGroup := range groupedResults {
		var diceSet DiceSet
		var minRoll = 1
		var maxRoll = sizes[i]
		diceSet.total = totals[i]

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
					value:    roll,
					rollType: rollType,
				}
				diceRolls = append(diceRolls, diceRoll)
			}
			diceSet.rolls = append(diceSet.rolls, diceRolls...)
		}
		diceSets = append(diceSets, diceSet)
	}
	return RollResult{
		sets:       diceSets,
		grandtotal: grandTotal,
	}
}