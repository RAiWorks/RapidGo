package validation

import (
	"fmt"
	"net"
	"net/mail"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Errors holds validation error messages keyed by field name.
type Errors map[string][]string

// HasErrors returns true if any validation errors exist.
func (e Errors) HasErrors() bool { return len(e) > 0 }

// Add appends an error message for a field.
func (e Errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// First returns the first error message for a field, or empty string.
func (e Errors) First(field string) string {
	if msgs, ok := e[field]; ok && len(msgs) > 0 {
		return msgs[0]
	}
	return ""
}

// Validator performs fluent, chainable input validation.
type Validator struct {
	errors Errors
}

// New creates a new Validator instance.
func New() *Validator {
	return &Validator{errors: make(Errors)}
}

// Required checks that the value is not empty after trimming whitespace.
func (v *Validator) Required(field, value string) *Validator {
	if strings.TrimSpace(value) == "" {
		v.errors.Add(field, field+" is required")
	}
	return v
}

// MinLength checks minimum string length (rune count).
func (v *Validator) MinLength(field, value string, min int) *Validator {
	if utf8.RuneCountInString(value) < min {
		v.errors.Add(field, fmt.Sprintf("%s must be at least %d characters", field, min))
	}
	return v
}

// MaxLength checks maximum string length (rune count).
func (v *Validator) MaxLength(field, value string, max int) *Validator {
	if utf8.RuneCountInString(value) > max {
		v.errors.Add(field, fmt.Sprintf("%s must be at most %d characters", field, max))
	}
	return v
}

// Email validates email format using net/mail.
func (v *Validator) Email(field, value string) *Validator {
	if _, err := mail.ParseAddress(value); err != nil {
		v.errors.Add(field, field+" must be a valid email")
	}
	return v
}

// URL validates that value starts with http:// or https://.
func (v *Validator) URL(field, value string) *Validator {
	if !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
		v.errors.Add(field, field+" must be a valid URL")
	}
	return v
}

// Matches checks value against a regex pattern.
func (v *Validator) Matches(field, value, pattern string) *Validator {
	if matched, _ := regexp.MatchString(pattern, value); !matched {
		v.errors.Add(field, field+" format is invalid")
	}
	return v
}

// In checks value is one of the allowed options.
func (v *Validator) In(field, value string, allowed []string) *Validator {
	for _, a := range allowed {
		if value == a {
			return v
		}
	}
	v.errors.Add(field, fmt.Sprintf("%s must be one of: %s", field, strings.Join(allowed, ", ")))
	return v
}

// Confirmed checks two values match (e.g. password confirmation).
func (v *Validator) Confirmed(field, value, confirmation string) *Validator {
	if value != confirmation {
		v.errors.Add(field, field+" confirmation does not match")
	}
	return v
}

// IP validates an IP address.
func (v *Validator) IP(field, value string) *Validator {
	if net.ParseIP(value) == nil {
		v.errors.Add(field, field+" must be a valid IP address")
	}
	return v
}

// Errors returns the collected validation errors.
func (v *Validator) Errors() Errors { return v.errors }

// Valid returns true if no validation errors exist.
func (v *Validator) Valid() bool { return !v.errors.HasErrors() }
