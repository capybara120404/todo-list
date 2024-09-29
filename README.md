# todo-list

A simple to-do list application that allows users to manage tasks with basic CRUD operations.

## Features

- Add new tasks
- Update tasks
- Mark tasks as completed
- Delete tasks
- View all tasks

## Technologies

- **Go**: Programming language used for the backend.
- **SQLite**: Database to store tasks.
- **Chi**: HTTP router used for handling requests.

## How to Run

Clone the repository:
   ```bash
   git clone https://github.com/capybara120404/todo-list
   cd todo-list
   go run cmd/api/main.go
   ```

## API Endpoints

- `GET /*`: Serve static files.
- `GET /api/nextdate`: Get the next scheduled date for a task.
- `GET /api/tasks`: Get all tasks.
- `GET /api/task`: Get a task by Id.
- `POST /api/task`: Add a new task.
- `POST /api/task/done`: Mark a task as completed.
- `PUT /api/task`: Update an existing task.
- `DELETE /api/task`: Delete a task.

## Go Version

This project was developed using **Go 1.23.1**. It is recommended to use this version for compatibility.
