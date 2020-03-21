package ytmp3

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rylio/ytdl"
)

const audioPath = "audio-cache"

func init() {
	_ = os.Mkdir(audioPath, 0777)
}

// Download a MP3 file based on youtube ID
func Download(youtubeID string) {
	vid, err := ytdl.GetVideoInfo("https://www.youtube.com/watch?v=" + youtubeID)
	if err != nil {
		fmt.Println("Failed to get video info")
		return
	}
	ffmpeg, err := exec.LookPath("ffmpeg")
	if err != nil {
		fmt.Println("ffmpeg not found")
		return
	}

	videoFname := filepath.Join(audioPath, youtubeID+".mp4")
	mp3Fname := filepath.Join(audioPath, youtubeID+".mp3")
	file, _ := os.Create(videoFname)
	defer file.Close()
	defer os.Remove(videoFname)

	vid.Download(vid.Formats[0], file)

	fmt.Println("Video is ready.")
	cmd := exec.Command(ffmpeg, "-y", "-loglevel", "quiet", "-i", videoFname, "-vn", mp3Fname)
	if err := cmd.Run(); err != nil {
		fmt.Println("Failed to extract audio:", err)
	} else {
		fmt.Println("Extracted audio:", mp3Fname)
	}
}
