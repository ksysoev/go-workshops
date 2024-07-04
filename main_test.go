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
	if result, err := Divide(10, 2); err == nil {
		fmt.Println("Result1:", result)
	}

	if _, err := Divide(10, 0); err != nil {
		fmt.Println("Error2:", err)
	}

	// Output:
	// Result1: 5
	// Error2: division by zero
}

// Expected flow errors are errors that are expected to happen in the normal flow of the program.
// For example:
// - SQL query returns record not found
// - while reading a file you reach the end of the file
// - user enters an invalid password
// - etc.

// To simplify the error handling of expected flow errors, we can define public variables for them.

// ErrUserNotFound is an error returned when a user is not found.
var ErrUserNotFound = errors.New("user not found")

// GetUser function returns a user by ID.
// It returns an error if the user is not found.
func GetUser(id int) (string, error) {
	return "", ErrUserNotFound
}

func TestExpectedFlowErrors(t *testing.T) {
	_, err := GetUser(1)
	if err != errors.New("record not found") {
		t.Errorf("expected ErrRecordNotFound, got %v", err)
	}
}

// To implement custom errors, we can create a new type that implements the error interface.
// error interface has only one method: Error() string
// Custom errors are useful when we need to add more context to the error.
// for example, we can add an error code, error message, etc.

// FieldValidationError is  field validation error.
type FieldValidationError struct {
	Field string
	Msg   string
}

// NewFieldValidationError creates a new field validation error.
func NewFieldValidationError(field, msg string) *FieldValidationError {
	return &FieldValidationError{
		Field: field,
		Msg:   msg,
	}
}

// ValidateField function validates a field value.
func ValidateField(field, value string) error {
	if len(value) > 10 {
		return fmt.Errorf("value is too long")
	}

	return nil
}

func ExampleCustomErrors() {
	err := ValidateField("username", "verylongvalue")

	if err != nil {
		fmt.Println("Error:", err)

		var field, msg string

		fmt.Println("Field:", field)
		fmt.Println("Message:", msg)
	}

	// Output:
	// Error: too long
	// Field: username
	// Message: value is too long
}

// Error wrapping is a technique to add more context to an error
// by wrapping the original error with a new error.
// To wrap an error, we can use fmt.Errorf() function with %w verb.

func ExampleErrorWrapping() {
	userID := 10
	_, err := GetUser(userID)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Output:
	// Error: Fail to fetch user 10 for update: user not found
}

// With error wrapping, we can create a chain of errors.
// To unwrap an error, we can use errors.Is() and errors.As() functions.

func ExampleErrorUnwrapping() {
	userID := 10
	_, err := GetUser(userID)
	fmt.Errorf("Fail to fetch user %d for update: %w", userID, err)

	if err == ErrUserNotFound {
		fmt.Println("Error:", err)
	}

	err = ValidateField("username", "verylongvalue")
	fmt.Errorf("Fail to validate signup form: %w", err)
	// if fieldErr, ok := err.(*FieldValidationError); ok {
	// 	fmt.Println("Field:", fieldErr.Field)
	// 	fmt.Println("Message:", fieldErr.Msg)
	// }

	// Output:
	// Error: user not found
	// Field: username
	// Message: value is too long
}
