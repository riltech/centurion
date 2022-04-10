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
	"github.com/riltech/centurion/core/player"
	"github.com/sirupsen/logrus"
)

// Describes the interface of the engine controller
type IConroller interface {
	Register(http.ResponseWriter, *http.Request, httprouter.Params)
	GetRouter() *httprouter.Router
}

// Controller implementation
type Controller struct {
	bus     bus.IBus
	service IService
}

// Constructor for the engine controller
func NewController(bus bus.IBus, service IService) IConroller {
	return &Controller{
		bus,
		service,
	}
}

// Cleans up and logs any potential panic that occours
// in the controller functions
func (c Controller) cleanUp(w http.ResponseWriter) {
	err := recover()
	if err == nil {
		return
	}
	logrus.Error(err)
	(*ResponseCreator)(nil).InternalServerError(w)
}

// Handles /team/register request
func (c Controller) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var response *ResponseCreator
	defer c.cleanUp(w)
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.BadRequest(w, map[string]interface{}{
			"reason": "Body is malformed",
		})
		return
	}
	var reqDTO dto.RegisterRequest
	if err = json.Unmarshal(b, &reqDTO); err != nil {
		response.BadRequest(w, map[string]interface{}{
			"reason": "Body is malformed",
		})
		return
	}
	if reqDTO.Name == "" {
		response.BadRequest(w, map[string]interface{}{
			"reason": "Name is required",
		})
		return
	}
	switch strings.ToLower(reqDTO.Team) {
	case "defender":
	case "attacker":
	default:
		response.BadRequest(w, map[string]interface{}{
			"reason": "Team has to be either attacker or defender",
		})
		return
	}
	information := bus.RegistrationEvent{
		Name: reqDTO.Name,
		Team: reqDTO.Team,
		ID:   string(uuid.NewString()),
	}
	if exists := c.service.IsPlayerExist(&player.Model{
		Name: reqDTO.Name,
		Team: reqDTO.Team,
		ID:   information.ID,
	}); exists {
		response.BadRequest(w, map[string]interface{}{
			"reason": "Player already exists",
		})
		return
	}
	c.bus.Send(&bus.BusEvent{
		Type:        bus.EventTypeRegistration,
		Information: information,
	})
	response.OK(w, dto.RegisterResponse{
		CenturionResponse: dto.CenturionResponse{
			Message: "Success",
			Code:    200,
			Meta:    nil,
		},
		ID: information.ID,
	})
}

// Creates a new router
func (c Controller) GetRouter() *httprouter.Router {
	router := httprouter.New()
	router.POST("/team/register", c.Register)
	return router
}
