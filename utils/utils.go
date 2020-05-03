package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// VideoDurationValid parses the duration of a YouTube video
// and checks if it valid
func VideoDurationValid(videoDuration time.Duration) (err error) {
	if videoDuration.Minutes() > 10 {
		err = errors.New("video is more than 10 minutes long")
	}
	return
}

func ExtractEmojiCounts(reactions []*discordgo.MessageReactions) map[string]int {
	res := make(map[string]int)

	for _, reaction := range reactions {
		res[reaction.Emoji.Name] = reaction.Count
	}

	return res
}

func Bold(s string) string {
	return fmt.Sprintf("**%s**", s)
}

func StringArrayMap(strs []string, f func(string) string) []string {
	var res []string

	for _, s := range strs {
		res = append(res, f(s))
	}

	return res
}
