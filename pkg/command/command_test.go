package command

import (
	"fmt"
	"testing"
)

func TestIsError1(t *testing.T) {
	output := Output{
		Error: nil,
	}

	actual := IsError(output)
	const expected = false

	if actual != expected {
		t.Errorf("Expected %t to be %t", actual, expected)
	}
}

func TestIsError2(t *testing.T) {
	output := Output{
		Error: fmt.Errorf("error"),
	}

	actual := IsError(output)
	const expected = true

	if actual != expected {
		t.Errorf("Expected %t to be %t", actual, expected)
	}
}
