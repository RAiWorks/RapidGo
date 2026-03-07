package service

import (
	"testing"
)

// TC-01: ParseMode — Valid Inputs
func TestParseMode_ValidInputs(t *testing.T) {
	tests := []struct {
		input    string
		expected Mode
	}{
		{"all", ModeAll},
		{"web", ModeWeb},
		{"api", ModeAPI},
		{"ws", ModeWS},
		{"api,ws", ModeAPI | ModeWS},
		{"web,api", ModeWeb | ModeAPI},
		{"web,api,ws", ModeAll},
		{" api , ws ", ModeAPI | ModeWS}, // whitespace
		{"ALL", ModeAll},                  // case insensitive
		{"Api", ModeAPI},                  // mixed case
		{"WEB,API", ModeWeb | ModeAPI},    // uppercase combo
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			mode, err := ParseMode(tt.input)
			if err != nil {
				t.Fatalf("ParseMode(%q) returned error: %v", tt.input, err)
			}
			if mode != tt.expected {
				t.Errorf("ParseMode(%q) = %d, want %d", tt.input, mode, tt.expected)
			}
		})
	}
}

// TC-02: ParseMode — Invalid Inputs
func TestParseMode_InvalidInputs(t *testing.T) {
	tests := []string{
		"",
		"invalid",
		"api,invalid",
		"worker",
		",",
		"   ",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			mode, err := ParseMode(input)
			if err == nil {
				t.Fatalf("ParseMode(%q) should return error, got mode=%d", input, mode)
			}
			if mode != 0 {
				t.Errorf("ParseMode(%q) returned mode=%d on error, want 0", input, mode)
			}
		})
	}
}

// TC-03: Mode Bitmask Operations
func TestMode_Has(t *testing.T) {
	if !ModeAPI.Has(ModeAPI) {
		t.Error("ModeAPI.Has(ModeAPI) should be true")
	}
	if ModeAPI.Has(ModeWeb) {
		t.Error("ModeAPI.Has(ModeWeb) should be false")
	}
	if !ModeAll.Has(ModeWeb) {
		t.Error("ModeAll.Has(ModeWeb) should be true")
	}
	if !ModeAll.Has(ModeAPI) {
		t.Error("ModeAll.Has(ModeAPI) should be true")
	}
	if !ModeAll.Has(ModeWS) {
		t.Error("ModeAll.Has(ModeWS) should be true")
	}

	combined := ModeAPI | ModeWS
	if !combined.Has(ModeAPI) {
		t.Error("(ModeAPI|ModeWS).Has(ModeAPI) should be true")
	}
	if !combined.Has(ModeWS) {
		t.Error("(ModeAPI|ModeWS).Has(ModeWS) should be true")
	}
	if combined.Has(ModeWeb) {
		t.Error("(ModeAPI|ModeWS).Has(ModeWeb) should be false")
	}
}

func TestMode_String(t *testing.T) {
	tests := []struct {
		mode     Mode
		expected string
	}{
		{ModeAll, "all"},
		{ModeWeb, "web"},
		{ModeAPI, "api"},
		{ModeWS, "ws"},
		{ModeAPI | ModeWS, "api,ws"},
		{ModeWeb | ModeAPI, "web,api"},
		{Mode(0), "none"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.mode.String(); got != tt.expected {
				t.Errorf("Mode(%d).String() = %q, want %q", tt.mode, got, tt.expected)
			}
		})
	}
}

func TestMode_Services(t *testing.T) {
	tests := []struct {
		mode     Mode
		expected []Mode
	}{
		{ModeWeb, []Mode{ModeWeb}},
		{ModeAPI, []Mode{ModeAPI}},
		{ModeWS, []Mode{ModeWS}},
		{ModeAPI | ModeWS, []Mode{ModeAPI, ModeWS}},
		{ModeAll, []Mode{ModeWeb, ModeAPI, ModeWS}},
	}

	for _, tt := range tests {
		t.Run(tt.mode.String(), func(t *testing.T) {
			got := tt.mode.Services()
			if len(got) != len(tt.expected) {
				t.Fatalf("Services() returned %d items, want %d", len(got), len(tt.expected))
			}
			for i, m := range got {
				if m != tt.expected[i] {
					t.Errorf("Services()[%d] = %d, want %d", i, m, tt.expected[i])
				}
			}
		})
	}
}

func TestMode_Services_Empty(t *testing.T) {
	got := Mode(0).Services()
	if len(got) != 0 {
		t.Errorf("Mode(0).Services() should return empty slice, got %d items", len(got))
	}
}

func TestMode_PortEnvKey(t *testing.T) {
	tests := []struct {
		mode     Mode
		expected string
	}{
		{ModeWeb, "WEB_PORT"},
		{ModeAPI, "API_PORT"},
		{ModeWS, "WS_PORT"},
		{ModeAll, "APP_PORT"},
		{ModeAPI | ModeWS, "APP_PORT"}, // combined falls back to APP_PORT
	}

	for _, tt := range tests {
		t.Run(tt.mode.String(), func(t *testing.T) {
			if got := tt.mode.PortEnvKey(); got != tt.expected {
				t.Errorf("Mode(%d).PortEnvKey() = %q, want %q", tt.mode, got, tt.expected)
			}
		})
	}
}
