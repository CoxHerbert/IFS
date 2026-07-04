package model

import "time"

type ChatSession struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"user_id" json:"userId"`
	Title     string    `db:"title" json:"title"`
	ModelName string    `db:"model_name" json:"modelName"`
	Summary   string    `db:"summary" json:"summary"`
	CreatedAt time.Time `db:"created_at" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
}

type ChatSessionVO struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	ModelName string `json:"modelName"`
	Summary   string `json:"summary"`
	UpdatedAt string `json:"updatedAt"`
}
