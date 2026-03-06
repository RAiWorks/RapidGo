# 🧪 Test Plan: Input Validation

> **Feature**: `23` — Input Validation
> **Total Test Cases**: 22

---

## Required

| ID | Test | Input | Expected |
|----|------|-------|----------|
| TC-01 | Required fails empty | `""` | Error: "name is required" |
| TC-02 | Required fails whitespace | `"   "` | Error: "name is required" |
| TC-03 | Required passes | `"alice"` | No error |

## MinLength / MaxLength

| ID | Test | Input | Expected |
|----|------|-------|----------|
| TC-04 | MinLength fails below | `"ab"`, min=3 | Error |
| TC-05 | MinLength passes at minimum | `"abc"`, min=3 | No error |
| TC-06 | MaxLength fails above | `"abcdef"`, max=5 | Error |
| TC-07 | MaxLength passes at max | `"abcde"`, max=5 | No error |

## Email

| ID | Test | Input | Expected |
|----|------|-------|----------|
| TC-08 | Email fails invalid | `"not-an-email"` | Error |
| TC-09 | Email passes valid | `"user@example.com"` | No error |

## URL

| ID | Test | Input | Expected |
|----|------|-------|----------|
| TC-10 | URL fails non-URL | `"ftp://example.com"` | Error |
| TC-11 | URL passes valid | `"https://example.com"` | No error |

## Matches

| ID | Test | Input | Expected |
|----|------|-------|----------|
| TC-12 | Matches fails | `"abc"`, pattern `^[0-9]+$` | Error |
| TC-13 | Matches passes | `"123"`, pattern `^[0-9]+$` | Error: none |

## In

| ID | Test | Input | Expected |
|----|------|-------|----------|
| TC-14 | In fails disallowed | `"purple"`, allowed `["red","green","blue"]` | Error |
| TC-15 | In passes allowed | `"red"`, allowed `["red","green","blue"]` | No error |

## Confirmed

| ID | Test | Input | Expected |
|----|------|-------|----------|
| TC-16 | Confirmed fails mismatch | `"pass1"` vs `"pass2"` | Error |
| TC-17 | Confirmed passes match | `"pass1"` vs `"pass1"` | No error |

## IP

| ID | Test | Input | Expected |
|----|------|-------|----------|
| TC-18 | IP fails invalid | `"999.999.999.999"` | Error |
| TC-19 | IP passes valid | `"192.168.1.1"` | No error |

## Errors Type & Integration

| ID | Test | Input | Expected |
|----|------|-------|----------|
| TC-20 | Errors.Add and Errors.First | Add two errors to same field | First returns first one |
| TC-21 | Chaining multiple rules | Required + MinLength + Email on empty | Multiple errors accumulated |
| TC-22 | Valid returns true | No rules fail | `Valid() == true` |

---

## Acceptance Criteria

1. All 22 tests pass
2. Full regression (`go test ./... -count=1`) — 0 failures
3. `go vet ./...` — clean
4. No new dependencies in `go.mod`
