package models

import (
	"time"

	"gorm.io/gorm"
)

// ReadingSession represents a single reading session
type ReadingSession struct {
	gorm.Model
	UserID    uint      `json:"userId" gorm:"index"`
	BookID    uint      `json:"bookId" gorm:"index"`
	ChapterID uint      `json:"chapterId" gorm:"index"`
	SectionID uint      `json:"sectionId" gorm:"index"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Duration  int       `json:"duration"` // Duration in seconds
	Source    string    `json:"source"`   // Web, mobile, etc.
}

// ReadingProgress tracks a user's progress through a book
type ReadingProgress struct {
	gorm.Model
	UserID        uint      `json:"userId" gorm:"index"`
	BookID        uint      `json:"bookId" gorm:"index"`
	ChapterID     uint      `json:"chapterId"`
	SectionID     uint      `json:"sectionId"`
	LastReadAt    time.Time `json:"lastReadAt"`
	LastReadDay   time.Time `json:"lastReadDay" gorm:"index"` // Truncated to day for streak tracking
	TimeSpent     int       `json:"timeSpent"`                // Total time spent in seconds
	SessionCount  int       `json:"sessionCount"`             // Number of reading sessions
	StreakDays    int       `json:"streakDays"`               // Current streak in days
	LongestStreak int       `json:"longestStreak"`            // Longest streak achieved
	LastStreak    int       `json:"lastStreak"`               // Last streak before current one
}

// ReadingAnalytics represents aggregated reading analytics
type ReadingAnalytics struct {
	gorm.Model
	UserID            uint      `json:"userId" gorm:"index"`
	BookID            uint      `json:"bookId" gorm:"index"`
	TotalTimeSpent    int       `json:"totalTimeSpent"`    // Total time spent in seconds
	AvgSessionTime    int       `json:"avgSessionTime"`    // Average session time in seconds
	LongestSession    int       `json:"longestSession"`    // Longest session in seconds
	SessionCount      int       `json:"sessionCount"`      // Total number of sessions
	LastActivity      time.Time `json:"lastActivity"`      // Last reading activity
	MostActiveDay     string    `json:"mostActiveDay"`     // Most active day of the week
	PreferredTimeSlot string    `json:"preferredTimeSlot"` // Preferred time slot (Morning, Afternoon, Evening, Night)
}

// RecentlyViewedSection represents a section that has been recently viewed
type RecentlyViewedSection struct {
	gorm.Model
	UserID    uint      `json:"userId" gorm:"index"`
	BookID    uint      `json:"bookId" gorm:"index"`
	ChapterID uint      `json:"chapterId"`
	SectionID uint      `json:"sectionId"`
	ViewedAt  time.Time `json:"viewedAt"`
}

// ReadingStreak represents a user's daily reading streak
type ReadingStreak struct {
	gorm.Model
	UserID      uint      `json:"userId" gorm:"uniqueIndex:idx_user_date"`
	StreakDate  time.Time `json:"streakDate" gorm:"uniqueIndex:idx_user_date"`
	TimeSpent   int       `json:"timeSpent"`   // Time spent reading on this date (seconds)
	SectionRead int       `json:"sectionRead"` // Number of sections read on this date
}

// DailyReadingGoal represents a user's daily reading goal
type DailyReadingGoal struct {
	gorm.Model
	UserID         uint `json:"userId" gorm:"index"`
	MinutesPerDay  int  `json:"minutesPerDay"`  // Target minutes per day
	SectionsPerDay int  `json:"sectionsPerDay"` // Target sections per day
}

// ReadingRecommendation represents content recommendations based on reading behavior
type ReadingRecommendation struct {
	gorm.Model
	UserID      uint      `json:"userId" gorm:"index"`
	BookID      uint      `json:"bookId"`
	ChapterID   uint      `json:"chapterId"`
	SectionID   uint      `json:"sectionId"`
	Reason      string    `json:"reason"`      // Why this content is recommended
	Score       float64   `json:"score"`       // Recommendation strength score
	GeneratedAt time.Time `json:"generatedAt"` // When was this recommendation created
	Clicked     bool      `json:"clicked"`     // Whether user clicked on this recommendation
	ClickedAt   time.Time `json:"clickedAt"`   // When user clicked on this recommendation
}

// UserReadingPreference represents a user's reading preferences
type UserReadingPreference struct {
	gorm.Model
	UserID             uint    `json:"userId" gorm:"index"`
	PreferredTimeOfDay string  `json:"preferredTimeOfDay"` // Morning, Afternoon, Evening, Night
	SessionDuration    int     `json:"sessionDuration"`    // Typical session duration in minutes
	ContentDifficulty  string  `json:"contentDifficulty"`  // Easy, Medium, Hard
	ContentFormat      string  `json:"contentFormat"`      // Text, Audio, Interactive
	ReadingSpeed       float64 `json:"readingSpeed"`       // Words per minute
}