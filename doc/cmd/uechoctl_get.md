## uechoctl get

Get property value from Echonet Lite device.

### Synopsis

Get property value from Echonet Lite device. Object and property codes must be specified in hexadecimal format.

```
uechoctl get <node-address> <object-code> <property-code> [flags]
```

### Examples

```
  uechoctl get 192.168.1.100 013001 80
  uechoctl get 10.0.0.50 028001 B0
```

### Options

```
  -h, --help   help for get
```

### Options inherited from parent commands

```
      --format string   output format: table|json|csv (default "table")
      --verbose         enable verbose output
```

### SEE ALSO

* [uechoctl](uechoctl.md)	 - Control Echonet Lite devices from command line.

