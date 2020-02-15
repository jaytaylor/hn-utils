package domain

import (
	"encoding/json"
	"time"
)

// CommentNestingWidthIncrement is the CSS unit, in pixels, of one level of
// nesting depth in the discussion.
const CommentNestingWidthIncrement = 40

// Comment is a representation of a HackerNews comment.
type Comment struct {
	ID        int64
	Author    string
	Timestamp time.Time
	Content   string
	Width     int // Width indicates nesting depth, divide by 40 to get depth.
	N         int // N is the number of children in thread, as reported by the original data source.
	Children  Threads
}

type Comments []Comment

// Threads is a group of comments.
type Threads []*Comment

// Depth returns the nesting depth level of a comment.
func (c Comment) Depth() int {
	return c.Width / CommentNestingWidthIncrement
}

// ConversationLen recursively returns the number of descendent comments under
// this one.
//
// Note: The calculation includes this comment, so minimum return value is 1.
func (c Comment) ConversationLen() int {
	total := 1
	for _, child := range c.Children {
		total += child.ConversationLen()
	}
	return total
}

// Len returns the total number of comments across all conversations.
func (t Threads) Len() int {
	total := 0
	for _, c := range t {
		total += c.ConversationLen()
	}
	return total
}

// String returns a pretty-format JSON string representation of a thread.
func (t Threads) String() string {
	bs, _ := json.MarshalIndent(t, "", "    ")
	return string(bs)
}
