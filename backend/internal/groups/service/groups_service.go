package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/groups/models"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/groups/repository"
)

// GroupService defines the interface for group-related business logic
type GroupService interface {
	// Group operations
	CreateGroup(group *models.Group) (*models.Group, error)
	GetGroupByID(id uint) (*models.Group, error)
	GetGroupsByUserID(userID uint, includeInvited bool) ([]models.Group, error)
	GetGroupsByLocation(latitude, longitude float64, radiusKm float64) ([]models.Group, error)
	UpdateGroup(id uint, group *models.Group, userID uint) (*models.Group, error)
	DeleteGroup(id uint, userID uint) error
	SearchGroups(query string, groupType string, tags []string, page, pageSize int) ([]models.Group, int64, error)

	// Member operations
	AddMember(groupID, userID uint, role models.MemberRole, status models.MemberStatus, invitedBy *uint) (*models.GroupMember, error)
	GetMemberByID(groupID, userID uint) (*models.GroupMember, error)
	GetMembersByGroupID(groupID uint, status string, page, pageSize int) ([]models.GroupMember, int64, error)
	UpdateMember(groupID, userID uint, role models.MemberRole, status models.MemberStatus, updatedByID uint) (*models.GroupMember, error)
	RemoveMember(groupID, userID, removedByID uint) error

	// Event operations
	CreateEvent(event *models.LocalEvent) (*models.LocalEvent, error)
	GetEventByID(id uint) (*models.LocalEvent, error)
	GetEventsByGroupID(groupID uint, status string, page, pageSize int) ([]models.LocalEvent, int64, error)
	GetUpcomingEvents(userID uint, page, pageSize int) ([]models.LocalEvent, int64, error)
	UpdateEvent(id uint, event *models.LocalEvent, userID uint) (*models.LocalEvent, error)
	DeleteEvent(id uint, userID uint) error

	// Attendee operations
	AddAttendee(eventID, userID uint, status string) (*models.EventAttendee, error)
	GetAttendeeByID(eventID, userID uint) (*models.EventAttendee, error)
	GetAttendeesByEventID(eventID uint, status string, page, pageSize int) ([]models.EventAttendee, int64, error)
	UpdateAttendee(eventID, userID uint, status string) (*models.EventAttendee, error)
	RemoveAttendee(eventID, userID uint) error
	CheckInAttendee(eventID, userID, checkedInByID uint) (*models.EventAttendee, error)

	// Resource operations
	CreateResource(resource *models.SharedResource) (*models.SharedResource, error)
	GetResourceByID(id uint) (*models.SharedResource, error)
	GetResourcesByGroupID(groupID uint, resourceType string, page, pageSize int) ([]models.SharedResource, int64, error)
	GetResourcesByEventID(eventID uint, page, pageSize int) ([]models.SharedResource, int64, error)
	GetResourcesByActionID(actionID uint, page, pageSize int) ([]models.SharedResource, int64, error)
	UpdateResource(id uint, resource *models.SharedResource, userID uint) (*models.SharedResource, error)
	DeleteResource(id uint, userID uint) error

	// Discussion operations
	CreateDiscussion(discussion *models.GroupDiscussion) (*models.GroupDiscussion, error)
	GetDiscussionByID(id uint) (*models.GroupDiscussion, error)
	GetDiscussionsByGroupID(groupID uint, page, pageSize int) ([]models.GroupDiscussion, int64, error)
	UpdateDiscussion(id uint, discussion *models.GroupDiscussion, userID uint) (*models.GroupDiscussion, error)
	DeleteDiscussion(id uint, userID uint) error

	// Comment operations
	CreateComment(comment *models.GroupComment) (*models.GroupComment, error)
	GetCommentByID(id uint) (*models.GroupComment, error)
	GetCommentsByDiscussionID(discussionID uint, page, pageSize int) ([]models.GroupComment, int64, error)
	UpdateComment(id uint, content string, userID uint) (*models.GroupComment, error)
	DeleteComment(id uint, userID uint) error

	// Action operations
	CreateAction(action *models.GroupAction) (*models.GroupAction, error)
	GetActionByID(id uint) (*models.GroupAction, error)
	GetActionsByGroupID(groupID uint, status string, page, pageSize int) ([]models.GroupAction, int64, error)
	GetActionsByAssigneeID(userID uint, status string, page, pageSize int) ([]models.GroupAction, int64, error)
	UpdateAction(id uint, action *models.GroupAction, userID uint) (*models.GroupAction, error)
	DeleteAction(id uint, userID uint) error

	// Invitation operations
	CreateInvitation(groupID uint, email string, userID *uint, invitedByID uint, message string) (*models.GroupInvitation, error)
	GetInvitationByID(id uint) (*models.GroupInvitation, error)
	GetInvitationByCode(code string) (*models.GroupInvitation, error)
	GetInvitationsByGroupID(groupID uint, status string, page, pageSize int) ([]models.GroupInvitation, int64, error)
	GetInvitationsByUserID(userID uint, status string, page, pageSize int) ([]models.GroupInvitation, int64, error)
	AcceptInvitation(code string, userID uint) error
	DeclineInvitation(code string, userID uint) error
	CancelInvitation(id uint, cancelledByID uint) error

	// Join request operations
	CreateJoinRequest(groupID, userID uint, message, code string) (*models.GroupJoinRequest, error)
	GetJoinRequestByID(id uint) (*models.GroupJoinRequest, error)
	GetJoinRequestsByGroupID(groupID uint, status string, page, pageSize int) ([]models.GroupJoinRequest, int64, error)
	GetJoinRequestsByUserID(userID uint, status string, page, pageSize int) ([]models.GroupJoinRequest, int64, error)
	ApproveJoinRequest(id uint, approvedByID uint) error
	RejectJoinRequest(id uint, rejectedByID uint) error
	CancelJoinRequest(id uint, userID uint) error
}

// GroupServiceImpl implements the GroupService interface
type GroupServiceImpl struct {
	groupRepo repository.GroupRepository
}

// NewGroupService creates a new group service
func NewGroupService(groupRepo repository.GroupRepository) GroupService {
	return &GroupServiceImpl{
		groupRepo: groupRepo,
	}
}

// CreateGroup creates a new group
func (s *GroupServiceImpl) CreateGroup(group *models.Group) (*models.Group, error) {
	// Validate group data
	if group.Name == "" {
		return nil, errors.New("group name is required")
	}

	if group.Type == "" {
		return nil, errors.New("group type is required")
	}

	if group.CreatedByID == 0 {
		return nil, errors.New("creator ID is required")
	}

	// Set default visibility if not provided
	if group.Visibility == "" {
		group.Visibility = models.PublicGroup
	}

	// Create the group
	err := s.groupRepo.CreateGroup(group)
	if err != nil {
		return nil, err
	}

	// Add the creator as an owner
	_, err = s.AddMember(group.ID, group.CreatedByID, models.OwnerRole, models.ActiveMember, nil)
	if err != nil {
		return nil, err
	}

	return group, nil
}

