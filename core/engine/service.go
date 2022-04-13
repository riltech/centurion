package engine

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/riltech/centurion/core/bus"
	"github.com/riltech/centurion/core/challenge"
	"github.com/riltech/centurion/core/combat"
	"github.com/riltech/centurion/core/engine/dto"
	"github.com/riltech/centurion/core/player"
	"github.com/sirupsen/logrus"
)

// Describes an engine service interface
type IService interface {
	// Handles join event for users
	Join(dto.JoinEvent, *websocket.Conn) error
}

// Service implementation
type Service struct {
	bus              bus.IBus
	playerService    player.IService
	challengeService challenge.IService
	combatService    combat.IService

	activeConnections map[string]*websocket.Conn
	mux               sync.RWMutex
}

// Interface check
var _ IService = (*Service)(nil)

func (s *Service) Join(event dto.JoinEvent, conn *websocket.Conn) error {
	if conn == nil {
		return fmt.Errorf("%s user socket is empty", event.ID)
	}
	updated, err := s.playerService.SetPlayerOnlineStatus(event.ID, true)
	if err != nil {
		return err
	}
	s.mux.Lock()
	s.activeConnections[event.ID] = conn
	s.mux.Unlock()
	s.bus.Send(&bus.BusEvent{
		Type: bus.EventTypePlayerJoined,
		Information: bus.PlayerJoinedEvent{
			Name: updated.Name,
			Team: updated.Team,
		},
	})
	if updated.Team == player.TeamTypeAttacker {
		return s.attacker(event.ID)
	}
	return s.defender(event.ID)
}

// Sends an error message via the socket connection
// or if the socket is closed cleans it up completely
// returns status [true] if the connection is alive
// returns [false] if the connection was terminated
func (s *Service) sendError(ID string, message string) (isConnectionStillAlive bool) {
	defer func(alive *bool) {
		if *alive {
			return
		}
		if _, err := s.playerService.SetPlayerOnlineStatus(ID, false); err != nil {
			logrus.Error(err)
		}
	}(&isConnectionStillAlive)
	if s == nil {
		isConnectionStillAlive = false
		return
	}
	s.mux.RLock()
	conn, ok := s.activeConnections[ID]
	if !ok {
		s.mux.RUnlock()
		isConnectionStillAlive = false
		return
	}
	s.mux.RUnlock()

	if err := conn.WriteJSON(dto.ErrorEvent{
		SocketEvent: dto.SocketEvent{
			Type: dto.SocketEventTypeError,
		},
		Message: message,
	}); err != nil {
		logrus.Error(err)
		s.closeConnection(ID)
		isConnectionStillAlive = false
		return
	}
	isConnectionStillAlive = true
	return
}

// This function sends a response to the socket
// or if it is not alive anymore it breaks the connection
func (s *Service) sendResponseOrBreakConnection(ID string, message interface{}) (isConnectionStillAlive bool) {
	defer func(alive *bool) {
		if *alive {
			return
		}
		if _, err := s.playerService.SetPlayerOnlineStatus(ID, false); err != nil {
			logrus.Error(err)
		}
	}(&isConnectionStillAlive)
	s.mux.RLock()
	conn, ok := s.activeConnections[ID]
	if !ok {
		s.mux.RUnlock()
		s.closeConnection(ID)
		isConnectionStillAlive = false
		return
	}
	s.mux.RUnlock()
	if err := conn.WriteJSON(message); err != nil {
		s.closeConnection(ID)
		isConnectionStillAlive = false
		return
	}
	isConnectionStillAlive = true
	return
}

