package repository

import (
	"fmt"
	"time"

	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/groups/models"
	"gorm.io/gorm"
)

// GroupRepository defines the interface for group data access
type GroupRepository interface {
	// Group operations
	CreateGroup(group *models.Group) error
	GetGroupByID(id uint) (*models.Group, error)
	GetGroupsByUserID(userID uint, includeInvited bool) ([]models.Group, error)
	GetGroupsByLocation(latitude, longitude float64, radiusKm float64) ([]models.Group, error)
	UpdateGroup(group *models.Group) error
	DeleteGroup(id uint) error
	SearchGroups(query string, groupType string, tags []string, page, pageSize int) ([]models.Group, int64, error)

	// Member operations
	AddMember(member *models.GroupMember) error
	GetMemberByID(groupID, userID uint) (*models.GroupMember, error)
	GetMembersByGroupID(groupID uint, status string, page, pageSize int) ([]models.GroupMember, int64, error)
	UpdateMember(member *models.GroupMember) error
	RemoveMember(groupID, userID uint) error

	// Event operations
	CreateEvent(event *models.LocalEvent) error
	GetEventByID(id uint) (*models.LocalEvent, error)
	GetEventsByGroupID(groupID uint, status string, page, pageSize int) ([]models.LocalEvent, int64, error)
	GetUpcomingEvents(userID uint, page, pageSize int) ([]models.LocalEvent, int64, error)
	UpdateEvent(event *models.LocalEvent) error
	DeleteEvent(id uint) error

	// Attendee operations
	AddAttendee(attendee *models.EventAttendee) error
	GetAttendeeByID(eventID, userID uint) (*models.EventAttendee, error)
	GetAttendeesByEventID(eventID uint, status string, page, pageSize int) ([]models.EventAttendee, int64, error)
	UpdateAttendee(attendee *models.EventAttendee) error
	RemoveAttendee(eventID, userID uint) error

	// Resource operations
	CreateResource(resource *models.SharedResource) error
	GetResourceByID(id uint) (*models.SharedResource, error)
	GetResourcesByGroupID(groupID uint, resourceType string, page, pageSize int) ([]models.SharedResource, int64, error)
	GetResourcesByEventID(eventID uint, page, pageSize int) ([]models.SharedResource, int64, error)
	GetResourcesByActionID(actionID uint, page, pageSize int) ([]models.SharedResource, int64, error)
	UpdateResource(resource *models.SharedResource) error
	DeleteResource(id uint) error

	// Discussion operations
	CreateDiscussion(discussion *models.GroupDiscussion) error
	GetDiscussionByID(id uint) (*models.GroupDiscussion, error)
	GetDiscussionsByGroupID(groupID uint, page, pageSize int) ([]models.GroupDiscussion, int64, error)
	UpdateDiscussion(discussion *models.GroupDiscussion) error
	DeleteDiscussion(id uint) error

	// Comment operations
	CreateComment(comment *models.GroupComment) error
	GetCommentByID(id uint) (*models.GroupComment, error)
	GetCommentsByDiscussionID(discussionID uint, page, pageSize int) ([]models.GroupComment, int64, error)
	UpdateComment(comment *models.GroupComment) error
	DeleteComment(id uint) error

	// Action operations
	CreateAction(action *models.GroupAction) error
	GetActionByID(id uint) (*models.GroupAction, error)
	GetActionsByGroupID(groupID uint, status string, page, pageSize int) ([]models.GroupAction, int64, error)
	GetActionsByAssigneeID(userID uint, status string, page, pageSize int) ([]models.GroupAction, int64, error)
	UpdateAction(action *models.GroupAction) error
	DeleteAction(id uint) error

	// Invitation operations
	CreateInvitation(invitation *models.GroupInvitation) error
	GetInvitationByID(id uint) (*models.GroupInvitation, error)
	GetInvitationByCode(code string) (*models.GroupInvitation, error)
	GetInvitationsByGroupID(groupID uint, status string, page, pageSize int) ([]models.GroupInvitation, int64, error)
	GetInvitationsByUserID(userID uint, status string, page, pageSize int) ([]models.GroupInvitation, int64, error)
	UpdateInvitation(invitation *models.GroupInvitation) error
	DeleteInvitation(id uint) error

	// Join request operations
	CreateJoinRequest(request *models.GroupJoinRequest) error
	GetJoinRequestByID(id uint) (*models.GroupJoinRequest, error)
	GetJoinRequestsByGroupID(groupID uint, status string, page, pageSize int) ([]models.GroupJoinRequest, int64, error)
	GetJoinRequestsByUserID(userID uint, status string, page, pageSize int) ([]models.GroupJoinRequest, int64, error)
	UpdateJoinRequest(request *models.GroupJoinRequest) error
	DeleteJoinRequest(id uint) error
}

