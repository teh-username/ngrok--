package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

const defaultControlAddr = "localhost:60623"

func createConnListenChannel(conn net.Conn) <-chan []byte {
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

func init() {
	log.SetOutput(os.Stdout)
}

// Handles connection initiation and further correspondence
func runControlConn() {
	ctrlConn, err := net.Dial("tcp", defaultControlAddr)
	if err != nil {
		log.Println(err)
		return
	}
	defer ctrlConn.Close()
	ctrlChan := generateListenChannel(ctrlConn)

	// Initiate initial ping
	_, writeErr := ctrlConn.Write([]byte("ping\n"))
	if writeErr != nil {
		log.Println(writeErr)
		return
	}

	for {
		log.Println("Waiting for data")
		select {
		case ctrlData := <-ctrlChan:
			if ctrlData == nil {
				return
			}
			log.Printf("Message received: %s", string(ctrlData))
			_, pingErr := ctrlConn.Write([]byte("ping\n"))
			if pingErr != nil {
				log.Println(err)
				return
			}
		}
	}

}

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	runControlConn()
}
