package gtbackend

import (
	"github.com/Sirupsen/logrus"
	"github.com/lestrrat-go/file-rotatelogs"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// LoggingMiddleware - capture http.Handle.
type LoggingMiddleware func(http.Handler) http.Handler

// InternalLog internal log entry
type InternalLog struct {
	Origin string
	Method string
	Err    error
	Msg    string
}

// LogConfig enables users to dictate where log file are to be put.
type LogConfig struct {
	Path string
}

var config = LogConfig{
	LOGS,
}

// LogPath stores log directory path.
var LogPath string

// CONNECTOR log file prefix.
const CONNECTOR = "connector-log"

// APPLICATION log file prefix.
const APPLICATION = "application-log"

// LOGS log directory name.
const LOGS = "logs"

// LogIncomingRequests - Logs request traffic into our app.
func LogIncomingRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logIncomingRequests(r)
		next.ServeHTTP(w, r)
	})
}

// DebugLog - Log an error internally, will contain implementation specific information.
func DebugLog(entry InternalLog) {
	abs, _ := filepath.Abs(config.Path)
	setLogPath(abs)

	path := filepath.Join(abs, APPLICATION)

	writer, err := GetLogWriter(path)
	if err != nil {
		logrus.Error(err.Error())
	}

	writers := []io.Writer{os.Stderr, writer}
	logInternal(entry, writers, "Error thrown")
}

// LogFatal log to writers and exit app.
func LogFatal(entry InternalLog) {
	abs, _ := filepath.Abs(config.Path)
	setLogPath(abs)

	path := filepath.Join(abs, APPLICATION)

	writer, err := GetLogWriter(path)
	if err != nil {
		logrus.Error(err.Error())
	}

	writers := []io.Writer{os.Stderr, writer}

	logInternal(entry, writers, "Fatal")
	logrus.Fatal("Shutting down")
}

// GetLogWriter will return a file log writer.
func GetLogWriter(path string) (writer *rotatelogs.RotateLogs, err error) {
	writer, err = rotatelogs.New(
		path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(604800)*time.Second),
	)

	return
}

// GetLogConfig fetch logger config.
func GetLogConfig() LogConfig {
	return config
}

// SetLogConfig sets logger config.
func SetLogConfig(newConfig LogConfig) {
	config = newConfig
}

// SetLogConfigDefault restore logger initial config.
func SetLogConfigDefault() {
	config = LogConfig{
		LOGS,
	}
}

func logIncomingRequests(r *http.Request) {
	abs, _ := filepath.Abs(config.Path)
	setLogPath(abs)

	path := filepath.Join(abs, CONNECTOR)

	writer, err := GetLogWriter(path)
	if err != nil {
		logrus.Error(err.Error())
	}

	writers := []io.Writer{os.Stderr, writer}
	logRequest(writers, r, "Incoming traffic")
}

// logRequest message - logs incoming traffic
func logRequest(writers []io.Writer, r *http.Request, msg string) {
	logger := logrus.New()

	entry := logger.WithFields(logrus.Fields{
		"Address":   r.RequestURI,
		"method":    r.Method,
		"multiform": r.MultipartForm,
		"body":      r.Body,
		"form":      r.Form,
	})

	for i := range writers {
		setFormat(writers[i], logger)

		logger.SetOutput(writers[i])
		entry.Info(msg)
	}
}

// logInternal message - usually a error
func logInternal(entry InternalLog, writers []io.Writer, msg string) {
	logger := logrus.New()

	logEntry := logger.WithFields(logrus.Fields{
		"origin": entry.Origin,
		"method": entry.Method,
		"err":    entry.Err,
		"msg":    entry.Msg,
	})

	for i := range writers {
		setFormat(writers[i], logger)

		logger.SetOutput(writers[i])
		logEntry.Error(msg)
	}
}

// setFormat determine if logger should have color or not
func setFormat(writer io.Writer, logger *logrus.Logger) {
	if writer != os.Stderr {
		logger.SetFormatter(&logrus.TextFormatter{
			DisableColors: true,
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			DisableColors: false,
		})
	}
}

func setLogPath(path string) {
	LogPath = path
}
