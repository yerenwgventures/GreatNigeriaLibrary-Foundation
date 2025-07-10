package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/content/models"
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/content/repository"
)

// MediaGenerator defines the interface for generating media from content
type MediaGenerator interface {
	GenerateAudioFromText(sectionID uint) (string, int, error)
	GeneratePhotoCollection(sectionID uint) ([]string, error)
	GenerateVideoSlideshow(sectionID uint) (string, int, error)
	GeneratePDF(sectionID uint) (string, int, error)
	GetShareableLink(sectionID uint, mediaType string) (string, string, error)
}

// MediaGeneratorImpl implements the MediaGenerator interface
type MediaGeneratorImpl struct {
	bookRepo        repository.BookRepository
	contentRenderer ContentRenderer
	baseURL         string
	mediaBasePath   string
}

// NewMediaGenerator creates a new media generator instance
func NewMediaGenerator(
	bookRepo repository.BookRepository,
	contentRenderer ContentRenderer,
	baseURL string,
	mediaBasePath string,
) MediaGenerator {
	return &MediaGeneratorImpl{
		bookRepo:        bookRepo,
		contentRenderer: contentRenderer,
		baseURL:         baseURL,
		mediaBasePath:   mediaBasePath,
	}
}

// GenerateAudioFromText generates an audio file from section text
func (g *MediaGeneratorImpl) GenerateAudioFromText(sectionID uint) (string, int, error) {
	// Get section content
	section, err := g.bookRepo.GetSectionByID(sectionID)
	if err != nil {
		return "", 0, err
	}

	// Generate a unique identifier for this content
	contentHash := generateContentHash(section.Content)
	audioFileName := fmt.Sprintf("audio_%s.mp3", contentHash)
	audioDir := filepath.Join(g.mediaBasePath, "audio", "generated")
	audioFilePath := filepath.Join(audioDir, audioFileName)
	audioURL := fmt.Sprintf("%s/static/media/audio/generated/%s", g.baseURL, audioFileName)

	// Check if audio file already exists
	if fileExists(audioFilePath) {
		// Get audio duration (in seconds)
		duration := 180 // Placeholder: In a real implementation, get actual duration
		return audioURL, duration, nil
	}

	// Ensure directory exists
	os.MkdirAll(audioDir, 0755)

	// In a real implementation, this would call a text-to-speech service
	// For this implementation, we'll create a placeholder audio file
	err = createPlaceholderAudioFile(audioFilePath)
	if err != nil {
		return "", 0, err
	}

	// Return the URL and a placeholder duration
	return audioURL, 180, nil
}

// GeneratePhotoCollection generates a collection of images related to the section content
func (g *MediaGeneratorImpl) GeneratePhotoCollection(sectionID uint) ([]string, error) {
	// Get section content
	section, err := g.bookRepo.GetSectionByID(sectionID)
	if err != nil {
		return nil, err
	}

	// Generate a unique identifier for this content
	contentHash := generateContentHash(section.Content)
	photoDir := filepath.Join(g.mediaBasePath, "images", "generated", contentHash)
	
	// Check if photo collection already exists
	if dirExists(photoDir) {
		return getExistingPhotoURLs(photoDir, g.baseURL)
	}

	// Ensure directory exists
	os.MkdirAll(photoDir, 0755)

	// Extract keywords from content for image generation
	keywords := extractKeywords(section.Content)
	
	// In a real implementation, this would generate or fetch relevant images
	// For this implementation, we'll create placeholder images
	var photoURLs []string
	for i, keyword := range keywords {
		if i >= 6 { // Limit to 6 images
			break
		}
		
		imageName := fmt.Sprintf("image_%d_%s.jpg", i, sanitizeFilename(keyword))
		imagePath := filepath.Join(photoDir, imageName)
		imageURL := fmt.Sprintf("%s/static/media/images/generated/%s/%s", g.baseURL, contentHash, imageName)
		
		err := createPlaceholderImage(imagePath)
		if err != nil {
			continue
		}
		
		photoURLs = append(photoURLs, imageURL)
	}
	
	return photoURLs, nil
}

// GenerateVideoSlideshow generates a video slideshow from section content
func (g *MediaGeneratorImpl) GenerateVideoSlideshow(sectionID uint) (string, int, error) {
	// Get section content
	section, err := g.bookRepo.GetSectionByID(sectionID)
	if err != nil {
		return "", 0, err
	}

	// Generate a unique identifier for this content
	contentHash := generateContentHash(section.Content)
	videoFileName := fmt.Sprintf("video_%s.mp4", contentHash)
	videoDir := filepath.Join(g.mediaBasePath, "video", "generated")
	videoFilePath := filepath.Join(videoDir, videoFileName)
	videoURL := fmt.Sprintf("%s/static/media/video/generated/%s", g.baseURL, videoFileName)

	// Check if video file already exists
	if fileExists(videoFilePath) {
		// Get video duration (in seconds)
		duration := 240 // Placeholder: In a real implementation, get actual duration
		return videoURL, duration, nil
	}

	// Ensure directory exists
	os.MkdirAll(videoDir, 0755)

	// In a real implementation, this would generate a slideshow video
	// For this implementation, we'll create a placeholder video file
	err = createPlaceholderVideoFile(videoFilePath)
	if err != nil {
		return "", 0, err
	}

	// Return the URL and a placeholder duration
	return videoURL, 240, nil
}

