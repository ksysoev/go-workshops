package main

import (
	"errors"
	"testing"
)

func getError() error {
	return errors.New("an error")
}

func TestReturningError(t *testing.T) {
	expectedError := errors.New("an error")

	err := getError()
	if err != expectedError {
		t.Errorf("Expected error to be '%v', but got '%v'", expectedError, err)
	}
}
