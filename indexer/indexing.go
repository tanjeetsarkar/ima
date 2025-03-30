package indexer

import (
	"html/template"
	"image"
	"image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"golang.org/x/image/draw"
)

// PageData holds the data for our HTML template.
type PageData struct {
	Title       string
	SubDirs     []SubDir
	Images      []string
	CurrentPath string
	Thumbs      bool
}

// SubDir represents a subdirectory entry for the sidebar.
type SubDir struct {
	Name string // Display name
	Link string // Relative link to the subdirectory's index.html
}

// HTML template for index.html pages.
var indexTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>{{.Title}}</title>
  <style>
    /* Basic reset */
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      font-family: Arial, sans-serif;
      display: flex;
      height: 100vh;
    }
    /* Sidebar */
    .sidebar {
      width: 220px;
      background: #f0f0f0;
      overflow-y: auto;
      border-right: 1px solid #ccc;
      padding: 20px;
      resize: horizontal;
      overflow: hidden;
      min-width: 150px;
      max-width: 400px;
    }
    .sidebar ul { list-style: none; }
    .sidebar li { margin-bottom: 10px; }
    .sidebar a {
      text-decoration: none;
      color: #333;
      display: block;
      padding: 5px 10px;
      border-radius: 4px;
    }
    .sidebar a:hover, .sidebar a.active { background-color: #ddd; }
    /* Content */
    .content {
      flex: 1;
      padding: 20px;
      overflow-y: auto;
    }
    .grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
      grid-gap: 15px;
    }
    .grid img {
      width: 100%;
      height: 100%;
      object-fit: cover;
      display: block;
      cursor: pointer;
    }
     
    .modal {
      display: none;
      position: fixed;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background: rgba(0, 0, 0, 0.8);
      justify-content: center;
      align-items: center;
    }
    .modal img {
      max-width: 95vh;
      max-height: 95vh;
      object-fit: contain;
    }
    .modal:target {
      display: flex;
    }
  </style>
</head>
<body>
  {{if .SubDirs}}
  <div class="sidebar">
    <ul>
      {{range .SubDirs}}
      <li><a href="{{.Link}}/index.html">{{.Name}}</a></li>
      {{end}}
    </ul>
  </div>
  {{end}}
  <div class="content">
    <h1>{{.Title}}</h1>
    {{if .Images}}
    <div class="grid">
      {{range .Images}}
      <a href="#modal-{{.}}">
      {{if $.Thumbs}}
        <img loading="lazy" src=".thumbs/{{.}}" alt="">
      {{else}}
        <img loading="lazy" src="{{.}}" alt="">
      {{end}}
      </a>
      <div id="modal-{{.}}" class="modal">
        <img src="{{.}}" alt="">
      </div>
      {{end}}
    </div>
    {{else}}
    <p>No images in this folder.</p>
    {{end}}
  </div>
  <script>
    const links = document.querySelectorAll('.sidebar a');
    links.forEach(link => {
      link.addEventListener('click', function() {
        links.forEach(lnk => lnk.classList.remove('active'));
        this.classList.add('active');
      });
    });
    document.addEventListener('DOMContentLoaded', () => {
    const modals = document.querySelectorAll('.modal');
    const images = Array.from(document.querySelectorAll('.grid a'));
    const modalImages = Array.from(document.querySelectorAll('.modal img'));

    // Close modal when clicking outside the image
    modals.forEach(modal => {
      modal.addEventListener('click', (e) => {
        if (e.target === modal) {
          window.location.hash = ''; // Close the modal
        }
      });
    });

    // Navigate with left and right arrow keys
    document.addEventListener('keydown', (e) => {
      const currentHash = window.location.hash;
      if (!currentHash) return;

      const currentIndex = images.findIndex(link => ` + "`" + `#${link.getAttribute('href').substring(1)}` + "`" + `=== currentHash);

      if (e.key === 'ArrowRight') {
        const nextIndex = (currentIndex + 1) % images.length;
        window.location.hash = ` + "`" + `#${images[nextIndex].getAttribute('href').substring(1)}` + "`" + `;
      } else if (e.key === 'ArrowLeft') {
        const prevIndex = (currentIndex - 1 + images.length) % images.length;
        window.location.hash = ` + "`" + `#${images[prevIndex].getAttribute('href').substring(1)}` + "`" + `;
      }
    });
  });
  </script>
