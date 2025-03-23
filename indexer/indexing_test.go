package indexer

import (
	"os"
	"path/filepath"
	"strings"
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
