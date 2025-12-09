package domain

import "time"

// Domain models are independent of infrastructure and external libraries.

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password,omitempty"` // hashed
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	Status      string `json:"status"`
}
