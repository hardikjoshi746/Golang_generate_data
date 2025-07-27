package models

import "encoding/json"

type ProcessResponse struct {
	User_id    string `json:"user_id"`
	Result     string `json:"result"`
	Duration   int    `json:"duration"`
	Words_used int    `json:"words_used"`
	Words_left int    `json:"words_left"`
}

func (j *ProcessResponse) ToJSON() (string, error) {
	data, err := json.Marshal(j)
	return string(data), err
}

func JobFromJSON(data string) (*ProcessResponse, error) {
	var processResponse ProcessResponse
	err := json.Unmarshal([]byte(data), &processResponse)
	return &processResponse, err
}
