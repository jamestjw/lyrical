package matcher

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

const checkMark = "\u2713"
const ballotX = "\u2717"

var (
	testCommandWithArgRe         = regexp.MustCompile(`^!test-command(\s+(.*)$)?`)
	testCommandWithArgAndAliasRe = regexp.MustCompile(`^!(?:test-command|alternative-command)(\s+(.*)$)?`)
	testCommandWithNoArgRe       = regexp.MustCompile(`^!test-command`)
)

func TestNewMatcher(t *testing.T) {
	m := NewMatcher("name", "!nowplaying", "arg")
	assert.IsType(t, &Matcher{}, m, "should receive new pointer to Matcher")
	assert.Equal(t, "name", m.GetName(), "should return name of matcher")
}

// TestMatch validates that the Matcher function matches a command
// and returns the right arguments, or produces the right error when
// a required argument is not given.
func TestMatchOneCommand(t *testing.T) {
	command := "test-command"
	argumentName := []string{"expected-arg"}

	tables := []struct {
		input           string
		expectedMatched bool
		expectedArg     string
		expectedErr     error
	}{
		{"!test-command test-arg", true, "test-arg", *new(error)},
		{"!test-command", true, "", Error{}},
		{"!test-commandx test-arg", true, "", Error{}},
		{"!wrong-test-command test-arg", false, "", *new(error)},
	}

	matcher := &Matcher{
		testCommandWithArgRe,
		command,
		argumentName,
	}

	t.Log("Given the need to match commands and return right arguments or errors")
	for _, table := range tables {
		matched, arg, err := matcher.Match(table.input)
		{
			t.Logf("\tWhen checking \"%s\" for a single command \"%s\" and argument name \"%s\"", table.input, command, argumentName)
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

func TestMatchTwoCommands(t *testing.T) {
	command := "test-command"
	argumentName := []string{"expected-arg"}

	tables := []struct {
		input           string
		expectedMatched bool
		expectedArg     string
		expectedErr     error
	}{
		{"!test-command", true, "", Error{}},
		{"!test-command args", true, "args", *new(error)},
		{"!alternative-command", true, "", Error{}},
	}

	matcher := &Matcher{
		testCommandWithArgAndAliasRe,
		command,
		argumentName,
	}

	t.Log("Given the need to match two possible commands and return right arguments or errors")
	for _, table := range tables {
		matched, arg, err := matcher.Match(table.input)
		{
			t.Logf("\tWhen checking \"%s\" a single command \"%s\" and argument name \"%s\"", table.input, command, argumentName)
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

func TestErrorString(t *testing.T) {
	name := "command-name"
	args := []string{"arg-1", "arg-2"}
	expectedError := "whoops `command-name` requires the following parameter(s) `arg-1`, `arg-2` 😅"

	e := Error{name, args}
	assert.Equal(t, expectedError, e.Error())
}

func TestMatcherWithTwoArgsErrorMessage(t *testing.T) {
	matcher := NewMatcher("command-name", `^!command-name(\s+(.*)$)?`, "arg-1", "arg-2")

	expectedError := "whoops `command-name` requires the following parameter(s) `arg-1`, `arg-2` 😅"
	matched, _, error := matcher.Match("!command-name")
	assert.True(t, matched)
	assert.Equal(t, expectedError, error.Error())
}

func TestMatchCommandWithNoArgument(t *testing.T) {
	command := "test-command"
	var argumentName []string

	tables := []struct {
		input           string
		expectedMatched bool
		expectedArg     string
		expectedErr     error
	}{
		{"!test-command", true, "", *new(error)},
		{"!test-command useless args", true, "", *new(error)},
		{"!alternative-command", false, "", *new(error)},
	}

	matcher := &Matcher{
		testCommandWithNoArgRe,
		command,
		argumentName,
	}

	t.Log("Given the need to match two possible commands and return right arguments or errors")
	for _, table := range tables {
		matched, arg, err := matcher.Match(table.input)
		{
			t.Logf("\tWhen checking \"%s\" a single command \"%s\" and argument name \"%s\"", table.input, command, argumentName)
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
