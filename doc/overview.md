![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/cybergarage/uecho-go) [![Go](https://github.com/cybergarage/uecho-go/actions/workflows/make.yml/badge.svg)](https://github.com/cybergarage/uecho-go/actions/workflows/make.yml)
 [![Go Reference](https://pkg.go.dev/badge/github.com/cybergarage/uecho-go.svg)](https://pkg.go.dev/github.com/cybergarage/uecho-go)
 [![Go Report Card](https://img.shields.io/badge/go%20report-A%2B-brightgreen)](https://goreportcard.com/report/github.com/cybergarage/uecho-go) 
[![codecov](https://codecov.io/gh/cybergarage/uecho-go/graph/badge.svg?token=UJVU1MNHYD)](https://codecov.io/gh/cybergarage/uecho-go)

![logo](https://raw.githubusercontent.com/cybergarage/uecho-go/master/doc/img/logo.png)

The `uecho-go` is a portable and cross platform development framework for creating controller applications and devices of [ECHONET Lite][enet] for Go developers. [ECHONET][enet] is an open standard specification for IoT devices in Japan, it specifies more than 100 IoT devices such as crime prevention sensor, air conditioner and refrigerator.

## What is uEcho?

The `uecho-go` supports to control devices of [ECHONET Lite][enet] or create the standard devices of the specification easily. The `uecho-go` is designed in object-oriented programming, and the functions are object-oriented in their naming convention, and are grouped into classes such as `Controller`, `Node`, `Class` and `Object`.

![framwork](https://raw.githubusercontent.com/cybergarage/uecho-go/master/doc/img/framework.png)

To implement IoT controller or devices of [ECHONET Lite][enet], the developer had to understand and implement the communication middleware specification such as the message format and base sequences.

The `uecho-go` is inspired by reactive programming too. Using The `uecho-go`, developer have only to set basic listeners to implement the devices and controllers because uEcho handles other requests such as request and notification requests automatically.

# Table of Contents

- Controller
  - [Overview of Controller](https://github.com/cybergarage/uecho-go/blob/master/doc/controller_overview.md)
  - [Inside of Controller](https://github.com/cybergarage/uecho-go/blob/master/doc/controller_inside.md)
- Device
  - [Overview of Device](https://github.com/cybergarage/uecho-go/blob/master/doc/device_overview.md)
  - [Inside of Device](https://github.com/cybergarage/uecho-go/blob/master/doc/device_inside.md)
- [Examples](https://github.com/cybergarage/uecho-go/blob/master/doc/examples.md)
- Appendix
  - [Expanded Configurations for Go Platform](https://github.com/cybergarage/uecho-go/blob/master/doc/extensions.md)

[enet]:http://echonet.jp/english/
