# 🏠 CodersHouse

A **real-time audio room platform** inspired by Clubhouse — built for developers. Users can join or create audio rooms, talk in real-time using WebRTC peer-to-peer connections, and communicate via Socket.IO signalling. Authentication is phone-number-based with OTP verification.

---

## ✨ Features

- 📱 **Phone OTP Authentication** — Sign in with your phone number via a secure, hash-based OTP flow
- 🎙️ **Live Audio Rooms** — Create and join real-time audio rooms powered by WebRTC
- 🔊 **Speaker / Listener Roles** — Rooms have distinct speaker and listener roles
- ⚡ **Real-time Signalling** — Socket.IO handles WebRTC handshakes, room join/leave events
- 🔐 **JWT Auth** — Stateless auth using access + refresh token pair stored in httpOnly cookies
- 🐳 **Fully Dockerized** — One command to spin up the entire stack

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
- (For local dev) [Go 1.25+](https://go.dev/) and [Node.js 24+](https://nodejs.org/)

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

DB_URL=mongodb://mongo:27017
DB_NAME=codershouse

#jwt env
JWT_SECRET=your-secret-key-here
JWT_ISSUER=your-issuer-here
JWT_AUDIENCE=your-audience-here
COOKIE_DOMAIN=localhost

```

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

---

## 📄 License

This project is open source and available under the [MIT License](LICENSE).
