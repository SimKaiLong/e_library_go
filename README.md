# e-Library API

This is a system for managing e-book loans. It is built to be reliable, easy to change, and easy to monitor.

## Overview

This project provides a way to manage a library's books and loans through a web interface (API). It is organized in a way that makes it easy to grow and maintain.

### Features

- **Organized Structure**: The code is separated into logical parts.
- **Flexible Storage**: It can easily switch between different types of data storage.
- **Two Storage Options**: It can store data in temporary memory or a permanent database.
- **Easy Settings**: Settings can be adjusted without changing the code.
- **Safe Stopping**: The system finishes its work before shutting down.
- **Activity Records**: Every action is recorded in a clear format.
- **Testing**: Includes checks to ensure everything works as expected.
- **System Monitoring**: Includes a way to check if the system and its storage are working correctly.

## Tech Stack

- **Language**: Go (1.23+)
- **Web Framework**: [Gin Gonic](https://github.com/gin-gonic/gin)
- **Logging**: [zerolog](https://github.com/rs/zerolog)
- **Settings**: [env](https://github.com/caarlos0/env) & [godotenv](https://github.com/joho/godotenv)
- **Database**: PostgreSQL (Driver: `lib/pq`)
- **Testing**: [testify](https://github.com/stretchr/testify)

## Project Structure

```text
├── cmd/api/            # Application startup logic
├── internal/
│   ├── config/         # Settings loader
│   ├── errors/         # Error definitions
│   ├── handlers/       # Web interface logic
│   ├── middleware/     # Activity tracking and recovery
│   ├── models/         # Data definitions
│   ├── repository/     # Data storage logic
│   └── service/        # Business rules
├── .env.example        # Settings template
└── README.md
```

## Setup & Installation

### Prerequisites

- Go 1.23 or higher
- PostgreSQL (optional, uses memory by default)

### Local Development

1. **Get the code**
   ```bash
   git clone <repository-url>
   cd e-library-api
   ```

2. **Prepare settings**
   Copy the example settings file:
   ```bash
   cp .env.example .env
   ```

3. **Install tools**
   ```bash
   go mod tidy
   ```

4. **Run the system**
   ```bash
   go run cmd/api/main.go
   ```
   *The system will start on port 3000.*

5. **Run checks**
   ```bash
   go test -v ./...
   ```

## PostgreSQL Local Setup

To use PostgreSQL as your storage backend:

1. **Start PostgreSQL**: Ensure your local PostgreSQL server is running.

2. **Create Database and User**:
   > **Note**: Replace `<password>` with a strong password of your choice.
   ```sql
   CREATE USER e_library_user WITH PASSWORD '<password>';
   CREATE DATABASE e_library_db OWNER e_library_user;
   ```

3. **Initialize Schema**:
   Connect to `e_library_db` and run the following commands:
   ```sql
   CREATE TABLE books (
       title TEXT PRIMARY KEY,
       available_copies INT NOT NULL CHECK (available_copies >= 0)
   );

   CREATE TABLE loans (
       borrower TEXT NOT NULL,
       title TEXT NOT NULL REFERENCES books(title),
       loan_date TIMESTAMP NOT NULL,
       return_date TIMESTAMP NOT NULL,
       PRIMARY KEY (borrower, title)
   );

   -- Seed initial data
   INSERT INTO books (title, available_copies) VALUES 
   ('The Go Programming Language', 5),
   ('Clean Code', 2),
   ('Design Patterns', 1);
   ```

4. **Update Environment Settings**:
   In your `.env` file, change the following:
   ```env
   DB_TYPE=postgres
   DATABASE_URL=host=localhost user=e_library_user password=<password> dbname=e_library_db sslmode=disable
   ```

## Configuration

The system uses environment settings. These can be placed in a `.env` file for local use.

| Variable | Description | Default |
| :--- | :--- | :--- |
| `PORT` | The port the system uses | `3000` |
| `DB_TYPE` | Where to store data (`memory` or `postgres`) | `memory` |
| `DATABASE_URL` | Database connection details | `host=localhost user=user password=<password> dbname=lib sslmode=disable` |
| `APP_ENV` | Mode (`development` or `production`) | `development` |

## How to use the API

### Look for a book
- **GET** `/Book?title={title}`
  - Shows if a book is available.
  - **Example**: `200 OK` with `{"title": "...", "available_copies": 5}`

### Borrow a book
- **POST** `/Borrow`
  - Starts a 28-day loan.
  - **Body**: `{"name_of_borrower": "Alice", "book_title": "Clean Code"}`

### Extend a loan
- **POST** `/Extend`
  - Adds 21 days to a loan.
  - **Body**: `{"name_of_borrower": "Alice", "book_title": "Clean Code"}`

### Return a book
- **POST** `/Return`
  - Ends a loan and puts the book back.
  - **Body**: `{"name_of_borrower": "Alice", "book_title": "Clean Code"}`

### Check system status
- **GET** `/health`
  - Shows if the system and its storage are working correctly.
  - **Example**: `200 OK` with `{"status": "UP"}`

## Design Principles

- **Separation of Logic**: Business rules are kept separate from how data is stored.
- **Reliable Updates**: Database changes happen safely to prevent errors.
- **Safety**: Handles multiple requests at the same time without issues.
- **Safe Exit**: Stops gracefully to avoid losing work.
- **Visibility**: Records every transaction to help with troubleshooting.

## Future Plans

- **Better Timing**: Improve how the system handles long-running tasks.
- **Automatic Database Setup**: Automate how the database is prepared.
- **Automatic Documentation**: Generate technical guides automatically.
- **Security**: Add login requirements.
