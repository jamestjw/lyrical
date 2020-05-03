package playlist

import (
	"fmt"
	"strings"
)

// FormatNowPlayingText accepts a list of songs and formats it into a string
// that contains a list of songs names.
func FormatNowPlayingText(songs []*Song, header string) string {
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
func LimitSongsArrayLengths(songs1 []*Song, songs2 []*Song, limit int) []*Song {
	res := []*Song{}
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
