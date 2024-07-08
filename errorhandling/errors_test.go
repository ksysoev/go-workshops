package errorhandling

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgconn"
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
// when to wrap an error:
// - when we need to add more context to the error
// - when we get an error from a third-party library it's a good practice to wrap it with a more descriptive error

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
	err = fmt.Errorf("Fail to fetch user %d for update: %w", userID, err)

	if err == ErrUserNotFound {
		fmt.Println("Error:", err)
	}

	err = ValidateField("username", "verylongvalue")
	err = fmt.Errorf("Fail to validate signup form: %w", err)
	if err != nil {
		fmt.Println("Error:", err)
		// 	fmt.Println("Field:", fieldErr.Field)
		// 	fmt.Println("Message:", fieldErr.Msg)
	}

	// Output:
	// Error: user not found
	// Field: username
	// Message: value is too long
}

// Let's imaging we have multiple errors and we want to return them all.
// If we know exactly how many errors we have, we can use fmt.Errorf() function with %w verb.

func ExampleJoiningErrors1() {
	_ = errors.New("error1")
	_ = errors.New("error2")
	err := fmt.Errorf("multiple errors")

	fmt.Println("Error:", err)

	// Output:
	// Error: multiple errors: error1, error2
}

// If we don't know how many errors we have, we can use errors.Join() function to achieve the simular result.

func ExampleJoiningErrors2() {
	_ = errors.New("error1")
	_ = errors.New("error2")

	var errs error

	fmt.Println("Error:", errs)

	// Output:
	// Errors: error1
	// error2
}

// In addition to standard errors, Go provides a way to throw exceptions like errors using panic() function.
// to catch exceptions we can use recover() function.
// But Rule of thumb: DON'T PANIC and always carry a towel. :)
// But if we have panic and recover as part of the language, probably there are some valid use cases.
// Let's try to think what are they:
// 1. panic when it's clearly a developer mistake and it could put the application in incorect state.

// here is an example from standard library:
// - Creating ticker with duration that doesn't make any sense https://cs.opensource.google/go/go/+/refs/tags/go1.22.5:src/time/tick.go;l=20-23
// - Creating context from parent context, but got nil instead https://cs.opensource.google/go/go/+/master:src/context/context.go;l=269-271
// - Creating a new channel with a negative buffer size https://cs.opensource.google/go/go/+/refs/tags/go1.22.5:src/runtime/chan.go;l=33-35

// But if your function already returns an error, it's better to return an error instead of panicking even if it's a developer mistake.

// 2. When you have a function with deep call stack and you want to stop the execution and return an error to the top level function.
//    for this case it's crucial to recover the panic and return an error at top level of your function.
//    The rule of thumb here, internal panic should never cross boundaries of your package.

// 3. When we fail to initialize dependncies that are crucial for the application logic.
//    For example:
//    - Regular expression pattern is invalid https://cs.opensource.google/go/go/+/master:src/net/http/cgi/host.go;l=36?q=MustCompile&ss=go%2Fgo:src%2Fnet%2F&start=1
//    - Some template is crucial for the application to work is missing or invalid

// Let's try to fix the code below by using panic() and recover() functions to pass the test.

func ExamplePanicAndRecover() {
	panic("something went wrong")

	// Output:
	// Panic: something went wrong
}

// I think we mostly covered the error handling in Go.
// Let's try to consider some pitfalls that could be related to error handling.

// let't  try to analyse the code below and find what's wrong with it.

// Pitfall 1:

type Client struct {
	Name string
	Age  uint16
}

type InvalidClientError struct {
	Msg string
}

func (e *InvalidClientError) Error() string {
	return e.Msg
}

func ValidateClient(client Client) error {
	var err *InvalidClientError

	if client.Name == "" {
		err = &InvalidClientError{Msg: "name is required"}
	} else if client.Age < 18 {
		err = &InvalidClientError{Msg: "age should be greater than 18"}
	}

	return err
}

func ExampleReturningNilInterface() {
	client := Client{Name: "Vasia Pupkin", Age: 42}
	err := ValidateClient(client)

	if err != nil {
		fmt.Println("Error:", err)
	}

	// Output:
}

// Pitfall 2:
func logError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func ExampleLoggingErrors() {
	var err error
	defer logError(err)

	err = errors.New("something went wrong")

	// Output:
	// Error: something went wrong
}

// Pitfall 3:
func ExampleHandllingDbError() {
	err := GetUsers()

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			fmt.Println(pgErr.Message)
			fmt.Println(pgErr.Code)
		}
	}

	// Output:
	// relation "users" does not exist
	// 42P01
}
