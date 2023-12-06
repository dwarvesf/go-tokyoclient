package tokyoclient

// Config is the configuration for the client
type Config struct {
	// ServerURL is the URL of the server
	ServerURL string
	// RoomToken is the token of the room
	RoomToken string
	// APIKey is the unique key
	APIKey string
	// UserName is the display name
	UserName string
}

const (
	socketPath = "/socket"
)

func (c *Config) GetServerURL() string {
	return "ws://" + c.ServerURL + socketPath + "?key=" + c.APIKey + "&name=" + c.UserName + "&room_token=" + c.RoomToken
}
