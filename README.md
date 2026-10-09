![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/cybergarage/uecho-go) [![Go](https://github.com/cybergarage/uecho-go/actions/workflows/make.yml/badge.svg)](https://github.com/cybergarage/uecho-go/actions/workflows/make.yml)
 [![Go Reference](https://pkg.go.dev/badge/github.com/cybergarage/uecho-go.svg)](https://pkg.go.dev/github.com/cybergarage/uecho-go)
 [![Go Report Card](https://img.shields.io/badge/go%20report-A%2B-brightgreen)](https://goreportcard.com/report/github.com/cybergarage/uecho-go) 
[![codecov](https://codecov.io/gh/cybergarage/uecho-go/graph/badge.svg?token=UJVU1MNHYD)](https://codecov.io/gh/cybergarage/uecho-go)

![logo](https://raw.githubusercontent.com/cybergarage/uecho-go/master/doc/img/logo.png)

`uecho-go` is a portable, cross-platform framework for developing [ECHONET Lite][enet] controllers and devices in Go. ECHONET Lite is an open standard for IoT devices in Japan that defines more than 100 IoT device types, including security sensors, air conditioners, and refrigerators.

## What is uEcho?

`uecho-go` provides APIs for controlling [ECHONET Lite][enet] devices and implementing device applications. It follows object-oriented naming conventions, with components such as `Controller`, `Node`, `Class`, and `Object`.

![framework](https://raw.githubusercontent.com/cybergarage/uecho-go/master/doc/img/framework.png)

Implementing ECHONET Lite controllers and devices from scratch requires handling protocol details such as message formats and communication sequences.

`uecho-go` is also inspired by reactive programming principles. It handles property read/write requests and notifications, allowing developers to focus on application logic by configuring listeners.

## Table of Contents

- **Controller**
  - [Overview of Controller](https://github.com/cybergarage/uecho-go/blob/master/doc/controller_overview.md)
  - [Inside of Controller](https://github.com/cybergarage/uecho-go/blob/master/doc/controller_inside.md)
- **Device**
  - [Overview of Device](https://github.com/cybergarage/uecho-go/blob/master/doc/device_overview.md)
  - [Inside of Device](https://github.com/cybergarage/uecho-go/blob/master/doc/device_inside.md)
- **Examples**
  - [uechoctl](https://github.com/cybergarage/uecho-go/blob/master/doc/cmd/uechoctl.md)
  - [Examples](https://github.com/cybergarage/uecho-go/blob/master/doc/examples.md)
- **Appendix**
  - [Extended Configurations for Go Platform](https://github.com/cybergarage/uecho-go/blob/master/doc/extensions.md)

## Related projects

[uecho-simulator](https://github.com/cybergarage/uecho-simulator) is a small ECHONET Lite development simulator with virtual lighting, air conditioning, and temperature sensing. It provides a full-screen terminal UI and a live, read-only browser preview, runs offline by default, and implements limited device profiles. It is built with `uecho-go` and serves as an example application.

[enet]:https://echonet.jp/english/

## Fullscreen developer controller

Run `uechoctl tui` (or `make tui`) to discover all ECHONET Lite devices on
one local IPv4 interface and open the fullscreen controller. A sole eligible
interface/address is selected automatically; multiple addresses open a picker.
Use `uechoctl tui --demo` for the socket-free fixture.

The [TUI guide](cmd/uechotui/README.md) covers device selection, MRA settings
with raw Get values, property Enter editing, Enter-to-send SetC/fresh Get verification, and protocol/INF
logs. `/` filters the listed devices by IP, EOJ or class; an empty filter shows
all devices. `d` or F5 repeats discovery (F5 also works while filtering). No filter or peer IP is needed to discover
all devices. Specify both `--interface` and `--bind` to choose local values
explicitly; optional `--peer` selects unicast discovery instead of multicast.
See the guide for isolated simulator v1.0.0 verification and limitations.

Install the command-line tools with `make install`. This runs `go install` for
`uechoctl`, `uechopost`, `uechosearch`, the compatibility alias `uechotui`, and
the existing `uecholight`/`uechobench` examples (six binaries). Set `GOBIN` to
choose the destination (for example, `GOBIN=/tmp/uecho-bin make install`).
`uechoctl tui` and `uechotui` share one implementation in the root module;
existing `uechoctl` subcommands remain available.
