package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config represents Redis configuration
type Config struct {
	Host         string
	Port         int
	Password     string
	Database     int
	PoolSize     int
	MinIdleConns int
	MaxRetries   int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Client wraps Redis client with additional functionality
type Client struct {
	*redis.Client
	config *Config
}

// NewClient creates a new Redis client
func NewClient(config *Config) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:     config.Password,
		DB:           config.Database,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Client{
		Client: rdb,
		config: config,
	}, nil
}

// Close closes the Redis connection
func (c *Client) Close() error {
	return c.Client.Close()
}

// Health checks Redis health
func (c *Client) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return c.Client.Ping(ctx).Err()
}

// GetStats returns Redis connection statistics
func (c *Client) GetStats() map[string]interface{} {
	stats := c.Client.PoolStats()
	return map[string]interface{}{
		"hits":         stats.Hits,
		"misses":       stats.Misses,
		"timeouts":     stats.Timeouts,
		"total_conns":  stats.TotalConns,
		"idle_conns":   stats.IdleConns,
		"stale_conns":  stats.StaleConns,
	}
}

// SetWithExpiration sets a key with expiration
func (c *Client) SetWithExpiration(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.Client.Set(ctx, key, value, expiration).Err()
}

// GetString gets a string value
func (c *Client) GetString(ctx context.Context, key string) (string, error) {
	return c.Client.Get(ctx, key).Result()
}

// DeleteKeys deletes multiple keys
func (c *Client) DeleteKeys(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return c.Client.Del(ctx, keys...).Err()
}

// KeyExists checks if a key exists
func (c *Client) KeyExists(ctx context.Context, key string) (bool, error) {
	result, err := c.Client.Exists(ctx, key).Result()
	return result > 0, err
}

// SetNX sets a key only if it doesn't exist
func (c *Client) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return c.Client.SetNX(ctx, key, value, expiration).Result()
}

// Increment increments a key's value
func (c *Client) Increment(ctx context.Context, key string) (int64, error) {
	return c.Client.Incr(ctx, key).Result()
}

// IncrementWithExpiration increments a key and sets expiration
func (c *Client) IncrementWithExpiration(ctx context.Context, key string, expiration time.Duration) (int64, error) {
	pipe := c.Client.TxPipeline()
	incrCmd := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, expiration)
	
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	
	return incrCmd.Val(), nil
}

// GetKeysWithPattern gets all keys matching a pattern
func (c *Client) GetKeysWithPattern(ctx context.Context, pattern string) ([]string, error) {
	return c.Client.Keys(ctx, pattern).Result()
}

// DefaultConfig returns a default Redis configuration
func DefaultConfig() *Config {
	return &Config{
		Host:         "localhost",
		Port:         6379,
		Password:     "",
		Database:     0,
		PoolSize:     10,
		MinIdleConns: 5,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
}

// TestConnection tests Redis connectivity
func TestConnection(config *Config) error {
	client, err := NewClient(config)
	if err != nil {
		return err
	}
	defer client.Close()
	
	return client.Health()
}

// TokenStore provides Redis-based token storage operations
type TokenStore struct {
	client *Client
}

// NewTokenStore creates a new token store
func NewTokenStore(client *Client) *TokenStore {
	return &TokenStore{
		client: client,
	}
}

// StoreToken stores a token with metadata
func (ts *TokenStore) StoreToken(ctx context.Context, tokenHash string, metadata map[string]interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("token:%s", tokenHash)
	return ts.client.SetWithExpiration(ctx, key, metadata, expiration)
}

// GetToken retrieves token metadata
func (ts *TokenStore) GetToken(ctx context.Context, tokenHash string) (map[string]interface{}, error) {
	key := fmt.Sprintf("token:%s", tokenHash)
	result, err := ts.client.Client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	
	metadata := make(map[string]interface{})
	for k, v := range result {
		metadata[k] = v
	}
	return metadata, nil
}

// RevokeToken marks a token as revoked
func (ts *TokenStore) RevokeToken(ctx context.Context, tokenHash string, expiration time.Duration) error {
	key := fmt.Sprintf("revoked_token:%s", tokenHash)
	return ts.client.SetWithExpiration(ctx, key, "revoked", expiration)
}

// IsTokenRevoked checks if a token is revoked
func (ts *TokenStore) IsTokenRevoked(ctx context.Context, tokenHash string) (bool, error) {
	key := fmt.Sprintf("revoked_token:%s", tokenHash)
	return ts.client.KeyExists(ctx, key)
}

// IncrementUserTokenVersion increments user's token version
func (ts *TokenStore) IncrementUserTokenVersion(ctx context.Context, userID uint) (int64, error) {
	key := fmt.Sprintf("user_token_version:%d", userID)
	return ts.client.IncrementWithExpiration(ctx, key, 24*time.Hour)
}

// GetUserTokenVersion gets user's current token version
func (ts *TokenStore) GetUserTokenVersion(ctx context.Context, userID uint) (int64, error) {
	key := fmt.Sprintf("user_token_version:%d", userID)
	result, err := ts.client.Client.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil // Default version
	}
	return result, err
}

// StoreSession stores session information
func (ts *TokenStore) StoreSession(ctx context.Context, sessionID string, sessionData map[string]interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return ts.client.Client.HMSet(ctx, key, sessionData).Err()
}

// GetSession retrieves session information
func (ts *TokenStore) GetSession(ctx context.Context, sessionID string) (map[string]string, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	return ts.client.Client.HGetAll(ctx, key).Result()
}

// DeleteSession deletes a session
func (ts *TokenStore) DeleteSession(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return ts.client.Client.Del(ctx, key).Err()
}

// GetUserSessions gets all sessions for a user
func (ts *TokenStore) GetUserSessions(ctx context.Context, userID uint) ([]string, error) {
	pattern := fmt.Sprintf("session:*:user:%d", userID)
	return ts.client.GetKeysWithPattern(ctx, pattern)
}

// RevokeAllUserSessions revokes all sessions for a user
func (ts *TokenStore) RevokeAllUserSessions(ctx context.Context, userID uint) error {
	sessions, err := ts.GetUserSessions(ctx, userID)
	if err != nil {
		return err
	}
	
	if len(sessions) > 0 {
		return ts.client.DeleteKeys(ctx, sessions...)
	}
	
	return nil
}
