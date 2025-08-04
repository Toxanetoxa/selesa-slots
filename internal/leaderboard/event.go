package leaderboard

import "time"

type Event struct {
	UserID    int64     `json:"user_id"`
	Position  int       `json:"position"`
	Score     int64     `json:"score"`
	UpdatedAt time.Time `json:"updated_at"`
}
