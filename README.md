# pigeon
Pigeon is a simple message queue service written in Go, intended to be used in prototyping distributed services.

Currently, there are three parts to pigeon - the broker, the subscriber, and the publisher.

The broker handles incoming TCP connections from publishers, and routes them to the appropriate subscribers via websockets. 

The subscriber connects to the broker via websockets, and registers which topic it subscribes to. Then it waits to consume messages from the broker.

The publisher creats a TCP connection to the broker, and sends a message. It currently does not wait for acknowledgement of the message.

# Examples
A sample implementation exists in the CMD folder structure.
