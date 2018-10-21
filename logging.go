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

// LoggingMiddleware - capture http.Handle
type LoggingMiddleware func(http.Handler) http.Handler

// InternalLog internal log entry
type InternalLog struct {
	Origin string
	Method string
	Err    error
	Msg    string
}

// CONNECTOR log file prefix
const CONNECTOR = "connector-log"

// APPLICATION log file prefix
const APPLICATION = "application-log"

// LOGS log directory name
const LOGS = "logs"

// LogIncomingRequests - Logs request traffic into our app
func LogIncomingRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(LOGS, CONNECTOR)
		writer, err := GetLogWriter(path)
		if err != nil {
			logrus.Error(err.Error())
		}

		writers := []io.Writer{os.Stderr, writer}
		logRequest(writers, r, "Incoming traffic")

		next.ServeHTTP(w, r)
	})
}

// DebugLog - Log an error internally, will contain implementation specific information
func DebugLog(entry InternalLog) {
	abs, _ := filepath.Abs("../" + LOGS)
	path := filepath.Join(abs, APPLICATION)
	writer, err := GetLogWriter(path)
	if err != nil {
		logrus.Error(err.Error())
	}

	writers := []io.Writer{os.Stderr, writer}
	logInternal(entry, writers, "Error thrown")
}

// FatalLog log to writers and exit app
func FatalLog(entry InternalLog) {
	path := LOGS + "/" + APPLICATION
	writer, err := GetLogWriter(path)
	if err != nil {
		logrus.Error(err.Error())
	}

	writers := []io.Writer{os.Stderr, writer}

	logInternal(entry, writers, "Fatal")
	logrus.Fatal("Shutting down")
}

// GetLogWriter will return a file log writer
func GetLogWriter(path string) (writer *rotatelogs.RotateLogs, err error) {
	writer, err = rotatelogs.New(
		path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(604800)*time.Second),
	)

	return
}

// logRequest message - logs incoming traffic
func logRequest(writers []io.Writer, r *http.Request, msg string) {
	logger := logrus.New()

	entry := logger.WithFields(logrus.Fields{
		"Address": r.RequestURI,
		"method":  r.Method,
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
