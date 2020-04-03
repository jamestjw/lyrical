package matcher

var (
	JoinChannelRequestMatcher = NewMatcher("join-voice", "channel-name", `^!join-voice(\s+(.*)$)?`)
	AddPlaylistRequestMatcher = NewMatcher("add-playlist", "song-name", `^!(?:add-playlist|add-music)(\s+(.*)$)?`)
	SkipMusicMatcher          = NewMatcher("skip-music", "", `^!skip-music`)
	NowPlayingMatcher         = NewMatcher("now-playing", "", `^!now-playing`)
	StopMusicMatcher          = NewMatcher("stop-music", "", `^!stop-music`)
	PlayMusicMatcher          = NewMatcher("play-music", "", `^!play-music`)
	HelpMatcher               = NewMatcher("help", "", `^!help`)
	LeaveVoiceMatcher         = NewMatcher("leave-voice", "", `^!leave-voice`)
)
