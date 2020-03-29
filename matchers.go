package main

import "github.com/jamestjw/lyrical/matcher"

var (
	joinChannelRequestMatcher = matcher.NewMatcher("join-voice", "channel-name", `^!join-voice(\s+(.*)$)?`)
	addPlaylistRequestMatcher = matcher.NewMatcher("add-playlist", "song-name", `^!(?:add-playlist|add-music)(\s+(.*)$)?`)
	skipMusicMatcher          = matcher.NewMatcher("skip-music", "", `^!skip-music`)
	nowPlayingMatcher         = matcher.NewMatcher("now-playing", "", `^!now-playing`)
	stopMusicMatcher          = matcher.NewMatcher("stop-music", "", `^!stop-music`)
	playMusicMatcher          = matcher.NewMatcher("play-music", "", `^!play-music`)
	helpMatcher               = matcher.NewMatcher("help", "", `^!help`)
	leaveVoiceMatcher         = matcher.NewMatcher("leave-voice", "", `^!leave-voice`)
)
