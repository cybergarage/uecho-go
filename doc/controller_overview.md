![logo](img/logo.png)

# Overview of uEcho Controller

The controller is a special node of [ECHONET Lite][enet] that controls other nodes. It can discover nodes in the local area network and send messages to the found devices.

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

Next, use `Controller::Search()` to search for other nodes in the local area network as follows:

```
ctrl := echonet.NewController()
....
err = ctrl.Search(context.Background())
if err != nil {
    ....
}
```

### 3. Getting Nodes and Objects

After searching, use `Controller::GetNodes()` to get all found nodes. An [ECHONET Lite](http://www.echonet.gr.jp/english/index.htm) node can have multiple objects; use `Node::Objects()` to get all objects in the node.

```
ctrl := echonet.NewController()
....
for _, node := range ctrl.GetNodes() {
    ....
    objs := node.Objects()
    for _, obj := range objs {
    ....
    }
}
```

### 4. Creating Control Message

To control the found objects, create a control message using `NewMessage()` as follows:

```
ctrl := echonet.NewController()
....
msg := echonet.NewMessage()
msg.SetDEOJ(0xXXXXXX)
msg.SetESV(0xXX)
...
prop := echonet.NewProperty()
prop.SetCode(echonet.PropertyCode(0x00))
prop.SetData(....)
msg.AddProperty((prop.toProtocolProperty())
```

To create the message, developers should only set the following message objects using the `Message::SetDEOJ()`, `SetESV()`, and `AddProperty()` functions.

- DEOJ : Destination ECHONET Lite object specification
- ESV : ECHONET Lite service
- EPC : ECHONET Lite Property
- PDC : Property data counter
- EDT : Property value data

The `uecho-go` controller sets the following message objects automatically when the message is sent.

- EHD1 : ECHONET Lite message header 1
- EHD2 : ECHONET Lite message header 2
- TID : Transaction ID
- SEOJ : Source ECHONET Lite object specification
- OPC  : Number of processing properties

### 5. Sending Messages

To send the created message, use `Controller::SendMessage()` as follows:

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

Basically, all messages in [ECHONET Lite](http://www.echonet.gr.jp/english/index.htm) are asynchronous. To handle the asynchronous response of the message request, use `Controller::PostMessage()` as follows:

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
