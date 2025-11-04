# Go Student API

![Go Version](https://img.shields.io/badge/Go-1.19+-00ADD8?style=flat&logo=go)
![Last Updated](https://img.shields.io/badge/Last%20Updated-2025--11--04-brightgreen)
![License](https://img.shields.io/badge/License-MIT-blue.svg)
A Go API project for managing student data with a clean architecture approach.

A RESTful API service for managing student data built with Go, following clean architecture principles.

## Project Overview

This API service provides the following functionality:
- Student CRUD operations
- Course enrollment/management
- Student data validation
- Database persistence (PostgreSQL/SQLite support)
- RESTful API endpoints

## Project Structure

```
go-student-api/
├── cmd/
│   └── student-api/
│       └── main.go          # Application entry point
├── configs/
│   ├── config.dev.yaml      # Development configuration
│   └── config.prod.yaml     # Production configuration
├── internal/
│   ├── config/              # Configuration management
│   │   └── config.go        # Config structs and loading logic
│   ├── http/
│   │   └── handlers/        # HTTP request handlers
│   │       └── student/     # Student-related handlers
│   ├── storage/             # Data access layer
│   │   ├── factory/         # Storage implementation factory
│   │   ├── sqlite/          # SQLite implementation
│   │   └── postgres/        # PostgreSQL implementation
│   ├── types/               # Domain types/models
│   │   └── types.go         # Student and Course types
│   └── utils/               # Utility packages
│       └── response/        # HTTP response helpers
├── migrations/              # Database migration files
│   ├── 000001_init.up.sql
│   └── 000001_init.down.sql
├── pkg/                     # Public libraries
│   └── validator/           # Shared validation utilities
├── scripts/                 # Utility scripts
│   ├── migrations.sh        # Database migration helper
│   └── setup.sh             # Development setup script
├── .env.example             # Environment variables template
├── docker-compose.yml       # Docker compose configuration
├── Dockerfile               # Docker build configuration
├── go.mod                   # Go modules file
└── README.md                # Project documentation
```

## Prerequisites

- Go 1.19 or higher
- PostgreSQL or SQLite
- Docker (optional)
- Make (optional)

## Getting Started

1. Clone the repository
```bash
git clone https://github.com/manish-npx/go-student-api.git
cd go-student-api
```

2. Set up environment variables
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. Start the database (using Docker)
```bash
docker-compose up -d postgres
```

4. Run database migrations
```bash
./scripts/migrations.sh up
```

5. Build and run the application
```bash
go build -o bin/api ./cmd/student-api
./bin/api
```

## Configuration

Create a `.env`  or `yaml` file with the following variables:

```env
# Server
PORT=8080
ENV=development
LOG_LEVEL=debug

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=postgres
DB_NAME=student_api
DB_SSLMODE=disable

# Connection Pool
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=25
DB_CONN_MAX_LIFETIME=300s
```

## API Endpoints

### Students
- `GET /api/students` - List all students
- `GET /api/student/{id}` - Get a specific student
- `POST /api/student` - Create a new student
- `PUT /api/student/{id}` - Update a student
- `DELETE /api/student/{id}` - Delete a student

### Courses
- `GET /api/courses` - List all courses
- `POST /api/courses` - Create a new course
- `POST /api/students/{id}/enroll` - Enroll student in courses

## Database Schema

### Students Table
```sql
CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    age INTEGER NOT NULL
);
```

### Courses Table
```sql
CREATE TABLE courses (
    id SERIAL PRIMARY KEY,
    code VARCHAR(10) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT
);
```

### Enrollments Table
```sql
CREATE TABLE enrollments (
    student_id INTEGER REFERENCES students(id),
    course_id INTEGER REFERENCES courses(id),
    PRIMARY KEY (student_id, course_id)
);
```

## Development

### Running with Docker
```bash
docker build -t student-api .
docker run -p 8080:8080 --env-file .env student-api
```

### Running Tests
```bash
go test ./...
```

### Database Migrations
```bash
# Apply migrations
./scripts/migrations.sh up

# Rollback migrations
./scripts/migrations.sh down
```

## Project Layout Explanation

- `cmd/`: Contains the main application entry points
- `configs/`: Configuration files for different environments
- `internal/`: Private application code
  - `config/`: Configuration management
  - `http/handlers/`: HTTP request handlers
  - `storage/`: Data access layer implementations
  - `types/`: Domain models
  - `utils/`: Utility packages
- `migrations/`: Database migration files
- `pkg/`: Public libraries that can be used by external projects
- `scripts/`: Utility scripts for development and deployment

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contact

Manish - [@manish-npx](https://github.com/manish-npx)

Project Link: [https://github.com/manish-npx/go-student-api](https://github.com/manish-npx/go-student-api)