</body>
</html>
`

var noThumb bool // Global variable to track the --nothumb flag

func AddFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&noThumb, "nothumb", false, "Disable thumbnail generation")
}

// isImageFile checks if a file extension is an image type.
func isImageFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif":
		return true
	}
	return false
}

func GenerateIndexHTML(dir string) error {
	// List items in the directory.
	items, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var subDirs []SubDir
	var images []string

	// Create .thumbs directory if thumbnails are enabled
	thumbsDir := filepath.Join(dir, ".thumbs")
	if !noThumb {
		if err := os.MkdirAll(thumbsDir, os.ModePerm); err != nil {
			return err
		}
	}

	parentDir := filepath.Dir(dir)

	if parentDir != dir && parentDir != "." && parentDir != "/" { // Avoid adding ".." for the root directory.
		subDirs = append(subDirs, SubDir{
			Name: "..",
			Link: "..",
		})
	}

	// Channel to collect errors from goroutines
	errChan := make(chan error, len(items))
	// WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	for _, item := range items {
		log.Printf("Processing %s", item.Name())
		if item.IsDir() {
			// Add subdirectory link.
			if item.Name() != ".thumbs" {
				subDirs = append(subDirs, SubDir{
					Name: item.Name(),
					Link: item.Name(),
				})
			}
		} else if isImageFile(item.Name()) {
			// Add image file.
			images = append(images, item.Name())

			if !noThumb {
				imagePath := filepath.Join(dir, item.Name())
				thumbnailPath := filepath.Join(thumbsDir, item.Name())

				// Check if thumbnail already exists
				if _, err := os.Stat(thumbnailPath); os.IsNotExist(err) {
					wg.Add(1)
					go func(imagePath, thumbnailPath string) {
						defer wg.Done()
						if err := generateThumbnail(imagePath, thumbnailPath); err != nil {
							log.Printf("Failed to generate thumbnail for %s: %v", imagePath, err)
							errChan <- err
						}
					}(imagePath, thumbnailPath)
				}
			}
		}
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errChan)

	// Check if there were any errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	// Prepare template data.
	data := PageData{
		Title:       filepath.Base(dir),
		SubDirs:     subDirs,
		Images:      images,
		CurrentPath: dir,
		Thumbs:      !noThumb,
	}

	tmpl, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		return err
	}

	// Create or overwrite index.html.
	f, err := os.Create(filepath.Join(dir, "index.html"))
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}

func SplitCreate(rootDir string) {
	rootDir = strings.TrimSuffix(rootDir, string(os.PathSeparator))
	log.Printf("Indexing directory : %s", rootDir)
	// Walk through each directory and generate an index.html.
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// Create/update the index.html for this directory.
			if filepath.Base(path) == ".thumbs" {
				return nil // Skip the .thumbs directory
			}
			if err := GenerateIndexHTML(path); err != nil {
				// Handle the error as needed (e.g., log it).
				return err
			}
		}
		return nil
	})
	if err != nil {
		// Handle the error from filepath.Walk.
		return
	}
}

func generateThumbnail(imagePath, thumbnailPath string) error {
	// Open the original image file.
	file, err := os.Open(imagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode the image.
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	// Resize the image to a thumbnail (e.g., 150x150).
	thumbnail := resizeImage(img, 150, 150)

	// Create the thumbnail file.
	outFile, err := os.Create(thumbnailPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Save the thumbnail as a JPEG.
	return jpeg.Encode(outFile, thumbnail, &jpeg.Options{Quality: 80})
}

func resizeImage(img image.Image, width, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dst
}
