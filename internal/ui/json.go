package ui

import (
	"encoding/json"
	"fmt"
)

type JSONOutput struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
}

func PrintJSON(data interface{}, err error) {
	output := JSONOutput{Success: err == nil}

	if err != nil {
		output.Error = err.Error()
	} else {
		output.Data = data
	}

	jsonBytes, _ := json.MarshalIndent(output, "", "  ")
	fmt.Println(string(jsonBytes))
}

func PrintJSONWithMeta(data interface{}, meta map[string]interface{}, err error) {
	output := JSONOutput{
		Success: err == nil,
		Meta:    meta,
	}

	if err != nil {
		output.Error = err.Error()
	} else {
		output.Data = data
	}

	jsonBytes, _ := json.MarshalIndent(output, "", "  ")
	fmt.Println(string(jsonBytes))
}