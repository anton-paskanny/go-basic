package main

import "fmt"

const (
	USDToEUR = 0.92
	USDToRUB = 94.50
)

// readInput reads a float number from user input
func readInput(prompt string) float64 {
	var value float64
	fmt.Print(prompt)
	fmt.Scanln(&value)
	return value
}

// convertCurrency is a stub function that will convert from one currency to another
func convertCurrency(amount float64, from string, to string) float64 {
	// TODO: implement real conversion
	fmt.Printf("Converting %.2f from %s to %s...\n", amount, from, to)
	return 0
}

func main() {
	amount := readInput("Enter amount: ")

	result := convertCurrency(amount, "USD", "EUR")

	fmt.Printf("Result: %.2f\n", result)
}
