package watcher

import (
	"io"
	"log"
	"os"
)

// Logger logger interface
type Logger interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

var (
	// DefaultLogger is used when Config.Logger == nil
	DefaultLogger Logger = log.New(os.Stderr, "", log.LstdFlags)
	// DiscardLogger can be used to disable logging output
	DiscardLogger Logger = log.New(io.Discard, "", 0)
)
