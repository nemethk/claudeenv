package upgrade

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadAndExtract(t *testing.T) {
	// Create a test tar.gz archive with a claudeenv binary
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	tarWriter := tar.NewWriter(gzWriter)

	// Add a test binary
	header := &tar.Header{
		Name: "claudeenv",
		Size: 5,
		Mode: 0755,
	}
	if err := tarWriter.WriteHeader(header); err != nil {
		t.Fatalf("failed to write tar header: %v", err)
	}

	if _, err := tarWriter.Write([]byte("hello")); err != nil {
		t.Fatalf("failed to write tar body: %v", err)
	}

	if err := tarWriter.Close(); err != nil {
		t.Fatalf("failed to close tar: %v", err)
	}

	if err := gzWriter.Close(); err != nil {
		t.Fatalf("failed to close gzip: %v", err)
	}

	// Start a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/gzip")
		w.Write(buf.Bytes())
	}))
	defer server.Close()

	// Create temp directory for extraction
	tmpDir := t.TempDir()

	// Test download and extract
	binaryPath, err := downloadAndExtract(server.URL, tmpDir)
	if err != nil {
		t.Fatalf("downloadAndExtract failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(binaryPath); err != nil {
		t.Fatalf("binary file not found: %v", err)
	}

	// Verify content
	content, err := os.ReadFile(binaryPath)
	if err != nil {
		t.Fatalf("failed to read binary: %v", err)
	}

	if string(content) != "hello" {
		t.Errorf("expected 'hello', got '%s'", string(content))
	}

	// Verify permissions
	info, _ := os.Stat(binaryPath)
	if info.Mode()&0111 == 0 {
		t.Errorf("binary is not executable")
	}
}

func TestFetchLatestReleaseFromURL(t *testing.T) {
	// Create a test server that returns a release
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		release := Release{TagName: "v1.0.0"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(release)
	}))
	defer server.Close()

	release, err := fetchLatestReleaseFromURL(server.URL)
	if err != nil {
		t.Fatalf("fetchLatestReleaseFromURL failed: %v", err)
	}

	if release.TagName != "v1.0.0" {
		t.Errorf("expected 'v1.0.0', got '%s'", release.TagName)
	}
}

func TestFetchLatestReleaseFromURL_NetworkError(t *testing.T) {
	_, err := fetchLatestReleaseFromURL("http://invalid-url-that-does-not-exist:99999")
	if err == nil {
		t.Errorf("expected error for network failure")
	}
}

func TestReplaceBinary_RequiresSudo(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("test should run as non-root user")
	}

	tmpSrc := t.TempDir()
	srcFile := filepath.Join(tmpSrc, "claudeenv")
	if err := os.WriteFile(srcFile, []byte("test"), 0755); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	err := replaceBinary(srcFile, "/usr/local/bin/claudeenv")
	if err == nil {
		t.Errorf("expected error when not running as root")
	}
}

func TestExtractGzipTar(t *testing.T) {
	// Create a gzip tar archive manually
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	tarWriter := tar.NewWriter(gzWriter)

	testContent := []byte("test binary content")
	header := &tar.Header{
		Name: "claudeenv",
		Size: int64(len(testContent)),
		Mode: 0755,
	}

	if err := tarWriter.WriteHeader(header); err != nil {
		t.Fatalf("WriteHeader error: %v", err)
	}

	if _, err := tarWriter.Write(testContent); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	tarWriter.Close()
	gzWriter.Close()

	// Verify we can extract it
	reader := bytes.NewReader(buf.Bytes())
	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		t.Fatalf("gzip.NewReader error: %v", err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)
	header2, err := tarReader.Next()
	if err != nil {
		t.Fatalf("tar.Next error: %v", err)
	}

	if header2.Name != "claudeenv" {
		t.Errorf("expected 'claudeenv', got '%s'", header2.Name)
	}

	content, err := io.ReadAll(tarReader)
	if err != nil {
		t.Fatalf("ReadAll error: %v", err)
	}

	if !bytes.Equal(content, testContent) {
		t.Errorf("content mismatch: got %v", content)
	}
}
