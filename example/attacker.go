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
	"github.com/sirupsen/logrus"
)

// Describes the interface of an example client
type IAttacker interface {
	// Starts the client
	Start()
	// Stops the client
	Stop()
}

// Describes an example client implementation
type Attacker struct {
	// Address of the server
	address url.URL
	// Stops the client
	stop chan uint8
}

// Interface check
var _ IAttacker = (*Attacker)(nil)

// Constructor for the client
func NewAttacker() IAttacker {
	return &Attacker{
		address: url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/team/join"},
		stop:    make(chan uint8, 1),
	}
}

func (a Attacker) Stop() {
	a.stop <- 1
}

// Called when the default module is solved
func (a Attacker) attackDefendBotChallenge(ID string, conn *websocket.Conn) {
	logrus.Info("Attacker bot is waiting 15 seconds to solve defender bot challenge")
	<-time.After(15 * time.Second)
	logrus.Info("Attacker bot starts defender bot challenge")
	client := http.Client{}
	resp, err := client.Get("http://localhost:8080/challenges")
	if err != nil {
		logger.LogError(err)
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.LogError(err)
		return
	}
	var challenges dto.FetchChallengesResponse
	if err = json.Unmarshal(b, &challenges); err != nil {
		logger.LogError(err)
		return
	}
	for _, c := range challenges.Challenges {
		logrus.Info(c.Name)
		if c.Name == "Reverse sorter - 2" {
			logrus.Info("Attacker bot write")
			err = conn.WriteJSON(dto.AttackEvent{
				SocketEvent: dto.SocketEvent{
					Type: dto.SocketEventTypeAttack,
				},
				TargetID: c.ID,
			})
			logrus.Info("Attacker bot written")
			if err != nil {
				logger.LogError(err)
			}
			return
		}
	}
	logger.LogError(fmt.Errorf("Attacker did not find defender bot challenge"))
}

func (a Attacker) Start() {
	client := http.Client{}
	registerDTO := dto.RegisterRequest{
		Name: gofakeit.Name(),
		Team: "attacker",
	}
	b, err := json.Marshal(registerDTO)
	if err != nil {
		logger.LogError(err)
		return
	}
	resp, err := client.Post(
		"http://localhost:8080/team/register",
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
	conn, _, err := websocket.DefaultDialer.Dial(a.address.String(), nil)
	if err != nil {
		logger.LogError(err)
		return
	}
	defer conn.Close()
	resp, err = client.Get("http://localhost:8080/challenges")
	if err != nil {
		logger.LogError(err)
		return
	}
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.LogError(err)
		return
	}
	var challenges dto.FetchChallengesResponse
	if err = json.Unmarshal(b, &challenges); err != nil {
		logger.LogError(err)
		return
	}
	var selected *dto.ChallengeResponseDTO
	for _, c := range challenges.Challenges {
		if c.Name == "Reverse sorter" {
			selected = c
		}
	}
	if selected == nil {
		logger.LogError(fmt.Errorf("Could not find default module"))
		return
	}
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
	attack := dto.AttackEvent{
		SocketEvent: dto.SocketEvent{
			Type: dto.SocketEventTypeAttack,
		},
		TargetID: selected.ID,
	}
	if err = conn.WriteJSON(attack); err != nil {
		logger.LogError(err)
		return
	}
	firstStagePassed := false
	for {
		var event dto.SocketEvent
		_, message, err := conn.ReadMessage()
		if err != nil {
			logger.LogError(err)
			return
		}
		if err = json.Unmarshal(message, &event); err != nil {
			logger.LogError(err)
			return
		}
		if event.Type == dto.SocketEventTypeAttackChallenge {
			var detailedEvent dto.AttackChallengeEvent
			err = json.Unmarshal(message, &detailedEvent)
			if err != nil {
				logger.LogError(err)
				return
			}
			if len(detailedEvent.Hints) == 0 {
				logger.LogError(fmt.Errorf("Did not get hints in example client"))
				return
			}
			switch detailedEvent.Hints[0].(type) {
			case string:
				hint := detailedEvent.Hints[0].(string)
				solution := ""
				for i := range hint {
					solution += string(hint[len(hint)-(1+i)])
				}
				conn.WriteJSON(dto.AttackSolutionEvent{
					SocketEvent: dto.SocketEvent{
						Type: dto.SocketEventTypeAttackSolution,
					},
					Hints:     detailedEvent.Hints,
					Solutions: []interface{}{solution},
					TargetID:  selected.ID,
				})
				continue
			default:
				logger.LogError(fmt.Errorf("Hint was not a string in example client"))
				return
			}
		}
		if event.Type == dto.SocketEventTypeAttackResult {
			var detailedEvent dto.AttackResultEvent
			err = json.Unmarshal(message, &detailedEvent)
			if err != nil {
				logger.LogError(err)
				return
			}
			if detailedEvent.Success == false {
				logrus.Infof("Challenge (%s) failed by bot", selected.ID)
			} else {
				logrus.Infof("Challenge (%s) resolved by bot", selected.ID)
				if !firstStagePassed {
					firstStagePassed = true
					go a.attackDefendBotChallenge(regResponseDTO.ID, conn)
				} else {
					return
				}
			}
			continue
		}
	}
}
