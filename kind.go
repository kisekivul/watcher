package watcher

import (
	"os"
)

// Kind kind of path
type Kind int

const (
	// UNKNOWN unknown
	UNKNOWN Kind = iota
	// FOLDER folder
	FOLDER
	// FILE file
	FILE
)

var (
	kinds = [...]string{
		"unknown",
		"folder",
		"file",
	}
)

// String return kind string
func (k Kind) String() string {
	return kinds[int(k)%len(kinds)]
}

// Exist exist stat
func Exist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// IsDir is folder
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile is file
func IsFile(path string) bool {
	return !IsDir(path)
}
