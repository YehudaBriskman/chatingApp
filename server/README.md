# 📌 Chat Application - README

## 📖 Introduction
This project is a **real-time chat application** built using **Go, Gin, PostgreSQL, and WebSockets**. The application supports authentication, role-based access, and real-time messaging between users in different rooms.

## 🚀 Features
- User authentication with JWT
- Role-based access control (**User, Admin, Super Admin**)
- Real-time WebSocket communication
- Room creation and management
- System logs for admin tracking
- RESTful API endpoints with **Gin framework**
- PostgreSQL as the database

## 🔧 Prerequisites
Before running the application, ensure you have the following installed:

- **Go 1.20+**
- **PostgreSQL 13+**
- **Git**
- **Make** (optional but recommended for easier setup)
- **Docker** (if using Docker for PostgreSQL)

## 📥 Installation & Setup
### 1️⃣ Clone the repository
```sh
git clone https://github.com/YehudaBriskman/chatingApp.git
cd chatingApp/server
```

### 2️⃣ Set up environment variables
Create a **`.env`** file in the root directory and add the following:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=Yehuda@2004xyz
DB_NAME=mydb
DB_SSLMODE=disable
SUPERADMIN_EMAIL=yr0556772363@gmail.com
SUPERADMIN_PASSWORD=Yehuda@2004xyz
MODE=dev
```

### 3️⃣ Install dependencies
```sh
go mod tidy
```

### 4️⃣ Start PostgreSQL Database
**Option 1:** Run PostgreSQL locally
```sh
sudo systemctl start postgresql
```

**Option 2:** Use Docker (recommended)
```sh
docker-compose up -d
```

### 5️⃣ Run database migrations
```sh
go run migrations/migrate.go
```

### 6️⃣ Start the application
```sh
go run main.go
```

The server will start on **http://localhost:8080** 🚀

---

## 🎯 API Endpoints
### 🔑 Authentication
| Method | Endpoint       | Description |
|--------|---------------|-------------|
| POST   | `/auth/login` | User login  |
| POST   | `/auth/signup` | Register new user |

### 📌 Rooms Management
| Method | Endpoint       | Description |
|--------|---------------|-------------|
| GET    | `/rooms/`     | Get all rooms |
| POST   | `/rooms/`     | Create a new room (Admin only) |
| DELETE | `/rooms/:id`  | Delete a room (Admin only) |

### 💬 Messages
| Method | Endpoint       | Description |
|--------|---------------|-------------|
| GET    | `/messages/:room_id` | Get all messages in a room |
| POST   | `/messages/` | Send a new message |

---

## 🛠 Project Structure
```
chat-app/
├── handlers/        # HTTP request handlers
├── middleware/      # Authentication and validation middleware
├── models/          # Database models
├── repository/      # Database queries
├── services/        # Business logic
├── routes/          # API route definitions
├── migrations/      # Database migrations
├── main.go          # Application entry point
└── .env             # Environment variables
```

---

## 📡 WebSocket Support
The application supports real-time communication using WebSockets.
### WebSocket Connection Example
```javascript
const socket = new WebSocket("ws://localhost:8080/ws");
socket.onopen = () => console.log("Connected to WebSocket");
socket.onmessage = (event) => console.log("Message received: ", event.data);
socket.send(JSON.stringify({ action: "send_message", roomID: 1, content: "Hello!" }));
```

---

## ✅ Testing API Requests
For testing purposes, use **Postman** or **cURL**:
```sh
curl -X POST http://localhost:8080/auth/login -d '{"email": "yr0556772363@gmail.com", "password": "Yehuda@2004xyz"}' -H "Content-Type: application/json"
```

---

## 🛑 Stop the application
If running locally: 
```sh
CTRL + C
```

If running via Docker:
```sh
docker-compose down
```

---

## 👨‍💻 Contributors
- **Yehuda** - Full-stack developer 🔥

📧 Contact: [yr0556772363@gmail.com](mailto:yr0556772363@gmail.com)

---

## 🌟 Support
If you find this project useful, consider giving it a ⭐ on GitHub! 🚀

