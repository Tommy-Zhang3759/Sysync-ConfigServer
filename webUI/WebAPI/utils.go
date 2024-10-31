package WebAPI

import (
	"encoding/json"
)

type FunctionRequest struct {
	FName string `json:"f_name"`
}

func getFName(body *[]byte) (string, error) {
	var requestData FunctionRequest
	if err := json.Unmarshal(*body, &requestData); err != nil {
		return "", err
	}
	return requestData.FName, nil
}
