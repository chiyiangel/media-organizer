package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

type LoggerOptions struct {
	LogPath   string // 日志文件路径
	QuietMode bool   // 安静模式，只输出到文件
}

type Logger struct {
	*log.Logger
	logFile *os.File
	level   LogLevel
	quiet   bool
}

func NewLogger(opts *LoggerOptions) *Logger {
	if opts == nil {
		opts = &LoggerOptions{}
	}

	// 确定日志文件路径
	logPath := opts.LogPath
	if logPath == "" {
		logDir := "logs"
		if err := os.MkdirAll(logDir, 0755); err != nil {
			log.Fatal("无法创建日志目录:", err)
		}
		logPath = filepath.Join(logDir, fmt.Sprintf("media-organizer-%s.log",
			time.Now().Format("2006-01-02-15-04-05")))
	} else {
		// 确保日志文件的目录存在
		logDir := filepath.Dir(logPath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			log.Fatal("无法创建日志目录:", err)
		}
	}

	// 打开日志文件
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("无法创建日志文件:", err)
	}

	// 根据安静模式决定输出位置
	var writer io.Writer
	if opts.QuietMode {
		writer = logFile
	} else {
		writer = io.MultiWriter(os.Stdout, logFile)
	}

	return &Logger{
		Logger:  log.New(writer, "", log.Ldate|log.Ltime),
		logFile: logFile,
		level:   LevelInfo,
		quiet:   opts.QuietMode,
	}
}

func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *Logger) log(level LogLevel, prefix string, v ...interface{}) {
	if level >= l.level {
		l.Printf("%s %s", prefix, fmt.Sprint(v...))
	}
}

func (l *Logger) logf(level LogLevel, prefix, format string, v ...interface{}) {
	if level >= l.level {
		l.Printf("%s %s", prefix, fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Debug(v ...interface{}) {
	l.log(LevelDebug, "[DEBUG]", v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.logf(LevelDebug, "[DEBUG]", format, v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.log(LevelInfo, "[INFO]", v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.logf(LevelInfo, "[INFO]", format, v...)
}

func (l *Logger) Warn(v ...interface{}) {
	l.log(LevelWarn, "[WARN]", v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.logf(LevelWarn, "[WARN]", format, v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.log(LevelError, "[ERROR]", v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logf(LevelError, "[ERROR]", format, v...)
}

func (l *Logger) Fatal(v ...interface{}) {
	l.log(LevelFatal, "[FATAL]", v...)
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.logf(LevelFatal, "[FATAL]", format, v...)
	os.Exit(1)
}

func (l *Logger) Close() error {
	return l.logFile.Close()
}
