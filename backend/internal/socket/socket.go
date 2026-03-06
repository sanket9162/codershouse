package socket

import (
	"log/slog"
	"net/http"

	"github.com/zishang520/socket.io/v2/socket"
)

type SocketServer struct {
	io     *socket.Server
	logger *slog.Logger
}

func NewSocketServer(logger *slog.Logger) *SocketServer {
	io := socket.NewServer(nil, nil)

	s := &SocketServer{
		io:     io,
		logger: logger,
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
		client.On("join", func(args ...any) {
			s.logger.Info("Join event triggered", "args", args)
		})

		// relay ice
		client.On("relayICE", func(args ...any) {
			s.logger.Info("relayICE event triggered")
		})

		// relay sdp
		client.On("relaySDP", func(args ...any) {
			s.logger.Info("relaySDP event triggered")
		})

		// disconnect
		client.On("disconnect", func(args ...any) {
			s.logger.Info("Socket disconnected", "id", client.Id())
		})
	})
}
