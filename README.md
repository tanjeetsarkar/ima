# Image Archive Program

This idea was originally scoped by [Chris](https://github.com/chrissy-dev) in his [blog](https://www.scottishstoater.com/2025/03/scoping-a-local-first-image-archive/)

## About the Program
The Image Archive Program is a CLI-based tool designed to simplify the process of organizing and managing image collections. It generates `index.html` files for directories, creating a navigable gallery structure for images. The program also includes a watcher feature to monitor directories for changes and automatically update the gallery.

## Advantages of the Program
- **Automated Indexing**: Automatically generates `index.html` files with links to subdirectories and previews of images in the directory.
- **Minimal JavaScript**: Ensures lightweight and fast-loading HTML pages.
- **Resizable Sidebar**: Allows users to adjust the sidebar width for better visibility of directory names.
- **Keyboard Navigation**: Supports navigation between images using arrow keys.
- **Dynamic Updates**: Includes a watcher feature to detect changes in directories and update the gallery in real-time.

## Usage of the Program
1. **Build the Program**:
   - Run the following command to build the program for your current OS:
     ```sh
     make build
     ```
   - Alternatively, build for specific platforms:
     ```sh
     make build-linux
     make build-windows
     ```

2. **Run the Program**:
   - To generate an image gallery for a directory:
     ```sh
     ./image-archive [directory]
     ```
   - To enable the watcher feature:
     ```sh
     ./image-archive [directory] --watch
     ```
   - To disable thumbnail generation:
     ```sh
     ./image-archive [directory] --nothumbs
     ```

3. **Clean Build Artifacts**:
   - Use the following command to clean up build artifacts:
     ```sh
     make clean
     ```

4. **Run Tests**:
   - Execute tests to ensure the program works as expected:
     ```sh
     make test
     ```

## Roadmap of the Program
- **Enhanced Gallery Features**:
  - Add support for additional image formats (e.g., `.bmp`, `.tiff`).
  - Implement search functionality for images and directories.

- **Improved Watcher**:
  - Add support for filtering events based on file types or patterns.
  - Optimize performance for large directories.

- **Customization Options**:
  - Allow users to customize the appearance of the gallery (e.g., themes, layouts).
  - Add configuration options for excluding specific directories or files.

- **Advanced Features**:
  - Add support for generating thumbnails for faster loading.
  - Implement a feature to export the gallery as a standalone package.
  - Add support for adding and displaying metadata for images and directories (e.g., titles, descriptions, tags).
  - Metadata based searching and sorting 

- **Documentation and Tutorials**:
  - Provide detailed documentation and video tutorials for new users.
  - Add examples and templates for common use cases.
