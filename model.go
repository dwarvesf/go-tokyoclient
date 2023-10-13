package tokyoclient

import (
	"encoding/json"
	"fmt"
)

const (
	RotateAction   = "rotate"
	ThrottleAction = "throttle"
	FireAction     = "fire"
)

// GameState is the game state sent by the server
type GameState struct {
	Bounds     []float64    `json:"bounds"`
	Players    []Player     `json:"players"`
	Bullets    []Bullet     `json:"bullets"`
	Dead       []DeadPlayer `json:"dead"`
	Scoreboard map[int]int  `json:"scoreboard"`
}

// Player is the player state sent by the server
type Player struct {
	ID       int     `json:"id"`
	Angle    float64 `json:"angle"`
	Throttle float64 `json:"throttle"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
}

// Bullet is the bullet state sent by the server
type Bullet struct {
	ID       int     `json:"id"`
	PlayerID int     `json:"player_id"`
	Angle    float64 `json:"angle"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
}

// Bot is the interface that the bot must implement
type DeadPlayer struct {
	Respawn RespawnInfo `json:"respawn"`
	Player  Player      `json:"player"`
}

// RespawnInfo is the respawn info sent by the server
type RespawnInfo struct {
	SecsSinceEpoch  int64 `json:"secs_since_epoch"`
	NanosSinceEpoch int64 `json:"nanos_since_epoch"`
}

type EventData interface{}

type Event struct {
	Event string    `json:"e"`
	Data  EventData `json:"data"`
}

type IDEvent struct {
	ID int `json:"data"`
}

type TeamNamesEvent struct {
	Names map[int]string `json:"data"`
}

// ParseEvent parses the event data
func ParseEvent(eventData []byte) (*Event, error) {
	var event Event
	if err := json.Unmarshal(eventData, &event); err != nil {
		return nil, fmt.Errorf("error parsing event: %v", err)
	}

	switch event.Event {
	case "state":
		var stateData GameState
		if err := json.Unmarshal([]byte(eventData), &stateData); err != nil {
			return nil, fmt.Errorf("error parsing state event data: %v", err)
		}
		event.Data = stateData

	case "id":
		var idData IDEvent
		if err := json.Unmarshal([]byte(eventData), &idData); err != nil {
			return nil, fmt.Errorf("error parsing id event data: %v", err)
		}
		event.Data = idData

	case "teamnames":
		var teamNamesData TeamNamesEvent
		if err := json.Unmarshal([]byte(eventData), &teamNamesData); err != nil {
			return nil, fmt.Errorf("error parsing teamnames event data: %v", err)
		}
		event.Data = teamNamesData

	default:
		return nil, fmt.Errorf("unknown event type: %s", event.Event)
	}

	return &event, nil
}
