package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"godice/roller"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Type your dice roll stirng (e.g. 2d6+1):")

	diceString, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	result, err := roller.RollDiceString(strings.TrimSpace(diceString))
	if err != nil {
		fmt.Println("Error rolling dice:", err)
		return
	}

	fmt.Println("Result:", result)
}
