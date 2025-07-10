package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

// InteractiveElementType defines the type of interactive elements
type InteractiveElementType string

const (
	// Quiz represents a quiz element with questions and answers
	QuizType InteractiveElementType = "quiz"
	
	// ReflectionType represents a reflection exercise
	ReflectionType InteractiveElementType = "reflection"
	
	// CallToActionType represents a call-to-action element
	CallToActionType InteractiveElementType = "call_to_action"
	
	// PollType represents a poll element
	PollType InteractiveElementType = "poll"
	
	// DiscussionPromptType represents a discussion prompt
	DiscussionPromptType InteractiveElementType = "discussion_prompt"
	
	// CodeChallengeType represents a code challenge
	CodeChallengeType InteractiveElementType = "code_challenge"
	
	// CaseStudyType represents a case study exercise
	CaseStudyType InteractiveElementType = "case_study"
	
	// VisualizationType represents an interactive data visualization
	VisualizationType InteractiveElementType = "visualization"
	
	// TimelineType represents an interactive timeline
	TimelineType InteractiveElementType = "timeline"
)

// InteractiveElement represents an interactive element within book content
type InteractiveElement struct {
	gorm.Model
	SectionID      uint                   `json:"sectionId" gorm:"index"`
	Position       int                    `json:"position"` // Position within the section
	Type           InteractiveElementType `json:"type"`
	Title          string                 `json:"title"`
	Description    string                 `json:"description"`
	Content        string                 `json:"content" gorm:"type:text"` // JSON content specific to element type
	CompletionType string                 `json:"completionType"`           // graded, self-check, no-check
	PointsValue    int                    `json:"pointsValue"`              // Points awarded for completion
	RequiredStatus bool                   `json:"requiredStatus"`           // Whether completion is required to progress
	Section        Section                `json:"-" gorm:"foreignKey:SectionID"`
}

// QuizContent represents the content structure for a quiz
type QuizContent struct {
	Questions []QuizQuestion `json:"questions"`
	TimeLimit int            `json:"timeLimit,omitempty"` // Time limit in seconds, 0 means no limit
	Randomize bool           `json:"randomize"`           // Whether to randomize question order
	PassScore int            `json:"passScore"`           // Percentage needed to pass
}

// QuizQuestion represents a single quiz question
type QuizQuestion struct {
	ID            uint              `json:"id"`
	QuestionText  string            `json:"questionText"`
	QuestionType  string            `json:"questionType"` // multiple-choice, true-false, fill-blank, short-answer
	Options       []QuizOption      `json:"options,omitempty"`
	CorrectAnswer interface{}       `json:"correctAnswer"` // String or array of strings
	Explanation   string            `json:"explanation,omitempty"`
	Hints         []string          `json:"hints,omitempty"`
	Media         *MediaAttachment  `json:"media,omitempty"`
	Tags          []string          `json:"tags,omitempty"`
	Difficulty    string            `json:"difficulty,omitempty"` // easy, medium, hard
	Metadata      map[string]string `json:"metadata,omitempty"`
}

// QuizOption represents an option for a multiple-choice question
type QuizOption struct {
	ID    string          `json:"id"`
	Text  string          `json:"text"`
	Media *MediaAttachment `json:"media,omitempty"`
}

// MediaAttachment represents media attached to a quiz question or option
type MediaAttachment struct {
	Type string `json:"type"` // image, audio, video
	URL  string `json:"url"`
	Alt  string `json:"alt,omitempty"` // Description for accessibility
}

// ReflectionContent represents the content structure for a reflection exercise
type ReflectionContent struct {
	Prompt            string   `json:"prompt"`
	GuidingQuestions  []string `json:"guidingQuestions,omitempty"`
	MinResponseLength int      `json:"minResponseLength,omitempty"` // Minimum characters required
	MaxResponseLength int      `json:"maxResponseLength,omitempty"` // Maximum characters allowed
	SharingOptions    []string `json:"sharingOptions,omitempty"`    // private, peers, public
}

// CallToActionContent represents the content structure for a call-to-action element
type CallToActionContent struct {
	ActionType        string            `json:"actionType"` // click, download, register, submit, etc.
	Text              string            `json:"text"`       // The call to action text
	URL               string            `json:"url,omitempty"`
	ButtonText        string            `json:"buttonText"`
	ButtonStyle       map[string]string `json:"buttonStyle,omitempty"`
	SecondaryText     string            `json:"secondaryText,omitempty"`
	CompletionMessage string            `json:"completionMessage,omitempty"`
	TrackingID        string            `json:"trackingId,omitempty"`
}

