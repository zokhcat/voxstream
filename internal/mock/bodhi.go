package mock

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

const (
	ListenAddr = ":8081"
)

var (
	totalFramesReceived uint64
	totalBytesReceived  uint64
)

func Serve() {
	mux := http.NewServeMux()
	mux.HandleFunc("/bodhi/stream", handleAudioStream)
	mux.HandleFunc("/bodhi/metrics", handleMetrics)

	log.Printf("Bodhi Mock STT Engine listening on HTTP %s", ListenAddr)
	go func() {
		if err := http.ListenAndServe(ListenAddr, mux); err != nil {
			log.Printf("Mock Bodhi server error: %v", err)
		}
	}()
}

func handleAudioStream(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	callID := r.Header.Get("X-Call-ID")
	if callID == "" {
		callID = "unknown-session"
	}

	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "Invalid frame payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	frames := atomic.AddUint64(&totalFramesReceived, 1)
	bytes := atomic.AddUint64(&totalBytesReceived, uint64(len(body)))

	if frames%50 == 0 || frames == 1 {
		log.Printf("[Bodhi STT] CallID: %s | Frames Received: %d | Total Ingested: %d bytes",
			callID, frames, bytes)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := fmt.Sprintf(`{"status":"success","call_id":"%s","frame":%d,"transcript_chunk":""}`, callID, frames)
	w.Write([]byte(response))
}

func handleMetrics(w http.ResponseWriter, r *http.Request) {
	frames := atomic.LoadUint64(&totalFramesReceived)
	bytes := atomic.LoadUint64(&totalBytesReceived)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"online","total_frames":%d,"total_bytes":%d,"timestamp":"%s"}`,
		frames, bytes, time.Now().Format(time.RFC3339))
}
