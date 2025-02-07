package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CheckType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif":
		return "image"
	case ".mp4", ".avi", ".mov", ".mkv":
		return "video"
	case ".pdf", ".doc", ".docx":
		return "document"
	default:
		return "unknown"
	}
}

type ObjectHandle struct {
	ObjectAttrs struct {
		Metadata map[string]string
	}
}

func (o *ObjectHandle) NewWriter(ctx context.Context) *ObjectWriter {
	return &ObjectWriter{ObjectAttrs: &o.ObjectAttrs}
}

// type ObjectWriter struct {
// 	ObjectAttrs *struct {
// 		Metadata map[string]string
// 	}
// }

// func (o *ObjectWriter) Close() error {
// 	return nil
// }

func initializeObjectHandle(filename string) *ObjectHandle {
	return &ObjectHandle{}
}

type ObjectWriter struct {
	ObjectAttrs *struct {
		Metadata map[string]string
	}
	buffer []byte // Simulate storage buffer
}

// Implement the Write method for ObjectWriter
func (o *ObjectWriter) Write(p []byte) (n int, err error) {
	o.buffer = append(o.buffer, p...) // Simulate writing data to a buffer
	return len(p), nil
}

func (o *ObjectWriter) Close() error {
	// Simulate closing the writer
	fmt.Println("Data written:", string(o.buffer))
	return nil
}

// UploadHandler handles file uploads and saves them locally.
// @Summary Upload an Image
// @Description Upload an image file to the server and store it locally.
// @Tags File
// @Accept multipart/form-data
// @Produce application/json
// @Param file formData file true "Image File"
// @Success 200 {object} Response "File uploaded successfully"
// @Failure 400 {object} Response "Bad Request"
// @Failure 500 {object} Response "Server Error"
// @Router /uploadd [post]
func (h *handler) UploadHandler(c *gin.Context) {
	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is uploaded"})
		return
	}

	// Create a directory to store the files
	storageDir := "./uploads"
	if err := os.MkdirAll(storageDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create storage directory"})
		return
	}
	uuid := uuid.New()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate UUID"})
		return
	}
	filename := uuid.String() + filepath.Ext(file.Filename)
	// Generate a unique filename
	// timestamp := time.Now().Unix()
	// ext := filepath.Ext(file.Filename)
	// filename := fmt.Sprintf("%d%s", timestamp, ext)

	// Save the file to the disk
	savePath := filepath.Join(storageDir, filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the file"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"filename": filename,
		"path":     savePath,
	})
}

// ListImagesHandler lists all available images
// @Summary List all uploaded images except Profile Images
// @Description Get a list of all images that have been uploaded
// @Tags File
// @Produce json
// @Success 200 {object} map[string]interface{} "List of images"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /images [get]
func (h *handler) ListImagesHandler(c *gin.Context) {
	uploadDir := "./uploads/profile_images"

	files, err := os.ReadDir(uploadDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read images directory"})
		return
	}

	var images []string

	// Loop through the files and get the filenames of image files
	for _, file := range files {
		if !file.IsDir() {
			ext := strings.ToLower(filepath.Ext(file.Name()))
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".pdf" {
				images = append(images, file.Name())
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"images": images})
}

// GetImageHandler serves the requested image dynamically
// @Summary Get a specific image
// @Description Get an image by its filename
// @Tags File
// @Accept json
// @Produce image/png, image/jpeg, image/gif
// @Param filename path string true "Filename of the image"
// @Success 200 {file} file "The image file"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 404 {object} map[string]interface{} "Image Not Found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /image/{filename} [get]
func (h *handler) GetImageHandler(c *gin.Context) {
	// Get the filename from the URL parameter
	filename := c.Param("filename")
	if filename == "" {
		// Handle missing filename
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	// Use the absolute path to the uploads directory
	// Make sure to update the path according to your system's structure
	imagePath := "./uploads/" + filename

	// Check if the file exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		// If the file doesn't exist, return a 404 error
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	// Set the proper content type based on the file extension
	ext := filepath.Ext(filename)
	var contentType string
	switch ext {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	case ".pdf":
		contentType = "image/pdf"
	case ".docx":
		contentType = "image/docx"
	default:
		// If the extension is unsupported, return an error
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Unsupported image type"})
		return
	}

	// Set the response header to the content type
	c.Header("Content-Type", contentType)

	// Serve the image file
	c.File(imagePath)
}

// UploadHandlerProfile handles file uploads and saves them locally.
// @Summary Upload an Image
// @Description Upload an image file to the server and store it locally.
// @Tags File
// @Accept multipart/form-data
// @Produce application/json
// @Param file formData file true "Image File"
// @Success 200 {object} Response "File uploaded successfully"
// @Failure 400 {object} Response "Bad Request"
// @Failure 500 {object} Response "Server Error"
// @Router /upload_profile [post]
func (h *handler) UploadHandlerProfile(c *gin.Context) {
	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is uploaded"})
		return
	}

	// Create a directory to store the files
	storageDir := "./uploads/profile_images"
	if err := os.MkdirAll(storageDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create storage directory"})
		return
	}
	uuid := uuid.New()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate UUID"})
		return
	}
	filename := uuid.String() + filepath.Ext(file.Filename)
	// Generate a unique filename
	// timestamp := time.Now().Unix()
	// ext := filepath.Ext(file.Filename)
	// filename := fmt.Sprintf("%d%s", timestamp, ext)

	// Save the file to the disk
	savePath := filepath.Join(storageDir, filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the file"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"filename": filename,
		"path":     savePath,
	})
}

// FileDownloadFileHandler serves a file for download based on the filename provided in the URL.
// @Summary File Download a specific file
// @Description File Download a file by its filename
// @Tags File
// @Param filename path string true "Filename of the file to download"
// @Produce application/octet-stream
// @Success 200 {file} file "The file for download"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 404 {object} map[string]interface{} "File Not Found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /file_download/{filename} [get]
func (h *handler) FileDownloadFileHandler(c *gin.Context) {
	// Get the filename from the URL parameter
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	// Construct the file path
	filePath := "./uploads/" + filename

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Set the headers to indicate a file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/octet-stream")

	// Serve the file
	c.File(filePath)
}

// GetProfileImageHandler serves the requested image dynamically
// @Summary Get a specific Profile image
// @Description Get an image by its filename
// @Tags File
// @Accept json
// @Produce image/png, image/jpeg, image/gif
// @Param filename path string true "Filename of the image"
// @Success 200 {file} file "The image file"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 404 {object} map[string]interface{} "Image Not Found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /profile_image/{filename} [get]
func (h *handler) GetProfileImageHandler(c *gin.Context) {
	// Get the filename from the URL parameter
	filename := c.Param("filename")
	if filename == "" {
		// Handle missing filename
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	// Use the absolute path to the uploads directory
	// Make sure to update the path according to your system's structure
	imagePath := "./uploads/profile_images/" + filename

	// Check if the file exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		// If the file doesn't exist, return a 404 error
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	// Set the proper content type based on the file extension
	ext := filepath.Ext(filename)
	var contentType string
	switch ext {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	case ".pdf":
		contentType = "image/pdf"
	case ".docx":
		contentType = "image/docx"
	default:
		// If the extension is unsupported, return an error
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Unsupported image type"})
		return
	}

	// Set the response header to the content type
	c.Header("Content-Type", contentType)

	// Serve the image file
	c.File(imagePath)
}
