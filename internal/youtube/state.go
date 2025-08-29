package youtube

import (
	"encoding/json"
	"os"
)

// State holds the processing state
type State struct {
	Mode   string `json:"mode"`   // "init" or "backfillDone"
	LastID string `json:"lastId"` // last processed comment ID
}

// LoadState loads state.json from configs directory
func (c *Client) LoadState() State {
	data, err := os.ReadFile("configs/state.json")
	if err != nil {
		return State{Mode: "init", LastID: ""}
	}

	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return State{Mode: "init", LastID: ""}
	}
	return s
}

// SaveState saves state.json to configs directory
func (c *Client) SaveState(s State) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("configs/state.json", data, 0644)
}
