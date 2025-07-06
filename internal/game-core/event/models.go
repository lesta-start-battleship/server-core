package event

import "time"

type MatchResult struct {
	WinnerID   string      `json:"winner_id"`
	LoserID    string      `json:"loser_id"`
	MatchID    string      `json:"match_id"`
	MatchDate  time.Time   `json:"match_date"`
	MatchType  string      `json:"match_type"`
	Experience *Experience `json:"experience,omitempty"`
}

type Experience struct {
	WinnerGain int `json:"winner_gain,omitempty"`
	LoserGain  int `json:"loser_gain,omitempty"`
}

type Item struct {
	PlayerID string `json:"player_id"`
	ItemID   int    `json:"item_id"`
}
