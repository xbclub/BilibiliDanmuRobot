package logexport

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Export 把 baseDir/logs 打包到 outDir 下的 zip，并附带一份脱敏配置。
// baseDir 通常是程序工作目录，outDir 通常是 exe 所在目录。
func Export(baseDir, outDir string) (string, error) {
	logDir := filepath.Join(baseDir, "logs")
	info, err := os.Stat(logDir)
	if err != nil || !info.IsDir() {
		return "", fmt.Errorf("日志目录不存在: %s", logDir)
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return "", err
	}
	zipPath := filepath.Join(outDir, "BilibiliDanmuRobot-logs-"+time.Now().Format("20060102-150405")+".zip")
	f, err := os.Create(zipPath)
	if err != nil {
		return "", err
	}
	zw := zip.NewWriter(f)
	walkErr := filepath.Walk(logDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		rel, err := filepath.Rel(baseDir, path)
		if err != nil {
			return err
		}
		hdr, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		hdr.Name = filepath.ToSlash(rel)
		hdr.Method = zip.Deflate
		w, err := zw.CreateHeader(hdr)
		if err != nil {
			return err
		}
		src, err := os.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()
		_, err = io.Copy(w, src)
		return err
	})
	if walkErr != nil {
		_ = zw.Close()
		_ = f.Close()
		_ = os.Remove(zipPath)
		return "", walkErr
	}
	if cfg, err := os.ReadFile(filepath.Join(baseDir, "etc", "bilidanmaku-api.yaml")); err == nil {
		if w, err := zw.Create("config-redacted.yaml"); err == nil {
			_, _ = w.Write([]byte(redactConfig(string(cfg))))
		}
	}
	if err := zw.Close(); err != nil {
		_ = f.Close()
		_ = os.Remove(zipPath)
		return "", err
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(zipPath)
		return "", err
	}
	abs, err := filepath.Abs(zipPath)
	if err != nil {
		return zipPath, nil
	}
	return abs, nil
}

func redactConfig(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		low := strings.ToLower(line)
		if strings.Contains(low, "apitoken") || strings.Contains(low, "sessdata") || strings.Contains(low, "bili_jct") || strings.Contains(low, "cookie") {
			if idx := strings.Index(line, ":"); idx >= 0 {
				lines[i] = line[:idx+1] + " ***"
			}
		}
	}
	return strings.Join(lines, "\n")
}
