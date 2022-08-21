![logo](img/logo.png)

# Expanded Configurations for Go Platform

`uecho-go` supports some extension specifications for [ECHONET Lite][enet] to use the IoT protocol in a variety of situations such as Cloud computing or testing.


# TCP Unicast

[ECHONET Lite][enet] does not strictly specify the TCP specification, and the TCP unicast is not mandatory. However `uecho-go` supports it the unicast protocol to send the messages more reliably. The TCP extention is disabled as default, and so use `Node::SetTCPEnabled() to enable it as the following.


```
node := NewLocalNode()
node.SetTCPEnabled(true)
```

According to [ECHONET Lite System Design Guidelines][enet_guideline_tcp], `uecho-go` tries to send any request messages using TCP connection at first when the option is enabled and send the request messages using UDP connection again when the TCP requests are failed. 
In addition, `uecho-go` returns all response messages using UDP connection when the request massages are received from UDP or multicast connection.

# Automatic Port Binding

 An [ECHONET Lite][enet] node must listen the UDP unicast, UDP multicast and TCP unicast packets always at port number 3610, but `echo-go` supports automatic port mapping to bind at any port to be able to run the [ECHONET Lite][enet] nodes in the same machine at the same time. The extention is also disabled as default, and so use `Node::SetAutoPortBindingEnabled() to enable it as the following.

```
node := NewLocalNode()
node.SetAutoPortBindingEnabled(true)
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

# References

- [Part V ECHONET Lite System Design Guidelines v1.12 : Chapter 5 - Guidelines on TCP][enet_guideline_tcp]

[enet]:http://echonet.jp/english/
[enet_guideline_tcp]:https://echonet.jp/spec_en/
