package main

import (
	"encoding/json"
	"net/http"
	"time"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var clients []*websocket.Conn

var msg string

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func sendUpdate(fromClient *websocket.Conn) {
	for i, client := range clients {
		err := client.WriteMessage(1, []byte(msg))

		if err != nil {
			clients = append(clients[:i], clients[i+1:]...)
			log.Println(err)
			continue
		}
	}

}

func closeClient(c *websocket.Conn) {
	log.Println("close requested")
	c.Close()
}

func addClient(c *websocket.Conn) {
	clients = append(clients, c)
	go processMessages(c)
}

func processMessages(c *websocket.Conn) {
	defer closeClient(c)
		sendUpdate(c)
	for {
		sendUpdate(c)
		time.Sleep(1000 * time.Millisecond)

	}
}

func ReceiveClient(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	} else {
		go addClient(c)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func WriteJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	w.Write(b)
	return nil
}

func SetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v",r)
}

func updateMessage(){
	messages := [...]string{
		"Ooowhee",
		"It's a scorcher out there"}
	i := 0
	for {
		time.Sleep(1000 * time.Millisecond)
		i ++;
		if(i >= len(messages)){
			i = 0;
		}
		msg = messages[i]
	}
}

func main() {
	go updateMessage()
	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("static"))
	r.HandleFunc("/", IndexHandler).Methods("GET")

	r.Handle("/static/app.js", http.StripPrefix("/static/", fs))
	r.Handle("/static/styles.css", http.StripPrefix("/static/", fs))
	r.HandleFunc("/ws", ReceiveClient)

	// Serve requests
	http.Handle("/", r)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Unable to ListenAndServe: %v", err)
	}
}
