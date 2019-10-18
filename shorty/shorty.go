package shorty

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Config is config
type Config struct {
	BaseURL string `json:"baseUrl"`
}

type shortenResponse struct {
	Error  *string `json:"error"`
	Result *struct {
		Suffix string `json:"suffix"`
		URL    string `json:"url"`
	} `json:"result"`
}

// Shorten takes a URL, shortens it and returns the
// shortened link.
func Shorten(longURL string, conf *Config) (string, error) {
	form := url.Values{}
	form.Add("url", longURL)

	req, err := http.NewRequest("POST", conf.BaseURL+"/api/shorten", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var resp shortenResponse
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&resp)
	if err != nil {
		return "", err
	}

	if resp.Error != nil {
		return "", errors.New(*resp.Error)
	}

	return fmt.Sprintf("%s/%s", conf.BaseURL, resp.Result.Suffix), nil
}
