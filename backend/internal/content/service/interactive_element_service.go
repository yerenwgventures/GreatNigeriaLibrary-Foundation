package service

import (
        "encoding/json"
        "time"

        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/content/models"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/content/repository"
)

// InteractiveElementService defines the interface for interactive element business logic
type InteractiveElementService interface {
        GetInteractiveElementsBySection(sectionID uint) ([]models.InteractiveElement, error)
        GetInteractiveElementByID(id uint) (*models.InteractiveElement, error)
        CreateQuiz(sectionID uint, title, description string, content *models.QuizContent, completionType string, pointsValue int, requiredStatus bool) (*models.InteractiveElement, error)
        CreateReflection(sectionID uint, title, description string, content *models.ReflectionContent, completionType string, pointsValue int, requiredStatus bool) (*models.InteractiveElement, error)
        CreateCallToAction(sectionID uint, title, description string, content *models.CallToActionContent, completionType string, pointsValue int, requiredStatus bool) (*models.InteractiveElement, error)
        CreateDiscussionPrompt(sectionID uint, title, description string, content *models.DiscussionPromptContent, completionType string, pointsValue int, requiredStatus bool) (*models.InteractiveElement, error)
        UpdateInteractiveElement(element *models.InteractiveElement) error
        DeleteInteractiveElement(id uint) error
        
        // Response handling
        SubmitQuizResponse(userID, elementID uint, answers []QuizAnswer) (*models.InteractiveElementResponse, error)
        SubmitReflectionResponse(userID, elementID uint, response string) (*models.InteractiveElementResponse, error)
        SubmitCallToActionResponse(userID, elementID uint, actionType, actionData string) (*models.InteractiveElementResponse, error)
        SubmitDiscussionResponse(userID, elementID uint, response string, topicID uint) (*models.InteractiveElementResponse, error)
        GetUserResponsesForElement(userID, elementID uint) ([]models.InteractiveElementResponse, error)
        GetUserProgress(userID, bookID uint) (*models.UserInteractiveElementProgress, error)
}

// QuizAnswer represents a user's answer to a quiz question
type QuizAnswer struct {
        QuestionID uint        `json:"questionId"`
        Answer     interface{} `json:"answer"` // String or array of strings
        TimeSpent  int         `json:"timeSpent"`
}

// InteractiveElementServiceImpl implements the InteractiveElementService interface
type InteractiveElementServiceImpl struct {
        elementRepo repository.InteractiveElementRepository
        bookRepo    repository.BookRepository
        pointsRepo  repository.PointsRepository
}

// NewInteractiveElementService creates a new interactive element service instance
func NewInteractiveElementService(
        elementRepo repository.InteractiveElementRepository,
        bookRepo repository.BookRepository,
        pointsRepo repository.PointsRepository,
) InteractiveElementService {
        return &InteractiveElementServiceImpl{
                elementRepo: elementRepo,
                bookRepo:    bookRepo,
                pointsRepo:  pointsRepo,
        }
}

// GetInteractiveElementsBySection retrieves all interactive elements for a section
func (s *InteractiveElementServiceImpl) GetInteractiveElementsBySection(sectionID uint) ([]models.InteractiveElement, error) {
        return s.elementRepo.GetInteractiveElementsBySection(sectionID)
}

// GetInteractiveElementByID retrieves an interactive element by ID
func (s *InteractiveElementServiceImpl) GetInteractiveElementByID(id uint) (*models.InteractiveElement, error) {
        return s.elementRepo.GetInteractiveElementByID(id)
}

// CreateQuiz creates a new quiz interactive element
func (s *InteractiveElementServiceImpl) CreateQuiz(sectionID uint, title, description string, content *models.QuizContent, completionType string, pointsValue int, requiredStatus bool) (*models.InteractiveElement, error) {
        // Get the current max position for this section
        elements, err := s.elementRepo.GetInteractiveElementsBySection(sectionID)
        if err != nil {
                return nil, err
        }
        
        position := 0
        if len(elements) > 0 {
                position = elements[len(elements)-1].Position + 1
        }
        
        // Convert content to JSON
        contentBytes, err := json.Marshal(content)
        if err != nil {
                return nil, err
        }
        
        // Create new element
        element := &models.InteractiveElement{
                SectionID:      sectionID,
                Position:       position,
                Type:           models.QuizType,
                Title:          title,
                Description:    description,
                Content:        string(contentBytes),
                CompletionType: completionType,
                PointsValue:    pointsValue,
                RequiredStatus: requiredStatus,
        }
        
        // Save the element
        err = s.elementRepo.CreateInteractiveElement(element)
        if err != nil {
                return nil, err
        }
        
        return element, nil
}

