package main

import (
	"encoding/json"
	"flag"
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

const (
	DATA_PATH  = "PackSync\\"
	MOD_PATH   = "Code\\"
	PURGE_PATH = "content\\"
)

func CheckUpdateAvailable(repo string) (update bool, version string, release Release) {
	fmt.Println("Checking latest version...")

	resp, err := http.Get("https://api.github.com/repos/" + repo + "/releases/latest")
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

	versionFile, err := os.ReadFile(core.GetAppDataDir(DATA_PATH) + ".version")
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
	if err := os.MkdirAll(core.GetAppDataDir(DATA_PATH)+"temp", 0755); err != nil {
		fmt.Println("Error creating temp directory:", err)
		return
	}

	if len(release.Assets) > 0 {
		downloadURL := release.Assets[0].BrowserDownloadURL
		dest := filepath.Join(core.GetAppDataDir(DATA_PATH), "temp", filepath.Base(downloadURL))
		err := core.DownloadPackage(downloadURL, dest)
		if err != nil {
			fmt.Println("Error downloading package:", err)
			return
		}

		fmt.Println("Package downloaded for version:", version)
		os.WriteFile(core.GetAppDataDir(DATA_PATH)+".version", []byte(version), 0644)
	} else {
		fmt.Println("No assets found for the latest release.")
	}
}

func InstallPack(force bool) {
	update, version, release := CheckUpdateAvailable("cutil-inv/gopher-lua")
	dataPath := core.GetAppDataDir(DATA_PATH)

	if force {
		fmt.Println("Forcing update...")
	} else if !update {
		fmt.Println("You are using the latest version:", version)
		return
	} else {
		fmt.Println("A newer version is available:", version)
	}

	RetrieveLatestVersion(release, version)

	location, err := core.FindProgramLocation(MOD_PATH)
	if err != nil {
		fmt.Println("Error finding directory location:", err)
	}
	location += "test\\"

	zipFile := dataPath + "temp\\" + filepath.Base(release.Assets[0].BrowserDownloadURL)
	destDir := dataPath + "content\\"
	err = core.Unzip(zipFile, destDir)
	if err != nil {
		fmt.Println("Error unzipping file:", err)
	}

	os.RemoveAll(location + "\\" + PURGE_PATH)

	err = core.CopyFiles(destDir, location)
	if err != nil {
		fmt.Println("Error copying files:", err)
	}

	os.RemoveAll(core.GetAppDataDir(DATA_PATH) + "temp\\")
}

func UpdateSelf() {
	update, version, release := CheckUpdateAvailable("cutil-inv/PackSync")
	dataPath := core.GetAppDataDir(DATA_PATH)

	if !update {
		fmt.Println("You are using the latest version:", version)
		return
	}

	RetrieveLatestVersion(release, version)

	newFile := dataPath + "temp\\" + filepath.Base(release.Assets[0].BrowserDownloadURL)
	core.CopyFiles(newFile, os.Args[0])

	os.RemoveAll(core.GetAppDataDir(DATA_PATH) + "temp\\")
}

func main() {
	force := flag.Bool("force", false, "force update")
	f := flag.Bool("f", false, "force update")

	flag.Parse()

	switch flag.Arg(0) {
	case "self-update":
		UpdateSelf()
	default:
		InstallPack(*force || *f)
	}
}
