package tokyoclient

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	pingPeriod  = 30 * time.Millisecond
	writeWait   = 10 * time.Second
	readTimeout = 60 * time.Second
)

// EventConsumer is the interface that the client must implement
type EventConsumer interface {
	// HandleEvent handles the event from the server
	HandleEvent(userID int, teammates map[int]string, controller Controller, state GameState) error
}

// Controller is the interface that lib provides to the client to interact with the game client
type Controller interface {
	// Rotate rotates the ship
	Rotate(angle float64) error
	// Throttle throttles the ship
	Throttle(speed float64) error
	// Fire fires the bullet
	Fire() error
}

// Client is the client to interact with the server
type Client struct {
	userID     int
	teammates  map[int]string
	conn       *websocket.Conn
	dispatcher EventConsumer
	mu         sync.Mutex
}

// NewClient creates a new client
func NewClient(dispatch EventConsumer) *Client {
	return &Client{
		teammates:  make(map[int]string),
		dispatcher: dispatch,
		mu:         sync.Mutex{},
	}
}

func (c *Client) Controller() (Controller, error) {
	if c.conn == nil {
		return nil, ErrConnNotInitialized
	}
	return c, nil
}

func (c *Client) writeMessage(message string) error {
	if c.conn == nil {
		return ErrConnNotInitialized
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (c *Client) Start(cfg Config) (Controller, error) {
	startedSignal := make(chan struct{})
	go func() {
		c.start(cfg, startedSignal)
	}()
	<-startedSignal

	return c, nil
}

func (c *Client) start(cfg Config, startedSignal chan struct{}) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u, err := url.Parse(cfg.GetServerURL())
	if err != nil {
		log.Fatal(err)
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	c.conn = conn
	defer conn.Close()
	startedSignal <- struct{}{}

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			err = c.HandleMessage(message)
			if err != nil {
				log.Println(err)
				log.Println(string(message))
				return
			}
		}
	}()

	for {
		select {
		case <-done:
			return nil
		case <-interrupt:
			log.Println("Interrupt received. Closing connection.")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return nil
			}
			select {
			case <-done:
			case <-time.After(writeWait):
			}
			return nil
		}
	}
}

// UpdateTeammates updates the teammates
func (c *Client) UpdateTeammates(teammates map[int]string) {
	c.teammates = teammates
}

// UpdateID updates the ID
func (c *Client) UpdateID(id int) {
	c.userID = id
}

// HandleMessage handles the message from the server
func (c *Client) HandleMessage(message []byte) error {
	e, err := ParseEvent(message)
	if err != nil {
		return err
	}

	switch e.Event {
	case "id":
		dt, ok := e.Data.(IDEvent)
		if !ok {
			return ErrInvalidID
		}
		c.UpdateID(dt.ID)
	case "teamnames":
		dt, ok := e.Data.(TeamNamesEvent)
		if !ok {
			return ErrInvalidTeamNames
		}
		c.UpdateTeammates(dt.Names)

	case "state":
		state, ok := e.Data.(GameState)
		if !ok {
			return ErrInvalidState
		}

		if c.dispatcher != nil {
			if err := c.dispatcher.HandleEvent(c.userID, c.teammates, c, state); err != nil {
				return err
			}
		}
	}

	return nil
}

// Rotate rotates the ship
func (c *Client) Rotate(angle float64) error {
	return c.writeMessage(`{"e":"rotate","data":` + strconv.FormatFloat(angle, 'f', -1, 64) + `}`)
}

// Throttle throttles the ship
func (c *Client) Throttle(speed float64) error {
	return c.writeMessage(`{"e":"throttle","data":` + strconv.FormatFloat(speed, 'f', -1, 64) + `}`)
}

// Fire fires the bullet
func (c *Client) Fire() error {
	return c.writeMessage(`{"e":"fire"}`)
}
