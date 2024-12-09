package pages

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const FilePathsJSON = "file_paths.json"

// normalizePath ensures paths are canonical and consistent.
func normalizePath(path string) string {
	return filepath.Clean(path)
}

func loadFilePaths() (map[string]bool, error) {
	paths := make(map[string]bool)

	if _, err := os.Stat(FilePathsJSON); os.IsNotExist(err) {
		return paths, nil
	}

	data, err := ioutil.ReadFile(FilePathsJSON)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &paths)
	if err != nil {
		return nil, err
	}

	return paths, nil
}

func saveFilePaths(paths map[string]bool) error {
	data, err := json.MarshalIndent(paths, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(FilePathsJSON, data, 0644)
}

func CopyFile(src string, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func SyncDir(sourceDir, destDir string) error {
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		err = os.MkdirAll(destDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create destination directory: %v", err)
		}
	}

	destFiles := make(map[string]bool)
	err := filepath.Walk(destDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, err := filepath.Rel(destDir, path)
			if err != nil {
				return err
			}
			destFiles[relPath] = true
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to scan destination directory: %v", err)
	}

	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(destDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, os.ModePerm)
		}

		if _, err := os.Stat(destPath); os.IsNotExist(err) || isModified(path, destPath) {
			err = CopyFile(path, destPath)
			if err != nil {
				return fmt.Errorf("failed to copy file: %v", err)
			}
		}

		destFiles[relPath] = false
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to sync files: %v", err)
	}

	for relPath, toDelete := range destFiles {
		if toDelete {
			err := os.Remove(filepath.Join(destDir, relPath))
			if err != nil {
				return fmt.Errorf("failed to remove file: %v", err)
			}
		}
	}

	return nil
}

func isModified(sourcePath, destPath string) bool {
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return true
	}
	destInfo, err := os.Stat(destPath)
	if err != nil {
		return true
	}
	return sourceInfo.ModTime() != destInfo.ModTime() || sourceInfo.Size() != destInfo.Size()
}

// ValidateAndSync validates the source path and syncs it with the destination directory.
func ValidateAndSync(sourcePath string, destDir string) error {
	sourcePath = normalizePath(sourcePath)

	paths, err := loadFilePaths()
	if err != nil {
		return fmt.Errorf("error loading file paths: %v", err)
	}

	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return fmt.Errorf("source path does not exist: %s", sourcePath)
	}

	info, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("error stating source path: %s", sourcePath)
	}

	if !info.IsDir() {
		return fmt.Errorf("source path must be a directory: %s", sourcePath)
	}

	destPath := filepath.Join(destDir, filepath.Base(sourcePath))

	err = SyncDir(sourcePath, destPath)
	if err != nil {
		return fmt.Errorf("error syncing directories: %v", err)
	}

	paths[sourcePath] = true
	err = saveFilePaths(paths)
	if err != nil {
		return fmt.Errorf("error saving file paths: %v", err)
	}

	return nil
}

type Config struct {
	Destination string `json:"destination"`
}

func getDestinationDir() string {
	var config Config

	data, err := ioutil.ReadFile("settings.json")
	if err != nil {
		log.Fatalf("Error reading settings.json: %v", err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	return config.Destination
}
