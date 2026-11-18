package logexport

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExport(t *testing.T) {
	base := t.TempDir()
	if err := os.MkdirAll(filepath.Join(base, "logs", "applog"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(base, "logs", "applog", "access.log"), []byte("hello log\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(base, "etc"), 0o755); err != nil {
		t.Fatal(err)
	}
	cfg := "ChatGPT:\n  APIToken: sk-secret\n  APIUrl: https://example.com\nRoomId: 123\n"
	if err := os.WriteFile(filepath.Join(base, "etc", "bilidanmaku-api.yaml"), []byte(cfg), 0o644); err != nil {
		t.Fatal(err)
	}

	out := t.TempDir()
	zipPath, err := Export(base, out)
	if err != nil {
		t.Fatalf("Export() error = %v", err)
	}
	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		t.Fatalf("open zip: %v", err)
	}
	defer zr.Close()

	foundLog, foundConfig := false, false
	for _, f := range zr.File {
		rc, err := f.Open()
		if err != nil {
			t.Fatalf("open %s: %v", f.Name, err)
		}
		data, _ := io.ReadAll(rc)
		_ = rc.Close()
		switch f.Name {
		case "logs/applog/access.log":
			foundLog = true
			if string(data) != "hello log\n" {
				t.Fatalf("unexpected log content %q", string(data))
			}
		case "config-redacted.yaml":
			foundConfig = true
			text := string(data)
			if strings.Contains(text, "sk-secret") {
				t.Fatal("API token was not redacted")
			}
			if !strings.Contains(text, "APIToken: ***") || !strings.Contains(text, "RoomId: 123") {
				t.Fatalf("unexpected redacted config: %s", text)
			}
		}
	}
	if !foundLog || !foundConfig {
		t.Fatalf("missing zip entries, log=%v config=%v", foundLog, foundConfig)
	}
}
