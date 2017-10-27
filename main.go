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

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "update.html")
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

type Data struct {
	Message   string `json: "message"`
	Password  string `json: "password"`
}

func SetHandler(w http.ResponseWriter, r *http.Request) {
	test := Data{}
	err := json.NewDecoder(r.Body).Decode(&test)
	if(err != nil){
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("%v",test.Password)
	if(test.Password == "password"){
		msg = test.Message
		fmt.Println("good job")
	}
}


func main() {
	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("static"))
	r.HandleFunc("/", IndexHandler).Methods("GET")
	r.HandleFunc("/update", UpdateHandler).Methods("GET")
	r.HandleFunc("/update/submit/", SetHandler).Methods("POST")

	r.Handle("/static/app.js", http.StripPrefix("/static/", fs))
	r.Handle("/static/styles.css", http.StripPrefix("/static/", fs))
	r.HandleFunc("/ws", ReceiveClient)

	// Serve requests
	http.Handle("/", r)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Unable to ListenAndServe: %v", err)
	}
}
