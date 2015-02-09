package models

type Lesson struct {
	Number      int
	Description string
	HtmlClass   string
}

type Result struct {
	Ok      bool
	Message string
}

type Player struct {
	Nickname string
	Position int
}

type Attempt struct {
	Name   string
	Answer int
}
