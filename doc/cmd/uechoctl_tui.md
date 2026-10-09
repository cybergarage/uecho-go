## uechoctl tui

Discover all ECHONET Lite devices and open the terminal controller

### Synopsis

Open the network controller and discover all devices on the selected interface. Use --demo for a socket-free fixture. / filters the listed devices; d/F5 repeats discovery. Writes always require confirmation.

```
uechoctl tui [flags]
```

### Options

```
      --bind string        local IPv4; specify together with --interface
      --demo               use the offline demo; open no sockets
  -h, --help               help for tui
      --interface string   local interface; specify together with --bind
      --network            network mode (default); --network=false selects demo (default true)
      --offline            alias for --demo
      --peer string        optional literal IPv4 for unicast discovery instead of all-device multicast
```

### Options inherited from parent commands

```
      --format string   output format: table|json|csv (default "table")
      --verbose         enable verbose output
```

### SEE ALSO

* [uechoctl](uechoctl.md)	 - Control Echonet Lite devices from command line.

