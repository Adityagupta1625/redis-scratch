package main

import (
	"encoding/binary"
	"fmt"
	"syscall"
)

const kMaxMsg = 4096

// read_full ensures that exactly 'len' bytes are read from the file descriptor
// This function handles partial reads by continuing to read until all requested bytes are received
//
// Parameters:
//   - fd (int): File descriptor to read from
//   - buf ([]byte): Buffer to store the read data
//   - len (int): Number of bytes to read
//
// Returns:
//   - error: nil on success, error if read fails or EOF is encountered
//
// Example usage:
//   buffer := make([]byte, 1024)
//   err := read_full(connfd, buffer, 512)  // Read exactly 512 bytes
func read_full(fd int, buf []byte, len int) error {
	
	offset:=0

	for len > 0 {
		n, err:= syscall.Read(fd, buf[offset:len])

		if err != nil {
			return fmt.Errorf("Error reading from socket: %v", err)
		}

		if n==0{
			return fmt.Errorf("EOF reading from socket")
		}

		offset+=n
		len-=n
	}

	return nil
}

// write_full ensures that all bytes in the buffer are written to the file descriptor
// This function handles partial writes by continuing to write until all data is sent
//
// Parameters:
//   - fd (int): File descriptor to write to
//   - buf ([]byte): Buffer containing data to write
//
// Returns:
//   - error: nil on success, error if write fails or returns 0 unexpectedly
//
// Example usage:
//   message := []byte("Hello Redis")
//   err := write_full(connfd, message)  // Write entire message
func write_full(fd int, buf []byte) error{
	total:=len(buf)
	offset:=0

	for total > 0 {
		n, err:= syscall.Write(fd,buf[offset:total])

		if err!=nil{
			return fmt.Errorf("Error writing to socket: %v", err)
		}

		if n==0{
			return fmt.Errorf("write returned 0, unexpected")
		}

		offset+=n
		total-=n
	}

	return nil
}

// one_request handles a single client request using a length-prefixed protocol
// This function implements the Redis wire protocol for message framing
//
// Protocol format:
//   - 4 bytes: message length (little endian)
//   - N bytes: message body
//
// Parameters:
//   - connfd (int): File descriptor of the client connection
//
// Returns:
//   - error: nil on success, error if protocol violation or I/O error occurs
//
// Example usage:
//   connfd, _, err := syscall.Accept(serverfd)
//   if err == nil {
//       err = one_request(connfd)  // Process one client request
//   }
func one_request(connfd int) error {

	// Step 1: Read 4 byte length header
	rbuf:= make([]byte,4+kMaxMsg)
	err:= read_full(connfd,rbuf,4);

	if err!=nil{
		fmt.Println("Error reading from socket:", err)
		return err
	}

	// Step 2: Decode length (little endian)
	length:= binary.LittleEndian.Uint32(rbuf[:4])
	if length > kMaxMsg {
		fmt.Println("too long")
		return fmt.Errorf("message too long: %d", length)
	}

	// Step 3: Read message body
	err = read_full(connfd, rbuf[4:4+length],int(length)); 
	
	if err != nil {
		fmt.Println("read() error:", err)
		return err
	}

	fmt.Printf("client says: %s\n", string(rbuf[4:4+length]))

	// Step 4: Send response
	reply := []byte("Hello world!!")
	replyLen := uint32(len(reply))
	wbuf := make([]byte, 4+replyLen)
	binary.LittleEndian.PutUint32(wbuf[:4], replyLen)
	copy(wbuf[4:], reply)

	err = write_full(connfd, wbuf); 
	
	if err != nil {
		fmt.Println("write_all() error:", err)
		return err
	}

	return nil
}

// handleConnection processes incoming client connections by reading data and sending a response
// This function demonstrates basic socket I/O operations for a Redis-like server
//
// Parameters:
//   - connfd (int): File descriptor of the accepted client connection
//
// Example usage:
//   connfd, _, err := syscall.Accept(serverfd)
//   if err == nil {
//       handleConnection(connfd)
//   }
//
// The function performs the following operations:
// 1. Creates a 64-byte buffer to read incoming data
// 2. Reads data from the client socket
// 3. Prints the received data to console
// 4. Sends a "Hello World" response back to the client
func handleConnection(connfd int) {
	// make creates a slice with specified length and capacity
	// Parameters: make([]type, length, capacity)
	// Example: make([]byte, 64) creates a byte slice of length 64
	rbuf := make([]byte, 64)
	
	// syscall.Read reads data from a file descriptor into a buffer
	// Parameters:
	//   - fd (int): File descriptor to read from
	//   - buf ([]byte): Buffer to store the read data
	// Returns:
	//   - n (int): Number of bytes read
	//   - err (error): Error if any occurred during reading
	// Example: n, err := syscall.Read(connfd, buffer)
	n, err := syscall.Read(connfd, rbuf)

	if err != nil {
		fmt.Println("Error reading from socket:", err)
		return
	}
	
	// fmt.Println prints values to standard output with a newline
	// Parameters: fmt.Println(a ...interface{})
	// Example: fmt.Println("Message:", data)
	fmt.Println("Read", n, "bytes:", string(rbuf[:n]))

	// Convert string to byte slice for network transmission
	// []byte() converts string to byte slice
	// Example: []byte("Hello") creates [72 101 108 108 111]
	wbuf := []byte("Hello World")
	
	// syscall.Write writes data from a buffer to a file descriptor
	// Parameters:
	//   - fd (int): File descriptor to write to
	//   - buf ([]byte): Buffer containing data to write
	// Returns:
	//   - n (int): Number of bytes written
	//   - err (error): Error if any occurred during writing
	// Example: n, err := syscall.Write(connfd, []byte("response"))
	_, err = syscall.Write(connfd, wbuf)

	if err != nil {
		fmt.Println("Error writing to socket:", err)
		return
	}
}

