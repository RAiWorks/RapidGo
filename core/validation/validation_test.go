package validation

import (
	"testing"
)

// TC-01: Required fails on empty string
func TestRequiredEmpty(t *testing.T) {
	v := New().Required("name", "")
	if v.Valid() {
		t.Fatal("expected error for empty string")
	}
	if v.Errors().First("name") != "name is required" {
		t.Fatalf("unexpected message: %s", v.Errors().First("name"))
	}
}

// TC-02: Required fails on whitespace-only
func TestRequiredWhitespace(t *testing.T) {
	v := New().Required("name", "   ")
	if v.Valid() {
		t.Fatal("expected error for whitespace-only")
	}
}

// TC-03: Required passes on non-empty
func TestRequiredValid(t *testing.T) {
	v := New().Required("name", "alice")
	if !v.Valid() {
		t.Fatal("expected no error")
	}
}

// TC-04: MinLength fails below minimum
func TestMinLengthFail(t *testing.T) {
	v := New().MinLength("name", "ab", 3)
	if v.Valid() {
		t.Fatal("expected error for short string")
	}
}

// TC-05: MinLength passes at minimum
func TestMinLengthPass(t *testing.T) {
	v := New().MinLength("name", "abc", 3)
	if !v.Valid() {
		t.Fatal("expected no error at exact minimum")
	}
}

// TC-06: MaxLength fails above maximum
func TestMaxLengthFail(t *testing.T) {
	v := New().MaxLength("name", "abcdef", 5)
	if v.Valid() {
		t.Fatal("expected error for long string")
	}
}

// TC-07: MaxLength passes at maximum
func TestMaxLengthPass(t *testing.T) {
	v := New().MaxLength("name", "abcde", 5)
	if !v.Valid() {
		t.Fatal("expected no error at exact maximum")
	}
}

// TC-08: Email fails on invalid
func TestEmailFail(t *testing.T) {
	v := New().Email("email", "not-an-email")
	if v.Valid() {
		t.Fatal("expected error for invalid email")
	}
}

// TC-09: Email passes on valid
func TestEmailPass(t *testing.T) {
	v := New().Email("email", "user@example.com")
	if !v.Valid() {
		t.Fatal("expected no error for valid email")
	}
}

// TC-10: URL fails on non-URL
func TestURLFail(t *testing.T) {
	v := New().URL("website", "ftp://example.com")
	if v.Valid() {
		t.Fatal("expected error for non-http URL")
	}
}

// TC-11: URL passes on valid URL
func TestURLPass(t *testing.T) {
	v := New().URL("website", "https://example.com")
	if !v.Valid() {
		t.Fatal("expected no error for valid URL")
	}
}

// TC-12: Matches fails on non-match
func TestMatchesFail(t *testing.T) {
	v := New().Matches("code", "abc", `^[0-9]+$`)
	if v.Valid() {
		t.Fatal("expected error for non-matching pattern")
	}
}

// TC-13: Matches passes on match
func TestMatchesPass(t *testing.T) {
	v := New().Matches("code", "123", `^[0-9]+$`)
	if !v.Valid() {
		t.Fatal("expected no error for matching pattern")
	}
}

// TC-14: In fails on disallowed value
func TestInFail(t *testing.T) {
	v := New().In("color", "purple", []string{"red", "green", "blue"})
	if v.Valid() {
		t.Fatal("expected error for disallowed value")
	}
}

// TC-15: In passes on allowed value
func TestInPass(t *testing.T) {
	v := New().In("color", "red", []string{"red", "green", "blue"})
	if !v.Valid() {
		t.Fatal("expected no error for allowed value")
	}
}

// TC-16: Confirmed fails on mismatch
func TestConfirmedFail(t *testing.T) {
	v := New().Confirmed("password", "pass1", "pass2")
	if v.Valid() {
		t.Fatal("expected error for mismatched confirmation")
	}
}

// TC-17: Confirmed passes on match
func TestConfirmedPass(t *testing.T) {
	v := New().Confirmed("password", "pass1", "pass1")
	if !v.Valid() {
		t.Fatal("expected no error for matching confirmation")
	}
}

// TC-18: IP fails on invalid IP
func TestIPFail(t *testing.T) {
	v := New().IP("ip", "999.999.999.999")
	if v.Valid() {
		t.Fatal("expected error for invalid IP")
	}
}

// TC-19: IP passes on valid IP
func TestIPPass(t *testing.T) {
	v := New().IP("ip", "192.168.1.1")
	if !v.Valid() {
		t.Fatal("expected no error for valid IP")
	}
}

// TC-20: Errors.Add and Errors.First
func TestErrorsAddAndFirst(t *testing.T) {
	e := make(Errors)
	e.Add("field", "first error")
	e.Add("field", "second error")
	if e.First("field") != "first error" {
		t.Fatalf("expected 'first error', got %q", e.First("field"))
	}
	if e.First("missing") != "" {
		t.Fatal("expected empty string for missing field")
	}
}

// TC-21: Chaining multiple rules accumulates errors
func TestChainingAccumulatesErrors(t *testing.T) {
	v := New()
	v.Required("email", "").MinLength("email", "", 5).Email("email", "")
	errs := v.Errors()["email"]
	if len(errs) < 2 {
		t.Fatalf("expected multiple errors, got %d", len(errs))
	}
}

// TC-22: Valid returns true when no errors
func TestValidNoErrors(t *testing.T) {
	v := New()
	if !v.Valid() {
		t.Fatal("expected Valid() to be true for fresh validator")
	}
}
