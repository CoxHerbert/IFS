package model

import "time"

type ChatMemory struct {
	ID          int64     `db:"id" json:"id"`
	UserID      int64     `db:"user_id" json:"userId"`
	MemoryKey   string    `db:"memory_key" json:"memoryKey"`
	MemoryValue string    `db:"memory_value" json:"memoryValue"`
	MemoryType  string    `db:"memory_type" json:"memoryType"`
	CreatedAt   time.Time `db:"created_at" json:"-"`
	UpdatedAt   time.Time `db:"updated_at" json:"-"`
}
