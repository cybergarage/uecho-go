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

Run `make tui` for a socket-free controller demo. The separate
[uechotui command](cmd/uechotui/README.md) provides device selection,
interface-selected multicast discovery, MRA settings with raw Get values,
confirmed SetC/fresh Get verification, and
protocol/INF logs. Enable networking with `--network` to choose the local
interface/address in the UI, or specify both `--interface` and `--bind`.
`--peer` is optional and selects unicast discovery instead of multicast.
See its README for isolated simulator v1.0.0 verification and limitations.

Install the command-line tools with `make install`. This runs `go install` for
`uechoctl`, `uechopost`, `uechosearch`, `uechotui`, and the existing
`uecholight`/`uechobench` examples. Set `GOBIN` to choose the destination
(for example, `GOBIN=/tmp/uecho-bin make install`). The TUI uses a separate
module under `cmd/uechotui` to keep terminal dependencies out of the library.
