package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/BESTSELLER/terraform-provider-servicenow-data/internal/models"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Client holds the client info for netbox
type Client struct {
	ctx  context.Context
	url  string
	user string
	pass string
}

var httpClient *http.Client
var once sync.Once

// NewClient creates common settings
func NewClient(ctx context.Context, url, user, pass string) *Client {
	once.Do(func() {
		httpClient = &http.Client{
			Timeout: time.Duration(30) * time.Second,
		}
	})

	return &Client{
		ctx:  ctx,
		url:  url,
		user: user,
		pass: pass,
	}
}

func (client *Client) GetTableRow(tableID string, params map[string]interface{}) (*models.ParsedResult, error) {
	tflog.Info(client.ctx, fmt.Sprintf("GetTableRow: tableID=%s, params=%s", tableID, params))
	if params == nil || len(params) == 0 {
		return nil, fmt.Errorf("sys_id and params cannot be both empty")
	}
	query := "?"
	for k, v := range params {
		query = fmt.Sprintf("%s&%s=%s", query, k, v)
	}
	rowPath := fmt.Sprintf("/api/now/table/%s%s", tableID, query)

	rawData, err := client.SendRequest(http.MethodGet, rowPath, nil, 200)
	if err != nil {
		return nil, err
	}
	result, err := parseRawListData(rawData)
	tflog.Info(client.ctx, fmt.Sprintf("GetTableRow: err=%d, result=%s", err, result))
	return result, err

}

func (client *Client) InsertTableRow(tableID string, tableData interface{}) (*models.ParsedResult, error) {
	rowPath := fmt.Sprintf("/api/now/table/%s", tableID)
	rawData, err := client.SendRequest(http.MethodPost, rowPath, tableData, 201)
	if err != nil {
		return nil, err
	}
	return parseRawData(rawData)
}

func (client *Client) DeleteTableRow(tableID string, sysID string) error {
	rowPath := fmt.Sprintf("/api/now/table/%s/%s", tableID, sysID)
	_, err := client.SendRequest(http.MethodDelete, rowPath, nil, 204, 404)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) SendRequest(method, path string, payload interface{}, expectedStatusCodes ...int) (value *[]byte, err error) {
	url := client.url + path
	tflog.Info(client.ctx, fmt.Sprintf("SendRequest: url=%s, method=%s", url, method))

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
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if len(expectedStatusCodes) > 0 {
		for _, code := range expectedStatusCodes {
			if resp.StatusCode == code {
				tflog.Info(client.ctx, fmt.Sprintf("SendRequest response: %d %s", resp.StatusCode, string(body)))
				return &body, nil
			}
		}
		return nil, fmt.Errorf("[ERROR] unexpected status code got: %v expected: %v  \n %v  \n %v", resp.StatusCode, expectedStatusCodes, string(body), url)
	}
	tflog.Info(client.ctx, fmt.Sprintf("SendRequest response: %d %s", resp.StatusCode, string(body)))
	return &body, nil
}

func (client *Client) UpdateTableRow(tableID, sysID string, payload interface{}) (*models.ParsedResult, error) {
	rowPath := fmt.Sprintf("/api/now/table/%s/%s", tableID, sysID)
	rawData, err := client.SendRequest(http.MethodPut, rowPath, payload, 200)
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

	switch len(rawResult.Result) {
	case 0:
		return &models.ParsedResult{SysData: map[string]string{}, RowData: map[string]string{}}, nil
	case 1:
		break
	default:
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
