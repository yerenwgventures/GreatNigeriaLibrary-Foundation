package database

import (
	"fmt"

	"gorm.io/gorm"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/models"
)

// MigrationService handles database migrations for all services
type MigrationService struct {
	db     *gorm.DB
	logger Logger
}

// Logger interface for migration logging
type Logger interface {
	Info(msg string)
	Error(msg string)
	WithField(key string, value interface{}) Logger
}

// NewMigrationService creates a new migration service
func NewMigrationService(db *gorm.DB, logger Logger) *MigrationService {
	return &MigrationService{
		db:     db,
		logger: logger,
	}
}

// MigrateAuthService runs migrations for authentication service
func (m *MigrationService) MigrateAuthService() error {
	m.logger.Info("Running auth service migrations...")
	
	models := []interface{}{
		&models.User{},
		&models.Session{},
		&models.PasswordResetToken{},
		&models.EmailVerificationToken{},
		&models.UserTrustLevel{},
	}
	
	if err := m.db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("auth service migration failed: %w", err)
	}
	
	// Create indexes for auth service
	if err := m.createAuthIndexes(); err != nil {
		return fmt.Errorf("failed to create auth indexes: %w", err)
	}
	
	m.logger.Info("Auth service migrations completed successfully")
	return nil
}

// MigrateContentService runs migrations for content service
func (m *MigrationService) MigrateContentService() error {
	m.logger.Info("Running content service migrations...")
	
	// Import content models when they exist
	// For now, we'll create basic content tables that match init.sql
	if err := m.createContentTables(); err != nil {
		return fmt.Errorf("content service migration failed: %w", err)
	}
	
	m.logger.Info("Content service migrations completed successfully")
	return nil
}

// MigrateDiscussionService runs migrations for discussion service
func (m *MigrationService) MigrateDiscussionService() error {
	m.logger.Info("Running discussion service migrations...")
	
	// Create discussion tables that match init.sql
	if err := m.createDiscussionTables(); err != nil {
		return fmt.Errorf("discussion service migration failed: %w", err)
	}
	
	m.logger.Info("Discussion service migrations completed successfully")
	return nil
}

// MigrateGroupsService runs migrations for groups service
func (m *MigrationService) MigrateGroupsService() error {
	m.logger.Info("Running groups service migrations...")
	
	// Create groups tables that match init.sql
	if err := m.createGroupsTables(); err != nil {
		return fmt.Errorf("groups service migration failed: %w", err)
	}
	
	m.logger.Info("Groups service migrations completed successfully")
	return nil
}

// MigrateFoundation runs all foundation migrations
func (m *MigrationService) MigrateFoundation() error {
	m.logger.Info("Running foundation migrations...")
	
	// Run all service migrations
	if err := m.MigrateAuthService(); err != nil {
		return err
	}
	
	if err := m.MigrateContentService(); err != nil {
		return err
	}
	
	if err := m.MigrateDiscussionService(); err != nil {
		return err
	}
	
	if err := m.MigrateGroupsService(); err != nil {
		return err
	}
	
	// Insert demo data
	if err := m.insertDemoData(); err != nil {
		return fmt.Errorf("failed to insert demo data: %w", err)
	}
	
	m.logger.Info("Foundation migrations completed successfully")
	return nil
}

// createAuthIndexes creates indexes for auth service
func (m *MigrationService) createAuthIndexes() error {
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)",
		"CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)",
		"CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at)",
		"CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_user_id ON password_reset_tokens(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_expires_at ON password_reset_tokens(expires_at)",
		"CREATE INDEX IF NOT EXISTS idx_email_verification_tokens_user_id ON email_verification_tokens(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_user_trust_levels_user_id ON user_trust_levels(user_id)",
	}
	
	for _, index := range indexes {
		if err := m.db.Exec(index).Error; err != nil {
			m.logger.WithField("index", index).Error("Failed to create index")
			// Continue with other indexes even if one fails
		}
	}
	
	return nil
}

