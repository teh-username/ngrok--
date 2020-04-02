package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

const defaultClientPort = "54286"

// Allows us to interleave blocking IOs by converting them into channels
func generateListenChannel(conn net.Conn) <-chan []byte {
	listener := make(chan []byte)
	go func() {
		buf := make([]byte, 1024)
		for {
			size, err := conn.Read(buf)
			if err != nil {
				listener <- nil
				log.Println(err)
				break
			}
			staging := make([]byte, size)
			copy(staging, buf[:size])
			listener <- staging
		}
	}()
	return listener
}

func handleConnection(sourceConn net.Conn, targetAddr string) {
	targetConn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		log.Println(err)
		return
	}
	defer sourceConn.Close()
	defer targetConn.Close()
	log.Printf("Connection to %s established", targetConn.RemoteAddr().String())

	sourceChan := generateListenChannel(sourceConn)
	targetChan := generateListenChannel(targetConn)

	for {
		select {
		case srcData := <-sourceChan:
			if srcData == nil {
				return
			} else {
				log.Printf("Forwarding data from %s to %s:\n %s",
					sourceConn.RemoteAddr().String(),
					targetConn.RemoteAddr().String(),
					string(srcData),
				)
				_, err := targetConn.Write(srcData)
				if err != nil {
					log.Println(err)
					break
				}
			}
		case tarData := <-targetChan:
			if tarData == nil {
				return
			} else {
				log.Printf("Forwarding data from %s to %s:\n %s",
					targetConn.RemoteAddr().String(),
					sourceConn.RemoteAddr().String(),
					string(tarData),
				)
				_, err := sourceConn.Write(tarData)
				if err != nil {
					log.Println(err)
					break
				}
			}
		}
	}
}

func clientListenAndServe(targetAddr string) {
	clientConn, err := net.Listen("tcp", fmt.Sprintf(":%s", defaultClientPort))
	if err != nil {
		log.Println(err)
		return
	}
	defer clientConn.Close()
	log.Printf("Client is now listening on: %s", clientConn.Addr().String())

	for {
		sourceConn, err := clientConn.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("Connection accepted on %s", sourceConn.RemoteAddr().String())
		go handleConnection(sourceConn, targetAddr)
	}
}

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	args := os.Args
	if len(args) == 1 {
		log.Println("Please specify address to tunnel. Example: golang.org:80 | localhost:8080")
		return
	}

	clientListenAndServe(args[1])
}
