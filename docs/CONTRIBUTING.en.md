# Contributing

[简体中文](CONTRIBUTING.md) | [English](CONTRIBUTING.en.md) | [Русский](CONTRIBUTING.ru.md)

Thank you for contributing to `go-dtls`. Changes must preserve DTLS 1.3 datagram semantics, protocol security, and existing performance. Do not prebuild frameworks for extensions without a concrete requirement.

## Development Environment

- Go 1.26 or later.
- golangci-lint v2.12.2, matching CI.
- Windows race tests require Zig 0.17 and a working CGO toolchain.
- When a change affects a feature supported by wolfSSL, wolfSSL interoperability tests and focused benchmarks are required; other changes do not require wolfSSL.

## Required Checks

All Go source must be formatted with `gofmt`. Every commit and pull request must pass lint, test, vet, and race. Linux, macOS, and CI use:

```sh
go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run ./...
go test ./... -count=1
go vet ./...
go test -race ./... -count=1
```

Protocol, concurrency, or resource-lifecycle changes must also run repeated tests:

```sh
go test ./... -count=10 -timeout=10m
```

On Windows, run the first three checks and replace the race command above with the repository script:

```powershell
.\tools\test-race.ps1
```

GitHub Actions runs the four required checks independently on every push and pull request. Fix every failure; do not bypass it by weakening assertions, skipping tests, or broadening lint exclusions.

## Protocol Changes

Protocol behavior follows RFC 9147, RFC 9846, and applicable related RFCs:

- `MUST`, `MUST NOT`, `REQUIRED`, `SHALL`, and `SHALL NOT` must be enforced in both send and receive directions by clients and servers.
- The client actively implements `SHOULD`-class requirements. The server may tolerate peer deviations only when authentication, confidentiality, replay protection, amplification limits, and state consistency are not weakened.
- Once a `MAY` or `OPTIONAL` capability is negotiated, its conditional mandatory requirements must be fully implemented.
- Preserve the mapping of one `WriteDatagram` to one Application Data record. Do not introduce implicit application-data fragmentation, ordering, or retransmission.
- Reuse the existing record, flight, ACK, transcript, key-schedule, CID, and error abstractions instead of duplicating their state machines.

Protocol changes must cover at least:

- The normal path and malformed, truncated, duplicate, unknown, and out-of-range inputs.
- Every affected branch for clients and servers, send and receive, full handshakes and resumption.
- Loss, delay, reordering, and duplication; retransmission or state-machine changes require high-repeat validation.
- Related security boundaries, including anti-replay, amplification limits, AEAD usage limits, secret clearing, and resource limits.
- Resource lifecycles for timeouts, closure, goroutines, memory, and Listener sessions.

For retransmission, handshake, or protocol-state-machine changes, run the relevant weak-network and end-to-end tests with `-count=100`. Run resource-lifecycle tests with at least `-count=10`.

Optional extensions require a concrete deployment need and a testable protocol design. Do not add placeholder APIs or empty frameworks in advance.

## Performance Validation

Changes to handshake, record, parsing, allocation, or concurrency hot paths must check both microbenchmarks and full connections. Full connections are the primary performance gate; local microbenchmarks cannot establish overall performance.

```sh
go test -run '^$' -bench . -benchmem
go test -run '^$' -bench '^BenchmarkConnectionHandshakeLifecycle$' -benchmem -benchtime=2000x -cpu=1
go test -run '^$' -bench '^BenchmarkMutualTLSHandshakeLifecycle/(Full|Resumed)$' -benchmem -count=10
go test -run '^$' -bench '^BenchmarkCertificateCompression' -benchmem
go test -run '^$' -bench '^BenchmarkProtectedRecord(Seal|RoundTripInPlace)$' -benchmem -count=5
```

When comparing performance before and after a change:

