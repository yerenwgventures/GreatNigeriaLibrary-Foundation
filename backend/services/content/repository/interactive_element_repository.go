package repository

import (
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/services/content/models"
        "gorm.io/gorm"
)

// InteractiveElementRepository defines the interface for interactive element operations
type InteractiveElementRepository interface {
        GetInteractiveElementsBySection(sectionID uint) ([]models.InteractiveElement, error)
        GetInteractiveElementByID(id uint) (*models.InteractiveElement, error)
        CreateInteractiveElement(element *models.InteractiveElement) error
        UpdateInteractiveElement(element *models.InteractiveElement) error
        DeleteInteractiveElement(id uint) error
        
        // Response related methods
        SaveElementResponse(response *models.InteractiveElementResponse) error
        GetUserResponsesForElement(userID, elementID uint) ([]models.InteractiveElementResponse, error)
        GetLatestUserResponseForElement(userID, elementID uint) (*models.InteractiveElementResponse, error)
        GetUserProgress(userID, bookID uint) (*models.UserInteractiveElementProgress, error)
        UpdateUserProgress(progress *models.UserInteractiveElementProgress) error
}

// GormInteractiveElementRepository implements the InteractiveElementRepository interface with GORM
type GormInteractiveElementRepository struct {
        db *gorm.DB
}

// NewGormInteractiveElementRepository creates a new interactive element repository instance
func NewGormInteractiveElementRepository(db *gorm.DB) *GormInteractiveElementRepository {
        return &GormInteractiveElementRepository{db: db}
}

// GetInteractiveElementsBySection retrieves all interactive elements for a section
func (r *GormInteractiveElementRepository) GetInteractiveElementsBySection(sectionID uint) ([]models.InteractiveElement, error) {
        var elements []models.InteractiveElement
        
        result := r.db.Where("section_id = ?", sectionID).
                Order("position ASC").
                Find(&elements)
        
        return elements, result.Error
}

// GetInteractiveElementByID retrieves an interactive element by ID
func (r *GormInteractiveElementRepository) GetInteractiveElementByID(id uint) (*models.InteractiveElement, error) {
        var element models.InteractiveElement
        
        result := r.db.First(&element, id)
        if result.Error != nil {
                return nil, result.Error
        }
        
        return &element, nil
}

// CreateInteractiveElement creates a new interactive element
func (r *GormInteractiveElementRepository) CreateInteractiveElement(element *models.InteractiveElement) error {
        return r.db.Create(element).Error
}

// UpdateInteractiveElement updates an existing interactive element
func (r *GormInteractiveElementRepository) UpdateInteractiveElement(element *models.InteractiveElement) error {
        return r.db.Save(element).Error
}

// DeleteInteractiveElement deletes an interactive element
func (r *GormInteractiveElementRepository) DeleteInteractiveElement(id uint) error {
        // Use transaction to handle deleting the element and its responses
        return r.db.Transaction(func(tx *gorm.DB) error {
                // Delete all responses to this element
                if err := tx.Where("interactive_element_id = ?", id).Delete(&models.InteractiveElementResponse{}).Error; err != nil {
                        return err
                }
                
                // Delete the element itself
                return tx.Delete(&models.InteractiveElement{}, id).Error
        })
}

