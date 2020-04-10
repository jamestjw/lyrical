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
	name      string
	arguments []string
}

// Matcher is a struct that will match commands from users
type Matcher struct {
	matchRegex *regexp.Regexp
	name       string
	arguments  []string
}

// NewMatcher creates a new matcher
func NewMatcher(name string, matchRegex string, arguments ...string) *Matcher {
	r := regexp.MustCompile(matchRegex)
	return &Matcher{
		matchRegex: r,
		name:       name,
		arguments:  arguments,
	}
}

func (e Error) Error() string {
	var prettyArgs []string
	for _, arg := range e.arguments {
		prettyArgs = append(prettyArgs, codify(arg))
	}

	joinedArguments := strings.Join(prettyArgs, ", ")
	return fmt.Sprintf("whoops %s requires the following parameter(s) %s 😅", codify(e.name), joinedArguments)
}

func codify(s string) string {
	return fmt.Sprintf("`%s`", s)
}

// Match will match a message to a regex
func (m *Matcher) Match(message string) (matched bool, arg string, err error) {
	matches := m.matchRegex.FindStringSubmatch(message)
	if matches == nil {
		matched = false
		return
	}

	matched = true

	if len(m.arguments) == 0 {
		return
	} else if len(matches) == 0 {
		log.Fatal("Expected argument but did not set it up in the regex.")
	}
	fmt.Println(matches)
	arg = strings.TrimSpace(matches[1])
	if arg == "" {
		err = Error{name: m.name, arguments: m.arguments}
	}
	return
}

// GetName returns name of the Matcher
func (m *Matcher) GetName() string {
	return m.name
}