// CreateReflection creates a new reflection interactive element
func (s *InteractiveElementServiceImpl) CreateReflection(sectionID uint, title, description string, content *models.ReflectionContent, completionType string, pointsValue int, requiredStatus bool) (*models.InteractiveElement, error) {
        // Get the current max position for this section
        elements, err := s.elementRepo.GetInteractiveElementsBySection(sectionID)
        if err != nil {
                return nil, err
        }
        
        position := 0
        if len(elements) > 0 {
                position = elements[len(elements)-1].Position + 1
        }
        
        // Convert content to JSON
        contentBytes, err := json.Marshal(content)
        if err != nil {
                return nil, err
        }
        
        // Create new element
        element := &models.InteractiveElement{
                SectionID:      sectionID,
                Position:       position,
                Type:           models.ReflectionType,
                Title:          title,
                Description:    description,
                Content:        string(contentBytes),
                CompletionType: completionType,
                PointsValue:    pointsValue,
                RequiredStatus: requiredStatus,
        }
        
        // Save the element
        err = s.elementRepo.CreateInteractiveElement(element)
        if err != nil {
                return nil, err
        }
        
        return element, nil
}

// CreateCallToAction creates a new call-to-action interactive element
func (s *InteractiveElementServiceImpl) CreateCallToAction(sectionID uint, title, description string, content *models.CallToActionContent, completionType string, pointsValue int, requiredStatus bool) (*models.InteractiveElement, error) {
        // Get the current max position for this section
        elements, err := s.elementRepo.GetInteractiveElementsBySection(sectionID)
        if err != nil {
                return nil, err
        }
        
        position := 0
        if len(elements) > 0 {
                position = elements[len(elements)-1].Position + 1
        }
        
        // Convert content to JSON
        contentBytes, err := json.Marshal(content)
        if err != nil {
                return nil, err
        }
        
        // Create new element
        element := &models.InteractiveElement{
                SectionID:      sectionID,
                Position:       position,
                Type:           models.CallToActionType,
                Title:          title,
                Description:    description,
                Content:        string(contentBytes),
                CompletionType: completionType,
                PointsValue:    pointsValue,
                RequiredStatus: requiredStatus,
        }
        
        // Save the element
        err = s.elementRepo.CreateInteractiveElement(element)
        if err != nil {
                return nil, err
        }
        
        return element, nil
}

// CreateDiscussionPrompt creates a new discussion prompt interactive element
func (s *InteractiveElementServiceImpl) CreateDiscussionPrompt(sectionID uint, title, description string, content *models.DiscussionPromptContent, completionType string, pointsValue int, requiredStatus bool) (*models.InteractiveElement, error) {
        // Get the current max position for this section
        elements, err := s.elementRepo.GetInteractiveElementsBySection(sectionID)
        if err != nil {
                return nil, err
        }
        
        position := 0
        if len(elements) > 0 {
                position = elements[len(elements)-1].Position + 1
        }
        
        // Convert content to JSON
        contentBytes, err := json.Marshal(content)
        if err != nil {
                return nil, err
        }
        
        // Create new element
        element := &models.InteractiveElement{
                SectionID:      sectionID,
                Position:       position,
                Type:           models.DiscussionPromptType,
                Title:          title,
                Description:    description,
                Content:        string(contentBytes),
                CompletionType: completionType,
                PointsValue:    pointsValue,
                RequiredStatus: requiredStatus,
        }
        
        // Save the element
        err = s.elementRepo.CreateInteractiveElement(element)
        if err != nil {
                return nil, err
        }
        
        return element, nil
}

// UpdateInteractiveElement updates an existing interactive element
func (s *InteractiveElementServiceImpl) UpdateInteractiveElement(element *models.InteractiveElement) error {
        return s.elementRepo.UpdateInteractiveElement(element)
}

