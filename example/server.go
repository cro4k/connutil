package example

import (
	"github.com/cro4k/connutil"
	"github.com/cro4k/connutil/wsutil"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"time"
)

func Listen() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		handler(id).ServeHTTP(w, r)
	})
	addr := ":4567"
	log.Println("started on " + addr)
	http.ListenAndServe(addr, nil)
}

func handler(id string) websocket.Handler {
	return func(conn *websocket.Conn) {
		c, _ := connutil.NewClient(id, wsutil.NewBytesConn(conn), listener)
		go func(cli *connutil.Client) {
			t := time.Tick(time.Second * 3)
			for {
				<-t
				c.Write([]byte("Hello, " + id))
			}
		}(c)
		c.Run()
	}
}

var listener = connutil.NewListener(
	func(c *connutil.Client) {
		log.Println(c.ID(), "connected")
	},
	func(c *connutil.Client) {
		log.Println(c.ID(), "disconnected")
	},
	func(c *connutil.Client) {
		log.Println(c.ID(), "reconnected")
	},
	func(c *connutil.Client, data []byte) {
		log.Println(c.ID(), "ping")
	},
	func(c *connutil.Client) {
		log.Println(c.ID(), "removed")
	},
)
