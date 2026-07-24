module github.com/headcanonverse/proto-common

go 1.23.0

require google.golang.org/protobuf v1.36.1

// v0.25.0 was briefly published with conflicting immutable module content.
// Consumers must use v0.25.1, which contains the canonical WorkPurgedEvent.
retract v0.25.0