// GetGroupByID retrieves a group by its ID
func (s *GroupServiceImpl) GetGroupByID(id uint) (*models.Group, error) {
	return s.groupRepo.GetGroupByID(id)
}

// GetGroupsByUserID retrieves groups that a user is a member of
func (s *GroupServiceImpl) GetGroupsByUserID(userID uint, includeInvited bool) ([]models.Group, error) {
	return s.groupRepo.GetGroupsByUserID(userID, includeInvited)
}

// GetGroupsByLocation retrieves groups near a specific location
func (s *GroupServiceImpl) GetGroupsByLocation(latitude, longitude float64, radiusKm float64) ([]models.Group, error) {
	return s.groupRepo.GetGroupsByLocation(latitude, longitude, radiusKm)
}

// UpdateGroup updates a group
func (s *GroupServiceImpl) UpdateGroup(id uint, group *models.Group, userID uint) (*models.Group, error) {
	// Check if the group exists
	existingGroup, err := s.groupRepo.GetGroupByID(id)
	if err != nil {
		return nil, err
	}

	// Check if the user is an owner or admin
	member, err := s.groupRepo.GetMemberByID(id, userID)
	if err != nil || (member.Role != models.OwnerRole && member.Role != models.AdminRole) {
		return nil, errors.New("unauthorized: only owners and admins can update the group")
	}

	// Update only allowed fields
	existingGroup.Name = group.Name
	existingGroup.Description = group.Description
	existingGroup.Location = group.Location
	existingGroup.Latitude = group.Latitude
	existingGroup.Longitude = group.Longitude
	existingGroup.IsVirtual = group.IsVirtual
	existingGroup.JoinApprovalRequired = group.JoinApprovalRequired
	existingGroup.MembershipCode = group.MembershipCode
	existingGroup.MaxMembers = group.MaxMembers
	existingGroup.BannerImage = group.BannerImage
	existingGroup.AvatarImage = group.AvatarImage
	existingGroup.Tags = group.Tags

	// Only owners can change visibility
	if member.Role == models.OwnerRole {
		existingGroup.Visibility = group.Visibility
	}

	err = s.groupRepo.UpdateGroup(existingGroup)
	if err != nil {
		return nil, err
	}

	return existingGroup, nil
}

// DeleteGroup deletes a group
func (s *GroupServiceImpl) DeleteGroup(id uint, userID uint) error {
	// Check if the group exists
	_, err := s.groupRepo.GetGroupByID(id)
	if err != nil {
		return err
	}

	// Check if the user is an owner
	member, err := s.groupRepo.GetMemberByID(id, userID)
	if err != nil || member.Role != models.OwnerRole {
		return errors.New("unauthorized: only owners can delete the group")
	}

	return s.groupRepo.DeleteGroup(id)
}

// SearchGroups searches for groups based on criteria
func (s *GroupServiceImpl) SearchGroups(query string, groupType string, tags []string, page, pageSize int) ([]models.Group, int64, error) {
	return s.groupRepo.SearchGroups(query, groupType, tags, page, pageSize)
}

// AddMember adds a member to a group
func (s *GroupServiceImpl) AddMember(groupID, userID uint, role models.MemberRole, status models.MemberStatus, invitedBy *uint) (*models.GroupMember, error) {
	// Check if the group exists
	group, err := s.groupRepo.GetGroupByID(groupID)
	if err != nil {
		return nil, err
	}

	// Check if the user is already a member
	existingMember, err := s.groupRepo.GetMemberByID(groupID, userID)
	if err == nil && existingMember != nil {
		return nil, errors.New("user is already a member of this group")
	}

	// Check if the group has reached its maximum members
	if group.MaxMembers > 0 {
		members, _, err := s.groupRepo.GetMembersByGroupID(groupID, string(models.ActiveMember), 1, 1000)
		if err != nil {
			return nil, err
		}

		if len(members) >= group.MaxMembers {
			return nil, errors.New("group has reached its maximum number of members")
		}
	}

	// Create the member
	now := time.Now()
	member := &models.GroupMember{
		GroupID:   groupID,
		UserID:    userID,
		Role:      role,
		Status:    status,
		InvitedBy: invitedBy,
	}

	// Set joined time if status is active
	if status == models.ActiveMember {
		member.JoinedAt = &now
	}

	err = s.groupRepo.AddMember(member)
	if err != nil {
		return nil, err
	}

	return member, nil
}

// GetMemberByID retrieves a group member by group ID and user ID
func (s *GroupServiceImpl) GetMemberByID(groupID, userID uint) (*models.GroupMember, error) {
	return s.groupRepo.GetMemberByID(groupID, userID)
}

// GetMembersByGroupID retrieves members of a group
func (s *GroupServiceImpl) GetMembersByGroupID(groupID uint, status string, page, pageSize int) ([]models.GroupMember, int64, error) {
	return s.groupRepo.GetMembersByGroupID(groupID, status, page, pageSize)
}

// UpdateMember updates a group member
func (s *GroupServiceImpl) UpdateMember(groupID, userID uint, role models.MemberRole, status models.MemberStatus, updatedByID uint) (*models.GroupMember, error) {
	// Check if the group exists
	_, err := s.groupRepo.GetGroupByID(groupID)
	if err != nil {
		return nil, err
	}

	// Check if the member exists
	member, err := s.groupRepo.GetMemberByID(groupID, userID)
	if err != nil {
		return nil, err
	}

	// Check if the updater is an owner or admin
	updater, err := s.groupRepo.GetMemberByID(groupID, updatedByID)
	if err != nil || (updater.Role != models.OwnerRole && updater.Role != models.AdminRole) {
		return nil, errors.New("unauthorized: only owners and admins can update members")
	}

	// Admins cannot update owners
	if updater.Role == models.AdminRole && member.Role == models.OwnerRole {
		return nil, errors.New("unauthorized: admins cannot update owners")
	}

	// Update the member
	member.Role = role
	member.Status = status

	// Set joined time if status is changing to active
	if status == models.ActiveMember && member.JoinedAt == nil {
		now := time.Now()
		member.JoinedAt = &now
	}

	err = s.groupRepo.UpdateMember(member)
	if err != nil {
		return nil, err
	}

	return member, nil
}

