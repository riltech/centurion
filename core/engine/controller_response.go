package engine

import (
	"encoding/json"
	"net/http"

	"github.com/riltech/centurion/core/engine/dto"
	"github.com/sirupsen/logrus"
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
		logrus.Error(err)
		return
	}
	w.Write(b)
}

// Generic 500
func (rc *ResponseCreator) InternalServerError(w http.ResponseWriter) {
	rc.jsonResponse(w)
	b, err := json.Marshal(dto.CenturionResponse{
		Message: "Internal server error",
		Code:    500,
		Meta:    nil,
	})
	if err != nil {
		logrus.Error(err)
		return
	}
	w.Write(b)
}

// Describes an empty 200 response
func (rc *ResponseCreator) Empty200(w http.ResponseWriter) {
	rc.jsonResponse(w)
	b, err := json.Marshal(dto.CenturionResponse{
		Message: "",
		Code:    200,
		Meta:    nil,
	})
	if err != nil {
		logrus.Error(err)
		return
	}
	w.Write(b)
}

// Describes a 200 with a given body
func (rc *ResponseCreator) OK(w http.ResponseWriter, body interface{}) {
	rc.jsonResponse(w)
	b, err := json.Marshal(body)
	if err != nil {
		logrus.Error(err)
		return
	}
	w.Write(b)
}
