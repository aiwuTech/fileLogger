// Package: fileLogger
// File: util.go
// Created by: mint(mint.zhao.chiu@gmail.com)_aiwuTech
// Useage: some useful utils
// DATE: 14-8-23 17:03
package fileLogger

import (
	"os"
	"path/filepath"
)

// Determine a file or a path exists in the os
func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// joinFilePath joins path & file into a single string
func joinFilePath(path, file string) string {
	return filepath.Join(path, file)
}

// return length in bytes for regular files
func fileSize(file string) int64 {
	f, e := os.Stat(file)
	if e != nil {
		return 0
	}

	return f.Size()
}

// return file name without dir
func shortFileName(file string) string {
	return filepath.Base(file)
}
