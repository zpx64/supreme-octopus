package genderize

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

const (
	Url = "https://api.genderize.io/"
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

type Gender int

const (
	Unknown Gender = iota
	Male
	Female
)

func (a Gender) String() string {
	switch a {
	default:
		return "unknown"
	case Male:
		return "male"
	case Female:
		return "female"
	}
}

func (a *Gender) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	default:
		*a = Unknown
	case "male":
		*a = Male
	case "female":
		*a = Female
	}

	return nil
}

func (a Gender) MarshalJSON() ([]byte, error) {
	var s string
	switch a {
	default:
		s = "unknown"
	case Male:
		s = "male"
	case Female:
		s = "female"
	}

	return json.Marshal(s)
}

type Output struct {
	Name        string  `json:"name"`
	Gender      Gender  `json:"gender"`
	Probability float64 `json:"probability"`
	Count       int     `json:"count"`
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
