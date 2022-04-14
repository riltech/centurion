package example

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/websocket"
	"github.com/riltech/centurion/core/engine/dto"
	"github.com/riltech/centurion/core/logger"
)

// Describes the interface of an example client
type IDefender interface {
	// Starts the client
	Start()
	// Stops the client
	Stop()
}

// Describes an example client
type Defender struct {
	// Host of the server
	host string
	// Websocket address of the server
	address url.URL
	// Stops the client
	stop chan uint8
}

// Interface check
var _ IDefender = (*Defender)(nil)

// Constructor for the client
func NewDefender(host string) IDefender {
	return &Defender{
		host:    host,
		address: url.URL{Scheme: "ws", Host: host, Path: "/team/join"},
		stop:    make(chan uint8, 1),
	}
}

func (d Defender) Stop() {
	d.stop <- 1
}

func (d Defender) Start() {
	client := http.Client{}
	registerDTO := dto.RegisterRequest{
		Name: gofakeit.Name(),
		Team: "defender",
	}
	b, err := json.Marshal(registerDTO)
	if err != nil {
		logger.LogError(err)
		return
	}
	resp, err := client.Post(
		fmt.Sprintf("http://%s/team/register", d.host),
		"application/json",
		bytes.NewBuffer(b))
	if err != nil {
		logger.LogError(err)
		return
	}
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.LogError(err)
		return
	}
	var regResponseDTO dto.RegisterResponse
	if err = json.Unmarshal(b, &regResponseDTO); err != nil {
		logger.LogError(err)
		return
	}
	if regResponseDTO.Code != 200 {
		logger.LogError(fmt.Errorf("Client response was not 200 %s", spew.Sdump(regResponseDTO)))
		return
	}
	<-time.After(time.Second * 2)
	conn, _, err := websocket.DefaultDialer.Dial(d.address.String(), nil)
	if err != nil {
		logger.LogError(err)
		return
	}
	defer conn.Close()
	join := dto.JoinEvent{
		SocketEvent: dto.SocketEvent{
			Type: dto.SocketEventTypeJoin,
		},
		ID: regResponseDTO.ID,
	}
	if err = conn.WriteJSON(join); err != nil {
		logger.LogError(err)
		return
	}
	<-time.After(time.Second * 5)

	customChallengeDTO := dto.InstallChallengeRequest{
		DefenderID:  regResponseDTO.ID,
		Name:        "Reverse sorter - 2",
		Description: "You do the same as in reverse sorter 1, this is a demo challenge",
		Example: dto.ChallengeExampleDTO{
			Hints:     []interface{}{"123456"},
			Solutions: []interface{}{"654321"},
		},
	}
	b, err = json.Marshal(customChallengeDTO)
	if err != nil {
		logger.LogError(err)
		return
	}
	resp, err = client.Post(fmt.Sprintf("http://%s/challenges", d.host), "application/json", bytes.NewBuffer(b))
	if err != nil {
		logger.LogError(err)
		return
	}
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.LogError(err)
		return
	}
	var challengeRespDTO dto.InstallChallengeResponse
	if err = json.Unmarshal(b, &challengeRespDTO); err != nil {
		logger.LogError(err)
		return
	}
	for {
		var event dto.SocketEvent
		_, message, err := conn.ReadMessage()
		if err != nil {
			logger.LogError(err)
			return
		}
		if err = json.Unmarshal(message, &event); err != nil {
			logger.LogError(err)
			continue
		}
		if event.Type == dto.SocketEventTypeDefendActionRequest {
			var detailedEvent dto.DefendActionRequestEvent
			if err = json.Unmarshal(message, &detailedEvent); err != nil {
				logger.LogError(err)
				continue
			}
			if detailedEvent.TargetID != challengeRespDTO.ID {
				logger.LogError(fmt.Errorf("Received target ID invalid (created, received) (%s, %s)", challengeRespDTO.ID, detailedEvent.TargetID))
				continue
			}
			if err = conn.WriteJSON(dto.DefendActionEvent{
				SocketEvent: dto.SocketEvent{
					Type: dto.SocketEventTypeDefendAction,
				},
				Hints:    []interface{}{"12345678910"}, // Ideally this should be generated
				CombatID: detailedEvent.CombatID,
			}); err != nil {
				logger.LogError(err)
			}
			continue
		}
		if event.Type == dto.SocketEventTypeSolutionEvaluationRequest {
			var detailedEvent dto.SolutionEvaluationRequestEvent
			if err = json.Unmarshal(message, &detailedEvent); err != nil {
				logger.LogError(err)
				continue
			}
			if detailedEvent.TargetID != challengeRespDTO.ID {
				logger.LogError(fmt.Errorf("Received target ID invalid (created, received) (%s, %s)", challengeRespDTO.ID, detailedEvent.TargetID))
				continue
			}
			if len(detailedEvent.Hints) != 1 || len(detailedEvent.Solutions) != 1 {
				logger.LogError(fmt.Errorf("Received hints and solutions are not what expected"))
				if err = conn.WriteJSON(dto.SolutionEvaluationEvent{
					SocketEvent: dto.SocketEvent{
						Type: dto.SocketEventTypeSolutionEvaluation,
					},
					TargetID: detailedEvent.TargetID,
					CombatID: detailedEvent.CombatID,
					Success:  false,
					Message:  "Length of hints or solutions is not 1",
				}); err != nil {
					logger.LogError(err)
				}
				continue
			}
			hint, ok := detailedEvent.Hints[0].(string)
			if !ok {
				if err = conn.WriteJSON(dto.SolutionEvaluationEvent{
					SocketEvent: dto.SocketEvent{
						Type: dto.SocketEventTypeSolutionEvaluation,
					},
					TargetID: detailedEvent.TargetID,
					CombatID: detailedEvent.CombatID,
					Success:  false,
					Message:  "Hints has to contain exactly 1 string",
				}); err != nil {
					logger.LogError(err)
				}
				continue
			}
			solution, ok := detailedEvent.Solutions[0].(string)
			if !ok {
				if err = conn.WriteJSON(dto.SolutionEvaluationEvent{
					SocketEvent: dto.SocketEvent{
						Type: dto.SocketEventTypeSolutionEvaluation,
					},
					TargetID: detailedEvent.TargetID,
					CombatID: detailedEvent.CombatID,
					Success:  false,
					Message:  "Solutions has to contain exactly 1 string",
				}); err != nil {
					logger.LogError(err)
				}
				continue
			}
			if hint != "12345678910" || solution != "01987654321" {
				if err = conn.WriteJSON(dto.SolutionEvaluationEvent{
					SocketEvent: dto.SocketEvent{
						Type: dto.SocketEventTypeSolutionEvaluation,
					},
					TargetID: detailedEvent.TargetID,
					CombatID: detailedEvent.CombatID,
					Success:  false,
					Message:  "Invalid solution or hint",
				}); err != nil {
					logger.LogError(err)
				}
			} else {
				if err = conn.WriteJSON(dto.SolutionEvaluationEvent{
					SocketEvent: dto.SocketEvent{
						Type: dto.SocketEventTypeSolutionEvaluation,
					},
					TargetID: detailedEvent.TargetID,
					CombatID: detailedEvent.CombatID,
					Success:  true,
					Message:  "",
				}); err != nil {
					logger.LogError(err)
				}
			}
			continue
		}
	}
}
