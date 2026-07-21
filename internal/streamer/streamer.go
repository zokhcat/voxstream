package streamer

import (
	"io"
	"log"
	"net"
	"os"
	"time"
)

const (
	TargetAddr     = "localhost:8080"
	ChunkSizeBytes = 320
	FrameInterval  = 20 * time.Millisecond
)

func StreamAudio(filePath string) {
	log.Println("Initializing VoxStream Telephony Client Simulator...")

	conn, err := net.Dial("tcp", TargetAddr)
	if err != nil {
		log.Fatalf("Failed to connect to proxy at %s: %v", TargetAddr, err)
	}
	defer conn.Close()

	log.Printf("Connected to Proxy at %s. Streaming initiated...", TargetAddr)

	var audioSource io.Reader

	if filePath != "" {
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatalf("Failed to open audio file %s: %v", filePath, err)
		}
		defer file.Close()
		log.Printf("Streaming audio from file: %s", filePath)
		audioSource = file
	} else {
		log.Println("No file provided. Generating synthetic 8kHz PCM audio frames...")
		audioSource = &syntheticPCMGenerator{}
	}

	ticker := time.NewTicker(FrameInterval)
	defer ticker.Stop()

	buf := make([]byte, ChunkSizeBytes)
	var totalFrames int
	var totalBytes int
	startTime := time.Now()

	for range ticker.C {
		_, err := io.ReadFull(audioSource, buf)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				log.Println("End of audio stream reached.")
			} else {
				log.Printf("Stream read error: %v", err)
			}
			break
		}

		n, err := conn.Write(buf)
		if err != nil {
			log.Printf("Failed to transmit frame over TCP socket: %v", err)
			break
		}

		totalFrames++
		totalBytes += n

		if totalFrames%50 == 0 {
			duration := time.Since(startTime).Truncate(time.Millisecond)
			log.Printf("[Streamer] Transmitted Frame #%04d | Pacing: %s | Total Sent: %d bytes",
				totalFrames, duration, totalBytes)
		}
	}

	callDuration := time.Since(startTime).Truncate(time.Millisecond)
	log.Printf("Stream completed | Simulated Call Duration: %s | Total Frames: %d",
		callDuration, totalFrames)
}

type syntheticPCMGenerator struct{}

func (g *syntheticPCMGenerator) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = byte(i % 256)
	}
	return len(p), nil
}
