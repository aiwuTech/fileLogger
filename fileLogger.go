// Package: fileLogger
// File: fileLogger.go
// Created by: mint(mint.zhao.chiu@gmail.com)_aiwuTech
// Useage:
// DATE: 14-8-23 17:20
package fileLogger

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	DATEFORMAT         = "2006-01-02"
	DEFAULT_FILE_COUNT = 10
	DEFAULT_FILE_SIZE  = 50
	DEFAULT_FILE_UNIT  = MB
	DEFAULT_LOG_SCAN   = 300
	DEFAULT_LOG_SEQ    = 5000
	DEFAULT_LOG_LEVEL  = TRACE
)

type UNIT int64

const (
	_       = iota
	KB UNIT = 1 << (iota * 10)
	MB
	GB
	TB
)

type SplitType byte

const (
	SplitType_Size SplitType = iota
	SplitType_Daily
)

type LEVEL byte

const (
	TRACE LEVEL = iota
	INFO
	WARN
	ERROR
	OFF
)

type FileLogger struct {
	splitType SplitType
	mu        *sync.RWMutex
	fileDir   string
	fileName  string
	suffix    int
	fileCount int
	fileSize  int64
	prefix    string

	date *time.Time

	logFile *os.File
	lg      *log.Logger

	logScan int64

	logChan chan string

	logLevel   LEVEL
	logConsole bool
}

// NewDefaultLogger return a logger split by fileSize by default
func NewDefaultLogger(fileDir, fileName string) *FileLogger {
	return NewSizeLogger(fileDir, fileName, "",
		DEFAULT_FILE_COUNT, DEFAULT_FILE_SIZE, DEFAULT_FILE_UNIT, DEFAULT_LOG_SCAN, DEFAULT_LOG_SEQ)
}

// NewSizeLogger return a logger split by fileSize
// Parameters:
// 		file directory
// 		file name
// 		log's prefix
// 		fileCount holds maxCount of bak file
//		fileSize holds each of bak file's size
// 		unit stands for kb, mb, gb, tb
//		logScan after a logScan time will check fileLogger isMustSplit, default is 300s
func NewSizeLogger(fileDir, fileName, prefix string, fileCount int, fileSize int64, unit UNIT,
	logScan int64, logSeq int) *FileLogger {
	sizeLogger := &FileLogger{
		splitType:  SplitType_Size,
		mu:         new(sync.RWMutex),
		fileDir:    fileDir,
		fileName:   fileName,
		fileCount:  fileCount,
		fileSize:   fileSize * int64(unit),
		prefix:     prefix,
		logScan:    logScan,
		logChan:    make(chan string, logSeq),
		logLevel:   DEFAULT_LOG_LEVEL,
		logConsole: false,
	}

	sizeLogger.initLogger()

	return sizeLogger
}

// NewDailyLogger return a logger split by daily
// Parameters:
// 		file directory
// 		file name
// 		log's prefix
func NewDailyLogger(fileDir, fileName, prefix string, logScan int64, logSeq int) *FileLogger {
	dailyLogger := &FileLogger{
		splitType:  SplitType_Daily,
		mu:         new(sync.RWMutex),
		fileDir:    fileDir,
		fileName:   fileName,
		prefix:     prefix,
		logScan:    logScan,
		logChan:    make(chan string, logSeq),
		logLevel:   DEFAULT_LOG_LEVEL,
		logConsole: false,
	}

	dailyLogger.initLogger()

	return dailyLogger
}

func (f *FileLogger) initLogger() {

	switch f.splitType {
	case SplitType_Size:
		f.initLoggerBySize()
	case SplitType_Daily:
		f.initLoggerByDaily()
	}

}

