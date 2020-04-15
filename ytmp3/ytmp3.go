package ytmp3

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"

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
	utils.LogInfo("download", utils.KvForEvent("ytmp3", utils.KVs("youtubeID", youtubeID)))

	vid, err := ytdl.GetVideoInfo("https://www.youtube.com/watch?v=" + youtubeID)
	if err != nil {
		log.Error("Failed to get video info: " + youtubeID)
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

	utils.LogInfo("video_ready", utils.KVs("name", title, "youtubeID", youtubeID, "event", "ytmp3"))

	cmd := exec.Command(ffmpeg, "-y", "-loglevel", "quiet", "-i", videoFname, "-vn", mp3Fname)
	if err := cmd.Run(); err != nil {
		utils.LogError("failed audio_extraction", utils.KVs("err", err.Error(), "youtubeID", youtubeID, "event", "ytmp3"))
		return "", err
	}

	utils.LogInfo("audio_extracted", utils.KVs("filename", mp3Fname, "youtubeID", youtubeID, "event", "ytmp3"))

	return title, nil
}

// PathToAudio returns a path to an audio file
func PathToAudio(youtubeID string) string {
	return filepath.Join(AudioPath, youtubeID+".mp3")
}
