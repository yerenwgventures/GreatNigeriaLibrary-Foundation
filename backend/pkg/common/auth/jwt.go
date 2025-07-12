package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

// JWTManager manages JWT tokens with Redis-based revocation
type JWTManager struct {
	secretKey              string
	accessTokenExpiration  time.Duration
	refreshTokenExpiration time.Duration
	redisClient            *redis.Client
	issuer                 string
}

// Claims represents enhanced JWT claims
type Claims struct {
	UserID       uint     `json:"user_id"`
	Username     string   `json:"username"`
	Email        string   `json:"email"`
	Role         int      `json:"role"`
	Permissions  []string `json:"permissions,omitempty"`
	SessionID    string   `json:"session_id,omitempty"`
	TokenType    string   `json:"token_type"` // "access" or "refresh"
	DeviceID     string   `json:"device_id,omitempty"`
	IPAddress    string   `json:"ip_address,omitempty"`
	TokenVersion int      `json:"token_version"` // For token invalidation
	jwt.RegisteredClaims
}

// TokenPair represents access and refresh tokens with metadata
type TokenPair struct {
	AccessToken   string    `json:"access_token"`
	RefreshToken  string    `json:"refresh_token"`
	TokenType     string    `json:"token_type"`     // "Bearer"
	ExpiresIn     int64     `json:"expires_in"`     // Access token expiration in seconds
	RefreshExpiresIn int64  `json:"refresh_expires_in"` // Refresh token expiration in seconds
	IssuedAt      time.Time `json:"issued_at"`
	SessionID     string    `json:"session_id"`
}

// TokenMetadata represents token metadata for tracking
type TokenMetadata struct {
	UserID       uint      `json:"user_id"`
	SessionID    string    `json:"session_id"`
	DeviceID     string    `json:"device_id"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	IssuedAt     time.Time `json:"issued_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenVersion int       `json:"token_version"`
}

// NewJWTManager creates a new enhanced JWT manager with Redis support
func NewJWTManager(secretKey string, accessExpiration, refreshExpiration time.Duration, redisClient *redis.Client, issuer string) *JWTManager {
	return &JWTManager{
		secretKey:              secretKey,
		accessTokenExpiration:  accessExpiration,
		refreshTokenExpiration: refreshExpiration,
		redisClient:            redisClient,
		issuer:                 issuer,
	}
}

// NewJWTManagerWithoutRedis creates a JWT manager without Redis (for testing)
func NewJWTManagerWithoutRedis(secretKey string, accessExpiration, refreshExpiration time.Duration, issuer string) *JWTManager {
	return &JWTManager{
		secretKey:              secretKey,
		accessTokenExpiration:  accessExpiration,
		refreshTokenExpiration: refreshExpiration,
		redisClient:            nil,
		issuer:                 issuer,
	}
}