// SaveElementResponse saves a user's response to an interactive element
func (r *GormInteractiveElementRepository) SaveElementResponse(response *models.InteractiveElementResponse) error {
        // Use transaction to save the response and update the user progress
        return r.db.Transaction(func(tx *gorm.DB) error {
                // Save the response
                if err := tx.Create(response).Error; err != nil {
                        return err
                }
                
                // Get the interactive element to find the book ID
                var element models.InteractiveElement
                if err := tx.First(&element, response.InteractiveElementID).Error; err != nil {
                        return err
                }
                
                // Get the section to find the book ID
                var section models.Section
                if err := tx.First(&section, element.SectionID).Error; err != nil {
                        return err
                }
                
                // Get the chapter to find the book ID
                var chapter models.Chapter
                if err := tx.First(&chapter, section.ChapterID).Error; err != nil {
                        return err
                }
                
                bookID := chapter.BookID
                
                // Now update the user progress
                var progress models.UserInteractiveElementProgress
                result := tx.Where("user_id = ? AND book_id = ?", response.UserID, bookID).First(&progress)
                
                // If progress doesn't exist, create it
                if result.Error != nil {
                        if result.Error == gorm.ErrRecordNotFound {
                                // Count total elements in the book
                                var totalElementsCount int64
                                tx.Model(&models.InteractiveElement{}).
                                        Joins("JOIN sections ON sections.id = interactive_elements.section_id").
                                        Joins("JOIN chapters ON chapters.id = sections.chapter_id").
                                        Where("chapters.book_id = ?", bookID).
                                        Count(&totalElementsCount)
                                
                                // Count total points available
                                var totalPointsAvailable int
                                tx.Model(&models.InteractiveElement{}).
                                        Joins("JOIN sections ON sections.id = interactive_elements.section_id").
                                        Joins("JOIN chapters ON chapters.id = sections.chapter_id").
                                        Where("chapters.book_id = ?", bookID).
                                        Select("SUM(points_value)").
                                        Scan(&totalPointsAvailable)
                                
                                // Create new progress entry
                                progress = models.UserInteractiveElementProgress{
                                        UserID:               response.UserID,
                                        BookID:               bookID,
                                        TotalElements:        int(totalElementsCount),
                                        CompletedElements:    1,                              // First completion
                                        CompletionPercentage: int(1 * 100 / totalElementsCount), // Calculate percentage
                                        TotalPointsAvailable: totalPointsAvailable,
                                        TotalPointsEarned:    response.PointsAwarded,
                                        AvgScorePercentage:   response.Score, // Initial score is also the average
                                        RequiredCompleted:    false,          // Will be calculated below
                                }
                                
                                // Check if all required elements are completed
                                var requiredElementsCount int64
                                tx.Model(&models.InteractiveElement{}).
                                        Joins("JOIN sections ON sections.id = interactive_elements.section_id").
                                        Joins("JOIN chapters ON chapters.id = sections.chapter_id").
                                        Where("chapters.book_id = ? AND required_status = ?", bookID, true).
                                        Count(&requiredElementsCount)
                                
                                // Count completed required elements
                                var completedRequiredCount int64
                                tx.Model(&models.InteractiveElementResponse{}).
                                        Joins("JOIN interactive_elements ON interactive_elements.id = interactive_element_responses.interactive_element_id").
                                        Joins("JOIN sections ON sections.id = interactive_elements.section_id").
                                        Joins("JOIN chapters ON chapters.id = sections.chapter_id").
                                        Where("chapters.book_id = ? AND interactive_elements.required_status = ? AND interactive_element_responses.user_id = ? AND completion_status = ?",
                                                bookID, true, response.UserID, "completed").
                                        Count(&completedRequiredCount)
                                
                                // Check if all required elements are completed
                                progress.RequiredCompleted = completedRequiredCount >= requiredElementsCount
                                
                                // Save the new progress
                                return tx.Create(&progress).Error
                        }
                        return result.Error
                }
                
                // Update existing progress
                progress.CompletedElements++
                progress.CompletionPercentage = int(progress.CompletedElements * 100 / progress.TotalElements)
                progress.TotalPointsEarned += response.PointsAwarded
                
                // Recalculate average score
                var totalScore int
                var responseCount int64
                tx.Model(&models.InteractiveElementResponse{}).
                        Joins("JOIN interactive_elements ON interactive_elements.id = interactive_element_responses.interactive_element_id").
                        Joins("JOIN sections ON sections.id = interactive_elements.section_id").
                        Joins("JOIN chapters ON chapters.id = sections.chapter_id").
                        Where("chapters.book_id = ? AND interactive_element_responses.user_id = ?", bookID, response.UserID).
                        Select("SUM(score)").
                        Scan(&totalScore)
                
                tx.Model(&models.InteractiveElementResponse{}).
                        Joins("JOIN interactive_elements ON interactive_elements.id = interactive_element_responses.interactive_element_id").
                        Joins("JOIN sections ON sections.id = interactive_elements.section_id").
                        Joins("JOIN chapters ON chapters.id = sections.chapter_id").
                        Where("chapters.book_id = ? AND interactive_element_responses.user_id = ?", bookID, response.UserID).
                        Count(&responseCount)
                
                if responseCount > 0 {
                        progress.AvgScorePercentage = totalScore / int(responseCount)
                }
                
                // Check if all required elements are completed
                var requiredElementsCount int64
                tx.Model(&models.InteractiveElement{}).
                        Joins("JOIN sections ON sections.id = interactive_elements.section_id").
                        Joins("JOIN chapters ON chapters.id = sections.chapter_id").
                        Where("chapters.book_id = ? AND required_status = ?", bookID, true).
                        Count(&requiredElementsCount)
                
                // Count completed required elements
                var completedRequiredCount int64
                tx.Model(&models.InteractiveElementResponse{}).
                        Joins("JOIN interactive_elements ON interactive_elements.id = interactive_element_responses.interactive_element_id").
                        Joins("JOIN sections ON sections.id = interactive_elements.section_id").
                        Joins("JOIN chapters ON chapters.id = sections.chapter_id").
                        Where("chapters.book_id = ? AND interactive_elements.required_status = ? AND interactive_element_responses.user_id = ? AND completion_status = ?",
                                bookID, true, response.UserID, "completed").
                        Count(&completedRequiredCount)
                
                // Check if all required elements are completed
                progress.RequiredCompleted = completedRequiredCount >= requiredElementsCount
                
                // Save the updated progress
                return tx.Save(&progress).Error
        })
}

// GetUserResponsesForElement retrieves all of a user's responses to an interactive element
func (r *GormInteractiveElementRepository) GetUserResponsesForElement(userID, elementID uint) ([]models.InteractiveElementResponse, error) {
        var responses []models.InteractiveElementResponse
        
        result := r.db.Where("user_id = ? AND interactive_element_id = ?", userID, elementID).
                Order("created_at DESC").
                Find(&responses)
        
        return responses, result.Error
}

// GetLatestUserResponseForElement retrieves a user's latest response to an interactive element
func (r *GormInteractiveElementRepository) GetLatestUserResponseForElement(userID, elementID uint) (*models.InteractiveElementResponse, error) {
        var response models.InteractiveElementResponse
        
        result := r.db.Where("user_id = ? AND interactive_element_id = ?", userID, elementID).
                Order("created_at DESC").
                First(&response)
        
        if result.Error != nil {
                return nil, result.Error
        }
        
        return &response, nil
}

// GetUserProgress retrieves a user's progress with interactive elements for a book
func (r *GormInteractiveElementRepository) GetUserProgress(userID, bookID uint) (*models.UserInteractiveElementProgress, error) {
        var progress models.UserInteractiveElementProgress
        
        result := r.db.Where("user_id = ? AND book_id = ?", userID, bookID).
                First(&progress)
        
        if result.Error != nil {
                return nil, result.Error
        }
        
        return &progress, nil
}

// UpdateUserProgress updates a user's progress with interactive elements
func (r *GormInteractiveElementRepository) UpdateUserProgress(progress *models.UserInteractiveElementProgress) error {
        return r.db.Save(progress).Error
}