package models

import "time"

type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UserID    int       `json:"user_id"`
}

type UserTodo struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type ResTodo struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}
