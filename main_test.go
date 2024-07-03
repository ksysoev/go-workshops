package main

import (
	"fmt"
)

// Go standard library provides 2 ways to create errors:
// 1. errors.New() function
// 2. fmt.Errorf() function
// Both functions return an error interface, which is a built-in interface in Go.
// Let's try to fix the code below by creating errors using errors.New() and fmt.Errorf() functions.
func ExampleCreatingErrors() {
	var err error
	// err :=
	if err != nil {
		fmt.Println("Error1:", err)
	}

	// expectedValue := 42
	// err =
	if err != nil {
		fmt.Println("Error2:", err)
	}
	// Output:
	// Error1: an error
	// Error2: an error with value 42
}

// Now, let's try to return an error from a function.

// ValidatePasswordLen function validates the length of the password.
// It returns an error if the password is shorter than 8 characters.
func ValidatePasswordLen(password string) error {
	return nil
}

func ExampleReturningError() {
	if err := ValidatePasswordLen("longpassword"); err != nil {
		fmt.Println("Error:", err)
	}

	if err := ValidatePasswordLen("short"); err != nil {
		fmt.Println("Error2:", err)
	}

	// Output:
	// Error2: password is too short: short
}

// What if we need to return a value along or an error from the function?
// by convention, first we return success values and last we return error.
// Let's try to fix the code below by returning a value along with an error from the function.

// Divide function divides two numbers.
// It returns the result of the division and an error if the denominator is zero.
func Divide(a, b int) (int, error) {
	return a / b, nil
}

func ExampleReturningValueAndError() {
	if result, err := Divide(10, 2); err != nil {
		fmt.Println("Error1:", err)
	} else {
		fmt.Println("Result1:", result)
	}

	if result, err := Divide(10, 0); err != nil {
		fmt.Println("Error2:", err)
	} else {
		fmt.Println("Result2:", result)
	}

	// Output:
	// Result: 5
	// Error2: division by zero
}
