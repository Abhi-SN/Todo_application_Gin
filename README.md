# Full-Stack Todo List Application with Go and Gin

This repository contains a complete, full-stack CRUD (Create, Read, Update, Delete) Todo List application. The backend is a RESTful API built with the Gin framework in Go, and the frontend is a single-page application. All components are containerized using Docker, with PostgreSQL providing persistent data storage.

## ✨ Features

- **Full CRUD Functionality:** Add, read, update, and delete todos.
- **RESTful API:** Consistent, clean endpoints powered by Gin (Go).
- **Persistent Storage:** PostgreSQL ensures data durability across restarts.
- **Full-Stack UI:** Intuitive web interface for easy todo management.
- **Containerized:** One-command startup using Docker Compose.

## 🗂 Directory Structure

```
gin-todo-api/
│
├── frontend/
│   ├── index.html
│   └── Dockerfile
│
├── main.go
├── Dockerfile
└── docker-compose.yml
```

## 🚦 Prerequisites

- **Go** (version 1.22 or later)
- **Docker & Docker Compose**

## 🚀 How to Run the Application

1. **Clone the repository**  
   ```
   git clone 
   cd gin-todo-api
   ```

2. **Install Go dependencies**  
   ```
   go mod tidy
   ```

3. **Build and launch all services**  
   ```
   docker-compose up --build
   ```
   This command starts the PostgreSQL database, the Gin API server, and the Nginx-based frontend.

4. **Open your browser:**  
   ```
   http://localhost:3000
   ```

## 🖥 How to Use

- **Add a Todo:** Type a task and click "Add".
- **Mark as Complete:** Click a todo’s text to toggle its completed status (crosses out when finished).
- **Edit a Todo:** Click the "Edit" button, update the text, then press "Save".
- **Delete a Todo:** Click the "×" icon to remove an item.

All your actions are synchronized with the database—refreshing the page or restarting the app will retain all completed, pending, or deleted todos.
