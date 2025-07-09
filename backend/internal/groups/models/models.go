package models

import (
	"time"

	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/pkg/models"
	"gorm.io/gorm"
)

// Group represents a local group
type Group struct {
	ID                  uint           `json:"id" gorm:"primaryKey"`
	Name                string         `json:"name" gorm:"size:255;not null"`
	Description         string         `json:"description" gorm:"type:text;not null"`
	Type                string         `json:"type" gorm:"size:50;not null;default:'local_action'"` // local_action, interest, project, study
	Visibility          string         `json:"visibility" gorm:"size:50;not null;default:'public'"` // public, private, secret
	Location            string         `json:"location" gorm:"size:255"`
	Latitude            float64        `json:"latitude" gorm:"type:decimal(10,8)"`
	Longitude           float64        `json:"longitude" gorm:"type:decimal(11,8)"`
	IsVirtual           bool           `json:"isVirtual" gorm:"default:false"`
	JoinApprovalRequired bool           `json:"joinApprovalRequired" gorm:"default:true"`
	MembershipCode      *string        `json:"membershipCode,omitempty" gorm:"size:50"`
	MaxMembers          int            `json:"maxMembers" gorm:"default:0"` // 0 means unlimited
	CreatedByID         uint           `json:"createdById" gorm:"not null"`
	BannerImage         *string        `json:"bannerImage,omitempty" gorm:"size:255"`
	AvatarImage         *string        `json:"avatarImage,omitempty" gorm:"size:255"`
	Tags                *string        `json:"tags,omitempty" gorm:"size:255"`
	CreatedAt           time.Time      `json:"createdAt"`
	UpdatedAt           time.Time      `json:"updatedAt"`
	DeletedAt           gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Members             []Member       `json:"members,omitempty" gorm:"foreignKey:GroupID"`
	Events              []LocalEvent   `json:"events,omitempty" gorm:"foreignKey:GroupID"`
	Resources           []SharedResource `json:"resources,omitempty" gorm:"foreignKey:GroupID"`
	Discussions         []Discussion   `json:"discussions,omitempty" gorm:"foreignKey:GroupID"`
	Actions             []Action       `json:"actions,omitempty" gorm:"foreignKey:GroupID"`
	Invitations         []Invitation   `json:"invitations,omitempty" gorm:"foreignKey:GroupID"`
	JoinRequests        []JoinRequest  `json:"joinRequests,omitempty" gorm:"foreignKey:GroupID"`
}

// MemberRole defines the role of a member in a group
type MemberRole string

const (
	OwnerRole     MemberRole = "owner"
	AdminRole     MemberRole = "admin"
	ModeratorRole MemberRole = "moderator"
	MemberRole    MemberRole = "member"
)

// MemberStatus defines the status of a member in a group
type MemberStatus string

const (
	ActiveMemberStatus  MemberStatus = "active"
	PendingMemberStatus MemberStatus = "pending"
	InvitedMemberStatus MemberStatus = "invited"
	BlockedMemberStatus MemberStatus = "blocked"
)

