# 💬 Real-Time Chat Backend

A real-time one-to-one chat backend built with **Go**, **Gorilla WebSocket**, **MySQL**, and **JWT Authentication**. Users can register, log in, establish a WebSocket connection, and exchange messages instantly.

## 🚀 Features

- User Registration
- User Login
- JWT Authentication
- Protected Routes
- Real-Time Messaging using WebSockets
- One-to-One Chat
- Message Persistence in MySQL
- Online User Tracking
- Docker Support

---

## 🛠 Tech Stack

- **Language:** Go
- **Router:** Gorilla Mux
- **WebSocket:** Gorilla WebSocket
- **Database:** MySQL
- **Authentication:** JWT
- **Password Hashing:** bcrypt
- **Containerization:** Docker & Docker Compose

---

## 📂 Project Structure

```
.
├── controllers/
├── database/
├── middleware/
├── models/
├── router/
├── utils/
├── ws/
├── main.go
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md
```

---

## 🔐 Authentication Flow

1. User registers.
2. Password is hashed using bcrypt.
3. User logs in.
4. Server generates a JWT token.
5. Client sends the JWT in the `Authorization` header.
6. Middleware validates the token before allowing access to protected endpoints.

---

## 💬 WebSocket Flow

1. Client authenticates using JWT.
2. WebSocket connection is established.
3. User is registered in the Hub.
4. Messages are stored in MySQL.
5. If the receiver is online, the message is delivered instantly.
6. On disconnect, the user is removed from the active connections.

---

## 📡 API Endpoints

### Authentication

| Method | Endpoint | Description |
|---------|----------|-------------|
| POST | `/register` | Register a new user |
| POST | `/login` | Login user |

### Protected

| Method | Endpoint | Description |
|---------|----------|-------------|
| GET | `/messages/{userID}` | Fetch chat history |

### WebSocket

```
ws://localhost:8080/ws
```

---

## 📨 WebSocket Message Format

### Send Message

```json
{
    "receiver_id": 2,
    "content": "Hello!"
}
```

### Receive Message

```json
{
    "sender_id": 1,
    "receiver_id": 2,
    "content": "Hello!",
    "created_at": "2026-07-24T12:00:00Z"
}
```

---

## 🗄 Database

### Users

- id
- name
- email
- password

### Messages

- id
- sender_id
- receiver_id
- content
- created_at

---

## 🐳 Running with Docker

### Clone the repository

```bash
git clone https://github.com/yourusername/chat-backend.git
cd chat-backend
```

### Build and start

```bash
docker compose up --build
```

---

## ▶️ Running Locally

Install dependencies

```bash
go mod tidy
```

Run the server

```bash
go run main.go
```

---

## 🔑 Environment Variables

Create a `.env` file.

```env
PORT=8080

MYSQL_USER=root
MYSQL_PASSWORD=your_password
MYSQL_DATABASE=chat_db
MYSQL_HOST=mysql
MYSQL_PORT=3306

SECRET=your_jwt_secret
```

---

## 📈 Future Improvements

- Group Chats
- Typing Indicators
- Read Receipts
- Message Reactions
- File Sharing
- User Presence
- Redis for Horizontal Scaling
- Push Notifications

---

## 📸 Architecture

```
Client A
     │
     │ WebSocket
     ▼
Go Server (Hub)
     │
     ├────────► Client B
     │
     ▼
   MySQL
```

---

## 👨‍💻 Author

**Anurag Jha**

- LinkedIn: https://www.linkedin.com/in/anurag-shekhar-jha-39945b27a/
