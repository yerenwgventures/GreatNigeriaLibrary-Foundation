package service

import (
	"errors"
	"testing"

	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/logger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository mocks the UserRepository interface
type MockUserRepositoryForDelete struct {
	mock.Mock
}

// DeleteUser mocks the DeleteUser method
func (m *MockUserRepositoryForDelete) DeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// VerifyPassword mocks the VerifyPassword method
func (m *MockUserRepositoryForDelete) VerifyPassword(id uint, password string) (bool, error) {
	args := m.Called(id, password)
	return args.Bool(0), args.Error(1)
}

// GetByID mocks the GetByID method required by the service
func (m *MockUserRepositoryForDelete) GetByID(id uint) (interface{}, error) {
	args := m.Called(id)
	return args.Get(0), args.Error(1)
}

// TestDeleteUser tests the DeleteUser method
func TestDeleteUser(t *testing.T) {
	// Create a test logger that doesn't output anything
	testLogger := &logger.Logger{Logger: logrus.New()}
	testLogger.Logger.SetOutput(nil)

	// Test cases
	testCases := []struct {
		name           string
		userID         uint
		password       string
		setupMock      func(*MockUserRepositoryForDelete)
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name:     "Successful deletion",
			userID:   1,
			password: "validPassword123",
			setupMock: func(m *MockUserRepositoryForDelete) {
				// Password verification succeeds
				m.On("GetByID", uint(1)).Return(struct{}{}, nil)
				m.On("VerifyPassword", uint(1), "validPassword123").Return(true, nil)
				
				// Deletion succeeds
				m.On("DeleteUser", uint(1)).Return(nil)
			},
			expectedError: false,
		},
		{
			name:     "Invalid password",
			userID:   2,
			password: "invalidPassword",
			setupMock: func(m *MockUserRepositoryForDelete) {
				// Password verification fails
				m.On("GetByID", uint(2)).Return(struct{}{}, nil)
				m.On("VerifyPassword", uint(2), "invalidPassword").Return(false, nil)
			},
			expectedError:  true,
			expectedErrMsg: "Invalid password",
		},
		{
			name:     "User not found",
			userID:   3,
			password: "password123",
			setupMock: func(m *MockUserRepositoryForDelete) {
				// User not found
				m.On("GetByID", uint(3)).Return(nil, errors.New("user not found"))
				m.On("VerifyPassword", uint(3), "password123").Return(false, errors.New("user not found"))
			},
			expectedError:  true,
			expectedErrMsg: "user not found",
		},
		{
			name:     "Database error during deletion",
			userID:   4,
			password: "validPassword123",
			setupMock: func(m *MockUserRepositoryForDelete) {
				// Password verification succeeds
				m.On("GetByID", uint(4)).Return(struct{}{}, nil)
				m.On("VerifyPassword", uint(4), "validPassword123").Return(true, nil)
				
				// Deletion fails
				m.On("DeleteUser", uint(4)).Return(errors.New("database error"))
			},
			expectedError:  true,
			expectedErrMsg: "database error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock repository
			mockRepo := new(MockUserRepositoryForDelete)
			
			// Setup mock expectations
			tc.setupMock(mockRepo)
			
			// Create a service with the mock repository
			service := &UserService{
				userRepo: mockRepo,
				logger:   testLogger,
			}
			
			// Call the method
			err := service.DeleteUser(tc.userID, tc.password)
			
			// Check expectations
			if tc.expectedError {
				assert.Error(t, err)
				if tc.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tc.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err)
			}
			
			// Verify that all expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}