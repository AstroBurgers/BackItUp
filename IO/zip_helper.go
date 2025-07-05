package IO

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ProgressMsg struct {
	Done  int
	Total int
}

// ZipWithExtensions ZipWithExtensionsProgress zips files from the current directory "."
// filtering by allowed extensions, and sends progress updates to the channel.
func ZipWithExtensions(
	outputZip string,
	allowedExts []string,
	progress chan<- ProgressMsg,
) error {
	inputDir := "."

	outFile, err := os.Create(outputZip)
	if err != nil {
		return err
	}
	defer func(outFile *os.File) {
		err := outFile.Close()
		if err != nil {

		}
	}(outFile)

	zipWriter := zip.NewWriter(outFile)
	defer func(zipWriter *zip.Writer) {
		err := zipWriter.Close()
		if err != nil {

		}
	}(zipWriter)

	extMap := make(map[string]bool)
	for _, ext := range allowedExts {
		extMap[strings.ToLower(ext)] = true
	}

	// Step 1: Count total files matching the extensions
	var totalFiles int
	err = filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if extMap[strings.ToLower(filepath.Ext(path))] {
			totalFiles++
		}
		return nil
	})
	if err != nil {
		return err
	}

	if totalFiles == 0 {
		return fmt.Errorf("no files matching the specified extensions were found; backup not created")
	}

	filesAdded := 0

	// Step 2: Walk again and zip files, reporting progress
	err = filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if !extMap[strings.ToLower(filepath.Ext(path))] {
			return nil
		}

		zipPath, err := filepath.Rel(inputDir, path)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {

			}
		}(file)

		w, err := zipWriter.Create(zipPath)
		if err != nil {
			return err
		}

		if _, err = io.Copy(w, file); err != nil {
			return err
		}

		filesAdded++

		// Non-blocking send progress update
		select {
		case progress <- ProgressMsg{Done: filesAdded, Total: totalFiles}:
		default:
		}

		return nil
	})

	defer close(progress)

	return err
}
