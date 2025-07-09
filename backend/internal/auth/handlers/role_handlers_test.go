package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock user service for testing
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUserByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestGetEngagedUserFeatures(t *testing.T) {
	// Set up test cases
	testCases := []struct {
		name           string
		userID         uint
		userRole       int
		expectedStatus int
		mockErr        error
	}{
		{
			name:           "Success - Engaged User",
			userID:         1,
			userRole:       models.RoleEngagedUser,
			expectedStatus: http.StatusOK,
			mockErr:        nil,
		},
		{
			name:           "Success - Higher Role (Active)",
			userID:         2,
			userRole:       models.RoleActiveUser,
			expectedStatus: http.StatusOK,
			mockErr:        nil,
		},
		{
			name:           "Forbidden - Basic User",
			userID:         3,
			userRole:       models.RoleBasicUser,
			expectedStatus: http.StatusForbidden,
			mockErr:        nil,
		},
		{
			name:           "Error - Service Error",
			userID:         4,
			userRole:       models.RoleEngagedUser,
			expectedStatus: http.StatusInternalServerError,
			mockErr:        errors.New("database error"),
		},
	}

	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock user service
			mockService := new(MockUserService)
			
			// Set expectations
			if tc.mockErr == nil {
				user := &models.User{
					ID:   tc.userID,
					Role: tc.userRole,
				}
				mockService.On("GetUserByID", tc.userID).Return(user, nil)
			} else {
				mockService.On("GetUserByID", tc.userID).Return(nil, tc.mockErr)
			}
			
			// Create handler with mock service
			handler := NewRoleHandlers(mockService)
			
			// Create a response recorder
			w := httptest.NewRecorder()
			
			// Create a request context
			c, _ := gin.CreateTestContext(w)
			
			// Set user ID in context (mimicking middleware)
			c.Set("user_id", float64(tc.userID))
			
			// Call the handler function
			handler.GetEngagedUserFeatures(c)
			
			// Assert status code
			assert.Equal(t, tc.expectedStatus, w.Code)
			
			// Assert response body for success cases
			if tc.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				
				// Check that features are present
				features, exists := response["features"]
				assert.True(t, exists)
				assert.IsType(t, []interface{}{}, features)
				assert.NotEmpty(t, features)
				
				// Check role name
				role, exists := response["role"]
				assert.True(t, exists)
				assert.NotEmpty(t, role)
			}
			
			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetModeratorTools(t *testing.T) {
	// Set up test cases
	testCases := []struct {
		name           string
		userID         uint
		userRole       int
		expectedStatus int
		mockErr        error
	}{
		{
			name:           "Success - Moderator",
			userID:         1,
			userRole:       models.RoleModerator,
			expectedStatus: http.StatusOK,
			mockErr:        nil,
		},
		{
			name:           "Success - Higher Role (Admin)",
			userID:         2,
			userRole:       models.RoleAdmin,
			expectedStatus: http.StatusOK,
			mockErr:        nil,
		},
		{
			name:           "Forbidden - Premium User",
			userID:         3,
			userRole:       models.RolePremiumUser,
			expectedStatus: http.StatusForbidden,
			mockErr:        nil,
		},
		{
			name:           "Error - Service Error",
			userID:         4,
			userRole:       models.RoleModerator,
			expectedStatus: http.StatusInternalServerError,
			mockErr:        errors.New("database error"),
		},
	}

	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock user service
			mockService := new(MockUserService)
			
			// Set expectations
			if tc.mockErr == nil {
				user := &models.User{
					ID:   tc.userID,
					Role: tc.userRole,
				}
				mockService.On("GetUserByID", tc.userID).Return(user, nil)
			} else {
				mockService.On("GetUserByID", tc.userID).Return(nil, tc.mockErr)
			}
			
			// Create handler with mock service
			handler := NewRoleHandlers(mockService)
			
			// Create a response recorder
			w := httptest.NewRecorder()
			
			// Create a request context
			c, _ := gin.CreateTestContext(w)
			
			// Set user ID in context (mimicking middleware)
			c.Set("user_id", float64(tc.userID))
			
			// Call the handler function
			handler.GetModeratorTools(c)
			
			// Assert status code
			assert.Equal(t, tc.expectedStatus, w.Code)
			
			// Assert response body for success cases
			if tc.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				
				// Check that tools are present
				tools, exists := response["tools"]
				assert.True(t, exists)
				assert.IsType(t, []interface{}{}, tools)
				assert.NotEmpty(t, tools)
				
				// Check role name
				role, exists := response["role"]
				assert.True(t, exists)
				assert.NotEmpty(t, role)
			}
			
			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}