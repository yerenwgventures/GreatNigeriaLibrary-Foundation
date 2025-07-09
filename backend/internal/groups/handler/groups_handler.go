package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/groups/models"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/groups/service"
)

// GroupHandler handles HTTP requests for group operations
type GroupHandler struct {
	groupService service.GroupService
}

// NewGroupHandler creates a new group handler
func NewGroupHandler(groupService service.GroupService) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
	}
}

// RegisterRoutes registers the group routes
func (h *GroupHandler) RegisterRoutes(router *gin.RouterGroup) {
	groups := router.Group("/groups")
	{
		// Group operations
		groups.GET("", h.GetGroups)
		groups.GET("/:id", h.GetGroupByID)
		groups.POST("", h.CreateGroup)
		groups.PUT("/:id", h.UpdateGroup)
		groups.DELETE("/:id", h.DeleteGroup)
		groups.GET("/nearby", h.GetNearbyGroups)
		groups.GET("/search", h.SearchGroups)

		// Member operations
		groups.GET("/:id/members", h.GetGroupMembers)
		groups.POST("/:id/members", h.AddMember)
		groups.GET("/:id/members/:userId", h.GetGroupMember)
		groups.PUT("/:id/members/:userId", h.UpdateMember)
		groups.DELETE("/:id/members/:userId", h.RemoveMember)

		// Event operations
		groups.GET("/:id/events", h.GetGroupEvents)
		groups.POST("/:id/events", h.CreateEvent)
		groups.GET("/events/:eventId", h.GetEventByID)
		groups.PUT("/events/:eventId", h.UpdateEvent)
		groups.DELETE("/events/:eventId", h.DeleteEvent)

		// Attendee operations
		groups.GET("/events/:eventId/attendees", h.GetEventAttendees)
		groups.POST("/events/:eventId/attendees", h.AddAttendee)
		groups.GET("/events/:eventId/attendees/:userId", h.GetEventAttendee)
		groups.PUT("/events/:eventId/attendees/:userId", h.UpdateAttendee)
		groups.DELETE("/events/:eventId/attendees/:userId", h.RemoveAttendee)
		groups.POST("/events/:eventId/attendees/:userId/checkin", h.CheckInAttendee)

		// Resource operations
		groups.GET("/:id/resources", h.GetGroupResources)
		groups.POST("/:id/resources", h.CreateGroupResource)
		groups.GET("/events/:eventId/resources", h.GetEventResources)
		groups.POST("/events/:eventId/resources", h.CreateEventResource)
		groups.GET("/actions/:actionId/resources", h.GetActionResources)
		groups.POST("/actions/:actionId/resources", h.CreateActionResource)
		groups.GET("/resources/:resourceId", h.GetResourceByID)
		groups.PUT("/resources/:resourceId", h.UpdateResource)
		groups.DELETE("/resources/:resourceId", h.DeleteResource)

		// Discussion operations
		groups.GET("/:id/discussions", h.GetGroupDiscussions)
		groups.POST("/:id/discussions", h.CreateDiscussion)
		groups.GET("/discussions/:discussionId", h.GetDiscussionByID)
		groups.PUT("/discussions/:discussionId", h.UpdateDiscussion)
		groups.DELETE("/discussions/:discussionId", h.DeleteDiscussion)

		// Comment operations
		groups.GET("/discussions/:discussionId/comments", h.GetDiscussionComments)
		groups.POST("/discussions/:discussionId/comments", h.CreateComment)
		groups.GET("/comments/:commentId", h.GetCommentByID)
		groups.PUT("/comments/:commentId", h.UpdateComment)
		groups.DELETE("/comments/:commentId", h.DeleteComment)

		// Action operations
		groups.GET("/:id/actions", h.GetGroupActions)
		groups.POST("/:id/actions", h.CreateAction)
		groups.GET("/actions/:actionId", h.GetActionByID)
		groups.PUT("/actions/:actionId", h.UpdateAction)
		groups.DELETE("/actions/:actionId", h.DeleteAction)

		// Invitation operations
		groups.GET("/:id/invitations", h.GetGroupInvitations)
		groups.POST("/:id/invitations", h.CreateInvitation)
		groups.GET("/invitations/:invitationId", h.GetInvitationByID)
		groups.GET("/invitations/code/:code", h.GetInvitationByCode)
		groups.POST("/invitations/code/:code/accept", h.AcceptInvitation)
		groups.POST("/invitations/code/:code/decline", h.DeclineInvitation)
		groups.DELETE("/invitations/:invitationId", h.CancelInvitation)

		// Join request operations
		groups.GET("/:id/join-requests", h.GetGroupJoinRequests)
		groups.POST("/:id/join-requests", h.CreateJoinRequest)
		groups.GET("/join-requests/:requestId", h.GetJoinRequestByID)
		groups.POST("/join-requests/:requestId/approve", h.ApproveJoinRequest)
		groups.POST("/join-requests/:requestId/reject", h.RejectJoinRequest)
		groups.DELETE("/join-requests/:requestId", h.CancelJoinRequest)

		// User-specific operations
		groups.GET("/user/:userId", h.GetUserGroups)
		groups.GET("/user/:userId/events/upcoming", h.GetUserUpcomingEvents)
		groups.GET("/user/:userId/actions", h.GetUserActions)
		groups.GET("/user/:userId/invitations", h.GetUserInvitations)
		groups.GET("/user/:userId/join-requests", h.GetUserJoinRequests)
	}
}

