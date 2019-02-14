package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"golang.org/x/net/websocket"
)

type Message struct {
	Text string `json:"text"`
	//Ip   string `json:"text"`
}

var (
	port = flag.String("port", "8000", "puerto usado")
)

func main() {
	//
	fmt.Println("conectando...")
	//
	flag.Parse()

	// connect

	ws, err := connect()
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()

	// receive
	var m Message

	go func() {
		for {
			err := websocket.JSON.Receive(ws, &m)
			if err != nil {
				fmt.Println("Error para recibir message: ", err.Error())
				break
			}

			//log.Println(ws)
			fmt.Println("Mensaje:  ", m)
		}
	}()

	// send
	//variable para guardar la direccion ip
	var direcip string
	addrs, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}
	if ipnet, ok := addrs.LocalAddr().(*net.UDPAddr); ok && !ipnet.IP.IsLoopback() {
		if ipnet.IP.To4() != nil {
			//os.Stdout.WriteString(ipnet.IP.String() + "\n")
			direcip = ipnet.IP.String()
		}
	}
	fmt.Println("conectado")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		rest := direcip + ": " + text
		//direcip = scanner.Text()
		if text == "" {
			continue
		}
		m := Message{
			Text: rest,
		}

		err = websocket.JSON.Send(ws, m)

		if err != nil {
			fmt.Println("Error al enviar mensage: ", err.Error())
			break
		}
	}
}

// connect connects to the local chat server at port <port>
func connect() (*websocket.Conn, error) {
	return websocket.Dial(fmt.Sprintf("ws://192.168.13.3:%s", *port), "", mockedIP())
}

// mockedIP is a demo-only utility that generates a random IP address for this client
func mockedIP() string {
	var arr [4]int
	for i := 0; i < 4; i++ {
		rand.Seed(time.Now().UnixNano())
		arr[i] = rand.Intn(256)
	}
	return fmt.Sprintf("http://%d.%d.%d.%d", arr[0], arr[1], arr[2], arr[3])
}
