# YouTube Audio Streamer

YouTube Audio Streamer is a web application that allows users to download audio from YouTube videos and stream it via a unique URL. The application is built using Go and utilizes `yt-dlp` for downloading audio content.

## Features

- Download audio from YouTube videos.
- Stream audio via a unique URL.
- Support for various audio formats (default to MP3).
- Simple and intuitive user interface.

## Prerequisites

Before you begin, ensure you have met the following requirements:

- You have installed the latest version of [Go](https://golang.org/dl/).
- You have `yt-dlp` and `ffmpeg` installed on your system.

## Installation

To install YouTube Audio Streamer, follow these steps:

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/youtube-audio-streamer.git
   cd youtube-audio-streamer
   docker build youtube-audio-streamer .
   docker run -p 8080:8080 youtube-audio-streamer
   ```
2. Open a web browser and go to http://localhost:8080.

## Downloading a Song
1. Get the Video URL: Copy the URL of the YouTube video from which you want to download audio.
2. Download the Audio: Use the /download endpoint to download the audio.
   - Format: http://localhost:8080/download?url=[YouTube Video ID] 
   - Example: http://localhost:8080/download?url=dQw4w9WgXcQ
3. Receive Audio URL: Upon successful download, the service will return a URL where the audio file is accessible.

## Listing All Songs
- To see a list of all downloaded songs, access the /songs endpoint.
   - URL: http://localhost:8080/songs
