package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

const defaultControlPort = ":60623"
const defaultPublicPort = ":60624"
const defaultProxyPort = ":60625"

func createConnListenChannel(conn net.Conn, name string) <-chan []byte {
	listener := make(chan []byte)
	go func() {
		for {
			data, err := bufio.NewReader(conn).ReadBytes('\n')
			if err != nil {
				listener <- nil
				log.Printf("[%s] %s", name, err)
				break
			}
			listener <- data
		}
	}()
	return listener
}

// Pipes traffic between the public and proxy connection
func handlePubTraffic(pubConn, proxyConn net.Conn) {
	defer func() {
		pubConn.Close()
		proxyConn.Close()
		log.Println("Public and Proxy connection closed")
	}()

	pubDataChan := createConnListenChannel(pubConn, "public")
	proxyDataChan := createConnListenChannel(proxyConn, "proxy")

	for {
		select {
		case pubData := <-pubDataChan:
			if pubData == nil {
				return
			}
			_, proxyConnErr := proxyConn.Write(pubData)
			if proxyConnErr != nil {
				log.Println(proxyConnErr)
				return
			}
		case proxyData := <-proxyDataChan:
			if proxyData == nil {
				return
			}
			_, pubConnErr := pubConn.Write(proxyData)
			if pubConnErr != nil {
				log.Println(pubConnErr)
				return
			}
		}
	}
}

// Handles the synchronization of two asynchronously created connections
func handleProxyOrchestration(
	ctrlConn net.Conn,
	pubConnChan,
	proxyConnChan <-chan net.Conn,
) {
	for {
		var pubConn, proxyConn net.Conn

		select {
		case pubConn = <-pubConnChan:
			log.Println("Requesting proxy creation")
		}

		_, err := ctrlConn.Write([]byte("create_proxy\n"))
		if err != nil {
			// Client will need to redial
			pubConn.Close()
			ctrlConn.Close()
			log.Println(err)
			return
		}

		select {
		case proxyConn = <-proxyConnChan:
			go handlePubTraffic(pubConn, proxyConn)
		}
	}
}

// Sets up the control listener. Spawns a goroutine on successful connection
func listenAndServeControlConn(pubConnChan, proxyConnChan <-chan net.Conn) {
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
		go handleProxyOrchestration(ctrlConn, pubConnChan, proxyConnChan)
	}
}

// Encapsulates connection accept events under channels
func createConnectionChan(svcPort, connType string) <-chan net.Conn {
	connChan := make(chan net.Conn)
	go func() {
		listener, err := net.Listen("tcp", svcPort)
		if err != nil {
			log.Println(err)
			return
		}
		defer listener.Close()
		log.Printf("Listening on %s port: %s", connType, svcPort)

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println(err)
				break
			}
			log.Printf("New %s conn accepted: %s", connType, conn.RemoteAddr().String())
			connChan <- conn
		}
	}()
	return connChan
}

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	pubConnChan := createConnectionChan(defaultPublicPort, "public")
	proxyConnChan := createConnectionChan(defaultProxyPort, "proxy")
	listenAndServeControlConn(
		pubConnChan,
		proxyConnChan,
	)
}
