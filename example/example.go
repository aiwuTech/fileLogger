// Package: example
// File: example.go
// Created by: mint(mint.zhao.chiu@gmail.com)_aiwuTech
// Useage: example
// DATE: 14-8-24 14:08
package main

import (
	"github.com/aiwuTech/fileLogger"
	"sync"
)

var (
	logFile *fileLogger.FileLogger
	TRACE   *fileLogger.FileLogger
	INFO    *fileLogger.FileLogger
	WARN    *fileLogger.FileLogger
	ERROR   *fileLogger.FileLogger
)

func init() {

	TRACE = fileLogger.NewDefaultLogger("/usr/local/aiwuTech/log", "trace.log")
	INFO = fileLogger.NewDefaultLogger("/usr/local/aiwuTech/log", "info.log")
	WARN = fileLogger.NewDefaultLogger("/usr/local/aiwuTech/log", "warn.log")
	ERROR = fileLogger.NewDefaultLogger("/usr/local/aiwuTech/log", "error.log")

	TRACE.SetPrefix("[TRACE] ")
	INFO.SetPrefix("[INFO] ")
	WARN.SetPrefix("[WARN] ")
	ERROR.SetPrefix("[ERROR] ")

	logFile = fileLogger.NewDefaultLogger("/usr/local/aiwuTech/log", "test.log")
	logFile.SetLogLevel(fileLogger.INFO) //trace log will not be print
}

func main() {

	wg := new(sync.WaitGroup)

	wg.Add(1)
	go traceTest(wg)

	wg.Add(1)
	go infoTest(wg)

	wg.Add(1)
	go warnTest(wg)

	wg.Add(1)
	go errorTest(wg)

	wg.Add(1)
	go logTest(wg)

	wg.Wait()

	TRACE.Close()
	INFO.Close()
	WARN.Close()
	ERROR.Close()
	logFile.Close()
}

func traceTest(wg *sync.WaitGroup) {
	for i := 1; i <= 1000; i++ {
		TRACE.Printf("This is the No[%v] TRACE log using fileLogger that written by aiwuTech.", i)
	}

	wg.Done()
}

func infoTest(wg *sync.WaitGroup) {
	for i := 1; i <= 1000; i++ {
		INFO.Printf("This is the No[%v] INFO log using fileLogger that written by aiwuTech.", i)
	}

	wg.Done()
}

func warnTest(wg *sync.WaitGroup) {
	for i := 1; i <= 1000; i++ {
		WARN.Printf("This is the No[%v] WARN log using fileLogger that written by aiwuTech.", i)
	}

	wg.Done()
}

func errorTest(wg *sync.WaitGroup) {
	for i := 1; i <= 1000; i++ {
		ERROR.Printf("This is the No[%v] ERROR log using fileLogger that written by aiwuTech.", i)
	}

	wg.Done()
}

func logTest(wg *sync.WaitGroup) {

	for i := 1; i <= 1000; i++ {
		logFile.T("This is the No[%v] TRACE log using fileLogger that written by aiwuTech.", i)
		logFile.I("This is the No[%v] INFO log using fileLogger that written by aiwuTech.", i)
		logFile.W("This is the No[%v] WARN log using fileLogger that written by aiwuTech.", i)
		logFile.E("This is the No[%v] ERROR log using fileLogger that written by aiwuTech.", i)
	}

	wg.Done()
}
