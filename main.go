package main

import (
	"fmt"

	"godice/roller"
)

func main() {
	result, err := roller.RollDiceString("2d6")
	if err == nil {
		fmt.Println(result)
	} else {
		fmt.Println(err)
	}
}
