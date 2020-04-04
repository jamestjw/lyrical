package ytmp3

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jamestjw/lyrical/utils"
	"github.com/jamestjw/ytdl"
)

// AudioPath is the path that contains all audio files
const AudioPath = "audio-cache"

func init() {
	_ = os.Mkdir(AudioPath, 0777)
}

// Download a MP3 file based on youtube ID
func Download(youtubeID string) (title string, err error) {
	vid, err := ytdl.GetVideoInfo("https://www.youtube.com/watch?v=" + youtubeID)
	if err != nil {
		log.Println("Failed to get video info")
		return "", errors.New("video ID is invalid")
	}
	ffmpeg, err := exec.LookPath("ffmpeg")
	if err != nil {
		log.Println("ffmpeg not found")
		return "", err
	}

	title = vid.Title
	if err := utils.VideoDurationValid(vid.Duration); err != nil {
		return "", err
	}

	videoFname := filepath.Join(AudioPath, youtubeID+".mp4")
	mp3Fname := filepath.Join(AudioPath, youtubeID+".mp3")
	file, _ := os.Create(videoFname)
	defer file.Close()
	defer os.Remove(videoFname)

	vid.Download(vid.Formats[0], file)

	log.Println("Video is ready.")
	cmd := exec.Command(ffmpeg, "-y", "-loglevel", "quiet", "-i", videoFname, "-vn", mp3Fname)
	if err := cmd.Run(); err != nil {
		log.Println("Failed to extract audio:", err)
		return "", err
	}
	log.Println("Extracted audio:", mp3Fname)
	return title, nil
}

// PathToAudio returns a path to an audio file
func PathToAudio(youtubeID string) string {
	return filepath.Join(AudioPath, youtubeID+".mp3")
}
