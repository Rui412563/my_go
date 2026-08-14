package logwriter

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type SmartRotationWriter struct {
	mu          sync.Mutex
	basePath    string
	maxSize     int64
	currentFile *os.File
	currentSize int64
}

func NewSmartRotationWriter(basePath string, maxSizeMB int) *SmartRotationWriter {
	return &SmartRotationWriter{
		basePath: basePath,
		maxSize:  int64(maxSizeMB) * 1024 * 1024,
	}
}

func (w *SmartRotationWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// 触发切分：当前无文件 或 加上这行日志后超过 5MB
	if w.currentFile == nil || w.currentSize+int64(len(p)) > w.maxSize {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = w.currentFile.Write(p)
	w.currentSize += int64(n)
	return n, err
}

func (w *SmartRotationWriter) rotate() error {
	if w.currentFile != nil {
		w.currentFile.Close()
	}

	now := time.Now()
	// 文件夹：logs/2026-08-11/10-01-56
	subDir := filepath.Join(now.Format("2006-01-02"), now.Format("15-04-05"))
	fullDirPath := filepath.Join(w.basePath, subDir)

	if err := os.MkdirAll(fullDirPath, 0755); err != nil {
		return err
	}

	// 文件名：app-2026-08-11-10-01-56.log
	fileName := fmt.Sprintf("app-%s.log", now.Format("2006-01-02-15-04-05"))
	f, err := os.OpenFile(filepath.Join(fullDirPath, fileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	w.currentFile = f
	w.currentSize = 0
	return nil
}

func (w *SmartRotationWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.currentFile != nil {
		return w.currentFile.Close()
	}
	return nil
}