// Command set for attackers
func (s *Service) attacker(ID string) error {
	s.mux.RLock()
	conn, ok := s.activeConnections[ID]
	if !ok {
		s.mux.RUnlock()
		return fmt.Errorf("%s user is not available in active players", ID)
	}
	s.mux.RUnlock()
	for {
		// Acquire message
		t, b, err := conn.ReadMessage()
		if t == websocket.CloseMessage {
			s.closeConnection(ID)
			break
		}
		if err != nil {
			logrus.Error(err)
			s.closeConnection(ID)
			break
		}
		// Deserialize message
		var event dto.SocketEvent
		if err = json.Unmarshal(b, &event); err != nil {
			if stillActive := s.sendError(ID, "Could not parse Socket Event"); !stillActive {
				break
			}
		}
		// Process of valid events
		if event.Type == dto.SocketEventTypeAttack {
			var detailedEvent dto.AttackEvent
			err = json.Unmarshal(b, &detailedEvent)
			if err != nil {
				logrus.Error(err)
				if stillActive := s.sendError(ID, "Could not parse Attack Event"); !stillActive {
					break
				}
				continue
			}
			target, err := s.challengeService.FindByID(detailedEvent.TargetID)
			if err != nil {
				logrus.Error(err)
				if stillActive := s.sendError(ID, "Invalid challenge ID"); !stillActive {
					break
				}
				continue
			}
			if target.Type == challenge.ChallengeTypeDefault {
				hints, err := s.challengeService.GenerateHintForDefault(target)
				if err != nil {
					logrus.Error(err)
					if stillActive := s.sendError(ID, err.Error()); !stillActive {
						break
					}
				}
				if isConnectionStillAlive := s.sendResponseOrBreakConnection(ID, dto.AttackChallengeEvent{
					SocketEvent: dto.SocketEvent{
						Type: dto.SocketEventTypeAttackChallenge,
					},
					TargetID: target.ID,
					Hints:    hints,
				}); !isConnectionStillAlive {
					break
				}
				continue
			}
			creator, err := s.playerService.FindByID(target.CreatorID)
			if err != nil {
				logrus.Error(err)
				if stillActive := s.sendError(ID, "Challenge owner could not be retrieved"); !stillActive {
					break
				}
				continue
			}
			newCombat := combat.Model{
				ID:          uuid.NewString(),
				ChallengeID: target.ID,
				AttackerID:  ID,
				DefenderID:  creator.ID,
				CombatState: combat.CombatStateAttackInitiated,
			}
			err = s.combatService.AddCombat(newCombat)
			if err != nil {
				logrus.Error(err)
				if stillActive := s.sendError(ID, "Combat could not be created, please try again"); !stillActive {
					break
				}
				continue
			}
			if !creator.Online {
				if _, err = s.combatService.UpdateCombatState(newCombat.ID, combat.CombatStateDefenseFailed); err != nil {
					logrus.Error(err)
				}
				// TODO: There should be a point reduction or increase
				if isConnectionStillAlive := s.sendResponseOrBreakConnection(ID, dto.DefenderFailedToDefendEvent{
					SocketEvent: dto.SocketEvent{
						Type: dto.SocketEventTypeDefenderFailedToDefend,
					},
					TargetID: target.ID,
				}); !isConnectionStillAlive {
					break
				}
			} else {
				if _, err = s.combatService.UpdateCombatState(newCombat.ID, combat.CombatStateDefenseRequested); err != nil {
					logrus.Error(err)
				}
				// if the connection is not alive here that's the problem of the potential
				// go routine handling the given defender
				s.sendResponseOrBreakConnection(creator.ID, dto.DefendActionRequestEvent{
					SocketEvent: dto.SocketEvent{
						Type: dto.SocketEventTypeDefenderFailedToDefend,
					},
					TargetID: target.ID,
					CombatID: newCombat.ID,
				})
			}
			continue
		}
		if event.Type == dto.SocketEventTypeAttackSolution {
			var detailedEvent dto.AttackSolutionEvent
			err = json.Unmarshal(b, &detailedEvent)
			if err != nil {
				if stillActive := s.sendError(ID, "Could not parse Attack Solution Event"); !stillActive {
					break
				}
				continue
			}
			target, err := s.challengeService.FindByID(detailedEvent.TargetID)
			if err != nil {
				if stillActive := s.sendError(ID, "Invalid challenge ID"); !stillActive {
					break
				}
				continue
			}
			if target.Type == challenge.ChallengeTypeDefault {
				isValid, err := s.challengeService.IsValidSolutionToDefaultModule(
					target,
					detailedEvent.Hints,
					detailedEvent.Solutions)
				if err != nil {
					if stillActive := s.sendError(ID, err.Error()); !stillActive {
						break
					}
					continue
				}
				if isConnectionStillAlive := s.sendResponseOrBreakConnection(ID, dto.AttackResultEvent{
					SocketEvent: dto.SocketEvent{
						Type: dto.SocketEventTypeAttackResult,
					},
					TargetID: target.ID,
					Success:  isValid,
				}); !isConnectionStillAlive {
					break
				}
				if isValid {
					attacker, _ := s.playerService.FindByID(ID)
					s.bus.Send(&bus.BusEvent{
						Type: bus.EventTypeAttackStateUpdate,
						Information: bus.AttackStateUpdateEvent{
							AttackerName:  attacker.Name,
							ChallengeName: target.Name,
							Success:       isValid,
						},
					})
				}
				continue
			}
			creator, err := s.playerService.FindByID(target.CreatorID)
			if err != nil {
				logrus.Error(err)
				if stillActive := s.sendError(ID, "Challenge owner could not be retrieved"); !stillActive {
					break
				}
				continue
			}
			ongoingCombat, err := s.combatService.FindByAttackerAndChallenge(ID, target.ID)
			if err != nil {
				logrus.Error(err)
				if stillActive := s.sendError(ID, "No combat found, first initiate an attack"); !stillActive {
					break
				}
				continue
			}
			if !creator.Online {
				if _, err = s.combatService.UpdateCombatState(ongoingCombat.ID, combat.CombatStateDefenseFailed); err != nil {
					logrus.Error(err)
				}
				if isConnectionStillAlive := s.sendResponseOrBreakConnection(ID, dto.DefenderFailedToDefendEvent{
					SocketEvent: dto.SocketEvent{
						Type: dto.SocketEventTypeDefenderFailedToDefend,
					},
					TargetID: target.ID,
				}); !isConnectionStillAlive {
					break
				}
			} else {
				if _, err = s.combatService.UpdateCombatState(ongoingCombat.ID, combat.CombatStateSolutionEvaluationRequested); err != nil {
					logrus.Error(err)
				}
				// if the connection is not alive here that's the problem of the potential
				// go routine handling the given defender
				s.sendResponseOrBreakConnection(creator.ID, dto.SolutionEvaluationRequestEvent{
					SocketEvent: dto.SocketEvent{
						Type: dto.SocketEventTypeSolutionEvaluationRequest,
					},
					TargetID:  target.ID,
					Solutions: detailedEvent.Solutions,
					CombatID:  ongoingCombat.ID,
				})
			}
			continue
		}
	}
	return nil
}

