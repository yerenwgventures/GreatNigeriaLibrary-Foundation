# Database Migrations Guide

This document explains the GORM auto-migration system implemented in the Great Nigeria Library Foundation backend.

## Overview

The database migration system has been migrated from static SQL files to GORM auto-migration, providing:

- **Automatic Schema Management**: Tables and indexes created automatically from Go models
- **Service-Specific Migrations**: Each service manages its own database schema
- **Consistent Schema**: Database structure matches GORM model definitions
- **Demo Data Insertion**: Automatic insertion of demo data for development

## Migration Architecture

### Core Components

1. **MigrationService** (`database/migrations.go`): Centralized migration management
2. **Service-Specific Migrations**: Each service runs its own migrations on startup
3. **GORM Models**: Database schema defined in Go structs with GORM tags
4. **Demo Data**: Automatic insertion of sample data for development

### Migration Flow

```
Application Startup
       ↓
Database Connection
       ↓
Service Migration
       ↓
GORM AutoMigrate
       ↓
Index Creation
       ↓
Demo Data Insertion
```

## Service Migrations

### Auth Service Migration

**Models Migrated:**
- `User` - User accounts and authentication
- `Session` - User sessions
- `PasswordResetToken` - Password reset tokens
- `EmailVerificationToken` - Email verification tokens
- `UserTrustLevel` - User trust levels

**Indexes Created:**
- `idx_users_email` - Email lookup
- `idx_users_username` - Username lookup
- `idx_sessions_user_id` - Session user lookup
- `idx_sessions_expires_at` - Session expiration cleanup

### Content Service Migration

**Tables Created:**
- `content_categories` - Content categorization
- `demo_content` - Sample content for demonstration

**Features:**
- UUID primary keys
- Foreign key relationships
- Automatic timestamps
- Content categorization support

### Discussion Service Migration

**Tables Created:**
- `forum_categories` - Discussion categories
- `forum_topics` - Discussion topics
- `forum_replies` - Topic replies

**Features:**
- Threaded discussions
- Topic pinning and locking
- Reply counting
- View tracking

### Groups Service Migration

**Tables Created:**
- `user_groups` - User groups
- `group_members` - Group membership

**Features:**
- Public/private groups
- Member roles
- Group statistics
- Membership tracking

## Usage

### Service-Level Migration

Each service automatically runs its migrations on startup:

```go
// In service main.go
db, err := database.NewDatabase(cfg)
if err != nil {
    logger.Fatal("Failed to connect to database: " + err.Error())
}

// Run service-specific migration
migrationService := database.NewMigrationService(db, logger)
if err := migrationService.MigrateAuthService(); err != nil {
    logger.Fatal("Failed to run auth service migrations: " + err.Error())
}
```

### Foundation-Wide Migration

The main application runs all migrations:

```go
// In main.go
migrationService := database.NewMigrationService(dbConn.DB, appLogger)
if err := migrationService.MigrateFoundation(); err != nil {
    appLogger.Fatal("Failed to run foundation migrations: " + err.Error())
}
```

### Available Migration Methods

```go
// Individual service migrations
migrationService.MigrateAuthService()
migrationService.MigrateContentService()
migrationService.MigrateDiscussionService()
migrationService.MigrateGroupsService()

// Complete foundation migration
migrationService.MigrateFoundation()
```

## Model Definitions

### GORM Model Example

```go
type User struct {
    ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Username        string    `gorm:"uniqueIndex;size:255;not null" json:"username"`
    Email           string    `gorm:"uniqueIndex;size:255;not null" json:"email"`
    Password        string    `gorm:"size:255;not null" json:"-"`
    FullName        string    `gorm:"size:255" json:"full_name"`
    IsActive        bool      `gorm:"default:true" json:"is_active"`
    CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
```

### GORM Tags Reference

- `primaryKey` - Primary key field
- `autoIncrement` - Auto-incrementing field
- `uniqueIndex` - Unique index
- `size:255` - Field size
- `not null` - NOT NULL constraint
- `default:true` - Default value
- `autoCreateTime` - Automatic creation timestamp
- `autoUpdateTime` - Automatic update timestamp

## Demo Data

### Automatic Demo Data Insertion

The migration system automatically inserts demo data:

