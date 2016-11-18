package common

import "encoding/json"

type Result struct {
	StatusCode int         `json:"status_code"`
	Msg        string      `json:"msg"`
	Data       interface{} `json:"data"`
}

func (result *Result) String() string {
	b, _ := json.Marshal(result)
	return string(b)
}