// main function initializes and runs a basic Redis-like TCP server
// This demonstrates low-level socket programming using system calls
//
// The server performs the following operations:
// 1. Creates a TCP socket
// 2. Sets socket options for address reuse
// 3. Binds to localhost:8000
// 4. Listens for incoming connections
// 5. Accepts and handles client connections in a loop
//
// Example server lifecycle:
//   socket() -> setsockopt() -> bind() -> listen() -> accept() -> handle -> close()
func main() {
	// syscall.Socket creates a new socket and returns its file descriptor
	// This is a low-level interface to the operating system's socket API
	//
	// Parameters:
	//   - domain (int): Address family (syscall.AF_INET for IPv4, syscall.AF_INET6 for IPv6)
	//   - typ (int): Socket type (syscall.SOCK_STREAM for TCP, syscall.SOCK_DGRAM for UDP)
	//   - proto (int): Protocol (syscall.IPPROTO_TCP for TCP, syscall.IPPROTO_UDP for UDP)
	//
	// Returns:
	//   - fd (int): File descriptor of the created socket
	//   - err (error): Error if socket creation failed
	//
	// Example usage:
	//   fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	//   // Creates a TCP socket for IPv4 communication
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)

	if err != nil {
		// fmt.Printf formats and prints to standard output
		// Parameters: fmt.Printf(format string, a ...interface{})
		// Example: fmt.Printf("Error: %v\n", err)
		fmt.Printf("Error creating socket: %v\n", err)
		return
	}

	// defer schedules a function call to be run when the surrounding function returns
	// Parameters: defer function_call
	// Example: defer file.Close() - ensures file is closed when function exits
	defer syscall.Close(fd)

	// syscall.SetsockoptInt sets an integer socket option
	// This allows the socket to reuse the address immediately after closing
	//
	// Parameters:
	//   - fd (int): Socket file descriptor
	//   - level (int): Protocol level (syscall.SOL_SOCKET for socket-level options)
	//   - name (int): Option name (syscall.SO_REUSEADDR allows address reuse)
	//   - value (int): Option value (1 to enable, 0 to disable)
	//
	// Returns:
	//   - err (error): Error if setting the option failed
	//
	// Example usage:
	//   err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	//   // Enables address reuse to avoid "address already in use" errors
	err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)

	if err != nil {
		fmt.Printf("Error setting socket option: %v\n", err)
		return
	}

	// syscall.SockaddrInet4 represents an IPv4 socket address
	// Fields:
	//   - Port (int): Port number (8000 in this case)
	//   - Addr ([4]byte): IPv4 address as 4-byte array
	//
	// Example addresses:
	//   - [4]byte{127, 0, 0, 1} = localhost (127.0.0.1)
	//   - [4]byte{0, 0, 0, 0} = all interfaces (0.0.0.0)
	//   - [4]byte{192, 168, 1, 100} = 192.168.1.100
	addr := syscall.SockaddrInet4{
		Port: 8000,
		Addr: [4]byte{127, 0, 0, 1},
	}

	// syscall.Bind associates a socket with a specific address and port
	//
	// Parameters:
	//   - fd (int): Socket file descriptor
	//   - sa (syscall.Sockaddr): Socket address structure
	//
	// Returns:
	//   - err (error): Error if binding failed
	//
	// Example usage:
	//   addr := syscall.SockaddrInet4{Port: 8080, Addr: [4]byte{0, 0, 0, 0}}
	//   err := syscall.Bind(fd, &addr)
	//   // Binds socket to port 8080 on all interfaces
	err = syscall.Bind(fd, &addr)

	if err != nil {
		fmt.Printf("Error binding socket: %v\n", err)
		return
	}

	// syscall.Listen marks the socket as a passive socket for accepting connections
	//
	// Parameters:
	//   - fd (int): Socket file descriptor
	//   - backlog (int): Maximum number of pending connections (syscall.SOMAXCONN for system maximum)
	//
	// Returns:
	//   - err (error): Error if listen failed
	//
	// Example usage:
	//   err := syscall.Listen(fd, 10)  // Allow up to 10 pending connections
	//   err := syscall.Listen(fd, syscall.SOMAXCONN)  // Use system maximum
	err = syscall.Listen(fd, syscall.SOMAXCONN)

	if err != nil {
		fmt.Printf("Error listening on socket: %v\n", err)
		return
	}

	fmt.Println("Redis server listening on 127.0.0.1:8000")

	// Infinite loop to continuously accept and handle client connections
	for {
		// syscall.Accept accepts an incoming connection on a listening socket
		//
		// Parameters:
		//   - fd (int): Listening socket file descriptor
		//
		// Returns:
		//   - connfd (int): File descriptor for the new connection
		//   - sa (syscall.Sockaddr): Address of the connecting client
		//   - err (error): Error if accept failed
		//
		// Example usage:
		//   connfd, clientAddr, err := syscall.Accept(serverfd)
		//   // connfd is used to communicate with the specific client
		//   // clientAddr contains the client's IP and port information
		connfd, sa, err := syscall.Accept(fd)

		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			return
		}

		fmt.Printf("Accepted connection from %v\n", sa)

		// handleConnection(connfd)
		one_request(connfd)

		// syscall.Close closes a file descriptor
		// Parameters: syscall.Close(fd int)
		// Example: syscall.Close(connfd) - closes the client connection
		syscall.Close(connfd)
	}
}