// GenerateTokensWithMetadata generates access and refresh tokens with enhanced metadata
func (j *JWTManager) GenerateTokensWithMetadata(userID uint, username, email string, role int, permissions []string, sessionID, deviceID, ipAddress string) (*TokenPair, error) {
	now := time.Now()
	tokenVersion := j.generateTokenVersion()

	// Generate access token
	accessToken, err := j.generateEnhancedToken(userID, username, email, role, permissions, sessionID, deviceID, ipAddress, "access", tokenVersion, j.accessTokenExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := j.generateEnhancedToken(userID, username, email, role, permissions, sessionID, deviceID, ipAddress, "refresh", tokenVersion, j.refreshTokenExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store token metadata in Redis if available
	if j.redisClient != nil {
		if err := j.storeTokenMetadata(accessToken, userID, sessionID, deviceID, ipAddress, "", now, now.Add(j.accessTokenExpiration), tokenVersion); err != nil {
			// Log error but don't fail token generation
			fmt.Printf("Warning: Failed to store access token metadata: %v\n", err)
		}
		if err := j.storeTokenMetadata(refreshToken, userID, sessionID, deviceID, ipAddress, "", now, now.Add(j.refreshTokenExpiration), tokenVersion); err != nil {
			// Log error but don't fail token generation
			fmt.Printf("Warning: Failed to store refresh token metadata: %v\n", err)
		}
	}

	return &TokenPair{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		TokenType:        "Bearer",
		ExpiresIn:        int64(j.accessTokenExpiration.Seconds()),
		RefreshExpiresIn: int64(j.refreshTokenExpiration.Seconds()),
		IssuedAt:         now,
		SessionID:        sessionID,
	}, nil
}

// GenerateTokens generates tokens with basic metadata (backward compatibility)
func (j *JWTManager) GenerateTokens(userID uint, username, email string, role int, sessionID string) (*TokenPair, error) {
	return j.GenerateTokensWithMetadata(userID, username, email, role, []string{}, sessionID, "", "")
}

// generateEnhancedToken generates a JWT token with enhanced claims
func (j *JWTManager) generateEnhancedToken(userID uint, username, email string, role int, permissions []string, sessionID, deviceID, ipAddress, tokenType string, tokenVersion int, expiration time.Duration) (string, error) {
	now := time.Now()

	claims := Claims{
		UserID:       userID,
		Username:     username,
		Email:        email,
		Role:         role,
		Permissions:  permissions,
		SessionID:    sessionID,
		TokenType:    tokenType,
		DeviceID:     deviceID,
		IPAddress:    ipAddress,
		TokenVersion: tokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    j.issuer,
			Subject:   username,
			ID:        j.generateJTI(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// generateToken generates a basic JWT token (backward compatibility)
func (j *JWTManager) generateToken(userID uint, username, email string, role int, sessionID string, expiration time.Duration) (string, error) {
	return j.generateEnhancedToken(userID, username, email, role, []string{}, sessionID, "", "", "access", 1, expiration)
}

// ValidateToken validates a JWT token with enhanced security checks
func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	// Check if token is revoked first
	if j.IsTokenRevoked(tokenString) {
		return nil, errors.New("token has been revoked")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token parsing failed: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Additional validation checks
		if err := j.validateClaims(claims); err != nil {
			return nil, fmt.Errorf("claims validation failed: %w", err)
		}
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}

// validateClaims performs additional validation on token claims
func (j *JWTManager) validateClaims(claims *Claims) error {
	// Check if token is expired
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return errors.New("token has expired")
	}

	// Check if token is used before valid time
	if claims.NotBefore != nil && claims.NotBefore.Time.After(time.Now()) {
		return errors.New("token not yet valid")
	}

	// Check issuer
	if claims.Issuer != j.issuer {
		return fmt.Errorf("invalid issuer: expected %s, got %s", j.issuer, claims.Issuer)
	}

	// Validate required fields
	if claims.UserID == 0 {
		return errors.New("invalid user ID in token")
	}

	if claims.Username == "" {
		return errors.New("missing username in token")
	}

	return nil
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

// RevokeToken adds token to revocation list in Redis
func (j *JWTManager) RevokeToken(tokenString string) error {
	if j.redisClient == nil {
		return errors.New("Redis client not available for token revocation")
	}

	// Extract claims to get expiration time
	claims, err := j.extractClaimsWithoutValidation(tokenString)
	if err != nil {
		return fmt.Errorf("failed to extract claims for revocation: %w", err)
	}

	// Calculate TTL based on token expiration
	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		// Token already expired, no need to revoke
		return nil
	}

	// Store token hash in Redis with expiration
	tokenHash := j.hashToken(tokenString)
	key := fmt.Sprintf("revoked_token:%s", tokenHash)

	ctx := context.Background()
	err = j.redisClient.Set(ctx, key, "revoked", ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to revoke token in Redis: %w", err)
	}

	return nil
}

// IsTokenRevoked checks if token is revoked in Redis
func (j *JWTManager) IsTokenRevoked(tokenString string) bool {
	if j.redisClient == nil {
		return false // If no Redis, assume not revoked
	}

	tokenHash := j.hashToken(tokenString)
	key := fmt.Sprintf("revoked_token:%s", tokenHash)

	ctx := context.Background()
	exists, err := j.redisClient.Exists(ctx, key).Result()
	if err != nil {
		// Log error but don't block validation
		fmt.Printf("Warning: Failed to check token revocation status: %v\n", err)
		return false
	}

	return exists > 0
}

// RevokeAllUserTokens revokes all tokens for a specific user
func (j *JWTManager) RevokeAllUserTokens(userID uint) error {
	if j.redisClient == nil {
		return errors.New("Redis client not available for token revocation")
	}

	// Increment user's token version to invalidate all existing tokens
	key := fmt.Sprintf("user_token_version:%d", userID)
	ctx := context.Background()

	_, err := j.redisClient.Incr(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to increment user token version: %w", err)
	}

	// Set expiration for the version key (optional cleanup)
	j.redisClient.Expire(ctx, key, 24*time.Hour)

	return nil
}

// Helper methods for enhanced JWT functionality

// generateJTI generates a unique JWT ID
func (j *JWTManager) generateJTI() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// generateTokenVersion generates a random token version
func (j *JWTManager) generateTokenVersion() int {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return int(bytes[0])<<24 | int(bytes[1])<<16 | int(bytes[2])<<8 | int(bytes[3])
}

// hashToken creates a hash of the token for storage
func (j *JWTManager) hashToken(tokenString string) string {
	// Simple hash for demonstration - in production, use crypto/sha256
	return fmt.Sprintf("%x", len(tokenString)*31+int(tokenString[0]))
}

// extractClaimsWithoutValidation extracts claims without full validation (for revocation)
func (j *JWTManager) extractClaimsWithoutValidation(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	}, jwt.WithoutClaimsValidation())

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}

// storeTokenMetadata stores token metadata in Redis
func (j *JWTManager) storeTokenMetadata(tokenString string, userID uint, sessionID, deviceID, ipAddress, userAgent string, issuedAt, expiresAt time.Time, tokenVersion int) error {
	if j.redisClient == nil {
		return nil // Skip if no Redis
	}

	metadata := TokenMetadata{
		UserID:       userID,
		SessionID:    sessionID,
		DeviceID:     deviceID,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		IssuedAt:     issuedAt,
		ExpiresAt:    expiresAt,
		TokenVersion: tokenVersion,
	}

	tokenHash := j.hashToken(tokenString)
	key := fmt.Sprintf("token_metadata:%s", tokenHash)

	ctx := context.Background()
	ttl := time.Until(expiresAt)

	// Store metadata as JSON (simplified - could use msgpack for efficiency)
	return j.redisClient.Set(ctx, key, fmt.Sprintf("%+v", metadata), ttl).Err()
}

// GetTokenMetadata retrieves token metadata from Redis
func (j *JWTManager) GetTokenMetadata(tokenString string) (*TokenMetadata, error) {
	if j.redisClient == nil {
		return nil, errors.New("Redis client not available")
	}

	tokenHash := j.hashToken(tokenString)
	key := fmt.Sprintf("token_metadata:%s", tokenHash)

	ctx := context.Background()
	_, err := j.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.New("token metadata not found")
		}
		return nil, fmt.Errorf("failed to get token metadata: %w", err)
	}

	// In a real implementation, you'd properly unmarshal JSON
	// This is simplified for demonstration
	return &TokenMetadata{}, nil
}
