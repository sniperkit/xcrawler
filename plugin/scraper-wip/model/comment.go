package model

type Comment struct {
	Author  string
	URL     string
	Comment string
	Replies []*Comment
	depth   int
}
