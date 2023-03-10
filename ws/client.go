package ws

import (
	"log"
	"net/http"

	"time"

	// "strings"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// allow all connections by default
		return true
	},
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

type Iot struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

// readPump pumps messages from the websocket connection to the hub.
func (s subscription) readPump() {
	c := s.conn
	defer func() {
		H.unregister <- s
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.ws.ReadMessage()
	
		// empty := make([]byte,3)
		// builder := strings.Builder{}
		// builder.Write(msg)
		// log.Println(builder.String())
		// if msg != nil {
		// 	log.Println("lllllllll")
		// 	log.Println(string(msg))
		// 	data := &model.Message{}
		// 	go func(){
		// 		if err := json.Unmarshal(msg, data); err != nil {
		// 			log.Printf("No se pudo leer el mensaje:err: %s", err.Error())
		// 			} else {
		// 				log.Println(data)
		// 			}
		// 			query := `INSERT INTO messages (caso_id,from_user,from_user_id,to_user,content,created) values(?,?,?,?,?,?)`
		// 			database.ExecuteQuery(query,data.CasoId,data.FromUser,data.FromUserId,data.ToUser,data.Content,time.Now())
		// 			}()
		// 		}
		// // messag := builder.String()
		// in := []byte(`{"id":1,"name":"test","context":{"key1":"value1","key2":2}}`)
		// var iot Iot
		// err := json.Unmarshal(in, &iot)
		// if err != nil {
		// 	panic(err)
		// }
		// fmt.Println("ctx:", string(iot.Type))
		// in := []byte(`{"id":1,"name":"test","context":{"key1":"value1","key2":2}}`)

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		m := message{msg, s.room}
		H.broadcast <- m
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (s *subscription) writePump() {
	c := s.conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				log.Println("send-event1....")
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				log.Println("send-event2....")
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}


// serveWs handles websocket requests from the peer.
func ServeWs(w http.ResponseWriter, r *http.Request, roomId string) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws}
	s := subscription{c, roomId}
	H.register <- s
	go s.writePump()
	go s.readPump()
}
