package matcher

var (
	JoinChannelRequestMatcher = NewMatcher("join-voice", `^!join-voice(\s+(.*)$)?`, "channel-name")
	AddPlaylistRequestMatcher = NewMatcher("add-playlist", `^!(?:add-playlist|add-music)(\s+(.*)$)?`, "song-name")
	SkipMusicMatcher          = NewMatcher("skip-music", `^!skip-music`)
	NowPlayingMatcher         = NewMatcher("now-playing", `^!now-playing`)
	StopMusicMatcher          = NewMatcher("stop-music", `^!stop-music`)
	PlayMusicMatcher          = NewMatcher("play-music", `^!play-music`)
	HelpMatcher               = NewMatcher("help", `^!help`)
	LeaveVoiceMatcher         = NewMatcher("leave-voice", `^!leave-voice`)
	UpNextMatcher             = NewMatcher("up-next", `^!up-next`)
	VoteMatcher               = NewMatcher("poll", `^!poll(\s+(.*)$)?`, "title", "option1", "option2", "...")
)