// RemoveMember removes a member from a group
func (s *GroupServiceImpl) RemoveMember(groupID, userID, removedByID uint) error {
	// Check if the group exists
	_, err := s.groupRepo.GetGroupByID(groupID)
	if err != nil {
		return err
	}

	// Check if the member exists
	member, err := s.groupRepo.GetMemberByID(groupID, userID)
	if err != nil {
		return err
	}

	// Check if the remover is an owner or admin, or if the user is removing themselves
	if removedByID != userID {
		remover, err := s.groupRepo.GetMemberByID(groupID, removedByID)
		if err != nil || (remover.Role != models.OwnerRole && remover.Role != models.AdminRole) {
			return errors.New("unauthorized: only owners and admins can remove members")
		}

		// Admins cannot remove owners
		if remover.Role == models.AdminRole && member.Role == models.OwnerRole {
			return errors.New("unauthorized: admins cannot remove owners")
		}
	}

	return s.groupRepo.RemoveMember(groupID, userID)
}

// CreateEvent creates a new local event
func (s *GroupServiceImpl) CreateEvent(event *models.LocalEvent) (*models.LocalEvent, error) {
	// Validate event data
	if event.Title == "" {
		return nil, errors.New("event title is required")
	}

	if event.GroupID == 0 {
		return nil, errors.New("group ID is required")
	}

	if event.CreatedByID == 0 {
		return nil, errors.New("creator ID is required")
	}

	// Check if the group exists
	_, err := s.groupRepo.GetGroupByID(event.GroupID)
	if err != nil {
		return nil, err
	}

	// Check if the creator is a member of the group
	member, err := s.groupRepo.GetMemberByID(event.GroupID, event.CreatedByID)
	if err != nil || member.Status != models.ActiveMember {
		return nil, errors.New("unauthorized: only active members can create events")
	}

	// Set default status if not provided
	if event.Status == "" {
		event.Status = models.PlannedEvent
	}

	// Create the event
	err = s.groupRepo.CreateEvent(event)
	if err != nil {
		return nil, err
	}

	// Add the creator as an attendee
	_, err = s.AddAttendee(event.ID, event.CreatedByID, "going")
	if err != nil {
		return nil, err
	}

	return event, nil
}

// GetEventByID retrieves an event by its ID
func (s *GroupServiceImpl) GetEventByID(id uint) (*models.LocalEvent, error) {
	return s.groupRepo.GetEventByID(id)
}

// GetEventsByGroupID retrieves events for a group
func (s *GroupServiceImpl) GetEventsByGroupID(groupID uint, status string, page, pageSize int) ([]models.LocalEvent, int64, error) {
	return s.groupRepo.GetEventsByGroupID(groupID, status, page, pageSize)
}

// GetUpcomingEvents retrieves upcoming events for a user
func (s *GroupServiceImpl) GetUpcomingEvents(userID uint, page, pageSize int) ([]models.LocalEvent, int64, error) {
	return s.groupRepo.GetUpcomingEvents(userID, page, pageSize)
}

// UpdateEvent updates an event
func (s *GroupServiceImpl) UpdateEvent(id uint, event *models.LocalEvent, userID uint) (*models.LocalEvent, error) {
	// Check if the event exists
	existingEvent, err := s.groupRepo.GetEventByID(id)
	if err != nil {
		return nil, err
	}

	// Check if the user is the creator or a group admin/owner
	if existingEvent.CreatedByID != userID {
		member, err := s.groupRepo.GetMemberByID(existingEvent.GroupID, userID)
		if err != nil || (member.Role != models.OwnerRole && member.Role != models.AdminRole) {
			return nil, errors.New("unauthorized: only the creator or group admins/owners can update the event")
		}
	}

	// Update only allowed fields
	existingEvent.Title = event.Title
	existingEvent.Description = event.Description
	existingEvent.Location = event.Location
	existingEvent.Latitude = event.Latitude
	existingEvent.Longitude = event.Longitude
	existingEvent.IsVirtual = event.IsVirtual
	existingEvent.VirtualLink = event.VirtualLink
	existingEvent.StartTime = event.StartTime
	existingEvent.EndTime = event.EndTime
	existingEvent.Status = event.Status
	existingEvent.MaxAttendees = event.MaxAttendees

	err = s.groupRepo.UpdateEvent(existingEvent)
	if err != nil {
		return nil, err
	}

	return existingEvent, nil
}

// DeleteEvent deletes an event
func (s *GroupServiceImpl) DeleteEvent(id uint, userID uint) error {
	// Check if the event exists
	event, err := s.groupRepo.GetEventByID(id)
	if err != nil {
		return err
	}

	// Check if the user is the creator or a group admin/owner
	if event.CreatedByID != userID {
		member, err := s.groupRepo.GetMemberByID(event.GroupID, userID)
		if err != nil || (member.Role != models.OwnerRole && member.Role != models.AdminRole) {
			return errors.New("unauthorized: only the creator or group admins/owners can delete the event")
		}
	}

	return s.groupRepo.DeleteEvent(id)
}

// AddAttendee adds an attendee to an event
func (s *GroupServiceImpl) AddAttendee(eventID, userID uint, status string) (*models.EventAttendee, error) {
	// Check if the event exists
	event, err := s.groupRepo.GetEventByID(eventID)
	if err != nil {
		return nil, err
	}

	// Check if the user is a member of the group
	member, err := s.groupRepo.GetMemberByID(event.GroupID, userID)
	if err != nil || member.Status != models.ActiveMember {
		return nil, errors.New("unauthorized: only active group members can attend events")
	}

	// Check if the user is already an attendee
	existingAttendee, err := s.groupRepo.GetAttendeeByID(eventID, userID)
	if err == nil && existingAttendee != nil {
		return nil, errors.New("user is already an attendee of this event")
	}

	// Check if the event has reached its maximum attendees
	if event.MaxAttendees > 0 && status == "going" {
		attendees, _, err := s.groupRepo.GetAttendeesByEventID(eventID, "going", 1, 1000)
		if err != nil {
			return nil, err
		}

		if len(attendees) >= event.MaxAttendees {
			return nil, errors.New("event has reached its maximum number of attendees")
		}
	}

	// Create the attendee
	attendee := &models.EventAttendee{
		EventID:  eventID,
		UserID:   userID,
		Status:   status,
		RSVPTime: time.Now(),
	}

	err = s.groupRepo.AddAttendee(attendee)
	if err != nil {
		return nil, err
	}

	return attendee, nil
}

// GetAttendeeByID retrieves an attendee by event ID and user ID
func (s *GroupServiceImpl) GetAttendeeByID(eventID, userID uint) (*models.EventAttendee, error) {
	return s.groupRepo.GetAttendeeByID(eventID, userID)
}

