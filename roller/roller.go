package roller

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
)

//parse the roll string (2d6) and then process the roll

func RollDice(roll string) (int, error) {
	var result int
	parsedRoll := strings.Split(roll, "d")
	quantity, errQuantity := strconv.Atoi(parsedRoll[0])
	sides, errSides := strconv.Atoi(parsedRoll[1])

	if errQuantity != nil || errSides != nil {
		return 0, errors.New("Invalid roll")
	}

	for i := 0; i < quantity; i++ {
		result += Roll(sides)
	}
	return result, nil
}

func Roll(sides int) int {
	return rand.Intn(sides) + 1
}
