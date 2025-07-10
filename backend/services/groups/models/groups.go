package models

import (
	"time"

	"gorm.io/gorm"
)

// GroupType represents the type of group
type GroupType string

const (
	// LocalActionGroup represents a local action group
	LocalActionGroup GroupType = "local_action"
	// InterestGroup represents a special interest group
	InterestGroup GroupType = "interest"
	// ProjectGroup represents a project-focused group
	ProjectGroup GroupType = "project"
	// StudyGroup represents a study or learning group
	StudyGroup GroupType = "study"
)

// GroupVisibility represents the visibility of a group
type GroupVisibility string

const (
	// PublicGroup is visible to everyone and anyone can request to join
	PublicGroup GroupVisibility = "public"
	// PrivateGroup is visible in listings but requires approval to join
	PrivateGroup GroupVisibility = "private"
	// SecretGroup is not visible in listings and requires invitation to join
	SecretGroup GroupVisibility = "secret"
)

// MemberRole represents the role of a member in a group
type MemberRole string

const (
	// OwnerRole is the creator and primary administrator of the group
	OwnerRole MemberRole = "owner"
	// AdminRole has administrative privileges in the group
	AdminRole MemberRole = "admin"
	// ModeratorRole has moderation privileges in the group
	ModeratorRole MemberRole = "moderator"
	// MemberRole is a regular member of the group
	RegularMemberRole MemberRole = "member"
)

// MemberStatus represents the status of a member in a group
type MemberStatus string

const (
	// ActiveMember is a current active member
	ActiveMember MemberStatus = "active"
	// PendingMember has requested to join but hasn't been approved yet
	PendingMember MemberStatus = "pending"
	// InvitedMember has been invited but hasn't accepted yet
	InvitedMember MemberStatus = "invited"
	// BlockedMember has been blocked from the group
	BlockedMember MemberStatus = "blocked"
)

// EventStatus represents the status of a local event
type EventStatus string

const (
	// PlannedEvent is in the planning stage
	PlannedEvent EventStatus = "planned"
	// ConfirmedEvent is confirmed to happen
	ConfirmedEvent EventStatus = "confirmed"
	// InProgressEvent is currently happening
	InProgressEvent EventStatus = "in_progress"
	// CompletedEvent has been completed
	CompletedEvent EventStatus = "completed"
	// CancelledEvent has been cancelled
	CancelledEvent EventStatus = "cancelled"
)

// ResourceType represents the type of shared resource
type ResourceType string

const (
	// DocumentResource is a document (PDF, Word, etc.)
	DocumentResource ResourceType = "document"
	// ImageResource is an image
	ImageResource ResourceType = "image"
	// VideoResource is a video
	VideoResource ResourceType = "video"
	// AudioResource is an audio file
	AudioResource ResourceType = "audio"
	// LinkResource is a URL link
	LinkResource ResourceType = "link"
	// BookResource is a reference to a book in the system
	BookResource ResourceType = "book"
)

// Group represents a local group for coordination and collaboration
type Group struct {
	gorm.Model
	Name        string          `json:"name" gorm:"size:255;not null"`
	Description string          `json:"description" gorm:"type:text"`
	Type        GroupType       `json:"type" gorm:"size:50;not null"`
	Visibility  GroupVisibility `json:"visibility" gorm:"size:50;not null;default:'public'"`
	
	// Location information
	Location     string  `json:"location" gorm:"size:255"`
	Latitude     float64 `json:"latitude" gorm:"default:0"`
	Longitude    float64 `json:"longitude" gorm:"default:0"`
	IsVirtual    bool    `json:"isVirtual" gorm:"default:false"`
	
	// Group settings
	JoinApprovalRequired bool   `json:"joinApprovalRequired" gorm:"default:true"`
	MembershipCode       string `json:"membershipCode" gorm:"size:50"` // Optional code to join
	MaxMembers           int    `json:"maxMembers" gorm:"default:0"` // 0 means unlimited
	
	// Group metadata
	CreatedByID uint      `json:"createdById" gorm:"not null"`
	BannerImage string    `json:"bannerImage" gorm:"size:255"`
	AvatarImage string    `json:"avatarImage" gorm:"size:255"`
	Tags        string    `json:"tags" gorm:"size:255"` // Comma-separated tags
	
	// Relationships
	Members     []GroupMember     `json:"members,omitempty" gorm:"foreignKey:GroupID"`
	Events      []LocalEvent      `json:"events,omitempty" gorm:"foreignKey:GroupID"`
	Resources   []SharedResource  `json:"resources,omitempty" gorm:"foreignKey:GroupID"`
	Discussions []GroupDiscussion `json:"discussions,omitempty" gorm:"foreignKey:GroupID"`
	Actions     []GroupAction     `json:"actions,omitempty" gorm:"foreignKey:GroupID"`
}

