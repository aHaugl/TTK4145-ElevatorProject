# Network module for Go (UDP broadcast only)

This network module is written by [TTK4145](https://github.com/TTK4145)
and can be found [here](https://github.com/TTK4145/Network-go). We've modified to code to suit our needs a little bit, but it is mostly the same as the original module.


Features
--------

Channel-in/channel-out pairs of (almost) any custom or built-in datatype can be supplied to a pair of transmitter/receiver functions. Data sent to the transmitter function is automatically serialized and broadcasted on the specified port. Any messages received on the receiver's port are deserialized (as long as they match any of the receiver's supplied channel datatypes) and sent on the corresponding channel. See [bcast.Transmitter and bcast.Receiver](network/bcast/bcast.go).

Peers on the local network can be detected by supplying your own ID to a transmitter and receiving peer updates (new, current and lost peers) from the receiver. See [peers.Transmitter and peers.Receiver](network/peers/peers.go).

Finding your own local IP address can be done with the [LocalIP](network/localip/localip.go) convenience function, but only when you are connected to the internet.


*Note: This network module does not work on Windows: Creating proper broadcast sockets on Windows is just not implemented yet in the Go libraries. See issues listed in [the comments here](network/conn/bcast_conn.go) for details.*