// GroupRepositoryImpl implements the GroupRepository interface
type GroupRepositoryImpl struct {
	db *gorm.DB
}

// NewGroupRepository creates a new group repository
func NewGroupRepository(db *gorm.DB) GroupRepository {
	return &GroupRepositoryImpl{
		db: db,
	}
}

// CreateGroup creates a new group
func (r *GroupRepositoryImpl) CreateGroup(group *models.Group) error {
	return r.db.Create(group).Error
}

// GetGroupByID retrieves a group by its ID
func (r *GroupRepositoryImpl) GetGroupByID(id uint) (*models.Group, error) {
	var group models.Group
	err := r.db.Preload("Members", func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", models.ActiveMember)
	}).First(&group, id).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// GetGroupsByUserID retrieves groups that a user is a member of
func (r *GroupRepositoryImpl) GetGroupsByUserID(userID uint, includeInvited bool) ([]models.Group, error) {
	var groups []models.Group

	query := r.db.Joins("JOIN group_members ON groups.id = group_members.group_id").
		Where("group_members.user_id = ? AND group_members.status = ?", userID, models.ActiveMember)

	if includeInvited {
		query = query.Or("group_members.user_id = ? AND group_members.status = ?", userID, models.InvitedMember)
	}

	err := query.Preload("Members", func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", models.ActiveMember)
	}).Find(&groups).Error

	if err != nil {
		return nil, err
	}

	return groups, nil
}

// GetGroupsByLocation retrieves groups near a specific location
func (r *GroupRepositoryImpl) GetGroupsByLocation(latitude, longitude float64, radiusKm float64) ([]models.Group, error) {
	var groups []models.Group

	// Haversine formula to calculate distance
	// This is a simplified version and might not be accurate for very large distances
	// For production, consider using PostGIS or a similar spatial database extension
	distanceQuery := fmt.Sprintf(
		"(6371 * acos(cos(radians(%f)) * cos(radians(latitude)) * cos(radians(longitude) - radians(%f)) + sin(radians(%f)) * sin(radians(latitude))))",
		latitude, longitude, latitude,
	)

	err := r.db.Where(fmt.Sprintf("%s <= ?", distanceQuery), radiusKm).
		Where("visibility = ? OR visibility = ?", models.PublicGroup, models.PrivateGroup).
		Preload("Members", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", models.ActiveMember)
		}).Find(&groups).Error

	if err != nil {
		return nil, err
	}

	return groups, nil
}

// UpdateGroup updates a group
func (r *GroupRepositoryImpl) UpdateGroup(group *models.Group) error {
	return r.db.Save(group).Error
}

// DeleteGroup deletes a group
func (r *GroupRepositoryImpl) DeleteGroup(id uint) error {
	return r.db.Delete(&models.Group{}, id).Error
}

// SearchGroups searches for groups based on criteria
func (r *GroupRepositoryImpl) SearchGroups(query string, groupType string, tags []string, page, pageSize int) ([]models.Group, int64, error) {
	var groups []models.Group
	var total int64

	db := r.db.Model(&models.Group{}).
		Where("visibility = ? OR visibility = ?", models.PublicGroup, models.PrivateGroup)

	if query != "" {
		db = db.Where("name LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%")
	}

	if groupType != "" {
		db = db.Where("type = ?", groupType)
	}

	if len(tags) > 0 {
		for _, tag := range tags {
			db = db.Where("tags LIKE ?", "%"+tag+"%")
		}
	}

	// Count total results
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err = db.Offset(offset).Limit(pageSize).
		Preload("Members", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", models.ActiveMember)
		}).Find(&groups).Error

	if err != nil {
		return nil, 0, err
	}

	return groups, total, nil
}

