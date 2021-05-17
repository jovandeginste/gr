package common

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/mitchellh/ioprogress"
	"github.com/sirupsen/logrus"
)

func Fetch(l *logrus.Logger, url string, target interface{}) error {
	_, res, err := httpClient(l, url)
	if err != nil {
		return err
	}

	if res == nil {
		return nil
	}

	body, readErr := ioutil.ReadAll(res)
	if readErr != nil {
		return err
	}

	if err := json.Unmarshal(body, target); err != nil {
		return err
	}

	return nil
}

func Download(l *logrus.Logger, url string, target string) error {
	outFile, err := os.Create(target)
	if err != nil {
		return err
	}

	defer outFile.Close()

	s, res, err := httpClient(l, url)
	if err != nil {
		return err
	}

	l.Infof("Downloading '%s' (%d bytes)...", url, s)
	l.Infof("To '%s'...", target)

	if res == nil {
		return nil
	}

	if s <= 0 {
		_, err = io.Copy(outFile, res)

		return err
	}

	progressR := &ioprogress.Reader{
		Reader:       res,
		Size:         s,
		DrawFunc:     ioprogress.DrawTerminalf(os.Stdout, ioprogress.DrawTextFormatBar(100)),
		DrawInterval: 50 * time.Millisecond,
	}

	_, err = io.Copy(outFile, progressR)

	return err
}

func httpClient(l *logrus.Logger, url string) (int64, io.ReadCloser, error) {
	c := retryablehttp.NewClient()
	c.RetryMax = 3
	c.Logger = l

	res, err := c.Get(url)
	if err != nil {
		return 0, nil, err
	}

	contentLengthHeader := res.Header.Get("Content-Length")

	size, err := strconv.ParseInt(contentLengthHeader, 10, 64)
	if err != nil {
		size = 0
	}

	return size, res.Body, nil
}
