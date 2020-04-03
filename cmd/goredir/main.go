package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

var (
	target string
	host   string
	port   int
)

func init() {
	flag.StringVar(&target, "target", "", "target to forward connections to (<host>:<port>)")
	flag.StringVar(&host, "host", "", "host to listen for tcp connections on")
	flag.IntVar(&port, "port", 5555, "port to listen for tcp connections on")
	flag.Parse()
	if target == "" {
		log.Fatal("-target must be provided")
	}
}

func main() {
	incoming, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalf("Could not start server on %d: %v", port, err)
	}
	fmt.Printf("Listening for connections on :%d and forwarding to %s\n", port, target)

	for {
		client, err := incoming.Accept()
		if err != nil {
			log.Println("Error while accepting incoming client connection:", err)
			continue
		}
		go handleConnection(client)
	}
}

func handleConnection(client net.Conn) {
	fmt.Printf("Received connection from client at %v\n", client.RemoteAddr())
	target, err := net.Dial("tcp", target)
	if err != nil {
		log.Println("Could not connect to target server, dropping connection:", err)
		client.Close()
		return
	}
	fmt.Printf("Connection to server %v established, proxying data\n", target.RemoteAddr())
	go copyIO(client, target)
	go copyIO(target, client)
}

func copyIO(src, dest net.Conn) {
	defer src.Close()
	defer dest.Close()
	if _, err := io.Copy(src, dest); err != nil {
		if !strings.Contains(strings.ToLower(err.Error()), "closed") {
			fmt.Println("Error copying stream:", err)
		}
	}
}
