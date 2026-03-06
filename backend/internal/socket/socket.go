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
)

type SocketUser struct {
	RoomID string
	User   map[string]interface{}
}

type SocketServer struct {
	io     *socket.Server
	logger *slog.Logger
	mu     sync.RWMutex
	users  map[socket.SocketId]*SocketUser
}

func NewSocketServer(logger *slog.Logger) *SocketServer {
	io := socket.NewServer(nil, nil)

	s := &SocketServer{
		io:     io,
		logger: logger,
		users:  make(map[socket.SocketId]*SocketUser),
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
			s.mu.Unlock()

			// Get all sockets currently in this room natively
			s.mu.RLock()
			var existingSockets []socket.SocketId
			for sid, u := range s.users {
				if u.RoomID == roomId && sid != client.Id() {
					existingSockets = append(existingSockets, sid)
				}
			}
			s.mu.RUnlock()

			for _, existingSocketId := range existingSockets {
				// Tell existing client about the new client
				s.mu.RLock()
				newUserProfile := s.users[client.Id()].User
				existingUserProfile := s.users[existingSocketId].User
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

		// disconnect (or leave)
		leave := func() {
			s.mu.Lock()
			user, ok := s.users[client.Id()]
			if ok {
				delete(s.users, client.Id())
			}
			s.mu.Unlock()

			if ok && user != nil {
				s.mu.RLock()
				var remainingSockets []socket.SocketId
				for sid, u := range s.users {
					if u.RoomID == user.RoomID {
						remainingSockets = append(remainingSockets, sid)
					}
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