// DiscussionPromptContent represents the content structure for a discussion prompt
type DiscussionPromptContent struct {
	Topic             string   `json:"topic"`
	InitialPrompt     string   `json:"initialPrompt"`
	SupportingPoints  []string `json:"supportingPoints,omitempty"`
	DiscussionForumID uint     `json:"discussionForumId,omitempty"`
	Tags              []string `json:"tags,omitempty"`
	Guidelines        string   `json:"guidelines,omitempty"`
}

// PollContent represents the content structure for a poll
type PollContent struct {
	Question      string      `json:"question"`
	Options       []PollOption `json:"options"`
	AllowMultiple bool        `json:"allowMultiple"`
	ShowResults   string      `json:"showResults"` // always, after-vote, after-close, never
	ClosingDate   string      `json:"closingDate,omitempty"`
	AllowComments bool        `json:"allowComments"`
}

// PollOption represents an option in a poll
type PollOption struct {
	ID    string `json:"id"`
	Text  string `json:"text"`
	Image string `json:"image,omitempty"`
}

// GetQuizContent parses the interactive element's content as QuizContent
func (e *InteractiveElement) GetQuizContent() (*QuizContent, error) {
	if e.Type != QuizType {
		return nil, ErrInvalidElementType
	}
	
	var content QuizContent
	err := json.Unmarshal([]byte(e.Content), &content)
	if err != nil {
		return nil, err
	}
	
	return &content, nil
}

// GetReflectionContent parses the interactive element's content as ReflectionContent
func (e *InteractiveElement) GetReflectionContent() (*ReflectionContent, error) {
	if e.Type != ReflectionType {
		return nil, ErrInvalidElementType
	}
	
	var content ReflectionContent
	err := json.Unmarshal([]byte(e.Content), &content)
	if err != nil {
		return nil, err
	}
	
	return &content, nil
}

// GetCallToActionContent parses the interactive element's content as CallToActionContent
func (e *InteractiveElement) GetCallToActionContent() (*CallToActionContent, error) {
	if e.Type != CallToActionType {
		return nil, ErrInvalidElementType
	}
	
	var content CallToActionContent
	err := json.Unmarshal([]byte(e.Content), &content)
	if err != nil {
		return nil, err
	}
	
	return &content, nil
}

// GetDiscussionPromptContent parses the interactive element's content as DiscussionPromptContent
func (e *InteractiveElement) GetDiscussionPromptContent() (*DiscussionPromptContent, error) {
	if e.Type != DiscussionPromptType {
		return nil, ErrInvalidElementType
	}
	
	var content DiscussionPromptContent
	err := json.Unmarshal([]byte(e.Content), &content)
	if err != nil {
		return nil, err
	}
	
	return &content, nil
}

// GetPollContent parses the interactive element's content as PollContent
func (e *InteractiveElement) GetPollContent() (*PollContent, error) {
	if e.Type != PollType {
		return nil, ErrInvalidElementType
	}
	
	var content PollContent
	err := json.Unmarshal([]byte(e.Content), &content)
	if err != nil {
		return nil, err
	}
	
	return &content, nil
}

// InteractiveElementResponse represents a user's response to an interactive element
type InteractiveElementResponse struct {
	gorm.Model
	UserID               uint   `json:"userId" gorm:"index"`
	InteractiveElementID uint   `json:"interactiveElementId" gorm:"index"`
	Response             string `json:"response" gorm:"type:text"` // JSON response content
	Score                int    `json:"score,omitempty"`           // Score for graded elements
	TimeSpent            int    `json:"timeSpent"`                 // Time spent in seconds
	CompletionStatus     string `json:"completionStatus"`          // completed, in-progress, abandoned
	PointsAwarded        int    `json:"pointsAwarded"`             // Points awarded for this response
	InteractiveElement   InteractiveElement `json:"-" gorm:"foreignKey:InteractiveElementID"`
}

// UserInteractiveElementProgress tracks a user's progress with interactive elements
type UserInteractiveElementProgress struct {
	gorm.Model
	UserID                uint `json:"userId" gorm:"uniqueIndex:idx_user_book"`
	BookID                uint `json:"bookId" gorm:"uniqueIndex:idx_user_book"`
	TotalElements         int  `json:"totalElements"`
	CompletedElements     int  `json:"completedElements"`
	CompletionPercentage  int  `json:"completionPercentage"`
	TotalPointsAvailable  int  `json:"totalPointsAvailable"`
	TotalPointsEarned     int  `json:"totalPointsEarned"`
	AvgScorePercentage    int  `json:"avgScorePercentage"`
	RequiredCompleted     bool `json:"requiredCompleted"`
}