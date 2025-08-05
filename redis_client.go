package main

import (
	"fmt"
	"syscall"
)

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
	wbuf := []byte("Hello World")
	
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
	_, err = syscall.Write(fd, wbuf)

	if err != nil {
		fmt.Println("Error writing to socket:", err)
		return
	}

	// make creates a slice with specified length and capacity
	// Parameters: make([]type, length, capacity)
	// Examples:
	//   - make([]byte, 1024) creates 1KB buffer for large responses
	//   - make([]byte, 64) creates 64-byte buffer for small responses
	//   - make([]byte, 4096) creates 4KB buffer for bulk data
	rbuf := make([]byte, 64)
	
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
	n, err := syscall.Read(fd, rbuf)

	if err != nil {
		fmt.Println("Error reading from socket:", err)
		return
	}
	
	// string() converts byte slice to string
	// rbuf[:n] creates a slice from index 0 to n (excluding n)
	// Examples:
	//   - string([]byte{72, 101, 108, 108, 111}) = "Hello"
	//   - string(rbuf[:n]) converts only the bytes that were actually read
	fmt.Println("Read", n, "bytes:", string(rbuf[:n]))

	// syscall.Close closes the socket connection
	// Parameters: syscall.Close(fd int)
	// Examples:
	//   - syscall.Close(fd) - closes client socket
	//   - defer syscall.Close(fd) - ensures socket is closed when function exits
	syscall.Close(fd)
	fmt.Println("Connection closed")
}
