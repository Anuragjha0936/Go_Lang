package WS

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"main.go/database"
	"main.go/models"
	"main.go/utils"
)

// upgrader will convert the http request to web Socket connection
// check origin only allow my frontend to start the ws connection
var upgrader=websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
    return true
},
}

// the method serveWs will return a handler func when Get /ws being called
// we send the token as a url to the server due to the limitation of the browser 
// Notice that new WebSocket() only accepts a URL (and optionally subprotocols). It does not let you set custom HTTP headers 
func ServeWS(hub *Hub) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		// fetching token from url
		token:=r.URL.Query().Get("token")
		if token==""{
			http.Error(w,"missing token",http.StatusUnauthorized)
			return 
		}
		// validate the token
		claims,err:=utils.ValidateToken(token)
		if err!=nil{
			http.Error(w,"invalid or expired token",http.StatusUnauthorized)
			return
		}

		userId:=claims.Id
		conn,err:=upgrader.Upgrade(w,r,nil)
		if err!=nil{
			fmt.Println("ws upgrade failed",err)
			return 
		}
		
		client:=hub.Register(userId,conn)
		defer func(){
			hub.UnRegister(userId,client)
			conn.Close()
		}()

		// infinite loop so the server always waiting for the client 
		// to receive the message and process
		for{
			var incomeMessage models.Message
			err:=conn.ReadJSON(&incomeMessage)
			if err!=nil{
				fmt.Println("WebSocket closed for user",userId,":",err)
				return
			}

			// now save the incoming message into the db
			saved,err:=database.SaveMessage(userId,incomeMessage.ReceiverId,incomeMessage.Content)
			if err!=nil{
				fmt.Println("failed to save message",err)
				_ = client.WriteJSON(map[string]string{"error": "failed to send message"})
				continue
			}

			// deliever the message to receiver
			_,err=hub.SendToUser(incomeMessage.ReceiverId,saved)
			if err!=nil{
				http.Error(w,"Unable to Send Message",http.StatusBadRequest)
				return
			}
		   
			// when receiver is offline which means the ws connection of receiver is closed 
			// the server not able to send the messag eto the receiver then it save the message in he db
			// and echo to the sender

			if err := client.WriteJSON(saved); err != nil {
				log.Println("failed to echo message to sender:", err)
				return
			}
		}
	}

}