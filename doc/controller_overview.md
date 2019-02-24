![logo](img/logo.png)

# Overview of uEcho Controller

The controller is a special node of [ECHONETLite][enet] to control other nodes, it can find the nodes in the local area network and send any messages into the found devices.

## Creating Controller

### 1. Starting Controller

To start a controller, create a controller using `uecho_controller_new` and start the controller using `uecho_controller_start` as the following:

```
import (
	"github.com/cybergarage/uecho-go/net/echonet"
)

ctrl := echonet.NewController()

err := ctrl.Start()
    ....
}
```

### 2. Searching Nodes

Next, use `Controller::SearchAllObjects()` to search other nodes in the local area network as the following:

```
ctrl := echonet.NewController()
....
err = ctrl.SearchAllObjects()
if err != nil {
    ....
}
```

### 3. Getting Nodes and Objects

After the searching, use `Controller::GetNodes()` to get all found nodes. [ECHONETLite](http://www.echonet.gr.jp/english/index.htm) node can have multiple objects, use `Node::GetObjects()` to get the all objects in the node.

```
ctrl := echonet.NewController()
....
for _, node := range ctrl.GetNodes() {
    ....
    objs := node.GetObjects()
    for _, obj := range objs {
    ....
    }
}
```

### 4. Creating Control Message

To control the found objects, create the control message using `NewMessage()` as the following.

```
ctrl := echonet.NewController()
....
msg := echonet.NewMessage()
msg.SetDestinationObjectCode(0xXXXXXX)
msg.SetESV(0xXX)
...
prop := echonet.NewProperty()
prop.SetCode(echonet.PropertyCode(0x00))
prop.SetData(....)
msg.AddProperty((prop.toProtocolProperty())
```

To create the message, developer should only set the following message objects using the `Message::SetDestinationObjectCode()`, `SetESV()` and `AddProperty()` functions.

- DEOJ : Destination ECHONET Lite object specification
- ESV : ECHONET Lite service
- EPC : ECHONET Lite Property
- PDC : Property data counter
- EDT : Property value data

The `uecho-go` controller sets the following message objects automatically when the message is sent.

- EHD1 : ECHONET Lite message header 1
- EHD2 : ECHONET Lite message header 2
- TID￼  : Transaction ID
- SEOJ : Source ECHONET Lite object specification
- OPC  : Number of processing properties

### 5. Sending Messages

To send the created message, use `Controller::SendMessage()` as the following:

```
ctrl := echonet.NewController()
....
msg := echonet.NewMessage()
....
err := ctrl.SendMessage(dstNode, reqMsg)
if err != nil {
    ....
}
```

Basically, all messages of [ECHONETLite](http://www.echonet.gr.jp/english/index.htm) is async. To handle the async response of the message request, use `Controller::PostMessage()` as the following:

```
ctrl := echonet.NewController()
....
reqMsg := echonet.NewMessage()
....
resMsg, err := ctrl.PostMessage(dstNode, reqMsg)
if err != nil {
    ....
}
....
```

## Next Steps

Let's check the following documentation to know the controller functions of uEcho in more detail.

- [Inside of uEcho Controller](./controller_inside.md)
- [Examples of uEcho](./examples.md)

[enet]:http://echonet.jp/english/
