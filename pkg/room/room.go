package room

import (
	"github.com/louis296/p2p-server/pkg/log"
	"github.com/louis296/p2p-server/pkg/server"
	"github.com/louis296/p2p-server/pkg/util"
	"net/http"
)

const (
	JoinRoom       = "joinRoom"
	Offer          = "offer"
	Answer         = "answer"
	Candidate      = "candidate"
	HangUp         = "hangUp"
	LeaveRoom      = "leaveRoom"
	UpdateUserList = "updateUserList"
	Talk           = "talk"
)

type Room struct {
	users    map[string]User
	sessions map[string]Session
	ID       string
}

func NewRoom(id string) *Room {
	room := &Room{
		users:    make(map[string]User),
		sessions: make(map[string]Session),
		ID:       id,
	}
	return room
}

type Manager struct {
	rooms map[string]*Room
}

func NewManager() *Manager {
	manager := &Manager{rooms: make(map[string]*Room)}
	return manager
}

func (m *Manager) getRoom(id string) *Room {
	return m.rooms[id]
}

func (m *Manager) createRoom(id string) *Room {
	m.rooms[id] = NewRoom(id)
	return m.rooms[id]
}

func (m *Manager) deleteRoom(id string) {
	delete(m.rooms, id)
}

func (m *Manager) HandleNewWebSocket(conn *server.WebSocketConn, request *http.Request) {
	log.Info("On Open %v", request)
	conn.On("message", func(message []byte) {
		req, err := util.Unmarshal(string(message))
		if err != nil {
			log.Error("Not a json data: %v", err)
			return
		}
		var data map[string]interface{} = nil
		tmp, ok := req["data"]
		if !ok {
			log.Error("No data from message: %v", err)
		}
		data = tmp.(map[string]interface{})
		roomId := data["roomId"].(string)
		log.Info("Room Id: %v", roomId)

		room := m.getRoom(roomId)
		if room == nil {
			room = m.createRoom(roomId)
		}

		switch req["type"] {
		case JoinRoom:
			onJoinRoom(conn, data, room, m)
		case Talk:
			onTalk(conn, data, room)
		case Offer:
			fallthrough
		case Answer:
			fallthrough
		case Candidate:
			onCandidate(data, room, req)
		case HangUp:
			onHangUp()
		default:
			log.Warn("Unknown request type %v", req)
		}
	})

	conn.On("close", func(code int, text string) {
		onClose(conn, m)
	})
}

func (m *Manager) notifyUsersUpdate(users map[string]User) {
	var infos []UserInfo
	for _, user := range users {
		infos = append(infos, user.info)
	}
	request := make(map[string]interface{})
	request["type"] = UpdateUserList
	request["data"] = infos
	for _, user := range users {
		user.conn.Send(util.Marshal(request))
	}
}

func onJoinRoom(conn *server.WebSocketConn, data map[string]interface{}, room *Room, manager *Manager) {
	user := User{
		info: UserInfo{
			ID:   data["id"].(string),
			Name: data["name"].(string),
		},
		conn: conn,
	}
	room.users[user.info.ID] = user
	manager.notifyUsersUpdate(room.users)
}

func onTalk(conn *server.WebSocketConn, data map[string]interface{}, room *Room) {
	var userInfo UserInfo
	for _, user := range room.users {
		if user.conn == conn {
			userInfo = user.info
		}
	}
	talkMessage := map[string]interface{}{
		"senderInfo": userInfo,
		"content":    data["content"],
	}
	msg := map[string]interface{}{
		"type": Talk,
		"data": talkMessage,
	}
	for _, user := range room.users {
		user.conn.Send(util.Marshal(msg))
	}
}

func onCandidate(data map[string]interface{}, room *Room, request map[string]interface{}) {
	to := data["to"].(string)
	if user, ok := room.users[to]; !ok {
		log.Error("Not found user [%v]", to)
		return
	} else {
		user.conn.Send(util.Marshal(request))
	}
}

func onHangUp() {

}

func onClose(conn *server.WebSocketConn, manager *Manager) {
	log.Info("Connection close %v", conn)
	userId, roomId := "", ""
	for _, room := range manager.rooms {
		for _, user := range room.users {
			if user.conn == conn {
				userId = user.info.ID
				roomId = room.ID
				break
			}
		}
	}

	if roomId == "" {
		log.Error("Not enter any room")
		return
	}

	log.Info("Close conn with roomId [%v], userId [%v]", roomId, userId)

	for _, user := range manager.getRoom(roomId).users {
		if user.conn != conn {
			leave := map[string]interface{}{
				"type": LeaveRoom,
				"data": userId,
			}
			user.conn.Send(util.Marshal(leave))
		}
	}
	log.Info("Delete user %v", userId)
	delete(manager.getRoom(roomId).users, userId)

	manager.notifyUsersUpdate(manager.getRoom(roomId).users)
}
