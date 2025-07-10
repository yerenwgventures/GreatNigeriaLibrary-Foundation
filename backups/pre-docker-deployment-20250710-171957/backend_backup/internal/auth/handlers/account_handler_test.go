package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/logger"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/models"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of the UserService interface
type MockUserService struct {
	mock.Mock
}

// Implementation of UserService methods needed for DeleteAccount handler
func (m *MockUserService) DeleteUser(userID uint, password string) error {
	args := m.Called(userID, password)
	return args.Error(0)
}

func TestDeleteAccount(t *testing.T) {
	// Set up Gin in test mode
	gin.SetMode(gin.TestMode)

	// Create a logger for testing
	loggerInst := &logger.Logger{Logger: logrus.New()}
	loggerInst.Logger.SetOutput(bytes.NewBuffer(nil)) // Suppress log output

	// Create mock service
	mockService := new(MockUserService)

	// Create the handler with the mock service
	accountHandler := NewAccountHandler(mockService, loggerInst)

	// Test cases
	testCases := []struct {
		name           string
		userID         uint
		requestBody    models.AccountDeleteRequest
		setupMock      func()
		expectedStatus int
	}{
		{
			name:   "Successful deletion",
			userID: 123,
			requestBody: models.AccountDeleteRequest{
				Password: "password123",
			},
			setupMock: func() {
				mockService.On("DeleteUser", uint(123), "password123").Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Invalid password",
			userID: 123,
			requestBody: models.AccountDeleteRequest{
				Password: "wrongpassword",
			},
			setupMock: func() {
				mockService.On("DeleteUser", uint(123), "wrongpassword").Return(errors.New("invalid password"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock expectations
			tc.setupMock()

			// Create a new HTTP recorder and request
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("DELETE", "/account/delete", bytes.NewBuffer(reqBody))

			// Create a new Gin context with the request and response recorder
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Set the userID in the context (as would be done by middleware)
			c.Set("userID", tc.userID)

			// Call the DeleteAccount handler
			accountHandler.DeleteAccount(c)

			// Check response status
			assert.Equal(t, tc.expectedStatus, w.Code)

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}