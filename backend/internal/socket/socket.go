package socket

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/zishang520/socket.io/v2/socket"
)

const (
	ActionJoin               = "join"
	ActionLeave              = "leave"
	ActionAddPeer            = "add-peer"
	ActionRelayICE           = "relay-ice"
	ActionRelaySDP           = "relay-sdp"
	ActionSessionDescription = "session-description"
	ActionIceCandidate       = "ice-candidate"
	ActionRemovePeer         = "remove-peer"
	ActionMute               = "mute"
	ActionUnMute             = "unmute"
)

type SocketUser struct {
	RoomID string
	User   map[string]interface{}
}

type SocketServer struct {
	io        *socket.Server
	logger    *slog.Logger
	mu        sync.RWMutex
	users     map[socket.SocketId]*SocketUser
	roomUsers map[string]map[socket.SocketId]map[string]interface{}
}

func NewSocketServer(logger *slog.Logger) *SocketServer {
	io := socket.NewServer(nil, nil)

	s := &SocketServer{
		io:        io,
		logger:    logger,
		users:     make(map[socket.SocketId]*SocketUser),
		roomUsers: make(map[string]map[socket.SocketId]map[string]interface{}),
	}

	s.setupRoutes()
	return s
}

func (s *SocketServer) Handler() http.Handler {
	return s.io.ServeHandler(nil)
}

func (s *SocketServer) setupRoutes() {
	s.io.On("connection", func(clients ...any) {
		client := clients[0].(*socket.Socket)
		s.logger.Info("New socket connection", "id", client.Id())

		// join handler
		client.On(ActionJoin, func(args ...any) {
			if len(args) == 0 {
				return
			}
			data, ok := args[0].(map[string]interface{})
			if !ok {
				return
			}

			roomId, _ := data["roomId"].(string)

			// Store mapping
			s.mu.Lock()
			s.users[client.Id()] = &SocketUser{
				RoomID: roomId,
				User:   data,
			}
			if s.roomUsers[roomId] == nil {
				s.roomUsers[roomId] = make(map[socket.SocketId]map[string]interface{})
			}
			s.roomUsers[roomId][client.Id()] = data
			s.mu.Unlock()

			// Get all sockets currently in this room efficiently via local map
			s.mu.RLock()
			var existingSockets []socket.SocketId
			for sid := range s.roomUsers[roomId] {
				if sid != client.Id() {
					existingSockets = append(existingSockets, sid)
				}
			}
			s.mu.RUnlock()

			for _, existingSocketId := range existingSockets {
				// Tell existing client about the new client
				s.mu.RLock()
				newUserProfile := s.roomUsers[roomId][client.Id()]
				existingUserProfile := s.roomUsers[roomId][existingSocketId]
				s.mu.RUnlock()

				s.io.To(socket.Room(existingSocketId)).Emit(ActionAddPeer, map[string]interface{}{
					"peerId":      client.Id(),
					"createOffer": false,
					"user":        newUserProfile,
				})

				// Tell the new client about the existing client
				client.Emit(ActionAddPeer, map[string]interface{}{
					"peerId":      string(existingSocketId),
					"createOffer": true,
					"user":        existingUserProfile,
				})
			}

			client.Join(socket.Room(roomId))
			s.logger.Info("User joined room", "roomId", roomId, "peerId", client.Id())
		})

		// relay ice
		client.On(ActionRelayICE, func(args ...any) {
			if len(args) == 0 {
				return
			}
			data, ok := args[0].(map[string]interface{})
			if !ok {
				return
			}

			peerId, _ := data["peerId"].(string)
			iceCandidate := data["iceCandidate"]

			s.io.To(socket.Room(peerId)).Emit(ActionIceCandidate, map[string]interface{}{
				"peerId":       client.Id(),
				"iceCandidate": iceCandidate,
			})
		})

		// relay sdp
		client.On(ActionRelaySDP, func(args ...any) {
			if len(args) == 0 {
				return
			}
			data, ok := args[0].(map[string]interface{})
			if !ok {
				return
			}

			peerId, _ := data["peerId"].(string)
			sessionDescription := data["sessionDescription"]

			s.io.To(socket.Room(peerId)).Emit(ActionSessionDescription, map[string]interface{}{
				"peerId":             client.Id(),
				"sessionDescription": sessionDescription,
			})
		})

		// mute
		client.On(ActionMute, func(args ...any) {
			if len(args) == 0 {
				return
			}
			data, ok := args[0].(map[string]interface{})
			if !ok {
				return
			}

			roomId, _ := data["roomId"].(string)
			userId, _ := data["userId"].(string)

			s.mu.Lock()
			if u, exists := s.users[client.Id()]; exists {
				u.User["muted"] = true
				if s.roomUsers[u.RoomID] != nil && s.roomUsers[u.RoomID][client.Id()] != nil {
					s.roomUsers[u.RoomID][client.Id()]["muted"] = true
				}
			}
			s.mu.Unlock()

			s.io.To(socket.Room(roomId)).Emit(ActionMute, map[string]interface{}{
				"peerId": client.Id(),
				"userId": userId,
			})
		})

		// unmute
		client.On(ActionUnMute, func(args ...any) {
			if len(args) == 0 {
				return
			}
			data, ok := args[0].(map[string]interface{})
			if !ok {
				return
			}

			roomId, _ := data["roomId"].(string)
			userId, _ := data["userId"].(string)

			s.mu.Lock()
			if u, exists := s.users[client.Id()]; exists {
				u.User["muted"] = false
				if s.roomUsers[u.RoomID] != nil && s.roomUsers[u.RoomID][client.Id()] != nil {
					s.roomUsers[u.RoomID][client.Id()]["muted"] = false
				}
			}
			s.mu.Unlock()

			s.io.To(socket.Room(roomId)).Emit(ActionUnMute, map[string]interface{}{
				"peerId": client.Id(),
				"userId": userId,
			})
		})

		// disconnect (or leave)
		leave := func() {
			s.mu.Lock()
			user, ok := s.users[client.Id()]
			if ok {
				delete(s.users, client.Id())
				if s.roomUsers[user.RoomID] != nil {
					delete(s.roomUsers[user.RoomID], client.Id())
					if len(s.roomUsers[user.RoomID]) == 0 {
						delete(s.roomUsers, user.RoomID)
					}
				}
			}
			s.mu.Unlock()

			if ok && user != nil {
				s.mu.RLock()
				var remainingSockets []socket.SocketId
				for sid := range s.roomUsers[user.RoomID] {
					remainingSockets = append(remainingSockets, sid)
				}
				s.mu.RUnlock()

				for _, remainingClientSid := range remainingSockets {
					s.io.To(socket.Room(remainingClientSid)).Emit(ActionRemovePeer, map[string]interface{}{
						"peerId": client.Id(),
						"userId": user.User["peerId"],
					})
				}
				client.Leave(socket.Room(user.RoomID))
			}
			s.logger.Info("Socket disconnected/left", "id", client.Id())
		}

		client.On(ActionLeave, func(args ...any) {
			leave()
		})

		client.On("disconnect", func(args ...any) {
			leave()
		})
	})
}
