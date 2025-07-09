# API Reference
## Great Nigeria Library Platform

**Document Version**: 1.0  
**Last Updated**: January 2025  
**API Version**: v1  
**Base URL**: `https://api.greatnigeria.net/api/v1`

---

## Table of Contents

1. [Overview](#overview)
2. [Authentication](#authentication)
3. [Error Handling](#error-handling)
4. [Pagination](#pagination)
5. [Rate Limiting](#rate-limiting)
6. [Auth Service API](#auth-service-api)
7. [Content Service API](#content-service-api)
8. [User Service API](#user-service-api)
9. [Discussion Service API](#discussion-service-api)
10. [Points Service API](#points-service-api)
11. [Payment Service API](#payment-service-api)
12. [Celebration Service API](#celebration-service-api)
13. [Livestream Service API](#livestream-service-api)
14. [Search Service API](#search-service-api)
15. [Notification Service API](#notification-service-api)

---

## Overview

The Great Nigeria Library API is a RESTful API that provides access to all platform functionality. The API is organized around resources and uses standard HTTP verbs, response codes, and authentication patterns.

### API Principles

- **RESTful Design**: Resources are accessed via standard HTTP verbs
- **JSON Format**: All requests and responses use JSON
- **Stateless**: Each request contains all necessary information
- **Versioned**: API versions are maintained for backward compatibility
- **Secure**: All endpoints require proper authentication and authorization

### Base Response Format

All API responses follow this consistent format:

```json
{
  "status": "success" | "error",
  "data": any | null,
  "error": {
    "code": "string",
    "message": "string",
    "details": object
  } | null,
  "meta": {
    "timestamp": "2025-01-15T10:30:00Z",
    "request_id": "req_123456",
    "version": "v1"
  }
}
```

---

## Authentication

### Bearer Token Authentication

All protected endpoints require a JWT bearer token in the Authorization header:

```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Token Structure

JWT tokens contain the following claims:

```json
{
  "sub": "user_id",
  "email": "user@example.com",
  "role": "user|admin|creator",
  "permissions": ["read:books", "write:comments"],
  "exp": 1735689600,
  "iat": 1735603200
}
```

### Authentication Flow

1. **Login**: `POST /auth/login` - Get access and refresh tokens
2. **Refresh**: `POST /auth/refresh` - Get new access token
3. **Logout**: `POST /auth/logout` - Invalidate tokens

---

## Error Handling

### HTTP Status Codes

| Code | Description |
|------|-------------|
| `200` | OK - Request succeeded |
| `201` | Created - Resource created successfully |
| `400` | Bad Request - Invalid request format or parameters |
| `401` | Unauthorized - Authentication required or invalid |
| `403` | Forbidden - Access denied for authenticated user |
| `404` | Not Found - Resource doesn't exist |
| `422` | Unprocessable Entity - Validation failed |
| `429` | Too Many Requests - Rate limit exceeded |
| `500` | Internal Server Error - Server error |

### Error Response Format

```json
{
  "status": "error",
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Request validation failed",
    "details": {
      "email": "Valid email address is required",
      "password": "Password must be at least 8 characters"
    }
  },
  "meta": {
    "timestamp": "2025-01-15T10:30:00Z",
    "request_id": "req_123456"
  }
}
```

### Common Error Codes

| Code | Description |
|------|-------------|
| `VALIDATION_ERROR` | Request data validation failed |
| `AUTHENTICATION_REQUIRED` | User must be authenticated |
| `PERMISSION_DENIED` | User lacks required permissions |
| `RESOURCE_NOT_FOUND` | Requested resource doesn't exist |
| `RATE_LIMIT_EXCEEDED` | Too many requests from client |
| `SERVER_ERROR` | Internal server error occurred |

---

## Pagination

### Request Parameters

```http
GET /books?page=1&limit=20&sort=created_at&order=desc
```

| Parameter | Description | Default | Maximum |
|-----------|-------------|---------|---------|
| `page` | Page number (1-indexed) | 1 | 1000 |
| `limit` | Items per page | 20 | 100 |
| `sort` | Sort field | `created_at` | - |
| `order` | Sort order (`asc`, `desc`) | `desc` | - |

### Response Format

```json
{
  "status": "success",
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "pages": 8,
    "has_next": true,
    "has_prev": false,
    "next_page": 2,
    "prev_page": null
  }
}
```

---

## Rate Limiting

### Rate Limits by User Tier

| Tier | Requests per Minute | Burst Limit |
|------|-------------------|-------------|
| Anonymous | 60 | 100 |
| Basic | 120 | 200 |
| Premium | 300 | 500 |
| Creator | 500 | 1000 |
| Admin | 1000 | 2000 |

### Rate Limit Headers

```http
X-RateLimit-Limit: 120
X-RateLimit-Remaining: 115
X-RateLimit-Reset: 1640995200
X-RateLimit-Retry-After: 60
```

---

## Auth Service API

### Register User

Create a new user account.

**Endpoint**: `POST /auth/register`

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "securepassword123",
  "name": "John Doe",
  "role": "user"
}
```

**Response** (201):
```json
{
  "status": "success",
  "data": {
    "user": {
      "id": "usr_123456",
      "email": "user@example.com",
      "name": "John Doe",
      "role": "user",
      "verified": false,
      "created_at": "2025-01-15T10:30:00Z"
    },
    "tokens": {
      "access_token": "eyJhbGciOiJIUzI1NiIs...",
      "refresh_token": "refresh_token_here",
      "expires_in": 3600
    }
  }
}
```

### Login User

Authenticate user and get tokens.

**Endpoint**: `POST /auth/login`

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response** (200):
```json
{
  "status": "success",
  "data": {
    "user": {
      "id": "usr_123456",
      "email": "user@example.com",
      "name": "John Doe",
      "role": "user",
      "verified": true,
      "last_login": "2025-01-15T10:30:00Z"
    },
    "tokens": {
      "access_token": "eyJhbGciOiJIUzI1NiIs...",
      "refresh_token": "refresh_token_here",
      "expires_in": 3600
    }
  }
}
```

### Refresh Token

Get a new access token using refresh token.

**Endpoint**: `POST /auth/refresh`

**Request Body**:
```json
{
  "refresh_token": "refresh_token_here"
}
```

**Response** (200):
```json
{
  "status": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 3600
  }
}
```

### Logout

Invalidate user tokens.

**Endpoint**: `POST /auth/logout`
**Authentication**: Required

**Response** (200):
```json
{
  "status": "success",
  "data": {
    "message": "Successfully logged out"
  }
}
```

### Get Current User

Get current authenticated user information.

**Endpoint**: `GET /auth/me`
**Authentication**: Required

**Response** (200):
```json
{
  "status": "success",
  "data": {
    "id": "usr_123456",
    "email": "user@example.com",
    "name": "John Doe",
    "role": "user",
    "verified": true,
    "preferences": {
      "theme": "light",
      "language": "en",
      "notifications": true
    },
    "created_at": "2025-01-15T10:30:00Z",
    "updated_at": "2025-01-15T10:30:00Z"
  }
}
```

---

## Content Service API

### List Books

Get a paginated list of books.

**Endpoint**: `GET /books`

**Query Parameters**:
- `page` (integer): Page number
- `limit` (integer): Items per page
- `category` (string): Filter by category
- `difficulty` (string): Filter by difficulty level
- `search` (string): Search query

**Response** (200):
```json
{
  "status": "success",
  "data": [
    {
      "id": "book_123456",
      "title": "Introduction to Computer Science",
      "author": "Dr. Adebayo Ogundimu",
      "description": "A comprehensive guide to computer science fundamentals",
      "cover_image": "https://cdn.greatnigeria.net/books/cs-intro-cover.jpg",
      "category": "Technology",
      "difficulty": "beginner",
      "chapters_count": 12,
      "estimated_duration": "15 hours",
      "price": 0,
      "is_free": true,
      "created_at": "2025-01-15T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 45,
    "pages": 3,
    "has_next": true,
    "has_prev": false
  }
}
```

### Get Book Details

Get detailed information about a specific book.

**Endpoint**: `GET /books/{book_id}`

**Response** (200):
```json
{
  "status": "success",
  "data": {
    "id": "book_123456",
    "title": "Introduction to Computer Science",
    "author": "Dr. Adebayo Ogundimu",
    "description": "A comprehensive guide to computer science fundamentals",
    "full_description": "This book covers all fundamental concepts...",
    "cover_image": "https://cdn.greatnigeria.net/books/cs-intro-cover.jpg",
    "category": "Technology",
    "difficulty": "beginner",
    "chapters": [
      {
        "id": "chapter_001",
        "title": "Introduction to Programming",
        "order": 1,
        "estimated_duration": "1 hour",
        "is_locked": false
      }
    ],
    "metadata": {
      "isbn": "978-3-16-148410-0",
      "publisher": "Nigerian Academic Press",
      "publication_date": "2024-01-15",
      "language": "en",
      "pages": 340
    },
    "stats": {
      "views": 1250,
      "completions": 89,
      "average_rating": 4.7,
      "total_ratings": 156
    },
    "price": 0,
    "is_free": true,
    "created_at": "2025-01-15T10:30:00Z",
    "updated_at": "2025-01-15T10:30:00Z"
  }
}
```

### Get Chapter Content

Get the content of a specific chapter.

**Endpoint**: `GET /books/{book_id}/chapters/{chapter_id}`
**Authentication**: Required

**Response** (200):
```json
{
  "status": "success",
  "data": {
    "id": "chapter_001",
    "title": "Introduction to Programming",
    "content": "# Introduction to Programming\n\nProgramming is...",
    "content_type": "markdown",
    "order": 1,
    "estimated_duration": "1 hour",
    "sections": [
      {
        "id": "section_001",
        "title": "What is Programming?",
        "content": "Programming is the process...",
        "order": 1
      }
    ],
    "quiz": {
      "id": "quiz_001",
      "questions_count": 5,
      "time_limit": 600
    },
    "navigation": {
      "previous_chapter": null,
      "next_chapter": {
        "id": "chapter_002",
        "title": "Variables and Data Types"
      }
    }
  }
}
```

### Track Reading Progress

Update user's reading progress for a chapter.

**Endpoint**: `POST /books/{book_id}/chapters/{chapter_id}/progress`
**Authentication**: Required

**Request Body**:
```json
{
  "progress_percentage": 75,
  "time_spent": 45,
  "completed": false
}
```

**Response** (200):
```json
{
  "status": "success",
  "data": {
    "user_id": "usr_123456",
    "book_id": "book_123456",
    "chapter_id": "chapter_001",
    "progress_percentage": 75,
    "time_spent": 45,
    "completed": false,
    "last_read": "2025-01-15T10:30:00Z"
  }
}
```

### Create Note

Create a note for a specific section.

**Endpoint**: `POST /books/{book_id}/chapters/{chapter_id}/notes`
**Authentication**: Required

**Request Body**:
```json
{
  "content": "This is an important concept to remember",
  "section_id": "section_001",
  "highlight_text": "Programming is the process",
  "position": {
    "start": 125,
    "end": 158
  }
}
```

**Response** (201):
```json
{
  "status": "success",
  "data": {
    "id": "note_123456",
    "content": "This is an important concept to remember",
    "section_id": "section_001",
    "highlight_text": "Programming is the process",
    "position": {
      "start": 125,
      "end": 158
    },
    "created_at": "2025-01-15T10:30:00Z"
  }
}
```

---

## User Service API

### Get User Profile

Get user profile information.

**Endpoint**: `GET /users/{user_id}`

**Response** (200):
```json
{
  "status": "success",
  "data": {
    "id": "usr_123456",
    "username": "johndoe",
    "name": "John Doe",
    "bio": "Software engineer and lifelong learner",
    "avatar": "https://cdn.greatnigeria.net/avatars/usr_123456.jpg",
    "location": "Lagos, Nigeria",
    "website": "https://johndoe.dev",
    "social_links": {
      "twitter": "@johndoe",
      "linkedin": "johndoe"
    },
    "stats": {
      "books_completed": 15,
      "total_reading_time": 45.5,
      "points": 2750,
      "badges_earned": 8,
      "streak_days": 12
    },
    "preferences": {
      "is_public": true,
      "show_reading_stats": true,
      "show_achievements": true
    },
    "joined_at": "2024-08-15T10:30:00Z"
  }
}
```

### Update User Profile

Update current user's profile.

**Endpoint**: `PUT /users/me`
**Authentication**: Required

**Request Body**:
```json
{
  "name": "John Doe",
  "bio": "Software engineer and lifelong learner",
  "location": "Lagos, Nigeria",
  "website": "https://johndoe.dev",
  "social_links": {
    "twitter": "@johndoe",
    "linkedin": "johndoe"
  }
}
```

**Response** (200):
```json
{
  "status": "success",
  "data": {
    "id": "usr_123456",
    "name": "John Doe",
    "bio": "Software engineer and lifelong learner",
    "location": "Lagos, Nigeria",
    "website": "https://johndoe.dev",
    "updated_at": "2025-01-15T10:30:00Z"
  }
}
```

---

## Discussion Service API

### List Discussions

Get a list of discussion topics.

**Endpoint**: `GET /discussions`

**Query Parameters**:
- `category` (string): Filter by category
- `sort` (string): Sort by (`recent`, `popular`, `oldest`)
- `tag` (string): Filter by tag

**Response** (200):
```json
{
  "status": "success",
  "data": [
    {
      "id": "disc_123456",
      "title": "Best practices for learning programming",
      "content": "What are your favorite techniques...",
      "author": {
        "id": "usr_123456",
        "name": "John Doe",
        "avatar": "https://cdn.greatnigeria.net/avatars/usr_123456.jpg"
      },
      "category": "Technology",
      "tags": ["programming", "learning", "tips"],
      "stats": {
        "views": 245,
        "replies": 18,
        "likes": 32,
        "bookmarks": 8
      },
      "last_reply": {
        "author": "Jane Smith",
        "timestamp": "2025-01-14T15:45:00Z"
      },
      "created_at": "2025-01-10T08:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 156,
    "pages": 8
  }
}
```

### Create Discussion

Create a new discussion topic.

**Endpoint**: `POST /discussions`
**Authentication**: Required

**Request Body**:
```json
{
  "title": "Best practices for learning programming",
  "content": "What are your favorite techniques for learning new programming languages?",
  "category": "Technology",
  "tags": ["programming", "learning", "tips"]
}
```

**Response** (201):
```json
{
  "status": "success",
  "data": {
    "id": "disc_123456",
    "title": "Best practices for learning programming",
    "content": "What are your favorite techniques...",
    "category": "Technology",
    "tags": ["programming", "learning", "tips"],
    "author_id": "usr_123456",
    "created_at": "2025-01-15T10:30:00Z"
  }
}
```

### Get Discussion Details

Get detailed information about a discussion.

**Endpoint**: `GET /discussions/{discussion_id}`

**Response** (200):
```json
{
  "status": "success",
  "data": {
    "id": "disc_123456",
    "title": "Best practices for learning programming",
    "content": "What are your favorite techniques...",
    "author": {
      "id": "usr_123456",
      "name": "John Doe",
      "avatar": "https://cdn.greatnigeria.net/avatars/usr_123456.jpg",
      "badges": ["Top Contributor", "Programming Expert"]
    },
    "category": "Technology",
    "tags": ["programming", "learning", "tips"],
    "stats": {
      "views": 245,
      "replies": 18,
      "likes": 32,
      "bookmarks": 8
    },
    "replies": [
      {
        "id": "reply_123456",
        "content": "I recommend starting with Python...",
        "author": {
          "id": "usr_789012",
          "name": "Jane Smith",
          "avatar": "https://cdn.greatnigeria.net/avatars/usr_789012.jpg"
        },
        "likes": 12,
        "created_at": "2025-01-14T15:45:00Z"
      }
    ],
    "created_at": "2025-01-10T08:00:00Z",
    "updated_at": "2025-01-14T15:45:00Z"
  }
}
```

---

## Points Service API

### Get User Points

Get current user's points and level information.

**Endpoint**: `GET /points/me`
**Authentication**: Required

**Response** (200):
```json
{
  "status": "success",
  "data": {
    "user_id": "usr_123456",
    "current_points": 2750,
    "lifetime_points": 3980,
    "level": 5,
    "level_name": "Active Learner",
    "next_level": {
      "level": 6,
      "name": "Knowledge Seeker",
      "points_required": 3500,
      "points_remaining": 750
    },
    "recent_activities": [
      {
        "type": "chapter_completed",
        "points": 50,
        "description": "Completed Chapter 3: Variables and Data Types",
        "timestamp": "2025-01-15T10:30:00Z"
      }
    ]
  }
}
```

### Get Leaderboard

Get points leaderboard.

**Endpoint**: `GET /points/leaderboard`

**Query Parameters**:
- `period` (string): Time period (`daily`, `weekly`, `monthly`, `all-time`)
- `limit` (integer): Number of users to return

**Response** (200):
```json
{
  "status": "success",
  "data": {
    "period": "monthly",
    "rankings": [
      {
        "rank": 1,
        "user": {
          "id": "usr_456789",
          "name": "Sarah Johnson",
          "avatar": "https://cdn.greatnigeria.net/avatars/usr_456789.jpg"
        },
        "points": 8950,
        "level": 12
      }
    ],
    "current_user_rank": {
      "rank": 45,
      "points": 2750
    }
  }
}
```

---

## Payment Service API

### Create Payment Intent

Create a payment intent for subscription or content purchase.

**Endpoint**: `POST /payments/intents`
**Authentication**: Required

**Request Body**:
```json
{
  "amount": 2500,
  "currency": "NGN",
  "payment_method": "card",
  "description": "Premium subscription - Monthly",
  "metadata": {
    "subscription_type": "premium",
    "duration": "monthly"
  }
}
```

**Response** (201):
```json
{
  "status": "success",
  "data": {
    "id": "pi_123456",
    "amount": 2500,
    "currency": "NGN",
    "status": "requires_payment_method",
    "client_secret": "pi_123456_secret_xyz",
    "payment_methods": ["card", "bank_transfer", "mobile_money"],
    "expires_at": "2025-01-15T11:30:00Z",
    "created_at": "2025-01-15T10:30:00Z"
  }
}
```

### Get Payment History

Get user's payment history.

**Endpoint**: `GET /payments/history`
**Authentication**: Required

**Response** (200):
```json
{
  "status": "success",
  "data": [
    {
      "id": "payment_123456",
      "amount": 2500,
      "currency": "NGN",
      "status": "succeeded",
      "description": "Premium subscription - Monthly",
      "payment_method": "card",
      "receipt_url": "https://payments.greatnigeria.net/receipts/payment_123456",
      "created_at": "2025-01-01T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 12,
    "pages": 1
  }
}
```

---

## Celebration Service API

### List Celebration Entries

Get a list of celebration entries (people, places, events).

**Endpoint**: `GET /celebrate/entries`

**Query Parameters**:
- `type` (string): Entry type (`person`, `place`, `event`)
- `category` (string): Filter by category
- `featured` (boolean): Show only featured entries

**Response** (200):
```json
{
  "status": "success",
  "data": [
    {
      "id": "cel_123456",
      "type": "person",
      "title": "Chinua Achebe",
      "short_description": "Renowned Nigerian novelist and critic",
      "image": "https://cdn.greatnigeria.net/celebrate/people/chinua-achebe.jpg",
      "category": "Literature",
      "featured": true,
      "stats": {
        "views": 1250,
        "likes": 89,
        "shares": 24
      },
      "created_at": "2025-01-15T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 156,
    "pages": 8
  }
}
```

### Get Celebration Entry

Get detailed information about a celebration entry.

**Endpoint**: `GET /celebrate/entries/{entry_id}`

**Response** (200):
```json
{
  "status": "success",
  "data": {
    "id": "cel_123456",
    "type": "person",
    "title": "Chinua Achebe",
    "short_description": "Renowned Nigerian novelist and critic",
    "full_description": "Chinua Achebe was a Nigerian novelist, poet, and critic...",
    "image": "https://cdn.greatnigeria.net/celebrate/people/chinua-achebe.jpg",
    "category": "Literature",
    "featured": true,
    "facts": [
      {
        "label": "Born",
        "value": "November 16, 1930"
      },
      {
        "label": "Notable Work",
        "value": "Things Fall Apart"
      }
    ],
    "media": [
      {
        "type": "image",
        "url": "https://cdn.greatnigeria.net/celebrate/people/achebe-1.jpg",
        "caption": "Chinua Achebe at a literary event"
      }
    ],
    "related_entries": [
      {
        "id": "cel_789012",
        "title": "Wole Soyinka",
        "type": "person"
      }
    ],
    "stats": {
      "views": 1250,
      "likes": 89,
      "shares": 24,
      "comments": 15
    },
    "created_at": "2025-01-15T10:30:00Z",
    "updated_at": "2025-01-15T10:30:00Z"
  }
}
```

---

## Search Service API

### Search Content

Search across all platform content.

**Endpoint**: `GET /search`

**Query Parameters**:
- `q` (string): Search query
- `type` (string): Content type (`books`, `discussions`, `celebrate`, `users`)
- `filters` (object): Additional filters

**Response** (200):
```json
{
  "status": "success",
  "data": {
    "query": "programming",
    "total_results": 245,
    "results": [
      {
        "type": "book",
        "id": "book_123456",
        "title": "Introduction to Programming",
        "description": "Learn programming fundamentals...",
        "relevance_score": 0.95,
        "highlighted_text": "...fundamentals of <mark>programming</mark>..."
      },
      {
        "type": "discussion",
        "id": "disc_789012",
        "title": "Best Programming Languages for Beginners",
        "description": "Discussion about programming languages...",
        "relevance_score": 0.87,
        "highlighted_text": "...<mark>programming</mark> languages for beginners..."
      }
    ],
    "facets": {
      "type": {
        "books": 45,
        "discussions": 32,
        "celebrate": 8,
        "users": 12
      },
      "category": {
        "Technology": 78,
        "Education": 34,
        "Science": 12
      }
    }
  }
}
```

---

## Notification Service API

### List Notifications

Get user's notifications.

**Endpoint**: `GET /notifications`
**Authentication**: Required

**Query Parameters**:
- `unread` (boolean): Show only unread notifications
- `type` (string): Filter by notification type

**Response** (200):
```json
{
  "status": "success",
  "data": [
    {
      "id": "notif_123456",
      "type": "comment_reply",
      "title": "New reply to your comment",
      "message": "John Doe replied to your comment on 'Best Programming Languages'",
      "data": {
        "discussion_id": "disc_789012",
        "comment_id": "comment_456789"
      },
      "read": false,
      "created_at": "2025-01-15T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 23,
    "pages": 2
  },
  "unread_count": 5
}
```

### Mark Notifications as Read

Mark notifications as read.

**Endpoint**: `PUT /notifications/read`
**Authentication**: Required

**Request Body**:
```json
{
  "notification_ids": ["notif_123456", "notif_789012"]
}
```

**Response** (200):
```json
{
  "status": "success",
  "data": {
    "marked_read": 2,
    "total_unread": 3
  }
}
```

---

## WebSocket Events

### Real-time Connection

Connect to real-time events via WebSocket.

**Endpoint**: `wss://api.greatnigeria.net/ws`
**Authentication**: Required (via query parameter or header)

### Event Types

#### Discussion Updates
```json
{
  "type": "discussion_reply",
  "data": {
    "discussion_id": "disc_123456",
    "reply": {
      "id": "reply_789012",
      "content": "Great point!",
      "author": "John Doe"
    }
  }
}
```

#### Progress Updates
```json
{
  "type": "progress_updated",
  "data": {
    "user_id": "usr_123456",
    "book_id": "book_123456",
    "chapter_id": "chapter_001",
    "progress_percentage": 75
  }
}
```

#### Livestream Events
```json
{
  "type": "stream_started",
  "data": {
    "stream_id": "stream_123456",
    "title": "Introduction to React",
    "streamer": "Dr. Adebayo Ogundimu",
    "viewers": 125
  }
}
```

---

*This API reference is automatically generated and updated. For the most current information, refer to the interactive API documentation at https://api.greatnigeria.net/docs* 