// GroupMember represents a member of a group
type GroupMember struct {
	gorm.Model
	GroupID     uint         `json:"groupId" gorm:"not null"`
	UserID      uint         `json:"userId" gorm:"not null"`
	Role        MemberRole   `json:"role" gorm:"size:50;not null;default:'member'"`
	Status      MemberStatus `json:"status" gorm:"size:50;not null;default:'active'"`
	JoinedAt    *time.Time   `json:"joinedAt"`
	InvitedBy   *uint        `json:"invitedBy"`
	ApprovedBy  *uint        `json:"approvedBy"`
	LastActive  *time.Time   `json:"lastActive"`
	Notes       string       `json:"notes" gorm:"type:text"`
}

// LocalEvent represents an event organized by a local group
type LocalEvent struct {
	gorm.Model
	GroupID      uint        `json:"groupId" gorm:"not null"`
	Title        string      `json:"title" gorm:"size:255;not null"`
	Description  string      `json:"description" gorm:"type:text"`
	Location     string      `json:"location" gorm:"size:255"`
	Latitude     float64     `json:"latitude" gorm:"default:0"`
	Longitude    float64     `json:"longitude" gorm:"default:0"`
	IsVirtual    bool        `json:"isVirtual" gorm:"default:false"`
	VirtualLink  string      `json:"virtualLink" gorm:"size:255"`
	StartTime    time.Time   `json:"startTime" gorm:"not null"`
	EndTime      time.Time   `json:"endTime" gorm:"not null"`
	Status       EventStatus `json:"status" gorm:"size:50;not null;default:'planned'"`
	CreatedByID  uint        `json:"createdById" gorm:"not null"`
	MaxAttendees int         `json:"maxAttendees" gorm:"default:0"` // 0 means unlimited
	
	// Relationships
	Attendees []EventAttendee `json:"attendees,omitempty" gorm:"foreignKey:EventID"`
	Resources []SharedResource `json:"resources,omitempty" gorm:"foreignKey:EventID"`
}

// EventAttendee represents a user attending a local event
type EventAttendee struct {
	gorm.Model
	EventID    uint      `json:"eventId" gorm:"not null"`
	UserID     uint      `json:"userId" gorm:"not null"`
	Status     string    `json:"status" gorm:"size:50;not null;default:'going'"` // going, maybe, not_going
	RSVPTime   time.Time `json:"rsvpTime" gorm:"not null"`
	CheckedIn  bool      `json:"checkedIn" gorm:"default:false"`
	CheckInTime *time.Time `json:"checkInTime"`
	Notes      string    `json:"notes" gorm:"type:text"`
}

