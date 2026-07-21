package proxy

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"voxstream/internal/config"
	"voxstream/internal/downstream"
	"voxstream/internal/frame"
)

func NewBufferPool() *sync.Pool {
	return &sync.Pool{
		New: func() any {
			b := make([]byte, config.ChunkSizeBytes)
			return &b
		},
	}
}

func HandleCallStream(ctx context.Context, conn net.Conn, pool *sync.Pool, client *http.Client) {
	defer conn.Close()

	callID := fmt.Sprintf("call-%d", time.Now().UnixNano()%100000)
	log.Printf("[%s] Stream connected from %s", callID, conn.RemoteAddr())

	var totalFrames int
	var totalOverheadUs int64

	for {
		select {
		case <-ctx.Done():
			log.Printf("[%s] Context canceled, closing stream loop", callID)
			return
		default:
		}

		bufPtr := pool.Get().(*[]byte)
		buf := *bufPtr

		_, err := io.ReadFull(conn, buf)
		if err != nil {
			pool.Put(bufPtr)
			if err != io.EOF && err != io.ErrUnexpectedEOF {
				log.Printf("[%s] Read error: %v", callID, err)
			}
			break
		}

		startTime := time.Now()

		f := frame.AudioFrame{
			CallID:    callID,
			Data:      buf,
			Length:    config.ChunkSizeBytes,
			ArrivedAt: startTime,
		}

		dispatched := downstream.Dispatch(client, &f)

		overhead := time.Since(startTime)
		totalFrames++
		totalOverheadUs += overhead.Microseconds()

		log.Printf(" [%s] Frame #%04d | Proxy Overhead: %4µs | Bodhi Dispatched: %t",
			callID, totalFrames, overhead.Microseconds(), dispatched)

		pool.Put(bufPtr)
	}

	avgLatency := float64(0)
	if totalFrames > 0 {
		avgLatency = float64(totalOverheadUs) / float64(totalFrames)
	}
	log.Printf(" [%s] Stream ended | Total Frames: %d | Avg Proxy Latency: %.2fµs",
		callID, totalFrames, avgLatency)
}