// GeneratePDF generates a PDF from section content
func (g *MediaGeneratorImpl) GeneratePDF(sectionID uint) (string, int, error) {
	// Get section content
	section, err := g.bookRepo.GetSectionByID(sectionID)
	if err != nil {
		return "", 0, err
	}

	// Generate a unique identifier for this content
	contentHash := generateContentHash(section.Content)
	pdfFileName := fmt.Sprintf("pdf_%s.pdf", contentHash)
	pdfDir := filepath.Join(g.mediaBasePath, "pdf", "generated")
	pdfFilePath := filepath.Join(pdfDir, pdfFileName)
	pdfURL := fmt.Sprintf("%s/static/media/pdf/generated/%s", g.baseURL, pdfFileName)

	// Check if PDF file already exists
	if fileExists(pdfFilePath) {
		// Get page count
		pageCount := 5 // Placeholder: In a real implementation, get actual page count
		return pdfURL, pageCount, nil
	}

	// Ensure directory exists
	os.MkdirAll(pdfDir, 0755)

	// In a real implementation, this would generate a PDF
	// For this implementation, we'll create a placeholder PDF file
	err = createPlaceholderPDFFile(pdfFilePath)
	if err != nil {
		return "", 0, err
	}

	// Return the URL and a placeholder page count
	return pdfURL, 5, nil
}

// GetShareableLink generates a shareable link for media content
func (g *MediaGeneratorImpl) GetShareableLink(sectionID uint, mediaType string) (string, string, error) {
	// Get section
	section, err := g.bookRepo.GetSectionByID(sectionID)
	if err != nil {
		return "", "", err
	}
	
	// Get chapter and book info for the title
	chapter, err := g.bookRepo.GetChapterByID(section.ChapterID)
	if err != nil {
		return "", "", err
	}
	
	book, err := g.bookRepo.GetBookByID(chapter.BookID)
	if err != nil {
		return "", "", err
	}
	
	// Generate a shareable link based on the media type
	var mediaURL string
	var shareableLink string
	
	switch mediaType {
	case "audio":
		audioURL, _, err := g.GenerateAudioFromText(sectionID)
		if err != nil {
			return "", "", err
		}
		mediaURL = audioURL
		shareableLink = fmt.Sprintf("%s/share/audio/%d", g.baseURL, sectionID)
	case "photo":
		photoURLs, err := g.GeneratePhotoCollection(sectionID)
		if err != nil || len(photoURLs) == 0 {
			return "", "", err
		}
		mediaURL = strings.Join(photoURLs, ",")
		shareableLink = fmt.Sprintf("%s/share/photos/%d", g.baseURL, sectionID)
	case "video":
		videoURL, _, err := g.GenerateVideoSlideshow(sectionID)
		if err != nil {
			return "", "", err
		}
		mediaURL = videoURL
		shareableLink = fmt.Sprintf("%s/share/video/%d", g.baseURL, sectionID)
	case "pdf":
		pdfURL, _, err := g.GeneratePDF(sectionID)
		if err != nil {
			return "", "", err
		}
		mediaURL = pdfURL
		shareableLink = fmt.Sprintf("%s/share/pdf/%d", g.baseURL, sectionID)
	default:
		return "", "", fmt.Errorf("invalid media type: %s", mediaType)
	}
	
	return shareableLink, mediaURL, nil
}

// Helper functions

// generateContentHash generates a hash from content for caching
func generateContentHash(content string) string {
	hash := md5.Sum([]byte(content))
	return hex.EncodeToString(hash[:])
}

// fileExists checks if a file exists
func fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// dirExists checks if a directory exists
func dirExists(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// extractKeywords extracts keywords from content
func extractKeywords(content string) []string {
	// Simple implementation: extract words longer than 5 characters
	words := regexp.MustCompile(`\b\w{5,}\b`).FindAllString(content, -1)
	
	// Deduplicate
	seen := make(map[string]bool)
	var keywords []string
	
	for _, word := range words {
		word = strings.ToLower(word)
		if !seen[word] {
			seen[word] = true
			keywords = append(keywords, word)
		}
	}
	
	return keywords
}

// sanitizeFilename sanitizes a string for use as a filename
func sanitizeFilename(s string) string {
	// Replace invalid characters with underscores
	return regexp.MustCompile(`[^a-zA-Z0-9_-]`).ReplaceAllString(s, "_")
}

// getExistingPhotoURLs gets URLs for existing photos
func getExistingPhotoURLs(photoDir string, baseURL string) ([]string, error) {
	var photoURLs []string
	
	err := filepath.Walk(photoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() && (strings.HasSuffix(info.Name(), ".jpg") || 
							 strings.HasSuffix(info.Name(), ".jpeg") || 
							 strings.HasSuffix(info.Name(), ".png")) {
			relPath, err := filepath.Rel(photoDir, path)
			if err != nil {
				return err
			}
			
			dirName := filepath.Base(photoDir)
			imageURL := fmt.Sprintf("%s/static/media/images/generated/%s/%s", baseURL, dirName, relPath)
			photoURLs = append(photoURLs, imageURL)
		}
		
		return nil
	})
	
	return photoURLs, err
}

// Placeholder file creation functions
// In a real implementation, these would be replaced with actual media generation

func createPlaceholderAudioFile(filePath string) error {
	// Create an empty file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	// Write some placeholder content
	_, err = file.WriteString("This is a placeholder audio file")
	return err
}

func createPlaceholderImage(filePath string) error {
	// Create an empty file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	// Write some placeholder content
	_, err = file.WriteString("This is a placeholder image file")
	return err
}

func createPlaceholderVideoFile(filePath string) error {
	// Create an empty file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	// Write some placeholder content
	_, err = file.WriteString("This is a placeholder video file")
	return err
}

func createPlaceholderPDFFile(filePath string) error {
	// Create an empty file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	// Write some placeholder content
	_, err = file.WriteString("%PDF-1.4\nThis is a placeholder PDF file\n%%EOF")
	return err
}
