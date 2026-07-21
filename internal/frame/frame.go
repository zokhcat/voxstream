package frame

import "time"

type AudioFrame struct {
	CallID    string
	Data      []byte
	Length    int
	ArrivedAt time.Time
}