// GetGroups handles GET /groups
func (h *GroupHandler) GetGroups(c *gin.Context) {
	// Get query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// Search for groups
	groups, total, err := h.groupService.SearchGroups("", "", nil, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": groups,
		"meta": gin.H{
			"page":       page,
			"pageSize":   pageSize,
			"total":      total,
			"totalPages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetGroupByID handles GET /groups/:id
func (h *GroupHandler) GetGroupByID(c *gin.Context) {
	// Get group ID from path
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// Get group
	group, err := h.groupService.GetGroupByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	c.JSON(http.StatusOK, group)
}

// CreateGroup handles POST /groups
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Bind request body
	var group models.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set creator ID
	group.CreatedByID = userID.(uint)

	// Create group
	createdGroup, err := h.groupService.CreateGroup(&group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdGroup)
}

// UpdateGroup handles PUT /groups/:id
func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get group ID from path
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// Bind request body
	var group models.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update group
	updatedGroup, err := h.groupService.UpdateGroup(uint(id), &group, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedGroup)
}

// DeleteGroup handles DELETE /groups/:id
func (h *GroupHandler) DeleteGroup(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get group ID from path
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// Delete group
	err = h.groupService.DeleteGroup(uint(id), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group deleted successfully"})
}

// GetNearbyGroups handles GET /groups/nearby
func (h *GroupHandler) GetNearbyGroups(c *gin.Context) {
	// Get query parameters
	latitude, _ := strconv.ParseFloat(c.Query("latitude"), 64)
	longitude, _ := strconv.ParseFloat(c.Query("longitude"), 64)
	radius, _ := strconv.ParseFloat(c.DefaultQuery("radius", "10"), 64) // Default 10km

	// Validate coordinates
	if latitude == 0 && longitude == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Latitude and longitude are required"})
		return
	}

	// Get nearby groups
	groups, err := h.groupService.GetGroupsByLocation(latitude, longitude, radius)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, groups)
}

