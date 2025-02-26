package core

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadPackage(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error fetching package: %w", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}

func FindProgramLocation(directory string) (string, error) {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		return "", fmt.Errorf("APPDATA environment variable not set")
	}

	dirPath := filepath.Join(appData, directory)
	if stat, err := os.Stat(dirPath); err == nil && stat.IsDir() {
		return dirPath, nil
	}

	return "", fmt.Errorf("directory %s not found in AppData", directory)
}
