package cachemgr

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type imageDownloader struct {
	RanobeProvider

	UniqueName string
	Url        string
	Filename   string
}

func cImageDownloader(ranobeProvider RanobeProvider, uniqueName string, url string, filename string) *imageDownloader {
	return &imageDownloader{
		RanobeProvider: ranobeProvider,
		UniqueName:     uniqueName,
		Url:            url,
		Filename:       filename,
	}
}
func (self *imageDownloader) httpGet() (*http.Response, error) {
	if response, err := http.Get(self.Url); err != nil {
		return nil, err
	} else {
		if response.StatusCode == http.StatusTooManyRequests {
			defer response.Body.Close()

			time.Sleep(time.Second)
			return self.httpGet()
		}
		if response.StatusCode != http.StatusOK {
			defer response.Body.Close()

			return nil, fmt.Errorf("Status code not 200, %s", response.Status)
		}
		return response, nil
	}
}
func (self *imageDownloader) openFile() (*os.File, string, error) {
	if ranobeDir, err := ConstructPath(self.RanobeProvider, self.UniqueName); err != nil {
		return nil, "", err
	} else {

		path := filepath.Join(ranobeDir, self.Filename)
		file, err := os.Create(path)

		return file, path, err
	}
}
func (self *imageDownloader) saveImage(stream io.Reader) (string, error) {
	if file, path, err := self.openFile(); err != nil {
		return "", err
	} else {

		if _, err := io.Copy(file, stream); err != nil {
			return "", err
		} else {
			return path, nil
		}
	}
}
func (self *imageDownloader) Download() (string, error) {
	if response, err := self.httpGet(); err != nil {
		return "", err
	} else {
		defer response.Body.Close()

		return self.saveImage(response.Body)
	}
}
func DownloadImage(ranobeProvider RanobeProvider, uniqueName string, url string, filename string) (string, error) {
	return cImageDownloader(
		ranobeProvider,
		uniqueName,
		url,
		filename,
	).Download()
}
