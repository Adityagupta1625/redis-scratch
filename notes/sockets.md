Network Layers

IP Layer (layer of small discrete messages):

- lowest layer deals with generating small data packets

The layer of Multiplexing (Port Number):

- layer above IP
- multiple apps share network with a same device so how they know to which packet belongs to which app. This is called demultiplexing.This layer adds a 16 bit number over the packet to differentiate between apps. Each app claim an unused local port where it can send or receive data. (src_ip, src_port, dst_ip, dst_port)

The layer of reliable & ordered bytes (TCP):

- TCP provides a layer of reliable & ordered bytes on top of IP packets, it handles retransmission, reordering automatically.

3 layers represent 3 needs in networking. They are mapped nicely to TCP/IP concepts. There are other models, such as the TCP/IP model:

Application -> Transport layer (TCP/UDP) -> IP layer -> Link layer (below IP)

Packet vs. stream

TCP provides a byte stream, but typical apps expect messages; few apps use the byte stream without interpreting it. Thus, we either need to add a message layer to TCP, or add reliability & order to UDP. The former is far easier, so most apps use TCP, either by using a well-known protocol on top of TCP, or by rolling their own protocol.

TCP and UDP are not only functionally different, their semantics are incompatible. TCP or UDP is the first thing to decide for networked applications.

Socket

A socket is a handle to refer to a connection or something else. The API for networking is called the socket API, which is similar on different operating systems. The name “socket” has nothing to do with sockets on the wall.

The socket() method allocates and returns a socket fd (handle), which is used later to actually create connections.

A handle must be closed when you’re done to free the associated resources on the OS side. This is the only thing in common between different types of handles.

Listening is telling the OS that an app is ready to accept TCP connections from a given port. The OS then returns a socket handle for apps to refer to that port. From the listening socket, apps can retrieve (accept) incoming TCP connections, which is also represented as a socket handle. So there are 2 types of handles: listening socket & connection socket.

Creating a listening socket requires at least 3 API calls:

Obtain a socket handle via socket().
Set the listening IP:port via bind().
Create the listening socket via listen().
Then use the accept() API to wait for incoming TCP connections

A connection socket is created from the client side with 2 API calls:

Obtain a socket handle via socket().
Create the connection socket via connect().

socket() creates a typeless socket; the socket type (listening or connection) is determined after the listen() or connect() call. The bind() between socket() and listen() merely sets a parameter.

Although TCP and UDP provide different types of services, they share the same socket API, including the send() and recv() methods. For message-based sockets (UDP), each send/recv corresponds to a single packet. For byte-stream-based sockets (TCP), each send/recv appends to/consumes from the byte stream.
