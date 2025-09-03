package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

var operations = map[string]struct {
	calc  func([]float64) float64
	label string
}{
	"SUM": {calc: calculateSum, label: "Sum"},
	"AVG": {calc: calculateAverage, label: "Arithmetic mean"},
	"MED": {calc: calculateMedian, label: "Median"},
}

func operationKeys() []string {
	keys := make([]string, 0, len(operations))
	for k := range operations {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func isValidOperation(operation string) bool {
	operation = strings.ToUpper(operation)
	_, ok := operations[operation]
	return ok
}

func readOperation() string {
	var operation string
	for {
		fmt.Printf("Enter operation (%s): ", strings.Join(operationKeys(), ", "))
		fmt.Scanln(&operation)

		operation = strings.ToUpper(strings.TrimSpace(operation))

		if isValidOperation(operation) {
			return operation
		}

		fmt.Printf("Error: operation '%s' is not supported. Available operations: %s\n",
			operation, strings.Join(operationKeys(), ", "))
	}
}

func parseNumbers(input string) ([]float64, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty string")
	}

	parts := strings.Split(input, ",")
	numbers := make([]float64, 0, len(parts))

	for i, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			return nil, fmt.Errorf("empty value at position %d", i+1)
		}

		num, err := strconv.ParseFloat(part, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number '%s' at position %d", part, i+1)
		}

		numbers = append(numbers, num)
	}

	if len(numbers) == 0 {
		return nil, fmt.Errorf("no numbers entered")
	}

	return numbers, nil
}

func readNumbers() []float64 {
	var input string
	for {
		fmt.Print("Enter numbers separated by commas (e.g., 2,10,9 - no spaces between numbers): ")
		fmt.Scanln(&input)

		numbers, err := parseNumbers(input)
		if err != nil {
			fmt.Printf("Error: %v. Please try again.\n", err)
			continue
		}

		return numbers
	}
}

func calculateSum(numbers []float64) float64 {
	sum := 0.0
	for _, num := range numbers {
		sum += num
	}
	return sum
}

func calculateAverage(numbers []float64) float64 {
	if len(numbers) == 0 {
		return 0
	}
	return calculateSum(numbers) / float64(len(numbers))
}

func calculateMedian(numbers []float64) float64 {
	if len(numbers) == 0 {
		return 0
	}

	sorted := make([]float64, len(numbers))
	copy(sorted, numbers)
	sort.Float64s(sorted)

	length := len(sorted)

	if length%2 == 0 {
		mid1 := sorted[length/2-1]
		mid2 := sorted[length/2]
		return (mid1 + mid2) / 2.0
	} else {
		return sorted[length/2]
	}
}

// map-based dispatch replaces performCalculation

func formatNumbers(numbers []float64) string {
	strNumbers := make([]string, len(numbers))
	for i, num := range numbers {
		if num == float64(int64(num)) {
			strNumbers[i] = fmt.Sprintf("%.0f", num)
		} else {
			strNumbers[i] = fmt.Sprintf("%.2f", num)
		}
	}
	return strings.Join(strNumbers, ", ")
}

func formatResult(result float64) string {
	if result == float64(int64(result)) {
		return fmt.Sprintf("%.0f", result)
	}
	return fmt.Sprintf("%.2f", result)
}

func showWelcome() {
	fmt.Println("=== MATHEMATICAL CALCULATOR ===")
	fmt.Println("Supported operations:")
	fmt.Println("  SUM - sum of numbers")
	fmt.Println("  AVG - arithmetic mean")
	fmt.Println("  MED - median")
	fmt.Println()
}

func askForContinue() bool {
	var answer string
	fmt.Print("\nDo you want to perform another operation? (y/n): ")
	fmt.Scanln(&answer)
	answer = strings.ToLower(strings.TrimSpace(answer))
	return answer == "y" || answer == "yes"
}

func main() {
	showWelcome()

	for {
		operation := readOperation()

		numbers := readNumbers()

		result := operations[operation].calc(numbers)

		fmt.Println("\n=== RESULT ===")
		fmt.Printf("Operation: %s\n", operation)
		fmt.Printf("Numbers: %s\n", formatNumbers(numbers))

		fmt.Printf("%s: %s\n", operations[operation].label, formatResult(result))

		fmt.Printf("Number count: %d\n", len(numbers))

		if !askForContinue() {
			break
		}

		fmt.Println()
	}

	fmt.Println("Thank you for using the calculator!")
}
