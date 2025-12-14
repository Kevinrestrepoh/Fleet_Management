package types

import "fmt"

type CommandRequest struct {
	Type  string `json:"type"`
	Value uint32 `json:"value"`
}

func (c *CommandRequest) Validate() error {
	switch c.Type {
	case "UPDATE_RATE", "PING", "SHUTDOWN":
		return nil
	default:
		return fmt.Errorf("invalid command type")
	}
}
