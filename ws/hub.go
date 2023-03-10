package ws

import (
	"context"
	"encoding/json"
	"log"
	"soporte-go/core/model/ws"
)

type message struct {
	data []byte
	room string
}

type subscription struct {
	conn *connection
	room string
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	rooms map[string]map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan message

	// Register requests from the connections.
	register chan subscription

	// Unregister requests from connections.
	unregister chan subscription
}

var H = hub{
	broadcast:  make(chan message),
	register:   make(chan subscription),
	unregister: make(chan subscription),
	rooms:      make(map[string]map[*connection]bool),
}

func (h *hub) Run(re ws.WsRepository) {
	ctx := context.Background()
	for {
		select {
		case s := <-h.register:
			log.Println("registrando")
			// log.Println(s.room)
			// log.Println(s.room)
			// query := `SELECT * FROM messages WHERE caso_id = ?`
			// value := database.ExecuteQuery3(query, s.room)
			// jsonStr, err := json.Marshal(value)
			// if err != nil {
			// 	log.Printf("Error: %s", err.Error())
			// } else {
			// 	s.conn.send <- jsonStr
			// }
			// log.Println(h.rooms)
			connections := h.rooms[s.room]
			if connections == nil {
				// log.Println(connections)
				// log.Println("no room")
				connections = make(map[*connection]bool)
				h.rooms[s.room] = connections
			}
			// log.Println(connections)
			// log.Println("si room")
			h.rooms[s.room][s.conn] = true
			res,err := re.GetMessages(ctx,s.room)
			if err != nil {
				log.Println(err)
			}
			jsonStr ,err := json.Marshal(res)
			if err != nil {
				log.Println(err)
			}else{
				s.conn.send <- jsonStr
			}
		case s := <-h.unregister:
			log.Println("Unregister connection")
			connections := h.rooms[s.room]
			if connections != nil {
				if _, ok := connections[s.conn]; ok {
					delete(connections, s.conn)
					close(s.conn.send)
					if len(connections) == 0 {
						delete(h.rooms, s.room)
					}
				}
			}
		case m := <-h.broadcast:
			connections := h.rooms[m.room]
			log.Println("Sending data")
			log.Println(string(m.data))
		
			msgData := &ws.MessageData{}
			err := json.Unmarshal(m.data,msgData)
			if err != nil {
				log.Println(err)
			}
			_,err = re.SaveMessage(ctx,msgData)
			if err != nil{
				log.Println(err)
			}
			log.Println(msgData)
			// log.Println(connections)
			for c := range connections {
				select {
				case c.send <- m.data:
				default:
					// log.Println("close - delete")
					close(c.send)
					delete(connections, c)
					if len(connections) == 0 {
						// log.Println("delete room")
						delete(h.rooms, m.room)
					}
				}
			}
		}
	}
}

// create table messages(
// 	caso_id text,
// 	content text,
// 	created timestamp,
// 	from_user text,
// 	from_user_id text,
// 	to_user text,
// 	to_user_id text,
// 	read int,
// 	PRIMARY KEY (caso_id,content));
