package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTManager manages JWT tokens
type JWTManager struct {
	secretKey              string
	accessTokenExpiration  time.Duration
	refreshTokenExpiration time.Duration
}

// Claims represents JWT claims
type Claims struct {
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      int    `json:"role"`
	SessionID string `json:"session_id,omitempty"`
	jwt.RegisteredClaims
}

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(secretKey string, accessExpiration, refreshExpiration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:              secretKey,
		accessTokenExpiration:  accessExpiration,
		refreshTokenExpiration: refreshExpiration,
	}
}

// GenerateTokens generates access and refresh tokens
func (j *JWTManager) GenerateTokens(userID uint, username, email string, role int, sessionID string) (*TokenPair, error) {
	// Generate access token
	accessToken, err := j.generateToken(userID, username, email, role, sessionID, j.accessTokenExpiration)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := j.generateToken(userID, username, email, role, sessionID, j.refreshTokenExpiration)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(j.accessTokenExpiration.Seconds()),
	}, nil
}

// generateToken generates a JWT token
func (j *JWTManager) generateToken(userID uint, username, email string, role int, sessionID string, expiration time.Duration) (string, error) {
	claims := Claims{
		UserID:    userID,
		Username:  username,
		Email:     email,
		Role:      role,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "great-nigeria-library",
			Subject:   username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// ValidateToken validates a JWT token and returns claims
func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken refreshes an access token using a refresh token
func (j *JWTManager) RefreshToken(refreshTokenString string) (*TokenPair, error) {
	claims, err := j.ValidateToken(refreshTokenString)
	if err != nil {
		return nil, err
	}

	// Generate new tokens
	return j.GenerateTokens(claims.UserID, claims.Username, claims.Email, claims.Role, claims.SessionID)
}

// ExtractUserID extracts user ID from token
func (j *JWTManager) ExtractUserID(tokenString string) (uint, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// ExtractClaims extracts all claims from token
func (j *JWTManager) ExtractClaims(tokenString string) (*Claims, error) {
	return j.ValidateToken(tokenString)
}

// IsTokenExpired checks if a token is expired
func (j *JWTManager) IsTokenExpired(tokenString string) bool {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return true
	}
	return claims.ExpiresAt.Time.Before(time.Now())
}

// GetTokenExpiration returns token expiration time
func (j *JWTManager) GetTokenExpiration(tokenString string) (time.Time, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return time.Time{}, err
	}
	return claims.ExpiresAt.Time, nil
}

// RevokeToken adds token to revocation list (would need Redis implementation)
func (j *JWTManager) RevokeToken(tokenString string) error {
	// This would typically store the token in Redis with expiration
	// For now, we'll just validate the token exists
	_, err := j.ValidateToken(tokenString)
	return err
}

// IsTokenRevoked checks if token is revoked (would need Redis implementation)
func (j *JWTManager) IsTokenRevoked(tokenString string) bool {
	// This would typically check Redis for the token
	// For now, we'll just return false
	return false
}
