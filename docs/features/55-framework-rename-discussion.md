# 💬 Discussion: Framework Rename / Rebrand

> **Feature**: `55` — Framework Rename / Rebrand
> **Status**: � COMPLETE
> **Branch**: `feature/55-framework-rename`
> **Depends On**: All 41 features complete (no active feature branches)
> **Date Started**: 2026-03-07
> **Date Completed**: 2026-03-07

---

## Summary

Rename the RGo framework to a new identity that better communicates the framework's value proposition and secures proper domain presence. The rename touches every layer — module path, package references, CLI binary, documentation, templates, Docker config, and repository metadata. This is a cross-cutting change that must be executed atomically to avoid a half-renamed state.

---

## Functional Requirements

- As a framework author, I want a name with strong domain availability so that the framework has a professional web presence
- As a developer discovering the framework, I want the name to immediately communicate what it is (Go + speed) so I can evaluate it quickly
- As a developer using the framework, I want the CLI command and module path to reflect the new name consistently
- As a developer reading docs, I want all references updated so there's no confusion between old and new names

## Current State / Reference

### What Exists

- **Module path**: `github.com/RAiWorks/RGo`
- **Repository**: `https://github.com/RAiWorks/RGo`
- **CLI binary**: `bin/rgo.exe`
- **Version constant**: `const Version = "0.1.0"` in `core/cli/root.go`
- **Root command**: `rootCmd` with `Use: "rgo"` in `core/cli/root.go`
- **References**: "RGo" appears across all 41 feature docs, README.md, project-context.md, project-roadmap.md, welcome page template, controller comments, provider comments, test files
- **Docker**: `Dockerfile`, `docker-compose.yml`, `Caddyfile` reference the binary and app name
- **Go module**: `go.mod` declares `module github.com/RAiWorks/RGo`
- **All imports**: Every `.go` file imports from `github.com/RAiWorks/RGo/...`

### What Works Well

- The framework is fully built (41 features shipped)
- Clear project structure and documentation
- Consistent naming throughout codebase

### What Needs Improvement

- "RGo" is generic and doesn't communicate speed or purpose
- No domain presence — no `.dev`, `.org`, or `.io` domains secured
- The name doesn't differentiate from other Go frameworks
- "RGo" could be confused with R language + Go combination

## Proposed Approach

### Name Candidates Evaluated

| # | Name | Full Form | Domains Available | Recommendation |
|---|---|---|---|---|
| 1 | **RAMISE** | RApid MIcroSErvices | .dev .org .io | ❌ Pronunciation ambiguous, forced acronym, no Go association |
| 2 | **RapidMS** | Rapid MicroServices | .dev .org .io | ❌ "MS" = Microsoft confusion, sounds like protocol/format, not brandable |
| 3 | **GoRapid** | Go + Rapid | .org only | ❌ Only 1 TLD available — no domain strategy |
| 4 | **RapidGo** | Rapid + Go | .dev .org .io | ✅ **RECOMMENDED** — see rationale below |

### Recommendation: RapidGo

**Rationale**:

1. **Domain strategy**: All three premium TLDs available (.dev for primary site, .org and .io for protective registration)
2. **Instant recognition**: "Rapid" communicates speed; "Go" communicates the language — both key selling points
3. **Follows proven naming patterns**: FastAPI named speed-first and became the #1 Python API framework. "RapidGo" follows the same pattern
4. **Pronunciation**: Unambiguous — everyone will say "rapid-go" correctly
5. **Searchability**: Unique enough to rank well in search, descriptive enough to be found by intent
6. **CLI ergonomics**: `rapidgo serve`, `rapidgo migrate`, `rapidgo make:controller` — reads naturally
7. **Module path**: `github.com/RAiWorks/RapidGo` — clean and professional
8. **Tag line potential**: "RapidGo — Build Fast, Ship Faster" or "RapidGo — The Rapid Go Framework"

### Rejected Alternatives — Detailed Reasoning

**RAMISE (RApid MIcroSErvices)**
- The acronym is a stretch — "RA" from Rapid, "MI" from Micro, "SE" from Services — requires explanation
- Pronunciation is ambiguous: ra-MEEZ? RA-mise? ra-MIZE? This kills word-of-mouth growth
- No "Go" in the name — a Go framework should identify its language
- Could be phonetically confused with "premise" or "demise" (negative association)
- However: most unique and brandable of all options

