package matcher

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// JoinChannelRequestRe is a regex to match request to join voice channels
	JoinChannelRequestRe = regexp.MustCompile(`^!join-voice\s?(.*)$`)
)

// Match will match a message to a regex
func Match(matchRegex *regexp.Regexp, message string, name string, argument string) (matched bool, arg string, err error) {
	matches := matchRegex.FindStringSubmatch(message)
	if matches == nil {
		matched = false
		return
	}

	matched = true

	arg = strings.TrimSpace(matches[1])
	if arg == "" {
		err = fmt.Errorf("whoops `%s` requires a parameter of `%s` 😅", name, argument)
	}
	return
}
