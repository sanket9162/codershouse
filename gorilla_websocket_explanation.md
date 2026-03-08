# Implementing WebRTC Signaling with `gorilla/websocket`

If you decide to migrate away from `socket.io` (a Node.js ecosystem standard) and want to write idiomatic, high-performance Go, you should use **`github.com/gorilla/websocket`**.

Unlike `socket.io`, native WebSockets do not have built-in concepts like "Rooms" or "Emitting specific event names". You have to build that architecture yourself using JSON and Channels. This is the Go standard.

Here is how you structure a `gorilla/websocket` signaling server for Codershouse:

---

## 1. The Frontend (Standard Browser API)

You don't need any `npm install` libraries. All modern browsers support `WebSocket` natively.
Instead of emitting events, you `send()` stringified JSON payloads with an `action` type.

```javascript
// frontend/src/socket/index.js
export const socketInit = () => {
    // Note the `ws://` protocol instead of `http://`
    const ws = new WebSocket("ws://localhost:8080/ws");
    
    ws.onopen = () => {
        console.log("Connected to Go WebSocket");
    };

    return ws;
};

// Inside useWebRTC.jsx
socket.current = socketInit();

socket.current.onopen = () => {
    // Instead of socket.emit("join", data)
    socket.current.send(JSON.stringify({
        action: "join",
        data: { roomId: "123", peerId: "456" }
    }));
};

socket.current.onmessage = (event) => {
    const payload = JSON.parse(event.data);
    
    switch (payload.action) {
        case "add-peer":
            console.log("New peer joined!", payload.data);
            break;
        case "session-description":
            // Handle SDP offer/answer...
            break;
    }
};
```

---

## 2. The Go Backend Architecture

To replicate Socket.IO's "Rooms" and broadcast capabilities, you must build a **Hub** pattern.

### A. The Client Struct
Each connected user becomes a `Client` running two internal goroutines: one for reading exactly what WebSockets send, and one for writing back to the browser.

```go
type Client struct {
    hub    *Hub
    conn   *websocket.Conn // The actual gorilla websocket
    send   chan []byte     // Buffered channel for outgoing messages
    peerID string          // React user ID
    roomID string          // The voice room they joined
}
```

### B. The Hub Config
The Hub acts as your centralized `sync.RWMutex` router. It keeps track of which Client is in which Room.

```go
type Hub struct {
    // Maps RoomID -> Map of connected Clients
    rooms      map[string]map[*Client]bool
    broadcast  chan Message
    register   chan *Client
    unregister chan *Client
}

type Message struct {
    Action string      `json:"action"`
    Data   interface{} `json:"data"`
}
```

### C. The Hub Pipeline (The Event Loop)
You run the Hub exactly once when the server boots inside a massive infinite `go handler.run()` loop.

```go
func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            // Someone connected, add them to h.rooms[client.roomID]
            if h.rooms[client.roomID] == nil {
                h.rooms[client.roomID] = make(map[*Client]bool)
            }
            h.rooms[client.roomID][client] = true

        case client := <-h.unregister:
            // Someone disconnected, remove them from the room
            delete(h.rooms[client.roomID], client)
            
            // Replicate ActionRemovePeer broadcasting
            h.BroadcastToRoom(client.roomID, Message{
                Action: "remove-peer",
                Data: map[string]string{"peerId": client.peerID},
            })

        case message := <-h.broadcast:
            // Send a customized JSON payload to everyone in a specific room
            // Example: Looping over h.rooms[roomID] and doing `client.send <- jsonBytes`
        }
    }
}
```

### D. Upgrading HTTP to WebSockets
To accept connections from the React app, you mount the Gorilla Upgrader to an endpoint (`/ws`):

```go
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true }, // Allow React CORS
}

func (h *Handler) ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
    // 1. Upgrade the standard HTTP request into a persistent TCP WebSocket
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }

    // 2. Wrap it inside your architecture
    client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
    
    // 3. Register the client into the Hub background loop
    client.hub.register <- client

    // 4. Spin up concurrent Go routines for non-blocking I/O
    go client.writePump() // Constantly takes things from `send` channel -> browser
    go client.readPump()  // Constantly takes browser events -> `hub.broadcast` channel
}
```

## 3. Summary of Tradeoffs

**Gorilla WebSocket**
- ✅ Pure, Idiomatic Go (extremely performant, standard routing pipelines)
- ✅ No `npm` libraries required on the frontend!
- ❌ You must build "Rooms", event routing (`switch payload.action`), and JSON handling entirely from scratch.
- ❌ Re-connection logic (if wifi drops) must be manually written in React.

**zishang520/socket.io**
- ✅ Has "Rooms" (`io.To(room).Emit()`) completely out-of-the-box.
- ✅ Emits specific event names beautifully (`socket.On("join")`).
- ✅ Frontend library handles automatic reconnections gracefully.
- ❌ Non-standard, large third-party Go library ported from NodeJS.
