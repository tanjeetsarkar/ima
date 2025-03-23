package indexer

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/schollz/progressbar/v3"
)

const outputFile = "gallery.html"

type Folder struct {
	Name     string
	Path     string
	Children []Folder
	Images   []string
}

func CreateIndexHtml(root string) {
	folders := buildFolderStructure(root)
	generateHTML(folders)
}

func buildFolderStructure(root string) []Folder {
	var folders []Folder

	// Count total directories for progress bar
	totalDirs := 0
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err == nil && d.IsDir() && !shouldSkipDir(path) {
			totalDirs++
		}
		return nil
	})

	// Initialize progress bar
	bar := progressbar.NewOptions(totalDirs,
		progressbar.OptionSetDescription("Indexing directories..."),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(30),
		progressbar.OptionClearOnFinish(),
	)

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if shouldSkipDir(path) {
				return fs.SkipDir
			}

			// Update progress bar and print current directory
			bar.Add(1)
			fmt.Printf("\rProcessing: %s\n", path)

			relPath, _ := filepath.Rel(root, path)
			folder := Folder{
				Name: filepath.Base(path),
				Path: relPath,
			}

			// Get images
			files, _ := os.ReadDir(path)
			for _, f := range files {
				if !f.IsDir() && isImage(f.Name()) {
					folder.Images = append(folder.Images, f.Name())
				}
			}

			folders = addToStructure(folders, strings.Split(relPath, string(filepath.Separator)), folder)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return folders
}

func generateHTML(folders []Folder) {
	tmpl := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>Image Gallery</title>
        <style>
            body { margin: 0; padding: 0; display: flex; }
            .sidebar { width: 300px; height: 100vh; overflow-y: auto; padding: 20px; }
            .content { flex: 1; padding: 20px; display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 10px; }
            .folder { cursor: pointer; padding: 5px; }
            .folder:hover { background: #eee; }
            img { width: 100%; height: 200px; object-fit: cover; }
            ul { list-style: none; padding-left: 15px; }
        </style>
    </head>
    <body>
        <div class="sidebar">
            {{- range .}}
                {{template "folderTemplate" .}}
            {{- end}}
        </div>
        <div class="content" id="content"></div>

        <script>
            function showImages(path, images) {
                if (images) {
                const content = document.getElementById('content');
                content.innerHTML = images.map(img =>` +
		"`<img" + ` src="${path}/${img}" alt="${img}">` + "`" + `
                ).join('');
                }
            }
        </script>
    </body>
    </html>

    {{define "folderTemplate"}}
        <li class="folder" onclick="showImages('{{.Path}}', {{.Images}})">
            📁 {{.Name}}
            {{if .Children}}
                <ul>
                    {{range .Children}}
                        {{template "folderTemplate" .}}
                    {{end}}
                </ul>
            {{end}}
        </li>
    {{end}}`

	t := template.Must(template.New("html").Parse(tmpl))

	f, err := os.Create(outputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = t.Execute(f, folders)
	if err != nil {
		log.Fatal(err)
	}
}

// Helper functions
func shouldSkipDir(path string) bool {
	return filepath.Base(path)[0] == '.' // Skip hidden directories
}

func isImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif"
}

func addToStructure(folders []Folder, path []string, newFolder Folder) []Folder {
	if len(path) == 0 {
		return append(folders, newFolder)
	}

	for i, f := range folders {
		if f.Name == path[0] {
			folders[i].Children = addToStructure(f.Children, path[1:], newFolder)
			return folders
		}
	}

	return append(folders, newFolder)
}
