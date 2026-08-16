# Freshness Phase 11B changelog

## 2026-08-16 — catalog first-publish fact contract (write-ahead)

Recorded before code.

- `WorkUpsertedEvent` must carry catalog-owned `first_publish_year` as an
  explicit fact. Zero means unknown.
- recommendation-service must use this fact for recent-release selection;
  event occurrence/synchronization time is never a publication-year proxy.
- The change is a clean current contract with no legacy fallback. Publish a
  new proto-common tag and consume that tag normally; local `replace` is
  forbidden.

Implemented:

- added `WorkUpsertedEvent.first_publish_year = 20` with zero-as-unknown
  semantics;
- regenerated `events/events.pb.go` from the proto;
- `go test ./...`, `go vet ./...`, and `git diff --check` passed.
