// Package: fileLogger
// File: writer.go
// Created by: mint(mint.zhao.chiu@gmail.com)_aiwuTech
// Useage: 
// DATE: 14-8-24 12:40
package fileLogger

import (
	"sync"
	"log"
	"fmt"
	"time"
	"os"
)

const (
	DEFAULT_PRINT_INTERVAL = 300
)

func (f *FileLogger) logWriter() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("FileLogger's LogWritter() catch panic: %v\n", err.Error())
		}
	}()

	//TODO let printInterVal can be configure, without change sourceCode
	printInterval := DEFAULT_PRINT_INTERVAL

	seqTimer := time.NewTicker(time.Duration(printInterval) * time.Second)
	for {
		select {
		case str := <- f.logChan:

			f.p(str)
		case <- seqTimer.C:
			log.Printf("================ LOG SEQ SIZE:%v ==================\n", len(f.logChan))
		}
	}
}

func (f *FileLogger) p(str string) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	f.lg.Output(2, str)
}

// Printf throw logstr to channel to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (f *FileLogger) Printf(format string, v ...interface {}) {
	f.logChan <- fmt.Sprintf(format, v...)
}

// Print throw logstr to channel to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (f *FileLogger) Print(v ...interface {}) {
	f.logChan <- fmt.Sprint(v...)
}

// Println throw logstr to channel to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func (f *FileLogger) Println(v ...interface {}) {
	f.logChan <- fmt.Sprintln(v...)
}

// Fatal is equivalent to f.Print() followed by a call to os.Exit(1).
func (f *FileLogger) Fatal(v ...interface {}) {
	f.logChan <- fmt.Sprint(v...)

	//TODO current goroutine fatal, other goroutine what to do? what about logstr already in logchan?
	os.Exit(1)
}

// Fatalf is equivalent to f.Printf() followed by a call to os.Exit(1).
func (f *FileLogger) Fatalf(format string, v ...interface {}) {
	f.logChan <- fmt.Sprintf(format, v...)

	//TODO current goroutine fatal, other goroutine what to do? what about logstr already in logchan?
	os.Exit(1)
}

// Fatalln is equivalent to f.Println() followed by a call to os.Exit(1).
func (f *FileLogger) Fatalln(v ...interface {}) {
	f.logChan <- fmt.Sprintln(v...)

	//TODO current goroutine fatal, other goroutine what to do? what about logstr already in logchan?
	os.Exit(1)
}

// Panic is equivalent to f.Print() followed by a call to panic().
func (f *FileLogger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	f.logChan <- s

	//TODO current goroutine panic, other goroutine what to do? what about logstr already in logchan?
	panic(s)
}

// Panicf is equivalent to f.Printf() followed by a call to panic().
func (f *FileLogger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	f.logChan <- s

	//TODO current goroutine panic, other goroutine what to do? what about logstr already in logchan?
	panic(s)
}

// Panicln is equivalent to f.Println() followed by a call to panic().
func (f *FileLogger) Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	f.logChan <- s

	//TODO current goroutine panic, other goroutine what to do? what about logstr already in logchan?
	panic(s)
}

