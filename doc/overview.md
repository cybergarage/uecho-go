![logo](img/logo.png)

uEcho for Go, `uecho-go`, is a portable and cross platform development framework for creating controller applications and devices of [ECHONET Lite][enet]. 

[ECHONET][enet] is an open standard specification for IoT devices in Japan, it specifies more than 100 IoT devices such as crime prevention sensor, air conditioner and refrigerator.

# Overview

The `uecho-go` supports to control any [ECHONET Lite][enet] devices and create the standard devices easily for Go.

![framwork](img/framework.png)

To implement IoT controller or devices of [ECHONET Lite][enet], the developer had to understand and implement the communication middleware specification such as the message format and base sequences.

Using the `uecho-go`, developer have only to set basic listeners to implement the devices and controllers because `uecho-go` handles other requests such as request and notification requests automatically.

# Table of Contents

- [What is uEcho ?](overview.md)
- Controller
  - [Overview of Controller](controller_overview.md)
  - [Inside of Controller](controller_inside.md)
- Device
  - [Overview of Device](device_overview.md)
  - [Inside of Device](device_inside.md)
- Appendix
  - [Extension Specifications](extension.md)
- [Examples](examples.md)

[enet]:http://echonet.jp/english/
