package indexer

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// PageData holds the data for our HTML template.
type PageData struct {
	Title       string
	SubDirs     []SubDir
	Images      []string
	CurrentPath string
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
      grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
      grid-gap: 15px;
    }
    .grid img {
      width: 300px;
      height: 300px;
      display: block;
      border: 1px solid #ccc;
      border-radius: 4px;
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
      <img src="{{.}}" alt="">
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
  </script>
</body>
</html>
`

// isImageFile checks if a file extension is an image type.
func isImageFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif":
		return true
	}
	return false
}

// GenerateIndexHTML creates or updates an index.html file in the given directory.
func GenerateIndexHTML(dir string) error {
	// List items in the directory.
	items, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var subDirs []SubDir
	var images []string

	for _, item := range items {
		if item.IsDir() {
			// Add subdirectory link.
			subDirs = append(subDirs, SubDir{
				Name: item.Name(),
				Link: item.Name(),
			})
		} else if isImageFile(item.Name()) {
			// Add image file.
			images = append(images, item.Name())
		}
	}

	// Prepare template data.
	data := PageData{
		Title:       filepath.Base(dir),
		SubDirs:     subDirs,
		Images:      images,
		CurrentPath: dir,
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
	log.Printf("Indexing directory : %s", rootDir)
	// Walk through each directory and generate an index.html.
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// Create/update the index.html for this directory.
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
