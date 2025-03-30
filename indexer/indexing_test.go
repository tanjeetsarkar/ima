package indexer

import (
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestIsImageFile(t *testing.T) {
	tests := []struct {
		filename string
		expected bool
	}{
		{"image.jpg", true},
		{"image.jpeg", true},
		{"image.png", true},
		{"image.gif", true},
		{"document.pdf", false},
		{"archive.zip", false},
		{"", false},
	}

	for _, test := range tests {
		result := isImageFile(test.filename)
		if result != test.expected {
			t.Errorf("isImageFile(%q) = %v; want %v", test.filename, result, test.expected)
		}
	}
}

func TestGenerateIndexHTML(t *testing.T) {
	// Create a temporary directory for testing.
	tempDir := t.TempDir()

	// Create mock subdirectories and image files.
	os.Mkdir(filepath.Join(tempDir, "subdir1"), 0755)
	os.Mkdir(filepath.Join(tempDir, "subdir2"), 0755)
	os.WriteFile(filepath.Join(tempDir, "image1.jpg"), []byte{}, 0644)
	os.WriteFile(filepath.Join(tempDir, "image2.png"), []byte{}, 0644)
	os.WriteFile(filepath.Join(tempDir, "document.txt"), []byte{}, 0644)

	// Call GenerateIndexHTML.
	err := GenerateIndexHTML(tempDir)
	if err != nil {
		t.Fatalf("GenerateIndexHTML failed: %v", err)
	}

	// Verify that index.html was created.
	indexPath := filepath.Join(tempDir, "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Fatalf("index.html was not created")
	}

	// Verify the contents of index.html.
	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("Failed to read index.html: %v", err)
	}

	if !strings.Contains(string(content), "subdir1") || !strings.Contains(string(content), "subdir2") {
		t.Errorf("index.html does not contain expected subdirectory links")
	}
	if !strings.Contains(string(content), "image1.jpg") || !strings.Contains(string(content), "image2.png") {
		t.Errorf("index.html does not contain expected image links")
	}
	if strings.Contains(string(content), "document.txt") {
		t.Errorf("index.html contains unexpected non-image file")
	}
}

func TestSplitCreate(t *testing.T) {
	// Create a temporary directory for testing.
	tempDir := t.TempDir()

	// Create nested directories and files.
	os.Mkdir(filepath.Join(tempDir, "subdir1"), 0755)
	os.Mkdir(filepath.Join(tempDir, "subdir1", "nested"), 0755)
	os.WriteFile(filepath.Join(tempDir, "subdir1", "image1.jpg"), []byte{}, 0644)
	os.WriteFile(filepath.Join(tempDir, "subdir1", "nested", "image2.png"), []byte{}, 0644)

	// Call SplitCreate.
	SplitCreate(tempDir)

	// Verify that index.html was created in all directories.
	pathsToCheck := []string{
		filepath.Join(tempDir, "index.html"),
		filepath.Join(tempDir, "subdir1", "index.html"),
		filepath.Join(tempDir, "subdir1", "nested", "index.html"),
	}

	for _, path := range pathsToCheck {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("index.html was not created in %s", path)
		}
	}
}
func TestGenerateThumbnail(t *testing.T) {
	// Create a temporary directory for testing.
	tempDir := t.TempDir()

	// Create a mock image file.
	imagePath := filepath.Join(tempDir, "test.jpg")
	thumbnailPath := filepath.Join(tempDir, "test_thumbnail.jpg")

	// Generate a simple mock image and save it as a JPEG.
	img := image.NewRGBA(image.Rect(0, 0, 300, 300)) // Create a 300x300 image.
	outFile, err := os.Create(imagePath)
	if err != nil {
		t.Fatalf("Failed to create mock image file: %v", err)
	}
	defer outFile.Close()
	if err := jpeg.Encode(outFile, img, &jpeg.Options{Quality: 80}); err != nil {
		t.Fatalf("Failed to encode mock image: %v", err)
	}

	// Call generateThumbnail.
	err = generateThumbnail(imagePath, thumbnailPath)
	if err != nil {
		t.Fatalf("generateThumbnail failed: %v", err)
	}

	// Verify that the thumbnail file was created.
	if _, err := os.Stat(thumbnailPath); os.IsNotExist(err) {
		t.Fatalf("Thumbnail file was not created: %s", thumbnailPath)
	}

	// Verify the dimensions of the generated thumbnail.
	thumbnailFile, err := os.Open(thumbnailPath)
	if err != nil {
		t.Fatalf("Failed to open thumbnail file: %v", err)
	}
	defer thumbnailFile.Close()

	thumbnailImg, _, err := image.Decode(thumbnailFile)
	if err != nil {
		t.Fatalf("Failed to decode thumbnail image: %v", err)
	}

	// Check if the thumbnail dimensions are 150x150.
	if thumbnailImg.Bounds().Dx() != 150 || thumbnailImg.Bounds().Dy() != 150 {
		t.Errorf("Thumbnail dimensions are incorrect: got %dx%d, want 150x150",
			thumbnailImg.Bounds().Dx(), thumbnailImg.Bounds().Dy())
	}
}

