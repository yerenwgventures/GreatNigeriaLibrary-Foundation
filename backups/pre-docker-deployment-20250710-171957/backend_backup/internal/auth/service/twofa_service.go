package service

import (
        "crypto/rand"
        "fmt"
        "math/big"
        "time"

        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/auth/repository"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/common/logger"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/models"
        "github.com/pquerna/otp/totp"
)

// TwoFAService defines methods for two-factor authentication service operations
type TwoFAService interface {
        SetupTwoFA(userID uint) (*models.TwoFactorSetupResponse, error)
        VerifyTwoFA(userID uint, token string) (bool, error)
        EnableTwoFA(userID uint, token string) error
        DisableTwoFA(userID uint, token string, password string) error
        ValidateToken(userID uint, token string) (bool, error)
        ValidateBackupCode(userID uint, backupCode string) (bool, error)
        GetTwoFAStatus(userID uint) (*models.TwoFactorAuthStatus, error)
        GenerateBackupCodes(userID uint) ([]string, error)
}

// TwoFAServiceImpl implements TwoFAService interface
type TwoFAServiceImpl struct {
        twoFARepo   repository.TwoFARepository
        userService UserService
        logger      *logger.Logger
}

// NewTwoFAService creates a new TwoFAService instance
func NewTwoFAService(twoFARepo repository.TwoFARepository, userService UserService, logger *logger.Logger) TwoFAService {
        return &TwoFAServiceImpl{
                twoFARepo:   twoFARepo,
                userService: userService,
                logger:      logger,
        }
}

// SetupTwoFA initializes 2FA for a user and returns the secret and QR code URL
func (s *TwoFAServiceImpl) SetupTwoFA(userID uint) (*models.TwoFactorSetupResponse, error) {
        // Get user details to include in the key
        user, err := s.userService.GetUserByID(userID)
        if err != nil {
                return nil, fmt.Errorf("failed to get user: %w", err)
        }
        
        // Generate a new TOTP key
        key, err := totp.Generate(totp.GenerateOpts{
                Issuer:      "Great Nigeria Platform",
                AccountName: user.Email,
        })
        if err != nil {
                return nil, fmt.Errorf("failed to generate TOTP key: %w", err)
        }
        
        // Store the secret in the database
        if err := s.twoFARepo.CreateTwoFA(userID, key.Secret()); err != nil {
                return nil, fmt.Errorf("failed to store 2FA secret: %w", err)
        }
        
        // Return the setup response
        return &models.TwoFactorSetupResponse{
                Secret:    key.Secret(),
                QRCodeURL: key.URL(),
        }, nil
}

// VerifyTwoFA verifies a 2FA token without enabling 2FA
func (s *TwoFAServiceImpl) VerifyTwoFA(userID uint, token string) (bool, error) {
        // Get 2FA settings
        twoFA, err := s.twoFARepo.GetTwoFA(userID)
        if err != nil {
                return false, fmt.Errorf("failed to get 2FA settings: %w", err)
        }
        
        // Verify the token
        valid := totp.Validate(token, twoFA.Secret)
        
        // If valid, mark as verified but don't enable yet
        if valid {
                if err := s.twoFARepo.UpdateTwoFA(userID, false, true); err != nil {
                        return false, fmt.Errorf("failed to update 2FA verification status: %w", err)
                }
        }
        
        return valid, nil
}

// EnableTwoFA enables 2FA for a user after verification
func (s *TwoFAServiceImpl) EnableTwoFA(userID uint, token string) error {
        // First verify the token
        valid, err := s.VerifyTwoFA(userID, token)
        if err != nil {
                return err
        }
        
        if !valid {
                return fmt.Errorf("invalid token")
        }
        
        // Enable 2FA
        if err := s.twoFARepo.UpdateTwoFA(userID, true, true); err != nil {
                return fmt.Errorf("failed to enable 2FA: %w", err)
        }
        
        // Generate backup codes
        _, err = s.GenerateBackupCodes(userID)
        if err != nil {
                return fmt.Errorf("failed to generate backup codes: %w", err)
        }
        
        return nil
}

