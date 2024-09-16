package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"spoker/frame"
	"strings"
)

const (
	STATIC = "static"
	INDEX  = "index.html"

	CONN      = "Connection"
	UPGRADE   = "Upgrade"
	WEBSOCKET = "websocket"
	WSKEY     = "Sec-WebSocket-Key"
	WSVERSION = "Sec-WebSocket-Version"

	NOT_FOUND_RESPONSE       = "<h1>Not Found</h1><p>Could not find what you were looking for :(</p>"
	INTERNAL_SERVER_RESPONSE = "<h1>Internal Server Error</h1><p>An unknown error occurred</p>"

	ACCEPT_BASE = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
)

const (
	ControlFrameClose = 0x8
	ControlFramePing  = 0x9
	ControlFramePong  = 0xA
)

func main() {
	fmt.Println("This is Spoker!")

	http.HandleFunc("/static/", staticHandler)
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/", handler)

	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	index := STATIC + "/" + INDEX

	body, err := os.ReadFile(index)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, INTERNAL_SERVER_RESPONSE)
	}

	fmt.Fprintf(w, "%s", body)
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	filename := STATIC + "/" + r.URL.Path[len("/"+STATIC+"/"):]

	body, err := os.ReadFile(filename)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, NOT_FOUND_RESPONSE)
	}

	fmt.Fprintf(w, "%s", body)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("received websocket connection request")

	if r.ProtoMajor < 1 || r.ProtoMajor == 1 && r.ProtoMinor < 1 {
		w.Header().Add(UPGRADE, r.Proto)
		w.WriteHeader(http.StatusUpgradeRequired)
		fmt.Fprint(w, "This service requires use of the HTTP/1.1 protocol.")
		log.Println("bad HTTP version")
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Add("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(nil)
		log.Println("bad method")
		return
	}

	if !strings.Contains(strings.ToLower(r.Header.Get(CONN)), strings.ToLower(UPGRADE)) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "The request lacks the %s header with an %s value", CONN, UPGRADE)
		log.Println("no connection -> upgrade")
		return
	}

	if strings.ToLower(r.Header.Get(UPGRADE)) != WEBSOCKET {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "The request lacks the %s header with a %s value", CONN, WEBSOCKET)
		log.Println("no upgrade -> websocket")
		return
	}

	wskey := r.Header.Get(WSKEY)

	if len(wskey) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "The request lacks the %s header", WSKEY)
		log.Println("no key")
		return
	}

	key, err := base64.StdEncoding.DecodeString(wskey)
	if err != nil || len(key) != 16 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "The request does not have a valid value for the %s header", WSKEY)
		log.Println("invalid key")
		return
	}

	if r.Header.Get(WSVERSION) != "13" {
		w.Header().Add(WSVERSION, "13")
		w.WriteHeader(http.StatusUpgradeRequired)
		fmt.Fprintf(w, "The request is for an unsupported WebSocket version")
		log.Println("bad WS version")
		return
	}

	log.Println("request validated")

	conn, brw, err := http.NewResponseController(w).Hijack()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, INTERNAL_SERVER_RESPONSE)
		log.Println("error hijacking connection")
	}

	hasher := sha1.New()
	fmt.Fprint(hasher, wskey)
	fmt.Fprint(hasher, ACCEPT_BASE)

	buf := brw.Writer.AvailableBuffer()
	buf = append(buf, "HTTP/1.1 101 Switching Protocols\r\n"...)
	buf = append(buf, fmt.Sprintf("%s: %s\r\n", UPGRADE, WEBSOCKET)...)
	buf = append(buf, fmt.Sprintf("%s: %s\r\n", CONN, UPGRADE)...)
	buf = append(buf, fmt.Sprintf("Sec-WebSocket-Accept: %s\r\n", base64.StdEncoding.EncodeToString(hasher.Sum(nil)))...)
	buf = append(buf, "\r\n"...)
	if _, err := conn.Write(buf); err != nil {
		log.Fatal("couldn't complete socket")
	}

	go handleWSConn(conn, brw)
}

func handleWSConn(conn net.Conn, brw *bufio.ReadWriter) {
	_ = conn
	br := brw.Reader
	bw := brw.Writer
	frame.NewString(true, "hey, this is a message from the server").WriteToBuffer(bw)

	for {
		f, err := frame.ReadToNewFrame(br)
		if err != nil {
			continue
		}

		switch f.Opcode {
		case ControlFrameClose:
			log.Println("received close")
		case ControlFramePing:
			log.Println("received ping")
		case ControlFramePong:
			log.Println("received pong")
		default:
			f = frame.NewString(true, fmt.Sprintf("you sent: %s", f.PayloadData))
			f.WriteToBuffer(bw)
			log.Printf("%x: %s", f.Opcode, f.PayloadData)
		}
	}
}
