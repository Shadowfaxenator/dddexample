# DDD Library Code + Architecture Review (V3)

Date: 2026-02-16

Scope: `/Users/shadowfax/Developer/gw/ddd` only.

## Review method

1. Read all library source files under `/Users/shadowfax/Developer/gw/ddd/pkg` and `/Users/shadowfax/Developer/gw/ddd/internal`.
2. Validated build hygiene:
   - `go test ./...` (passes)
   - `go vet ./...` (passes)
   - `gofmt -l` (1 file reported)
3. Attempted `staticcheck`, but the installed binary was built with an older Go toolchain and could not analyze this module (`go 1.26`).

## Executive summary

The library has a good core shape for DDD + Event Sourcing + CQRS:
- typed aggregate API,
- stream and snapshot abstractions,
- separate driver packages,
- convenience constructor for NATS.

However, there are critical runtime and architecture issues that should be resolved before treating the package as stable public infrastructure. The highest-risk issues are:
- panic in aggregate rebuild after snapshot,
- broken `AtMostOnce` subscription path,
- poison-message redelivery loop,
- event kind parsing coupled to subject tokenization.

## Findings (ordered by severity)

## P0 (must fix)

### P0-1: Panic when snapshot exists and no newer events are returned

Evidence:
- `/Users/shadowfax/Developer/gw/ddd/pkg/aggregate/aggregate.go:133`
- `/Users/shadowfax/Developer/gw/ddd/pkg/aggregate/aggregate.go:142`

What happens:
- In `build`, when load yields `ErrNotExists` and a snapshot is present, code does not return or continue.
- The loop then dereferences `event.Body` while `event` is `nil`, which can panic.

What to change:
1. In `build`, when `errors.Is(err, ErrNotExists)` and `sn != nil`, return the current snapshot aggregate immediately.
2. Add a defensive `if event == nil { ... }` guard before accessing `event.Body`.

Reason:
- This is a hard correctness bug that can crash command handling for snapshotted aggregates with no new events.

### P0-2: `AtMostOnce` subscription path is effectively non-functional

Evidence:
- `/Users/shadowfax/Developer/gw/ddd/pkg/drivers/stream/esnats/eventstream.go:253`
- `/Users/shadowfax/Developer/gw/ddd/pkg/drivers/stream/esnats/eventstream.go:256`
- `/Users/shadowfax/Developer/gw/ddd/pkg/drivers/stream/esnats/adapters.go:56`
- `/Users/shadowfax/Developer/gw/ddd/pkg/drivers/stream/esnats/adapters.go:57`

What happens:
- `AtMostOnce` uses core `Conn().Subscribe(...)` and then `newNatsMessageAdapter` calls `msg.Metadata()`.
- `Metadata()` returns error for non-JetStream messages.
- Handler path exits early on adapter creation error, so messages are not processed in this mode.

What to change:
1. Implement `AtMostOnce` using JetStream consumer settings (`AckNone`) instead of core subscription.
2. Or provide a separate adapter for core NATS messages that does not depend on JetStream metadata.
3. Remove `Ack/Nak` behavior from core-subscribe path if you keep it.

Reason:
- Current behavior breaks advertised QoS semantics and silently drops processing in a public API mode.

### P0-3: Poison-message loop risk on envelope parse failures

Evidence:
- `/Users/shadowfax/Developer/gw/ddd/pkg/drivers/stream/esnats/eventstream.go:216`
- `/Users/shadowfax/Developer/gw/ddd/pkg/drivers/stream/esnats/eventstream.go:220`
- `/Users/shadowfax/Developer/gw/ddd/pkg/drivers/stream/esnats/eventstream.go:221`

What happens:
- If `streamMsgFromNatsMsg` fails, code logs and returns without `Ack`/`Nak`.
- For explicit-ack consumer mode, that message can be redelivered repeatedly.

What to change:
1. Classify envelope parse failure as non-retriable at this layer.
2. Explicitly ack/term such messages so one malformed message does not stall progress.
3. Include structured logging fields (`subject`, headers) for operator diagnosis.

Reason:
- Prevents redelivery storms and consumer starvation.

### P0-4: Event kind parsing is incompatible with dot-containing event names

