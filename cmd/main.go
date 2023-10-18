package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/signal"

	"github.com/dwarvesf/go-tokyoclient"
)

const (
	serverHost = "localhost:8091" // Replace with the actual server host
	userKey    = "webuild"        // Replace with your unique user key
	userName   = "andy"           // Replace with your display name
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	fmt.Println("Starting client:", userKey, userName)
	cfg := &tokyoclient.Config{
		ServerURL: serverHost,
		APIKey:    userKey,
		UserName:  userName,
	}
	client := tokyoclient.NewClient(&SimpleBot{})
	gamePad, err := client.Start(*cfg)
	if err != nil {
		fmt.Println("Error starting client:", err)
		return
	}
	gamePad.Fire()

	<-interrupt
}

// SimpleBot is a simple implementation of EventConsumer
type SimpleBot struct {
	ReadyShot bool
}

// Consume consumes the event
func (c *SimpleBot) HandleEvent(userID int, teammates map[int]string, gamePad tokyoclient.GamePad, state tokyoclient.GameState) error {
	fmt.Println("User ID:", userID)

	// get user info
	var user tokyoclient.Player
	for _, u := range state.Players {
		if u.ID == userID {
			user = u
			break
		}
	}
	enemy := findNearestEnemy(user, state.Players)
	angle := tokyoclient.CalculateAngle(tokyoclient.Point{X: user.X, Y: user.Y}, tokyoclient.Point{X: enemy.X, Y: enemy.Y})
	distance := tokyoclient.DistanceBetween(tokyoclient.Point{X: user.X, Y: user.Y}, tokyoclient.Point{X: enemy.X, Y: enemy.Y})

	if distance <= 400 {
		c.ReadyShot = false
		return gamePad.Rotate(angle + math.Pi/2)
	}

	// random action for each stick
	n := rand.Intn(3)
	switch n {
	case 0:
		c.ReadyShot = true
		// random add small angle
		angle += (rand.Float64() - 0.5) * math.Pi / 4
		return gamePad.Rotate(angle)
	case 1:
		if c.ReadyShot {
			return gamePad.Fire()
		}
	case 2:
		return gamePad.Throttle(1)
	}

	return nil
}

func findNearestEnemy(user tokyoclient.Player, enemies []tokyoclient.Player) tokyoclient.Player {
	var nearestEnemy tokyoclient.Player
	minDistance := math.MaxFloat64
	for _, enemy := range enemies {
		distance := tokyoclient.DistanceBetween(tokyoclient.Point{X: user.X, Y: user.Y}, tokyoclient.Point{X: enemy.X, Y: enemy.Y})
		if distance < minDistance && enemy.ID != user.ID {
			minDistance = distance
			nearestEnemy = enemy
		}
	}

	return nearestEnemy
}
