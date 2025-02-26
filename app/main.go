package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	core "PackSync.core"
)

type Release struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func CheckUpdateAvailable() (update bool, version string, release Release) {
	fmt.Println("Checking latest version...")

	resp, err := http.Get("https://api.github.com/repos/cutil-inv/gopher-lua/releases/latest")
	if err != nil {
		fmt.Println("Error fetching version:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	err = json.Unmarshal(body, &release)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	latestVersion := release.TagName
	currentVersion := "0.0.0"

	versionFile, err := os.ReadFile(".version")
	if err == nil {
		currentVersion = string(versionFile)
	}

	if latestVersion != currentVersion {
		version = latestVersion
		update = true
	} else {
		version = currentVersion
	}

	return
}

func RetrieveLatestVersion(release Release, version string) {
	if err := os.MkdirAll("downloads", 0755); err != nil {
		fmt.Println("Error creating downloads directory:", err)
		return
	}

	if len(release.Assets) > 0 {
		downloadURL := release.Assets[0].BrowserDownloadURL
		dest := filepath.Join("downloads", filepath.Base(downloadURL))
		err := core.DownloadPackage(downloadURL, dest)
		if err != nil {
			fmt.Println("Error downloading package:", err)
			return
		}

		fmt.Println("Success! Package downloaded for version:", version)
		os.WriteFile(".version", []byte(version), 0644)
	} else {
		fmt.Println("No assets found for the latest release.")
	}
}

func main() {
	update, version, release := CheckUpdateAvailable()

	if update {
		fmt.Println("A newer version is available:", version)

		RetrieveLatestVersion(release, version)
	} else {
		fmt.Println("You are using the latest version:", version)
	}

	directory := "Code" // Replace with the directory you want to find in AppData
	location, err := core.FindProgramLocation(directory)
	if err != nil {
		fmt.Println("Error finding directory location:", err)
	} else {
		fmt.Printf("The location of %s in AppData is: %s\n", directory, location)
	}
}
