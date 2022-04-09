package engine

import (
	"encoding/json"
	"net/http"

	"github.com/riltech/centurion/core/engine/dto"
)

// Static pointer struct to provide functions for consistent responses
type ResponseCreator struct{}

// Adds header for JSON response
func (rc *ResponseCreator) jsonResponse(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
}

// Generic bad request
func (rc *ResponseCreator) BadRequest(w http.ResponseWriter, meta map[string]interface{}) {
	rc.jsonResponse(w)
	b, err := json.Marshal(dto.CenturionResponse{
		Message: "Bad request",
		Code:    400,
		Meta:    meta,
	})
	if err != nil {
		panic(err)
	}
	w.Write(b)
}
