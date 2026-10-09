# Verification record

2026-10-09; repository base `7cf62d1`; released simulator tag `v1.0.0` (`9954b99`). No existing ~/Src checkout was edited. Existing uecho-go untracked MRA files were preserved. The implementation changes no core source or existing example source.

- macOS arm64 (Go 1.27.1): `make tui-check` passed formatting, vet, tests, race and both darwin/arm64 and linux/arm64 builds. The modules declare Go 1.25. Test widgets run in actual tcell simulation screens, including a live application loop, background snapshot update, forms/confirmation, search, resize and shutdown with an active worker.
- Native macOS PTY: offline discovery confirmation, keyboard navigation and q/Ctrl-C termination exercised. Terminal input/echo/signal/canonical settings were checked on exit. Darwin's transient PENDIN kernel marker is excluded from comparison. Non-terminal stdin returns a normal `interactive terminal required` error.
- Docker Linux arm64, official `golang:1.25`, **`--network none`**: released-simulator integration passed under race detection, with controller `127.0.0.1:3610` and device `127.0.0.2:3610`. Observed 3 instances; loaded all maps/Get values; lighting SetC 80=30; separate readback TID 002F, local TX/RX timestamps; required status-change INF; process/socket shutdown. Example fixtures also exercise mismatch and unknown outcomes, late/duplicate responses and notifications, and parallel requests.
- Existing root `go vet ./...` passed in a network-disabled Linux container. Existing non-race core/protocol/encoding/transport tests passed using a multicast-capable dummy interface at documentation IP `192.0.2.1/24`, entirely inside that container. Existing search/post/bench examples passed there.

Existing checks have limitations independent of this change:

- Root `go test ./...` fails in `examples/uecholight/light_node_test.go:25`, where the existing test calls `os.Exit(0)` (forbidden by current Go testing). That file is unchanged.
- Root/core race checks detect existing races in controller discovery/post-response and transport socket shutdown, including `net/echonet/controller_message_listener.go:89` and `net/echonet/transport/udp_socket.go:55` versus `ReadMessage:128`. No core race cleanup is included. The new controller uses public transport binding/sending and owns its correlated waits/read lifecycle, avoiding those core Controller/manager paths; its fixtures and released-simulator test pass `-race`.

No household LAN, physical appliances, physical Raspberry Pi or macOS multicast was tested. The current scope includes explicitly selected IPv4 multicast discovery/INF and MRA state/number/raw controls; complex MRA schemas, IPv6 and certified appliance support remain unsupported. The dedicated workflow checks the new example and v1.0.0 loopback and isolated multicast interoperability. See the PR checks for final hosted-CI results.

## PR #7 discovery and MRA follow-up

- Interface picker reads the OS inventory only; it never chooses a network or sends automatically. The multicast response collector checks fresh TID/window, source port, node-profile EOJs, service and complete D6 list; duplicate replies and canceled/late responses cannot complete another request. D5/device INF can add observations without replacing Get values.
- Official MRA 1.3.0 ZIP SHA256 and all 57 class records, definitions and license compared with embedded snapshot; schema editing tests cover enum/range/scaling/signedness, unknown/historical/future release and unsupported types. Actual maps are intersected with static definitions; typed form Review/cancel and invalid-input rejection are exercised on tcell screens.
- Docker Linux arm64 `--network none` + dummy `simtest0` only: released simulator v1.0.0 multicast Get D6 discovered three instances; maps/Get loaded; lighting B0=4B decoded as 75%; SetC followed by fresh Get TID 002F (TX 2026-10-09T06:11:38.499088178Z, RX 2026-10-09T06:11:38.499186803Z); INF received. Both multicast and loopback tests passed `-race`, including normal released-process/private-PTY shutdown. The workflow repeats both tests in `unshare --net`. No physical/household/macOS multicast test was run.

- Follow-up native macOS PTY: confirmed demo discovery/maps, typed enum form, Esc cancellation and q/Ctrl-C exits. A session-local wrapper verified terminal restoration (with transient PENDIN masked) after each exit; no sockets were opened.

## uechoctl tui integration and network default

The earlier records above describe the prior offline-default implementation. The
shared code now lives under `internal/tui` in the root module; the existing
terminal dependency versions are preserved. `uechoctl tui` and compatibility
`uechotui` default to interface-scoped all-device discovery. A sole eligible
address is selected automatically; multiple addresses require picker selection.
`--demo`/`--offline` explicitly selects the socket-free fixture. `/` filters the
listed IP/EOJ/class and never sends discovery; empty means all. `d` repeats
discovery. SetC still requires Review and Confirm.

Local checks: format/vet, fake-client mode selection and cancellation, simulated
screen startup discovery without filter or writes, all existing TUI fixtures
and race checks, root vet/build, encoding/protocol race checks, and Darwin/Linux
arm64 controller builds passed. Temporary GOBIN installed all six existing tools;
help kept get/scan/set and added tui. Native PTY `uechoctl tui --demo` and
`uechotui --offline` discovered the fixture automatically and exited with q.
The network-default CLI was never started on the household LAN. Hosted CI runs
released-simulator loopback/multicast tests only in its isolated namespace.
