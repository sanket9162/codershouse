# Comprehensive Guide to WebRTC and Socket.io

This document provides a detailed explanation of the core files powering the real-time voice rooms: the React frontend hook (`useWebRTC.jsx`) and the Go backend server (`socket.go`).

## 1. The React Frontend: `useWebRTC.jsx`

This React Hook manages the entire lifecycle of a user inside a voice room. It handles microphone access, Peer-to-Peer connections, mapping audio to HTML elements, and communicating with the signaling server.

### State & References
```javascript
export const useWebRTC = (roomId, user) => {
    // Custom state that holds all users currently in the room. 
    // We use a custom hook to handle immediate callbacks when state updates.
    const [clients, setClients] = useStateWithCallback([])
    
    // Stores references to HTML <audio> output tags mapped by User IDs
    const audioElements = useRef({})
    
    // Inter-browser connections mapping. Key: Socket ID, Value: RTCPeerConnection object
    const connections = useRef({})
    
    // Holds the actual hardware microphone feed from your computer
    const localMediaStream = useRef(null)
    
    // The active socket.io connection to the Go server
    const socket = useRef(null)
    
    // A mutable reference to bypass React Closure bugs inside stale event listeners
    const clientsRef = useRef([])
```

### Initialization & Cleanup
```javascript
    useEffect(() => {
        // Connect to the Socket.io Backend instantly when the component mounts
        socket.current = socketInit()

        // Cleanup Function: Runs when you leave the room (/rooms)
        return () => {
            // Turn off the hardware microphone light
            localMediaStream.current?.getTracks().forEach((track) => track.stop())
            
            // Tell the server we are leaving
            socket.current.emit(ACTIONS.LEAVE, { roomId })
            
            // Unplug the websocket connection cleanly
            if (socket.current && socket.current.connected) {
                socket.current.disconnect()
            }
        }
    }, [])
```

### Helper Functions
```javascript
    // Attaches the <audio> tag rendered in Room.jsx to the audioElements map
    const provideRef = (instance, userId) => {
        audioElements.current[userId] = instance
    }

    // Safely adds a user object to the `clients` React State array so they show up on screen
    const addNewClient = useCallback((newClient, cb) => {
        const lookingFor = clients.find((client) => client.id === newClient.id)

        if (!lookingFor) {
            setClients((existingClients) => {
...
```

### Starting the Hardware Media (Microphone)
```javascript
    useEffect(() => {
        const startMedia = async () => {
            try {
                // Request microphone access from the user's browser securely
                const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
                localMediaStream.current = stream
            } catch (error) { ... }
        }
        
        startMedia().then(() => {
            // Inject OURSELVES into the UI array so our avatar appears safely muted
            addNewClient({ ...user, muted: true }, () => {
                
                // Mute our local <audio> tag so we don't hear our own echo playback!
                const localAudio = audioElements.current[user.id]
                if (localAudio) {
                    localAudio.volume = 0 
                    localAudio.srcObject = localMediaStream.current
                }

                // Announce our presence to the Go Server
                socket.current.emit(ACTIONS.JOIN, { roomId, peerId: user.id, ...user })
            })
        })
    }, [])
```

### Handling New Peers & Creating Offers (The Core WebRTC Engine)
```javascript
    useEffect(() => {
        // The server told us someone is already here (or joined after us)
        const handleNewPeer = async ({ peerId, createOffer, user: peerUser }) => {
            // Instantly render them on our screen
            addNewClient(peerUser, () => { })

            // 1. Create a raw network bridge config (RTCPeerConnection)
            connections.current[peerId] = new RTCPeerConnection({
                iceServers: [ // Google's free servers to punch through household routers (NAT traversal)
                    { urls: ['stun:stun.l.google.com:19302'] },
                ],
            })

            // 2. We discovered our own IP route (ICE Candidate). Send it to the peer!
            connections.current[peerId].onicecandidate = (event) => {
                if (event.candidate) {
                    socket.current.emit(ACTIONS.RELAY_ICE, { peerId, iceCandidate: event.candidate })
                }
            }

            // 3. The magic moment: The other person's audio track successfully arrived!
            connections.current[peerId].ontrack = ({ streams: [remoteStream] }) => {
                addNewClient({ ...peerUser, muted: true }, () => {
                    // Pipe their internet audio stream into our HTML <audio> tag to hear them!
                    const remoteAudio = audioElements.current[peerUser.id]
                    if (remoteAudio) remoteAudio.srcObject = remoteStream
                })
            }

            // 4. Pipe OUR microphone feed INTO the new connection tunnel
            localMediaStream.current.getTracks().forEach(track => {
                connections.current[peerId].addTrack(track, localMediaStream.current)
            });

            // 5. If we are the initiating client, generate an SDP Offer (our media capabilities)
            if (createOffer) {
                const offer = await connections.current[peerId].createOffer()
                await connections.current[peerId].setLocalDescription(offer)
                socket.current.emit(ACTIONS.RELAY_SDP, { peerId, sessionDescription: offer })
            }
        }
        socket.current.on(ACTIONS.ADD_PEER, handleNewPeer)
    }, []);
```

