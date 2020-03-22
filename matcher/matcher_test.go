package matcher

import (
	"reflect"
	"regexp"
	"testing"
)

const checkMark = "\u2713"
const ballotX = "\u2717"

var (
	testRe = regexp.MustCompile(`^!test-command(\s+(.*)$)?`)
)

// TestMatch validates that the Matcher function matches a command
// and returns the right arguments, or produces the right error when
// a required argument is not given.
func TestMatch(t *testing.T) {
	command := "test-command"
	argumentName := "expected-arg"

	tables := []struct {
		input           string
		expectedMatched bool
		expectedArg     string
		expectedErr     error
	}{
		{"!test-command test-arg", true, "test-arg", *new(error)},
		{"!test-commandx test-arg", true, "", Error{}},
		{"!wrong-test-command test-arg", false, "", *new(error)},
	}

	t.Log("Given the need to match commands and return right arguments or errors")
	for _, table := range tables {
		matched, arg, err := Match(testRe, table.input, command, argumentName)
		{
			t.Logf("\tWhen checking \"%s\" for command  \"%s\" and argument name \"%s\"", table.input, command, argumentName)
			if matched != table.expectedMatched {
				t.Error("\t\tShould be able identify match", ballotX)
			} else {
				t.Log("\t\tShould be able identify match", checkMark)
			}

			if arg != table.expectedArg {
				t.Error("\t\tShould be able identify arg", ballotX)
				t.Errorf("\t\t\tExpected: %v, Got: %v", table.expectedArg, arg)
			} else {
				t.Log("\t\tShould be able identify arg", checkMark)
			}

			if !IsInstanceOf(err, table.expectedErr) {
				t.Error("\t\tShould return right error", ballotX, err)
			} else {
				t.Log("\t\tShould return right error", checkMark)
			}
		}
	}
}

func IsInstanceOf(objectPtr, typePtr interface{}) bool {
	return reflect.TypeOf(objectPtr) == reflect.TypeOf(typePtr)
}
