package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const controlPort = "60623"
const proxyPort = "60625"

func createCtrlListenChan(conn net.Conn) <-chan []byte {
	listener := make(chan []byte)
	go func() {
		for {
			data, err := bufio.NewReader(conn).ReadBytes('\n')
			if err != nil {
				listener <- nil
				log.Printf("[Control]: %s", err)
				break
			}
			listener <- data
		}
	}()
	return listener
}

func createGenericConnListenChan(conn net.Conn, connName string) <-chan []byte {
	listener := make(chan []byte)
	go func() {
		buf := make([]byte, 1024)
		for {
			size, err := conn.Read(buf)
			if err != nil {
				listener <- nil
				log.Printf("[%s]: %s", connName, err)
				break
			}
			staging := make([]byte, size)
			copy(staging, buf[:size])
			listener <- staging
		}
	}()
	return listener
}

// Sets up and handles proxy and private connection
func orchestrateProxyAndPrivateTraffic(privateAddr, serverAddr string) {
	proxyConn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", serverAddr, proxyPort))
	if err != nil {
		log.Println(err)
		return
	}
	defer proxyConn.Close()
	privateConn, privConnErr := net.Dial("tcp", privateAddr)
	if privConnErr != nil {
		log.Println(privConnErr)
		return
	}
	defer privateConn.Close()
	defer log.Println("Client-side proxy and private connections closed")
	log.Println("Proxy and private connections established")

	proxyDataChan := createGenericConnListenChan(proxyConn, "proxy")
	privateDataChan := createGenericConnListenChan(privateConn, "private")
	for {
		select {
		case proxyData := <-proxyDataChan:
			if proxyData == nil {
				return
			}
			_, privateConnWriteErr := privateConn.Write(proxyData)
			if privateConnWriteErr != nil {
				log.Println(privateConnWriteErr)
				return
			}
		case privateData := <-privateDataChan:
			if privateData == nil {
				return
			}
			_, proxyConnWriteErr := proxyConn.Write(privateData)
			if proxyConnWriteErr != nil {
				log.Println(proxyConnWriteErr)
				return
			}
		}
	}
}

// Sets up and handles control connection traffic
func runControlConn(privateAddr, serverAddr string) {
	ctrlConn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", serverAddr, controlPort))
	if err != nil {
		log.Println(err)
		return
	}
	defer ctrlConn.Close()
	log.Printf("Control connection established")
	ctrlChan := createCtrlListenChan(ctrlConn)

	for {
		select {
		case ctrlData := <-ctrlChan:
			if ctrlData == nil {
				return
			}
			if string(ctrlData) == "create_proxy\n" {
				go orchestrateProxyAndPrivateTraffic(privateAddr, serverAddr)
			} else {
				log.Printf("Server sent unknown command: %d", ctrlData)
			}
		}
	}

}

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	args := os.Args
	serverAddr := "localhost"
	if len(args) == 1 {
		log.Println("Please specify address to tunnel. Example: golang.org:80 | localhost:8080")
		return
	} else if len(args) == 2 {
		log.Println("Server defaulted to localhost. To specify server address, provide as second argument. Example: localhost:8080 127.0.0.1")
	} else if len(args) == 3 {
		log.Printf("Setting %s as server address", args[2])
		serverAddr = args[2]
	}

	runControlConn(args[1], serverAddr)
}