### Processing Network Answers (SDP and ICE)
```javascript
    // Handshake Part 2: Handling the other person's media specs (SDP)
    useEffect(() => {
        const handleRelaySDP = async ({ peerId, sessionDescription: remoteOffer }) => {
            const connection = connections.current[peerId]
            
            // Absorb their specs
            await connection.setRemoteDescription(new RTCSessionDescription(remoteOffer))

            // Flush the trickled ICE queues to prevent InvalidStateError crashes
            if (connection.iceCandidatesQueue) {
                connection.iceCandidatesQueue.forEach(async (candidate) => {
                    await connection.addIceCandidate(new RTCIceCandidate(candidate))
                })
            }

            // If they sent us an offer, reply backward with our Answer
            if (remoteOffer.type === 'offer') {
                const answer = await connection.createAnswer()
                await connection.setLocalDescription(answer)
                socket.current.emit(ACTIONS.RELAY_SDP, { peerId, sessionDescription: answer })
            }
        }
        socket.current.on(ACTIONS.SESSION_DESCRIPTION, handleRelaySDP)
    }, [])
```

---

## 2. The Go Backend: `socket.go`

This file is a real-time message relayer. Its only job is to receive metadata from User A, and instantly bounce it precisely back to User B without ever actually touching the audio data.

### Structures & Memory
```go
// Stores room association and any dynamic state (like name, avatar, mute status)
type SocketUser struct {
	RoomID string
	User   map[string]interface{}
}

type SocketServer struct {
	io     *socket.Server
	mu     sync.RWMutex                     // Thread-safe lock to prevent map crashes on concurrent map writes
	users  map[socket.SocketId]*SocketUser  // Global dictionary of everyone connected
}
```

### The Join Flow (Handing out introductions)
```go
		client.On(ActionJoin, func(args ...any) {
            // ... parsing data ...
            
			// 1. Save user to server RAM securely
			s.mu.Lock()
			s.users[client.Id()] = &SocketUser{ RoomID: roomId, User: data }
			s.mu.Unlock()

			// 2. Find everyone else already sitting in this room ID
			s.mu.RLock()
			var existingSockets []socket.SocketId
			for sid, u := range s.users {
				if u.RoomID == roomId && sid != client.Id() {
					existingSockets = append(existingSockets, sid)
				}
			}
			s.mu.RUnlock()

			for _, existingSocketId := range existingSockets {
				// 3. Inform the veteran user: "Hey, someone new just arrived!"
				s.io.To(socket.Room(existingSocketId)).Emit(ActionAddPeer, map[string]interface{}{
					"peerId":      client.Id(),
					"createOffer": false,        // The person already here waits
					"user":        newUserProfile,
				})

				// 4. Inform the new user: "Hey, someone is already sitting here!"
				client.Emit(ActionAddPeer, map[string]interface{}{
					"peerId":      string(existingSocketId),
					"createOffer": true,         // The new person must initiate the WebRTC handshake
					"user":        existingUserProfile,
				})
			}
            
            // Physically subscribe the socket to the Room channel for bulk broadcasts
			client.Join(socket.Room(roomId))
		})
```

### The Muting Engine (Dynamic state modification)
```go
		client.On(ActionMute, func(args ...any) {
            // ... parsing payload ...
            
            // 1. Thread-safe mutation of the User's memory profile to "Muted"
            // This guarantees that if a new person joins 5 minutes from now, 
            // the `ActionAddPeer` payload will correctly tell them we are muted.
			s.mu.Lock()
			if u, exists := s.users[client.Id()]; exists {
				u.User["muted"] = true
			}
			s.mu.Unlock()

            // 2. Broadcast the Mute event to the entire room instantly
			s.io.To(socket.Room(roomId)).Emit(ActionMute, ...)
		})
```

### The Teardown Flow (Server memory cleanup)
```go
		// Triggers gracefully on route change, or forcefully on internet disconnection
		leave := func() {
            // 1. Nuke them from the server-wide map
			s.mu.Lock()
			user, ok := s.users[client.Id()]
			if ok { delete(s.users, client.Id()) }
			s.mu.Unlock()

			if ok && user != nil {
                // 2. Figure out who else is still in the room
				s.mu.RLock()
				// ... lookup loop ...
				s.mu.RUnlock()

                // 3. Emits ActionRemovePeer to every remaining survivor so React can unmount the audio element
				for _, remainingClientSid := range remainingSockets {
					s.io.To(socket.Room(remainingClientSid)).Emit(ActionRemovePeer, ...)
				}
                // 4. Unsubscribe from the socket channel
				client.Leave(socket.Room(user.RoomID))
			}
		}

		client.On(ActionLeave, func(args ...any) { leave() })      // Graceful tab navigation
		client.On("disconnect", func(args ...any) { leave() }) // Forceful tab close / wifi drops
```