- Keep the Go version, CPU, `-cpu`, and `-benchtime` fixed.
- Precompile test binaries before and after the change. For full connections, use `-cpu=1`, 500 iterations per sample, alternate the binaries for at least nine rounds, and compare paired medians.
- Check `ns/op`, `B/op`, `allocs/op`, full-connection throughput, and resource lifecycles together.
- When allocations or memory differ, use profiles to identify the exact object and call path.
- Store benchmark output, profiles, and temporary binaries outside the repository. Do not commit generated artifacts.

The pull request benchmark workflow runs this base/head paired gate automatically. It fails when the paired median time increases by more than 5% and most pairs are slower, or when median `B/op` or `allocs/op` increases with at least 75% of pairs in the same direction. Benchmark workload source changes are incomparable and fail by default; after maintainer review, the `performance-workload-approved` label may rerun the workflow, but the report does not generate deltas between different workloads. Missing benchmarks, metrics, or paired samples also fail, and the PR comment and artifact retain the decision and raw output.

README files record only representative performance for the current environment. A/B data, frozen baselines, and before/after narratives belong in pull request validation evidence.

## Interoperability Validation

When the peer supports a feature, add bidirectional third-party end-to-end tests. Run wolfSSL tests with:

```powershell
$env:WOLFSSL_ROOT = 'C:\path\to\wolfssl'
go test -run TestInteropWolfSSL -v -count=10
```

When wolfSSL supports the target feature, also add and run a focused real-UDP benchmark covering go-dtls -> go-dtls, go-dtls -> wolfSSL, wolfSSL -> go-dtls, and wolfSSL -> wolfSSL. Include the peer commit/build configuration, workload, timing boundaries, five-run median, and capability limits in the pull request validation evidence; do not put cross-implementation benchmark data in README files.

The formal four-way benchmark uses `go test -json`. Every applicable workload direction must produce five results or match an exact skip reason and evidence registered for the full wolfSSL SHA. A new SHA, recovered capability, changed reason, unknown skip, missing direction, or insufficient sample count fails report generation and the workflow. Allowed peer limitations must remain visible in the report.

The benchmark workflow uses the exact wolfSSL SHA resolved for the run and executes the complete interoperability matrix before performance tests: one iteration for pull requests and ten for `master` pushes, scheduled runs, and manual dispatches. Any unexpected failure blocks reporting and publication. The wolfSSL build cache is saved before the matrix, so a test failure does not force recompilation on the next run.

The benchmark workflow updates the separate `benchmark-results` branch after a push to `master`. The scheduled check triggers and publishes only when the go-dtls or wolfSSL `master` SHA changed and no benchmark for that go-dtls SHA is queued or running. A pull request only updates its benchmark comment and never writes to `master` or the results branch. Do not maintain performance numbers manually in README files.

Adding or renaming any benchmark or sub-benchmark requires updating its category, order, and multilingual labels in `tools/benchmark-report`. The report generator rejects benchmark output without explicit coverage, and the benchmark workflow must pass.

When a third-party implementation does not enable or support the target extension, state the interoperability boundary explicitly. Do not report a skipped test as passing.

## Documentation

- Public API, configuration, error, or protocol-behavior changes must update Go documentation and all README language versions.
- Keep the Chinese, English, and Russian content synchronized and all language-switch links valid.
- README files describe only existing capabilities, RFC completion, current limitations, usage, representative performance, and test coverage.
- List unsupported capabilities briefly. Development plans, priorities, research history, and historical performance comparisons do not belong in README files.

## Commits and Pull Requests

Every commit subject must follow Conventional Commits:

```text
type(scope): description
```

Examples:

```text
feat(cid): implement RFC 9853 return routability checks
fix(record): reject invalid truncated sequence numbers
docs(api): clarify datagram truncation behavior
```

The pull request title must use the same format. The Conventional Commits workflow checks both the PR title and every commit subject in the branch, and reruns when the title is edited.

A pull request must identify the applicable RFC sections, behavior changes, test results, race results, performance data, and interoperability scope. Use `!` for a breaking change and explain migration in the body.