// Closes a connection gracefully
func (s *Service) closeConnection(ID string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	if conn, ok := s.activeConnections[ID]; ok {
		err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(200, "OK"))
		if err != nil {
			logrus.Error(err)
		}
		conn.Close()
	}
	s.activeConnections[ID] = nil
}

// Command set for defenders
func (s *Service) defender(ID string) error {
	s.mux.RLock()
	conn, ok := s.activeConnections[ID]
	if !ok {
		return fmt.Errorf("%s user is not available in active players", ID)
	}
	s.mux.RUnlock()
	for {
		// Acquire message
		t, b, err := conn.ReadMessage()
		if t == websocket.CloseMessage {
			s.closeConnection(ID)
			break
		}
		if err != nil {
			logrus.Error(err)
			s.closeConnection(ID)
			break
		}
		// Deserialize message
		var event dto.SocketEvent
		if err = json.Unmarshal(b, &event); err != nil {
			if stillActive := s.sendError(ID, "Could not parse Socket Event"); !stillActive {
				break
			}
			continue
		}
		// Process of valid events
		if event.Type == dto.SocketEventTypeDefendAction {
			var detailedEvent dto.DefendActionEvent
			if err = json.Unmarshal(b, &detailedEvent); err != nil {
				logrus.Error(err)
				if stillActive := s.sendError(ID, "Could not parse Defend Action Event"); !stillActive {
					break
				}
				continue
			}
			ongoingCombat, err := s.combatService.FindByID(detailedEvent.CombatID)
			if err != nil {
				logrus.Error(err)
				if stillActive := s.sendError(ID, "Invalid combat ID"); !stillActive {
					break
				}
				continue
			}
			if ongoingCombat.IsInFinalState() {
				if stillActive := s.sendError(ID, "Combat is already over, state is: "+ongoingCombat.CombatState); !stillActive {
					break
				}
				continue
			}
			attacker, err := s.playerService.FindByID(ongoingCombat.AttackerID)
			if err != nil {
				logrus.Error(err)
				if stillActive := s.sendError(ID, "Invalid attacker ID in Combat"); !stillActive {
					break
				}
				continue
			}
			if !attacker.Online {
				if _, err = s.combatService.UpdateCombatState(ongoingCombat.ID, combat.CombatStateAttackFailed); err != nil {
					logrus.Error(err)
				}
				if isConnectionStillAlive := s.sendResponseOrBreakConnection(ID, dto.AttackerFailedToAttackEvent{
					SocketEvent: dto.SocketEvent{
						Type: dto.SocketEventTypeAttackerFailedToAttack,
					},
					TargetID: ongoingCombat.ChallengeID,
					CombatID: ongoingCombat.ID,
				}); !isConnectionStillAlive {
					break
				}
			} else {
				if _, err = s.combatService.UpdateCombatState(ongoingCombat.ID, combat.CombatStateAttackerChallenged); err != nil {
					logrus.Error(err)
				}
				// if the connection is not alive here that's the problem of the potential
				// go routine handling the given defender
				s.sendResponseOrBreakConnection(attacker.ID, dto.AttackChallengeEvent{
					SocketEvent: dto.SocketEvent{
						Type: dto.SocketEventTypeSolutionEvaluationRequest,
					},
					TargetID: ongoingCombat.ChallengeID,
					Hints:    detailedEvent.Hints,
				})
			}
			continue
		}
		if event.Type == dto.SocketEventTypeSolutionEvaluation {
			var detailedEvent dto.SolutionEvaluationEvent
			if err = json.Unmarshal(b, &detailedEvent); err != nil {
				logrus.Error(err)
				if stillActive := s.sendError(ID, "Could not parse Defend Action Event"); !stillActive {
					break
				}
				continue
			}
			ongoingCombat, err := s.combatService.FindByID(detailedEvent.CombatID)
			if err != nil {
				logrus.Error(err)
				if stillActive := s.sendError(ID, "Invalid combat ID"); !stillActive {
					break
				}
				continue
			}
			if ongoingCombat.IsInFinalState() {
				if stillActive := s.sendError(ID, "Combat is already over, state is: "+ongoingCombat.CombatState); !stillActive {
					break
				}
				continue
			}
			attacker, err := s.playerService.FindByID(ongoingCombat.AttackerID)
			if err != nil {
				logrus.Error(err)
				if stillActive := s.sendError(ID, "Invalid attacker ID in Combat"); !stillActive {
					break
				}
				continue
			}
			// Here it does not really matter if the attacker is not online
			// worst case scenario the attacker does not receive the result
			// of the combat
			stateToUpdate := combat.CombatStateDefenseSucceeded
			if detailedEvent.Success {
				stateToUpdate = combat.CombatStateAttackSucceeded
			}
			if _, err = s.combatService.UpdateCombatState(ongoingCombat.ID, stateToUpdate); err != nil {
				logrus.Error(err)
			}
			if attacker.Online {
				s.sendResponseOrBreakConnection(attacker.ID, dto.AttackResultEvent{
					SocketEvent: dto.SocketEvent{
						Type: dto.SocketEventTypeSolutionEvaluationRequest,
					},
					TargetID: ongoingCombat.ChallengeID,
					Success:  detailedEvent.Success,
				})
			}
			continue
		}
		continue
	}
	return nil
}

// Constructor for engine service
func NewService(
	bus bus.IBus,
	playerService player.IService,
	challengeService challenge.IService,
	combatService combat.IService,
) IService {
	return &Service{
		bus,
		playerService,
		challengeService,
		combatService,
		make(map[string]*websocket.Conn),
		sync.RWMutex{},
	}
}
