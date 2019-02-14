package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

//Struct que contendra el ms
type Message struct {
	Text string `json:"text"`
}

//Strcut que con las funciones que se ocupan en el servidor
type hub struct {
	clients          map[string]*websocket.Conn
	addClientChan    chan *websocket.Conn
	removeClientChan chan *websocket.Conn
	broadcastChan    chan Message
}

//Variable estatica para el puerto que se conectara el servidor
var (
	port = flag.String("port", "8000", "puerto usado por ws connection")
)

func main() {
	flag.Parse()
	//Verifica si la conexion es correcta
	log.Fatal(server(*port))

}

// Funcion que permite crear el servidor y ejecutar el websocket
func server(port string) error {
	h := newHub()
	mux := http.NewServeMux()
	mux.Handle("/", websocket.Handler(func(ws *websocket.Conn) {
		handler(ws, h)
	}))
	s := http.Server{Addr: "192.168.13.3" + ":" + port, Handler: mux}
	return s.ListenAndServe()
}

//Función handler encargada de ejecutar la go rutina h y verifica los mensajes para los clientes
func handler(ws *websocket.Conn, h *hub) {
	go h.run()
	h.addClientChan <- ws

	for {
		var m Message
		err := websocket.JSON.Receive(ws, &m)
		if err != nil {
			h.broadcastChan <- Message{err.Error()}
			h.removeClient(ws)
			return
		}
		h.broadcastChan <- m
	}
}

//NewHub retorna las funciones que ocupa el servidor para interactuar
func newHub() *hub {
	return &hub{
		clients:          make(map[string]*websocket.Conn),
		addClientChan:    make(chan *websocket.Conn),
		removeClientChan: make(chan *websocket.Conn),
		broadcastChan:    make(chan Message),
	}
}

//run selecciona que debe de hacer de acuerdo a lo recibido de la funcion
func (h *hub) run() {
	for {
		select {
		case conn := <-h.addClientChan:
			h.addClient(conn)
		case conn := <-h.removeClientChan:
			h.removeClient(conn)
		case m := <-h.broadcastChan:
			h.broadcastMessage(m)
		}
	}
}

//Quita los clientes del servidor
func (h *hub) removeClient(conn *websocket.Conn) {
	delete(h.clients, conn.LocalAddr().String())
}

//Agrega un nuevo cliente al servidor
func (h *hub) addClient(conn *websocket.Conn) {
	h.clients[conn.RemoteAddr().String()] = conn
}

//Verifica a los clientes conectados y si un clente se desconecta manda ms que el cliente no puede recibir ms
func (h *hub) broadcastMessage(m Message) {
	for _, conn := range h.clients {
		err := websocket.JSON.Send(conn, m)
		if err != nil {
			fmt.Println("Error de message: ", err)
			return
		}
	}

}
