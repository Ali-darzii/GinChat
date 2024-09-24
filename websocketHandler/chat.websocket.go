package websocketHandler

import (
	"GinChat/db"
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
	Broadcast  chan serializer.MessageV2
	Register   chan *Client
	Unregister chan *Client
}

var (
	postDb *gorm.DB = db.ConnectPostgres()

	// in memory system
	Manager = ClientManager{
		Broadcast:  make(chan serializer.MessageV2),
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

			switch message.PvMessage.Type {
			// because u need a send message to create a pv_message --> pv_message and new_pv_message are same
			case "pv_message", "new_pv_message":
				jsonMessage, _ = json.Marshal(&message.PvMessage)

				//NEED UPDATE BECAUSE OF SERIALIZER and METHOD CHANGING !!!!!!
				/*/
				case "group_message":
					groupMessage := entity.GroupMessageRoom{
						Sender:  message.Sender,
						GroupID: message.RoomID,
						Body:    &message.Content,
					}
					if res := postDb.Save(&groupMessage); res.Error != nil {
						jsonMessage, _ = json.Marshal(&serializer.ServerMessage{Content: "can't save message in db", RoomID: message.RoomID, Status: false})
						message.Recipients = []uint{message.Sender}
					} else {
						gpMessage := serializer.SendPvMessage{
							Type:      message.Type,
							RoomID:    message.RoomID,
							Content:   message.Content,
							Sender:    message.Sender,
							TimeStamp: groupMessage.TimeStamp,
						}
						jsonMessage, _ = json.Marshal(&gpMessage)
					}

				case "new_group_message":
					gpMessage := serializer.NewGroupChat{
						Avatar:  message.Avatar,
						Type:    message.Type,
						RoomID:  message.RoomID,
						Members: message.Recipients,
					}
					jsonMessage, _ = json.Marshal(&gpMessage)
				/*/
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

// Reading Method because of file saving(image and voice) moved on http then on ws
func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c
		_ = c.Socket.Close()
	}()

	for {
		//Read message
		_, _, err := c.Socket.ReadMessage()

		if err != nil {
			c.Disconnect()
			break
		}
		serverMessage := serializer.ServerMessage{Content: "Read Message doesn't work, use api(cause of sending image and voice)"}
		jsonMessage, err := json.Marshal(serverMessage)
		if err != nil {
			c.Send <- jsonMessage
		}

		/*
			var readMessage serializer.Message
			if err = json.Unmarshal(message, &readMessage); err != nil {
				c.Disconnect()
				break
			}
			readMessage.Sender = c.Id

			switch readMessage.Type {
			case "pv_message":
				if ok := readMessage.PrivateMessageValidate(); !ok {
					c.Disconnect()
					break
				}

				if res := postDb.Table("pv_users").Select("user_id").Where("private_room_id = ?", readMessage.RoomID).Pluck("user_id", &readMessage.Recipients); res.Error != nil {
					c.Disconnect()
					break
				}
				var sameRoom bool
				//checking
				for _, item := range readMessage.Recipients {
					if item == c.Id {
						//readMessage.Recipients = append(readMessage.Recipients[:index], readMessage.Recipients[index+1:]...)
						sameRoom = true
					}
				}
				// if clients are not in the same room
				if !sameRoom {
					c.Disconnect()
					break
				}

				Manager.Broadcast <- readMessage

			case "group_message":
				if res := postDb.Table("group_users").Select("user_id").Where("group_room_id = ?", readMessage.RoomID).Pluck("user_id", &readMessage.Recipients); res.Error != nil {
					c.Disconnect()
					break
				}
				var sameRoom bool
				//checking
				for _, item := range readMessage.Recipients {
					if item == c.Id {
						//readMessage.Recipients = append(readMessage.Recipients[:index], readMessage.Recipients[index+1:]...)
						sameRoom = true
					}
				}
				// if clients are not in the same room
				if !sameRoom {
					c.Disconnect()
					break
				}

				Manager.Broadcast <- readMessage

			default:
				close(c.Send)
				delete(Manager.Clients, c)

			}
		*/
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
func (c *Client) Disconnect() {
	Manager.Unregister <- c
	_ = c.Socket.Close()

}
