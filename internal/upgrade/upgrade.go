package upgrade

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	gitHubAPI       = "https://api.github.com/repos/nemethk/claudeenv/releases/latest"
	releaseBaseURL  = "https://github.com/nemethk/claudeenv/releases/download"
	installPath     = "/usr/local/bin/claudeenv"
)

type Release struct {
	TagName string `json:"tag_name"`
}

func Run() error {
	// Fetch latest release
	fmt.Println("Fetching latest release...")
	release, err := fetchLatestReleaseFromURL(gitHubAPI)
	if err != nil {
		return fmt.Errorf("failed to fetch release: %w", err)
	}

	fmt.Printf("Latest version: %s\n", release.TagName)

	// Detect OS and architecture
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	// Map Go arch names to release arch names
	archName := goarch
	if goarch == "amd64" {
		archName = "amd64"
	} else if goarch == "arm64" {
		archName = "arm64"
	} else {
		return fmt.Errorf("unsupported architecture: %s", goarch)
	}

	osName := goos
	if goos != "linux" && goos != "darwin" {
		return fmt.Errorf("unsupported OS: %s", goos)
	}

	fmt.Printf("Detected: %s %s\n", osName, archName)

	// Download and extract
	tmpDir, err := os.MkdirTemp("", "claudeenv-upgrade-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	archiveName := fmt.Sprintf("claudeenv_%s_%s.tar.gz", osName, archName)
	downloadURL := fmt.Sprintf("%s/%s/%s", releaseBaseURL, release.TagName, archiveName)

	fmt.Printf("Downloading from: %s\n", downloadURL)
	binaryPath, err := downloadAndExtract(downloadURL, tmpDir)
	if err != nil {
		return fmt.Errorf("failed to download and extract: %w", err)
	}

	// Replace binary with sudo
	fmt.Printf("Installing to %s...\n", installPath)
	err = replaceBinary(binaryPath, installPath)
	if err != nil {
		return err
	}

	// Verify
	fmt.Println("Verifying installation...")
	cmd := exec.Command(installPath, "--version")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to verify: %w", err)
	}

	fmt.Printf("Upgrade complete: %s\n", string(output))
	return nil
}

func fetchLatestReleaseFromURL(url string) (*Release, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release Release
	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		return nil, err
	}

	if release.TagName == "" {
		return nil, fmt.Errorf("no release tag found")
	}

	return &release, nil
}

func downloadAndExtract(url, destDir string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Decompress gzip
	gzReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read gzip: %w", err)
	}
	defer gzReader.Close()

	// Extract tar
	tarReader := tar.NewReader(gzReader)
	binaryPath := ""

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("failed to read tar: %w", err)
		}

		// Extract claudeenv binary
		if header.Name == "claudeenv" {
			filePath := filepath.Join(destDir, "claudeenv")
			file, err := os.Create(filePath)
			if err != nil {
				return "", fmt.Errorf("failed to create file: %w", err)
			}

			_, err = io.Copy(file, tarReader)
			file.Close()
			if err != nil {
				return "", fmt.Errorf("failed to extract binary: %w", err)
			}

			// Make executable
			err = os.Chmod(filePath, 0755)
			if err != nil {
				return "", fmt.Errorf("failed to chmod: %w", err)
			}

			binaryPath = filePath
		}
	}

	if binaryPath == "" {
		return "", fmt.Errorf("binary not found in archive")
	}

	return binaryPath, nil
}

func replaceBinary(src, dest string) error {
	// Check if running as root or sudoed
	if os.Geteuid() != 0 {
		return fmt.Errorf("upgrade requires sudo: run 'sudo claudeenv upgrade'")
	}

	// Copy new binary to destination
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source binary: %w", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create destination binary: %w", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy binary: %w", err)
	}

	err = destFile.Chmod(0755)
	if err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	return nil
}