Evidence:
- `/Users/shadowfax/Developer/gw/ddd/pkg/drivers/stream/esnats/eventstream.go:93`
- `/Users/shadowfax/Developer/gw/ddd/pkg/drivers/stream/esnats/adapters.go:152`
- `/Users/shadowfax/Developer/gw/ddd/pkg/drivers/stream/esnats/adapters.go:156`

What happens:
- Save encodes kind in subject suffix (`...<aggregateID>.<kind>`).
- Load parses kind using `strings.Split(subject, ".")[2]`.
- If kind contains `.`, parsed kind is truncated and deserialization/filtering break.

What to change:
1. Store event kind in a dedicated header and read from header, not subject token index.
2. If subject encoding is retained, validate `WithEvent` names and reject `.` explicitly.

Reason:
- Public API should not have hidden naming traps that cause runtime decode failures.

## P1 (high priority)

### P1-1: Panic-heavy public API/configuration surface

Evidence:
- `/Users/shadowfax/Developer/gw/ddd/pkg/aggregate/options.go:43`
- `/Users/shadowfax/Developer/gw/ddd/pkg/stream/options.go:16`
- `/Users/shadowfax/Developer/gw/ddd/internal/typereg/registry.go:33`
- `/Users/shadowfax/Developer/gw/ddd/internal/typereg/registry.go:36`
- `/Users/shadowfax/Developer/gw/ddd/internal/serde/default.go:22`
- `/Users/shadowfax/Developer/gw/ddd/internal/typereg/registry.go:110`

What to change:
1. Replace panics in configuration/type-registration paths with returned errors from constructors/options.
2. Keep panic-only helpers as explicit `Must...` APIs if needed.

Reason:
- Public libraries should avoid terminating host processes for recoverable misconfiguration.

### P1-2: Layer boundary inversion (driver depends on aggregate package)

Evidence:
- `/Users/shadowfax/Developer/gw/ddd/pkg/drivers/stream/esnats/eventstream.go:11`
- `/Users/shadowfax/Developer/gw/ddd/pkg/drivers/stream/esnats/eventstream.go:171`
- `/Users/shadowfax/Developer/gw/ddd/pkg/drivers/stream/esnats/eventstream.go:217`

What happens:
- Low-level stream driver imports aggregate package and uses aggregate-domain errors (`ErrNotExists`, `InvariantViolationError`).

What to change:
1. Move storage/transport-level sentinel errors into `/pkg/stream` (or a shared transport package).
2. Keep aggregate-domain error interpretation in aggregate adapters, not in drivers.

Reason:
- This improves pluggability for future drivers and prevents cross-layer coupling.

### P1-3: Snapshot load API collapses backend errors into “not found”

Evidence:
- `/Users/shadowfax/Developer/gw/ddd/pkg/snapshot/store.go:61`
- `/Users/shadowfax/Developer/gw/ddd/pkg/snapshot/store.go:68`
- `/Users/shadowfax/Developer/gw/ddd/pkg/snapshot/store.go:69`
- `/Users/shadowfax/Developer/gw/ddd/pkg/snapshot/store.go:78`

What to change:
1. Change snapshot load contract from `(*Snapshot[T], bool)` to `(*Snapshot[T], error)`.
2. Return `ErrNoSnapshot` for misses and propagate real backend/unmarshal errors.

Reason:
- Callers need to distinguish normal cache miss from infrastructure/data-corruption failures.

### P1-4: `WithSnapshotCodec` changes stream codec too (API surprise)

Evidence:
- `/Users/shadowfax/Developer/gw/ddd/pkg/natsaggregate/options.go:64`
- `/Users/shadowfax/Developer/gw/ddd/pkg/natsaggregate/options.go:66`
- `/Users/shadowfax/Developer/gw/ddd/pkg/aggregate/options.go:61`
- `/Users/shadowfax/Developer/gw/ddd/pkg/aggregate/options.go:62`

What to change:
1. Either rename `WithSnapshotCodec` to `WithCodec` to match behavior.
2. Or split into two distinct options: snapshot codec and event codec.

Reason:
- Avoid accidental event wire-format changes from an option name that implies snapshot-only behavior.

### P1-5: Snapshot persistence strategy can create unbounded background work

Evidence:
- `/Users/shadowfax/Developer/gw/ddd/pkg/aggregate/aggregate.go:224`
- `/Users/shadowfax/Developer/gw/ddd/pkg/aggregate/aggregate.go:225`

