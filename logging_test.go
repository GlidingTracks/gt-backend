package gtbackend

import (
	"github.com/gorilla/mux"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

var testInternalLog = internalLog{
	"test",
	"test",
	errors.New("test"),
	"test",
}

var testPath = "dump" + strconv.FormatInt(time.Now().Unix(), 10)

func MockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		w.WriteHeader(200)
	}
}

func TestGetAndSetConfig(t *testing.T) {
	oldC := GetLogConfig()

	SetLogConfig(LogConfig{
		testPath,
	})

	newC := GetLogConfig()

	if newC == oldC {
		t.Error("Config not changed")
	}

	SetLogConfigDefault()

	newC = GetLogConfig()

	if newC != oldC {
		t.Error("Default not sat")
	}
}

func TestLogIncomingRequests(t *testing.T) {
	SetLogConfig(LogConfig{
		testPath,
	})

	server := mux.NewRouter()
	server.Use(LogIncomingRequests)

	req, err := http.NewRequest("GET", "/getTracks", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	server.HandleFunc("/getTracks", MockHandler)

	server.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Error("Wrong code returned")
	}

	if !exists(LogPath) {
		t.Fatal("No log directory created")
	}

	os.RemoveAll(LogPath)
}

func TestDebugLog(t *testing.T) {
	SetLogConfig(LogConfig{
		testPath,
	})

	testInternalLogHeader := DebugLogPrepareHeader("test", "test")
	err := errors.New("test")
	DebugLogErrNoMsg(testInternalLogHeader, err)
	DebugLogErrMsg(testInternalLogHeader, err, "test")
	DebugLogNoErrMsg(testInternalLogHeader, "test")

	if !exists(LogPath) {
		t.Fatal("No log directory created")
	}

	os.RemoveAll(LogPath)
}

func TestGetLogWriter(t *testing.T) {
	path := filepath.Join(config.Path, testPath)
	expected, err := rotatelogs.New(
		path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(604800)*time.Second),
	)
	if err != nil {
		t.Error(err)
	}

	defer expected.Close()

	actual, err := GetLogWriter(path)
	if err != nil {
		t.Error(err)
	}

	defer actual.Close()

	if actual.CurrentFileName() != expected.CurrentFileName() {
		t.Errorf("Assertion failed, expected: %s, actual %s", expected.CurrentFileName(), actual.CurrentFileName())
	}

}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}
