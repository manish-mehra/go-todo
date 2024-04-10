package models

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"` // Can hold any type of data
}
