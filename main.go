package main

import (
	"fmt"

	"godice/roller"
)

func main() {
	result := roller.Roll(6)
	fmt.Println("Result:", result)
}
