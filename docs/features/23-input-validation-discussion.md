# 💬 Discussion: Input Validation

> **Feature**: `23` — Input Validation
> **Status**: FINAL
> **Date**: 2026-03-06

---

## 1. What Are We Building?

A built-in, zero-dependency validation engine in `core/validation/` that provides a fluent, chainable API for validating user input. This covers the most common validation scenarios without requiring any external libraries.

## 2. Why?

Every web framework needs input validation. While Gin already integrates `go-playground/validator` for struct-based binding, the framework needs its own built-in validator for:

- **SSR form validation** — validate individual form fields without struct tags
- **Flexible validation** — chain rules fluently per field
- **Zero-dependency option** — no external libraries required for common cases
- **Consistent error format** — field-keyed error maps usable in JSON and template rendering

## 3. Scope

### In Scope

- `Errors` type — `map[string][]string` with `HasErrors()`, `Add()`, `First()`
- `Validator` struct — fluent builder with `New()`, `Errors()`, `Valid()`
- 9 validation methods: `Required`, `MinLength`, `MaxLength`, `Email`, `URL`, `Matches`, `In`, `Confirmed`, `IP`
- All methods return `*Validator` for chaining
- All stdlib imports only

### Out of Scope

- Struct-based validation (handled by Gin's binding system, already available)
- Request structs in `http/requests/` (separate feature or user-land code)
- Custom rule registration / extensibility (can be added later)
- Database-aware rules (e.g., unique email check)

## 4. Dependencies

- **#07 (Router)** — validation is used in route handlers
- **#15 (Controllers)** — controllers call the validator

Both are already shipped.

## 5. Key Decisions

| Decision | Choice | Rationale |
|---|---|---|
| Package location | `core/validation/` | Blueprint specifies this; directory already exists with `.gitkeep` |
| Single file | `validation.go` | 9 methods + Errors type fits in one file |
| String-only values | All methods take `string` values | Matches form/query input (always strings); struct-based validation covers typed fields |
| Error messages | Hardcoded, English | Simple first pass; i18n can be layered later |
| Chaining | Every method returns `*Validator` | Blueprint pattern; enables `v.Required().MinLength().Email()` |
