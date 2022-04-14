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
	"github.com/riltech/centurion/core/scoreboard"
	"github.com/sirupsen/logrus"
)

// Describes the interface of the engine controller
type IConroller interface {
	// Status endpoint
	Ping(http.ResponseWriter, *http.Request, httprouter.Params)
	// Endpoint for fetching currently available challenges
	FetchChallanges(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	// Endpoint for installing defense modules
	InstallChallenge(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	// Endpoint for registration
	Register(http.ResponseWriter, *http.Request, httprouter.Params)
	// Entry point for the websocket API
	PlayerJoin(w http.ResponseWriter, r *http.Request)
	// Boostrapping of the router
	GetRouter() *httprouter.Router
}

// Controller implementation
type Controller struct {
	bus              bus.IBus
	engineService    IService
	playerService    player.IService
	challengeService challenge.IService
	scoreService     scoreboard.IService

	// Websocket
	upgrader websocket.Upgrader
}

// Constructor for the engine controller
func NewController(
	bus bus.IBus,
	engineService IService,
	playerService player.IService,
	challengeService challenge.IService,
	scoreService scoreboard.IService,
) IConroller {
	return &Controller{
		bus,
		engineService,
		playerService,
		challengeService,
		scoreService,
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

func (c Controller) PlayerJoin(w http.ResponseWriter, r *http.Request) {
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

func (c Controller) InstallChallenge(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var response *ResponseCreator
	defer c.cleanUp(w)
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.BadRequest(w, map[string]interface{}{
			"reason": "Body is malformed",
		})
		return
	}
	var reqDTO dto.InstallChallengeRequest
	if err = json.Unmarshal(b, &reqDTO); err != nil {
		response.BadRequest(w, map[string]interface{}{
			"reason": "Body is malformed",
		})
		return
	}
	defender, err := c.playerService.FindByID(reqDTO.DefenderID)
	if err != nil || defender.Team != player.TeamTypeDefender {
		response.BadRequest(w, map[string]interface{}{
			"reason": "Defender ID is invalid or not found",
		})
		return
	}
	toCreate := challenge.Model{
		Type:        challenge.ChallengeTypePlayerCreated,
		ID:          uuid.NewString(),
		CreatorID:   reqDTO.DefenderID,
		Name:        reqDTO.Name,
		Description: reqDTO.Description,
		Example: challenge.Example{
			Hints:     reqDTO.Example.Hints,
			Solutions: reqDTO.Example.Solutions,
		},
	}
	isFirstModule := c.challengeService.IsFirstModule(toCreate)
	err = c.challengeService.AddChallenge(toCreate)
	if err != nil {
		response.BadRequest(w, map[string]interface{}{
			"reason": err.Error(),
		})
		return
	}
	if isFirstModule {
		if err = c.scoreService.AddPoint(reqDTO.DefenderID, 1); err != nil {
			logrus.Error(err)
		}
	}
	c.bus.Send(&bus.BusEvent{
		Type: bus.EventTypeDefenseModuleInstalled,
		Information: bus.DefenseModuleInstalledEvent{
			Name:        reqDTO.Name,
			CreatorName: defender.Name,
		},
	})
	response.OK(w, dto.InstallChallengeResponse{
		CenturionResponse: dto.CenturionResponse{
			Message: "Success",
			Code:    200,
			Meta:    nil,
		},
		ID: toCreate.ID,
	})
}

func (c Controller) Ping(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	(*ResponseCreator)(nil).Empty200(w)
}

// Creates a new router
func (c Controller) GetRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/", c.Ping)
	router.POST("/team/register", c.Register)
	router.GET("/challenges", c.FetchChallanges)
	router.POST("/challenges", c.InstallChallenge)
	router.HandlerFunc("GET", "/team/join", c.PlayerJoin)
	return router
}
