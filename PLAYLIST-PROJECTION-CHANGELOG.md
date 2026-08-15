# Playlist search projection contract

## Scope

- Add one versioned protobuf event carrying the complete catalog-owned playlist
  search projection.
- Represent online projection changes explicitly as `UPSERT` or `DELETE`; do
  not preserve or introduce legacy event variants.
- Include playlist identity, owner identity, title, description, cover URL,
  work count, author user IDs, public/deleted state, source update time,
  monotonic source version, and event ID.
- Keep this repository limited to the cross-service wire contract. Catalog
  publication and search consumption are separate implementation checkpoints.

## Invariants

- Kafka message key is the playlist ID.
- `event_id`, `playlist_id`, `owner_id`, `source_updated_at`, and a positive
  `projection_version` are mandatory.
- `UPSERT` is a complete snapshot and never a patch.
- `DELETE` is a tombstone. It carries source ordering metadata and sets
  `deleted=true`; consumers remove the indexed document.
- A newer `projection_version` wins. Equal versions are idempotent; older
  versions cannot overwrite or resurrect newer state.
- No backward-compatibility envelope, fallback parser, or dual-write contract
  is added.

## Checkpoints

- [x] Record scope and invariants before changing protobuf code.
- [x] Add and generate the versioned protobuf event.
- [x] Validate generated code and repository build/tests.
- [x] Commit locally without tagging or pushing.

## Validation

- `PATH="$HOME/go/bin:$PATH" make gen`
- `go test ./...`
- `go build -o /tmp/proto-common-events.a ./events`
- `git diff --check`

Coordinator handoff: publish the contract as the next `proto-common` tag, then
bump `search-service` and `catalog-service` to that tag without a local module
replacement.