// GetAttendeesByEventID retrieves attendees for an event
func (s *GroupServiceImpl) GetAttendeesByEventID(eventID uint, status string, page, pageSize int) ([]models.EventAttendee, int64, error) {
	return s.groupRepo.GetAttendeesByEventID(eventID, status, page, pageSize)
}

// UpdateAttendee updates an event attendee
func (s *GroupServiceImpl) UpdateAttendee(eventID, userID uint, status string) (*models.EventAttendee, error) {
	// Check if the event exists
	event, err := s.groupRepo.GetEventByID(eventID)
	if err != nil {
		return nil, err
	}

	// Check if the attendee exists
	attendee, err := s.groupRepo.GetAttendeeByID(eventID, userID)
	if err != nil {
		return nil, err
	}

	// Check if the event has reached its maximum attendees
	if event.MaxAttendees > 0 && status == "going" && attendee.Status != "going" {
		attendees, _, err := s.groupRepo.GetAttendeesByEventID(eventID, "going", 1, 1000)
		if err != nil {
			return nil, err
		}

		if len(attendees) >= event.MaxAttendees {
			return nil, errors.New("event has reached its maximum number of attendees")
		}
	}

	// Update the attendee
	attendee.Status = status
	attendee.RSVPTime = time.Now()

	err = s.groupRepo.UpdateAttendee(attendee)
	if err != nil {
		return nil, err
	}

	return attendee, nil
}

// RemoveAttendee removes an attendee from an event
func (s *GroupServiceImpl) RemoveAttendee(eventID, userID uint) error {
	// Check if the event exists
	_, err := s.groupRepo.GetEventByID(eventID)
	if err != nil {
		return err
	}

	// Check if the attendee exists
	_, err = s.groupRepo.GetAttendeeByID(eventID, userID)
	if err != nil {
		return err
	}

	return s.groupRepo.RemoveAttendee(eventID, userID)
}

// CheckInAttendee checks in an attendee to an event
func (s *GroupServiceImpl) CheckInAttendee(eventID, userID, checkedInByID uint) (*models.EventAttendee, error) {
	// Check if the event exists
	event, err := s.groupRepo.GetEventByID(eventID)
	if err != nil {
		return nil, err
	}

	// Check if the attendee exists
	attendee, err := s.groupRepo.GetAttendeeByID(eventID, userID)
	if err != nil {
		return nil, err
	}

	// Check if the user checking in attendees is the event creator or a group admin/owner
	if checkedInByID != event.CreatedByID {
		member, err := s.groupRepo.GetMemberByID(event.GroupID, checkedInByID)
		if err != nil || (member.Role != models.OwnerRole && member.Role != models.AdminRole && member.Role != models.ModeratorRole) {
			return nil, errors.New("unauthorized: only the event creator or group admins/moderators can check in attendees")
		}
	}

	// Update the attendee
	now := time.Now()
	attendee.CheckedIn = true
	attendee.CheckInTime = &now

	err = s.groupRepo.UpdateAttendee(attendee)
	if err != nil {
		return nil, err
	}

	return attendee, nil
}

