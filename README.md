# P2P Chat Application

This is a simple peer-to-peer (P2P) chat application built using Go, Gin, WebSockets, and a PostgreSQL database. It features user authentication with JWT and real-time messaging between two users.

## Features

* **User Authentication:**
    * User registration and login with JWT authentication.
    * Secure password hashing using bcrypt.
* **Real-time Messaging:**
    * WebSocket-based P2P communication.
    * Instant message delivery.
* **Message Persistence:**
    * Messages are stored in a PostgreSQL database.
* **CORS Support:**
    * CORS middleware for cross-origin requests.

## Technologies Used

* **Go:** Programming language.
* **Gin:** Web framework.
* **WebSockets:** Real-time communication.
* **PostgreSQL:** Database.
* **GORM:** ORM (Object-Relational Mapping).
* **JWT:** JSON Web Tokens for authentication.
* **bcrypt:** Password hashing.
* **godotenv:** Load environment variables from a `.env` file.

## Prerequisites

* Go (version 1.21 or later)
* PostgreSQL database
* `.env` file with database credentials and JWT secret.

## Setup

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/BishowShrestha50/chat
    cd chat
    ```

2.  **Create a `.env` file:**

    Create a `.env` file in the project root with the following content, replacing the placeholders with your actual values:

    ```
    DB_USER=your_db_user
    DB_PASSWORD=your_db_password
    DB_NAME=your_db_name
    DB_HOST=localhost
    DB_PORT=5432
    JWT_SECRET=your_jwt_secret
    ```

3.  **Install dependencies:**

    ```bash
    go mod tidy
    ```

4.  **Run the application:**

    ```bash
    go run main.go
    ```

5.  **Frontend:**

    Place the `index.html` file into the static folder. The application will serve the static files from the static folder.

## API Endpoints

* **`POST /register`:** Register a new user.
* **`POST /login`:** Login and get a JWT token.
* **`GET /chat`:** Establish a WebSocket connection. (requires token as query parameter)
* **`GET /chat/id`:** Get messages between two users. (requires JWT token in the Authorization header. and userID2 as query parameter)
* **`GET /chat/online`:** Get online status of user

## Future Improvements

* Add support for group chats.
* Improve error handling and logging.
* Implement more robust security measures.
* Add more comprehensive frontend.
* Implement better user interface for receiver selection.
* Implement file sharing.
* Implement image sharing.

## Author

Bishow Shrestha