// DeleteInteractiveElement deletes an interactive element
func (s *InteractiveElementServiceImpl) DeleteInteractiveElement(id uint) error {
        return s.elementRepo.DeleteInteractiveElement(id)
}

// SubmitQuizResponse submits a user's response to a quiz
func (s *InteractiveElementServiceImpl) SubmitQuizResponse(userID, elementID uint, answers []QuizAnswer) (*models.InteractiveElementResponse, error) {
        // Get the interactive element
        element, err := s.elementRepo.GetInteractiveElementByID(elementID)
        if err != nil {
                return nil, err
        }
        
        // Ensure element is a quiz
        if element.Type != models.QuizType {
                return nil, models.ErrInvalidElementType
        }
        
        // Get the quiz content
        quizContent, err := element.GetQuizContent()
        if err != nil {
                return nil, err
        }
        
        // Calculate score
        score := 0
        totalQuestions := len(quizContent.Questions)
        correctAnswers := 0
        
        // Convert answers to a map for easy lookup
        answerMap := make(map[uint]interface{})
        for _, answer := range answers {
                answerMap[answer.QuestionID] = answer.Answer
        }
        
        // Check each question
        for _, question := range quizContent.Questions {
                userAnswer, exists := answerMap[question.ID]
                if !exists {
                        continue // User didn't answer this question
                }
                
                // Compare answers - different comparison based on question type
                if isCorrectAnswer(userAnswer, question.CorrectAnswer, question.QuestionType) {
                        correctAnswers++
                }
        }
        
        // Calculate score as a percentage
        if totalQuestions > 0 {
                score = (correctAnswers * 100) / totalQuestions
        }
        
        // Determine if the user passed
        passed := score >= quizContent.PassScore
        
        // Determine points awarded
        pointsAwarded := 0
        if passed {
                pointsAwarded = element.PointsValue
        }
        
        // Calculate total time spent
        totalTimeSpent := 0
        for _, answer := range answers {
                totalTimeSpent += answer.TimeSpent
        }
        
        // Create response object
        response := &models.InteractiveElementResponse{
                UserID:               userID,
                InteractiveElementID: elementID,
                Response:             "", // Will be set below
                Score:                score,
                TimeSpent:            totalTimeSpent,
                CompletionStatus:     "completed",
                PointsAwarded:        pointsAwarded,
        }
        
        // Convert answers to JSON for storage
        responseData := map[string]interface{}{
                "answers":        answers,
                "score":          score,
                "correctAnswers": correctAnswers,
                "totalQuestions": totalQuestions,
                "passed":         passed,
                "timeSpent":      totalTimeSpent,
                "completedAt":    time.Now(),
        }
        
        responseBytes, err := json.Marshal(responseData)
        if err != nil {
                return nil, err
        }
        response.Response = string(responseBytes)
        
        // Save the response
        err = s.elementRepo.SaveElementResponse(response)
        if err != nil {
                return nil, err
        }
        
        // Award points if passed
        if passed && s.pointsRepo != nil {
                err = s.pointsRepo.AwardPoints(userID, "interactive_element", elementID, pointsAwarded, "Completed quiz: "+element.Title)
                if err != nil {
                        // Log the error but don't fail the response
                        // TODO: Add logging
                }
        }
        
        return response, nil
}