// SearchGroups handles GET /groups/search
func (h *GroupHandler) SearchGroups(c *gin.Context) {
	// Get query parameters
	query := c.Query("q")
	groupType := c.Query("type")
	tagsStr := c.Query("tags")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// Parse tags
	var tags []string
	if tagsStr != "" {
		tags = strings.Split(tagsStr, ",")
	}

	// Search for groups
	groups, total, err := h.groupService.SearchGroups(query, groupType, tags, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": groups,
		"meta": gin.H{
			"page":       page,
			"pageSize":   pageSize,
			"total":      total,
			"totalPages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetUserGroups handles GET /groups/user/:userId
func (h *GroupHandler) GetUserGroups(c *gin.Context) {
	// Get user ID from path
	id, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get include invited parameter
	includeInvited, _ := strconv.ParseBool(c.DefaultQuery("includeInvited", "false"))

	// Get user groups
	groups, err := h.groupService.GetGroupsByUserID(uint(id), includeInvited)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, groups)
}

// GetGroupMembers handles GET /groups/:id/members
func (h *GroupHandler) GetGroupMembers(c *gin.Context) {
	// Get group ID from path
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// Get query parameters
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// Get group members
	members, total, err := h.groupService.GetMembersByGroupID(uint(id), status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": members,
		"meta": gin.H{
			"page":       page,
			"pageSize":   pageSize,
			"total":      total,
			"totalPages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// AddMember handles POST /groups/:id/members
func (h *GroupHandler) AddMember(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get group ID from path
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// Bind request body
	var request struct {
		UserID uint                `json:"userId" binding:"required"`
		Role   models.MemberRole   `json:"role" binding:"required"`
		Status models.MemberStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Add member
	uid := userID.(uint)
	member, err := h.groupService.AddMember(uint(groupID), request.UserID, request.Role, request.Status, &uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, member)
}

// GetGroupMember handles GET /groups/:id/members/:userId
func (h *GroupHandler) GetGroupMember(c *gin.Context) {
	// Get group ID and user ID from path
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get group member
	member, err := h.groupService.GetMemberByID(uint(groupID), uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}

	c.JSON(http.StatusOK, member)
}

// UpdateMember handles PUT /groups/:id/members/:userId
func (h *GroupHandler) UpdateMember(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get group ID and member user ID from path
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	memberUserID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Bind request body
	var request struct {
		Role   models.MemberRole   `json:"role" binding:"required"`
		Status models.MemberStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update member
	member, err := h.groupService.UpdateMember(uint(groupID), uint(memberUserID), request.Role, request.Status, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, member)
}

// RemoveMember handles DELETE /groups/:id/members/:userId
func (h *GroupHandler) RemoveMember(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get group ID and member user ID from path
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	memberUserID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Remove member
	err = h.groupService.RemoveMember(uint(groupID), uint(memberUserID), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member removed successfully"})
}

// GetGroupEvents handles GET /groups/:id/events
func (h *GroupHandler) GetGroupEvents(c *gin.Context) {
	// Get group ID from path
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// Get query parameters
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// Get group events
	events, total, err := h.groupService.GetEventsByGroupID(uint(id), status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": events,
		"meta": gin.H{
			"page":       page,
			"pageSize":   pageSize,
			"total":      total,
			"totalPages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// CreateEvent handles POST /groups/:id/events
func (h *GroupHandler) CreateEvent(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get group ID from path
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// Bind request body
	var event models.LocalEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set group ID and creator ID
	event.GroupID = uint(groupID)
	event.CreatedByID = userID.(uint)

	// Create event
	createdEvent, err := h.groupService.CreateEvent(&event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdEvent)
}

// GetEventByID handles GET /groups/events/:eventId
func (h *GroupHandler) GetEventByID(c *gin.Context) {
	// Get event ID from path
	id, err := strconv.ParseUint(c.Param("eventId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	// Get event
	event, err := h.groupService.GetEventByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	c.JSON(http.StatusOK, event)
}

// UpdateEvent handles PUT /groups/events/:eventId
func (h *GroupHandler) UpdateEvent(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get event ID from path
	id, err := strconv.ParseUint(c.Param("eventId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	// Bind request body
	var event models.LocalEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update event
	updatedEvent, err := h.groupService.UpdateEvent(uint(id), &event, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedEvent)
}

// DeleteEvent handles DELETE /groups/events/:eventId
func (h *GroupHandler) DeleteEvent(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get event ID from path
	id, err := strconv.ParseUint(c.Param("eventId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	// Delete event
	err = h.groupService.DeleteEvent(uint(id), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}

// GetUserUpcomingEvents handles GET /groups/user/:userId/events/upcoming
func (h *GroupHandler) GetUserUpcomingEvents(c *gin.Context) {
	// Get user ID from path
	id, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// Get upcoming events
	events, total, err := h.groupService.GetUpcomingEvents(uint(id), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": events,
		"meta": gin.H{
			"page":       page,
			"pageSize":   pageSize,
			"total":      total,
			"totalPages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetEventAttendees handles GET /groups/events/:eventId/attendees
func (h *GroupHandler) GetEventAttendees(c *gin.Context) {
	// Get event ID from path
	id, err := strconv.ParseUint(c.Param("eventId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	// Get query parameters
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// Get event attendees
	attendees, total, err := h.groupService.GetAttendeesByEventID(uint(id), status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": attendees,
		"meta": gin.H{
			"page":       page,
			"pageSize":   pageSize,
			"total":      total,
			"totalPages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// AddAttendee handles POST /groups/events/:eventId/attendees
func (h *GroupHandler) AddAttendee(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get event ID from path
	eventID, err := strconv.ParseUint(c.Param("eventId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	// Bind request body
	var request struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Add attendee
	attendee, err := h.groupService.AddAttendee(uint(eventID), userID.(uint), request.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, attendee)
}

// GetEventAttendee handles GET /groups/events/:eventId/attendees/:userId
func (h *GroupHandler) GetEventAttendee(c *gin.Context) {
	// Get event ID and user ID from path
	eventID, err := strconv.ParseUint(c.Param("eventId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get event attendee
	attendee, err := h.groupService.GetAttendeeByID(uint(eventID), uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attendee not found"})
		return
	}

	c.JSON(http.StatusOK, attendee)
}

// UpdateAttendee handles PUT /groups/events/:eventId/attendees/:userId
func (h *GroupHandler) UpdateAttendee(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get event ID and attendee user ID from path
	eventID, err := strconv.ParseUint(c.Param("eventId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	attendeeUserID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if the user is updating their own attendance
	if userID.(uint) != uint(attendeeUserID) {
		// Check if the user is the event creator or a group admin/owner
		event, err := h.groupService.GetEventByID(uint(eventID))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
			return
		}

		if event.CreatedByID != userID.(uint) {
			member, err := h.groupService.GetMemberByID(event.GroupID, userID.(uint))
			if err != nil || (member.Role != models.OwnerRole && member.Role != models.AdminRole) {
				c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized: only the attendee, event creator, or group admins/owners can update attendance"})
				return
			}
		}
	}

	// Bind request body
	var request struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update attendee
	attendee, err := h.groupService.UpdateAttendee(uint(eventID), uint(attendeeUserID), request.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, attendee)
}

// RemoveAttendee handles DELETE /groups/events/:eventId/attendees/:userId
func (h *GroupHandler) RemoveAttendee(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get event ID and attendee user ID from path
	eventID, err := strconv.ParseUint(c.Param("eventId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	attendeeUserID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if the user is removing their own attendance
	if userID.(uint) != uint(attendeeUserID) {
		// Check if the user is the event creator or a group admin/owner
		event, err := h.groupService.GetEventByID(uint(eventID))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
			return
		}

		if event.CreatedByID != userID.(uint) {
			member, err := h.groupService.GetMemberByID(event.GroupID, userID.(uint))
			if err != nil || (member.Role != models.OwnerRole && member.Role != models.AdminRole) {
				c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized: only the attendee, event creator, or group admins/owners can remove attendance"})
				return
			}
		}
	}

	// Remove attendee
	err = h.groupService.RemoveAttendee(uint(eventID), uint(attendeeUserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendee removed successfully"})
}

// CheckInAttendee handles POST /groups/events/:eventId/attendees/:userId/checkin
func (h *GroupHandler) CheckInAttendee(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get event ID and attendee user ID from path
	eventID, err := strconv.ParseUint(c.Param("eventId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	attendeeUserID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check in attendee
	attendee, err := h.groupService.CheckInAttendee(uint(eventID), uint(attendeeUserID), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, attendee)
}