// createContentTables creates content tables to match init.sql
func (m *MigrationService) createContentTables() error {
	// Create content_categories table
	if err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS content_categories (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			slug VARCHAR(255) UNIQUE NOT NULL,
			description TEXT,
			parent_id UUID REFERENCES content_categories(id),
			sort_order INTEGER DEFAULT 0,
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`).Error; err != nil {
		return err
	}
	
	// Create demo_content table
	if err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS demo_content (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			title VARCHAR(255) NOT NULL,
			slug VARCHAR(255) UNIQUE NOT NULL,
			content TEXT NOT NULL,
			content_type VARCHAR(50) DEFAULT 'article',
			category_id UUID REFERENCES content_categories(id),
			author_id BIGINT REFERENCES users(id),
			status VARCHAR(20) DEFAULT 'published',
			featured BOOLEAN DEFAULT FALSE,
			view_count INTEGER DEFAULT 0,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`).Error; err != nil {
		return err
	}
	
	// Create indexes
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_demo_content_category ON demo_content(category_id)",
		"CREATE INDEX IF NOT EXISTS idx_demo_content_author ON demo_content(author_id)",
		"CREATE INDEX IF NOT EXISTS idx_demo_content_status ON demo_content(status)",
	}
	
	for _, index := range indexes {
		if err := m.db.Exec(index).Error; err != nil {
			m.logger.WithField("index", index).Error("Failed to create content index")
		}
	}
	
	return nil
}

// createDiscussionTables creates discussion tables to match init.sql
func (m *MigrationService) createDiscussionTables() error {
	// Create forum_categories table
	if err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS forum_categories (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			description TEXT,
			sort_order INTEGER DEFAULT 0,
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`).Error; err != nil {
		return err
	}
	
	// Create forum_topics table
	if err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS forum_topics (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			category_id UUID REFERENCES forum_categories(id),
			user_id BIGINT REFERENCES users(id),
			title VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			is_pinned BOOLEAN DEFAULT FALSE,
			is_locked BOOLEAN DEFAULT FALSE,
			view_count INTEGER DEFAULT 0,
			reply_count INTEGER DEFAULT 0,
			last_reply_at TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`).Error; err != nil {
		return err
	}
	
	// Create forum_replies table
	if err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS forum_replies (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			topic_id UUID REFERENCES forum_topics(id) ON DELETE CASCADE,
			user_id BIGINT REFERENCES users(id),
			content TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`).Error; err != nil {
		return err
	}
	
	// Create indexes
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_forum_topics_category ON forum_topics(category_id)",
		"CREATE INDEX IF NOT EXISTS idx_forum_topics_user ON forum_topics(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_forum_replies_topic ON forum_replies(topic_id)",
		"CREATE INDEX IF NOT EXISTS idx_forum_replies_user ON forum_replies(user_id)",
	}
	
	for _, index := range indexes {
		if err := m.db.Exec(index).Error; err != nil {
			m.logger.WithField("index", index).Error("Failed to create discussion index")
		}
	}
	
	return nil
}

// createGroupsTables creates groups tables to match init.sql
func (m *MigrationService) createGroupsTables() error {
	// Create user_groups table
	if err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS user_groups (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			description TEXT,
			creator_id BIGINT REFERENCES users(id),
			is_public BOOLEAN DEFAULT TRUE,
			member_count INTEGER DEFAULT 0,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`).Error; err != nil {
		return err
	}
	
	// Create group_members table
	if err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS group_members (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			group_id UUID REFERENCES user_groups(id) ON DELETE CASCADE,
			user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
			role VARCHAR(50) DEFAULT 'member',
			joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			UNIQUE(group_id, user_id)
		)
	`).Error; err != nil {
		return err
	}
	
	// Create indexes
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_group_members_group ON group_members(group_id)",
		"CREATE INDEX IF NOT EXISTS idx_group_members_user ON group_members(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_user_groups_creator ON user_groups(creator_id)",
	}
	
	for _, index := range indexes {
		if err := m.db.Exec(index).Error; err != nil {
			m.logger.WithField("index", index).Error("Failed to create groups index")
		}
	}
	
	return nil
}

// insertDemoData inserts demo data to match init.sql
func (m *MigrationService) insertDemoData() error {
	m.logger.Info("Inserting demo data...")
	
	// Insert content categories
	if err := m.db.Exec(`
		INSERT INTO content_categories (name, slug, description) VALUES
		('Platform Guide', 'platform-guide', 'How to use the Great Nigeria Library platform'),
		('Educational Content', 'educational', 'Educational materials and resources'),
		('Community Guidelines', 'community', 'Community rules and guidelines')
		ON CONFLICT (slug) DO NOTHING
	`).Error; err != nil {
		return fmt.Errorf("failed to insert content categories: %w", err)
	}
	
	// Insert forum categories
	if err := m.db.Exec(`
		INSERT INTO forum_categories (name, description, sort_order) VALUES
		('General Discussion', 'General topics and discussions', 1),
		('Platform Help', 'Help and support for using the platform', 2),
		('Feature Requests', 'Suggest new features for the platform', 3)
		ON CONFLICT DO NOTHING
	`).Error; err != nil {
		return fmt.Errorf("failed to insert forum categories: %w", err)
	}
	
	m.logger.Info("Demo data inserted successfully")
	return nil
}