// Member represents a member of a group
type Member struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	GroupID    uint           `json:"groupId" gorm:"not null;index:idx_group_user,unique"`
	UserID     uint           `json:"userId" gorm:"not null;index:idx_group_user,unique"`
	Role       MemberRole     `json:"role" gorm:"size:20;not null;default:'member'"`
	Status     MemberStatus   `json:"status" gorm:"size:20;not null;default:'pending'"`
	JoinedAt   *time.Time     `json:"joinedAt,omitempty"`
	InvitedBy  *uint          `json:"invitedBy,omitempty"`
	ApprovedBy *uint          `json:"approvedBy,omitempty"`
	LastActive *time.Time     `json:"lastActive,omitempty"`
	Notes      *string        `json:"notes,omitempty" gorm:"type:text"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Group      *Group         `json:"group,omitempty" gorm:"foreignKey:GroupID"`
	User       *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// LocalEvent represents an event organized by a group
type LocalEvent struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	GroupID      uint           `json:"groupId" gorm:"not null;index"`
	Title        string         `json:"title" gorm:"size:255;not null"`
	Description  string         `json:"description" gorm:"type:text;not null"`
	Location     string         `json:"location" gorm:"size:255"`
	Latitude     float64        `json:"latitude" gorm:"type:decimal(10,8)"`
	Longitude    float64        `json:"longitude" gorm:"type:decimal(11,8)"`
	IsVirtual    bool           `json:"isVirtual" gorm:"default:false"`
	VirtualLink  *string        `json:"virtualLink,omitempty" gorm:"size:255"`
	StartTime    time.Time      `json:"startTime" gorm:"not null"`
	EndTime      time.Time      `json:"endTime" gorm:"not null"`
	Status       string         `json:"status" gorm:"size:20;not null;default:'planned'"` // planned, confirmed, in_progress, completed, cancelled
	CreatedByID  uint           `json:"createdById" gorm:"not null"`
	MaxAttendees int            `json:"maxAttendees" gorm:"default:0"` // 0 means unlimited
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Group        *Group         `json:"group,omitempty" gorm:"foreignKey:GroupID"`
	CreatedBy    *User          `json:"createdBy,omitempty" gorm:"foreignKey:CreatedByID"`
	Attendees    []Attendee     `json:"attendees,omitempty" gorm:"foreignKey:EventID"`
	Resources    []SharedResource `json:"resources,omitempty" gorm:"foreignKey:EventID"`
}

// Attendee represents a user attending an event
type Attendee struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	EventID     uint           `json:"eventId" gorm:"not null;index:idx_event_user,unique"`
	UserID      uint           `json:"userId" gorm:"not null;index:idx_event_user,unique"`
	Status      string         `json:"status" gorm:"size:20;not null;default:'going'"` // going, maybe, not_going
	RSVPTime    time.Time      `json:"rsvpTime" gorm:"not null"`
	CheckedIn   bool           `json:"checkedIn" gorm:"default:false"`
	CheckInTime *time.Time     `json:"checkInTime,omitempty"`
	Notes       *string        `json:"notes,omitempty" gorm:"type:text"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Event       *LocalEvent    `json:"event,omitempty" gorm:"foreignKey:EventID"`
	User        *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// SharedResource represents a resource shared within a group, event, or action
type SharedResource struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	GroupID     *uint          `json:"groupId,omitempty" gorm:"index"`
	EventID     *uint          `json:"eventId,omitempty" gorm:"index"`
	ActionID    *uint          `json:"actionId,omitempty" gorm:"index"`
	Title       string         `json:"title" gorm:"size:255;not null"`
	Description string         `json:"description" gorm:"type:text;not null"`
	Type        string         `json:"type" gorm:"size:20;not null"` // document, image, video, audio, link, book
	URL         *string        `json:"url,omitempty" gorm:"size:255"`
	FilePath    *string        `json:"filePath,omitempty" gorm:"size:255"`
	FileSize    *int64         `json:"fileSize,omitempty"`
	MimeType    *string        `json:"mimeType,omitempty" gorm:"size:100"`
	CreatedByID uint           `json:"createdById" gorm:"not null"`
	IsPublic    bool           `json:"isPublic" gorm:"default:false"`
	Tags        *string        `json:"tags,omitempty" gorm:"size:255"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Group       *Group         `json:"group,omitempty" gorm:"foreignKey:GroupID"`
	Event       *LocalEvent    `json:"event,omitempty" gorm:"foreignKey:EventID"`
	Action      *Action        `json:"action,omitempty" gorm:"foreignKey:ActionID"`
	CreatedBy   *User          `json:"createdBy,omitempty" gorm:"foreignKey:CreatedByID"`
}

// Discussion represents a discussion thread in a group
type Discussion struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	GroupID     uint           `json:"groupId" gorm:"not null;index"`
	Title       string         `json:"title" gorm:"size:255;not null"`
	Content     string         `json:"content" gorm:"type:text;not null"`
	CreatedByID uint           `json:"createdById" gorm:"not null"`
	IsPinned    bool           `json:"isPinned" gorm:"default:false"`
	IsLocked    bool           `json:"isLocked" gorm:"default:false"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Group       *Group         `json:"group,omitempty" gorm:"foreignKey:GroupID"`
	CreatedBy   *User          `json:"createdBy,omitempty" gorm:"foreignKey:CreatedByID"`
	Comments    []Comment      `json:"comments,omitempty" gorm:"foreignKey:DiscussionID"`
}

