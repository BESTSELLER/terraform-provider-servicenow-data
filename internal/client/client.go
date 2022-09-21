package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func init() {

}

// Client holds the client info for netbox
type Client struct {
	url  string
	user string
	pass string
}

var httpClient *http.Client
var once sync.Once

// NewClient creates common settings
func NewClient(url, user, pass string) *Client {
	once.Do(func() {
		httpClient = &http.Client{
			Timeout: time.Duration(30) * time.Second,
		}
	})

	return &Client{
		url:  url,
		user: user,
		pass: pass,
	}
}

func (Client *Client) sendRequest(method string, path string, payload interface{}, statusCode int) (value string, err error) {
	url := Client.url + path

	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(method, url, b)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(Client.user, Client.pass)
	req.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	strbody := string(body)

	if statusCode != 0 {
		if resp.StatusCode != statusCode {
			return "", fmt.Errorf("[ERROR] unexpected status code got: %v expected: %v  \n %v  \n %v", resp.StatusCode, statusCode, strbody, url)
		}
	}

	return strbody, nil
}
