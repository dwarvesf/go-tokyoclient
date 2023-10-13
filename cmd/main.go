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
type SimpleBot struct{}

// Consume consumes the event
func (c *SimpleBot) HandleEvent(userID int, teammates map[int]string, gamePad tokyoclient.GamePad, state tokyoclient.GameState) error {
	fmt.Println("User ID:", userID)

	err := gamePad.Throttle(1)
	if err != nil {
		fmt.Println("Error throttling:", err)
	}
	// random rotation
	angle := rand.Float64() * 2 * math.Pi

	fmt.Println("Rotating by:", angle)
	gamePad.Rotate(angle)
	if err != nil {
		fmt.Println("Error throttling:", err)
	}
	err = gamePad.Fire()
	if err != nil {
		fmt.Println("Error firing:", err)
		// return err
	}
	return nil
}
