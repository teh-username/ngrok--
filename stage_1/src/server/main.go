package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

const defaultControlPort = ":60623"

func generateListenChannel(conn net.Conn) <-chan []byte {
	listener := make(chan []byte)
	go func() {
		for {
			data, err := bufio.NewReader(conn).ReadBytes('\n')
			if err != nil {
				listener <- nil
				log.Println(err)
				break
			}
			listener <- data
		}
	}()
	return listener
}

// Handles each request received by control connection
func handleControlConn(ctrlConn net.Conn) {
	defer func() {
		log.Printf("Connection Closed: %s", ctrlConn.RemoteAddr().String())
		ctrlConn.Close()
	}()
	ctrlChan := generateListenChannel(ctrlConn)
	ctr := 10

	for {
		log.Println("Waiting for data")
		select {
		case ctrlData := <-ctrlChan:
			if ctrlData == nil {
				return
			}
			if ctr > 0 {
				log.Printf("[%d] Message received: %s", ctr, string(ctrlData))
				_, err := ctrlConn.Write([]byte("pong\n"))
				if err != nil {
					log.Println(err)
					return
				}
			} else {
				return
			}
		}
		ctr--
	}
}

// Handles setup and service of the control connection initiated by the client
func listenAndServeControlConn() {
	ctrlListener, err := net.Listen("tcp", defaultControlPort)
	if err != nil {
		log.Println(err)
		return
	}
	defer ctrlListener.Close()
	log.Printf("Listening on control port: %s", defaultControlPort)

	for {
		ctrlConn, err := ctrlListener.Accept()
		if err != nil {
			log.Println(err)
			break
		}
		log.Printf("New connection accepted: %s", ctrlConn.RemoteAddr().String())
		go handleControlConn(ctrlConn)
	}
}

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	listenAndServeControlConn()
}
