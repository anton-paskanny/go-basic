package main

import (
	"fmt"
	"strings"
)

// Exchange rates map where key is "FROM_TO" format and value is the exchange rate
var exchangeRates = map[string]float64{
	"USD_EUR": 0.92,
	"USD_RUB": 94.50,
	"EUR_USD": 1.09,
	"EUR_RUB": 102.89,
	"RUB_USD": 0.0106,
	"RUB_EUR": 0.0097,
}

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

// convertCurrency converts from one currency to another using the exchange rates map pointer
func convertCurrency(amount float64, from string, to string, rates *map[string]float64) float64 {
	from = strings.ToUpper(from)
	to = strings.ToUpper(to)

	if from == to {
		return amount
	}

	// Create the key for the exchange rate map
	rateKey := from + "_" + to

	// Check if the exchange rate exists in the map
	if rate, exists := (*rates)[rateKey]; exists {
		return amount * rate
	}

	// If the direct conversion doesn't exist, show error
	fmt.Printf("Error: conversion from %s to %s is not supported\n", from, to)
	return 0
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

	result := convertCurrency(amount, sourceCurrency, targetCurrency, &exchangeRates)
	showResult(amount, sourceCurrency, targetCurrency, result)
}
