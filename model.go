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
	Items      []Item       `json:"items"`
	Dead       []DeadPlayer `json:"dead"`
	Bullets    []Bullet     `json:"bullets"`
	Scoreboard map[int]int  `json:"scoreboard"`
}

// Point is a point in the game
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Player is the player state sent by the server
type Player struct {
	ID int `json:"id"`
	Point
	Angle        float64 `json:"angle"`
	Throttle     float64 `json:"throttle"`
	Radius       float64 `json:"radius"`
	BulletRadius float64 `json:"bullet_radius"`
	BulletSpeed  float64 `json:"bullet_speed"`
	BulletLimit  int     `json:"bullet_limit"`
}

type Item struct {
	ID int `json:"id"`
	Point
	Radius   float64 `json:"radius"`
	ItemType string  `json:"item_type"`
}

// Bullet is the bullet state sent by the server
type Bullet struct {
	ID int `json:"id"`
	Point
	PlayerID int     `json:"player_id"`
	Angle    float64 `json:"angle"`
	Radius   float64 `json:"radius"`
	Speed    float64 `json:"speed"`
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

type StateEvent struct {
	Data GameState `json:"data"`
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
		var stateData StateEvent
		if err := json.Unmarshal([]byte(eventData), &stateData); err != nil {
			return nil, fmt.Errorf("error parsing state event data: %v", err)
		}
		event.Data = stateData.Data

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
