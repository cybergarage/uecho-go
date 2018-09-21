// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type LoggerOutpter func(file string, level LogLevel, msg string) (int, error)

type LogLevel int

type Logger struct {
	File     string
	Level    LogLevel
	outputer LoggerOutpter
}

const (
	Format           = "%s %s %s"
	LF               = "\n"
	FilePerm         = 0644
	LoggerLevelTrace = (1 << 5)
	LoggerLevelInfo  = (1 << 4)
	LoggerLevelWarn  = (1 << 3)
	LoggerLevelError = (1 << 2)
	LoggerLevelFatal = (1 << 1)
	LoggerLevelNone  = 0

	loggerLevelUnknownString = "UNKNOWN"
	loggerStdout             = "stdout"
)

var sharedLogger *Logger

func SetSharedLogger(logger *Logger) {
	sharedLogger = logger
}

func GetSharedLogger() *Logger {
	return sharedLogger
}

var logLevelStrings = map[LogLevel]string{
	LoggerLevelTrace: "TRACE",
	LoggerLevelInfo:  "INFO",
	LoggerLevelWarn:  "WARN",
	LoggerLevelError: "ERROR",
	LoggerLevelFatal: "FATAL",
}

func getLogLevelString(logLevel LogLevel) string {
	logString, hasString := logLevelStrings[logLevel]
	if !hasString {
		return loggerLevelUnknownString
	}
	return logString
}

func (logger *Logger) SetLevel(level LogLevel) {
	logger.Level = level
}

func (logger *Logger) GetLevel() LogLevel {
	return logger.Level
}

func NewStdoutLogger(level LogLevel) *Logger {
	logger := &Logger{
		File:     loggerStdout,
		Level:    level,
		outputer: outputStrount}
	return logger
}

func outputStrount(file string, level LogLevel, msg string) (int, error) {
	fmt.Println(msg)
	return len(msg), nil
}

func NewFileLogger(file string, level LogLevel) *Logger {
	logger := &Logger{
		File:     file,
		Level:    level,
		outputer: outputToFile}
	return logger
}

func outputToFile(file string, level LogLevel, msg string) (int, error) {
	msgBytes := []byte(msg + LF)
	fd, err := os.OpenFile(file, (os.O_WRONLY | os.O_CREATE | os.O_APPEND), FilePerm)
	if err != nil {
		return 0, err
	}

	writer := bufio.NewWriter(fd)
	writer.Write(msgBytes)
	writer.Flush()

	fd.Close()

	return len(msgBytes), nil
}

func output(outputLevel LogLevel, msg string) int {
	if sharedLogger == nil {
		return 0
	}

	logLevel := sharedLogger.GetLevel()
	if (logLevel < outputLevel) || (logLevel <= LoggerLevelFatal) || (LoggerLevelTrace < logLevel) {
		return 0
	}

	t := time.Now()
	logDate := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	headerString := fmt.Sprintf("[%s]", getLogLevelString(outputLevel))
	logMsg := fmt.Sprintf(Format, logDate, headerString, msg)
	logMsgLen := len(logMsg)

	if 0 < logMsgLen {
		logMsgLen, _ = sharedLogger.outputer(sharedLogger.File, logLevel, logMsg)
	}

	return logMsgLen
}

func Trace(msg string) int {
	return output(LoggerLevelTrace, msg)
}

func Info(msg string) int {
	return output(LoggerLevelInfo, msg)
}

func Warn(msg string) int {
	return output(LoggerLevelWarn, msg)
}

func Error(msg string) int {
	return output(LoggerLevelError, msg)
}

func Fatal(msg string) int {
	return output(LoggerLevelFatal, msg)
}

func Output(outputLevel LogLevel, msg string) int {
	return output(outputLevel, msg)
}
