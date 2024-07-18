package events

import "time"

type StartProcessingMessage struct {
	ID      uint64 `json:"id"`
	Content string `json:"content"`
}

type CompleteProcessingMessage struct {
	StartProcessingMessage
	ProcessedAt time.Time
}
