package downstream_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"voxstream/internal/downstream"
	"voxstream/internal/frame"
)

func TestDispatch_Success(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "audio/raw" {
			t.Errorf("expected audio/raw, got %s", r.Header.Get("Content-Type"))
		}
		if r.Header.Get("X-Call-ID") != "test-call" {
			t.Errorf("expected test-call, got %s", r.Header.Get("X-Call-ID"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer svr.Close()

	client := &http.Client{}
	f := &frame.AudioFrame{CallID: "test-call", Data: make([]byte, 320)}

	if !downstream.Dispatch(client, svr.URL, f) {
		t.Fatal("expected true")
	}
}

func TestDispatch_NonOKStatus(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	client := &http.Client{}
	f := &frame.AudioFrame{CallID: "test", Data: make([]byte, 320)}

	if downstream.Dispatch(client, svr.URL, f) {
		t.Fatal("expected false for 500")
	}
}

func TestDispatch_ConnectionError(t *testing.T) {
	client := &http.Client{}
	f := &frame.AudioFrame{CallID: "test", Data: make([]byte, 320)}

	if downstream.Dispatch(client, "http://localhost:19999", f) {
		t.Fatal("expected false for connection error")
	}
}
