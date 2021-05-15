package common

import (
	"encoding/json"
	"io/ioutil"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

func Fetch(url string, target interface{}) error {
	c := retryablehttp.NewClient()
	c.RetryMax = 3

	res, err := c.Get(url)
	if err != nil {
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return err
	}

	if err := json.Unmarshal(body, target); err != nil {
		return err
	}

	return nil
}
