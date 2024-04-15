package models

import "time"

type Todo struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UserID    int64     `json:"user_id"`
}

type UserTodo struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type ResTodo struct {
	Id        int64  `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}
