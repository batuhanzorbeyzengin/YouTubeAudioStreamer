FROM golang:1.21-bullseye

# Install wget and ffmpeg
RUN apt-get update && \
    apt-get install -y wget ffmpeg

# Install yt-dlp
RUN wget https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -O /usr/local/bin/yt-dlp \
    && chmod a+rx /usr/local/bin/yt-dlp

WORKDIR /app

COPY . .

# Build the Go application
RUN go build -o youtube-audio-server

# Expose port 8080
EXPOSE 8080

# Command to run the application
CMD ["./youtube-audio-server"]
