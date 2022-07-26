package utils

import "encoding/json"

type ErrorResponse struct {
	Error string `json:"error"`
}

func PrintError(err error) []byte {
	response := ErrorResponse{
		Error: err.Error(),
	}
	errJson, _ := json.Marshal(response)
	return errJson
}