// CreateResource creates a new shared resource
func (s *GroupServiceImpl) CreateResource(resource *models.SharedResource) (*models.SharedResource, error) {
	// Validate resource data
	if resource.Title == "" {
		return nil, errors.New("resource title is required")
	}

	if resource.Type == "" {
		return nil, errors.New("resource type is required")
	}

	if resource.CreatedByID == 0 {
		return nil, errors.New("creator ID is required")
	}

	// Check if the resource belongs to a group, event, or action
	if resource.GroupID == 0 && resource.EventID == nil && resource.ActionID == nil {
		return nil, errors.New("resource must belong to a group, event, or action")
	}

	// If the resource belongs to a group, check if the creator is a member
	if resource.GroupID != 0 {
		member, err := s.groupRepo.GetMemberByID(resource.GroupID, resource.CreatedByID)
		if err != nil || member.Status != models.ActiveMember {
			return nil, errors.New("unauthorized: only active members can create group resources")
		}
	}

	// If the resource belongs to an event, check if the creator is an attendee or the event creator
	if resource.EventID != nil {
		event, err := s.groupRepo.GetEventByID(*resource.EventID)
		if err != nil {
			return nil, err
		}

		if event.CreatedByID != resource.CreatedByID {
			attendee, err := s.groupRepo.GetAttendeeByID(*resource.EventID, resource.CreatedByID)
			if err != nil || attendee.Status != "going" {
				return nil, errors.New("unauthorized: only the event creator or attendees can create event resources")
			}
		}
	}

	// If the resource belongs to an action, check if the creator is the action creator or assignee
	if resource.ActionID != nil {
		action, err := s.groupRepo.GetActionByID(*resource.ActionID)
		if err != nil {
			return nil, err
		}

		if action.CreatedByID != resource.CreatedByID && (action.AssignedToID == nil || *action.AssignedToID != resource.CreatedByID) {
			return nil, errors.New("unauthorized: only the action creator or assignee can create action resources")
		}
	}

	// Create the resource
	err := s.groupRepo.CreateResource(resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// GetResourceByID retrieves a resource by its ID
func (s *GroupServiceImpl) GetResourceByID(id uint) (*models.SharedResource, error) {
	return s.groupRepo.GetResourceByID(id)
}

// GetResourcesByGroupID retrieves resources for a group
func (s *GroupServiceImpl) GetResourcesByGroupID(groupID uint, resourceType string, page, pageSize int) ([]models.SharedResource, int64, error) {
	return s.groupRepo.GetResourcesByGroupID(groupID, resourceType, page, pageSize)
}

// GetResourcesByEventID retrieves resources for an event
func (s *GroupServiceImpl) GetResourcesByEventID(eventID uint, page, pageSize int) ([]models.SharedResource, int64, error) {
	return s.groupRepo.GetResourcesByEventID(eventID, page, pageSize)
}

// GetResourcesByActionID retrieves resources for an action
func (s *GroupServiceImpl) GetResourcesByActionID(actionID uint, page, pageSize int) ([]models.SharedResource, int64, error) {
	return s.groupRepo.GetResourcesByActionID(actionID, page, pageSize)
}

// UpdateResource updates a resource
func (s *GroupServiceImpl) UpdateResource(id uint, resource *models.SharedResource, userID uint) (*models.SharedResource, error) {
	// Check if the resource exists
	existingResource, err := s.groupRepo.GetResourceByID(id)
	if err != nil {
		return nil, err
	}

	// Check if the user is the creator
	if existingResource.CreatedByID != userID {
		// If the resource belongs to a group, check if the user is a group admin/owner
		if existingResource.GroupID != 0 {
			member, err := s.groupRepo.GetMemberByID(existingResource.GroupID, userID)
			if err != nil || (member.Role != models.OwnerRole && member.Role != models.AdminRole) {
				return nil, errors.New("unauthorized: only the creator or group admins/owners can update the resource")
			}
		} else {
			return nil, errors.New("unauthorized: only the creator can update the resource")
		}
	}

	// Update only allowed fields
	existingResource.Title = resource.Title
	existingResource.Description = resource.Description
	existingResource.URL = resource.URL
	existingResource.FilePath = resource.FilePath
	existingResource.FileSize = resource.FileSize
	existingResource.MimeType = resource.MimeType
	existingResource.IsPublic = resource.IsPublic
	existingResource.Tags = resource.Tags

	err = s.groupRepo.UpdateResource(existingResource)
	if err != nil {
		return nil, err
	}

	return existingResource, nil
}

// DeleteResource deletes a resource
func (s *GroupServiceImpl) DeleteResource(id uint, userID uint) error {
	// Check if the resource exists
	resource, err := s.groupRepo.GetResourceByID(id)
	if err != nil {
		return err
	}

	// Check if the user is the creator
	if resource.CreatedByID != userID {
		// If the resource belongs to a group, check if the user is a group admin/owner
		if resource.GroupID != 0 {
			member, err := s.groupRepo.GetMemberByID(resource.GroupID, userID)
			if err != nil || (member.Role != models.OwnerRole && member.Role != models.AdminRole) {
				return errors.New("unauthorized: only the creator or group admins/owners can delete the resource")
			}
		} else {
			return errors.New("unauthorized: only the creator can delete the resource")
		}
	}

	return s.groupRepo.DeleteResource(id)
}

// CreateDiscussion creates a new group discussion
func (s *GroupServiceImpl) CreateDiscussion(discussion *models.GroupDiscussion) (*models.GroupDiscussion, error) {
	// Validate discussion data
	if discussion.Title == "" {
		return nil, errors.New("discussion title is required")
	}

	if discussion.GroupID == 0 {
		return nil, errors.New("group ID is required")
	}

	if discussion.CreatedByID == 0 {
		return nil, errors.New("creator ID is required")
	}

	// Check if the group exists
	_, err := s.groupRepo.GetGroupByID(discussion.GroupID)
	if err != nil {
		return nil, err
	}

	// Check if the creator is a member of the group
	member, err := s.groupRepo.GetMemberByID(discussion.GroupID, discussion.CreatedByID)
	if err != nil || member.Status != models.ActiveMember {
		return nil, errors.New("unauthorized: only active members can create discussions")
	}

	// Create the discussion
	err = s.groupRepo.CreateDiscussion(discussion)
	if err != nil {
		return nil, err
	}

	return discussion, nil
}

// GetDiscussionByID retrieves a discussion by its ID
func (s *GroupServiceImpl) GetDiscussionByID(id uint) (*models.GroupDiscussion, error) {
	return s.groupRepo.GetDiscussionByID(id)
}

// GetDiscussionsByGroupID retrieves discussions for a group
func (s *GroupServiceImpl) GetDiscussionsByGroupID(groupID uint, page, pageSize int) ([]models.GroupDiscussion, int64, error) {
	return s.groupRepo.GetDiscussionsByGroupID(groupID, page, pageSize)
}

// UpdateDiscussion updates a discussion
func (s *GroupServiceImpl) UpdateDiscussion(id uint, discussion *models.GroupDiscussion, userID uint) (*models.GroupDiscussion, error) {
	// Check if the discussion exists
	existingDiscussion, err := s.groupRepo.GetDiscussionByID(id)
	if err != nil {
		return nil, err
	}

	// Check if the user is the creator or a group admin/owner/moderator
	if existingDiscussion.CreatedByID != userID {
		member, err := s.groupRepo.GetMemberByID(existingDiscussion.GroupID, userID)
		if err != nil || (member.Role != models.OwnerRole && member.Role != models.AdminRole && member.Role != models.ModeratorRole) {
			return nil, errors.New("unauthorized: only the creator or group admins/moderators can update the discussion")
		}
	}

	// Update only allowed fields
	existingDiscussion.Title = discussion.Title
	existingDiscussion.Content = discussion.Content

	// Only admins/owners/moderators can pin or lock discussions
	if userID != existingDiscussion.CreatedByID {
		member, err := s.groupRepo.GetMemberByID(existingDiscussion.GroupID, userID)
		if err == nil && (member.Role == models.OwnerRole || member.Role == models.AdminRole || member.Role == models.ModeratorRole) {
			existingDiscussion.IsPinned = discussion.IsPinned
			existingDiscussion.IsLocked = discussion.IsLocked
		}
	}

	err = s.groupRepo.UpdateDiscussion(existingDiscussion)
	if err != nil {
		return nil, err
	}

	return existingDiscussion, nil
}

// DeleteDiscussion deletes a discussion
func (s *GroupServiceImpl) DeleteDiscussion(id uint, userID uint) error {
	// Check if the discussion exists
	discussion, err := s.groupRepo.GetDiscussionByID(id)
	if err != nil {
		return err
	}

	// Check if the user is the creator or a group admin/owner
	if discussion.CreatedByID != userID {
		member, err := s.groupRepo.GetMemberByID(discussion.GroupID, userID)
		if err != nil || (member.Role != models.OwnerRole && member.Role != models.AdminRole) {
			return errors.New("unauthorized: only the creator or group admins/owners can delete the discussion")
		}
	}

	return s.groupRepo.DeleteDiscussion(id)
}

// CreateComment creates a new comment
func (s *GroupServiceImpl) CreateComment(comment *models.GroupComment) (*models.GroupComment, error) {
	// Validate comment data
	if comment.Content == "" {
		return nil, errors.New("comment content is required")
	}

	if comment.DiscussionID == 0 {
		return nil, errors.New("discussion ID is required")
	}

	if comment.CreatedByID == 0 {
		return nil, errors.New("creator ID is required")
	}

	// Check if the discussion exists
	discussion, err := s.groupRepo.GetDiscussionByID(comment.DiscussionID)
	if err != nil {
		return nil, err
	}

	// Check if the discussion is locked
	if discussion.IsLocked {
		return nil, errors.New("discussion is locked")
	}

	// Check if the creator is a member of the group
	member, err := s.groupRepo.GetMemberByID(discussion.GroupID, comment.CreatedByID)
	if err != nil || member.Status != models.ActiveMember {
		return nil, errors.New("unauthorized: only active members can create comments")
	}

	// Create the comment
	err = s.groupRepo.CreateComment(comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

// GetCommentByID retrieves a comment by its ID
func (s *GroupServiceImpl) GetCommentByID(id uint) (*models.GroupComment, error) {
	return s.groupRepo.GetCommentByID(id)
}

// GetCommentsByDiscussionID retrieves comments for a discussion
func (s *GroupServiceImpl) GetCommentsByDiscussionID(discussionID uint, page, pageSize int) ([]models.GroupComment, int64, error) {
	return s.groupRepo.GetCommentsByDiscussionID(discussionID, page, pageSize)
}

// UpdateComment updates a comment
func (s *GroupServiceImpl) UpdateComment(id uint, content string, userID uint) (*models.GroupComment, error) {
	// Check if the comment exists
	comment, err := s.groupRepo.GetCommentByID(id)
	if err != nil {
		return nil, err
	}

	// Check if the user is the creator
	if comment.CreatedByID != userID {
		return nil, errors.New("unauthorized: only the creator can update the comment")
	}

	// Get the discussion to check if it's locked
	discussion, err := s.groupRepo.GetDiscussionByID(comment.DiscussionID)
	if err != nil {
		return nil, err
	}

	// Check if the discussion is locked
	if discussion.IsLocked {
		return nil, errors.New("discussion is locked")
	}

	// Update the comment
	comment.Content = content

	err = s.groupRepo.UpdateComment(comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

// DeleteComment deletes a comment
func (s *GroupServiceImpl) DeleteComment(id uint, userID uint) error {
	// Check if the comment exists
	comment, err := s.groupRepo.GetCommentByID(id)
	if err != nil {
		return err
	}

	// Get the discussion to check the group ID
	discussion, err := s.groupRepo.GetDiscussionByID(comment.DiscussionID)
	if err != nil {
		return err
	}

	// Check if the user is the creator or a group admin/owner/moderator
	if comment.CreatedByID != userID {
		member, err := s.groupRepo.GetMemberByID(discussion.GroupID, userID)
		if err != nil || (member.Role != models.OwnerRole && member.Role != models.AdminRole && member.Role != models.ModeratorRole) {
			return errors.New("unauthorized: only the creator or group admins/moderators can delete the comment")
		}
	}

	return s.groupRepo.DeleteComment(id)
}

// CreateAction creates a new group action
func (s *GroupServiceImpl) CreateAction(action *models.GroupAction) (*models.GroupAction, error) {
	// Validate action data
	if action.Title == "" {
		return nil, errors.New("action title is required")
	}

	if action.GroupID == 0 {
		return nil, errors.New("group ID is required")
	}

	if action.CreatedByID == 0 {
		return nil, errors.New("creator ID is required")
	}

	// Check if the group exists
	_, err := s.groupRepo.GetGroupByID(action.GroupID)
	if err != nil {
		return nil, err
	}

	// Check if the creator is a member of the group
	member, err := s.groupRepo.GetMemberByID(action.GroupID, action.CreatedByID)
	if err != nil || member.Status != models.ActiveMember {
		return nil, errors.New("unauthorized: only active members can create actions")
	}

	// If there's an assignee, check if they are a member of the group
	if action.AssignedToID != nil {
		assignee, err := s.groupRepo.GetMemberByID(action.GroupID, *action.AssignedToID)
		if err != nil || assignee.Status != models.ActiveMember {
			return nil, errors.New("assignee must be an active member of the group")
		}
	}

	// Set default status if not provided
	if action.Status == "" {
		action.Status = "pending"
	}

	// Create the action
	err = s.groupRepo.CreateAction(action)
	if err != nil {
		return nil, err
	}

	return action, nil
}

// GetActionByID retrieves an action by its ID
func (s *GroupServiceImpl) GetActionByID(id uint) (*models.GroupAction, error) {
	return s.groupRepo.GetActionByID(id)
}

// GetActionsByGroupID retrieves actions for a group
func (s *GroupServiceImpl) GetActionsByGroupID(groupID uint, status string, page, pageSize int) ([]models.GroupAction, int64, error) {
	return s.groupRepo.GetActionsByGroupID(groupID, status, page, pageSize)
}

// GetActionsByAssigneeID retrieves actions assigned to a user
func (s *GroupServiceImpl) GetActionsByAssigneeID(userID uint, status string, page, pageSize int) ([]models.GroupAction, int64, error) {
	return s.groupRepo.GetActionsByAssigneeID(userID, status, page, pageSize)
}

// UpdateAction updates an action
func (s *GroupServiceImpl) UpdateAction(id uint, action *models.GroupAction, userID uint) (*models.GroupAction, error) {
	// Check if the action exists
	existingAction, err := s.groupRepo.GetActionByID(id)
	if err != nil {
		return nil, err
	}

	// Check if the user is the creator, assignee, or a group admin/owner
	isAuthorized := existingAction.CreatedByID == userID
	if !isAuthorized && existingAction.AssignedToID != nil && *existingAction.AssignedToID == userID {
		isAuthorized = true
	}

	if !isAuthorized {
		member, err := s.groupRepo.GetMemberByID(existingAction.GroupID, userID)
		if err != nil || (member.Role != models.OwnerRole && member.Role != models.AdminRole) {
			return nil, errors.New("unauthorized: only the creator, assignee, or group admins/owners can update the action")
		}
	}

	// If there's a new assignee, check if they are a member of the group
	if action.AssignedToID != nil && (existingAction.AssignedToID == nil || *existingAction.AssignedToID != *action.AssignedToID) {
		assignee, err := s.groupRepo.GetMemberByID(existingAction.GroupID, *action.AssignedToID)
		if err != nil || assignee.Status != models.ActiveMember {
			return nil, errors.New("assignee must be an active member of the group")
		}
	}

	// Update only allowed fields
	existingAction.Title = action.Title
	existingAction.Description = action.Description
	existingAction.Status = action.Status
	existingAction.Priority = action.Priority
	existingAction.DueDate = action.DueDate
	existingAction.AssignedToID = action.AssignedToID

	// Set completed time if status is changing to completed
	if existingAction.Status == "completed" && existingAction.CompletedAt == nil {
		now := time.Now()
		existingAction.CompletedAt = &now
	}

	err = s.groupRepo.UpdateAction(existingAction)
	if err != nil {
		return nil, err
	}

	return existingAction, nil
}

// DeleteAction deletes an action
func (s *GroupServiceImpl) DeleteAction(id uint, userID uint) error {
	// Check if the action exists
	action, err := s.groupRepo.GetActionByID(id)
	if err != nil {
		return err
	}

	// Check if the user is the creator or a group admin/owner
	if action.CreatedByID != userID {
		member, err := s.groupRepo.GetMemberByID(action.GroupID, userID)
		if err != nil || (member.Role != models.OwnerRole && member.Role != models.AdminRole) {
			return errors.New("unauthorized: only the creator or group admins/owners can delete the action")
		}
	}

	return s.groupRepo.DeleteAction(id)
}

// CreateInvitation creates a new group invitation
func (s *GroupServiceImpl) CreateInvitation(groupID uint, email string, userID *uint, invitedByID uint, message string) (*models.GroupInvitation, error) {
	// Check if the group exists
	group, err := s.groupRepo.GetGroupByID(groupID)
	if err != nil {
		return nil, err
	}

	// Check if the inviter is a member of the group with appropriate permissions
	inviter, err := s.groupRepo.GetMemberByID(groupID, invitedByID)
	if err != nil || (inviter.Role != models.OwnerRole && inviter.Role != models.AdminRole && inviter.Role != models.ModeratorRole) {
		return nil, errors.New("unauthorized: only owners, admins, and moderators can send invitations")
	}

	// Check if the group has reached its maximum members
	if group.MaxMembers > 0 {
		members, _, err := s.groupRepo.GetMembersByGroupID(groupID, string(models.ActiveMember), 1, 1000)
		if err != nil {
			return nil, err
		}

		if len(members) >= group.MaxMembers {
			return nil, errors.New("group has reached its maximum number of members")
		}
	}

	// If inviting an existing user, check if they are already a member
	if userID != nil {
		existingMember, err := s.groupRepo.GetMemberByID(groupID, *userID)
		if err == nil && existingMember != nil {
			return nil, errors.New("user is already a member of this group")
		}
	}

	// Generate a unique invitation code
	code, err := generateRandomCode(32)
	if err != nil {
		return nil, err
	}

	// Create the invitation
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 days
	invitation := &models.GroupInvitation{
		GroupID:     groupID,
		Email:       email,
		UserID:      userID,
		InvitedByID: invitedByID,
		Status:      "pending",
		ExpiresAt:   expiresAt,
		Code:        code,
		Message:     message,
	}

	err = s.groupRepo.CreateInvitation(invitation)
	if err != nil {
		return nil, err
	}

	return invitation, nil
}

// GetInvitationByID retrieves an invitation by its ID
func (s *GroupServiceImpl) GetInvitationByID(id uint) (*models.GroupInvitation, error) {
	return s.groupRepo.GetInvitationByID(id)
}

// GetInvitationByCode retrieves an invitation by its code
func (s *GroupServiceImpl) GetInvitationByCode(code string) (*models.GroupInvitation, error) {
	return s.groupRepo.GetInvitationByCode(code)
}

// GetInvitationsByGroupID retrieves invitations for a group
func (s *GroupServiceImpl) GetInvitationsByGroupID(groupID uint, status string, page, pageSize int) ([]models.GroupInvitation, int64, error) {
	return s.groupRepo.GetInvitationsByGroupID(groupID, status, page, pageSize)
}

// GetInvitationsByUserID retrieves invitations for a user
func (s *GroupServiceImpl) GetInvitationsByUserID(userID uint, status string, page, pageSize int) ([]models.GroupInvitation, int64, error) {
	return s.groupRepo.GetInvitationsByUserID(userID, status, page, pageSize)
}

// AcceptInvitation accepts a group invitation
func (s *GroupServiceImpl) AcceptInvitation(code string, userID uint) error {
	// Get the invitation
	invitation, err := s.groupRepo.GetInvitationByCode(code)
	if err != nil {
		return err
	}

	// Check if the invitation is pending
	if invitation.Status != "pending" {
		return errors.New("invitation is not pending")
	}

	// Check if the invitation has expired
	if invitation.ExpiresAt.Before(time.Now()) {
		return errors.New("invitation has expired")
	}

	// Check if the invitation is for the user
	if invitation.UserID != nil && *invitation.UserID != userID {
		return errors.New("invitation is not for this user")
	}

	// Check if the group exists
	group, err := s.groupRepo.GetGroupByID(invitation.GroupID)
	if err != nil {
		return err
	}

	// Check if the group has reached its maximum members
	if group.MaxMembers > 0 {
		members, _, err := s.groupRepo.GetMembersByGroupID(invitation.GroupID, string(models.ActiveMember), 1, 1000)
		if err != nil {
			return err
		}

		if len(members) >= group.MaxMembers {
			return errors.New("group has reached its maximum number of members")
		}
	}

	// Check if the user is already a member
	existingMember, err := s.groupRepo.GetMemberByID(invitation.GroupID, userID)
	if err == nil && existingMember != nil {
		return errors.New("user is already a member of this group")
	}

	// Add the user as a member
	_, err = s.AddMember(invitation.GroupID, userID, models.RegularMemberRole, models.ActiveMember, &invitation.InvitedByID)
	if err != nil {
		return err
	}

	// Update the invitation status
	invitation.Status = "accepted"
	err = s.groupRepo.UpdateInvitation(invitation)
	if err != nil {
		return err
	}

	return nil
}

// DeclineInvitation declines a group invitation
func (s *GroupServiceImpl) DeclineInvitation(code string, userID uint) error {
	// Get the invitation
	invitation, err := s.groupRepo.GetInvitationByCode(code)
	if err != nil {
		return err
	}

	// Check if the invitation is pending
	if invitation.Status != "pending" {
		return errors.New("invitation is not pending")
	}

	// Check if the invitation is for the user
	if invitation.UserID != nil && *invitation.UserID != userID {
		return errors.New("invitation is not for this user")
	}

	// Update the invitation status
	invitation.Status = "declined"
	err = s.groupRepo.UpdateInvitation(invitation)
	if err != nil {
		return err
	}

	return nil
}

// CancelInvitation cancels a group invitation
func (s *GroupServiceImpl) CancelInvitation(id uint, cancelledByID uint) error {
	// Get the invitation
	invitation, err := s.groupRepo.GetInvitationByID(id)
	if err != nil {
		return err
	}

	// Check if the invitation is pending
	if invitation.Status != "pending" {
		return errors.New("invitation is not pending")
	}

	// Check if the user is the inviter or a group admin/owner
	if invitation.InvitedByID != cancelledByID {
		member, err := s.groupRepo.GetMemberByID(invitation.GroupID, cancelledByID)
		if err != nil || (member.Role != models.OwnerRole && member.Role != models.AdminRole) {
			return errors.New("unauthorized: only the inviter or group admins/owners can cancel the invitation")
		}
	}

	// Update the invitation status
	invitation.Status = "cancelled"
	err = s.groupRepo.UpdateInvitation(invitation)
	if err != nil {
		return err
	}

	return nil
}

// CreateJoinRequest creates a new join request
func (s *GroupServiceImpl) CreateJoinRequest(groupID, userID uint, message, code string) (*models.GroupJoinRequest, error) {
	// Check if the group exists
	group, err := s.groupRepo.GetGroupByID(groupID)
	if err != nil {
		return nil, err
	}

	// Check if the group is public or private
	if group.Visibility == models.SecretGroup {
		return nil, errors.New("cannot request to join a secret group")
	}

	// Check if the user is already a member
	existingMember, err := s.groupRepo.GetMemberByID(groupID, userID)
	if err == nil && existingMember != nil {
		return nil, errors.New("user is already a member of this group")
	}

	// Check if the user already has a pending join request
	requests, _, err := s.groupRepo.GetJoinRequestsByUserID(userID, "pending", 1, 10)
	if err != nil {
		return nil, err
	}

	for _, request := range requests {
		if request.GroupID == groupID {
			return nil, errors.New("user already has a pending join request for this group")
		}
	}

	// Check if the group requires a membership code
	if group.JoinApprovalRequired && group.MembershipCode != "" && group.MembershipCode != code {
		return nil, errors.New("invalid membership code")
	}

	// Create the join request
	request := &models.GroupJoinRequest{
		GroupID: groupID,
		UserID:  userID,
		Status:  "pending",
		Message: message,
		Code:    code,
	}

	// If the group doesn't require approval and the code is correct, auto-approve
	if !group.JoinApprovalRequired || (group.MembershipCode != "" && group.MembershipCode == code) {
		// Check if the group has reached its maximum members
		if group.MaxMembers > 0 {
			members, _, err := s.groupRepo.GetMembersByGroupID(groupID, string(models.ActiveMember), 1, 1000)
			if err != nil {
				return nil, err
			}

			if len(members) >= group.MaxMembers {
				return nil, errors.New("group has reached its maximum number of members")
			}
		}

		// Add the user as a member
		_, err = s.AddMember(groupID, userID, models.RegularMemberRole, models.ActiveMember, nil)
		if err != nil {
			return nil, err
		}

		request.Status = "approved"
		now := time.Now()
		request.ApprovedAt = &now
	}

	err = s.groupRepo.CreateJoinRequest(request)
	if err != nil {
		return nil, err
	}

	return request, nil
}

// GetJoinRequestByID retrieves a join request by its ID
func (s *GroupServiceImpl) GetJoinRequestByID(id uint) (*models.GroupJoinRequest, error) {
	return s.groupRepo.GetJoinRequestByID(id)
}

// GetJoinRequestsByGroupID retrieves join requests for a group
func (s *GroupServiceImpl) GetJoinRequestsByGroupID(groupID uint, status string, page, pageSize int) ([]models.GroupJoinRequest, int64, error) {
	return s.groupRepo.GetJoinRequestsByGroupID(groupID, status, page, pageSize)
}

// GetJoinRequestsByUserID retrieves join requests for a user
func (s *GroupServiceImpl) GetJoinRequestsByUserID(userID uint, status string, page, pageSize int) ([]models.GroupJoinRequest, int64, error) {
	return s.groupRepo.GetJoinRequestsByUserID(userID, status, page, pageSize)
}

// ApproveJoinRequest approves a join request
func (s *GroupServiceImpl) ApproveJoinRequest(id uint, approvedByID uint) error {
	// Get the join request
	request, err := s.groupRepo.GetJoinRequestByID(id)
	if err != nil {
		return err
	}

	// Check if the request is pending
	if request.Status != "pending" {
		return errors.New("request is not pending")
	}

	// Check if the user is a group admin/owner
	member, err := s.groupRepo.GetMemberByID(request.GroupID, approvedByID)
	if err != nil || (member.Role != models.OwnerRole && member.Role != models.AdminRole) {
		return errors.New("unauthorized: only group admins/owners can approve join requests")
	}

	// Check if the group exists
	group, err := s.groupRepo.GetGroupByID(request.GroupID)
	if err != nil {
		return err
	}

	// Check if the group has reached its maximum members
	if group.MaxMembers > 0 {
		members, _, err := s.groupRepo.GetMembersByGroupID(request.GroupID, string(models.ActiveMember), 1, 1000)
		if err != nil {
			return err
		}

		if len(members) >= group.MaxMembers {
			return errors.New("group has reached its maximum number of members")
		}
	}

	// Add the user as a member
	_, err = s.AddMember(request.GroupID, request.UserID, models.RegularMemberRole, models.ActiveMember, &approvedByID)
	if err != nil {
		return err
	}

	// Update the request status
	request.Status = "approved"
	request.ApprovedBy = &approvedByID
	now := time.Now()
	request.ApprovedAt = &now

	err = s.groupRepo.UpdateJoinRequest(request)
	if err != nil {
		return err
	}

	return nil
}

// RejectJoinRequest rejects a join request
func (s *GroupServiceImpl) RejectJoinRequest(id uint, rejectedByID uint) error {
	// Get the join request
	request, err := s.groupRepo.GetJoinRequestByID(id)
	if err != nil {
		return err
	}

	// Check if the request is pending
	if request.Status != "pending" {
		return errors.New("request is not pending")
	}

	// Check if the user is a group admin/owner
	member, err := s.groupRepo.GetMemberByID(request.GroupID, rejectedByID)
	if err != nil || (member.Role != models.OwnerRole && member.Role != models.AdminRole) {
		return errors.New("unauthorized: only group admins/owners can reject join requests")
	}

	// Update the request status
	request.Status = "rejected"
	request.RejectedBy = &rejectedByID
	now := time.Now()
	request.RejectedAt = &now

	err = s.groupRepo.UpdateJoinRequest(request)
	if err != nil {
		return err
	}

	return nil
}

// CancelJoinRequest cancels a join request
func (s *GroupServiceImpl) CancelJoinRequest(id uint, userID uint) error {
	// Get the join request
	request, err := s.groupRepo.GetJoinRequestByID(id)
	if err != nil {
		return err
	}

	// Check if the request is pending
	if request.Status != "pending" {
		return errors.New("request is not pending")
	}

	// Check if the user is the requester
	if request.UserID != userID {
		return errors.New("unauthorized: only the requester can cancel the join request")
	}

	// Delete the request
	err = s.groupRepo.DeleteJoinRequest(id)
	if err != nil {
		return err
	}

	return nil
}

// Helper function to generate a random code
func generateRandomCode(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}
