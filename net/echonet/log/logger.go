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
	Format     = "%s %s %s"
	LF         = "\n"
	FilePerm   = 0644
	LevelTrace = (1 << 5)
	LevelInfo  = (1 << 4)
	LevelWarn  = (1 << 3)
	LevelError = (1 << 2)
	LevelFatal = (1 << 1)
	LevelNone  = 0

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
	LevelTrace: "TRACE",
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
	LevelFatal: "FATAL",
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
		outputer: outputStdout}
	return logger
}

func outputStdout(file string, level LogLevel, msg string) (int, error) {
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

func output(outputLevel LogLevel, msgFormat string, msgArgs ...interface{}) int {
	if sharedLogger == nil {
		return 0
	}

	logLevel := sharedLogger.GetLevel()
	if (logLevel < outputLevel) || (logLevel <= LevelFatal) || (LevelTrace < logLevel) {
		return 0
	}

	t := time.Now()
	logDate := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	headerString := fmt.Sprintf("[%s]", getLogLevelString(outputLevel))
	logMsg := fmt.Sprintf(Format, logDate, headerString, fmt.Sprintf(msgFormat, msgArgs...))
	logMsgLen := len(logMsg)

	if 0 < logMsgLen {
		logMsgLen, _ = sharedLogger.outputer(sharedLogger.File, logLevel, logMsg)
	}

	return logMsgLen
}

func Trace(format string, args ...interface{}) int {
	return output(LevelTrace, format, args...)
}

func Info(format string, args ...interface{}) int {
	return output(LevelInfo, format, args...)
}

func Warn(format string, args ...interface{}) int {
	return output(LevelWarn, format, args...)
}

func Error(format string, args ...interface{}) int {
	return output(LevelError, format, args...)
}

func Fatal(format string, args ...interface{}) int {
	return output(LevelFatal, format, args...)
}

func Output(outputLevel LogLevel, format string, args ...interface{}) int {
	return output(outputLevel, format, args...)
}
