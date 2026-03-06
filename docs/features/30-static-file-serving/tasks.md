# Feature #30 — Static File Serving: Tasks

## Implementation tasks

| # | Task | Status |
|---|------|--------|
| 1 | Verify `Static()` and `StaticFile()` exist on Router | ✅ Already done (#07) |
| 2 | Verify `/static` and `/uploads` routes in RouterProvider | ✅ Already done (#07, #17) |
| 3 | Verify tests TC-04, TC-05 cover functionality | ✅ Already done (#07) |
| 4 | Update roadmap | ⬜ |

## Acceptance criteria

All criteria met by Features #07 and #17:

- [x] `Router.Static()` serves directory contents.
- [x] `Router.StaticFile()` serves individual files.
- [x] `/static` route maps to `resources/static/`.
- [x] `/uploads` route maps to `storage/uploads/`.
- [x] Routes are conditional (only registered if directories exist).
- [x] Tests pass.
