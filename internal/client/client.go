package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BESTSELLER/terraform-provider-servicenow-data/internal/models"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

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

func (client *Client) GetTableRow(tableID string, params map[string]interface{}) (*models.ParsedResult, error) {
	if params == nil || len(params) == 0 {
		return nil, fmt.Errorf("sys_id and params cannot be both empty")
	}
	query := "?"
	for k, v := range params {
		query = fmt.Sprintf("%s&%s=%s", query, k, v)
	}
	rowPath := fmt.Sprintf("/table/%s%s", tableID, query)

	rawData, err := client.sendRequest(http.MethodGet, rowPath, nil, 200)
	if err != nil {
		return nil, err
	}
	return parseRawListData(rawData)
}

func (client *Client) InsertTableRow(tableID string, tableData interface{}) (*models.ParsedResult, error) {
	rowPath := fmt.Sprintf("/table/%s", tableID)
	rawData, err := client.sendRequest(http.MethodPost, rowPath, tableData, 201)
	if err != nil {
		return nil, err
	}
	return parseRawData(rawData)
}

func (client *Client) DeleteTableRow(tableID string, sysID string) error {
	rowPath := fmt.Sprintf("/table/%s/%s", tableID, sysID)
	_, err := client.sendRequest(http.MethodDelete, rowPath, nil, 204)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) sendRequest(method, path string, payload interface{}, expectedStatusCode int) (value *[]byte, err error) {
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

	if expectedStatusCode != 0 {
		if resp.StatusCode != expectedStatusCode {
			return nil, fmt.Errorf("[ERROR] unexpected status code got: %v expected: %v  \n %v  \n %v", resp.StatusCode, expectedStatusCode, string(body), url)
		}
	}

	return &body, nil
}

func (client *Client) UpdateTableRow(tableID, sysID string, payload interface{}) (*models.ParsedResult, error) {
	rowPath := fmt.Sprintf("/table/%s/%s", tableID, sysID)
	rawData, err := client.sendRequest(http.MethodPut, rowPath, payload, 200)
	if err != nil {
		return nil, err
	}
	return parseRawData(rawData)
}

func parseRawData(rawData *[]byte) (*models.ParsedResult, error) {
	var rawResult models.RawResult
	err := json.Unmarshal(*rawData, &rawResult)
	if err != nil {
		return nil, err
	}

	result, err := rawMapParse(rawResult.Result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func parseRawListData(rawData *[]byte) (*models.ParsedResult, error) {
	var rawResult models.RawResultList
	err := json.Unmarshal(*rawData, &rawResult)
	if err != nil {
		return nil, err
	}
	if len(rawResult.Result) != 1 {
		return nil, fmt.Errorf("received more than one row as result, make sure your query returns a single item")
	}

	result, err := rawMapParse(rawResult.Result[0])
	if err != nil {
		return nil, err
	}
	return result, nil
}

func rawMapParse(rawResult map[string]json.RawMessage) (*models.ParsedResult, error) {
	var parsedResult = models.ParsedResult{}
	rowData := make(map[string]string, len(rawResult)-7)
	sysData := make(map[string]string, 7)
	for k, message := range rawResult {
		rv, err := extractRowValue(message)
		if err != nil {
			return nil, err
		}
		//A small hack :), I'm sure nothing will go wrong here
		if strings.HasPrefix(k, "sys_") {
			sysData[k] = rv
		} else {
			rowData[k] = rv
		}
	}
	parsedResult.SysData = sysData
	parsedResult.RowData = rowData
	return &parsedResult, nil
}

func extractRowValue(rm json.RawMessage) (string, error) {
	var str string
	err := json.Unmarshal(rm, &str)
	if err == nil {
		return str, nil
	} else {
		var ai models.ReferenceItem
		err = json.Unmarshal(rm, &ai)
		if err != nil {
			return "", fmt.Errorf("Unmarshal exploded for result.%v", rm)
		}
		return ai.Value, nil
	}
}
