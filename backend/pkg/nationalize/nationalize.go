package nationalize

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
)

const (
	Url = "https://api.nationalize.io/"
)

type Client struct {
	httpClient *http.Client
}

func New() (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	return &Client{
		httpClient: &http.Client{
			Jar: jar,
		},
	}, nil
}

type Country struct {
	Id          string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

type Output struct {
	Name      string    `json:"name"`
	Countries []Country `json:"country"`
	Count     int       `json:"count"`
}

func (c *Client) Get(name string) (Output, error) {
	// i toooo laze to fix "Not found" error
	// with http.PostForm(Url, url.Values{"name": name})
	resp, err := c.httpClient.Get(
		Url + "?name=" + name,
	)
	if err != nil {
		return Output{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Output{}, err
	}

	out := Output{}
	err = json.Unmarshal(body, &out)
	if err != nil {
		return Output{}, err
	}

	return out, nil
}
