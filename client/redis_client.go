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
//   message := []byte("PING")
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

// query sends a command to the Redis server and reads the response
// This function implements the client side of the length-prefixed protocol
//
// Protocol format:
//   - Send: 4 bytes length + message body
//   - Receive: 4 bytes length + response body
//
// Parameters:
//   - fd (int): File descriptor of the server connection
//   - text (string): Command text to send to the server
//
// Returns:
//   - error: nil on success, error if message too long, I/O error, or protocol violation
//
// Example usage:
//   err := query(serverfd, "PING")        // Send PING command
//   err := query(serverfd, "GET mykey")   // Send GET command
func query(fd int, text string) error {

	length := len(text)
	if length > kMaxMsg {
		return fmt.Errorf("message too long")
	}

	// Prepare write buffer (length-prefixed message)
	wbuf := make([]byte, 4+length)
	binary.LittleEndian.PutUint32(wbuf[:4], uint32(length))
	copy(wbuf[4:], text)

	// Write full request
	err := write_full(fd, wbuf); 
	
	if err != nil {
		return fmt.Errorf("write_all error: %v", err)
	}

	// Read 4-byte response header
	rbuf := make([]byte, 4+kMaxMsg+1)

	err = read_full(fd, rbuf[:4],4); 
	
	if err != nil {
		return err
	}

	// Parse reply length
	replyLen := binary.LittleEndian.Uint32(rbuf[:4])
	if replyLen > kMaxMsg {
		fmt.Println("too long")
		return fmt.Errorf("response too long: %d", replyLen)
	}

	// Read response body
	err = read_full(fd, rbuf[4:4+replyLen],int(replyLen)); 
	
	if err != nil {
		fmt.Println("read() error:", err)
		return err
	}

	// Print response
	fmt.Printf("server says: %s\n", string(rbuf[4:4+replyLen]))
	return nil
}

// main function implements a basic Redis client that connects to a server
// This demonstrates client-side socket programming using system calls
//
// The client performs the following operations:
// 1. Creates a TCP socket
// 2. Connects to the Redis server at localhost:8000
// 3. Sends a "Hello World" message
// 4. Reads the server's response
// 5. Closes the connection
//
// Example client lifecycle:
//   socket() -> connect() -> write() -> read() -> close()
func main() {
	// syscall.Socket creates a new socket and returns its file descriptor
	// This creates a TCP socket for IPv4 communication
	//
	// Parameters:
	//   - domain (int): Address family (syscall.AF_INET for IPv4)
	//   - typ (int): Socket type (syscall.SOCK_STREAM for TCP)
	//   - proto (int): Protocol (syscall.IPPROTO_TCP for TCP)
	//
	// Returns:
	//   - fd (int): File descriptor of the created socket
	//   - err (error): Error if socket creation failed
	//
	// Example usage:
	//   fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	//   // Creates a TCP socket for client communication
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	
	if err != nil {
		// fmt.Printf formats and prints to standard output
		// Parameters: fmt.Printf(format string, a ...interface{})
		// Example: fmt.Printf("Connection failed: %v\n", err)
		fmt.Printf("Error creating socket: %v\n", err)
		return
	}

	// syscall.SockaddrInet4 represents an IPv4 socket address for the server
	// Fields:
	//   - Port (int): Server port number (8000)
	//   - Addr ([4]byte): Server IPv4 address as 4-byte array
	//
	// Common server addresses:
	//   - [4]byte{127, 0, 0, 1} = localhost (127.0.0.1)
	//   - [4]byte{192, 168, 1, 10} = 192.168.1.10
	//   - [4]byte{10, 0, 0, 1} = 10.0.0.1
	addr := syscall.SockaddrInet4{
		Port: 8000,
		Addr: [4]byte{127, 0, 0, 1},
	}

	// syscall.Connect establishes a connection to the server
	//
	// Parameters:
	//   - fd (int): Client socket file descriptor
	//   - sa (syscall.Sockaddr): Server address structure
	//
	// Returns:
	//   - err (error): Error if connection failed
	//
	// Example usage:
	//   serverAddr := syscall.SockaddrInet4{Port: 9000, Addr: [4]byte{192, 168, 1, 100}}
	//   err := syscall.Connect(clientfd, &serverAddr)
	//   // Connects to server at 192.168.1.100:9000
	err = syscall.Connect(fd, &addr)

	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		return
	}

	// fmt.Println prints values to standard output with a newline
	// Parameters: fmt.Println(a ...interface{})
	// Example: fmt.Println("Status:", "Connected", "Port:", 8000)
	fmt.Println("Connected to Redis server at 127.0.0.1:8000")

	// Convert string to byte slice for network transmission
	// []byte() converts string to byte slice
	// Examples:
	//   - []byte("GET key1") creates [71 69 84 32 107 101 121 49]
	//   - []byte("SET key value") for Redis SET command
	//   - []byte("PING") for Redis PING command
	// wbuf := []byte("Hello World")
	
	// syscall.Write sends data to the server through the socket
	//
	// Parameters:
	//   - fd (int): Socket file descriptor
	//   - buf ([]byte): Buffer containing data to send
	//
	// Returns:
	//   - n (int): Number of bytes written
	//   - err (error): Error if writing failed
	//
	// Example usage:
	//   message := []byte("PING")
	//   n, err := syscall.Write(fd, message)
	//   // Sends PING command to Redis server
	// _, err = syscall.Write(fd, wbuf)

	// if err != nil {
	// 	fmt.Println("Error writing to socket:", err)
	// 	return
	// }

	// make creates a slice with specified length and capacity
	// Parameters: make([]type, length, capacity)
	// Examples:
	//   - make([]byte, 1024) creates 1KB buffer for large responses
	//   - make([]byte, 64) creates 64-byte buffer for small responses
	//   - make([]byte, 4096) creates 4KB buffer for bulk data
	// rbuf := make([]byte, 64)
	
	// syscall.Read receives data from the server through the socket
	//
	// Parameters:
	//   - fd (int): Socket file descriptor
	//   - buf ([]byte): Buffer to store received data
	//
	// Returns:
	//   - n (int): Number of bytes read
	//   - err (error): Error if reading failed
	//
	// Example usage:
	//   buffer := make([]byte, 1024)
	//   n, err := syscall.Read(fd, buffer)
	//   response := string(buffer[:n])  // Convert bytes to string
	// n, err := syscall.Read(fd, rbuf)

	// if err != nil {
	// 	fmt.Println("Error reading from socket:", err)
	// 	return
	// }
	
	// string() converts byte slice to string
	// rbuf[:n] creates a slice from index 0 to n (excluding n)
	// Examples:
	//   - string([]byte{72, 101, 108, 108, 111}) = "Hello"
	//   - string(rbuf[:n]) converts only the bytes that were actually read
	// fmt.Println("Read", n, "bytes:", string(rbuf[:n]))


	query(fd, "PING")

	// syscall.Close closes the socket connection
	// Parameters: syscall.Close(fd int)
	// Examples:
	//   - syscall.Close(fd) - closes client socket
	//   - defer syscall.Close(fd) - ensures socket is closed when function exits
	syscall.Close(fd)
	fmt.Println("Connection closed")
}
