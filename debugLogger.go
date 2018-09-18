package gtbackend

import "github.com/Sirupsen/logrus"

// DebugLog - more verbose and understandable log output
func DebugLog(fileName string, method string, err error) {
	logrus.Debug(fileName, ", method: ", method, ", ", err)
}
