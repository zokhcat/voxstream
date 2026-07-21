package config

var MockBodhiURL = "http://localhost:8081/bodhi/stream"

const (
	TCPListenAddr  = ":8080"
	ChunkSizeBytes = 320 // 20ms frame of 8kHz 16-bit mono PCM telephony audio
)
