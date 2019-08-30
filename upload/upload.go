package upload

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// Uploader http client
type Uploader struct {
	APIServer string
	client    *http.Client
	Filename  string
}

// NewUploader create uploader
func NewUploader(filename string) *Uploader {

	return &Uploader{
		Filename: filename,
		client:   http.DefaultClient,
	}
}

// Execute upload recoding file to server
func (up *Uploader) Execute() {
	if err := up.upload(); err != nil {
		fmt.Printf("upload file error %s", err)
	}
}

func (up *Uploader) upload() error {
	file, err := os.Open(up.Filename)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(up.Filename))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return err
	}
	uploadURL := "http://127.0.0.1:8080/upload"
	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := up.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	println(string(bytes))
	return nil
}
