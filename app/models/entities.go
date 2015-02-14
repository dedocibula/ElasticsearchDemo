package models

import (
	"time"
)

type Lesson struct {
	Number      int
	Description string
	HtmlClass   string
}

type Result struct {
	Ok       bool
	Messages []string
}

type ELKRecord struct {
	Nickname  string
	Timestamp time.Time
}

type Attempt struct {
	Name   string
	Answer int
}
