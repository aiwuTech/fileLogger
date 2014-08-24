// Package: example
// File: example1.go
// Created by: mint(mint.zhao.chiu@gmail.com)_aiwuTech
// Useage: example
// DATE: 14-8-24 14:08
package main

import (
	"github.com/aiwuTech/fileLogger"
)

func main() {

	TRACE := fileLogger.NewDefaultLogger("/usr/local/aiwuTech/log", "trace.log")
	TRACE.SetPrefix("[TRACE] ")

	TRACE.Println("This is the first TRACE log using fileLogger that written by aiwuTech.")

	//because the fileLogger is asynchronous, so the main goroutine must be watting for logWriter
	c := make(chan bool)
	_ = <- c
}

