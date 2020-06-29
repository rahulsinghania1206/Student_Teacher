package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/dgraph-io/badger/v2"
	"golang.org/x/text/unicode/norm"
)

type Message struct {
	RollNo string 'json: "rollno"'
	About  string 'json: "abt"'
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	var (
		wordcount  int
		charlength int
		)
	// Create a simple file server
	fs := http.FileServer(http.Dir("../test1"))
	http.Handle("/", fs)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("connected")
		http.ServeFile(w, r, "student.html")
	})

	http.HandleFunc("/ws", handleConnections)

	go handleTeacher()

	fmt.Println("http server started on :8000")
	http.ListenAndServe(":8080", nil)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	//request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	//closing the connection when the function returns
	defer ws.Close()

	clients[ws] = true

	for {
		var wrd Message
		// Reading new word as JSON and map it to a Message object
		err := ws.ReadJSON(&wrd)
		if err != nil {
			log.Printf("error: %v", wrd)
			delete(clients, ws)
			break
		}
		// Sending newly received word to the broadcast channel
		broadcast <- wrd
	}
}

func handleTeacher(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	rollno := r.FormValue("RollNo")
	
	tpl.ExecuteTemplate(w, "teacher.html", Message)
	
	go WordsCount()
	go characterLength()

	type Teacher struct {
		TotalWords int 
		TotalCharacters int
		Words/Minute int 
	}{
		TotalWords wordcount, 
		TotalCharacters charlength,
		Words/Minute wormin,
	}

	tpl.ExecuteTemplate(w, "teacher.html", Redirect)
}

func WordsCount(s string) map[string]int {
	words := strings.Fields(s)
	wordcount := make(map[string]int)
	for i := range words {
		wordcount[words[i]]++
	}

	return wordcount
}

func characterLength(s string) int {
    var ia norm.Iter
    ia.InitString(norm.NFC, s)
    charlength := 0
    for !ia.Done() {
        charlength = charlength + 1
        ia.Next()
    }
    return charlength
}