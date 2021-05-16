package common

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/sirupsen/logrus"
)

func Fetch(l *logrus.Logger, url string, target interface{}) error {
	res, err := httpClient(l, url)
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

	res, err := httpClient(l, url)
	if err != nil {
		return err
	}

	if res == nil {
		return nil
	}

	_, err = io.Copy(outFile, res)

	return err
}

func httpClient(l *logrus.Logger, url string) (io.ReadCloser, error) {
	c := retryablehttp.NewClient()
	c.RetryMax = 3
	c.Logger = l

	res, err := c.Get(url)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}
