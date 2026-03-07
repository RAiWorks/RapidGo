package service

import (
	"fmt"
	"strings"
)

// Mode represents which services the application should run.
// Uses bitmask for easy combination.
type Mode uint8

const (
	ModeWeb Mode = 1 << iota // Web SSR (templates, static files)
	ModeAPI                  // JSON API endpoints
	ModeWS                   // WebSocket service

	ModeAll = ModeWeb | ModeAPI | ModeWS // Monolith — all HTTP services
)

// modeNames maps string identifiers to Mode constants.
var modeNames = map[string]Mode{
	"web": ModeWeb,
	"api": ModeAPI,
	"ws":  ModeWS,
	"all": ModeAll,
}

// ParseMode parses a comma-separated mode string into a Mode bitmask.
// Valid inputs: "all", "web", "api", "ws", "api,ws", "web,api", etc.
// Returns error for empty or invalid mode strings.
func ParseMode(s string) (Mode, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return 0, fmt.Errorf("service mode cannot be empty")
	}

	var m Mode
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		flag, ok := modeNames[part]
		if !ok {
			return 0, fmt.Errorf("invalid service mode: %q (valid: all, web, api, ws)", part)
		}
		m |= flag
	}

	if m == 0 {
		return 0, fmt.Errorf("service mode cannot be empty")
	}
	return m, nil
}

// Has returns true if the mode includes the given flag.
func (m Mode) Has(flag Mode) bool {
	return m&flag != 0
}

// String returns a human-readable representation of the mode.
func (m Mode) String() string {
	if m == ModeAll {
		return "all"
	}
	var parts []string
	if m.Has(ModeWeb) {
		parts = append(parts, "web")
	}
	if m.Has(ModeAPI) {
		parts = append(parts, "api")
	}
	if m.Has(ModeWS) {
		parts = append(parts, "ws")
	}
	if len(parts) == 0 {
		return "none"
	}
	return strings.Join(parts, ",")
}

// Services returns the list of individual modes active in this bitmask.
func (m Mode) Services() []Mode {
	var s []Mode
	if m.Has(ModeWeb) {
		s = append(s, ModeWeb)
	}
	if m.Has(ModeAPI) {
		s = append(s, ModeAPI)
	}
	if m.Has(ModeWS) {
		s = append(s, ModeWS)
	}
	return s
}

// PortEnvKey returns the environment variable name for the port of a single-mode constant.
func (m Mode) PortEnvKey() string {
	switch m {
	case ModeWeb:
		return "WEB_PORT"
	case ModeAPI:
		return "API_PORT"
	case ModeWS:
		return "WS_PORT"
	default:
		return "APP_PORT"
	}
}
