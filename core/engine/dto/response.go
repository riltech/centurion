package dto

type CenturionResponse struct {
	Message string                 `json:"message"`
	Code    int                    `json:"code"`
	Meta    map[string]interface{} `json:"meta"`
}
