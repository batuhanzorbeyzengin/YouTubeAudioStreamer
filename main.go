package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var audioFiles = make(map[string]string)

type SongInfo struct {
	UniqueID    string `json:"unique_id"`
	FilePath    string `json:"file_path"`
	DownloadURL string `json:"download_url"`
}

func main() {
	if err := os.Mkdir("audio", os.ModePerm); err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	http.HandleFunc("/download", handleDownload)
	http.HandleFunc("/audio/", handleAudio)
	http.HandleFunc("/songs", handleSongList)
	http.HandleFunc("/song", handleSongDetails)
	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	videoURL := r.URL.Query().Get("url")
	if videoURL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	videoURL = formatYouTubeURL(videoURL)
	audioFilePath, err := downloadAndExtractAudio(videoURL)
	if err != nil {
		http.Error(w, "Error processing video: "+err.Error(), http.StatusInternalServerError)
		return
	}

	uniqueID := filepath.Base(audioFilePath)
	audioFiles[uniqueID] = audioFilePath

	audioURL := "http://" + r.Host + "/audio/" + uniqueID
	songInfo := SongInfo{
		UniqueID:    uniqueID,
		FilePath:    audioFilePath,
		DownloadURL: audioURL,
	}

	saveSongInfo(songInfo)

	io.WriteString(w, audioURL)
}

func saveSongInfo(info SongInfo) {
	var allSongs []SongInfo
	data, err := os.ReadFile("songs.json")
	if err == nil {
		json.Unmarshal(data, &allSongs)
	}
	allSongs = append(allSongs, info)

	newData, err := json.Marshal(allSongs)
	if err != nil {
		log.Printf("Error saving song info: %s\n", err)
		return
	}
	os.WriteFile("songs.json", newData, 0644)
}

func handleSongDetails(w http.ResponseWriter, r *http.Request) {
	uniqueID := r.URL.Query().Get("id")
	if uniqueID == "" {
		http.Error(w, "Unique ID is required", http.StatusBadRequest)
		return
	}

	var allSongs []SongInfo
	data, err := os.ReadFile("songs.json")
	if err != nil {
		http.Error(w, "Could not read song data", http.StatusInternalServerError)
		return
	}
	json.Unmarshal(data, &allSongs)

	for _, song := range allSongs {
		if song.UniqueID == uniqueID {
			responseData, _ := json.Marshal(song)
			w.Header().Set("Content-Type", "application/json")
			w.Write(responseData)
			return
		}
	}
	http.NotFound(w, r)
}

func handleSongList(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("songs.json")
	if err != nil {
		http.Error(w, "Could not read song data", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func handleAudio(w http.ResponseWriter, r *http.Request) {
	uniqueID := strings.TrimPrefix(r.URL.Path, "/audio/")
	audioFilePath, exists := audioFiles[uniqueID]
	if !exists {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, audioFilePath)
}

func downloadAndExtractAudio(videoURL string) (string, error) {
	audioFilePath := "audio/" + generateUniqueID() + ".mp3"
	log.Println("Downloading audio from URL:", videoURL)
	if err := runCmd("yt-dlp", "-x", "--audio-format", "mp3", "-o", audioFilePath, videoURL); err != nil {
		return "", err
	}

	return audioFilePath, nil
}

func formatYouTubeURL(videoID string) string {
	if !strings.Contains(videoID, "youtu.be") {
		return "https://www.youtube.com/watch?v=" + videoID
	}
	return videoID
}

func generateUniqueID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
