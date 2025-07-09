# Development Guide
## Great Nigeria Library Platform

**Document Version**: 1.0  
**Last Updated**: January 2025  
**Maintained By**: Development Team  

---

## Table of Contents

1. [Getting Started](#getting-started)
2. [Development Environment Setup](#development-environment-setup)
3. [Project Structure](#project-structure)
4. [Coding Standards](#coding-standards)
5. [Git Workflow](#git-workflow)
6. [Testing Guidelines](#testing-guidelines)
7. [Database Management](#database-management)
8. [API Development](#api-development)
9. [Frontend Development](#frontend-development)
10. [Deployment & DevOps](#deployment--devops)
11. [Troubleshooting](#troubleshooting)
12. [Contributing Guidelines](#contributing-guidelines)

---

## Getting Started

### Prerequisites

#### Required Software
- **Go 1.21+** - Backend development ([Download](https://golang.org/dl/))
- **Node.js 18+** - Frontend development ([Download](https://nodejs.org/))
- **PostgreSQL 15+** - Primary database ([Download](https://postgresql.org/download/))
- **Redis 7+** - Caching and session storage ([Download](https://redis.io/download))
- **Docker & Docker Compose** - Local development ([Download](https://docker.com/))
- **Git** - Version control ([Download](https://git-scm.com/downloads))

#### Optional Tools
- **VS Code** - Recommended IDE with Go and React extensions
- **Postman/Insomnia** - API testing
- **pgAdmin** - PostgreSQL administration
- **Redis Desktop Manager** - Redis administration

### Quick Setup

```bash
# 1. Clone the repository
git clone https://github.com/greatnigeria/library.git
cd library

# 2. Start development dependencies
docker-compose up -d postgres redis elasticsearch

# 3. Install backend dependencies
go mod download

# 4. Install frontend dependencies
cd frontend
npm install
cd ..

# 5. Copy environment configuration
cp .env.example .env

# 6. Run database migrations
go run cmd/migrate/main.go up

# 7. Populate sample data
go run scripts/populate_sample_data.go

# 8. Start development servers (in separate terminals)
go run cmd/api-gateway/main.go          # Terminal 1
go run cmd/auth-service/main.go         # Terminal 2
go run cmd/content-service/main.go      # Terminal 3
cd frontend && npm start                # Terminal 4
```

---

## Development Environment Setup

### Local Environment Configuration

#### Environment Variables
Create a `.env` file in the project root:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=greatnigeria_dev
DB_SSL_MODE=disable

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT Configuration
JWT_SECRET=your_super_secret_jwt_key_here
JWT_EXPIRY=24h
REFRESH_TOKEN_EXPIRY=168h

# API Configuration
API_PORT=8080
API_HOST=localhost
FRONTEND_URL=http://localhost:3000

# External Services
PAYSTACK_SECRET_KEY=your_paystack_secret
FLUTTERWAVE_SECRET_KEY=your_flutterwave_secret
SENDGRID_API_KEY=your_sendgrid_key

# Development Settings
ENV=development
LOG_LEVEL=debug
ENABLE_PROFILING=true
```

#### Docker Development Setup

Use Docker Compose for consistent development environment:

```yaml
# docker-compose.dev.yml
version: '3.8'

services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: greatnigeria_dev
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: devpassword
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data

  elasticsearch:
    image: elasticsearch:8.10.0
    environment:
      - discovery.type=single-node
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
      - xpack.security.enabled=false
    ports:
      - "9200:9200"
    volumes:
      - es_data:/usr/share/elasticsearch/data

volumes:
  postgres_data:
  redis_data:
  es_data:
```

Start with: `docker-compose -f docker-compose.dev.yml up -d`

### IDE Configuration

#### VS Code Settings
Create `.vscode/settings.json`:

```json
{
  "go.toolsManagement.checkForUpdates": "local",
  "go.useLanguageServer": true,
  "go.formatTool": "goimports",
  "go.lintTool": "golangci-lint",
  "go.testFlags": ["-v"],
  "editor.formatOnSave": true,
  "typescript.preferences.importModuleSpecifier": "relative",
  "eslint.format.enable": true,
  "prettier.singleQuote": true,
  "prettier.trailingComma": "es5"
}
```

#### Recommended Extensions
- **Go** (golang.go)
- **ES7+ React/Redux/React-Native snippets**
- **Prettier - Code formatter**
- **ESLint**
- **GitLens**
- **Docker**
- **PostgreSQL** (ms-ossdata.vscode-postgresql)

---

## Project Structure

### Backend Structure

```
great-nigeria-library/
├── cmd/                      # Service entry points
│   ├── api-gateway/         # API Gateway service
│   ├── auth-service/        # Authentication service
│   ├── content-service/     # Content management service
│   ├── discussion-service/  # Community discussions
│   ├── livestream-service/  # Live streaming
│   ├── payment-service/     # Payment processing
│   ├── points-service/      # Gamification
│   └── migrate/            # Database migration tool
├── internal/                # Private application code
│   ├── auth/               # Auth service implementation
│   │   ├── handlers/       # HTTP handlers
│   │   ├── repository/     # Data access layer
│   │   ├── service/        # Business logic
│   │   └── models/         # Data models
│   ├── content/            # Content service implementation
│   ├── discussion/         # Discussion service implementation
│   └── ...                 # Other services
├── pkg/                    # Shared library code
│   ├── common/            # Common utilities
│   │   ├── config/        # Configuration management
│   │   ├── database/      # Database connections
│   │   ├── logger/        # Logging utilities
│   │   ├── middleware/    # HTTP middleware
│   │   └── utils/         # Helper functions
│   ├── auth/              # Authentication utilities
│   └── errors/            # Error handling
├── migrations/             # Database migrations
├── scripts/               # Utility scripts
├── docs/                  # Documentation
├── frontend/              # React frontend application
└── deployments/           # Deployment configurations
```

### Frontend Structure

```
frontend/
├── public/                 # Static assets
├── src/
│   ├── api/               # API service layer
│   ├── components/        # Reusable UI components
│   │   ├── common/       # Generic components
│   │   ├── auth/         # Authentication components
│   │   ├── book/         # Book-related components
│   │   └── ...           # Feature-specific components
│   ├── features/         # Redux slices and feature logic
│   │   ├── auth/         # Authentication state
│   │   ├── books/        # Books state
│   │   └── ...           # Other features
│   ├── hooks/            # Custom React hooks
│   ├── layouts/          # Page layouts
│   ├── pages/            # Page components
│   ├── store/            # Redux store configuration
│   ├── theme/            # UI theme and styling
│   ├── utils/            # Utility functions
│   └── types/            # TypeScript type definitions
├── package.json
└── tsconfig.json
```

---

## Coding Standards

### Go Coding Standards

#### Code Formatting
```bash
# Format code
go fmt ./...

# Import organization
goimports -w .

# Linting
golangci-lint run
```

#### Naming Conventions
```go
// Package names: lowercase, single word
package auth

// Interface names: noun or noun phrase
type UserRepository interface {
    GetUser(id string) (*User, error)
    CreateUser(user *User) error
}

// Struct names: PascalCase
type UserService struct {
    repo UserRepository
    log  logger.Logger
}

// Method names: PascalCase
func (s *UserService) RegisterUser(req *RegisterRequest) (*User, error) {
    // Implementation
}

// Constants: PascalCase or SCREAMING_SNAKE_CASE for exported
const (
    DefaultPageSize = 20
    MAX_FILE_SIZE   = 10 * 1024 * 1024 // 10MB
)

// Variables: camelCase
var userCache = make(map[string]*User)
```

#### Error Handling
```go
// Always handle errors explicitly
func (s *UserService) GetUser(id string) (*User, error) {
    user, err := s.repo.GetUser(id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user %s: %w", id, err)
    }
    
    if user == nil {
        return nil, errors.New("user not found")
    }
    
    return user, nil
}

// Use custom error types for business logic
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
}
```

#### API Handler Pattern
```go
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "invalid request format"})
        return
    }
    
    if err := req.Validate(); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    user, err := h.service.CreateUser(&req)
    if err != nil {
        h.logger.Error("failed to create user", "error", err)
        c.JSON(500, gin.H{"error": "internal server error"})
        return
    }
    
    c.JSON(201, gin.H{"data": user})
}
```

### TypeScript/React Standards

#### Component Structure
```typescript
// Component props interface
interface BookCardProps {
  book: Book;
  onSelect?: (book: Book) => void;
  className?: string;
}

// React component with TypeScript
export const BookCard: React.FC<BookCardProps> = ({ 
  book, 
  onSelect, 
  className = '' 
}) => {
  const handleClick = useCallback(() => {
    onSelect?.(book);
  }, [book, onSelect]);

  return (
    <div 
      className={`book-card ${className}`}
      onClick={handleClick}
    >
      <img src={book.coverImage} alt={book.title} />
      <h3>{book.title}</h3>
      <p>{book.author}</p>
    </div>
  );
};
```

#### Redux Slice Pattern
```typescript
// Redux slice with RTK
interface BooksState {
  books: Book[];
  currentBook: Book | null;
  loading: boolean;
  error: string | null;
}

const initialState: BooksState = {
  books: [],
  currentBook: null,
  loading: false,
  error: null,
};

export const booksSlice = createSlice({
  name: 'books',
  initialState,
  reducers: {
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },
    setError: (state, action: PayloadAction<string>) => {
      state.error = action.payload;
      state.loading = false;
    },
    setBooksSuccess: (state, action: PayloadAction<Book[]>) => {
      state.books = action.payload;
      state.loading = false;
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchBooks.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchBooks.fulfilled, (state, action) => {
        state.books = action.payload;
        state.loading = false;
      })
      .addCase(fetchBooks.rejected, (state, action) => {
        state.error = action.error.message || 'Failed to fetch books';
        state.loading = false;
      });
  },
});
```

#### API Service Pattern
```typescript
// API service with error handling
class BookService {
  private apiClient: AxiosInstance;

  constructor() {
    this.apiClient = createApiClient();
  }

  async getBooks(params?: GetBooksParams): Promise<PaginatedResponse<Book>> {
    try {
      const response = await this.apiClient.get<PaginatedResponse<Book>>('/books', {
        params,
      });
      return response.data;
    } catch (error) {
      if (isAxiosError(error)) {
        throw new ApiError(
          error.response?.data?.message || 'Failed to fetch books',
          error.response?.status || 500
        );
      }
      throw error;
    }
  }

  async getBook(id: string): Promise<Book> {
    try {
      const response = await this.apiClient.get<ApiResponse<Book>>(`/books/${id}`);
      return response.data.data;
    } catch (error) {
      if (isAxiosError(error)) {
        throw new ApiError(
          error.response?.data?.message || 'Failed to fetch book',
          error.response?.status || 500
        );
      }
      throw error;
    }
  }
}
```

---

## Git Workflow

### Branch Strategy

We use **Git Flow** with the following branch types:

#### Main Branches
- **`main`** - Production-ready code
- **`develop`** - Integration branch for features

#### Supporting Branches
- **`feature/*`** - New features (`feature/user-authentication`)
- **`bugfix/*`** - Bug fixes (`bugfix/login-validation`)
- **`hotfix/*`** - Critical production fixes (`hotfix/security-patch`)
- **`release/*`** - Release preparation (`release/v1.2.0`)

### Commit Convention

We follow **Conventional Commits** specification:

```bash
# Format: <type>[optional scope]: <description>
# 
# Types: feat, fix, docs, style, refactor, test, chore

# Examples:
git commit -m "feat(auth): add two-factor authentication"
git commit -m "fix(api): resolve user registration validation"
git commit -m "docs(readme): update installation instructions"
git commit -m "refactor(database): optimize query performance"
git commit -m "test(auth): add unit tests for login service"
git commit -m "chore(deps): update dependencies to latest versions"
```

### Pull Request Process

#### 1. Create Feature Branch
```bash
git checkout develop
git pull origin develop
git checkout -b feature/your-feature-name
```

#### 2. Development & Commits
```bash
# Make changes
git add .
git commit -m "feat(scope): add new feature"

# Push to remote
git push origin feature/your-feature-name
```

#### 3. Create Pull Request
- **Title**: Clear, descriptive title
- **Description**: What, why, and how
- **Reviewers**: Assign appropriate reviewers
- **Labels**: Add relevant labels (enhancement, bug, documentation)

#### 4. PR Template
```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Checklist
- [ ] Code follows project coding standards
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] No merge conflicts
```

### Code Review Guidelines

#### For Authors
- **Small PRs**: Keep changes focused and reviewable
- **Clear Description**: Explain what and why
- **Self-Review**: Review your own code first
- **Tests**: Include appropriate tests
- **Documentation**: Update relevant documentation

#### For Reviewers
- **Timely Reviews**: Review within 24 hours
- **Constructive Feedback**: Be specific and helpful
- **Code Quality**: Check for standards compliance
- **Security**: Look for security vulnerabilities
- **Performance**: Consider performance implications

---

## Testing Guidelines

### Backend Testing

#### Unit Tests
```go
// Example unit test
func TestUserService_CreateUser(t *testing.T) {
    // Setup
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo, logger.NewNoop())
    
    // Test data
    req := &CreateUserRequest{
        Email:    "test@example.com",
        Password: "securepassword",
        Name:     "Test User",
    }
    
    // Mock expectations
    mockRepo.On("CreateUser", mock.AnythingOfType("*User")).Return(nil)
    mockRepo.On("GetUserByEmail", req.Email).Return(nil, nil)
    
    // Execute
    user, err := service.CreateUser(req)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, req.Email, user.Email)
    assert.Equal(t, req.Name, user.Name)
    
    // Verify mocks
    mockRepo.AssertExpectations(t)
}
```

#### Integration Tests
```go
func TestAuthAPI_Integration(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    // Setup test server
    router := setupTestRouter(db)
    
    // Test user registration
    t.Run("RegisterUser", func(t *testing.T) {
        reqBody := `{
            "email": "test@example.com",
            "password": "securepassword",
            "name": "Test User"
        }`
        
        req := httptest.NewRequest("POST", "/auth/register", strings.NewReader(reqBody))
        req.Header.Set("Content-Type", "application/json")
        
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)
        
        assert.Equal(t, http.StatusCreated, w.Code)
        
        var response map[string]interface{}
        err := json.Unmarshal(w.Body.Bytes(), &response)
        assert.NoError(t, err)
        assert.Contains(t, response, "data")
    })
}
```

### Frontend Testing

#### Component Tests
```typescript
// Component test with React Testing Library
import { render, screen, fireEvent } from '@testing-library/react';
import { Provider } from 'react-redux';
import { BookCard } from './BookCard';
import { createTestStore } from '../test-utils';

const mockBook: Book = {
  id: '1',
  title: 'Test Book',
  author: 'Test Author',
  coverImage: 'test-cover.jpg',
};

describe('BookCard', () => {
  it('renders book information correctly', () => {
    const onSelect = jest.fn();
    
    render(
      <BookCard book={mockBook} onSelect={onSelect} />
    );
    
    expect(screen.getByText('Test Book')).toBeInTheDocument();
    expect(screen.getByText('Test Author')).toBeInTheDocument();
    expect(screen.getByAltText('Test Book')).toHaveAttribute('src', 'test-cover.jpg');
  });
  
  it('calls onSelect when clicked', () => {
    const onSelect = jest.fn();
    
    render(
      <BookCard book={mockBook} onSelect={onSelect} />
    );
    
    fireEvent.click(screen.getByRole('button', { name: /test book/i }));
    
    expect(onSelect).toHaveBeenCalledWith(mockBook);
  });
});
```

#### Redux Tests
```typescript
// Redux slice tests
import { configureStore } from '@reduxjs/toolkit';
import { booksSlice, fetchBooks } from './booksSlice';

describe('booksSlice', () => {
  let store: ReturnType<typeof configureStore>;
  
  beforeEach(() => {
    store = configureStore({
      reducer: {
        books: booksSlice.reducer,
      },
    });
  });
  
  it('should handle setLoading', () => {
    store.dispatch(booksSlice.actions.setLoading(true));
    
    const state = store.getState().books;
    expect(state.loading).toBe(true);
  });
  
  it('should handle fetchBooks.fulfilled', () => {
    const mockBooks = [mockBook];
    
    store.dispatch(fetchBooks.fulfilled(mockBooks, 'requestId', undefined));
    
    const state = store.getState().books;
    expect(state.books).toEqual(mockBooks);
    expect(state.loading).toBe(false);
    expect(state.error).toBeNull();
  });
});
```

### Test Commands

```bash
# Backend tests
go test ./...                          # Run all tests
go test -v ./internal/auth/...         # Run tests with verbose output
go test -cover ./...                   # Run tests with coverage
go test -race ./...                    # Run tests with race detection

# Frontend tests
npm test                               # Run tests in watch mode
npm run test:ci                        # Run tests once
npm run test:coverage                  # Run tests with coverage
npm run test:e2e                       # Run end-to-end tests
```

---

## Database Management

### Migration System

#### Creating Migrations
```bash
# Create new migration
go run cmd/migrate/main.go create add_user_preferences

# This creates:
# migrations/YYYYMMDD_HHMMSS_add_user_preferences.up.sql
# migrations/YYYYMMDD_HHMMSS_add_user_preferences.down.sql
```

#### Migration File Example
```sql
-- migrations/20250115_120000_add_user_preferences.up.sql
CREATE TABLE user_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    theme VARCHAR(20) DEFAULT 'light',
    language VARCHAR(10) DEFAULT 'en',
    notifications_enabled BOOLEAN DEFAULT true,
    email_notifications BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(user_id)
);

CREATE INDEX idx_user_preferences_user_id ON user_preferences(user_id);

-- migrations/20250115_120000_add_user_preferences.down.sql
DROP TABLE IF EXISTS user_preferences;
```

#### Running Migrations
```bash
# Apply all pending migrations
go run cmd/migrate/main.go up

# Rollback last migration
go run cmd/migrate/main.go down

# Check migration status
go run cmd/migrate/main.go status

# Reset database (development only)
go run cmd/migrate/main.go reset
```

### Database Schema Guidelines

#### Naming Conventions
- **Tables**: `snake_case`, plural nouns (`users`, `book_chapters`)
- **Columns**: `snake_case` (`created_at`, `user_id`)
- **Indexes**: `idx_tablename_columnname` (`idx_users_email`)
- **Foreign Keys**: `fk_tablename_referencedtable` (`fk_posts_users`)

#### Standard Columns
```sql
-- Every table should have these columns
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
```

#### Foreign Key Patterns
```sql
-- User reference
user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE

-- Soft delete pattern
deleted_at TIMESTAMP WITH TIME ZONE NULL
```

### Database Seeding

#### Sample Data Script
```go
// scripts/seed_development_data.go
func seedDevelopmentData() error {
    // Seed users
    users := []User{
        {
            ID:       uuid.New(),
            Email:    "admin@greatnigeria.net",
            Name:     "Admin User",
            Role:     "admin",
            Verified: true,
        },
        {
            ID:       uuid.New(),
            Email:    "user@example.com",
            Name:     "Test User",
            Role:     "user",
            Verified: true,
        },
    }
    
    for _, user := range users {
        if err := userRepo.CreateUser(&user); err != nil {
            return fmt.Errorf("failed to seed user: %w", err)
        }
    }
    
    // Seed books, etc.
    return nil
}
```

---

## API Development

### API Design Principles

#### RESTful Design
```
GET    /api/v1/books           # List books
POST   /api/v1/books           # Create book
GET    /api/v1/books/{id}      # Get specific book
PUT    /api/v1/books/{id}      # Update book
DELETE /api/v1/books/{id}      # Delete book

GET    /api/v1/books/{id}/chapters     # List book chapters
POST   /api/v1/books/{id}/chapters     # Add chapter to book
```

#### Request/Response Format
```json
// Request format
{
  "data": {
    "title": "Book Title",
    "author": "Author Name",
    "description": "Book description"
  }
}

// Success response format
{
  "status": "success",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "title": "Book Title",
    "author": "Author Name",
    "created_at": "2025-01-15T10:30:00Z"
  },
  "meta": {
    "timestamp": "2025-01-15T10:30:00Z",
    "request_id": "req_123456"
  }
}

// Error response format
{
  "status": "error",
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": {
      "title": "Title is required",
      "author": "Author must be at least 2 characters"
    }
  },
  "meta": {
    "timestamp": "2025-01-15T10:30:00Z",
    "request_id": "req_123456"
  }
}
```

#### Pagination
```json
// Paginated response
{
  "status": "success",
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "pages": 8,
    "has_next": true,
    "has_prev": false
  }
}
```

### OpenAPI Documentation

#### Swagger Annotations
```go
// @Summary Create a new book
// @Description Create a new book with the provided information
// @Tags books
// @Accept json
// @Produce json
// @Param book body CreateBookRequest true "Book information"
// @Success 201 {object} BookResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /books [post]
func (h *BookHandler) CreateBook(c *gin.Context) {
    // Implementation
}
```

### Authentication Middleware

```go
// JWT Authentication middleware
func JWTAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }
        
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        
        token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("JWT_SECRET")), nil
        })
        
        if err != nil || !token.Valid {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        claims := token.Claims.(*Claims)
        c.Set("user_id", claims.UserID)
        c.Set("user_role", claims.Role)
        
        c.Next()
    }
}
```

---

## Frontend Development

### Project Setup

#### Package.json Scripts
```json
{
  "scripts": {
    "start": "vite",
    "build": "tsc && vite build",
    "test": "jest",
    "test:watch": "jest --watch",
    "test:coverage": "jest --coverage",
    "lint": "eslint src --ext .ts,.tsx",
    "lint:fix": "eslint src --ext .ts,.tsx --fix",
    "type-check": "tsc --noEmit",
    "storybook": "start-storybook -p 6006"
  }
}
```

### State Management

#### Redux Store Setup
```typescript
// store/index.ts
import { configureStore } from '@reduxjs/toolkit';
import { authSlice } from '../features/auth/authSlice';
import { booksSlice } from '../features/books/booksSlice';

export const store = configureStore({
  reducer: {
    auth: authSlice.reducer,
    books: booksSlice.reducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        ignoredActions: ['persist/PERSIST', 'persist/REHYDRATE'],
      },
    }),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
```

#### Typed Hooks
```typescript
// hooks/redux.ts
import { TypedUseSelectorHook, useDispatch, useSelector } from 'react-redux';
import type { RootState, AppDispatch } from '../store';

export const useAppDispatch = () => useDispatch<AppDispatch>();
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector;
```

### UI Component Guidelines

#### Component Composition
```typescript
// Good: Composable component structure
export interface ButtonProps {
  variant?: 'primary' | 'secondary' | 'outline';
  size?: 'small' | 'medium' | 'large';
  disabled?: boolean;
  loading?: boolean;
  children: React.ReactNode;
  onClick?: () => void;
}

export const Button: React.FC<ButtonProps> = ({
  variant = 'primary',
  size = 'medium',
  disabled = false,
  loading = false,
  children,
  onClick,
}) => {
  const classes = cn(
    'button',
    `button--${variant}`,
    `button--${size}`,
    {
      'button--disabled': disabled,
      'button--loading': loading,
    }
  );

  return (
    <button 
      className={classes}
      disabled={disabled || loading}
      onClick={onClick}
    >
      {loading && <Spinner size="small" />}
      {children}
    </button>
  );
};
```

#### Form Handling
```typescript
// Form with validation using react-hook-form and zod
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';

const loginSchema = z.object({
  email: z.string().email('Invalid email address'),
  password: z.string().min(8, 'Password must be at least 8 characters'),
});

type LoginFormData = z.infer<typeof loginSchema>;

export const LoginForm: React.FC = () => {
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
  });

  const onSubmit = async (data: LoginFormData) => {
    try {
      await authService.login(data);
    } catch (error) {
      // Handle error
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <div>
        <label htmlFor="email">Email</label>
        <input
          {...register('email')}
          type="email"
          id="email"
          aria-invalid={errors.email ? 'true' : 'false'}
        />
        {errors.email && (
          <span role="alert">{errors.email.message}</span>
        )}
      </div>

      <div>
        <label htmlFor="password">Password</label>
        <input
          {...register('password')}
          type="password"
          id="password"
          aria-invalid={errors.password ? 'true' : 'false'}
        />
        {errors.password && (
          <span role="alert">{errors.password.message}</span>
        )}
      </div>

      <Button type="submit" loading={isSubmitting}>
        Login
      </Button>
    </form>
  );
};
```

---

## Deployment & DevOps

### Docker Configuration

#### Backend Dockerfile
```dockerfile
# Multi-stage build for Go services
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service cmd/auth-service/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/auth-service .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080
CMD ["./auth-service"]
```

#### Frontend Dockerfile
```dockerfile
# Multi-stage build for React app
FROM node:18-alpine AS builder

WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci --only=production

COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

### Kubernetes Deployment

#### Service Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auth-service
  namespace: greatnigeria
spec:
  replicas: 3
  selector:
    matchLabels:
      app: auth-service
  template:
    metadata:
      labels:
        app: auth-service
    spec:
      containers:
      - name: auth-service
        image: greatnigeria/auth-service:v1.0.0
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: database-secret
              key: url
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: auth-secret
              key: jwt-secret
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
```

### CI/CD Pipeline

#### GitLab CI Configuration
```yaml
# .gitlab-ci.yml
stages:
  - test
  - build
  - deploy

variables:
  DOCKER_REGISTRY: registry.greatnigeria.net
  GO_VERSION: "1.21"

# Backend tests
backend-test:
  stage: test
  image: golang:${GO_VERSION}
  script:
    - go mod download
    - go test -race -cover ./...
  coverage: '/coverage: \d+\.\d+% of statements/'

# Frontend tests
frontend-test:
  stage: test
  image: node:18
  script:
    - cd frontend
    - npm ci
    - npm run test:ci
    - npm run lint
  coverage: '/All files[^|]*\|[^|]*\s+([\d\.]+)/'

# Build and push Docker images
build-services:
  stage: build
  image: docker:latest
  services:
    - docker:dind
  script:
    - docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY
    - docker build -t $CI_REGISTRY_IMAGE/auth-service:$CI_COMMIT_SHA -f cmd/auth-service/Dockerfile .
    - docker push $CI_REGISTRY_IMAGE/auth-service:$CI_COMMIT_SHA
  only:
    - main
    - develop

# Deploy to staging
deploy-staging:
  stage: deploy
  image: bitnami/kubectl:latest
  script:
    - kubectl set image deployment/auth-service auth-service=$CI_REGISTRY_IMAGE/auth-service:$CI_COMMIT_SHA -n staging
    - kubectl rollout status deployment/auth-service -n staging
  environment:
    name: staging
    url: https://staging.greatnigeria.net
  only:
    - develop

# Deploy to production
deploy-production:
  stage: deploy
  image: bitnami/kubectl:latest
  script:
    - kubectl set image deployment/auth-service auth-service=$CI_REGISTRY_IMAGE/auth-service:$CI_COMMIT_SHA -n production
    - kubectl rollout status deployment/auth-service -n production
  environment:
    name: production
    url: https://greatnigeria.net
  only:
    - main
  when: manual
```

---

## Troubleshooting

### Common Issues

#### Backend Issues

##### Database Connection Issues
```bash
# Check database connectivity
psql -h localhost -U postgres -d greatnigeria_dev

# Verify environment variables
echo $DATABASE_URL

# Check if database exists
psql -h localhost -U postgres -c "\l"

# Reset database (development only)
dropdb greatnigeria_dev && createdb greatnigeria_dev
go run cmd/migrate/main.go up
```

##### Service Compilation Issues
```bash
# Missing pkg/ directory
mkdir -p pkg/common/{config,database,logger,middleware,utils}
mkdir -p pkg/{auth,errors}

# Create basic structure files
touch pkg/common/config/config.go
touch pkg/common/database/database.go
touch pkg/common/logger/logger.go
touch pkg/errors/types.go

# Update go.mod
go mod tidy
```

##### JWT Token Issues
```bash
# Verify JWT secret is set
echo $JWT_SECRET

# Check token format
curl -H "Authorization: Bearer your-token" http://localhost:8080/api/user

# Debug token claims
go run scripts/debug_jwt.go your-token
```

#### Frontend Issues

##### Build Issues
```bash
# Clear node_modules and reinstall
rm -rf node_modules package-lock.json
npm install

# Check for TypeScript errors
npx tsc --noEmit

# Verify environment variables
cat .env.local
```

##### API Connection Issues
```bash
# Check if backend is running
curl http://localhost:8080/health

# Verify API base URL
grep -r "baseURL" src/api/

# Check network tab in browser dev tools
# Look for CORS issues or network errors
```

##### Redux State Issues
```typescript
// Debug Redux state
import { store } from './store';

// Log current state
console.log('Current state:', store.getState());

// Listen to state changes
store.subscribe(() => {
  console.log('State changed:', store.getState());
});
```

### Performance Issues

#### Backend Performance
```bash
# Profile Go application
go tool pprof http://localhost:8080/debug/pprof/profile

# Check memory usage
go tool pprof http://localhost:8080/debug/pprof/heap

# Database query optimization
EXPLAIN ANALYZE SELECT * FROM users WHERE email = 'user@example.com';
```

#### Frontend Performance
```javascript
// React DevTools Profiler
// Use React.Profiler to identify slow components

// Bundle analysis
npm run build
npx webpack-bundle-analyzer build/static/js/*.js

// Lighthouse audit
npx lighthouse http://localhost:3000 --output html --output-path ./lighthouse-report.html
```

### Development Environment Issues

#### Port Conflicts
```bash
# Check what's using port 8080
lsof -i :8080

# Kill process using port
kill -9 $(lsof -t -i:8080)

# Use different port
export API_PORT=8081
```

#### Docker Issues
```bash
# Clean Docker system
docker system prune -a

# Restart Docker services
docker-compose down
docker-compose up -d

# Check container logs
docker-compose logs -f postgres
docker-compose logs -f redis
```

#### Environment Variable Issues
```bash
# Check all environment variables
env | grep -E "(DB_|REDIS_|JWT_)"

# Load environment file
source .env

# Verify variables are exported
echo $DATABASE_URL
```

### Getting Help

#### Internal Resources
1. **Documentation**: Check this guide and `/docs` directory
2. **Code Comments**: Look for inline documentation
3. **Test Files**: Examples in `*_test.go` and `*.test.ts` files
4. **Git History**: `git log --oneline --graph` for context

#### External Resources
1. **Go Documentation**: [golang.org/doc](https://golang.org/doc)
2. **React Documentation**: [reactjs.org](https://reactjs.org)
3. **PostgreSQL Docs**: [postgresql.org/docs](https://postgresql.org/docs)
4. **Docker Docs**: [docs.docker.com](https://docs.docker.com)

#### Team Communication
- **Slack/Discord**: #development channel
- **GitHub Issues**: For bug reports and feature requests
- **Weekly Standups**: Monday 9 AM WAT
- **Code Reviews**: Tag relevant team members

---

## Contributing Guidelines

### Getting Started

1. **Read this guide thoroughly**
2. **Set up development environment**
3. **Pick an issue from GitHub**
4. **Create feature branch**
5. **Make changes following coding standards**
6. **Add tests for new functionality**
7. **Update documentation**
8. **Submit pull request**

### Code of Conduct

- **Be respectful** in all interactions
- **Provide constructive feedback** in code reviews
- **Help new contributors** get started
- **Follow project standards** and guidelines
- **Prioritize user experience** in all decisions

### Recognition

Contributors will be recognized through:
- **Contributor list** in README
- **Release notes** acknowledgments
- **Community highlights** in newsletters
- **Maintainer privileges** for consistent contributors

---

*This development guide is a living document. Please suggest improvements and updates as the project evolves.* 