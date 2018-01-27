package domain

import (
	"time"
)

type Story struct {
	ID          int64
	Title       string
	URL         string
	Points      int64
	Comments    int64
	CommentsURL string
	Submitter   string
	Timestamp   time.Time
	Children    Comments
}

type Stories []Story
