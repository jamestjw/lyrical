package matcher

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

// Error is the error returned from matcher when invalid
// parameters are given
type Error struct {
	name     string
	argument string
}

// Matcher is a struct that will match commands from users
type Matcher struct {
	matchRegex *regexp.Regexp
	name       string
	argument   string
}

// NewMatcher creates a new matcher
func NewMatcher(name string, argument string, matchRegex string) *Matcher {
	r := regexp.MustCompile(matchRegex)
	return &Matcher{
		matchRegex: r,
		name:       name,
		argument:   argument,
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("whoops `%s` requires a parameter of `%s` 😅", e.name, e.argument)
}

// Match will match a message to a regex
func (m *Matcher) Match(message string) (matched bool, arg string, err error) {
	matches := m.matchRegex.FindStringSubmatch(message)
	if matches == nil {
		matched = false
		return
	}

	matched = true

	if m.argument == "" {
		return
	} else if len(matches) == 0 {
		log.Fatal("Expected argument but did not set it up in the regex.")
	}

	arg = strings.TrimSpace(matches[1])
	if arg == "" {
		err = Error{name: m.name, argument: m.argument}
	}
	return
}

// GetName returns name of the Matcher
func (m *Matcher) GetName() string {
	return m.name
}
