package todo

import (
	"fmt"
	"time"
)

type Todo struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	UserID      string    `json:"user_id"`
}

func (t Todo) String() string {
	return fmt.Sprintf("[%s] %s", t.Created, t.Title)
}