// AddMember adds a member to a group
func (r *GroupRepositoryImpl) AddMember(member *models.GroupMember) error {
	// Set joined time if status is active
	if member.Status == models.ActiveMember && member.JoinedAt == nil {
		now := time.Now()
		member.JoinedAt = &now
	}

	return r.db.Create(member).Error
}

// GetMemberByID retrieves a group member by group ID and user ID
func (r *GroupRepositoryImpl) GetMemberByID(groupID, userID uint) (*models.GroupMember, error) {
	var member models.GroupMember
	err := r.db.Where("group_id = ? AND user_id = ?", groupID, userID).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// GetMembersByGroupID retrieves members of a group
func (r *GroupRepositoryImpl) GetMembersByGroupID(groupID uint, status string, page, pageSize int) ([]models.GroupMember, int64, error) {
	var members []models.GroupMember
	var total int64

	query := r.db.Where("group_id = ?", groupID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count total results
	err := query.Model(&models.GroupMember{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Find(&members).Error
	if err != nil {
		return nil, 0, err
	}

	return members, total, nil
}

// UpdateMember updates a group member
func (r *GroupRepositoryImpl) UpdateMember(member *models.GroupMember) error {
	return r.db.Save(member).Error
}

// RemoveMember removes a member from a group
func (r *GroupRepositoryImpl) RemoveMember(groupID, userID uint) error {
	return r.db.Where("group_id = ? AND user_id = ?", groupID, userID).Delete(&models.GroupMember{}).Error
}

// CreateEvent creates a new local event
func (r *GroupRepositoryImpl) CreateEvent(event *models.LocalEvent) error {
	return r.db.Create(event).Error
}

// GetEventByID retrieves an event by its ID
func (r *GroupRepositoryImpl) GetEventByID(id uint) (*models.LocalEvent, error) {
	var event models.LocalEvent
	err := r.db.Preload("Attendees").First(&event, id).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// GetEventsByGroupID retrieves events for a group
func (r *GroupRepositoryImpl) GetEventsByGroupID(groupID uint, status string, page, pageSize int) ([]models.LocalEvent, int64, error) {
	var events []models.LocalEvent
	var total int64

	query := r.db.Where("group_id = ?", groupID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count total results
	err := query.Model(&models.LocalEvent{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).
		Preload("Attendees").
		Order("start_time ASC").
		Find(&events).Error

	if err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

// GetUpcomingEvents retrieves upcoming events for a user
func (r *GroupRepositoryImpl) GetUpcomingEvents(userID uint, page, pageSize int) ([]models.LocalEvent, int64, error) {
	var events []models.LocalEvent
	var total int64

	// Get events from groups the user is a member of
	query := r.db.Joins("JOIN group_members ON local_events.group_id = group_members.group_id").
		Where("group_members.user_id = ? AND group_members.status = ?", userID, models.ActiveMember).
		Where("local_events.start_time > ? AND local_events.status IN (?, ?)",
			time.Now(), models.PlannedEvent, models.ConfirmedEvent)

	// Count total results
	err := query.Model(&models.LocalEvent{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).
		Preload("Attendees").
		Order("start_time ASC").
		Find(&events).Error

	if err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

// UpdateEvent updates an event
func (r *GroupRepositoryImpl) UpdateEvent(event *models.LocalEvent) error {
	return r.db.Save(event).Error
}

// DeleteEvent deletes an event
func (r *GroupRepositoryImpl) DeleteEvent(id uint) error {
	return r.db.Delete(&models.LocalEvent{}, id).Error
}

// AddAttendee adds an attendee to an event
func (r *GroupRepositoryImpl) AddAttendee(attendee *models.EventAttendee) error {
	return r.db.Create(attendee).Error
}

// GetAttendeeByID retrieves an attendee by event ID and user ID
func (r *GroupRepositoryImpl) GetAttendeeByID(eventID, userID uint) (*models.EventAttendee, error) {
	var attendee models.EventAttendee
	err := r.db.Where("event_id = ? AND user_id = ?", eventID, userID).First(&attendee).Error
	if err != nil {
		return nil, err
	}
	return &attendee, nil
}

// GetAttendeesByEventID retrieves attendees for an event
func (r *GroupRepositoryImpl) GetAttendeesByEventID(eventID uint, status string, page, pageSize int) ([]models.EventAttendee, int64, error) {
	var attendees []models.EventAttendee
	var total int64

	query := r.db.Where("event_id = ?", eventID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count total results
	err := query.Model(&models.EventAttendee{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Find(&attendees).Error
	if err != nil {
		return nil, 0, err
	}

	return attendees, total, nil
}

// UpdateAttendee updates an event attendee
func (r *GroupRepositoryImpl) UpdateAttendee(attendee *models.EventAttendee) error {
	return r.db.Save(attendee).Error
}

// RemoveAttendee removes an attendee from an event
func (r *GroupRepositoryImpl) RemoveAttendee(eventID, userID uint) error {
	return r.db.Where("event_id = ? AND user_id = ?", eventID, userID).Delete(&models.EventAttendee{}).Error
}

// CreateResource creates a new shared resource
func (r *GroupRepositoryImpl) CreateResource(resource *models.SharedResource) error {
	return r.db.Create(resource).Error
}

// GetResourceByID retrieves a resource by its ID
func (r *GroupRepositoryImpl) GetResourceByID(id uint) (*models.SharedResource, error) {
	var resource models.SharedResource
	err := r.db.First(&resource, id).Error
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetResourcesByGroupID retrieves resources for a group
func (r *GroupRepositoryImpl) GetResourcesByGroupID(groupID uint, resourceType string, page, pageSize int) ([]models.SharedResource, int64, error) {
	var resources []models.SharedResource
	var total int64

	query := r.db.Where("group_id = ?", groupID)

	if resourceType != "" {
		query = query.Where("type = ?", resourceType)
	}

	// Count total results
	err := query.Model(&models.SharedResource{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&resources).Error

	if err != nil {
		return nil, 0, err
	}

	return resources, total, nil
}

// GetResourcesByEventID retrieves resources for an event
func (r *GroupRepositoryImpl) GetResourcesByEventID(eventID uint, page, pageSize int) ([]models.SharedResource, int64, error) {
	var resources []models.SharedResource
	var total int64

	query := r.db.Where("event_id = ?", eventID)

	// Count total results
	err := query.Model(&models.SharedResource{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&resources).Error

	if err != nil {
		return nil, 0, err
	}

	return resources, total, nil
}

// GetResourcesByActionID retrieves resources for an action
func (r *GroupRepositoryImpl) GetResourcesByActionID(actionID uint, page, pageSize int) ([]models.SharedResource, int64, error) {
	var resources []models.SharedResource
	var total int64

	query := r.db.Where("action_id = ?", actionID)

	// Count total results
	err := query.Model(&models.SharedResource{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&resources).Error

	if err != nil {
		return nil, 0, err
	}

	return resources, total, nil
}

// UpdateResource updates a resource
func (r *GroupRepositoryImpl) UpdateResource(resource *models.SharedResource) error {
	return r.db.Save(resource).Error
}

// DeleteResource deletes a resource
func (r *GroupRepositoryImpl) DeleteResource(id uint) error {
	return r.db.Delete(&models.SharedResource{}, id).Error
}

// CreateDiscussion creates a new group discussion
func (r *GroupRepositoryImpl) CreateDiscussion(discussion *models.GroupDiscussion) error {
	return r.db.Create(discussion).Error
}

// GetDiscussionByID retrieves a discussion by its ID
func (r *GroupRepositoryImpl) GetDiscussionByID(id uint) (*models.GroupDiscussion, error) {
	var discussion models.GroupDiscussion
	err := r.db.Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC")
	}).First(&discussion, id).Error
	if err != nil {
		return nil, err
	}
	return &discussion, nil
}

// GetDiscussionsByGroupID retrieves discussions for a group
func (r *GroupRepositoryImpl) GetDiscussionsByGroupID(groupID uint, page, pageSize int) ([]models.GroupDiscussion, int64, error) {
	var discussions []models.GroupDiscussion
	var total int64

	query := r.db.Where("group_id = ?", groupID)

	// Count total results
	err := query.Model(&models.GroupDiscussion{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).
		Order("is_pinned DESC, created_at DESC").
		Find(&discussions).Error

	if err != nil {
		return nil, 0, err
	}

	return discussions, total, nil
}

// UpdateDiscussion updates a discussion
func (r *GroupRepositoryImpl) UpdateDiscussion(discussion *models.GroupDiscussion) error {
	return r.db.Save(discussion).Error
}

// DeleteDiscussion deletes a discussion
func (r *GroupRepositoryImpl) DeleteDiscussion(id uint) error {
	return r.db.Delete(&models.GroupDiscussion{}, id).Error
}

// CreateComment creates a new comment
func (r *GroupRepositoryImpl) CreateComment(comment *models.GroupComment) error {
	return r.db.Create(comment).Error
}

// GetCommentByID retrieves a comment by its ID
func (r *GroupRepositoryImpl) GetCommentByID(id uint) (*models.GroupComment, error) {
	var comment models.GroupComment
	err := r.db.First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetCommentsByDiscussionID retrieves comments for a discussion
func (r *GroupRepositoryImpl) GetCommentsByDiscussionID(discussionID uint, page, pageSize int) ([]models.GroupComment, int64, error) {
	var comments []models.GroupComment
	var total int64

	query := r.db.Where("discussion_id = ? AND parent_id IS NULL", discussionID)

	// Count total results
	err := query.Model(&models.GroupComment{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).
		Order("created_at ASC").
		Find(&comments).Error

	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// UpdateComment updates a comment
func (r *GroupRepositoryImpl) UpdateComment(comment *models.GroupComment) error {
	return r.db.Save(comment).Error
}

// DeleteComment deletes a comment
func (r *GroupRepositoryImpl) DeleteComment(id uint) error {
	return r.db.Delete(&models.GroupComment{}, id).Error
}

// CreateAction creates a new group action
func (r *GroupRepositoryImpl) CreateAction(action *models.GroupAction) error {
	return r.db.Create(action).Error
}

// GetActionByID retrieves an action by its ID
func (r *GroupRepositoryImpl) GetActionByID(id uint) (*models.GroupAction, error) {
	var action models.GroupAction
	err := r.db.Preload("Resources").First(&action, id).Error
	if err != nil {
		return nil, err
	}
	return &action, nil
}

// GetActionsByGroupID retrieves actions for a group
func (r *GroupRepositoryImpl) GetActionsByGroupID(groupID uint, status string, page, pageSize int) ([]models.GroupAction, int64, error) {
	var actions []models.GroupAction
	var total int64

	query := r.db.Where("group_id = ?", groupID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count total results
	err := query.Model(&models.GroupAction{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).
		Order("due_date ASC, priority DESC").
		Find(&actions).Error

	if err != nil {
		return nil, 0, err
	}

	return actions, total, nil
}

// GetActionsByAssigneeID retrieves actions assigned to a user
func (r *GroupRepositoryImpl) GetActionsByAssigneeID(userID uint, status string, page, pageSize int) ([]models.GroupAction, int64, error) {
	var actions []models.GroupAction
	var total int64

	query := r.db.Where("assigned_to_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count total results
	err := query.Model(&models.GroupAction{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).
		Order("due_date ASC, priority DESC").
		Find(&actions).Error

	if err != nil {
		return nil, 0, err
	}

	return actions, total, nil
}

// UpdateAction updates an action
func (r *GroupRepositoryImpl) UpdateAction(action *models.GroupAction) error {
	return r.db.Save(action).Error
}

// DeleteAction deletes an action
func (r *GroupRepositoryImpl) DeleteAction(id uint) error {
	return r.db.Delete(&models.GroupAction{}, id).Error
}

// CreateInvitation creates a new group invitation
func (r *GroupRepositoryImpl) CreateInvitation(invitation *models.GroupInvitation) error {
	return r.db.Create(invitation).Error
}

// GetInvitationByID retrieves an invitation by its ID
func (r *GroupRepositoryImpl) GetInvitationByID(id uint) (*models.GroupInvitation, error) {
	var invitation models.GroupInvitation
	err := r.db.First(&invitation, id).Error
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

// GetInvitationByCode retrieves an invitation by its code
func (r *GroupRepositoryImpl) GetInvitationByCode(code string) (*models.GroupInvitation, error) {
	var invitation models.GroupInvitation
	err := r.db.Where("code = ?", code).First(&invitation).Error
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

// GetInvitationsByGroupID retrieves invitations for a group
func (r *GroupRepositoryImpl) GetInvitationsByGroupID(groupID uint, status string, page, pageSize int) ([]models.GroupInvitation, int64, error) {
	var invitations []models.GroupInvitation
	var total int64

	query := r.db.Where("group_id = ?", groupID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count total results
	err := query.Model(&models.GroupInvitation{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&invitations).Error

	if err != nil {
		return nil, 0, err
	}

	return invitations, total, nil
}

// GetInvitationsByUserID retrieves invitations for a user
func (r *GroupRepositoryImpl) GetInvitationsByUserID(userID uint, status string, page, pageSize int) ([]models.GroupInvitation, int64, error) {
	var invitations []models.GroupInvitation
	var total int64

	query := r.db.Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count total results
	err := query.Model(&models.GroupInvitation{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&invitations).Error

	if err != nil {
		return nil, 0, err
	}

	return invitations, total, nil
}

// UpdateInvitation updates an invitation
func (r *GroupRepositoryImpl) UpdateInvitation(invitation *models.GroupInvitation) error {
	return r.db.Save(invitation).Error
}

// DeleteInvitation deletes an invitation
func (r *GroupRepositoryImpl) DeleteInvitation(id uint) error {
	return r.db.Delete(&models.GroupInvitation{}, id).Error
}

// CreateJoinRequest creates a new join request
func (r *GroupRepositoryImpl) CreateJoinRequest(request *models.GroupJoinRequest) error {
	return r.db.Create(request).Error
}

// GetJoinRequestByID retrieves a join request by its ID
func (r *GroupRepositoryImpl) GetJoinRequestByID(id uint) (*models.GroupJoinRequest, error) {
	var request models.GroupJoinRequest
	err := r.db.First(&request, id).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

// GetJoinRequestsByGroupID retrieves join requests for a group
func (r *GroupRepositoryImpl) GetJoinRequestsByGroupID(groupID uint, status string, page, pageSize int) ([]models.GroupJoinRequest, int64, error) {
	var requests []models.GroupJoinRequest
	var total int64

	query := r.db.Where("group_id = ?", groupID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count total results
	err := query.Model(&models.GroupJoinRequest{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&requests).Error

	if err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

// GetJoinRequestsByUserID retrieves join requests for a user
func (r *GroupRepositoryImpl) GetJoinRequestsByUserID(userID uint, status string, page, pageSize int) ([]models.GroupJoinRequest, int64, error) {
	var requests []models.GroupJoinRequest
	var total int64

	query := r.db.Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count total results
	err := query.Model(&models.GroupJoinRequest{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&requests).Error

	if err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

// UpdateJoinRequest updates a join request
func (r *GroupRepositoryImpl) UpdateJoinRequest(request *models.GroupJoinRequest) error {
	return r.db.Save(request).Error
}

// DeleteJoinRequest deletes a join request
func (r *GroupRepositoryImpl) DeleteJoinRequest(id uint) error {
	return r.db.Delete(&models.GroupJoinRequest{}, id).Error
}
