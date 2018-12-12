package client

import (
	"fmt"

	"github.com/avalchev94/connect4"
)

func render(field connect4.Field) {
	fmt.Println("----------------------")
	for i := len(field[0]) - 1; i >= 0; i-- {
		for j := 0; j < len(field); j++ {
			fmt.Printf(" %s ", field[j][i])
		}
		fmt.Println()
	}
	fmt.Println("----------------------")

	for i := range field {
		fmt.Printf(" %d ", i+1)
	}
	fmt.Println("\n----------------------")
}
