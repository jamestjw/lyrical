package utils

import (
	"errors"
	"time"
)

// VideoDurationValid parses the duration of a YouTube video
// and checks if it valid
func VideoDurationValid(videoDuration time.Duration) (err error) {
	if videoDuration.Minutes() > 10 {
		err = errors.New("video is more than 10 minutes long")
	}
	return
}
