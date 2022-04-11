package engine

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/riltech/centurion/core/bus"
	"github.com/riltech/centurion/core/challenge"
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
	bus              bus.IBus
	engineService    IService
	playerService    player.IService
	challengeService challenge.IService

	// Websocket
	upgrader websocket.Upgrader
}

// Constructor for the engine controller
func NewController(
	bus bus.IBus,
	engineService IService,
	playerService player.IService,
	challengeService challenge.IService,
) IConroller {
	return &Controller{
		bus,
		engineService,
		playerService,
		challengeService,
		websocket.Upgrader{},
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

// Returns the available challenges
func (c Controller) FetchChallanges(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var response *ResponseCreator
	defer c.cleanUp(w)
	challenges := c.challengeService.GetChallenges()
	dtoChallenges := []*dto.ChallengeResponseDTO{}
	for _, challenge := range challenges {
		dtoChallenges = append(dtoChallenges, &dto.ChallengeResponseDTO{
			ID:          challenge.ID,
			Name:        challenge.Name,
			Description: challenge.Description,
			Example: dto.ChallengeExampleDTO{
				Hints:     challenge.Example.Hints,
				Solutions: challenge.Example.Solutions,
			},
		})
	}
	response.OK(w, dto.FetchChallengesResponse{
		CenturionResponse: dto.CenturionResponse{
			Message: "Success",
			Code:    200,
			Meta:    nil,
		},
		Challenges: dtoChallenges,
	})
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
	case player.TeamTypeDefender:
	case player.TeamTypeAttacker:
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
	newPlayer := player.Model{
		Name: reqDTO.Name,
		Team: reqDTO.Team,
		ID:   information.ID,
	}
	if exists := c.playerService.IsPlayerExist(&newPlayer); exists {
		response.BadRequest(w, map[string]interface{}{
			"reason": "Player already exists",
		})
		return
	}
	if err = c.playerService.AddPlayer(newPlayer); err != nil {
		panic(err)
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

// Handles players who join the live game
func (c Controller) PlayerJoin(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Got request")
	connection, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Error("upgrade:", err)
		return
	}
	defer connection.Close()
	_, message, err := connection.ReadMessage()
	if err != nil {
		logrus.Error("SERVER | read:", err)
		return
	}
	var join dto.JoinEvent
	if err = json.Unmarshal(message, &join); err != nil {
		logrus.Error("SERVER | Invalid join request")
		return
	}
	c.engineService.Join(join, connection)
}

// Creates a new router
func (c Controller) GetRouter() *httprouter.Router {
	router := httprouter.New()
	router.POST("/team/register", c.Register)
	router.GET("/challenges", c.FetchChallanges)
	router.HandlerFunc("GET", "/team/join", c.PlayerJoin)
	return router
}
