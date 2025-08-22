package main

import "fmt"

const (
	USDToEUR = 0.92
	USDToRUB = 94.50
)

func main() {
	eurToRub := USDToRUB / USDToEUR

	fmt.Println("Base conversion rates:")
	fmt.Printf("1 USD = %.4f EUR\n", USDToEUR)
	fmt.Printf("1 USD = %.4f RUB\n", USDToRUB)

	fmt.Println("\nCalculated conversion:")
	fmt.Printf("1 EUR = %.4f RUB\n", eurToRub)

	amountEUR := 50.0
	rub := amountEUR * eurToRub
	fmt.Printf("\nExample: %.2f EUR = %.2f RUB\n", amountEUR, rub)
}
