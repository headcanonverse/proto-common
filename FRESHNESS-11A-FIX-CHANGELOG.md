# Freshness 11A fix — write-ahead changelog

## Checkpoint 0 — 2026-08-16, before code

- Initial HEAD: `798349a` (`master`, published `v0.27.0`).
- Audit finding: `WorkUpsertedEvent` is simultaneously used as a rebuildable
  catalog projection and as a publication-domain event by home-feed-service.
  Catalog metadata replay therefore turns every published catalog row into a
  fresh release and cannot express ordered unpublish/restore state.
- Scope: add a dedicated, versioned `WorkPublicationChangedEvent` contract.
  `WorkUpsertedEvent` remains catalog facts for reco/search only.
- The publication event is a complete home-feed projection: work identity,
  current availability, stable first-publication time, display facts,
  catalog `metadata_version`, and occurrence time.
- No compatibility fields or dual-write fallback. Downstream home-feed must
  stop consuming `WorkUpsertedEvent` for releases.
- Planned gate: deterministic generation, Go test/vet, diff check, published
  tag `v0.28.0`; no local module replacement.

## Checkpoint 1 — contract generated and verified

- Added `WorkPublicationChangedEvent` as the complete ordered release
  projection. No field was added to `WorkUpsertedEvent`, so reco/search
  metadata payload hashes remain stable across ordinary delivery and replay.
- First `make gen` failed because `protoc-gen-go` was not on the default PATH;
  no output was produced by that attempt. Re-run with `$HOME/go/bin` used the
  repository's existing generator versions (`protoc-gen-go v1.36.1`, protoc
  35.1) and succeeded.
- `go test ./...`, `go vet ./...`, `git diff --check`: PASS.
- Next: commit, tag and publish `v0.28.0`, then downstream consumers may bump
  the exact public module tag. No downstream local replace is permitted.
