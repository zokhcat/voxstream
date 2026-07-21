package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"voxstream/internal/config"
	"voxstream/internal/downstream"
	"voxstream/internal/mock"
	"voxstream/internal/proxy"
	"voxstream/internal/streamer"
)

func main() {
	stream := flag.Bool("stream", false, "Start a synthetic test stream")
	file := flag.String("file", "", "PCM file for test stream (implies -stream)")
	flag.Parse()

	log.Println("Starting VoxStream Low-Latency Audio Proxy...")

	mock.Serve()

	bufferPool := proxy.NewBufferPool()
	httpClient := downstream.NewClient()

	listener, err := net.Listen("tcp", config.TCPListenAddr)
	if err != nil {
		log.Fatalf("Failed to bind TCP listener on %s: %v", config.TCPListenAddr, err)
	}
	defer listener.Close()

	log.Printf("Ingestion Proxy listening on TCP %s", config.TCPListenAddr)
	log.Printf("Downstream Target: %s", config.MockBodhiURL)

	if *stream || *file != "" {
		time.Sleep(100 * time.Millisecond)
		go streamer.StreamAudio(*file)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
					log.Printf("Connection accept error: %v", err)
					continue
				}
			}

			go proxy.HandleCallStream(ctx, conn, bufferPool, httpClient)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down VoxStream Proxy gracefully...")
}
