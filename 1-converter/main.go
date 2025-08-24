package main

import (
	"fmt"
	"strings"
)

const (
	USDToEUR = 0.92
	USDToRUB = 94.50
	EURToUSD = 1.09
	EURToRUB = 102.89
	RUBToUSD = 0.0106
	RUBToEUR = 0.0097
)

var availableCurrencies = []string{"USD", "EUR", "RUB"}

func showWelcome() {
	fmt.Println("=== Currency Converter ===")
	fmt.Printf("Available currencies: %s\n", strings.Join(availableCurrencies, ", "))
	fmt.Println()
}

func showResult(amount float64, from, to string, result float64) {
	fmt.Println("\n=== Conversion Result ===")
	fmt.Printf("%.2f %s = %.2f %s\n", amount, from, result, to)
}

func isValidCurrency(currency string) bool {
	for _, c := range availableCurrencies {
		if c == currency {
			return true
		}
	}
	return false
}

func clearInputBuffer() {
	var discard string
	fmt.Scanln(&discard)
}

// readInput reads a float number from user input
func readAmount(prompt string) float64 {
	var value float64

	for {
		fmt.Print(prompt)
		_, err := fmt.Scanln(&value)

		if err != nil {
			fmt.Println("Invalid input, please enter a valid number")

			// Clear the input buffer
			var discard string
			fmt.Scanln(&discard)
			continue
		}

		if value > 0 {
			break
		}

		fmt.Println("Invalid amount, please try again")
	}

	return value
}

// convertCurrency is a stub function that will convert from one currency to another
func convertCurrency(amount float64, from string, to string) float64 {
	from = strings.ToUpper(from)
	to = strings.ToUpper(to)

	if from == to {
		return amount
	}

	switch {
	case from == "USD" && to == "EUR":
		return amount * USDToEUR
	case from == "USD" && to == "RUB":
		return amount * USDToRUB
	case from == "EUR" && to == "USD":
		return amount * EURToUSD
	case from == "EUR" && to == "RUB":
		return amount * EURToRUB
	case from == "RUB" && to == "USD":
		return amount * RUBToUSD
	case from == "RUB" && to == "EUR":
		return amount * RUBToEUR
	default:
		fmt.Printf("Error: conversion from %s to %s is not supported", from, to)
		return 0
	}
}

func readCurrency(prompt, excludeCurrency string) string {
	var currency string

	for {
		fmt.Print(prompt)

		_, err := fmt.Scanln(&currency)

		if err != nil {
			fmt.Println("Invalid input, please enter a valid currency")
			clearInputBuffer()
			continue
		}

		currency = strings.ToUpper(currency)

		if !isValidCurrency(currency) {
			fmt.Printf("Invalid currency '%s'. Available currencies: %s\n",
				currency, strings.Join(availableCurrencies, ", "))
			continue
		}

		if excludeCurrency != "" && currency == strings.ToUpper(excludeCurrency) {
			fmt.Println("Target currency cannot be the same as source currency, please try again")
			continue
		}

		break
	}

	return currency
}

func main() {
	showWelcome()

	sourceCurrency := readCurrency(
		fmt.Sprintf("Enter source currency (%s): ", strings.Join(availableCurrencies, ", ")),
		"")

	amount := readAmount("Enter amount: ")

	targetCurrency := readCurrency(
		fmt.Sprintf("Enter target currency (%s, not %s): ",
			strings.Join(availableCurrencies, ", "), sourceCurrency),
		sourceCurrency)

	result := convertCurrency(amount, sourceCurrency, targetCurrency)
	showResult(amount, sourceCurrency, targetCurrency, result)
}
