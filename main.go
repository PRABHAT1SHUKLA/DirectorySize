package main

import (
	"fmt"
	"io/fs" // For file system abstraction, especially fs.FileInfo
	"os"    // For command-line arguments, file operations, Stat (file info)
	"path/filepath" // For walking directory trees
	// "strconv" // REMOVED: This import is no longer needed as we don't use strconv.Atoi or similar here.
)

// humanReadableBytes converts bytes to a human-readable string (KB, MB, GB, etc.)
func humanReadableBytes(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/TB)
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB) // FIXED: Changed float664 to float64
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}

func main() {
	// 1. Check for command-line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <directory_path>")
		fmt.Println("Example: go run main.go .") // '.' means current directory
		fmt.Println("Example: go run main.go /path/to/my/folder")
		os.Exit(1) // Exit with a non-zero status code to indicate an error
	}

	// Get the directory path from the command-line argument
	targetDir := os.Args[1]

	// 2. Start the calculation
	fmt.Printf("Calculating size for '%s'...\n", targetDir)

	totalSize, err := calculateDirSize(targetDir)
	if err != nil {
		fmt.Printf("Error calculating size: %v\n", err)
		os.Exit(1)
	}

	// 3. Print the result
	fmt.Printf("Total size of '%s': %s (%d bytes)\n", targetDir, humanReadableBytes(totalSize), totalSize)
}

// calculateDirSize recursively calculates the total size of files in a directory.
// It returns the total size in bytes and an error if one occurs during traversal.
func calculateDirSize(dirPath string) (int64, error) {
	var totalSize int64 = 0

	// filepath.Walk is a powerful function that walks the file tree rooted at 'dirPath'.
	// It calls a function (the 'walkFn') for each file or directory in the tree,
	// including the root itself.
	walkFn := func(path string, info fs.FileInfo, err error) error {
		// If there's an error accessing a file/directory (e.g., permissions),
		// we print it but continue walking the rest of the tree.
		if err != nil {
			fmt.Printf("Error accessing %s: %v\n", path, err)
			return nil // Don't return the error, just skip this item
		}

		// If the current item is a directory, we just skip it.
		// Its contents will be visited by subsequent calls to walkFn.
		if info.IsDir() {
			return nil // Continue walking
		}

		// If it's a file, add its size to the total.
		totalSize += info.Size()

		return nil // Continue walking
	}

	// Start walking the directory tree.
	// The walkFn will be called for each file and directory.
	err := filepath.Walk(dirPath, walkFn)
	if err != nil {
		// This error indicates a problem starting the walk (e.g., dirPath doesn't exist)
		return 0, fmt.Errorf("failed to walk directory '%s': %w", dirPath, err)
	}

	return totalSize, nil
}