**RapidMS**
- "MS" is universally associated with Microsoft (MS Office, MS Teams, MS-DOS)
- Sounds like a technical abbreviation (like RPC, HTTP) rather than a product name
- Verbally awkward: "rapid-em-ess" — four syllables for two letters
- Not memorable — hard to build a brand around initials

**GoRapid**
- Only `.org` available — cannot secure `.dev` or `.io`
- A framework needs at least `.dev` (developer-focused TLD) for credibility
- The "Go" prefix pattern (GoLand, GoReleaser) is heavily saturated
- Otherwise a strong name — killed by domain unavailability

## Edge Cases & Risks

- [ ] **GitHub repository rename**: Need to rename `github.com/RAiWorks/RGo` → `github.com/RAiWorks/RapidGo`. GitHub provides automatic redirects from the old URL, but all Go module consumers would need to update their imports
- [ ] **Go module path change**: Changing `go.mod` module path is a breaking change. Since the framework is pre-1.0 and has no external consumers yet, this is acceptable now but would be disruptive later
- [ ] **Import path find-and-replace scope**: Every `.go` file in the project imports from `github.com/RAiWorks/RGo/...`. A single missed replacement breaks compilation
- [ ] **Documentation references**: "RGo" appears in 100+ markdown files. All must be updated
- [ ] **Binary name change**: `rgo` → `rapidgo` changes the CLI UX. Makefile, Dockerfile, docker-compose.yml must all update
- [ ] **Welcome page**: The home template references "RGo" in title, heading, and branding
- [ ] **Git history**: The rename doesn't rewrite history — old commits will reference "RGo". This is acceptable and expected
- [ ] **Search/replace false positives**: Must avoid replacing "RGo" inside words or unrelated contexts (though "RGo" is unique enough that false positives are unlikely)
- [ ] **Case sensitivity**: Need to handle "RGo", "rgo", "Rgo" variants in different contexts (docs, code, filenames)

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| All 41 features complete | Feature | ✅ Done |
| No active feature branches | Process | ✅ Clean — all merged |
| GitHub repository rename capability | Infrastructure | ✅ Available (GitHub Settings → Rename) |
| Domain registration | Infrastructure | 🔴 Not started — domains need to be purchased |

## Open Questions

- [x] **Q1**: Should we rename the GitHub repository (`RAiWorks/RGo` → `RAiWorks/RapidGo`) as part of this feature, or defer it?
  - **Answer**: YES — rename as part of this feature. GitHub provides automatic redirects from old URL. No external consumers yet, so this is the ideal time.
- [x] **Q2**: Should the Go module path change from `github.com/RAiWorks/RGo` to `github.com/RAiWorks/RapidGo`, or could we use a vanity import path like `rapidgo.dev/framework`?
  - **Answer**: Change to `github.com/RAiWorks/RapidGo`. Vanity import path can be added later when `rapidgo.dev` is set up, but the module path should match the actual repository.
- [x] **Q3**: Should the CLI binary be named `rapidgo` (full name) or shortened to something like `rpg` or `rg`?
  - **Answer**: `rapidgo` (full name). Clear, unambiguous, searchable. Follows patterns like `docker`, `kubectl`, `terraform` — full names are standard for CLI tools.
