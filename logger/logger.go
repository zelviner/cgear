package logger

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"text/template"
	"time"

	"github.com/ZEL-30/zel/logger/colors"
)

var errInvalidLogLevel = errors.New("logger: invaild log level")

// 日志级别
const (
	debugLevel = iota
	errorLevel
	fatalLevel
	criticalLevel
	successLevel
	warnLevel
	infoLevel
	hintLevel
)

var (
	sequenceNo uint64
	instance   *Logger
	once       sync.Once // Once是只执行一次动作的对象
)

var debugMode = os.Getenv("DEBUG_ENABLED") == "1"

var logLevel = infoLevel

// Logger将日志记录到指定的io.Writer
type Logger struct {
	mu     sync.Mutex
	output io.Writer
}

// LogRecord表示日志记录，并包含记录创建了一个递增的ID、级别和实际格式化的日志行。
type LogRecord struct {
	ID       string
	Level    string
	Message  string
	Filename string
	LineNo   int
}

var Log = GetLogger(os.Stdout)

var (
	logRecordTemplate      *template.Template
	debugLogRecordTemplate *template.Template
)

func GetLogger(w io.Writer) *Logger {
	once.Do(func() {
		var (
			err             error
			simpleLogFormat = `{{Now "2006/01/02 15:04:05"}} {{.Level}} ▶ {{.ID}} {{.Message}}{{EndLine}}`
			debugLogFormat  = `{{Now "2006/01/02 15:04:05"}} {{.Level}} ▶ {{.ID}} {{.Filename}}:{{.LineNo}} {{.Message}}{{EndLine}}`
		)

		// 初始化并且解析日志模板
		funcs := template.FuncMap{
			"Now":     Now,
			"EndLine": EndLine,
		}

		logRecordTemplate, err = template.New("simpleLogFormat").Funcs(funcs).Parse(simpleLogFormat)
		if err != nil {
			panic(err)
		}

		debugLogRecordTemplate, err = template.New("debugLogFormat").Funcs(funcs).Parse(debugLogFormat)
		if err != nil {
			panic(err)
		}

		instance = &Logger{output: colors.NewColorWriter(w)}
	})

	return instance
}

// SetOutput 设置记录器输出目标
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = colors.NewColorWriter(w)
}

// Now 返回指定布局中的当前本地时间
func Now(layout string) string {
	return time.Now().Format(layout)
}

// Endline 返回换行符
func EndLine() string {
	return "\n"
}

func (l *Logger) getLevalTag(level int) string {
	switch level {
	case debugLevel:
		return "DEBUG   "
	case errorLevel:
		return "ERROR   "
	case fatalLevel:
		return "FATAL   "
	case criticalLevel:
		return "CRITICAL"
	case successLevel:
		return "SUCCESS "
	case warnLevel:
		return "WARN    "
	case infoLevel:
		return "INFO    "
	case hintLevel:
		return "HINT    "
	default:
		panic(errInvalidLogLevel)
	}
}

func (l *Logger) getColorLevel(level int) string {
	switch level {
	case debugLevel:
		return colors.YellowBold(l.getLevalTag(level))
	case errorLevel:
		return colors.RedBold(l.getLevalTag(level))
	case fatalLevel:
		return colors.RedBold(l.getLevalTag(level))
	case criticalLevel:
		return colors.RedBold(l.getLevalTag(level))
	case successLevel:
		return colors.GreenBold(l.getLevalTag(level))
	case warnLevel:
		return colors.YellowBold(l.getLevalTag(level))
	case infoLevel:
		return colors.BlueBold(l.getLevalTag(level))
	case hintLevel:
		return colors.CyanBold(l.getLevalTag(level))
	default:
		panic(errInvalidLogLevel)
	}
}

// mustLog 根据指定的级别和参数记录消息, 如果出现错误，它会恐慌。
func (l *Logger) mustLog(level int, message string, args ...interface{}) {
	if level > logLevel {
		return
	}

	// 获取锁
	l.mu.Lock()
	defer l.mu.Unlock()

	// 创建日志记录并传入输出
	record := LogRecord{
		ID:      fmt.Sprintf("%04d", atomic.AddUint64(&sequenceNo, 1)),
		Level:   l.getColorLevel(level),
		Message: fmt.Sprintf(message, args...),
	}

	err := logRecordTemplate.Execute(l.output, record)
	if err != nil {
		panic(err)
	}
}

// mustLogistics 仅在启用调试模式时记录调试消息。即DEBUG_ENABLED=“1”
func (l *Logger) mustLogDebug(message string, file string, line int, args ...interface{}) {
	if !debugMode {
		return
	}

	// 将输出改为 Stderr
	l.SetOutput(os.Stderr)

	// 创建日志记录
	record := LogRecord{
		ID:       fmt.Sprintf("%04d", atomic.AddUint64(&sequenceNo, 1)),
		Level:    l.getColorLevel(debugLevel),
		Message:  fmt.Sprintf(message, args...),
		LineNo:   line,
		Filename: filepath.Base(file),
	}

	err := debugLogRecordTemplate.Execute(l.output, record)
	if err != nil {
		panic(err)
	}
}

// Debug 输出调试日志消息
func (l *Logger) Debug(message string, file string, line int) {
	l.mustLogDebug(message, file, line)
}

// Debugf 输出格式化的调试日志消息
func (l *Logger) Debugf(message string, file string, line int, vars ...interface{}) {
	l.mustLogDebug(message, file, line, vars...)
}

// Error 输出错误日志消息
func (l *Logger) Error(message string) {
	l.mustLog(errorLevel, message)
}

// Errorf 输出格式化的错误日志消息
func (l *Logger) Errorf(message string, vars ...interface{}) {
	l.mustLog(errorLevel, message, vars...)
}

// Fatal 输出致命日志消息并存在
func (l *Logger) Fatal(message string) {
	l.mustLog(fatalLevel, message)
	os.Exit(255)
}

// Fatalf 输出格式化的致命日志消息并存在
func (l *Logger) Fatalf(message string, vars ...interface{}) {
	l.mustLog(fatalLevel, message, vars...)
	os.Exit(255)
}

// Critival 输出关键日志消息
func (l *Logger) Critical(message string) {
	l.mustLog(criticalLevel, message)
}

// Criticalf 输出格式化的关键日志消息
func (l *Logger) Criticalf(message string, vars ...interface{}) {
	l.mustLog(criticalLevel, message, vars...)
}

// Success 输出成功日志消息
func (l *Logger) Success(message string) {
	l.mustLog(successLevel, message)
}

// Successf 输出格式化的成功日志消息
func (l *Logger) Successf(message string, vars ...interface{}) {
	l.mustLog(successLevel, message, vars...)
}

// Warn 输出警告日志消息
func (l *Logger) Warn(message string) {
	l.mustLog(warnLevel, message)
}

// Warnf 输出格式化的警告日志消息
func (l *Logger) Warnf(message string, vars ...interface{}) {
	l.mustLog(warnLevel, message, vars...)
}

// Info 输出信息日志消息
func (l *Logger) Info(message string) {
	l.mustLog(infoLevel, message)
}

// Infof 输出格式化的信息日志消息
func (l *Logger) Infof(message string, vars ...interface{}) {
	l.mustLog(infoLevel, message, vars...)
}

// Hint 输出提示日志消息
func (l *Logger) Hint(message string) {
	l.mustLog(hintLevel, message)
}

// Hintf 输出格式化的提示日志消息
func (l *Logger) Hintf(message string, vars ...interface{}) {
	l.mustLog(hintLevel, message, vars...)
}
