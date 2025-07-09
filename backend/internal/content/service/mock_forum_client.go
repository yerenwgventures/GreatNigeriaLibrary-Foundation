package service

import (
	"errors"
	"fmt"
	"sync"
)

// MockForumClient implements the ForumClient interface for testing and development
type MockForumClient struct {
	topics map[string]map[string]interface{}
	mu     sync.RWMutex
}

// NewMockForumClient creates a new mock forum client with some predefined topics
func NewMockForumClient() *MockForumClient {
	client := &MockForumClient{
		topics: make(map[string]map[string]interface{}),
	}
	
	// Add some mock topics
	client.AddTopic("1", "Introduction to Nigerian Economy", "Discussion about the Nigerian economy chapter")
	client.AddTopic("2", "Agricultural Opportunities in Nigeria", "Discussing agricultural potential")
	client.AddTopic("3", "Tech Innovation Hubs in Nigeria", "Nigeria's growing tech scene")
	
	return client
}

// AddTopic adds a mock topic to the client
func (c *MockForumClient) AddTopic(id, title, description string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.topics[id] = map[string]interface{}{
		"id":          id,
		"title":       title,
		"description": description,
		"url":         fmt.Sprintf("/forum/topic/%s", id),
	}
}

// GetTopicByID retrieves a topic by its ID
func (c *MockForumClient) GetTopicByID(id string) (map[string]interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	topic, exists := c.topics[id]
	if !exists {
		return nil, errors.New("topic not found")
	}
	
	return topic, nil
}