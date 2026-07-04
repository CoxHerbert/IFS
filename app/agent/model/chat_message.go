package model

import "time"

type ChatMessage struct {
	ID          int64     `db:"id" json:"id"`
	SessionID   int64     `db:"session_id" json:"sessionId"`
	Role        string    `db:"role" json:"role"`
	Content     string    `db:"content" json:"content"`
	BlockResult string    `db:"block_result" json:"-"`
	ModelName   string    `db:"model_name" json:"modelName"`
	CreatedAt   time.Time `db:"created_at" json:"-"`
}

type ChatMessageVO struct {
	ID          int64  `json:"id"`
	SessionID   int64  `json:"sessionId"`
	Role        string `json:"role"`
	Content     string `json:"content"`
	BlockResult any    `json:"blockResult,omitempty"`
	CreatedAt   string `json:"createdAt"`
}