// SubmitReflectionResponse submits a user's response to a reflection exercise
func (s *InteractiveElementServiceImpl) SubmitReflectionResponse(userID, elementID uint, response string) (*models.InteractiveElementResponse, error) {
        // Get the interactive element
        element, err := s.elementRepo.GetInteractiveElementByID(elementID)
        if err != nil {
                return nil, err
        }
        
        // Ensure element is a reflection
        if element.Type != models.ReflectionType {
                return nil, models.ErrInvalidElementType
        }
        
        // Get the reflection content
        reflectionContent, err := element.GetReflectionContent()
        if err != nil {
                return nil, err
        }
        
        // Validate response length if required
        if reflectionContent.MinResponseLength > 0 && len(response) < reflectionContent.MinResponseLength {
                return nil, models.ErrInvalidContent
        }
        
        if reflectionContent.MaxResponseLength > 0 && len(response) > reflectionContent.MaxResponseLength {
                return nil, models.ErrInvalidContent
        }
        
        // Create response object
        responseObj := &models.InteractiveElementResponse{
                UserID:               userID,
                InteractiveElementID: elementID,
                Response:             "", // Will be set below
                Score:                100, // Reflections are usually not scored; give full points
                TimeSpent:            0,   // Not tracked for reflections
                CompletionStatus:     "completed",
                PointsAwarded:        element.PointsValue, // Full points for reflection
        }
        
        // Convert response to JSON for storage
        responseData := map[string]interface{}{
                "reflectionResponse": response,
                "completedAt":        time.Now(),
        }
        
        responseBytes, err := json.Marshal(responseData)
        if err != nil {
                return nil, err
        }
        responseObj.Response = string(responseBytes)
        
        // Save the response
        err = s.elementRepo.SaveElementResponse(responseObj)
        if err != nil {
                return nil, err
        }
        
        // Award points
        if s.pointsRepo != nil {
                err = s.pointsRepo.AwardPoints(userID, "interactive_element", elementID, element.PointsValue, "Completed reflection: "+element.Title)
                if err != nil {
                        // Log the error but don't fail the response
                        // TODO: Add logging
                }
        }
        
        return responseObj, nil
}

// SubmitCallToActionResponse submits a user's response to a call-to-action
func (s *InteractiveElementServiceImpl) SubmitCallToActionResponse(userID, elementID uint, actionType, actionData string) (*models.InteractiveElementResponse, error) {
        // Get the interactive element
        element, err := s.elementRepo.GetInteractiveElementByID(elementID)
        if err != nil {
                return nil, err
        }
        
        // Ensure element is a call-to-action
        if element.Type != models.CallToActionType {
                return nil, models.ErrInvalidElementType
        }
        
        // Get the call-to-action content
        ctaContent, err := element.GetCallToActionContent()
        if err != nil {
                return nil, err
        }
        
        // Validate action type
        if ctaContent.ActionType != actionType {
                return nil, models.ErrInvalidContent
        }
        
        // Create response object
        responseObj := &models.InteractiveElementResponse{
                UserID:               userID,
                InteractiveElementID: elementID,
                Response:             "", // Will be set below
                Score:                100, // Call-to-actions are not scored; give full points
                TimeSpent:            0,   // Not tracked for call-to-actions
                CompletionStatus:     "completed",
                PointsAwarded:        element.PointsValue, // Full points for call-to-action
        }
        
        // Convert response to JSON for storage
        responseData := map[string]interface{}{
                "actionType":  actionType,
                "actionData":  actionData,
                "completedAt": time.Now(),
        }
        
        responseBytes, err := json.Marshal(responseData)
        if err != nil {
                return nil, err
        }
        responseObj.Response = string(responseBytes)
        
        // Save the response
        err = s.elementRepo.SaveElementResponse(responseObj)
        if err != nil {
                return nil, err
        }
        
        // Award points
        if s.pointsRepo != nil {
                err = s.pointsRepo.AwardPoints(userID, "interactive_element", elementID, element.PointsValue, "Completed call-to-action: "+element.Title)
                if err != nil {
                        // Log the error but don't fail the response
                        // TODO: Add logging
                }
        }
        
        return responseObj, nil
}

// SubmitDiscussionResponse submits a user's response to a discussion prompt
func (s *InteractiveElementServiceImpl) SubmitDiscussionResponse(userID, elementID uint, response string, topicID uint) (*models.InteractiveElementResponse, error) {
        // Get the interactive element
        element, err := s.elementRepo.GetInteractiveElementByID(elementID)
        if err != nil {
                return nil, err
        }
        
        // Ensure element is a discussion prompt
        if element.Type != models.DiscussionPromptType {
                return nil, models.ErrInvalidElementType
        }
        
        // Create response object
        responseObj := &models.InteractiveElementResponse{
                UserID:               userID,
                InteractiveElementID: elementID,
                Response:             "", // Will be set below
                Score:                100, // Discussion prompts are not scored; give full points
                TimeSpent:            0,   // Not tracked for discussion prompts
                CompletionStatus:     "completed",
                PointsAwarded:        element.PointsValue, // Full points for discussion prompt
        }
        
        // Convert response to JSON for storage
        responseData := map[string]interface{}{
                "discussionResponse": response,
                "topicID":            topicID,
                "completedAt":        time.Now(),
        }
        
        responseBytes, err := json.Marshal(responseData)
        if err != nil {
                return nil, err
        }
        responseObj.Response = string(responseBytes)
        
        // Save the response
        err = s.elementRepo.SaveElementResponse(responseObj)
        if err != nil {
                return nil, err
        }
        
        // Award points
        if s.pointsRepo != nil {
                err = s.pointsRepo.AwardPoints(userID, "interactive_element", elementID, element.PointsValue, "Participated in discussion: "+element.Title)
                if err != nil {
                        // Log the error but don't fail the response
                        // TODO: Add logging
                }
        }
        
        return responseObj, nil
}

