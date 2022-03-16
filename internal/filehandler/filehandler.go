package filehandler

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const MAX_UPLOAD_SIZE = 1024 * 1024 * 10 // 1GB

type HTTPMultiPartFormFileHandler struct {
	Httpmf   *http.Request
	Success  bool
	Err      error
	Filepath *string
}

func NewMultiPartHandler(r *http.Request) HTTPMultiPartFormFileHandler {
	return HTTPMultiPartFormFileHandler{Httpmf: r}
}

func (mf *HTTPMultiPartFormFileHandler) WriteHTTPMultiPartFormFile() {

	// 32 MB is the default used by FormFile
	if err := mf.Httpmf.ParseMultipartForm(32 << 20); err != nil {
		mf.Success = false
		mf.Filepath = nil
		mf.Err = err
	}

	files := mf.Httpmf.MultipartForm.File["file"]

	for _, fileHeader := range files {

		// checks for file size
		if fileHeader.Size > MAX_UPLOAD_SIZE {
			mf.Success = false
			mf.Filepath = nil
			mf.Err = fmt.Errorf("filesize exceed the max size of 1Gb, current filesize %d", fileHeader.Size)
		}

		// opens file for reading
		file, err := fileHeader.Open()
		if err != nil {
			mf.Success = false
			mf.Filepath = nil
			mf.Err = err
		}
		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {
				mf.Success = false
				mf.Filepath = nil
				mf.Err = err
			}
		}(file)

		// creates a buffer(byte Slice) of size 512 and Reads the file contents into it
		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			mf.Success = false
			mf.Filepath = nil
			mf.Err = err
		}

		// sets the position of the file handle to start
		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			mf.Success = false
			mf.Filepath = nil
			mf.Err = err
		}

		// creates directory in the given path if it doesn't exist
		err = os.MkdirAll("./data", os.ModePerm)
		if err != nil {
			mf.Success = false
			mf.Filepath = nil
			mf.Err = err
		}
		// creates a file with the given name in the given path
		f, err := os.Create(fmt.Sprintf("./data/%s-%d%s", strings.TrimSuffix(fileHeader.Filename,
			filepath.Ext(fileHeader.Filename)), time.Now().UnixMilli(), filepath.Ext(fileHeader.Filename)))
		if err != nil {
			mf.Success = false
			mf.Filepath = nil
			mf.Err = err
		}
		defer f.Close()

		// copies the content to the created file
		_, err = io.Copy(f, file)
		if err != nil {
			mf.Success = false
			mf.Filepath = nil
			mf.Err = err
		}
		fp := f.Name()
		mf.Success = true
		mf.Filepath = &fp
		mf.Err = nil
	}
}
