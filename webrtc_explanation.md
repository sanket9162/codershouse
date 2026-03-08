# Complete Guide to WebRTC and Socket.IO in Codershouse

This document breaks down exactly how WebRTC (for voice audio) and Socket.IO (for signaling) work together in your application between the React frontend and Go backend.

---

## 1. The Core Concept: Why do we need both?
- **WebRTC (Web Real-Time Communication)** is a peer-to-peer technology. This means that audio travels *directly* from User A's browser to User B's browser without going through your backend server, saving massive amounts of bandwidth and ensuring very low latency.
- However, WebRTC has a major problem: **Browsers don't know how to find each other on the internet.** 
- To solve this, we use a **Signaling Server** (our Go Backend using **Socket.IO**). Its only job is to act as a middleman so browsers can swap their IP addresses (ICE Candidates) and multimedia capabilities (Session Description Protocol / SDP). Once the swap is done, the browsers talk to each other directly, and the signaling server is no longer involved in the audio transfer.

---

## 2. The Socket.IO Actions (`action.jsx`)

Here is the lifecycle of our websocket events in alphabetical order:

### `JOIN` (Frontend -> Backend)
Emitted by the frontend as soon as `navigator.mediaDevices.getUserMedia` successfully activates the microphone. It tells the backend *"I am ready to join the room, please introduce me to everyone else!"*

### `ADD_PEER` (Backend -> Frontend)
The backend responds by finding everyone in the requested room and triggering this event exactly twice for every connection:
- It tells an **existing user**: "A new user joined, get ready to receive a call!" (`createOffer: false`)
- It tells the **new user**: "Here is an existing user in the room, ring them!" (`createOffer: true`)

### `RELAY_SDP` and `SESSION_DESCRIPTION`
- **SDP (Session Description Protocol)** contains data like "I am transmitting audio using Opus codec".
- **Offer / Answer:** The new user will generate an SDP "Offer" and send it to the backend via `RELAY_SDP`. 
- The backend forwards it directly to the target peer using the `SESSION_DESCRIPTION` event.
- The target peer generates an SDP "Answer" and sends it back through the same pipeline.

### `RELAY_ICE` and `ICE_CANDIDATE`
- **ICE (Interactive Connectivity Establishment):** This is essentially how browsers safely bypass firewalls and routers to find each other's public IP address.
- Behind the scenes, the browser generates dozens of "Candidates" (potential network paths). The React app fires `RELAY_ICE` to the Go backend.
- The Go backend forwards these candidates instantly to the target peer using the `ICE_CANDIDATE` event.

### `LEAVE` and `REMOVE_PEER`
When you close the tab, press the back button, or explicitly leave the room, your browser triggers `LEAVE` (or standard `disconnect`). The Go backend cleans you out of its memory, and tells all remaining clients `REMOVE_PEER`, preventing ghost connections and immediately deleting your audio element from the UI.

---

## 3. The Go Backend Architecture (`socket.go`)

Because Go is highly concurrent, your backend architecture takes steps to be thread-safe:
- **`sync.RWMutex`:** Whenever the backend stores an incoming user in the `s.users` memory map (which tracks the socket ID, room ID, and frontend user object), it locks the map. This prevents the server from crashing if 1,000 users join the same room on the exact same millisecond.
- **`io.In(socket.Room(roomId))`:** `github.com/zishang520/socket.io` natively handles grouping WebSockets into logical rooms. We use this to effortlessly broadcast exclusively to the people who care about an event.

---

## 4. The React Frontend Architecture (`useWebRTC.jsx`)

The logic here handles joining the room and dynamically rendering dynamic `<audio>` elements:
1. React's `useEffect` mounts exactly once.
2. It requests microphone permissions via `navigator.mediaDevices.getUserMedia`.
3. It updates the `clients` array to include your own profile, rendering an `<audio>` tag that is forcefully muted (`volume = 0`) so you don't hear an echo of yourself.
4. It calls `socket.current.emit(ACTIONS.JOIN)`.

### Upcoming Implementation (Phase 2):
To complete the voice feature, `useWebRTC.jsx` will need listeners for all the signaling triggers. Every time `ADD_PEER` is received, it will instantiate a new `new RTCPeerConnection(freeStunServers)` instance, attach your microphone tracks to it, and map it into a generic hashmap `connections.current = {}`. This ensures your browser maintains a unique, direct pipeline to *every single* person inside the room simultaneously!