// GetUserResponsesForElement retrieves all of a user's responses to an interactive element
func (s *InteractiveElementServiceImpl) GetUserResponsesForElement(userID, elementID uint) ([]models.InteractiveElementResponse, error) {
        return s.elementRepo.GetUserResponsesForElement(userID, elementID)
}

// GetUserProgress retrieves a user's progress with interactive elements for a book
func (s *InteractiveElementServiceImpl) GetUserProgress(userID, bookID uint) (*models.UserInteractiveElementProgress, error) {
        return s.elementRepo.GetUserProgress(userID, bookID)
}

// Helper function to compare answers based on question type
func isCorrectAnswer(userAnswer, correctAnswer interface{}, questionType string) bool {
        switch questionType {
        case "multiple-choice":
                // For multiple-choice questions, compare strings
                return compareStringAnswers(userAnswer, correctAnswer)
        case "multiple-answer":
                // For multiple-answer questions, compare arrays of strings
                return compareArrayAnswers(userAnswer, correctAnswer)
        case "true-false":
                // For true-false questions, compare booleans as strings
                return compareStringAnswers(userAnswer, correctAnswer)
        case "fill-blank":
                // For fill-in-the-blank questions, compare strings with case insensitivity
                return compareStringAnswersCaseInsensitive(userAnswer, correctAnswer)
        case "short-answer":
                // For short-answer questions, compare strings with case insensitivity
                return compareStringAnswersCaseInsensitive(userAnswer, correctAnswer)
        default:
                // Default to simple string comparison
                return compareStringAnswers(userAnswer, correctAnswer)
        }
}

// Helper functions for answer comparison
func compareStringAnswers(userAnswer, correctAnswer interface{}) bool {
        // Convert to strings for comparison
        userStr, userOk := userAnswer.(string)
        correctStr, correctOk := correctAnswer.(string)
        
        if !userOk || !correctOk {
                return false
        }
        
        return userStr == correctStr
}

func compareStringAnswersCaseInsensitive(userAnswer, correctAnswer interface{}) bool {
        // Convert to strings for comparison
        userStr, userOk := userAnswer.(string)
        correctStr, correctOk := correctAnswer.(string)
        
        if !userOk || !correctOk {
                return false
        }
        
        // TODO: Use strings.EqualFold for case-insensitive comparison when Go implementation is ready
        return userStr == correctStr
}

func compareArrayAnswers(userAnswer, correctAnswer interface{}) bool {
        // Convert to arrays for comparison
        userArray, userOk := userAnswer.([]interface{})
        correctArray, correctOk := correctAnswer.([]interface{})
        
        if !userOk || !correctOk {
                return false
        }
        
        // Check that arrays have the same length
        if len(userArray) != len(correctArray) {
                return false
        }
        
        // Convert both arrays to maps for easier comparison
        userMap := make(map[string]bool)
        correctMap := make(map[string]bool)
        
        for _, v := range userArray {
                str, ok := v.(string)
                if !ok {
                        return false
                }
                userMap[str] = true
        }
        
        for _, v := range correctArray {
                str, ok := v.(string)
                if !ok {
                        return false
                }
                correctMap[str] = true
        }
        
        // Check that all items in correctMap are in userMap
        for k := range correctMap {
                if !userMap[k] {
                        return false
                }
        }
        
        // Check that all items in userMap are in correctMap
        for k := range userMap {
                if !correctMap[k] {
                        return false
                }
        }
        
        return true
}