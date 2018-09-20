package gtbackend

import (
	"github.com/GlidingTracks/gt-backend/constant"
	"os"
	"path/filepath"
	"strings"
)

// SaveFileToLocalStorage - Save the uploaded file in the filesystem. Path: .Records/{uId}/
func SaveFileToLocalStorage(uid string, fileNameRaw string) (file *os.File, err error) {
	path := createFilePath(constant.LSRoot, uid)
	os.MkdirAll(path, os.ModePerm)

	// CleanedFileName
	cfn := cleanFilePath(fileNameRaw)
	fileName := path + constant.Slash + cfn

	file, err = os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)

	return
}

// DeleteFileFromLocalStorage takes a uid and a filename and deletes a file determined
// by the params
func DeleteFileFromLocalStorage(uid string, fileName string) (err error) {
	path := createFilePath(constant.LSRoot, uid, fileName)

	err = os.Remove(path)
	if err != nil {
		return
	}

	return
}

// createFilePath - Method for creating a new path OS independent
func createFilePath(args ...string) (path string) {
	for _, k := range args {
		path = filepath.Join(path, k)
	}

	return
}

// cleanFilePath - If the user has supplied a filename with already existing filepath, clean it up
// and return only the filename
func cleanFilePath(filePath string) (fileName string) {
	parts := strings.Split(filePath, constant.Slash)

	fileName = parts[len(parts)-1]
	return
}

