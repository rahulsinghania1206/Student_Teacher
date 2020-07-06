package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	bolt "go.etcd.io/bbolt"
	"golang.org/x/net/websocket"
)

var (
	port = flag.Int("port", 4050, "The server port")
)
var rollNo, about string

type server struct {
	dbConnection *bolt.DB
	teacherConn  *websocket.Conn
}

type Message struct {
	X string `json:"rollNo"`
	Y string `json:"about"`
}

type checkRollNoStruct struct {
	T string `json:"checkRollNo"`
}

func handleWebsocketEchoMessage(msg Message, s *server) error {
	rollNo := (msg.X)
	about := (msg.Y)
	s.Put("student", rollNo, about)
	result := map[string]map[string]string{}
	result[rollNo] = map[string]string{}
	result[rollNo]["rollNo"] = rollNo
	result[rollNo]["charCount"] = strconv.Itoa(len(about))
	result[rollNo]["wordCount"] = strconv.Itoa(wordCount(string(about)))
	result[rollNo]["actualMessage"] = string(about)
	bb, _ := json.Marshal(result)
	sendDataToTeacher(s.teacherConn, string(bb))
	return nil
}

func websocketStudentConnection(ws *websocket.Conn, s *server) {
	log.Printf("Client connected from %s", ws.RemoteAddr())
	for {
		var msg Message
		err := websocket.JSON.Receive(ws, &msg)
		if err != nil {
			log.Printf("Receive failed: %s; closing connection...", err.Error())
			if err = ws.Close(); err != nil {
				log.Println("Error closing connection:", err.Error())
			}
			break
		} else {
			if err := handleWebsocketEchoMessage(msg, s); err != nil {
				log.Println(err.Error())
				break
			}
		}
	}
}

func newServer(filename string) (s *server, err error) {
	s = &server{}
	s.dbConnection, err = bolt.Open(filename, 0600, &bolt.Options{Timeout: 1 * time.Second})
	return
}

func (s *server) Put(bucket string, key string, val string) error {
	return s.dbConnection.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		return b.Put([]byte(key), []byte(val))
		//return b.Put([]byte("11"), []byte("ert"))
	})
}

func (s *server) Get(bucket string) (data map[string]map[string]string, err error) {
	s.dbConnection.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		c := b.Cursor()
		result := map[string]map[string]string{}
		for k, v := c.First(); k != nil; k, v = c.Next() {
			//fmt.Printf("key=%s, value=%s\n", k, v)
			result[string(k)] = map[string]string{}
			result[string(k)]["rollNo"] = string(k)
			result[string(k)]["charCount"] = strconv.Itoa(len(v))
			result[string(k)]["wordCount"] = strconv.Itoa(wordCount(string(v)))
			result[string(k)]["actualMessage"] = string(v)
		}
		data = result
		return nil
	})
	return
}

func wordCount(s string) int {
	words := strings.Fields(s)
	return len(words)
}

func studentHandler(s *server) websocket.Handler {
	return func(ws *websocket.Conn) {
		websocketStudentConnection(ws, s)
	}
}

func websocketTeacherConnection(ws *websocket.Conn, s *server) {
	x, err := s.Get("student")
	bb, _ := json.Marshal(x)
	if err != nil {
		log.Println("some error occured while sending fetching data on load", err)
	} else {
		sendDataToTeacher(ws, string(bb))
	}
}

func sendDataToTeacher(ws *websocket.Conn, data string) {
	websocket.Message.Send(ws, data)
}

func teacherHandler(s *server) websocket.Handler {
	return func(ws *websocket.Conn) {
		s.teacherConn = ws
		websocketTeacherConnection(ws, s)
		for {
		}
	}
}

func main() {
	server, err := newServer("student.db")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/student", studentHandler(server))
	http.Handle("/teacher", teacherHandler(server))
	http.Handle("/", http.FileServer(http.Dir("static/html/")))
	log.Printf("Server listening on port %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))

	//defer db.Close()

}
