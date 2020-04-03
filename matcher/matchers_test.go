package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Example struct {
	input           string
	expectedMatched bool
	expectedArg     string
	errorExpected   bool
}

func TestJoinChannelRequestMatcher(t *testing.T) {
	matcher := JoinChannelRequestMatcher

	tables := []Example{
		{"!join-voice channel name", true, "channel name", false},
		{"!join-voice    ", true, "", true},
		{"!unrelated-join-voice test-arg", false, "", false},
	}

	tableTest(t, tables, matcher)
}

func TestAddPlaylistRequestMatcher(t *testing.T) {
	matcher := AddPlaylistRequestMatcher

	tables := []Example{
		{"!add-playlist", true, "", true},
		{"!add-playlist  ", true, "", true},
		{"!add-playlist songname", true, "songname", false},
		{"!add-music", true, "", true},
		{"!add-music  ", true, "", true},
		{"!add-music songname", true, "songname", false},
		{"!join-voice    ", false, "", false},
		{"!unrelated-join-voice test-arg", false, "", false},
	}

	tableTest(t, tables, matcher)
}

func TestSkipMusicMatcher(t *testing.T) {
	matcher := SkipMusicMatcher

	tables := []Example{
		{"!skip-music", true, "", false},
		{"!skip-music    ", true, "", false},
		{"!skip-music  useless arguments", true, "", false},
		{"!unrelated-join-voice test-arg", false, "", false},
	}

	tableTest(t, tables, matcher)
}

func TestNowPlayingtMatcher(t *testing.T) {
	matcher := NowPlayingMatcher

	tables := []Example{
		{"!now-playing", true, "", false},
		{"!now-playing    ", true, "", false},
		{"!now-playing  useless arguments", true, "", false},
		{"!unrelated-join-voice test-arg", false, "", false},
	}

	tableTest(t, tables, matcher)
}

func TestStopMusicMatcher(t *testing.T) {
	matcher := StopMusicMatcher

	tables := []Example{
		{"!stop-music", true, "", false},
		{"!stop-music    ", true, "", false},
		{"!stop-music  useless arguments", true, "", false},
		{"!unrelated-join-voice test-arg", false, "", false},
	}

	tableTest(t, tables, matcher)
}

func TestPlayMusicMatcher(t *testing.T) {
	matcher := PlayMusicMatcher

	tables := []Example{
		{"!play-music", true, "", false},
		{"!play-music    ", true, "", false},
		{"!play-music  useless arguments", true, "", false},
		{"!unrelated-join-voice test-arg", false, "", false},
	}

	tableTest(t, tables, matcher)
}

func TestHelpMatcher(t *testing.T) {
	matcher := HelpMatcher

	tables := []Example{
		{"!help", true, "", false},
		{"!help    ", true, "", false},
		{"!help  useless arguments", true, "", false},
		{"!unrelated-join-voice test-arg", false, "", false},
	}

	tableTest(t, tables, matcher)
}

func TestLeaveVoiceMatcher(t *testing.T) {
	matcher := LeaveVoiceMatcher

	tables := []Example{
		{"!leave-voice", true, "", false},
		{"!leave-voice    ", true, "", false},
		{"!leave-voice  useless arguments", true, "", false},
		{"!unrelated-join-voice test-arg", false, "", false},
	}

	tableTest(t, tables, matcher)
}

func tableTest(t *testing.T, tables []Example, matcher *Matcher) {
	for _, table := range tables {
		matched, arg, err := matcher.Match(table.input)
		{
			assert.Equal(t, matched, table.expectedMatched, "query is matched")
			assert.Equal(t, arg, table.expectedArg, "expected argument is returned")

			if table.errorExpected {
				assert.Error(t, err, "error is returned")
			} else {
				assert.Nil(t, err, "no error is returned")
			}
		}
	}
}
