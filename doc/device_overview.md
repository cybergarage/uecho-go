![logo](img/logo.png)

# Overview of uEcho Device

## Making Devices

uEcho supports your original standard devices of [ECHONET Lite][enet] specification easily. This document explains to create your original  [ECHONET Lite][enet] device step by step.

## Creating Devices

### 1. Creating Node

To create your original device, use `NewLocalNode()` as the following at first.

```
import (
	"github.com/cybergarage/uecho-go/net/echonet"
)

node := echonet.NewLocalNode()
```

The new node has only a node profile class object, and it has no device object. The node profile object is updated automatically when new devices are added into the node or the any properties in the node are changed.

### 2. Creating Device Object

The new node has no device object. To add your device objects, create a new device object using `NewDeviceWithCode()`.  `NewDeviceWithCode()` create a new device object which is added some mandatory and default properties of ECHONET standard device specification [\[1\]][enet-spec]. Next, set your property data using `Device::SetPropertyrData()`. Then, add the device object into the node using `LocalNode::AddDevice()` as the following:

```
dev := echonet.NewDeviceWithCode(echonet.ObjectCode(0xXXXXXX))
dev.SetPropertyrData(0xXX, ....)

node.AddDevice(dev)
```

### 3. Setting Observers

To implement the device, you have only to handle write requests from other nodes because `uecho-go` handles other standard read and notification requests automatically. To handle the write requests, use `Object::SetListener()` as the following:

```
type ObjectListener interface {
    OnPropertyRequest(obj *Object, esv protocol.ESV, prop *protocol.Property) error
}

type MyNode struct {
    *echonet.LocalNode
}

func NewMyNode() *MyNode {

	node := &MyNode{
		LocalNode: echonet.NewLocalNode(),
	}

  dev := echonet.NewDevice()
  ....
  node.AddDevice(dev)
	dev.SetListener(node)

	return node
}

func (node *MyNode) OnPropertyRequest(obj *echonet.Object, esv protocol.ESV, reqProp *protocol.Property) error {
  // Check whether the property request is a write request
  if !protocol.IsWriteRequest(esv) {
    return nil
  }

  // Check whether the local object (device) has the requested property
  propCode := reqProp.Code()
  prop, ok := obj.GetProperty(propCode)
  if !ok {
    return nil
  }

  .....
  // Set the requested data to the local object (device)
  prop.SetData(reqProp.GetData())

  return nil
}
```

### 4. Start Node

Finally, start the node to use `LocalNode::Start()` as the following:

```
node := echonet.NewLocalNode()
....
err := node.Start()
if err != nil {
  ....
}
```

## Next Steps

Let's check the following documentation to know the device functions of uEcho in more detail.

- [Examples](./examples.md)
- [Inside of uEcho Device](./device_inside.md)

## References

- \[1\] [Detailed Requirements for ECHONET Device objects][enet-spec]

[enet]:http://echonet.jp/english/
[enet-spec]:http://www.echonet.gr.jp/english/spec/index.htm
