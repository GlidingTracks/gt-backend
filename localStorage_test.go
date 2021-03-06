package gtbackend

import (
	"github.com/GlidingTracks/gt-backend/constant"
	"os"
	"testing"
)

func TestLocalStorage(t *testing.T) {
	uid := "test"
	filePath := "./testdata/text.txt"
	fileName := "text.txt"

	t.Run("Save", func(t *testing.T) {
		f, err := os.Open(filePath)
		if err != nil {
			t.Error("Could not open test file: ", filePath)
		}

		_, _, err = SaveFileToLocalStorage(uid, filePath, f)
		if err != nil {
			t.Error("Could not save File")
		}

		s := checkFileExist(uid, fileName)
		if !s {
			t.Error("File not created")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		err := DeleteFileFromLocalStorage(uid, fileName)
		if err != nil {
			t.Error("Could not delete File")
		}

		s := checkFileExist(uid, fileName)
		if s {
			t.Error("File not deleted")
		}
	})
}

func checkFileExist(uid string, fileName string) (exist bool) {
	path := createFilePath(constant.LSRoot, uid, fileName)
	f, err := os.Stat(path)
	if err != nil {
		exist = false
		return
	}

	if f != nil {
		exist = true
		return
	}

	exist = false
	return
}
