package models

import "github.com/RAiWorks/RapidGo/core/audit"

// AuditLog is an alias for the canonical definition in core/audit.
// This keeps existing code that references models.AuditLog working.
type AuditLog = audit.AuditLog
