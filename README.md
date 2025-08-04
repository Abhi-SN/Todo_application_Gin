Full-Stack Todo List Application with Go and Gin
This project is a complete full-stack CRUD (Create, Read, Update, Delete) application. The backend is a RESTful API built with the Gin web framework for Go, and the frontend is a simple, dynamic single-page application.

The application uses a PostgreSQL database for persistent data storage, and all services are containerized with Docker.

Features
Full CRUD Functionality: Create, read, update, and delete todo items.

RESTful API: A clean, predictable API built with Gin.

Database Persistence: Data is stored in a PostgreSQL database and persists across application restarts.

Full-Stack: Includes both a backend API and a frontend UI.

Containerized: The entire application (API, database, frontend) is managed with Docker Compose for easy setup.

Directory Structure
gin-todo-api/
│
├── 📂 frontend/
│   ├── index.html
│   └── Dockerfile
│
├── 📜 main.go
├── 📜 Dockerfile
└── 📜 docker-compose.yml

Prerequisites
Go: Version 1.22 or later.

Docker & Docker Compose: For building and running the services.

How to Run the Application
Clone the repository and navigate to the project's root directory (gin-todo-api).

Install Go Dependencies: Run the following command in the root directory to download the necessary Go packages (Gin, CORS middleware, database driver).

go mod tidy

Build and Start Services: Use Docker Compose to build and run all the containers.

docker-compose up --build

This will start the PostgreSQL database, the Gin API server, and the Nginx frontend server.

Access the Application: Open your web browser and navigate to:
http://localhost:3000

How to Use
The application is straightforward to use via the web interface.

Add a Todo: Type a task into the input field and click "Add".

Mark as Complete: Click on the text of a todo to toggle its completed status (it will be crossed out).

Edit a Todo: Click the "Edit" button, change the text, and click "Save".

Delete a Todo: Click the "×" button to remove a todo from the list.

All changes are saved to the database and will be there if you refresh the page or restart the application.



