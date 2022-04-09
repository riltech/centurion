package engine

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/riltech/centurion/core/bus"
	"github.com/riltech/centurion/core/engine/dto"
)

// Describes the interface of the engine controller
type IConroller interface {
	Register(http.ResponseWriter, *http.Request, httprouter.Params)
	GetRouter() *httprouter.Router
}

// Controller implementation
type Controller struct {
	bus bus.IBus
}

// Constructor for the engine controller
func NewController(bus bus.IBus) IConroller {
	return &Controller{
		bus,
	}
}

// Handles /team/register request
func (c Controller) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var responseCreator *ResponseCreator
	body, err := r.GetBody()
	if err != nil {
		responseCreator.BadRequest(w, map[string]interface{}{
			"reason": "Body is malformed",
		})
		return
	}
	var b []byte
	b, err = ioutil.ReadAll(body)
	if err != nil {
		responseCreator.BadRequest(w, map[string]interface{}{
			"reason": "Body is malformed",
		})
		return
	}
	var dto dto.RegisterRequest
	if err = json.Unmarshal(b, &dto); err != nil {
		responseCreator.BadRequest(w, map[string]interface{}{
			"reason": "Body is malformed",
		})
		return
	}
	if dto.Name == "" {
		responseCreator.BadRequest(w, map[string]interface{}{
			"reason": "Name is required",
		})
		return
	}
	switch strings.ToLower(dto.Team) {
	case "defender":
	case "attacker":
	default:
		responseCreator.BadRequest(w, map[string]interface{}{
			"reason": "Team has to be either attacker or defender",
		})
		return
	}
	c.bus.Send(&bus.BusEvent{
		Type: bus.RegisterEventType,
		Information: bus.RegistrationEvent{
			Name: dto.Name,
			Team: dto.Team,
			ID:   string(uuid.NewString()),
		},
	})
}

// Creates a new router
func (c Controller) GetRouter() *httprouter.Router {
	router := httprouter.New()
	router.POST("/team/register", c.Register)
	return router
}
