# Verification record

2026-10-09; repository base `7cf62d1`; released simulator tag `v1.0.0` (`9954b99`). No existing ~/Src checkout was edited. Existing uecho-go untracked MRA files were preserved. The implementation changes no core source or existing example source.

- macOS arm64 (Go 1.27.1): `make tui-check` passed formatting, vet, tests, race and both darwin/arm64 and linux/arm64 builds. The modules declare Go 1.25. Test widgets run in actual tcell simulation screens, including a live application loop, background snapshot update, forms/confirmation, search, resize and shutdown with an active worker.
- Native macOS PTY: offline discovery confirmation, keyboard navigation and q/Ctrl-C termination exercised. Terminal input/echo/signal/canonical settings were checked on exit. Darwin's transient PENDIN kernel marker is excluded from comparison. Non-terminal stdin returns a normal `interactive terminal required` error.
- Docker Linux arm64, official `golang:1.25`, **`--network none`**: released-simulator integration passed under race detection, with controller `127.0.0.1:3610` and device `127.0.0.2:3610`. Observed 3 instances; loaded all maps/Get values; lighting SetC 80=30; separate readback TID 002F, local TX/RX timestamps; required status-change INF; process/socket shutdown. Example fixtures also exercise mismatch and unknown outcomes, late/duplicate responses and notifications, and parallel requests.
- Existing root `go vet ./...` passed in a network-disabled Linux container. Existing non-race core/protocol/encoding/transport tests passed using a multicast-capable dummy interface at documentation IP `192.0.2.1/24`, entirely inside that container. Existing search/post/bench examples passed there.

Existing checks have limitations independent of this change:

- Root `go test ./...` fails in `examples/uecholight/light_node_test.go:25`, where the existing test calls `os.Exit(0)` (forbidden by current Go testing). That file is unchanged.
- Root/core race checks detect existing races in controller discovery/post-response and transport socket shutdown, including `net/echonet/controller_message_listener.go:89` and `net/echonet/transport/udp_socket.go:55` versus `ReadMessage:128`. No core race cleanup is included. The new controller uses public transport binding/sending and owns its correlated waits/read lifecycle, avoiding those core Controller/manager paths; its fixtures and released-simulator test pass `-race`.

No household LAN, physical appliances, physical Raspberry Pi or macOS multicast was tested. The initial scope is explicit-peer unicast discovery/Get/SetC/INF diagnostics; standard multicast discovery/reception and schema-aware controls remain deferred. The dedicated workflow checks the new example and v1.0.0 loopback interoperability. See the PR checks for final hosted-CI results.
