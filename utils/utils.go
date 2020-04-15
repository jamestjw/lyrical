package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jamestjw/lyrical/playlist"
	log "github.com/sirupsen/logrus"
)

// VideoDurationValid parses the duration of a YouTube video
// and checks if it valid
func VideoDurationValid(videoDuration time.Duration) (err error) {
	if videoDuration.Minutes() > 10 {
		err = errors.New("video is more than 10 minutes long")
	}
	return
}

// FormatNowPlayingText accepts a list of songs and formats it into a string
// that contains a list of songs names.
func FormatNowPlayingText(songs []*playlist.Song, header string) string {
	message := []string{header}

	for i, song := range songs {
		message = append(message, fmt.Sprintf("%v. %s", i+1, song.Name))
	}

	return strings.Join(message, "\n")
}

// LimitSongsArrayLengths takes two arrays of songs and returns only a single array
// with a limited number of elements. This function will first try to populate the
// resulting from elements in the first array. You many not get as many elements as you
// ask for.
func LimitSongsArrayLengths(songs1 []*playlist.Song, songs2 []*playlist.Song, limit int) []*playlist.Song {
	res := []*playlist.Song{}
	var takeFromFirst int
	var takeFromSecond int

	if len(songs1) >= limit {
		takeFromFirst = limit
		takeFromSecond = 0
	} else {
		takeFromFirst = len(songs1)
		if len(songs2) >= (limit - len(songs1)) {
			takeFromSecond = limit - len(songs1)
		} else {
			takeFromSecond = len(songs2)
		}
	}

	for _, song := range songs1[:takeFromFirst] {
		res = append(res, song)
	}

	for _, song := range songs2[:takeFromSecond] {
		res = append(res, song)
	}
	return res
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

func LogInfo(m string, kvs []LogKV) {
	fields := log.Fields{}
	for _, kv := range kvs {
		fields[kv.key] = kv.value
	}
	log.WithFields(fields).Info(m)
}

func LogError(m string, kvs []LogKV) {
	fields := log.Fields{}
	for _, kv := range kvs {
		fields[kv.key] = kv.value
	}
	log.WithFields(fields).Error(m)
}

func KvForHandler(guild string, handler string, kvs []LogKV) []LogKV {
	kvs = append(kvs, LogKV{"guildID", guild}, LogKV{"event", handler})

	return kvs
}

func KvForEvent(event string, kvs []LogKV) []LogKV {
	kvs = append(kvs, LogKV{"event", event})

	return kvs
}

func SingleKV(k string, v string) []LogKV {
	return []LogKV{LogKV{k, v}}
}

func KVs(kvs ...string) []LogKV {
	var res []LogKV
	for i := 0; i < len(kvs); i += 2 {
		res = append(res, LogKV{kvs[i], kvs[i+1]})
	}
	return res
}

type LogKV struct {
	key   string
	value string
}
