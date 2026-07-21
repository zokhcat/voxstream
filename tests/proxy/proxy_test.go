package proxy_test

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"voxstream/internal/config"
	"voxstream/internal/proxy"
)

func TestHandleCallStream_ReadsAndDispatchesFrames(t *testing.T) {
	var mu sync.Mutex
	var receivedFrames [][]byte

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		body := make([]byte, 320)
		r.Body.Read(body)
		receivedFrames = append(receivedFrames, body)
		w.WriteHeader(http.StatusOK)
	}))
	defer svr.Close()

	originalURL := config.MockBodhiURL
	config.MockBodhiURL = svr.URL
	defer func() { config.MockBodhiURL = originalURL }()

	serverConn, clientConn := net.Pipe()
	defer clientConn.Close()

	pool := proxy.NewBufferPool()
	client := &http.Client{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		proxy.HandleCallStream(ctx, serverConn, pool, client)
	}()

	frame := make([]byte, 320)
	for i := range frame {
		frame[i] = byte(i)
	}
	if _, err := clientConn.Write(frame); err != nil {
		t.Fatal(err)
	}
	clientConn.Close()
	wg.Wait()

	mu.Lock()
	if len(receivedFrames) != 1 {
		mu.Unlock()
		t.Fatalf("expected 1 frame, got %d", len(receivedFrames))
	}
	mu.Unlock()
}
