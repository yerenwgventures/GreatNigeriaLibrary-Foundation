package models

// QuizQuestion represents a quiz question for a book section
type QuizQuestion struct {
	ID                string   `json:"id"`
	Question          string   `json:"question"`
	Options           []string `json:"options"`
	CorrectOptionIndex int      `json:"correctOptionIndex,omitempty"`
	Explanation       string   `json:"explanation,omitempty"`
}
