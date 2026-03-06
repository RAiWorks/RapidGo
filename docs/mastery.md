# 🏛️ MASTERY — Development Process Framework

> **A disciplined, structured approach to building software — from idea to production.**
> Every feature starts as a discussion, gets designed, becomes a plan, and ships as clean, tested code.
> No cowboy coding. No improvised architectures. No untested merges.

---

## 📋 Table of Contents

- [Philosophy](#-philosophy)
- [Document Ecosystem](#-document-ecosystem)
- [The Workflow — Feature Lifecycle](#-the-workflow--feature-lifecycle)
- [Document Naming Convention](#-document-naming-convention)
- [Definition of Done](#-definition-of-done)
- [Document Templates](#-document-templates)
  - [Discussion Doc](#1-discussion-document)
  - [Architecture Doc](#2-architecture-document)
  - [Tasks Doc](#3-tasks-document)
  - [Test Plan Doc](#4-test-plan-document)
  - [API Spec Doc](#5-api-spec-document)
  - [Changelog Doc](#6-changelog-document)
  - [Review Doc](#7-review-document)
- [Milestone Audit — Every 5 Features](#-milestone-audit--every-5-features)
- [Git Branching Strategy](#-git-branching-strategy)
- [Commit Message Convention](#-commit-message-convention)
- [Go Development Standards](#-go-development-standards)
- [Quick Reference](#-quick-reference)

---

## 💡 Philosophy

### Core Principles

1. **Think before you type.** No code is written until the discussion is complete.
2. **Design before you build.** Architecture decisions are documented, not improvised.
3. **Plan before you execute.** Every task is written down and checkable.
4. **Verify before you present.** Cross-check docs against blueprint, codebase, and prior features before review.
5. **Review before you build.** All docs are reviewed and approved before implementation starts.
6. **Test before you ship.** Every feature has a test plan with clear acceptance criteria.
7. **Verify before you ship.** Cross-check implementation against architecture and prior features before merging.
8. **Document as you go.** Changes are logged in real time, not reconstructed from memory.
9. **Review when you're done.** Reflect, learn, carry lessons forward.

### Why This Framework Exists

Most projects fail not because of bad code, but because of:

- Features built without understanding the full picture
- Architecture decisions made on the fly and forgotten
- Tasks that live in someone's head instead of a checklist
- Bugs shipped because no one defined what "working" means
- The same mistakes repeated because no one wrote them down
- Scope creep from features that were never properly scoped

**Mastery** solves this by making the process as important as the product. The framework is tool-agnostic, language-agnostic, and scales from solo development to full teams.

### The Golden Rule

> **If it's not written down, it doesn't exist.**

Verbal decisions evaporate. Slack messages get buried. Only docs persist.

---

## 📂 Document Ecosystem

Every project using Mastery has this documentation structure:

```
docs/
├── mastery.md                  # 🏛️ THIS — The process framework
├── project-context.md          # 🎯 Project identity, stack, architecture, scope
├── project-roadmap.md          # 🗺️ Feature list, priorities, dependencies, progress
│
├── features/                   # 📁 Per-feature working docs (one set per feature)
│   ├── XX-feature-discussion.md
│   ├── XX-feature-architecture.md
│   ├── XX-feature-tasks.md
│   ├── XX-feature-testplan.md
│   ├── XX-feature-api.md          # (only for features with API endpoints)
│   ├── XX-feature-changelog.md
│   └── XX-feature-review.md
│
└── framework/                  # 📁 Framework reference documentation
    ├── README.md               # Navigation hub for all framework docs
    ├── architecture/           # Architecture overview, diagrams, design principles
    ├── core/                   # Service container, providers, config, logging
    ├── http/                   # Routing, controllers, middleware, views
    ├── data/                   # Database, models, migrations, pagination
    ├── security/               # Auth, sessions, CSRF, CORS, rate limiting
    ├── infrastructure/         # Caching, mail, events, storage, i18n
    ├── cli/                    # CLI commands, code generation
    ├── guides/                 # Getting started, tutorials, walkthroughs
    ├── testing/                # Test strategy, unit tests, integration tests
    ├── deployment/             # Docker, health checks, build & run
    ├── reference/              # Env vars, helpers, middleware quick reference
    └── appendix/               # Glossary, roadmap, naming
```

### Document Roles

| Document | Scope | Purpose | When Created |
|---|---|---|---|
| **mastery.md** | Universal | Process framework — HOW you work | Once (project init) |
| **project-context.md** | Project | Project identity — WHAT you're building | Once (project init) |
| **project-roadmap.md** | Project | Feature plan — WHEN you build it | Once, updated continuously |
| **discussion.md** | Feature | Requirements & design conversation | Start of every feature |
| **architecture.md** | Feature | Technical design & file structure | After discussion, before coding |
| **tasks.md** | Feature | Phased implementation checklist | After architecture is designed |
| **testplan.md** | Feature | Test cases & acceptance criteria | Alongside or after tasks |
| **api.md** | Feature | API contracts (routes, payloads, status codes) | When feature has API endpoints |
| **changelog.md** | Feature | Running log of changes during implementation | During build phase |
| **review.md** | Feature | Post-implementation retrospective | After merge to main |

### Which Docs Are Required vs Optional?

| Document | Required? | Skip When... |
|---|---|---|
| **discussion** | ✅ Always | Never skip — this is the foundation |
| **architecture** | ✅ Always | Never skip — even simple features need file structure planning |
| **tasks** | ✅ Always | Never skip — this is your execution plan |
| **testplan** | ✅ Always | Never skip — define "done" before you start |
| **api** | ⚡ Conditional | Feature has no HTTP/API endpoints |
| **changelog** | ✅ Always | Never skip — tracks what actually happened vs what was planned |
| **review** | ✅ Always | Never skip — learning compounds over time |

---

## 🔄 The Workflow — Feature Lifecycle

Every feature flows through **6 stages plus two Cross-Checks and a mandatory Review Gate**. Each stage has a clear entry condition and exit condition. No stage may be skipped. Cross-Checks verify completeness and catch gaps early. The Review Gate separates documentation from implementation — no code is written until docs are reviewed and approved.

```
┌────────────────────────────────────────────────────────────────────────────────────────┐
│                               FEATURE LIFECYCLE                                        │
│                                                                                        │
│  ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌───────────┐   ┌──────────┐            │
│  │    1.    │   │    2.    │   │    3.    │   │  CROSS-   │   │  REVIEW  │            │
│  │ DISCUSS  │──▶│ DESIGN   │──▶│  PLAN    │──▶│  CHECK    │──▶│   GATE   │            │
│  │          │   │          │   │          │   │  (DOCS)   │   │   🚦     │            │
│  └──────────┘   └──────────┘   └──────────┘   │   🔍      │   └────┬─────┘            │
│       │              │              │          └───────────┘        │                   │
│   discussion    architecture      tasks       Verify docs vs       ⏸️ STOP              │
│   doc created   doc created     doc created   blueprint, scope,    Present verified     │
│                                 testplan      consistency, cross-  docs for user review. │
│                                 doc created   feature impact.      Wait for approval     │
│                                 changelog     Fix gaps before      before proceeding.    │
│                                 doc created   review gate.                               │
│                                 api doc            │                    │                │
│                                 (if needed)        │               User says             │
│                                                    │               "continue"            │
│                                                    │                    │                │
│                                                    ▼                    ▼                │
│                ┌──────────┐   ┌───────────┐   ┌──────────┐   ┌─────────┐               │
│                │    4.    │   │  CROSS-   │   │   5.    │   │   6.    │               │
│                │  BUILD   │──▶│  CHECK    │──▶│  SHIP   │──▶│ REFLECT │               │
│                │          │   │  (IMPL)   │   │         │   │         │               │
│                └──────────┘   │   🔍      │   └─────────┘   └─────────┘               │
│                     │         └───────────┘        │              │                     │
│                 changelog    Verify code vs     review doc     review doc                │
│                 updated      architecture,      created        completed                 │
│                              cross-feature                     roadmap                   │
│                              impact, tests.                    updated                   │
│                              Fix gaps before                                             │
│                              shipping.                                                   │
└────────────────────────────────────────────────────────────────────────────────────────┘
```

### Stage 1 — Discuss 💬

> **Entry**: Feature identified in roadmap
> **Exit**: Discussion doc marked COMPLETE with summary

**Purpose**: Fully understand the feature before any design work begins. Surface ambiguity, edge cases, and dependencies early — when changes are free, not expensive.

| Action | Detail |
|---|---|
| Create `XX-feature-discussion.md` | Use the discussion template |
| Define WHAT the feature does | Functional requirements, user stories, acceptance criteria |
| Understand current state | Reference existing code, prior art, or "greenfield" |
| Identify the approach | High-level "how" — not detailed architecture yet |
| Surface edge cases | What can go wrong? What's tricky? What's ambiguous? |
| List dependencies | What must exist before this feature can be built? |
| Resolve open questions | Discuss iteratively until all questions are answered |
| Mark COMPLETE | Add summary at top, note the date |

**Anti-patterns to avoid**:
- Rushing to design before understanding the problem
- Leaving open questions unresolved
- Skipping dependency analysis

### Stage 2 — Design 🏗️

> **Entry**: Discussion marked COMPLETE
> **Exit**: Architecture doc reviewed and finalized

**Purpose**: Translate understanding into a technical blueprint. Every file, interface, data model, and data flow is defined before any code is written.

| Action | Detail |
|---|---|
| Create `XX-feature-architecture.md` | Use the architecture template |
| Define file structure | Every file to create/modify, with full paths |
| Design data models | Schema, relationships, constraints, migrations |
| Design interfaces | Function signatures, struct definitions, interface contracts |
| Draw data flow | How data moves through the system (request → response) |
| Document trade-offs | Why this approach over alternatives, with pros/cons |
| Define config changes | Environment variables, settings, feature flags |
| Identify security surface | Auth, validation, encryption — what applies? |

**Anti-patterns to avoid**:
- Designing in code instead of in documentation
- Skipping trade-off analysis
- Underspecifying interfaces

### Stage 3 — Plan 📋

> **Entry**: Architecture doc finalized
> **Exit**: Tasks doc + test plan + API spec (if needed) created, docs committed

**Purpose**: Break the architecture into granular, checkable tasks organized by phase. No task should be ambiguous — if you can't check it off definitively, break it down further.

| Action | Detail |
|---|---|
| Create `XX-feature-tasks.md` | Break architecture into atomic, checkable tasks |
| Organize into phases | Group by layer: data → logic → HTTP → UI → test → docs |
| Add checkpoints | Verification points between phases |
| Create `XX-feature-testplan.md` | Define test cases and acceptance criteria |
| Create `XX-feature-api.md` | If feature has endpoints — define full contracts |
| Create `XX-feature-changelog.md` | Empty — ready for build phase logging |
| Create docs branch | `docs/XX-feature-name` from `main` |
| Commit all docs | Commit the complete doc set to the docs branch |
| Push docs branch | Push `docs/XX-feature-name` to remote |
| Merge to main | Merge docs branch into `main`, push `main` |
| Keep docs branch | Never delete — docs branches are historical records |

**Anti-patterns to avoid**:
- Vague tasks like "implement feature" — be specific
- Missing checkpoints between phases
- Not defining the test plan before building
- Committing docs directly to `main` without a branch

---

### � Cross-Check (Docs) — MANDATORY VERIFICATION

> **Entry**: All planning docs created, committed, pushed, and merged to `main`
> **Exit**: All gaps fixed, docs updated, ready for Review Gate

**Purpose**: Systematically verify the documentation set against the blueprint, the existing codebase, and all prior feature docs. Catch gaps, scope creep, and inconsistencies BEFORE presenting docs to the user for review. This ensures the user always receives clean, verified docs.

| # | Check | Detail |
|---|---|---|
| 1 | **Blueprint coverage** | Compare docs against the relevant blueprint section line by line. Every element the blueprint shows for this feature must be accounted for — either implemented or explicitly listed as deferred. |
| 2 | **Scope check** | Identify anything in the docs that goes BEYOND the blueprint section. If present, verify it's a justified adaptation (e.g., testability, using our config system) — not scope creep. |
| 3 | **Doc-to-doc consistency** | Discussion → Architecture → Tasks → Testplan must align. File counts, file lists, task counts, test case counts must be consistent. Every functional requirement should trace to architecture code, tasks, and tests. |
| 4 | **Existing codebase check** | Verify docs correctly reference existing code, APIs, and file paths. Confirm file stubs and directories mentioned actually exist. Check import paths are correct. |
| 5 | **Cross-feature impact** | Check if this feature modifies files owned by prior features. If so, verify those prior feature docs are either historically accurate (describing state at their time) or need updating. Update provider ordering comments, shared file references, etc. |
| 6 | **`.env` / config alignment** | Verify every env var referenced in code/docs exists in `.env` (or is added). Variable names must match between `.env`, `NewXConfig()`, and architecture docs. |
| 7 | **Architecture code review** | Verify all code blocks in the architecture doc compile conceptually — correct imports, correct function signatures, correct usage of existing framework APIs. |

**Actions**:
1. Run through the checklist above
2. Document findings with severity (gap / scope creep / minor)
3. Fix all gaps — update the relevant docs
4. Commit and push fixes to `main`
5. Present the cross-check verdict to the user

**Anti-patterns to avoid**:
- Skipping the cross-check because "the docs look fine"
- Fixing code/gaps without updating the related docs
- Presenting unverified docs at the Review Gate

---

### �🚦 Review Gate — MANDATORY STOP

> **Entry**: All planning docs (discussion, architecture, tasks, testplan, changelog, api) created and committed
> **Exit**: User has reviewed docs and explicitly says "continue"

**Purpose**: Ensure the human reviews and approves all documentation before any code is written. This prevents wasted implementation effort if the plan has gaps, scope issues, or misunderstandings. No code is written until the user gives the green light.

| Action | Detail |
|---|---|
| Present doc summary | List all created docs with key highlights |
| Highlight decisions | Surface important trade-offs and approach choices |
| Wait for user review | **STOP — do not proceed automatically** |
| User says "continue" | Only then move to Stage 4 — Build |

**What happens at the gate**:
1. All 5-6 docs are created on `docs/XX-feature-name` branch, pushed, and merged to `main`
2. A summary is presented: scope, approach, file structure, task count, test count
3. **Execution pauses** — the user reviews the docs at their own pace
4. The user may request changes to docs before approving
5. When the user says "continue" (or equivalent), implementation begins on `feature/XX-feature-name` branch

**Anti-patterns to avoid**:
- Skipping the gate and jumping straight to implementation
- Creating docs and building in the same step
- Treating the gate as optional — it is MANDATORY for every feature

### Stage 4 — Build 🔨

> **Entry**: All planning docs approved by user at Review Gate, feature branch created
> **Exit**: All task checkboxes checked, all tests pass

**Purpose**: Execute the plan methodically. Check off tasks as you go. Log everything that deviates from the plan. Commit frequently with clear messages. This stage only begins after the user has reviewed all docs and explicitly approved them.

| Action | Detail |
|---|---|
| Create feature branch | `feature/XX-feature-name` from `main` |
| Execute tasks phase by phase | Check off items as you complete them |
| Log changes in changelog | What was built, what deviated, decisions made |
| Commit frequently | Clear messages following the commit convention |
| Run tests at checkpoints | Verify each phase before moving to the next |
| Push to feature branch | Keep remote in sync — don't accumulate local commits |

**Anti-patterns to avoid**:
- Working without checking off tasks
- Forgetting to log deviations from the plan
- Large, infrequent commits

---

### 🔍 Cross-Check (Implementation) — MANDATORY VERIFICATION

> **Entry**: All tasks complete, all tests pass on feature branch
> **Exit**: All gaps fixed, related docs updated, ready to Ship

**Purpose**: Systematically verify the implementation against the architecture doc, all prior features, and the full test suite. Catch deviations, missing pieces, and cross-feature regressions BEFORE merging to `main`. This ensures only verified, complete code ships.

| # | Check | Detail |
|---|---|---|
| 1 | **Code vs. architecture doc** | Compare every code block in the architecture doc against the actual implementation. Every function, struct, and interface should match. Any deviation must be logged in the changelog with a reason. |
| 2 | **Task completion** | Verify every task in the tasks doc is genuinely complete — not just checked off. Each checkpoint should have been verified. |
| 3 | **Test coverage** | Verify every test case in the testplan has a corresponding test function. All tests pass. Run `go test ./...` for full regression. |
| 4 | **Scope check** | Ensure the implementation doesn't exceed what the architecture doc specifies. No extra files, no extra functions, no unplanned features. |
| 5 | **Cross-feature impact** | If this feature modified shared files (`cmd/main.go`, `.env`, etc.), verify existing features still work. Run full test suite. Check that provider ordering comments are accurate. |
| 6 | **Deviations logged** | Every difference between the architecture doc and the actual code must be recorded in the changelog's "Deviations from Plan" table. Zero deviations is ideal; undocumented deviations are unacceptable. |
| 7 | **Related docs updated** | If cross-feature impact was found, update the affected docs. If `.env` was changed, verify it's documented. If provider order changed, update comments in `main.go`. |

**Actions**:
1. Run through the checklist above
2. Document findings
3. Fix all gaps — update code and/or docs
4. Commit and push fixes
5. Present the cross-check verdict (can be combined with Ship summary)

**Anti-patterns to avoid**:
- Shipping without cross-checking
- Finding deviations but not logging them in the changelog
- Fixing code without updating the corresponding docs

### Stage 5 — Ship 🚀

> **Entry**: All tasks complete, all tests pass
> **Exit**: Feature merged to main, pushed to remote

**Purpose**: Final quality gate. Self-review every change, run the full test plan one last time, then merge with confidence.

| Action | Detail |
|---|---|
| Self-review all changes | Read your own diff — would you approve this PR? |
| Final test pass | Full test plan execution one last time |
| Merge to main | PR or direct merge (based on team workflow) |
| Push main | Trigger CI/CD pipeline if configured |
| Keep the feature branch | Never delete — branches are historical records |

**Anti-patterns to avoid**:
- Merging without self-review
- Skipping the final test pass
- Deleting feature branches

### Stage 6 — Reflect 🪞

> **Entry**: Feature merged to main
> **Exit**: Review doc completed, roadmap updated

**Purpose**: Learning compounds over time. Every shipped feature teaches you something — but only if you write it down.

| Action | Detail |
|---|---|
| Create/complete `XX-feature-review.md` | Use the review template |
| What went well? | Patterns to repeat in future features |
| What went wrong? | Blockers, time sinks, bugs, surprises |
| What was learned? | New knowledge, techniques, insights |
| What to change next time? | Concrete, actionable improvements |
| Update roadmap | Mark feature as complete in `project-roadmap.md` |

**Anti-patterns to avoid**:
- Skipping reflection because "it went fine"
- Writing vague lessons like "be more careful"
- Not updating the roadmap

---

## � Milestone Audit — Every 5 Features

> **Trigger**: After every 5th feature is shipped (#25, #30, #35, #41)
> **Purpose**: Systematic health check of the entire project — catch accumulated gaps, scope creep, missing docs, and stale state before they compound.

Per-feature cross-checks catch issues within a single feature. But some problems only emerge across features: drifting conventions, missing review docs, stale roadmap markers, orphaned code, forgotten deferrals. The Milestone Audit catches these at regular intervals.

### Schedule

| Milestone | After Feature | Covers |
|-----------|---------------|--------|
| Audit 1 | #20 | Features #01–#20 ✅ *(completed)* |
| Audit 2 | #25 | Features #21–#25 (+ verify #01–#20 fixes) |
| Audit 3 | #30 | Features #26–#30 |
| Audit 4 | #35 | Features #31–#35 |
| Audit 5 | #41 | Features #36–#41 (final) |

### Audit Checklist

| # | Check | Detail |
|---|---|---|
| 1 | **Doc completeness** | Every shipped feature has all 6 required docs (discussion, architecture, tasks, testplan, changelog, review). List any missing files. |
| 2 | **Roadmap accuracy** | All shipped features marked ✅ in `project-roadmap.md`. No stale ⬜ or 🟡 markers. |
| 3 | **Blueprint coverage** | For each feature in the batch, compare implementation against the blueprint section. Identify gaps (missing items) and creep (extra items). |
| 4 | **Gap classification** | Classify every gap as: (a) intentionally deferred to a named future feature, or (b) actionable now. Fix category (b) items immediately. |
| 5 | **Scope creep assessment** | For each creep item, decide: justified (keep) or unjustified (remove or document). |
| 6 | **Test health** | Run full `go test ./... -count=1`. All tests pass. Report total count by package. |
| 7 | **Dependency audit** | Review `go.mod` direct dependencies. No unnecessary additions. No outdated versions with known issues. |
| 8 | **Config/security check** | `.env` variables match code references. No secrets in version control. `.gitignore` covers build artifacts. |
| 9 | **Cross-feature consistency** | Provider boot order, import paths, shared file modifications are all consistent across features. |
| 10 | **Memory & process update** | Update mastery memory notes with any new lessons. Update process docs if audit reveals process gaps. |

### Output

Each milestone audit produces:
1. **Audit report** — presented to user with findings table
2. **Fixes committed** — all actionable items resolved and committed
3. **Deferred items documented** — gaps tied to specific future features
4. **Memory updated** — lessons and project state refreshed

### Rules

- **Never skip a milestone audit.** Even if everything seems fine — verify.
- **Fix before proceeding.** All actionable gaps must be resolved before starting the next feature.
- **Keep it lightweight.** The audit should take minutes, not hours. Use subagents for parallel investigation.
- **Audit the batch, not everything.** After Audit 1 (#01–#20), subsequent audits focus on the new batch plus spot-checking prior fixes.

---

## �📛 Document Naming Convention

### Feature Documents

All feature documents live in `docs/features/`:

```
docs/features/
├── 01-project-setup-discussion.md
├── 01-project-setup-architecture.md
├── 01-project-setup-tasks.md
├── 01-project-setup-testplan.md
├── 01-project-setup-changelog.md
├── 01-project-setup-review.md
├── 02-auth-system-discussion.md
├── 02-auth-system-architecture.md
├── 02-auth-system-tasks.md
├── 02-auth-system-testplan.md
├── 02-auth-system-api.md
├── 02-auth-system-changelog.md
├── 02-auth-system-review.md
└── ...
```

### Naming Rules

| Element | Format | Example |
|---|---|---|
| **Sequence** | 2-digit zero-padded | `01`, `02`, `10` |
| **Feature name** | lowercase, hyphen-separated | `auth-system`, `user-dashboard` |
| **Doc type** | suffix before `.md` | `-discussion`, `-architecture`, `-tasks`, `-testplan`, `-api`, `-changelog`, `-review` |
| **Branch name** | `feature/XX-feature-name` | `feature/02-auth-system` |

### Sequence Assignment Principle

Features are numbered in **dependency order** — what must exist first gets a lower number. Define the sequence in `project-roadmap.md` based on the project's dependency graph.

General ordering logic:

1. **Foundation** — project setup, module init, configuration, directory structure
2. **Core infrastructure** — service container, providers, error handling, logging
3. **Data layer** — database, models, migrations, seeders
4. **HTTP layer** — routing, controllers, middleware, request/response
5. **Security** — authentication, sessions, CSRF, CORS, rate limiting
6. **Business logic** — services, events, caching, mail
7. **Presentation** — views, templates, static assets
8. **Developer experience** — CLI, code generation, helpers
9. **Testing** — test infrastructure, integration tests
10. **Deployment** — Docker, CI/CD, health checks, monitoring

---

## ✅ Definition of Done

A feature is **DONE** when ALL of the following are true:

| # | Criterion | Verified By |
|---|---|---|
| 1 | Discussion doc is marked COMPLETE | Summary present, date noted |
| 2 | Architecture doc is FINALIZED | All sections filled, trade-offs documented |
| 3 | Docs cross-check passed | Blueprint coverage, scope, consistency verified |
| 4 | All tasks in tasks doc are checked off | Every `[ ]` is `[x]` |
| 5 | All test plan test cases pass | Test summary table filled |
| 6 | Implementation cross-check passed | Code vs. architecture, cross-feature, regression verified |
| 7 | No known bugs remain | Or documented as accepted/deferred |
| 8 | Changelog reflects actual implementation | Deviations logged with reasons |
| 9 | Code is self-reviewed | Diff read, code is clean |
| 10 | Feature branch merged to main | Fast-forward or merge commit |
| 11 | Main pushed to remote | CI/CD green (if configured) |
| 12 | Feature branch preserved | Not deleted |
| 13 | Review doc completed | Lessons captured |
| 14 | Roadmap updated | Feature marked complete |

If any criterion is not met, the feature is **NOT DONE** — regardless of whether the code works.

---

## 📝 Document Templates

Below are the templates for every document type. Copy the relevant template when starting a new document.

---

### 1. Discussion Document

**Filename**: `XX-feature-name-discussion.md`
**Purpose**: Understand the feature completely through structured conversation before any design or code.

````markdown
# 💬 Discussion: [Feature Name]

> **Feature**: `XX` — [Feature Name]
> **Status**: 🟡 IN PROGRESS | 🟢 COMPLETE
> **Branch**: `feature/XX-feature-name`
> **Depends On**: #XX, #XX (list prerequisite feature numbers)
> **Date Started**: YYYY-MM-DD
> **Date Completed**: —

---

## Summary

<!-- One paragraph: What does this feature do and why does it matter? -->

---

## Functional Requirements

<!-- What should this feature do from the user's / developer's perspective? -->

- As a [role], I want [action] so that [outcome]
- ...

## Current State / Reference

<!-- How does this work today? Existing code? Starting from scratch? -->

### What Exists
<!-- Describe current implementation or "Nothing — greenfield feature" -->

### What Works Well
<!-- Patterns to keep or replicate -->

### What Needs Improvement
<!-- What should be redesigned, removed, or rethought? -->

## Proposed Approach

<!-- High-level description of how we'll implement this -->
<!-- NOT detailed architecture — that comes in the architecture doc -->

## Edge Cases & Risks

<!-- What can go wrong? What's non-obvious? -->

- [ ] [Edge case or risk 1]
- [ ] [Edge case or risk 2]

## Dependencies

<!-- What must exist before this feature can be built? -->

| Dependency | Type | Status |
|---|---|---|
| Feature #XX — [Name] | Feature | ✅ Done / 🔴 Not started |
| [Package/Library] | External | ✅ Available / 🔴 Needs install |
| [Service/API] | Infrastructure | ✅ Ready / 🔴 Needs setup |

## Open Questions

<!-- Things we're unsure about — ALL must be resolved before marking COMPLETE -->

- [ ] Question 1?
- [ ] Question 2?

## Decisions Made

<!-- Running log — add entries as decisions happen -->

| Date | Decision | Rationale |
|---|---|---|
| YYYY-MM-DD | [Decision] | [Why] |

## Discussion Complete ✅

<!-- Fill this section when ALL open questions are resolved -->

**Summary**: [One-sentence final summary of what was agreed]
**Completed**: YYYY-MM-DD
**Next**: Create architecture doc → `XX-feature-name-architecture.md`
````

---

### 2. Architecture Document

**Filename**: `XX-feature-name-architecture.md`
**Purpose**: Technical design — file structure, data models, interfaces, data flow, and trade-offs.

````markdown
# 🏗️ Architecture: [Feature Name]

> **Feature**: `XX` — [Feature Name]
> **Discussion**: [`XX-feature-name-discussion.md`](XX-feature-name-discussion.md)
> **Status**: 🟡 DRAFT | 🟢 FINALIZED
> **Date**: YYYY-MM-DD

---

## Overview

<!-- One paragraph: Technical summary of the approach -->

## File Structure

<!-- Every file to create or modify, with full paths from project root -->

```
path/to/
├── new-file-1.go           # Purpose
├── new-file-2.go           # Purpose
├── new-file_test.go        # Tests for new-file
└── existing-file.go        # MODIFY — what changes
```

## Data Model

<!-- Database tables, schemas, relationships -->

### [Table/Entity Name]

| Field | Type | Constraints | Description |
|---|---|---|---|
| id | uint | PK, auto | Unique identifier |
| ... | ... | ... | ... |

### Relationships
<!-- Foreign keys, references, associations -->

## Component Design

<!-- Structs, interfaces, functions — the building blocks -->

### [Component Name]

**Responsibility**: [What this component does]
**Package**: `path/to/package`
**File**: `path/to/file.go`

```
Exported API:
├── NewComponent(deps) → *Component             # Constructor
├── (c *Component) Method(params) → (result, error)  # Description
└── (c *Component) Method(params) → error            # Description
```

### Interfaces

```go
type InterfaceName interface {
    Method(params) (returnType, error)
}
```

## Data Flow

<!-- How data moves through the system for this feature -->

```
[Trigger] → [Entry Point] → [Processing] → [Storage] → [Response]
```

<!-- Describe each step -->

## Configuration

<!-- Environment variables, config files -->

| Key | Type | Default | Description |
|---|---|---|---|
| `ENV_VAR_NAME` | string | `""` | What it configures |

## Security Considerations

<!-- Auth, authorization, input validation, encryption, rate limiting -->

## Trade-offs & Alternatives

<!-- Why this approach? What else was considered? -->

| Approach | Pros | Cons | Verdict |
|---|---|---|---|
| Chosen approach | ... | ... | ✅ Selected |
| Alternative A | ... | ... | ❌ Reason |

## Next

Create tasks doc → `XX-feature-name-tasks.md`
````

---

### 3. Tasks Document

**Filename**: `XX-feature-name-tasks.md`
**Purpose**: Phased implementation checklist with checkpoints between phases.

````markdown
# ✅ Tasks: [Feature Name]

> **Feature**: `XX` — [Feature Name]
> **Architecture**: [`XX-feature-name-architecture.md`](XX-feature-name-architecture.md)
> **Branch**: `feature/XX-feature-name`
> **Status**: 🔴 NOT STARTED | 🟡 IN PROGRESS | 🟢 COMPLETE
> **Progress**: 0/XX tasks complete

---

## Pre-Flight Checklist

- [ ] Discussion doc is marked COMPLETE
- [ ] Architecture doc is FINALIZED
- [ ] Feature branch created from latest `main`
- [ ] Dependent features are merged to `main`
- [ ] Test plan doc created
- [ ] Changelog doc created (empty)

---

## Phase A — Data Layer

> Database schema, migrations, models, seeds.

- [ ] **A.1** — [Specific task description]
  - [ ] Sub-step if needed
  - [ ] Sub-step if needed
- [ ] **A.2** — [Specific task description]
- [ ] 📍 **Checkpoint A** — Migrations run, models instantiate, seed data loads

---

## Phase B — Core Logic

> Business logic, services, helpers, internal packages.

- [ ] **B.1** — [Specific task description]
- [ ] **B.2** — [Specific task description]
- [ ] 📍 **Checkpoint B** — Core logic works independently (unit tests pass)

---

## Phase C — HTTP / API Layer

> Routes, controllers, middleware, request/response handling.

- [ ] **C.1** — [Specific task description]
- [ ] **C.2** — [Specific task description]
- [ ] 📍 **Checkpoint C** — All endpoints respond correctly, middleware applied

---

## Phase D — Presentation

> Views, templates, components, static assets, client-side logic.

- [ ] **D.1** — [Specific task description]
- [ ] **D.2** — [Specific task description]
- [ ] 📍 **Checkpoint D** — Visual review complete, responsive, accessible

---

## Phase E — Testing

> Execute the test plan, verify all acceptance criteria.

- [ ] **E.1** — Run test plan: happy path test cases
- [ ] **E.2** — Run test plan: error cases
- [ ] **E.3** — Run test plan: edge cases
- [ ] **E.4** — Run test plan: security tests
- [ ] 📍 **Checkpoint E** — All acceptance criteria met, test summary filled

---

## Phase F — Documentation & Cleanup

> Code comments, doc updates, self-review.

- [ ] **F.1** — Add inline comments where logic is non-obvious
- [ ] **F.2** — Update changelog doc with final summary
- [ ] **F.3** — Update project roadmap progress
- [ ] **F.4** — Self-review all diffs
- [ ] 📍 **Checkpoint F** — Clean code, complete docs, ready to ship

---

## Ship 🚀

- [ ] All phases complete
- [ ] All checkpoints verified
- [ ] Final commit with descriptive message
- [ ] Push to feature branch
- [ ] Merge to `main`
- [ ] Push `main`
- [ ] **Keep the feature branch** — do not delete
- [ ] Create review doc → `XX-feature-name-review.md`
````

---

### 4. Test Plan Document

**Filename**: `XX-feature-name-testplan.md`
**Purpose**: Define exactly what "working" means — test cases, acceptance criteria, edge cases.

````markdown
# 🧪 Test Plan: [Feature Name]

> **Feature**: `XX` — [Feature Name]
> **Tasks**: [`XX-feature-name-tasks.md`](XX-feature-name-tasks.md)
> **Date**: YYYY-MM-DD

---

## Acceptance Criteria

<!-- The feature is DONE when ALL of these are true -->

- [ ] [Criterion 1 — specific, measurable, verifiable]
- [ ] [Criterion 2]
- [ ] [Criterion 3]

---

## Test Cases

### TC-01: [Test Case Name]

| Property | Value |
|---|---|
| **Category** | Happy Path / Error / Edge Case / Security / Performance |
| **Precondition** | [What must be true before this test] |
| **Steps** | 1. [Step] → 2. [Step] → 3. [Step] |
| **Expected Result** | [What should happen] |
| **Status** | ⬜ Not Run / ✅ Pass / ❌ Fail |
| **Notes** | — |

### TC-02: [Test Case Name]

| Property | Value |
|---|---|
| **Category** | ... |
| **Precondition** | ... |
| **Steps** | ... |
| **Expected Result** | ... |
| **Status** | ⬜ Not Run |
| **Notes** | — |

<!-- Add more test cases as needed -->

---

## Edge Cases

| # | Scenario | Expected Behavior |
|---|---|---|
| 1 | [Edge case description] | [How system should handle it] |
| 2 | ... | ... |

## Security Tests

| # | Test | Expected |
|---|---|---|
| 1 | [Unauthorized access attempt] | [Rejected with proper status code] |
| 2 | [Malicious input / injection] | [Sanitized / rejected] |

## Performance Considerations

| Metric | Target | Actual |
|---|---|---|
| Response time (p95) | < Xms | — |
| Memory usage | < XMB | — |
| Throughput | > X req/s | — |

---

## Test Summary

<!-- Fill AFTER running all tests -->

| Category | Total | Pass | Fail | Skip |
|---|---|---|---|---|
| Happy Path | — | — | — | — |
| Error Cases | — | — | — | — |
| Edge Cases | — | — | — | — |
| Security | — | — | — | — |
| Performance | — | — | — | — |
| **Total** | — | — | — | — |

**Result**: ⬜ NOT RUN | ✅ ALL PASS | ❌ HAS FAILURES
````

---

### 5. API Spec Document

**Filename**: `XX-feature-name-api.md`
**Purpose**: HTTP API contracts — endpoints, payloads, status codes, auth requirements.

> **Only create this doc when the feature includes HTTP/API endpoints.**

````markdown
# 🔌 API Spec: [Feature Name]

> **Feature**: `XX` — [Feature Name]
> **Base Path**: `/api/v1/...`
> **Auth Required**: Yes / No / Mixed
> **Date**: YYYY-MM-DD

---

## Endpoints Overview

| Method | Path | Auth | Description |
|---|---|---|---|
| `GET` | `/resource` | 🔒 Yes | List resources |
| `POST` | `/resource` | 🔒 Yes | Create resource |
| `GET` | `/resource/:id` | 🔒 Yes | Get single resource |
| `PUT` | `/resource/:id` | 🔒 Yes | Update resource |
| `DELETE` | `/resource/:id` | 🔒 Yes | Delete resource |

---

## Endpoint Details

### `GET /resource`

**Description**: [What this endpoint does]
**Auth**: [Required / Optional / None]

**Query Parameters**:

| Param | Type | Required | Default | Description |
|---|---|---|---|---|
| `page` | int | No | 1 | Page number |
| `limit` | int | No | 20 | Items per page |
| `sort` | string | No | `id` | Sort field |
| `order` | string | No | `asc` | Sort direction (`asc` / `desc`) |

**Success Response** (`200`):

```json
{
  "status": "success",
  "data": [],
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 0
  }
}
```

**Error Responses**:

| Status | Body | When |
|---|---|---|
| `401` | `{ "status": "error", "message": "Unauthorized" }` | Missing or invalid auth token |
| `403` | `{ "status": "error", "message": "Forbidden" }` | Insufficient permissions |
| `500` | `{ "status": "error", "message": "Internal error" }` | Server error |

---

### `POST /resource`

**Description**: [What this endpoint does]
**Auth**: Required
**Content-Type**: `application/json`

**Request Body**:

```json
{
  "field1": "string (required)",
  "field2": 0
}
```

**Validation Rules**:

| Field | Rules |
|---|---|
| `field1` | Required, string, min 1, max 255 |
| `field2` | Optional, integer, min 0 |

**Success Response** (`201`):

```json
{
  "status": "success",
  "data": { "id": 1, "field1": "value", "field2": 0 }
}
```

**Error Responses**:

| Status | Body | When |
|---|---|---|
| `400` | `{ "status": "error", "errors": {...} }` | Validation failure |
| `401` | `{ "status": "error", "message": "Unauthorized" }` | No auth |
| `409` | `{ "status": "error", "message": "Already exists" }` | Duplicate resource |
| `422` | `{ "status": "error", "message": "Unprocessable" }` | Semantic error |

<!-- Repeat for each endpoint -->
````

---

### 6. Changelog Document

**Filename**: `XX-feature-name-changelog.md`
**Purpose**: Running log of what actually happened during implementation — changes, deviations, and decisions made during the build phase.

````markdown
# 📝 Changelog: [Feature Name]

> **Feature**: `XX` — [Feature Name]
> **Branch**: `feature/XX-feature-name`
> **Started**: YYYY-MM-DD
> **Completed**: —

---

## Log

<!-- Add entries as you work. Most recent first. -->

### YYYY-MM-DD

- **[Added/Changed/Fixed/Removed]**: [Description of what happened]
  - Detail or context if needed
  - Related file: `path/to/file`

### YYYY-MM-DD

- **[Added]**: [Description]

---

## Deviations from Plan

<!-- Things that went differently than the architecture/tasks docs planned -->

| What Changed | Original Plan | What Actually Happened | Why |
|---|---|---|---|
| [Component] | [Planned approach] | [Actual approach] | [Reason for deviation] |

## Key Decisions Made During Build

<!-- Runtime decisions NOT in the discussion/architecture docs -->

| Decision | Context | Date |
|---|---|---|
| [Decision] | [Why it came up and what was chosen] | YYYY-MM-DD |
````

---

### 7. Review Document

**Filename**: `XX-feature-name-review.md`
**Purpose**: Post-implementation retrospective — capture what happened, what was learned, and what to do differently.

````markdown
# 🪞 Review: [Feature Name]

> **Feature**: `XX` — [Feature Name]
> **Branch**: `feature/XX-feature-name`
> **Merged**: YYYY-MM-DD
> **Duration**: [Start date] → [End date]

---

## Result

**Status**: ✅ Shipped | ⚠️ Shipped with known issues | ❌ Abandoned

**Summary**: [One paragraph — what was built and delivered]

---

## What Went Well ✅

- [Pattern, approach, or decision that worked great]
- [Something to repeat in future features]

## What Went Wrong ❌

- [Problem] — [Impact] — [Resolution]
- [Problem] — [Impact] — [Resolution]

## What Was Learned 📚

- [Concrete lesson or insight gained]
- [New technique or approach discovered]

## What To Do Differently Next Time 🔄

- [Specific, actionable change for future features]
- [Process improvement to apply]

## Metrics

| Metric | Value |
|---|---|
| Tasks planned | XX |
| Tasks completed | XX |
| Tests planned | XX |
| Tests passed | XX |
| Deviations from plan | XX |
| Commits on branch | XX |

## Follow-ups

<!-- Anything spawned from this feature that needs future attention -->

- [ ] [Follow-up item — file as future roadmap entry if significant]
- [ ] [Follow-up item]
````

---

## 🌿 Git Branching Strategy

```
main ──●────●────●─────●────●─────●────●──────▶
        \   ↗    \    ↗     \    ↗     \    ↗
         ●─●      ●─●─●      ●─●        ●─●─●
        docs/    feature/    docs/      feature/
        01-name  01-name     02-name    02-name
        (kept)   (kept)      (kept)     (kept)
```

Every feature produces **two branches**: one for documentation, one for implementation.

### Branch Rules

| Rule | Detail |
|---|---|
| **`main`** | Always deployable. Only receives merges from docs and feature branches. |
| **`docs/XX-name`** | Created from latest `main`. Contains the feature's documentation set. Merged before implementation begins. |
| **`feature/XX-name`** | Created from latest `main` after docs approved at Review Gate. Contains implementation code. |
| **Never delete** | Both docs and feature branches are kept forever as historical records. |
| **One feature at a time** | Finish docs + implementation before starting the next. |

### Branch Naming

| Pattern | Example | When |
|---|---|---|
| `docs/XX-feature-name` | `docs/04-error-handling` | Documentation set (discussion, architecture, tasks, testplan, changelog) |
| `feature/XX-feature-name` | `feature/04-error-handling` | Implementation code (after docs approved) |
| `fix/XX-description` | `fix/03-login-redirect-loop` | Bug fix on a shipped feature |
| `documentation` | `documentation` | Docs-only changes (no feature code) |

### Git Commands — Feature Workflow

```bash
# ─── DOCUMENTATION PHASE ─────────────────────────────────
# Create docs branch
git checkout main
git pull origin main
git checkout -b docs/XX-feature-name

# Write all docs, commit
git add docs/features/XX-*
git commit -m "docs(scope): create Feature #XX documentation set"

# Push docs branch, merge to main
git push origin docs/XX-feature-name
git checkout main
git merge docs/XX-feature-name
git push origin main

# ─── 🚦 REVIEW GATE — wait for user approval ─────────────

# ─── IMPLEMENTATION PHASE ────────────────────────────────
# Create feature branch (after docs approved)
git checkout main
git pull origin main
git checkout -b feature/XX-feature-name

# Work on feature (repeat as needed)
git add .
git commit -m "feat(scope): description"

# Push feature branch, merge to main
git push origin feature/XX-feature-name
git checkout main
git pull origin main
git merge feature/XX-feature-name
git push origin main

# Both branches are KEPT — never delete
```

### Rules of Thumb

- **Commit early, commit often** — small commits tell a better story
- **Push before you stop working** — never leave unpushed work overnight
- **Pull before you merge** — always merge into a fresh `main`
- **Never force push `main`** — history is sacred on the default branch

---

## 📝 Commit Message Convention

### Format

```
type(scope): short description

[optional body — explain WHY, not WHAT]
[optional footer — references, breaking changes]
```

### Types

| Type | When to Use |
|---|---|
| `feat` | New feature or functionality |
| `fix` | Bug fix |
| `docs` | Documentation changes only |
| `style` | Formatting, whitespace — no logic change |
| `refactor` | Code restructure — no behavior change |
| `test` | Adding or updating tests |
| `chore` | Build, config, tooling, dependencies |
| `perf` | Performance improvement |

### Scope

The scope identifies which module or feature is affected. Use short, consistent names.

| Scope | Area |
|---|---|
| `core` | Service container, providers, config |
| `http` | Routing, controllers, middleware |
| `data` | Database, models, migrations |
| `auth` | Authentication, sessions, JWT |
| `security` | CSRF, CORS, rate limiting, crypto |
| `infra` | Caching, mail, events, storage, i18n |
| `cli` | CLI commands, code generation |
| `deploy` | Docker, health checks, build |
| `docs` | Any documentation change |

### Examples

```
feat(http): add route group middleware chaining
fix(auth): resolve session expiry race condition
docs(security): add Cookie session store documentation
refactor(core): extract provider registration into helper
test(data): add pagination edge case coverage
chore(deps): update GORM to v2.x
perf(http): cache compiled route patterns
```

### Commit Message Quality Checklist

- ✅ Imperative mood ("add", "fix", "update" — not "added", "fixed", "updated")
- ✅ Lowercase after the colon
- ✅ No period at the end
- ✅ Under 72 characters for the subject line
- ✅ Body explains WHY, not WHAT (the diff shows what)

---

## 🔧 Go Development Standards

### Project Layout

Follow the standard Go project layout conventions:

```
├── cmd/                    # Application entry points
│   └── app/
│       └── main.go
├── internal/               # Private packages (not importable externally)
│   ├── core/
│   ├── http/
│   └── ...
├── pkg/                    # Public packages (importable by other projects)
├── configs/                # Configuration files
├── migrations/             # Database migrations
├── docs/                   # Documentation (this directory)
├── tests/                  # Integration / E2E tests
├── go.mod
├── go.sum
└── Makefile
```

### Code Quality Gates

Before merging any feature, verify:

| Gate | Command | Must Pass |
|---|---|---|
| **Compile** | `go build ./...` | ✅ Zero errors |
| **Tests** | `go test ./...` | ✅ All pass |
| **Race detector** | `go test -race ./...` | ✅ No races |
| **Vet** | `go vet ./...` | ✅ No issues |
| **Lint** | `golangci-lint run` | ✅ Clean |
| **Format** | `gofmt -l .` | ✅ No output |

### Testing Conventions

| Convention | Detail |
|---|---|
| **Test file location** | Same package, `_test.go` suffix |
| **Test function naming** | `TestXxx`, `TestXxx_SubCase` |
| **Table-driven tests** | Preferred for functions with multiple input scenarios |
| **Test helpers** | Use `t.Helper()` for shared setup functions |
| **Benchmarks** | `BenchmarkXxx` for performance-sensitive code |

### Error Handling

| Principle | Detail |
|---|---|
| **Always handle errors** | Never use `_` for error returns unless justified with a comment |
| **Wrap with context** | `fmt.Errorf("doing X: %w", err)` |
| **Sentinel errors** | Define package-level `var ErrXxx = errors.New("...")` |
| **Don't panic** | Reserve `panic` for truly unrecoverable situations |

---

## ⚡ Quick Reference

### Starting a New Feature — Step by Step

```
 ─── DOCUMENTATION PHASE ───────────────────────────────────────────────
 1.  Check project-roadmap.md → identify next feature number
 2.  Create  git branch: docs/XX-feature-name             (from main)
 3.  Create  docs/features/XX-feature-discussion.md       (discuss)
 4.  Discuss until fully understood → mark COMPLETE
 5.  Create  docs/features/XX-feature-architecture.md     (design)
 6.  Finalize architecture → mark FINALIZED
 7.  Create  docs/features/XX-feature-tasks.md            (plan)
 8.  Create  docs/features/XX-feature-testplan.md         (define done)
 9.  Create  docs/features/XX-feature-api.md              (if has API)
10.  Create  docs/features/XX-feature-changelog.md        (empty, ready)
11.  Commit all docs to docs branch                       (checkpoint)
12.  Push docs branch to remote                           (push)
13.  Merge docs branch to main, push main                 (merge)
14.  Keep docs branch — do not delete                     (preserve)

 ─── � CROSS-CHECK (DOCS) ── MANDATORY VERIFICATION ──────────────────
15.  Verify docs vs blueprint section (line by line)      (completeness)
16.  Check for scope creep — nothing beyond blueprint     (scope)
17.  Verify doc-to-doc consistency (files, tasks, tests)  (consistency)
18.  Verify against existing codebase and .env            (accuracy)
19.  Check cross-feature impact on prior features/docs    (impact)
20.  Fix all gaps, commit and push fixes to main          (fix)

 ─── 🚦 REVIEW GATE ── MANDATORY STOP ─────────────────────────────────
21.  Present cross-check verdict + doc summary to user    (gate)
22.  ⏸️  WAIT for user to review and say "continue"        (gate)

 ─── IMPLEMENTATION PHASE ──────────────────────────────────────────────
23.  Create  git branch: feature/XX-feature-name          (from main)
24.  Execute tasks, log in changelog                      (build)
25.  Run test plan                                        (verify)

 ─── 🔍 CROSS-CHECK (IMPL) ── MANDATORY VERIFICATION ─────────────────
26.  Verify code matches architecture doc                 (accuracy)
27.  Verify all deviations logged in changelog            (traceability)
28.  Verify cross-feature impact, full regression passes  (stability)
29.  Fix all gaps, update code and/or docs                (fix)

 ─── SHIP & REFLECT ────────────────────────────────────────────────────
30.  Push feature branch, merge to main, push main        (ship)
31.  Keep feature branch — do not delete                  (preserve)
32.  Create  docs/features/XX-feature-review.md           (reflect)
33.  Update  project-roadmap.md progress tracker          (track)
```

### Document Quick Reference

| Need to... | Open this doc |
|---|---|
| Understand the process | `docs/mastery.md` (this file) |
| See what we're building | `docs/project-context.md` |
| See what's next | `docs/project-roadmap.md` |
| Start a feature | `docs/features/XX-feature-discussion.md` |
| Design a feature | `docs/features/XX-feature-architecture.md` |
| Plan implementation | `docs/features/XX-feature-tasks.md` |
| Define test cases | `docs/features/XX-feature-testplan.md` |
| Spec an API | `docs/features/XX-feature-api.md` |
| Log build progress | `docs/features/XX-feature-changelog.md` |
| Reflect on delivery | `docs/features/XX-feature-review.md` |
| Browse framework docs | `docs/framework/README.md` |

### Stage Gate Summary

```
DISCUSS  ──▶  Discussion doc marked COMPLETE?        Yes ──▶  DESIGN
DESIGN   ──▶  Architecture doc FINALIZED?            Yes ──▶  PLAN
PLAN     ──▶  Docs committed, pushed, merged?        Yes ──▶  🚦 REVIEW GATE
🚦 GATE  ──▶  User reviewed docs and said continue?  Yes ──▶  BUILD
BUILD    ──▶  All tasks checked, all tests pass?     Yes ──▶  SHIP
SHIP     ──▶  Merged to main, pushed, branch kept?   Yes ──▶  REFLECT
REFLECT  ──▶  Review doc completed, roadmap updated?  Yes ──▶  DONE ✅
```

---

> *"Think. Design. Plan. Build. Ship. Reflect. Repeat."*

---

*Mastery Framework v1.0*
*Works for any project. Any language. Any stack. Any team.*