- [x] **Q4**: Should we bump the version from `0.1.0` to `0.2.0` (or `1.0.0`) to mark the rename as a milestone?
  - **Answer**: Bump to `0.2.0`. The rename is a breaking change (module path), but 1.0.0 should be reserved for production-ready status after service mode (#56) and other architectural features ship.
- [x] **Q5**: Which domains should be registered first? Recommendation: `rapidgo.dev` (primary), `rapidgo.io` (redirect), `rapidgo.org` (protective)
  - **Answer**: Register all three. `rapidgo.dev` = framework docs/primary. `rapidgo.io` = commercial/premium services. `rapidgo.org` = redirect to `.dev`.
- [x] **Q6**: Should the organization name `RAiWorks` remain, or also be reconsidered?
  - **Answer**: Keep `RAiWorks`. The org name is the company/author identity — it transcends any single product. `RAiWorks/RapidGo` reads well.
- [x] **Q7**: Final name confirmation — is **RapidGo** the approved choice?
  - **Answer**: YES — **RapidGo** is confirmed.

## Decisions Made

| Date | Decision | Rationale |
|---|---|---|
| 2026-03-07 | Evaluated 4 name candidates | User provided options with domain availability data |
| 2026-03-07 | Recommended **RapidGo** | Best balance of domain availability, brandability, Go association, and pronunciation clarity |
| 2026-03-07 | **RapidGo** confirmed as final name | User approved |
| 2026-03-07 | Use standard **MIT License** with RAi Works copyright | Custom licenses block OSI approval, corporate adoption, community trust, and tooling compatibility. MIT + copyright line achieves brand presence without friction |
| 2026-03-07 | **Domain strategy**: `.dev` = primary, `.io` = commercial, `.org` = redirect | `.dev` follows Go ecosystem pattern. `.io` for premium services. `.org` as protective registration |
| 2026-03-07 | **Business model**: MIT open-source + premium custom development | Framework is free (MIT). Revenue from consulting, custom dev, enterprise support via `rapidgo.io` |
| 2026-03-07 | CLI binary name: `rapidgo` (full name) | Unambiguous, searchable, follows standard CLI naming patterns |
| 2026-03-07 | Version bump to `0.2.0` | Breaking change (module path) warrants minor bump. `1.0.0` reserved for production-ready |
| 2026-03-07 | Keep `RAiWorks` org name | Company identity transcends products. `RAiWorks/RapidGo` reads well |
| 2026-03-07 | Rename GitHub repo as part of this feature | No external consumers yet. GitHub provides automatic redirects |
| 2026-03-07 | `RGO_MODE` env var → stays as `RAPIDGO_MODE` after rename | Consistency with new name. No backward compat needed (feature not implemented yet) |

## Scope of Change (Preliminary Audit)

This rename is a **project-wide refactoring** that touches every layer:

### Code Changes

| Area | Files Affected | Change |
|---|---|---|
| Go module path | `go.mod` | `module github.com/RAiWorks/RGo` → `module github.com/RAiWorks/RapidGo` |
| All Go imports | Every `.go` file (~80+ files) | `github.com/RAiWorks/RGo/...` → `github.com/RAiWorks/RapidGo/...` |
| CLI root command | `core/cli/root.go` | `Use: "rgo"` → `Use: "rapidgo"` |
| CLI version | `core/cli/root.go` | Update framework name in version output |
| Binary output | `Makefile` | `bin/rgo` → `bin/rapidgo` |
| Docker binary | `Dockerfile` | Binary name reference |
| Docker Compose | `docker-compose.yml` | Service/container name |
| Welcome page | `resources/views/home.html` | Title, heading, branding text |
| Home controller | `http/controllers/home_controller.go` | Any hardcoded "RGo" strings |

### Documentation Changes

| Area | Files Affected | Change |
|---|---|---|
| Project context | `docs/project-context.md` | Name, repository URL, all references |
| Project roadmap | `docs/project-roadmap.md` | Header references |
| Mastery doc | `docs/mastery.md` | If any framework-specific references |
| Framework docs | `docs/framework/**/*.md` | All RGo references |
| Feature docs | `docs/features/**/*.md` (100+ files) | All RGo references |
| README | `README.md` | Full rebrand |
| Architecture doc | `docs/framework/service-mode-architecture.md` | All RGo references |

### Infrastructure Changes

| Area | Change |
|---|---|
| GitHub repo | Rename `RAiWorks/RGo` → `RAiWorks/RapidGo` |
| Domains | Register `rapidgo.dev`, `rapidgo.io`, `rapidgo.org` |
| `.gitignore` | Update binary name if referenced |

---

## Discussion Complete ✅

**Summary**: Rename the framework from "RGo" to "RapidGo" across the entire codebase (~110 files, ~400 references). Use standard MIT license with RAi Works copyright. Register `rapidgo.dev` (primary), `rapidgo.io` (commercial), `rapidgo.org` (redirect). Bump version to `0.2.0`. CLI binary becomes `rapidgo`. Module path becomes `github.com/RAiWorks/RapidGo`. Add LICENSE file. Rename GitHub repository.
**Completed**: 2026-03-07
**Next**: Create architecture doc → `55-framework-rename-architecture.md`
