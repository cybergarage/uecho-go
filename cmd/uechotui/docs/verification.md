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

## Ephemeral reply source port and rediscovery follow-up

Provided RX packets were replayed in memory: TID 0001, source
192.168.100.44:65468 with D6 instances 0F2001/029101, and
192.168.100.216:35449 with 05FF01. The regression failed before the fix with zero
listed instances and passes afterward with all three. These logs do not establish
which responder is the user's simulator. Official Part II section 1.2 leaves UDP
source ports unspecified and requires destination 3610; matching retains source
IP (for unicast requests), TID, EOJs, ESV/EPC/count and response windows. Get/SetC
responses and INF now permit ephemeral source ports too; outgoing requests always
use 3610 and never adopt a response port.

Fake/simulation tests cover invalid and late discovery replies, repeat fresh TIDs,
d from each normal pane, F5 from the filter, canceled confirmation, busy-operation
rejection, Esc worker cancellation, and clearing the filter on confirmed rediscovery.
The interface picker starts on its dropdown, supports selecting another address,
and Esc cancels before binding. Native macOS PTY explicitly used --demo: F5 while
filtering rediscovered the fixture with TID 0002; d confirmation cancellation and
repeat discovery were checked before q exit. No live LAN test or interference
with an existing simulator was performed. The existing isolated CI repeats
released-simulator loopback/multicast interoperability.

## Device layout and property editing follow-up

The supplied Library screenshot was resolved but its byte transfer returned 403;
its pixels were not inspected. Layout changes were grounded in the fixed 38-column
pane and verified with newly generated actual tcell screenshots. Device rows now
have separate IPv4/EOJ/name columns; wide terminals allocate 40% to devices while
preserving at least 64 columns for properties. Below 120 columns or 28 rows, panes
stack at full width. Resize tests cover 45/68/119/132/160 columns and ensure pane
bounds, full identities/names at supported widths, and compact key visibility.

Property Enter now opens the existing MRA editor (w remains an alias). Device
selection automatically loads fresh maps/Get, r refreshes the device, and g reads
a Get-permitted property. Switching cancels prior reads and queues only the latest
selection; delayed canceled TIDs cannot populate the new device. A confirmed
SetC/readback is not canceled by selection changes. Known write-only schemas use
Set permission but send no forbidden Get; acknowledgment remains unverified with
readback unavailable. Readback mismatch/failure still never becomes success.

Fake/simulation/race tests cover enum/numeric Enter, invalid number, cancel,
read-only/unknown schemas, write-only acknowledgment, delayed replies, rapid
selection changes and protected SET jobs. Native PTY explicitly used --demo:
startup automatically loaded maps/values; property Enter opened enum Off (31),
Review/Confirm sent fake SetC TID 000B and fresh Get TID 000C with readback success;
numeric 50% produced review EDT 32 and was canceled before q exit. Screenshots
were regenerated for normal/wide/compact layouts and enum/number forms.

Read-only simulator investigation found the existing process launched with only
--display 127.0.0.1:8080; make preview defaults to that command and no UDP listener.
Its local interface/IP inventory was en8/192.168.100.25 and en1/192.168.100.26;
192.168.100.24 was absent. The preview supports PREVIEW_ARGS with --udp,
--allow-lan and --multicast-interface, exposing all three modeled devices through
one chosen interface. No simulator process, LAN traffic, or appliance was used or
modified by this task. Hosted CI retains isolated released-simulator checks.

## Enum commit and direct-send workflow (PR after #13)

The supplied TID 0097 Set_Res acknowledges processing but contains no original
EDT. TID 0098 Get_Res returns lighting EPC 80=31 (OFF). It cannot establish that
30 was sent. Read-only inspection traced simulator Engine.HandleAll ->
Store.Write -> revision/STATE/notify -> /events snapshot -> room.html power label;
there is no demonstrated simulator defect requiring a simulator patch.

Before: a focused enum Down key displayed an OFF candidate while
GetCurrentOption still returned index 0 / On (30). The new regression test
TestEnumArrowCommitsDisplayedValue failed on the prior implementation. After:
arrows commit the displayed option, and leaving an open list with Tab commits
its visible candidate. Space opens a list; its Enter only commits, while the
next Enter sends. Editor and interface picker have no Cancel/Review buttons;
Esc cancels/exits and Return sends/selects. Editors initially focus the value
and show target, EPC, decoded value, EDT, readback availability and key actions.
Opening Enter does not cascade into SET, and rapid repeated Enter cannot create
a second write while the first operation is busy. Device changes invalidate an
open editor. Multiple enum choices, numeric invalid input, cancellation and
read-only/write-only regression coverage remain in the fake-screen tests.

Protocol logs now reserve approximately one third of the screen. `l` opens a
full-screen wrapped log that retains exact TX/RX frames, SET-ACK (which is not
verified success), and RESULT with both requested/readback EDT. Mismatch details
remain available after other operations. Updated screenshots use actual cells.

The opt-in TestSimulatorUIEnumToDisplay runs only in the CI Linux network
namespace, driving actual focused-widget UI keys OFF -> ON -> OFF through the
released v1 simulator. It checks separate SET/Get TIDs, 31/30/31 readback,
model power, and the GUI's actual /events SSE power/STATE payload. The existing
loopback/multicast regressions also run there. No household discovery/write or
existing simulator stop/restart is performed. Physical M4/M6 LAN operation and
browser rendering of the existing GUI require the user's follow-up test.

Local native PTY demo additionally verified direct Enter OFF: SetC TID 000B /
Get 000C EDT 31, then ON: SetC 000D / Get 000E EDT 30, with the full wrapped
RESULT panel. `q` restored the terminal and exited zero. Temporary GOBIN
`/tmp/uecho-enum-bin` contains all six tools; root vet/build and help succeeded.
The integration display requires its own PTY (plain/display are mutually
exclusive). Integration packages run serially because they share loopback 3610.
The released simulator starts OFF: its initial no-op OFF write has no STATE
change event; ON and subsequent OFF must both publish STATE in the GUI SSE.

## Normal network interface candidates (after PR #14)

Synthetic interface/address fixtures reproduced the old behavior: up multicast
loopback adapters and 127/8 addresses entered the normal picker. The same tests
now pass after filtering loopback flags and non-LAN IPv4 addresses. Down and
non-multicast adapters, address-enumeration failures, malformed CIDRs, IPv6,
unspecified/multicast/broadcast/link-local addresses are excluded. Multiple IPv4
addresses on one adapter and addresses on multiple eligible adapters are retained.
Existing sole-address auto-selection, multiple-choice Return/Esc and zero-choice
regressions pass. Zero choices explains --demo and explicit isolated bindings;
there is no localhost fallback. A fake-client command regression confirms that
explicit --interface lo --bind 127.0.0.1 --peer 127.0.0.2 bypasses normal selection
and retains the unicast path. No real LAN startup or existing simulator changes.
Local make tui-check passed format, vet, regressions/race and both portable builds.
