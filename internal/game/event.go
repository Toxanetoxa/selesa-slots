package game

import "time"

type Outcome string

const (
	Win  Outcome = "Win"
	Lose Outcome = "Lose"
)

type Event struct {
	GameID    string    `json:"game_id"`
	UserID    int64     `json:"user_id"`
	Outcome   Outcome   `json:"outcome"`
	Amount    int64     `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}
