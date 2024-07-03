package main

import (
	"errors"
	"fmt"
	"testing"
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

func getError() error {
	return errors.New("an error")
}

func ExampleReturningError() {
	err := getError()
	if err != nil {
		fmt.Println("Error:", err)
	}
	// Output: Error: an error
}

func TestReturningError(t *testing.T) {
	expectedError := errors.New("an error")

	err := getError()
	if err != expectedError {
		t.Errorf("Expected error to be '%v', but got '%v'", expectedError, err)
	}
}
