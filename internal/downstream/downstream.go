package downstream

import (
	"bytes"
	"net/http"
	"time"

	"voxstream/internal/frame"
)

func NewClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 20,
			IdleConnTimeout:     90 * time.Second,
		},
		Timeout: 2 * time.Second,
	}
}

func Dispatch(client *http.Client, url string, f *frame.AudioFrame) bool {
	req, err := http.NewRequest("POST", url, bytes.NewReader(f.Data))
	if err != nil {
		return false
	}

	req.Header.Set("Content-Type", "audio/raw")
	req.Header.Set("X-Call-ID", f.CallID)

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