// DisableTwoFA disables 2FA for a user
func (s *TwoFAServiceImpl) DisableTwoFA(userID uint, token string, password string) error {
        // Verify password
        user, err := s.userService.GetUserByID(userID)
        if err != nil {
                return fmt.Errorf("failed to get user: %w", err)
        }
        
        isValid, err := s.userService.VerifyPassword(userID, password)
        if err != nil || !isValid {
                return fmt.Errorf("invalid password")
        }
        
        // Verify the token
        twoFA, err := s.twoFARepo.GetTwoFA(userID)
        if err != nil {
                return fmt.Errorf("failed to get 2FA settings: %w", err)
        }
        
        valid := totp.Validate(token, twoFA.Secret)
        if !valid {
                return fmt.Errorf("invalid token")
        }
        
        // Disable 2FA
        if err := s.twoFARepo.DisableTwoFA(userID); err != nil {
                return fmt.Errorf("failed to disable 2FA: %w", err)
        }
        
        return nil
}

// ValidateToken validates a 2FA token
func (s *TwoFAServiceImpl) ValidateToken(userID uint, token string) (bool, error) {
        // Get 2FA settings
        twoFA, err := s.twoFARepo.GetTwoFA(userID)
        if err != nil {
                return false, fmt.Errorf("failed to get 2FA settings: %w", err)
        }
        
        // If 2FA is not enabled, return true (no validation needed)
        if !twoFA.Enabled {
                return true, nil
        }
        
        // Validate the token
        return totp.Validate(token, twoFA.Secret), nil
}

// ValidateBackupCode validates a backup code for 2FA
func (s *TwoFAServiceImpl) ValidateBackupCode(userID uint, backupCode string) (bool, error) {
        return s.twoFARepo.ValidateBackupCode(userID, backupCode)
}

// GetTwoFAStatus gets the current 2FA status for a user
func (s *TwoFAServiceImpl) GetTwoFAStatus(userID uint) (*models.TwoFactorAuthStatus, error) {
        twoFA, err := s.twoFARepo.GetTwoFA(userID)
        if err != nil {
                // If record not found, return default status (not enabled)
                return &models.TwoFactorAuthStatus{
                        Enabled:  false,
                        Verified: false,
                        Method:   "none",
                }, nil
        }
        
        method := "none"
        if twoFA.Enabled {
                method = "app" // TOTP app is currently the only method supported
        }
        
        return &models.TwoFactorAuthStatus{
                Enabled:  twoFA.Enabled,
                Verified: twoFA.Verified,
                Method:   method,
        }, nil
}

// GenerateBackupCodes generates backup codes for a user
func (s *TwoFAServiceImpl) GenerateBackupCodes(userID uint) ([]string, error) {
        backupCodes := make([]string, 8) // Generate 8 backup codes
        
        for i := 0; i < 8; i++ {
                // Generate a random 10-character backup code
                code, err := generateRandomString(10)
                if err != nil {
                        return nil, fmt.Errorf("failed to generate backup code: %w", err)
                }
                
                // Format as XXXX-XXXX-XX
                code = fmt.Sprintf("%s-%s-%s", code[0:4], code[4:8], code[8:10])
                backupCodes[i] = code
        }
        
        // Store the backup codes
        if err := s.twoFARepo.StoreBackupCodes(userID, backupCodes); err != nil {
                return nil, fmt.Errorf("failed to store backup codes: %w", err)
        }
        
        return backupCodes, nil
}

// Helper functions

// generateRandomString generates a random string of the specified length
func generateRandomString(length int) (string, error) {
        const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // Removed similar looking characters
        result := make([]byte, length)
        
        for i := 0; i < length; i++ {
                num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
                if err != nil {
                        return "", err
                }
                result[i] = charset[num.Int64()]
        }
        
        return string(result), nil
}