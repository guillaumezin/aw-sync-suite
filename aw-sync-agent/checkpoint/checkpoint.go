package checkpoint

import (
	"aw-sync-agent/settings"
	"encoding/json"
	"log"
	"os"
	"time"
)

var checkpointFile = "checkpoint.json" // Default checkpoint file

type CheckpointEntry struct {
	Timestamp string `json:"timestamp"`
	LastID    int    `json:"last_id"`
}

// InitializeCheckpoint sets up the checkpoint file from configuration
func InitializeCheckpoint(config *settings.Configuration) {
	checkpointFile = config.Settings.CheckpointFile
}

func readAll() map[string]CheckpointEntry {
	file, err := os.OpenFile(checkpointFile, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error: Failed to open checkpoint file: %v", err)
	}
	defer file.Close()

	var checkpoints map[string]CheckpointEntry
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&checkpoints); err != nil {
		checkpoints = make(map[string]CheckpointEntry)
	}
	return checkpoints
}

// Read reads the checkpoint for the given watcher.
// Returns (nil, 0) if no checkpoint exists yet for this watcher.
func Read(watcher string) (*time.Time, int) {
	checkpoints := readAll()

	entry, ok := checkpoints[watcher]
	if !ok {
		return nil, 0
	}
	timestamp, err := time.Parse(time.RFC3339Nano, entry.Timestamp)
	if err != nil {
		log.Fatalf("Error: Failed to parse timestamp: %v", err)
	}
	return &timestamp, entry.LastID
}

// Update updates the checkpoint for the given watcher, storing both the
// timestamp (informational / for humans) and the last pushed event ID
// (the actual, reliable deduplication key).
func Update(watcher string, timestamp time.Time, lastID int) {
	file, err := os.OpenFile(checkpointFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error: Failed to open checkpoint file: %v", err)
	}
	defer file.Close()

	var checkpoints map[string]CheckpointEntry
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&checkpoints); err != nil {
		checkpoints = make(map[string]CheckpointEntry)
	}

	checkpoints[watcher] = CheckpointEntry{
		Timestamp: timestamp.Format("2006-01-02T15:04:05.000000000-07:00"), // nanoseconde, informatif uniquement
		LastID:    lastID,
	}

	file.Seek(0, 0)
	encoder := json.NewEncoder(file)
	if err = encoder.Encode(checkpoints); err != nil {
		log.Fatalf("Error: Failed to write to checkpoint file: %v", err)
	}
	pos, err := file.Seek(0, 1)
	if err != nil {
		log.Fatalf("Error: Failed to get the current size of the checkpoint file: %v", err)
	}
	if err = file.Truncate(pos); err != nil {
		log.Fatalf("Error: Failed to truncate checkpoint file: %v", err)
	}
}

func GetCheckpointFile() string {
	return checkpointFile
}

func SetCheckpointFile(file string) {
	checkpointFile = file
}