// SharedResource represents a resource shared within a group
type SharedResource struct {
	gorm.Model
	GroupID     uint         `json:"groupId"`
	EventID     *uint        `json:"eventId"`
	ActionID    *uint        `json:"actionId"`
	Title       string       `json:"title" gorm:"size:255;not null"`
	Description string       `json:"description" gorm:"type:text"`
	Type        ResourceType `json:"type" gorm:"size:50;not null"`
	URL         string       `json:"url" gorm:"size:255"`
	FilePath    string       `json:"filePath" gorm:"size:255"`
	FileSize    int64        `json:"fileSize" gorm:"default:0"`
	MimeType    string       `json:"mimeType" gorm:"size:100"`
	CreatedByID uint         `json:"createdById" gorm:"not null"`
	IsPublic    bool         `json:"isPublic" gorm:"default:false"`
	Tags        string       `json:"tags" gorm:"size:255"` // Comma-separated tags
}

// GroupDiscussion represents a discussion thread within a group
type GroupDiscussion struct {
	gorm.Model
	GroupID     uint   `json:"groupId" gorm:"not null"`
	Title       string `json:"title" gorm:"size:255;not null"`
	Content     string `json:"content" gorm:"type:text"`
	CreatedByID uint   `json:"createdById" gorm:"not null"`
	IsPinned    bool   `json:"isPinned" gorm:"default:false"`
	IsLocked    bool   `json:"isLocked" gorm:"default:false"`
	
	// Relationships
	Comments []GroupComment `json:"comments,omitempty" gorm:"foreignKey:DiscussionID"`
}

// GroupComment represents a comment in a group discussion
type GroupComment struct {
	gorm.Model
	DiscussionID uint   `json:"discussionId" gorm:"not null"`
	Content      string `json:"content" gorm:"type:text"`
	CreatedByID  uint   `json:"createdById" gorm:"not null"`
	ParentID     *uint  `json:"parentId"` // For replies to comments
}

// GroupAction represents an action or task within a group
type GroupAction struct {
	gorm.Model
	GroupID      uint      `json:"groupId" gorm:"not null"`
	Title        string    `json:"title" gorm:"size:255;not null"`
	Description  string    `json:"description" gorm:"type:text"`
	Status       string    `json:"status" gorm:"size:50;not null;default:'pending'"` // pending, in_progress, completed, cancelled
	Priority     string    `json:"priority" gorm:"size:50;not null;default:'medium'"` // low, medium, high, urgent
	DueDate      *time.Time `json:"dueDate"`
	CompletedAt  *time.Time `json:"completedAt"`
	CreatedByID  uint      `json:"createdById" gorm:"not null"`
	AssignedToID *uint     `json:"assignedToId"`
	
	// Relationships
	Resources []SharedResource `json:"resources,omitempty" gorm:"foreignKey:ActionID"`
}

// GroupInvitation represents an invitation to join a group
type GroupInvitation struct {
	gorm.Model
	GroupID     uint      `json:"groupId" gorm:"not null"`
	Email       string    `json:"email" gorm:"size:255"`
	UserID      *uint     `json:"userId"` // Optional, if inviting an existing user
	InvitedByID uint      `json:"invitedById" gorm:"not null"`
	Status      string    `json:"status" gorm:"size:50;not null;default:'pending'"` // pending, accepted, declined, expired
	ExpiresAt   time.Time `json:"expiresAt" gorm:"not null"`
	Code        string    `json:"code" gorm:"size:100;not null"` // Unique invitation code
	Message     string    `json:"message" gorm:"type:text"`
}

// GroupJoinRequest represents a request to join a group
type GroupJoinRequest struct {
	gorm.Model
	GroupID    uint      `json:"groupId" gorm:"not null"`
	UserID     uint      `json:"userId" gorm:"not null"`
	Status     string    `json:"status" gorm:"size:50;not null;default:'pending'"` // pending, approved, rejected
	Message    string    `json:"message" gorm:"type:text"`
	ApprovedBy *uint     `json:"approvedBy"`
	ApprovedAt *time.Time `json:"approvedAt"`
	RejectedBy *uint     `json:"rejectedBy"`
	RejectedAt *time.Time `json:"rejectedAt"`
	Code       string    `json:"code" gorm:"size:100"` // Optional membership code
}
