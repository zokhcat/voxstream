package mock_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"voxstream/internal/mock"
)

func TestHandleAudioStream_Success(t *testing.T) {
	svr := httptest.NewServer(mock.NewMux())
	defer svr.Close()

	body := strings.NewReader(string(make([]byte, 320)))
	resp, err := http.Post(svr.URL+"/bodhi/stream", "audio/raw", body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result struct {
		Status string `json:"status"`
		CallID string `json:"call_id"`
		Frame  int    `json:"frame"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result.Status != "success" {
		t.Fatalf("expected success, got %s", result.Status)
	}
	if result.Frame != 1 {
		t.Fatalf("expected frame 1, got %d", result.Frame)
	}
}

func TestHandleAudioStream_EmptyBody(t *testing.T) {
	svr := httptest.NewServer(mock.NewMux())
	defer svr.Close()

	resp, err := http.Post(svr.URL+"/bodhi/stream", "audio/raw", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestHandleAudioStream_WrongMethod(t *testing.T) {
	svr := httptest.NewServer(mock.NewMux())
	defer svr.Close()

	resp, err := http.Get(svr.URL + "/bodhi/stream")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", resp.StatusCode)
	}
}

func TestHandleMetrics(t *testing.T) {
	svr := httptest.NewServer(mock.NewMux())
	defer svr.Close()

	resp, err := http.Get(svr.URL + "/bodhi/metrics")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result struct {
		Status      string `json:"status"`
		TotalFrames int    `json:"total_frames"`
		TotalBytes  int    `json:"total_bytes"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result.Status != "online" {
		t.Fatalf("expected online, got %s", result.Status)
	}
}
