# 🏠 CodersHouse

A **real-time audio room platform** inspired by Clubhouse — built for developers. Users can join or create audio rooms, talk in real-time using WebRTC peer-to-peer connections, and communicate via Socket.IO signalling. Authentication is phone-number-based with OTP verification.

---

## ✨ Features

- 📱 **Phone OTP Authentication** — Sign in with your phone number via a secure, hash-based OTP flow
- 🎙️ **Live Audio Rooms** — Create and join real-time audio rooms powered by WebRTC
- 🔊 **Speaker / Listener Roles** — Rooms have distinct speaker and listener roles
- 🔴 **Open & Private Rooms** — Choose room visibility when creating
- ⚡ **Real-time Signalling** — Socket.IO handles WebRTC handshakes, room join/leave events
- 🔐 **JWT Auth** — Stateless auth using access + refresh token pair stored in httpOnly cookies
- 🐳 **Fully Dockerized** — One command to spin up the entire stack

---

## 🛠️ Tech Stack

### Frontend
| Technology | Purpose |
|---|---|
| React 19 + Vite | UI framework & build tool |
| Redux Toolkit | Global state management |
| React Router v7 | Client-side routing |
| Socket.IO Client | Real-time signalling |
| Tailwind CSS v4 | Styling |
| Axios | HTTP client |
| Nginx | Static file serving in production |

### Backend
| Technology | Purpose |
|---|---|
| Go (Golang) | Server runtime |
| Chi Router | HTTP routing & middleware |
| MongoDB (v2 driver) | Database |
| Socket.IO (Go) | WebSocket / signalling server |
| Gorilla WebSocket | Low-level WebSocket support |
| JWT (golang-jwt) | Authentication tokens |
| go-playground/validator | Request validation |
| godotenv | Environment config |

### Infrastructure
| Technology | Purpose |
|---|---|
| Docker + Docker Compose | Containerisation & orchestration |
| MongoDB | Persistent data storage |

---

## 🗂️ Project Structure

```
coderhouse/
├── backend/                  # Go backend
│   ├── cmd/api/              # Application entrypoint
│   ├── internal/
│   │   ├── config/           # App configuration
│   │   ├── database/         # MongoDB connection
│   │   ├── handler/          # HTTP request handlers
│   │   ├── middleware/        # JWT auth middleware
│   │   ├── models/           # Data models (User, Room)
│   │   ├── repository/       # Database access layer
│   │   ├── socket/           # Socket.IO event handlers
│   │   └── utils/            # Helpers (JWT, encryption, JSON)
│   ├── Dockerfile
│   └── go.mod
│
├── frontend/                 # React + Vite frontend
│   ├── src/
│   │   ├── components/       # Reusable UI components
│   │   ├── pages/            # Route-level page components
│   │   ├── routes/           # Protected / guest route wrappers
│   │   ├── hooks/            # Custom React hooks
│   │   ├── store/            # Redux store & slices
│   │   ├── socket/           # Socket.IO client setup
│   │   └── http/             # Axios instance & API calls
│   ├── nginx.conf            # Nginx config (SPA fallback)
│   └── Dockerfile
│
└── docker-compose.yml        # Full-stack orchestration
```

---

## 🔌 API Endpoints

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `GET` | `/` | — | Health check |
| `POST` | `/api/send-otp` | — | Send OTP to phone number |
| `POST` | `/api/verify-otp` | — | Verify OTP & issue tokens |
| `POST` | `/api/activate` | ✅ JWT | Set display name & avatar |
| `GET` | `/api/refresh` | — | Refresh access token via cookie |
| `GET` | `/api/logout` | — | Clear auth cookies |
| `GET` | `/api/rooms` | ✅ JWT | Get all open rooms |
| `POST` | `/api/rooms` | ✅ JWT | Create a new room |
| `GET` | `/api/rooms/{roomId}` | ✅ JWT | Get a room by ID |

---

## 🚀 Getting Started

### Prerequisites

- [Docker](https://www.docker.com/) & Docker Compose
- (For local dev) [Go 1.21+](https://go.dev/) and [Node.js 18+](https://nodejs.org/)

---

### 🐳 Run with Docker (Recommended)

**1. Clone the repository**
```bash
git clone https://github.com/sanket9162/codershouse.git
cd coderhouse
```

**2. Configure backend environment**
```bash
cp backend/.env.example backend/.env
# Then edit backend/.env with your values:
```

```env
PORT=8080
DOMAIN=localhost
KEY=your-secret-key-here

# Optional: Twilio credentials for real SMS OTP
TWILIO_SID=your_twilio_sid
TWILIO_TOKEN=your_twilio_token
TWILIO_PHONE=+1234567890
```

> **Note:** Without Twilio credentials, the OTP is logged to the server console — useful for development.

**3. Start all services**
```bash
docker compose up --build
```

| Service | URL |
|---|---|
| Frontend | http://localhost:3000 |
| Backend API | http://localhost:8080 |
| MongoDB | `mongodb://localhost:27017` |

---

### 💻 Local Development

#### Backend
```bash
cd backend
cp .env.example .env
# Fill in your .env values
go mod download
go run ./cmd/api
```

#### Frontend
```bash
cd frontend
npm install
npm run dev
```
The dev server runs at **http://localhost:5173** and proxies `/api` and `/socket.io` requests to the backend.

---

## 🔒 Authentication Flow

```
1. User enters phone number
      ↓
2. POST /api/send-otp  →  OTP generated & (optionally) sent via Twilio SMS
      ↓
3. POST /api/verify-otp  →  OTP verified, JWT access + refresh tokens set as httpOnly cookies
      ↓
4. POST /api/activate  →  User sets name & avatar (one-time profile setup)
      ↓
5. Authenticated user can browse & join rooms
```

---

## 🎙️ WebRTC Flow

```
User A creates/joins room  →  Socket.IO "join-room" event
      ↓
Server notifies existing peers  →  "user-joined" event
      ↓
Peers exchange SDP Offers/Answers via Socket.IO signalling
      ↓
ICE candidates exchanged
      ↓
Direct peer-to-peer audio connection established
```

---

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📄 License

This project is open source and available under the [MIT License](LICENSE).
