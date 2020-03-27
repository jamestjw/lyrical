package matcher

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// JoinChannelRequestRe is a regex to match request to join voice channels
	JoinChannelRequestRe = regexp.MustCompile(`^!join-voice(\s+(.*)$)?`)
	// AddPlaylistRequestRe is a regex to match request to add songs to playlists
	AddPlaylistRequestRe = regexp.MustCompile(`^!(?:add-playlist|add-music)(\s+(.*)$)?`)
)

// Error is the error returned from matcher when invalid
// parameters are given
type Error struct {
	name     string
	argument string
}

func (e Error) Error() string {
	return fmt.Sprintf("whoops `%s` requires a parameter of `%s` 😅", e.name, e.argument)
}

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
		err = Error{name: name, argument: argument}
	}
	return
}