// init filelogger split by fileSize
func (f *FileLogger) initLoggerBySize() {

	f.mu.Lock()
	defer f.mu.Unlock()

	logFile := joinFilePath(f.fileDir, f.fileName)
	for i := 1; i <= f.fileCount; i++ {
		if !isExist(logFile + "." + strconv.Itoa(i)) {
			break
		}

		f.suffix = i
	}

	if !f.isMustSplit() {
		if !isExist(f.fileDir) {
			os.Mkdir(f.fileDir, 0755)
		}
		f.logFile, _ = os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		f.lg = log.New(f.logFile, f.prefix, log.LstdFlags|log.Lmicroseconds)
	} else {
		f.split()
	}

	go f.logWriter()
	go f.fileMonitor()
}

// init fileLogger split by daily
func (f *FileLogger) initLoggerByDaily() {

	t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))

	f.date = &t
	f.mu.Lock()
	defer f.mu.Unlock()

	logFile := joinFilePath(f.fileDir, f.fileName)
	if !f.isMustSplit() {
		if !isExist(f.fileDir) {
			os.Mkdir(f.fileDir, 0755)
		}
		f.logFile, _ = os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		f.lg = log.New(f.logFile, f.prefix, log.LstdFlags|log.Lmicroseconds)
	} else {
		f.split()
	}

	go f.logWriter()
	go f.fileMonitor()
}

// used for determine the fileLogger f is time to split.
// size: once the current fileLogger's fileSize >= config.fileSize need to split
// daily: once the current fileLogger stands for yesterday need to split
func (f *FileLogger) isMustSplit() bool {

	switch f.splitType {
	case SplitType_Size:
		logFile := joinFilePath(f.fileDir, f.fileName)
		if f.fileCount > 1 {
			if fileSize(logFile) >= f.fileSize {
				return true
			}
		}
	case SplitType_Daily:
		t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
		if t.After(*f.date) {
			return true
		}
	}

	return false
}

// Split fileLogger
func (f *FileLogger) split() {

	logFile := joinFilePath(f.fileDir, f.fileName)

	switch f.splitType {
	case SplitType_Size:
		f.suffix = int(f.suffix%f.fileCount + 1)
		if f.logFile != nil {
			f.logFile.Close()
		}

		logFileBak := logFile + "." + strconv.Itoa(f.suffix)
		if isExist(logFileBak) {
			os.Remove(logFileBak)
		}
		os.Rename(logFile, logFileBak)

		f.logFile, _ = os.Create(logFile)
		f.lg = log.New(f.logFile, f.prefix, log.LstdFlags|log.Lmicroseconds)

	case SplitType_Daily:
		logFileBak := logFile + "." + f.date.Format(DATEFORMAT)
		if !isExist(logFileBak) && f.isMustSplit() {
			if f.logFile != nil {
				f.logFile.Close()
			}

			err := os.Rename(logFile, logFileBak)
			if err != nil {
				f.lg.Printf("FileLogger rename error: %v", err.Error())
			}

			t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
			f.date = &t
			f.logFile, _ = os.Create(logFile)
			f.lg = log.New(f.logFile, f.prefix, log.LstdFlags|log.Lmicroseconds)
		}
	}
}

// After some interval time, goto check the current fileLogger's size or date
func (f *FileLogger) fileMonitor() {
	defer func() {
		if err := recover(); err != nil {
			f.lg.Printf("FileLogger's FileMonitor() catch panic: %v\n", err)
		}
	}()

	//TODO  load logScan interval from config file
	logScan := DEFAULT_LOG_SCAN

	timer := time.NewTicker(time.Duration(logScan) * time.Second)
	for {
		select {
		case <-timer.C:
			f.fileCheck()
		}
	}
}

// If the current fileLogger need to split, just split
func (f *FileLogger) fileCheck() {
	defer func() {
		if err := recover(); err != nil {
			f.lg.Printf("FileLogger's FileCheck() catch panic: %v\n", err)
		}
	}()

	if f.isMustSplit() {
		f.mu.Lock()
		defer f.mu.Unlock()

		f.split()
	}
}

// passive to close fileLogger
func (f *FileLogger) Close() error {

	close(f.logChan)
	f.lg = nil

	return f.logFile.Close()
}
