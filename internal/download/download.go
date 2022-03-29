package download

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var (
	cacheDir = os.TempDir()
)

func tmpFile() (*os.File, error) {
	return os.CreateTemp(cacheDir, "")
}

func fileForUrl(url string) (string, error) {
	h := md5.New()
	io.WriteString(h, url)
	if _, err := os.Stat("cache"); os.IsNotExist(err) {
		os.Mkdir("cache", os.ModePerm)
	}
	return filepath.Abs("./cache/" + fmt.Sprintf("%x.html", h.Sum(nil)))
}

func Download(url string) (string, error) {
	fpath, err := fileForUrl(url)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(fpath); err == nil {
		return fpath, nil
	} else if !os.IsNotExist(err) {
		return "", err
	}

	f, err := tmpFile()
	if err != nil {
		return "", err
	}

	if res, err := http.Get(url); err != nil {
		return "", nil
	} else {
		defer res.Body.Close()
		if res.StatusCode == http.StatusOK {
			if _, err := io.Copy(f, res.Body); err != nil {
				return "", err
			}
		} else {
			return "", errors.New(res.Status)
		}
	}
	if err := os.Rename(f.Name(), fpath); err != nil {
		return "", err
	}

	return fpath, nil
}
