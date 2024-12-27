package command

import (
	"errors"
	"testing"
)

func TestIsError1(t *testing.T) {
	t.Parallel()

	output := Output{ //nolint:exhaustruct
		Error: nil,
	}

	actual := IsError(output)

	const expected = false

	if actual != expected {
		t.Errorf("Expected %t to be %t", actual, expected)
	}
}

func TestIsError2(t *testing.T) {
	t.Parallel()

	output := Output{ //nolint:exhaustruct
		Error: errors.New("error"),
	}

	actual := IsError(output)

	const expected = true

	if actual != expected {
		t.Errorf("Expected %t to be %t", actual, expected)
	}
}
