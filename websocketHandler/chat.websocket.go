package websocketHandler

import (
	"GinChat/db"
	"GinChat/entity"
	"GinChat/serializer"
	"encoding/json"
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
	Broadcast  chan serializer.Message
	Register   chan *Client
	Unregister chan *Client
}

var (
	postDb *gorm.DB = db.ConnectPostgres()

	// in memory system (every time we restart the server --> it will delete all saved Clients)
	Manager = ClientManager{
		Broadcast:  make(chan serializer.Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
)

func (manager *ClientManager) Start() {

	for {
		select {
		//If there is a new connection access, pass the connection to conn through the channel
		case client := <-manager.Register:
			manager.Clients[client] = true
			jsonMessage, _ := json.Marshal(&serializer.ServerMessage{Content: "Connected from server !", Status: true})
			client.Send <- jsonMessage
			//If the connection is disconnected
		// disconnected clients
		case client := <-manager.Unregister:
			if _, ok := manager.Clients[client]; ok {
				close(client.Send)
				delete(manager.Clients, client)
			}
		//broadcast
		case message := <-manager.Broadcast:
			var jsonMessage []byte

			switch message.Type {
			case "pv_message":
				var privateMessage = entity.PrivateMessageRoom{
					Sender:    message.Sender,
					PrivateID: message.RoomID,
					Body:      message.Content,
				}
				if res := postDb.Save(&privateMessage); res.Error != nil {
					jsonMessage, _ = json.Marshal(&serializer.ServerMessage{Content: "can't save in db", Status: false})
					message.Recipients = []uint{message.Sender}
				} else {
					jsonMessage, _ = json.Marshal(&privateMessage)

				}

			case "group_message":
				continue
			case "new_pv_message":
				continue
			case "new_group_message":
				continue

			}

			for _, item := range message.Recipients {
				for client := range manager.Clients {
					if client.Id == item {
						select {
						case client.Send <- jsonMessage:
						default:
							close(client.Send)
							delete(manager.Clients, client)
						}
					}
				}
			}
		}
	}
}

func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c
		_ = c.Socket.Close()
	}()

	for {
		//Read message
		_, message, err := c.Socket.ReadMessage()

		if err != nil {
			Manager.Unregister <- c
			_ = c.Socket.Close()
			break
		}

		var readMessage serializer.Message
		if err = json.Unmarshal(message, &readMessage); err != nil {
			Manager.Unregister <- c
			_ = c.Socket.Close()
			break
		}

		if ok := readMessage.Validate(); !ok {
			Manager.Unregister <- c
			_ = c.Socket.Close()
			break
		}

		readMessage.Sender = c.Id
		switch readMessage.Type {
		case "pv_message":
			if res := postDb.Table("pv_users").Select("user_id").Where("private_room_id = ?", readMessage.RoomID).Pluck("user_id", &readMessage.Recipients); res.Error != nil {
				Manager.Unregister <- c
				_ = c.Socket.Close()
				break
			}
			var sameRoom bool
			for index, item := range readMessage.Recipients {
				if item == c.Id {
					readMessage.Recipients = append(readMessage.Recipients[:index], readMessage.Recipients[index+1:]...)
					sameRoom = true
				}
			}
			// client must be in same room
			if !sameRoom {
				Manager.Unregister <- c
				_ = c.Socket.Close()
				break
			}
			//jsonMessage, _ := json.Marshal(&serializer.Message{
			//	Type:       readMessage.Type,
			//	Content:    readMessage.Content,
			//	Recipients: readMessage.Recipients,
			//	RoomID:     readMessage.RoomID,
			//	Sender:     readMessage.Sender,
			//})

			Manager.Broadcast <- readMessage

		case "group_message":
			continue
		case "new_pv_message":
			continue
		case "new_group_message":
			continue
		default:
			close(c.Send)
			delete(Manager.Clients, c)

		}
		//If there is no error message, put the information in Broadcast
		//jsonMessage, _ := json.Marshal(&serializer.ServerMessage{Content: string(message)})
		//Manager.Broadcast <- jsonMessage
	}
}
func (c *Client) Write() {

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

//func (manager *ClientManager) Send(message []byte, ignore *Client) {
//	for client := range manager.Clients {
//		//Send messages not to the shielded connection
//		if client != ignore {
//			client.Send <- message
//		}
//	}
//}
