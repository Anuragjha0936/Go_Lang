package WS

import (
	"sync"

	"github.com/gorilla/websocket"
)

// client represent the user and his ws connection
// WriteMu states Only one goroutine should write to a connection at a time.
type Client struct{
	conn *websocket.Conn
	WriteMu sync.Mutex
}

// so that no two connection write to one connection simultaneously
func (c *Client) WriteJSON (v interface{}) error{
	c.WriteMu.Lock()
	defer c.WriteMu.Unlock()
	return c.conn.WriteJSON(v)
}

// Hub tracks currently-connected users: userId -> their live Client.
// A user with no entry here is simply offline -- messages sent to
// them are still saved to the DB (see handler.go), just not delivered
// live; they'll see them next time they call GET /api/messages/{id}.

type Hub struct{
	mu sync.RWMutex
	clients map[int]*Client
}

// we are using RWMutex so that multiple goroutines can read the same user but when one goroutine try to write 
// then their is lock so no other goroutine can read or write

// new Hub will create a empty map means nobody is online
func NewHub() *Hub{
	return &Hub{
		clients:make(map[int]*Client),
	}
}

// register -> storing the user in the map 
// the map contains the user id and the Ws connection of the user
func (h*Hub) Register(userID int,conn *websocket.Conn) *Client{
	client:=&Client{conn: conn}
	h.mu.Lock()
	h.clients[userID]=client
	h.mu.Unlock()
	return client
}

// Unregister-> when a user get offline remove the WS conn
func (h *Hub) UnRegister(UserID int ,client *Client){
	h.mu.Lock()
	defer h.mu.Unlock()
	existing,ok:=h.clients[UserID]
	if ok && existing==client{
		delete(h.clients,UserID)
	}
}

func (h *Hub)SendToUser(userID int,v interface{}) (delivered bool,err error){
	h.mu.RLock()
	client,ok:=h.clients[userID]
	h.mu.RUnlock()
	if !ok {
		return false, nil
	}
	if err := client.WriteJSON(v); err != nil {
		return false, err
	}
	return true, nil
}


