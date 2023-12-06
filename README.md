# go-tokyoclient

`go-tokyoclient` is a Go client library for connecting to a Tokyo game server and controlling a ship in the game. It provides functionality to interact with the game server, receive events, and control the ship using a simple bot implementation.

## Usage

To use `go-tokyoclient`, follow these steps:

1. Install the package using `go get`:
```sh
go get -u github.com/dwarvesf/go-tokyoclient
```

2. Import the package in your Go code:
```go
import "github.com/dwarvesf/go-tokyoclient"
```

3. Implement a Bot
```go
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
```

4. Use the provided example code to start clients and control the ships in the game. Replace the server host and user key with the actual values.

```go
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
	userName   = "andy"        // Replace with your display name
	roomToken  = "wCimjaw"        // Replace with your display name
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	fmt.Println("Starting client:", userKey, userName)
	cfg := &tokyoclient.Config{
		ServerURL: serverHost,
		RoomToken: roomToken,
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
```

5. Run your Go program and observe the clients interacting with the Tokyo game server.

## Interfaces

### EventConsumer

`EventConsumer` is an interface that the client must implement to handle events received from the server. It defines the `HandleEvent` method, which is responsible for processing events and taking appropriate actions based on the event information.

```go
type EventConsumer interface {
	// HandleEvent handles the event from the server
	HandleEvent(userID int, teammates map[int]string, gamePad GamePad, state GameState) error
}
```

### GamePad

`GamePad` is an interface provided by the library to the client. It offers methods to control the ship in the game. The `GamePad` interface includes methods to rotate the ship (`Rotate`), adjust throttle (`Throttle`), and fire bullets (`Fire`).

```go
type GamePad interface {
	// Rotate rotates the ship
	Rotate(angle float64) error

	// Throttle throttles the ship
	Throttle(speed float64) error

	// Fire fires the bullet
	Fire() error
}
```

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests for any improvements or bug fixes.
