package models

import "encoding/json"

type ApprovalItem struct {
	Link  string `json:"link"`
	Value string `json:"value"`
}

type RawResult struct {
	Result map[string]json.RawMessage `json:"result"`
}
