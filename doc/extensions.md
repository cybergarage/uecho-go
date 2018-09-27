![logo](img/logo.png)

# Extension Specifications

`uecho-go` supports some extension specifications for [ECHONET Lite][enet] to use the IoT protocol in a variety of situations such as Cloud computing or testing.


[enet]:http://echonet.jp/english/

# TCP Unicast

TCP unicast is not mandatory in [ECHONET Lite][enet], but `uecho-go` supports it the unicast protocol to send the messages more reliably. The TCP extention is disabled as default, and so use `Node::SetTCPEnabled() to enable it as the following.


```
node := NewLocalNode()
node.SetTCPEnabled(true)
```

# Automatic Port Binding

 An [ECHONET Lite][enet] node must listen the UDP unicast, UDP multicast and TCP unicast packets always at port number 3610, but `echo-go` supports automatic port mapping to bind at any port to be able to run the [ECHONET Lite][enet] nodes in the same machine at the same time. The extention is also disabled as default, and so use `Node::SetAutoBindingEnabled() to enable it as the following.

```
node := NewLocalNode()
node.SetAutoBindingEnabled(true)
 ```

 The automatic function binds to the specified port, 3610, for UDP multicast, but it searches an available UDP and TCP unicast ports to bind when the default port, 3610, is bound. Use `Node::GetPort()` to know the bound port after the node is started.
 
```
err := node.Start()
if err != nil {
    ...
}
boundPort := node.GetPort()
```

The bound UDP and TCP unicast ports are the same number. 
In addition, [ECHONET Lite][enet] does not specify the source port numbers of UDP multicast, UDP and TCP unicast, but `uecho-go` uses the bound port as the source port number for the all messaging.