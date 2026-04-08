# Comprehensive Guide to Dockerizing the Codershouse Application

Docker is a tool that packages your entire application—React, Go, Nginx, and MongoDB—into isolated "containers". A container includes your code, the exact version of Node or Go you need, and all dependencies, ensuring it runs exactly the same on your laptop as it does on a production server.

Here is a line-by-line explanation of the Docker configuration for Codershouse.

---

## 1. The Go Backend: `backend/Dockerfile`
We use a **Multi-Stage Build** to keep the final image extremely tiny. We use a heavy SDK to compile the code, then copy only the static compiled binary into a fresh, barebones Linux container.

### Stage 1: The Builder
```dockerfile
# Start from the official Go compiler image (Alpine is a super lightweight Linux distro)
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy ONLY the dependency files first. Docker caches this layer heavily!
COPY go.mod go.sum ./

# Download all the imported packages (saves time on future builds)
RUN go mod download

# Now, copy all your source code (`main.go`, `internal/`, etc.)
COPY . .

# Compile the Go code into a single executable binary named "main"
# CGO_ENABLED=0 ensures it doesn't rely on any external C-libraries (static linking)
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api
```

### Stage 2: The Production Image
```dockerfile
# Start completely fresh with a barebones, 5MB Linux image
FROM alpine:latest

WORKDIR /app

# Grab the compiled 15MB binary from the builder stage, leaving behind 500MB of Go source cache!
COPY --from=builder /app/main .

# Expose port so outside traffic can hit our server
EXPOSE 8080

# The command that actually launches the server when the container starts
CMD ["./main"]
```

---

## 2. The React Frontend: `frontend/Dockerfile`
Similar to Go, we don't need Node.js in production! Browsers only read static `.html` and `.js` files. We use Node to build the Vite bundle, and **Nginx** (an ultra-fast web server) to serve it.

### Stage 1: The Node Builder
```dockerfile
# Start with Node.js to install NPM packages
FROM node:18-alpine AS builder

WORKDIR /app
COPY package.json package-lock.json ./

# Ensure exact dependency versions are installed cleanly
RUN npm ci

COPY . .

# Run Vite to compile React JSX into raw Javascript inside the /dist folder
RUN npm run build
```

### Stage 2: Nginx Web Server
```dockerfile
# Use the official Nginx web server
FROM nginx:alpine

# Copy the raw compiled /dist folder from the Node stage into Nginx's public HTML folder
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy our custom Nginx config (crucial for React Router, see below)
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

# Launch Nginx in the foreground so the container doesn't exit immediately
CMD ["nginx", "-g", "daemon off;"]
```

---

## 3. The Nginx Configuration: `frontend/nginx.conf`
When you build a Single Page Application (SPA) like React, there is only one physical file: `index.html`. 
If a user tries to access `http://yoursite.com/rooms` directly, the web server looks for a physical folder named `/rooms`. If it doesn't find it, it throws a `404 Not Found`!

```nginx
server {
    listen 80;

    location / {
        root   /usr/share/nginx/html;
        index  index.html index.htm;
        
        # This is where the magic happens for React Router!
        # When attempting to load /rooms...
        # 1. Look for a file called /rooms (fails)
        # 2. Look for a directory called /rooms/ (fails)
        # 3. Fallback to /index.html! (success! React Router takes over the URL)
        try_files $uri $uri/ /index.html;
    }
}
```

---

## 4. The Orchestrator: `docker-compose.yml`
This file weaves the database, backend, and frontend together into an isolated private network so they can talk to each other cleanly without IP clashes.

```yaml
version: '3.8'

services:
  # ----------------------------------------
  # THE DATABASE
  # ----------------------------------------
  mongo:
    image: mongo:latest            # Pull the official MongoDB image
    container_name: coderhouse-mongo
    ports:
      - "27017:27017"              # Map your laptop's 27017 to the container's 27017
    volumes:
      - mongo-data:/data/db        # Persist data permanently. Even if container deletes, data survives!
    restart: always                # If docker crashes, auto-reboot mongo

  # ----------------------------------------
  # YOUR GO BACKEND
  # ----------------------------------------
  backend:
    build: ./backend               # Execute the backend/Dockerfile
    container_name: coderhouse-backend
    ports:
      - "8080:8080"
    depends_on:
      - mongo                      # Guarantee MongoDB boots *before* Go boots
    environment:
      # Automatically tell Go to connect to the "mongo" container over the private Docker network!
      - DB_URL=mongodb://mongo:27017
      - DB_NAME=coderhouse
    restart: always

  # ----------------------------------------
  # YOUR REACT FRONTEND
  # ----------------------------------------
  frontend:
    build: ./frontend              # Execute the frontend/Dockerfile
    container_name: coderhouse-frontend
    ports:
      - "3000:80"                  # Map your laptop's port 3000 to the container's Nginx port 80
    depends_on:
      - backend                    # Guarantee Backend boots *before* Frontend
    restart: always

# Define the persistent volume storage for Mongo to use securely
volumes:
  mongo-data:
```
