package session

import "time"

// Store defines the contract every session backend must satisfy.
type Store interface {
	// Read returns session data for the given ID, or empty map if not found.
	Read(id string) (map[string]interface{}, error)

	// Write persists session data with the given ID and lifetime.
	Write(id string, data map[string]interface{}, lifetime time.Duration) error

	// Destroy removes a session by ID.
	Destroy(id string) error

	// GC removes sessions older than the given max lifetime.
	GC(maxLifetime time.Duration) error
}