// Comment represents a comment on a discussion
type Comment struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	DiscussionID uint           `json:"discussionId" gorm:"not null;index"`
	Content      string         `json:"content" gorm:"type:text;not null"`
	CreatedByID  uint           `json:"createdById" gorm:"not null"`
	ParentID     *uint          `json:"parentId,omitempty" gorm:"index"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Discussion   *Discussion    `json:"discussion,omitempty" gorm:"foreignKey:DiscussionID"`
	CreatedBy    *User          `json:"createdBy,omitempty" gorm:"foreignKey:CreatedByID"`
	Parent       *Comment       `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Replies      []Comment      `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
}

// Action represents a task or action item for a group
type Action struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	GroupID      uint           `json:"groupId" gorm:"not null;index"`
	Title        string         `json:"title" gorm:"size:255;not null"`
	Description  string         `json:"description" gorm:"type:text;not null"`
	Status       string         `json:"status" gorm:"size:20;not null;default:'pending'"` // pending, in_progress, completed, cancelled
	Priority     string         `json:"priority" gorm:"size:20;not null;default:'medium'"` // low, medium, high, urgent
	DueDate      *time.Time     `json:"dueDate,omitempty"`
	CompletedAt  *time.Time     `json:"completedAt,omitempty"`
	CreatedByID  uint           `json:"createdById" gorm:"not null"`
	AssignedToID *uint          `json:"assignedToId,omitempty" gorm:"index"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Group        *Group         `json:"group,omitempty" gorm:"foreignKey:GroupID"`
	CreatedBy    *User          `json:"createdBy,omitempty" gorm:"foreignKey:CreatedByID"`
	AssignedTo   *User          `json:"assignedTo,omitempty" gorm:"foreignKey:AssignedToID"`
	Resources    []SharedResource `json:"resources,omitempty" gorm:"foreignKey:ActionID"`
}

// Invitation represents an invitation to join a group
type Invitation struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	GroupID      uint           `json:"groupId" gorm:"not null;index"`
	Email        string         `json:"email" gorm:"size:255;not null;index"`
	UserID       *uint          `json:"userId,omitempty" gorm:"index"`
	InvitedByID  uint           `json:"invitedById" gorm:"not null"`
	Status       string         `json:"status" gorm:"size:20;not null;default:'pending'"` // pending, accepted, declined, expired, cancelled
	ExpiresAt    time.Time      `json:"expiresAt" gorm:"not null"`
	Code         string         `json:"code" gorm:"size:50;not null;index"`
	Message      *string        `json:"message,omitempty" gorm:"type:text"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Group        *Group         `json:"group,omitempty" gorm:"foreignKey:GroupID"`
	User         *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	InvitedBy    *User          `json:"invitedBy,omitempty" gorm:"foreignKey:InvitedByID"`
}

// JoinRequest represents a request to join a group
type JoinRequest struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	GroupID      uint           `json:"groupId" gorm:"not null;index"`
	UserID       uint           `json:"userId" gorm:"not null;index"`
	Status       string         `json:"status" gorm:"size:20;not null;default:'pending'"` // pending, approved, rejected
	Message      *string        `json:"message,omitempty" gorm:"type:text"`
	ApprovedBy   *uint          `json:"approvedBy,omitempty"`
	ApprovedAt   *time.Time     `json:"approvedAt,omitempty"`
	RejectedBy   *uint          `json:"rejectedBy,omitempty"`
	RejectedAt   *time.Time     `json:"rejectedAt,omitempty"`
	Code         *string        `json:"code,omitempty" gorm:"size:50"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Group        *Group         `json:"group,omitempty" gorm:"foreignKey:GroupID"`
	User         *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// User model is now imported from shared models package
// import "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/pkg/models"
