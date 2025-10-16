![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/cybergarage/uecho-go) [![Go](https://github.com/cybergarage/uecho-go/actions/workflows/make.yml/badge.svg)](https://github.com/cybergarage/uecho-go/actions/workflows/make.yml)
 [![Go Reference](https://pkg.go.dev/badge/github.com/cybergarage/uecho-go.svg)](https://pkg.go.dev/github.com/cybergarage/uecho-go)
 [![Go Report Card](https://img.shields.io/badge/go%20report-A%2B-brightgreen)](https://goreportcard.com/report/github.com/cybergarage/uecho-go) 
[![codecov](https://codecov.io/gh/cybergarage/uecho-go/graph/badge.svg?token=UJVU1MNHYD)](https://codecov.io/gh/cybergarage/uecho-go)

![logo](https://raw.githubusercontent.com/cybergarage/uecho-go/master/doc/img/logo.png)

`uecho-go` is a portable and cross-platform development framework for creating controller applications and devices based on [ECHONET Lite][enet] for Go developers. [ECHONET][enet] is an open standard specification for IoT devices in Japan that defines more than 100 IoT device types, including crime prevention sensors, air conditioners, and refrigerators.

## What is uEcho?

`uecho-go` enables developers to easily control [ECHONET Lite][enet] devices or create standard-compliant devices. The framework is designed with object-oriented programming principles, featuring object-oriented naming conventions and functions grouped into classes such as `Controller`, `Node`, `Class`, and `Object`.

![framework](https://raw.githubusercontent.com/cybergarage/uecho-go/master/doc/img/framework.png)

Traditionally, implementing IoT controllers or devices for [ECHONET Lite][enet] required developers to understand and implement complex communication middleware specifications, including message formats and base sequences.

`uecho-go` is also inspired by reactive programming principles. Using `uecho-go`, developers only need to set up basic listeners to implement devices and controllers, as uEcho automatically handles other requests such as property read/write and notification requests.

# Table of Contents

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

[enet]:http://echonet.jp/english/
