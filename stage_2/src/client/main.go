package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

const defaultControlAddr = "localhost:60623"
const defaultProxyAddr = "localhost:60625"

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

// Sets up and handles proxy connection traffic
func dialAndServeProxyTraffic() {
	proxyConn, err := net.Dial("tcp", defaultProxyAddr)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		proxyConn.Close()
		log.Println("Client-side proxy closed")
	}()
	log.Println("Proxy connection established")

	proxyDataChan := createConnListenChannel(proxyConn)
	for {
		select {
		case proxyData := <-proxyDataChan:
			if proxyData == nil || string(proxyData) == "stop\n" {
				return
			}
			log.Printf("Proxy data received: %s", string(proxyData))
			_, proxyConnErr := proxyConn.Write(proxyData)
			if proxyConnErr != nil {
				log.Println(proxyConnErr)
				return
			}
		}
	}
}

// Sets up and handles control connection traffic
func runControlConn() {
	ctrlConn, err := net.Dial("tcp", defaultControlAddr)
	if err != nil {
		log.Println(err)
		return
	}
	defer ctrlConn.Close()
	log.Printf("Control connection established")
	ctrlChan := createConnListenChannel(ctrlConn)

	for {
		select {
		case ctrlData := <-ctrlChan:
			if ctrlData == nil {
				return
			}
			if string(ctrlData) == "create_proxy\n" {
				go dialAndServeProxyTraffic()
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
	runControlConn()
}