What to change:
1. Replace per-mutation goroutine with bounded worker/semaphore or synchronous save mode.
2. Make snapshot execution strategy configurable (`sync`/`async`/max-concurrency).
3. Avoid unconditional `context.Background()` for detached writes unless explicitly configured.

Reason:
- Limits goroutine growth and improves shutdown/control behavior for library consumers.

### P1-6: Default durable name can collide across distinct subscriptions

Evidence:
- `/Users/shadowfax/Developer/gw/ddd/pkg/stream/stream.go:167`
- `/Users/shadowfax/Developer/gw/ddd/pkg/aggregate/aggregate.go:158`
- `/Users/shadowfax/Developer/gw/ddd/pkg/aggregate/aggregate.go:278`

What to change:
1. Derive default durable name from handler type + stream + filter hash.
2. Or require explicit durable name whenever a filter option is provided.

Reason:
- Prevents accidental consumer update conflicts and subscription reuse surprises.

## P2 (quality and maintainability)

### P2-1: `context.WithValue` key is exported/mutable

Evidence:
- `/Users/shadowfax/Developer/gw/ddd/pkg/idempotency/key.go:7`

What to change:
1. Use unexported key type and unexported key value.
2. Keep only accessor functions exported.

Reason:
- This follows Go context-key best practice and avoids external key collision/mutation.

### P2-2: Unused/dead API surface

Evidence:
- `/Users/shadowfax/Developer/gw/ddd/pkg/aggregate/aggregate.go:25` (`idempotencyWindow` unused)
- `/Users/shadowfax/Developer/gw/ddd/pkg/aggregate/aggregate.go:114` (`codec` field unused)
- `/Users/shadowfax/Developer/gw/ddd/pkg/aggregate/aggregate.go:40` (`ctx` not used in constructor)
- `/Users/shadowfax/Developer/gw/ddd/pkg/stream/stream.go:51` (`ctx` not used in constructor)
- `/Users/shadowfax/Developer/gw/ddd/pkg/natsaggregate/store.go:17` (`Store` type currently unused)

What to change:
1. Remove unused fields/constants/types.
2. Remove unused constructor params or start using them for initialization/lifecycle.

Reason:
- Cleaner public API and lower maintenance cost.

### P2-3: Global logging side effect in internal type registry

Evidence:
- `/Users/shadowfax/Developer/gw/ddd/internal/typereg/registry.go:40`

What to change:
1. Remove direct `slog.Info` from `Register`.
2. If needed, route registration logs through injected logger at a higher layer.

Reason:
- Libraries should avoid unsolicited global logs by default.

### P2-4: Go formatting drift

Evidence:
- `gofmt -l` reports `/Users/shadowfax/Developer/gw/ddd/pkg/identity/id.go`.

What to change:
1. Run `gofmt` on the file.

Reason:
- Keep exported library code idiomatic and consistent.

## Architecture assessment

## What is strong

1. Clean package split between aggregate logic (`/pkg/aggregate`), transport-neutral stream API (`/pkg/stream`), snapshot API (`/pkg/snapshot`), and concrete drivers.
2. Generic aggregate typing gives strong compile-time domain constraints.
3. Event registration + codec abstraction gives a practical extension point.

## What needs architectural adjustment for long-term pluggability

1. Keep domain package (`aggregate`) strictly above infra drivers.
2. Standardize error model across packages.
3. Make event envelope format explicit and transport-independent.
4. Separate API options by concern (event codec vs snapshot codec, delivery options vs naming options).

Recommended dependency direction:

```text
pkg/aggregate  -> pkg/stream (interfaces) -> pkg/drivers/stream/*
pkg/aggregate  -> pkg/snapshot (interfaces) -> pkg/drivers/snapshot/*
pkg/aggregate  -> pkg/codec, pkg/idempotency, pkg/identity
```

No dependency in the reverse direction (drivers should not import aggregate-domain errors/types).

## Concrete change plan

1. Fix P0 correctness issues first (`aggregate.build`, `AtMostOnce`, parse-failure ack strategy, event-kind encoding).
2. Refactor error contracts next (remove panics from option/config paths, make snapshot load return errors).
3. Decouple layer boundaries (`esnats` from `aggregate` package).
4. Clean API ergonomics (durable naming strategy, codec option clarity, remove dead surface).
5. Apply formatting and logging cleanup for public-library polish.

