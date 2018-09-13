package main

import (
	"bytes"
	"net/http/httptest"
	"testing"
)

func TestProcessUploadRequestWrongContentType(t *testing.T) {


	req := httptest.NewRequest("POST", "/upload", bytes.NewBufferString("{\"uid\":\"123\",\"uploadfile\":\"Test.txt\"}"))
	req.Header.Set("Content-Type", "multipart/form-data")

	err, code := ProcessUploadRequest(req, "123")

	if err != nil && code != 400 {
		t.Error("Wrong file content type got through")
	}
}

func TestProcessUploadRequestWrongFileFormat(t *testing.T) {
	t.Skip("Not implemented")
}

func TestProcessUploadRequestNoUID(t *testing.T) {
	t.Skip("Not implemented")

}

func TestProcessUploadRequestPayload(t *testing.T) {
	t.Skip("Not implemented")
}