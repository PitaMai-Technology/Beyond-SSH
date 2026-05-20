package main

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"

	"nhooyr.io/websocket"
)

func main() {
	http.HandleFunc("/ssh", func(w http.ResponseWriter, r *http.Request) {
		host := r.URL.Query().Get("host")
		if host == "" {
			http.Error(w, "missing host parameter", http.StatusBadRequest)
			return
		}

		c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			InsecureSkipVerify: true, // Allow all origins for dev
		})
		if err != nil {
			log.Printf("websocket accept error: %v", err)
			return
		}
		defer c.Close(websocket.StatusInternalError, "internal error")

		log.Printf("Dialing tcp %s...", host)
		tcpConn, err := net.Dial("tcp", host)
		if err != nil {
			log.Printf("tcp dial error: %v", err)
			c.Close(websocket.StatusBadGateway, "tcp dial error")
			return
		}
		defer tcpConn.Close()

		wsConn := websocket.NetConn(context.Background(), c, websocket.MessageBinary)

		// Bi-directional copy
		errc := make(chan error, 2)
		go func() {
			_, err := io.Copy(wsConn, tcpConn)
			errc <- err
		}()
		go func() {
			_, err := io.Copy(tcpConn, wsConn)
			errc <- err
		}()

		<-errc
		log.Printf("Connection closed for %s", host)
	})

	log.Println("Proxy listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