func TestGenerateThumbnailInvalidInput(t *testing.T) {
	// Create a temporary directory for testing.
	tempDir := t.TempDir()

	// Define invalid image and thumbnail paths.
	invalidImagePath := filepath.Join(tempDir, "nonexistent.jpg")
	thumbnailPath := filepath.Join(tempDir, "test_thumbnail.jpg")

	// Call generateThumbnail with a nonexistent image file.
	err := generateThumbnail(invalidImagePath, thumbnailPath)
	if err == nil {
		t.Fatalf("Expected an error for nonexistent image file, but got none")
	}
}

func TestGenerateThumbnailInvalidOutputPath(t *testing.T) {
	// Create a temporary directory for testing.
	tempDir := t.TempDir()

	// Create a mock image file.
	imagePath := filepath.Join(tempDir, "test.jpg")
	img := image.NewRGBA(image.Rect(0, 0, 300, 300)) // Create a 300x300 image.
	outFile, err := os.Create(imagePath)
	if err != nil {
		t.Fatalf("Failed to create mock image file: %v", err)
	}
	defer outFile.Close()
	if err := jpeg.Encode(outFile, img, &jpeg.Options{Quality: 80}); err != nil {
		t.Fatalf("Failed to encode mock image: %v", err)
	}

	// Define an invalid thumbnail path (e.g., a directory instead of a file).
	invalidThumbnailPath := filepath.Join(tempDir, "invalid_dir")
	if err := os.Mkdir(invalidThumbnailPath, 0755); err != nil {
		t.Fatalf("Failed to create invalid directory: %v", err)
	}

	// Call generateThumbnail with an invalid output path.
	err = generateThumbnail(imagePath, invalidThumbnailPath)
	if err == nil {
		t.Fatalf("Expected an error for invalid thumbnail path, but got none")
	}
}
func TestParallelThumbnailGeneration(t *testing.T) {
	// Create a temporary directory for testing.
	tempDir := t.TempDir()

	// Create mock image files.
	imagePaths := []string{
		filepath.Join(tempDir, "image1.jpg"),
		filepath.Join(tempDir, "image2.jpg"),
		filepath.Join(tempDir, "image3.jpg"),
	}
	for _, imagePath := range imagePaths {
		img := image.NewRGBA(image.Rect(0, 0, 300, 300)) // Create a 300x300 image.
		outFile, err := os.Create(imagePath)
		if err != nil {
			t.Fatalf("Failed to create mock image file: %v", err)
		}
		if err := jpeg.Encode(outFile, img, &jpeg.Options{Quality: 80}); err != nil {
			t.Fatalf("Failed to encode mock image: %v", err)
		}
		outFile.Close()
	}

	// Create a .thumbs directory.
	thumbsDir := filepath.Join(tempDir, ".thumbs")
	if err := os.MkdirAll(thumbsDir, 0755); err != nil {
		t.Fatalf("Failed to create .thumbs directory: %v", err)
	}

	// Generate thumbnails in parallel.
	var wg sync.WaitGroup
	errChan := make(chan error, len(imagePaths))
	for _, imagePath := range imagePaths {
		thumbnailPath := filepath.Join(thumbsDir, filepath.Base(imagePath))
		wg.Add(1)
		go func(imagePath, thumbnailPath string) {
			defer wg.Done()
			if err := generateThumbnail(imagePath, thumbnailPath); err != nil {
				errChan <- err
			}
		}(imagePath, thumbnailPath)
	}

	// Wait for all goroutines to finish.
	wg.Wait()
	close(errChan)

	// Check for errors.
	for err := range errChan {
		t.Errorf("Error during thumbnail generation: %v", err)
	}

	// Verify that all thumbnails were created.
	for _, imagePath := range imagePaths {
		thumbnailPath := filepath.Join(thumbsDir, filepath.Base(imagePath))
		if _, err := os.Stat(thumbnailPath); os.IsNotExist(err) {
			t.Errorf("Thumbnail was not created for %s", imagePath)
		}
	}
}

func TestGenerateIndexHTMLWithParallelThumbnails(t *testing.T) {
	// Create a temporary directory for testing.
	tempDir := t.TempDir()

	// Create mock subdirectories and image files.
	os.Mkdir(filepath.Join(tempDir, "subdir1"), 0755)
	os.Mkdir(filepath.Join(tempDir, "subdir2"), 0755)
	os.WriteFile(filepath.Join(tempDir, "image1.jpg"), []byte{}, 0644)
	os.WriteFile(filepath.Join(tempDir, "image2.png"), []byte{}, 0644)
	os.WriteFile(filepath.Join(tempDir, "document.txt"), []byte{}, 0644)

	// Call GenerateIndexHTML.
	err := GenerateIndexHTML(tempDir)
	if err != nil {
		t.Fatalf("GenerateIndexHTML failed: %v", err)
	}

	// Verify that index.html was created.
	indexPath := filepath.Join(tempDir, "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Fatalf("index.html was not created")
	}

	// Verify that thumbnails were created in parallel.
	thumbsDir := filepath.Join(tempDir, ".thumbs")
	thumbnailPaths := []string{
		filepath.Join(thumbsDir, "image1.jpg"),
		filepath.Join(thumbsDir, "image2.png"),
	}
	for _, thumbnailPath := range thumbnailPaths {
		if _, err := os.Stat(thumbnailPath); os.IsNotExist(err) {
			t.Errorf("Thumbnail was not created: %s", thumbnailPath)
		}
	}
}
