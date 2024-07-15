package websocketHandler

import (
	"GinChat/db"
	"GinChat/serializer"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type Client struct {
	Id     uint
	Socket *websocket.Conn
	Send   chan []byte
}
type ClientManager struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

var (
	postDb *gorm.DB = db.ConnectPostgres()

	// in memory system (every time we restart the server --> it will delete all saved Clients)
	Manager = ClientManager{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
)

func (manager *ClientManager) Start() {

	for {
		select {
		//If there is a new connection access, pass the connection to conn through the channel
		case conn := <-manager.Register:
			fmt.Println("function Start manager.Register case")
			manager.Clients[conn] = true

			//If the connection is disconnected
		case conn := <-manager.Unregister:
			fmt.Println("function Start manager.Unregister case")

			if _, ok := manager.Clients[conn]; ok {
				close(conn.Send)
				delete(manager.Clients, conn)
			}
		//broadcast
		case message := <-manager.Broadcast:
			fmt.Println("function Start manager.Broadcast case")
			//Traversing the client that has been connected, send the message to them
			for client := range manager.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(manager.Clients, client)
				}
			}
		}
	}
}
func (manager *ClientManager) Send(message []byte, ignore *Client) {
	fmt.Println("function send")

	for client := range manager.Clients {
		//Send messages not to the shielded connection
		if client != ignore {
			client.Send <- message
		}
	}
}
func (c *Client) Read() {
	fmt.Println("function Read")
	defer func() {
		Manager.Unregister <- c
		_ = c.Socket.Close()
	}()

	for {
		//Read message
		_, message, err := c.Socket.ReadMessage()
		//If there is an error message, cancel this connection and then close it
		if err != nil {
			Manager.Unregister <- c
			_ = c.Socket.Close()
			break
		}
		//If there is no error message, put the information in Broadcast
		jsonMessage, _ := json.Marshal(serializer.Message{Sender: c.Id, Content: string(message)})
		Manager.Broadcast <- jsonMessage
	}
}
func (c *Client) Write() {
	fmt.Println("function Write")

	defer func() {
		_ = c.Socket.Close()
	}()

	for {
		select {
		//Read the message from send
		case message, ok := <-c.Send:
			//If there is no message
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			//Write it if there is news and send it to the web side
			_ = c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
