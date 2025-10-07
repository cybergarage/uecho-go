## uechoctl set

Set property value to Echonet Lite device.

### Synopsis

Set property value to Echonet Lite device. Object code, property code and property value must be specified in hexadecimal format.

```
uechoctl set <node-address> <object-code> <property-code> <property-value> [flags]
```

### Examples

```
set 192.168.1.100 013001 80 30
```

### Options

```
  -h, --help   help for set
```

### Options inherited from parent commands

```
      --format string   output format: table|json|csv (default "table")
      --verbose         enable verbose output
```

### SEE ALSO

* [uechoctl](uechoctl.md)	 - 