```go
// Content categories
INSERT INTO content_categories (name, slug, description) VALUES
('Platform Guide', 'platform-guide', 'How to use the platform'),
('Educational Content', 'educational', 'Educational materials'),
('Community Guidelines', 'community', 'Community rules');

// Forum categories
INSERT INTO forum_categories (name, description, sort_order) VALUES
('General Discussion', 'General topics', 1),
('Platform Help', 'Help and support', 2),
('Feature Requests', 'Suggest new features', 3);
```

### Demo Data Features

- **Conflict Handling**: Uses `ON CONFLICT DO NOTHING` to prevent duplicates
- **Realistic Data**: Meaningful sample data for development
- **Categorization**: Proper content and forum categories
- **Relationships**: Maintains foreign key relationships

## Migration Benefits

### Before (Static SQL)

- **Manual Schema Management**: Required manual SQL file updates
- **Deployment Complexity**: Had to manage SQL files separately
- **Model Drift**: Risk of models and database schema getting out of sync
- **Environment Setup**: Required running SQL scripts manually

### After (GORM Auto-Migration)

- **Automatic Schema**: Database schema automatically matches Go models
- **Zero Deployment Overhead**: Migrations run automatically on startup
- **Model Consistency**: Database always matches GORM model definitions
- **Environment Agnostic**: Works in development, staging, and production

## Best Practices

### 1. Model-First Development

```go
// Define models with proper GORM tags
type Article struct {
    ID        uint      `gorm:"primaryKey"`
    Title     string    `gorm:"size:255;not null"`
    Content   string    `gorm:"type:text"`
    AuthorID  uint      `gorm:"not null;index"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
    
    // Relationships
    Author User `gorm:"foreignKey:AuthorID"`
}
```

### 2. Index Management

```go
// Create indexes for performance
indexes := []string{
    "CREATE INDEX IF NOT EXISTS idx_articles_author ON articles(author_id)",
    "CREATE INDEX IF NOT EXISTS idx_articles_created ON articles(created_at)",
}

for _, index := range indexes {
    if err := m.db.Exec(index).Error; err != nil {
        m.logger.WithField("index", index).Error("Failed to create index")
    }
}
```

### 3. Safe Migration Practices

```go
// Always use IF NOT EXISTS for custom tables
CREATE TABLE IF NOT EXISTS custom_table (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL
);

// Handle migration errors gracefully
if err := m.db.AutoMigrate(models...); err != nil {
    return fmt.Errorf("migration failed: %w", err)
}
```

## Troubleshooting

### Common Issues

1. **Migration Failures**
   ```
   Error: migration failed: column "new_field" already exists
   ```
   **Solution**: GORM handles existing columns automatically

2. **Index Creation Errors**
   ```
   Error: index "idx_name" already exists
   ```
   **Solution**: Use `CREATE INDEX IF NOT EXISTS`

3. **Foreign Key Constraints**
   ```
   Error: foreign key constraint fails
   ```
   **Solution**: Ensure referenced tables exist first

### Debug Migration Issues

```go
// Enable GORM debug mode
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
})

// Check migration status
if db.Migrator().HasTable(&models.User{}) {
    logger.Info("User table exists")
}

// Check column existence
if db.Migrator().HasColumn(&models.User{}, "email") {
    logger.Info("Email column exists")
}
```

## Docker Integration

### Database Container

The PostgreSQL container no longer needs init.sql:

```yaml
foundation-db:
  image: postgres:15-alpine
  volumes:
    - foundation_db_data:/var/lib/postgresql/data
    - ./database/backups:/backups
  # No longer mounting init.sql
```

### Application Startup

Migrations run automatically when the application starts:

1. Database connection established
2. Migration service initialized
3. Service-specific migrations run
4. Demo data inserted
5. Application ready

## Performance Considerations

### Migration Speed

- **GORM AutoMigrate**: Fast for small schema changes
- **Index Creation**: May take time for large tables
- **Demo Data**: Minimal impact with conflict handling

### Production Deployment

- **Zero Downtime**: GORM migrations are generally safe
- **Backup First**: Always backup before major deployments
- **Monitor Logs**: Watch migration logs for any issues

For more information, see the [Configuration Guide](configuration.md) and [Docker Setup Guide](docker-setup.md).
