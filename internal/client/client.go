package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BESTSELLER/terraform-provider-servicenow-data/internal/models"
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

func (client *Client) GetTableRow(tableName, sysID string) (*map[string]string, error) {
	rowPath := fmt.Sprintf("/table/%s/%s", tableName, sysID)
	rawData, err := client.sendRequest(http.MethodGet, rowPath, nil, 200)
	if err != nil {
		return nil, err
	}
	var objMap map[string]json.RawMessage
	err = json.Unmarshal(*rawData, &objMap)
	if err == nil {
		return nil, err
	}
	rowData := make(map[string]string, len(objMap))

	for s, message := range objMap {
		var str string
		err = json.Unmarshal(message, &str)
		if err == nil {
			rowData[s] = str
		} else {
			var ai models.ApprovalItem
			err = json.Unmarshal(message, &ai)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Unmarshal exploded for result.%s.%v", s, message))
			}
			rowData[s] = ai.Value
		}
	}
	return &rowData, nil
}

func (client *Client) sendRequest(method, path string, payload interface{}, statusCode int) (value *[]byte, err error) {
	url := client.url + "/api/now" + path

	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, b)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(client.user, client.pass)
	req.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if statusCode != 0 {
		if resp.StatusCode != statusCode {
			return nil, fmt.Errorf("[ERROR] unexpected status code got: %v expected: %v  \n %v  \n %v", resp.StatusCode, statusCode, string(body), url)
		}
	}

	return &body, nil
}
