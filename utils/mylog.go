package utils

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// MyWriter is a customized wrapper for io.Writer used to support custom logging.
type MyWriter struct {
	Writer io.Writer
}

// Write is a wrapper for io.Write that only operates if the trace flag is on.
func (w MyWriter) Write(p []byte) (int, error) {
	if isTrace != 0 {
		return w.Writer.Write(p)
	}
	return ioutil.Discard.Write(p)
}

// These are four logging streams from least to most serious.
var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

type logStreamType int

// These are log stream output destinations.
const (
	LogNoChange logStreamType = -1
	LogNil      logStreamType = 0
	LogStdout   logStreamType = 1
	LogStderr   logStreamType = 2
	LogBoth     logStreamType = 3
	LogFile     logStreamType = 4
	LogTrace    logStreamType = 5
)

var isTrace int // local flag (0 or 1) to enable disable Trace

// SetTrace sets trace on and off (1, 0)
func SetTrace(intEnable int) {
	isTrace = intEnable
}

// GetTrace returns trace flag (1, 0)
func GetTrace() int {
	return isTrace
}

// InitLog assigns logging stream types to output streams.
func InitLog(traceCode logStreamType, infoCode logStreamType, warningCode logStreamType, errorCode logStreamType) {

	traceHandle := AssignHandle(traceCode)
	infoHandle := AssignHandle(infoCode)
	warningHandle := AssignHandle(warningCode)
	errorHandle := AssignHandle(errorCode)

	Trace = log.New(traceHandle, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(infoHandle, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(warningHandle, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(errorHandle, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// AssignHandle maps a log stream type to a writer handle
func AssignHandle(logCode logStreamType) io.Writer {
	switch logCode {
	case LogNil:
		return ioutil.Discard
	case LogStdout:
		return os.Stdout
	case LogStderr:
		return os.Stderr
	case LogTrace:
		tracer := MyWriter{
			Writer: os.Stdout,
		}
		return tracer
	case LogBoth:
		multi := CreateMultiLog("logs/errlog.txt")
		return multi
	default:
		return ioutil.Discard
	}
}

// LogPkg is struct for capturing json to reassign a log stream
type LogPkg struct {
	Log    string `json:"log"`
	Writer string `json:"writer"`
	Path   string `json:"path,omitempty"`
}

// ReassignLog re-assigns a log stream type to an output stream.
func ReassignLog(pkg LogPkg) bool {
	var logCode logStreamType

	pkg.Log = strings.ToUpper(pkg.Log)
	pkg.Writer = strings.ToLower(pkg.Writer)

	switch pkg.Writer {
	case "nil":
		logCode = LogNil
	case "stdout":
		logCode = LogStdout
	case "stdErr":
		logCode = LogStderr
	case "both":
		logCode = LogBoth
	default:
		return false
	}

	strPrefix := pkg.Log + ": "
	handle := AssignHandle(logCode)
	logger := log.New(handle, strPrefix, log.Ldate|log.Ltime|log.Lshortfile)

	switch pkg.Log {
	case "TRACE":
		Trace = logger
	case "INFO":
		Info = logger
	case "WARNING":
		Warning = logger
	case "ERROR":
		Error = logger
	}
	return true
}

// CreateMultiLog implements tee of multiple output writers
func CreateMultiLog(strPath string) io.Writer {
	fileLog, err := os.OpenFile(strPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", strPath, ":", err)
		return ioutil.Discard
	} else {
		multi := io.MultiWriter(fileLog, os.Stdout)
		return multi
	}
}

// CreateFileLog assigns a logging stream to a physical file.
func CreateFileLog(strPath string) io.Writer {
	fileLog, err := os.OpenFile(strPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", strPath, ":", err)
		return ioutil.Discard
	} else {
		return fileLog
	}
}
